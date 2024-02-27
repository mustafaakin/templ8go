package templ8go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestParseTemplateStringSuccess tests successful template parsing.
func TestParseTemplateStringSuccess(t *testing.T) {
	var tests = []struct {
		name     string
		vars     map[string]interface{}
		input    string
		expected string
	}{
		{
			"Simple substitution",
			map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Mustafa",
					"age":  32,
				},
			},
			"{{ user.name }} is {{ user.age }} years old.",
			"Mustafa is 32 years old.",
		},
		{
			"Arithmetic operation",
			map[string]interface{}{
				"user": map[string]interface{}{
					"age": 32,
				},
			},
			"Next year, you will be {{ user.age + 1 }}.",
			"Next year, you will be 33.",
		},
		{
			"Nested object access",
			map[string]interface{}{
				"user": map[string]interface{}{
					"profile": map[string]interface{}{
						"nickname": "Moose",
					},
				},
			},
			"Your nickname is {{ user.profile.nickname }}.",
			"Your nickname is Moose.",
		},
		{
			"Array access",
			map[string]interface{}{
				"favorites": []interface{}{"Pizza", "Ice Cream"},
			},
			"I love {{ favorites[0] }} and {{ favorites[1] }}.",
			"I love Pizza and Ice Cream.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := ParseTemplateString(test.vars, test.input)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, output)
		})
	}
}

// TestParseTemplateStringError tests error handling.
func TestParseTemplateStringError(t *testing.T) {
	vars := map[string]interface{}{
		"user": map[string]interface{}{},
	}

	// Unmatched expression delimiter
	input := "{{ user.name is unmatched."

	_, err := ParseTemplateString(vars, input)
	assert.Error(t, err, "Expected an error for unmatched expression delimiter, got none")
}
