package templ8go

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestResolveJSExpression(t *testing.T) {
	tests := []struct {
		name       string
		bindings   map[string]interface{}
		expression string
		expected   interface{}
		expectErr  bool
	}{
		{
			name:       "Basic",
			bindings:   map[string]interface{}{},
			expression: "2 + 2",
			expected:   4.0,
		},
		{
			name: "WithBindings",
			bindings: map[string]interface{}{
				"a": 10,
				"b": 5,
			},
			expression: "a * b",
			expected:   50.0,
		},
		{
			name: "Complex",
			bindings: map[string]interface{}{
				"obj": map[string]interface{}{
					"value": 5,
				},
			},
			expression: "obj.value + 15",
			expected:   20.0,
		},
		{
			name:       "ErrorHandling",
			bindings:   map[string]interface{}{},
			expression: "undeclaredVariable + 1",
			expectErr:  true,
		},
		{
			name:       "Timeout",
			bindings:   map[string]interface{}{},
			expression: "while(true){}",
			expectErr:  true,
		},
		{
			name: "JSFunction",
			bindings: map[string]interface{}{
				"x": 10,
				"y": 20,
			},
			expression: "Math.max(x, y)",
			expected:   20.0,
		},
		{
			name: "EdgeCases",
			bindings: map[string]interface{}{
				"emptyString": "",
				"nullValue":   nil,
			},
			expression: "emptyString === '' && nullValue === null",
			expected:   true,
		},
		{
			name: "TypeCoercion",
			bindings: map[string]interface{}{
				"stringValue": "10",
				"numValue":    10,
			},
			expression: "stringValue == numValue",
			expected:   true,
		},
		{
			name:       "SecurityConcerns",
			bindings:   map[string]interface{}{},
			expression: "this.constructor.constructor('return process')().exit()",
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
}

func TestResolveJSExpressionBindingsUnchanged(t *testing.T) {
	initialBindings := map[string]interface{}{
		"x": 5,
		"y": 10,
	}
	// Create a deep copy of the initialBindings for comparison after function execution
	expectedBindings := make(map[string]interface{})
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
	assert.True(t, reflect.DeepEqual(initialBindings, expectedBindings), "Bindings were modified after second execution")
}
