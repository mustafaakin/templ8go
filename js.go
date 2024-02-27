package templ8go

import (
	"encoding/json"
	"fmt"
	v8 "rogchap.com/v8go"
	"time"
)

func ResolveJSExpression(bindings map[string]interface{}, expression string) (interface{}, error) {
	ctx := v8.NewContext()
	defer ctx.Close()

	for key, val := range bindings {
		j, err := json.Marshal(val)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		// Directly use the marshaled JSON string instead of converting back to string.
		if err := ctx.Global().Set(key, string(j)); err != nil {
			return nil, fmt.Errorf("failed to set global property %s: %w", key, err)
		}

		script := fmt.Sprintf("%s = JSON.parse(%s)", key, key)
		if _, err := ctx.RunScript(script, ""); err != nil {
			return nil, fmt.Errorf("failed to parse global property %s: %w", key, err)
		}
	}

	resultChan := make(chan interface{}, 1)
	errorChan := make(chan error, 1)

	go func() {
		script := fmt.Sprintf("JSON.stringify(%s)", expression)
		val, err := ctx.RunScript(script, "")
		if err != nil {
			errorChan <- fmt.Errorf("failed to evaluate expression: %w", err)
			return
		}

		var result interface{}
		if err := json.Unmarshal([]byte(val.String()), &result); err != nil {
			errorChan <- fmt.Errorf("failed to unmarshal result: %w", err)
			return
		}

		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	case <-time.After(100 * time.Millisecond): // TODO: configurable
		ctx.Isolate().TerminateExecution()
		return nil, fmt.Errorf("execution timeout")
	}
}
