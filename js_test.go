package templ8go

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResolveJSExpression(t *testing.T) {
	tests := []struct {
		name       string
		bindings   map[string]any
		expression string
		expected   any
		expectErr  bool
	}{
		{
			name:       "basic math, addition expression",
			bindings:   nil,
			expression: "2 + 2",
			expected:   4.0,
		},
		{
			name: "basic math, multiplication expression with bindings",
			bindings: map[string]any{
				"a": 10,
				"b": 5,
			},
			expression: "a * b",
			expected:   50.0,
		},
		{
			name: "math, addition with complex object",
			bindings: map[string]any{
				"obj": map[string]any{
					"value": 5,
				},
			},
			expression: "obj.value + 15",
			expected:   20.0,
		},
		{
			name:       "non existing declaration example",
			bindings:   nil,
			expression: "undeclaredVariable + 1",
			expectErr:  true,
		},
		{
			name:       "syntax error",
			bindings:   nil,
			expression: "while(true){}",
			expectErr:  true,
		},
		{
			name: "js Math function example",
			bindings: map[string]any{
				"x": 10,
				"y": 20,
			},
			expression: "Math.max(x, y)",
			expected:   20.0,
		},
		{
			name: "js value compare",
			bindings: map[string]any{
				"emptyString": "",
				"nullValue":   nil,
			},
			expression: "emptyString === '' && nullValue === null",
			expected:   true,
		},
		{
			name: "type coercion",
			bindings: map[string]any{
				"stringValue": "10",
				"numValue":    10,
			},
			expression: "stringValue == numValue",
			expected:   true,
		},
		{
			name:       "js security concerns",
			bindings:   nil,
			expression: "this.constructor.constructor('return process')().exit()",
			expectErr:  true,
		},
		{
			name: "should raise json.Marshal error",
			bindings: map[string]any{
				"key": make(chan struct{}),
			},
			expression: "",
			expectErr:  true,
		},
		{
			name: "should raise Set error",
			bindings: map[string]any{
				"key": "while(true){}",
			},
			expression: "",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveJSExpression(tt.bindings, tt.expression)
			if tt.expectErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error occurred")
				assert.Equal(t, tt.expected, result, "The result does not match the expected value")
			}
		})
	}

	t.Run("timeout error", func(t *testing.T) {
		originalTimeout := defaultExecutionTimeout
		SetDefaultExecutionTimeout(10 * time.Microsecond)
		defer func() { SetDefaultExecutionTimeout(originalTimeout) }()

		_, err := ResolveJSExpression(map[string]any{"a": 1}, "a + 2")
		assert.ErrorIs(t, err, ErrResolveJSExpressionExecutionTimeout)
	})
}

func TestResolveJSExpressionBindingsUnchanged(t *testing.T) {
	initialBindings := map[string]any{
		"x": 5,
		"y": 10,
	}
	// Create a deep copy of the initialBindings for comparison after function execution
	expectedBindings := make(map[string]any)
	for k, v := range initialBindings {
		expectedBindings[k] = v
	}

	expression := "x + y"
	_, err := ResolveJSExpression(initialBindings, expression)
	assert.NoError(t, err, "Unexpected error occurred during first execution")

	// Verify bindings have not changed after the first execution
	assert.True(t, reflect.DeepEqual(initialBindings, expectedBindings), "Bindings were modified after first execution")

	// Execute again to ensure bindings remain unchanged even after multiple executions
	_, err = ResolveJSExpression(initialBindings, expression)
	assert.NoError(t, err, "Unexpected error occurred during second execution")

	// Verify bindings have not changed after the second execution
	assert.True(
		t,
		reflect.DeepEqual(initialBindings, expectedBindings),
		"Bindings were modified after second execution",
	)
}
