package templ8go

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// sentinel errors.
var (
	ErrResolveTemplateUnmatchedExpressionDelimeter = errors.New("unmatched expression delimiter")
)

// ResolveTemplate handles the resolve operation with a given template string and binding data as args.
func ResolveTemplate(args map[string]any, input string) (string, error) {
	var result strings.Builder
	leftDelimiter := "{{"
	rightDelimiter := "}}"

	start := 0
	for start < len(input) {
		leftIndex := strings.Index(input[start:], leftDelimiter)
		if leftIndex == -1 {
			// No more expressions to process
			result.WriteString(input[start:])
			break
		}
		leftIndex += start
		rightIndex := strings.Index(input[leftIndex:], rightDelimiter)
		if rightIndex == -1 {
			return "", ErrResolveTemplateUnmatchedExpressionDelimeter
		}
		rightIndex += leftIndex + len(rightDelimiter) - 1

		// Write the text before the expression
		result.WriteString(input[start:leftIndex])

		// Extract and evaluate the expression
		expression := input[leftIndex+len(leftDelimiter) : rightIndex-len(rightDelimiter)+1]
		evaluated, err := ResolveJSExpression(args, expression)
		if err != nil {
			return "", err
		}

		var out string
		switch v := evaluated.(type) {
		case int:
			out = fmt.Sprintf("%d", v)
		case float64:
			out = strconv.FormatFloat(v, 'f', -1, 64)
		case string:
			out = v
		default:
			out1, _ := json.Marshal(v)
			out = string(out1)
		}
		result.WriteString(out)

		// Update start position
		start = rightIndex + 1
	}

	return result.String(), nil
}
