package templ8go

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	v8 "rogchap.com/v8go"
)

var defaultExecutionTimeout = 100 * time.Millisecond

// sentinel errors.
var (
	ErrResolveJSExpressionExecutionTimeout = errors.New("execution timeout error")
)

// SetDefaultExecutionTimeout overrides default execution timeout.
func SetDefaultExecutionTimeout(d time.Duration) {
	defaultExecutionTimeout = d
}

// ResolveJSExpression handles the resolve operation with a given JS expression and binding data.
func ResolveJSExpression(bindings map[string]any, expression string) (any, error) {
	ctx := v8.NewContext()
	defer ctx.Close()

	for key, val := range bindings {
		j, err := json.Marshal(val)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		if err := ctx.Global().Set(key, string(j)); err != nil {
			return nil, fmt.Errorf("failed to set global property %s: %w", key, err)
		}

		if _, err := ctx.RunScript(key+" = JSON.parse("+key+")", ""); err != nil {
			return nil, fmt.Errorf("failed to run script, key: %s: %w", key, err) // JSError
		}
	}

	resultChan := make(chan any, 1)
	errorChan := make(chan error, 1)

	go func() {
		val, err := ctx.RunScript("JSON.stringify("+expression+")", "")
		if err != nil {
			errorChan <- fmt.Errorf("%w, expression was: %s", err, expression)
			return
		}

		var result any
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
	case <-time.After(defaultExecutionTimeout):
		ctx.Isolate().TerminateExecution()
		return nil, ErrResolveJSExpressionExecutionTimeout
	}
}
