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
			output, err := ResolveTemplate(test.vars, test.input)

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

	_, err := ResolveTemplate(vars, input)
	assert.Error(t, err, "Expected an error for unmatched expression delimiter, got none")
}

func TestReadmeExamples(t *testing.T) {
	tests := []struct {
		name     string
		template string
		bindings map[string]interface{}
		want     string
	}{
		{
			name:     "Simple Arithmetic",
			template: "The sum of 5 and 3 is {{ 5 + 3 }}.",
			bindings: map[string]interface{}{},
			want:     "The sum of 5 and 3 is 8.",
		},
		{
			name:     "Conditional Greetings",
			template: "Good {{ hour < 12 ? 'morning' : 'afternoon' }}, {{ user.name }}!",
			bindings: map[string]interface{}{
				"hour": 9,
				"user": map[string]interface{}{"name": "Alice"},
			},
			want: "Good morning, Alice!",
		},
		{
			name:     "Array Operations",
			template: "Users list: {{ users.map(user => user.name).join(', ') }}",
			bindings: map[string]interface{}{
				"users": []map[string]interface{}{
					{"name": "Alice"},
					{"name": "Bob"},
					{"name": "Charlie"},
				},
			},
			want: "Users list: Alice, Bob, Charlie",
		},
		{
			name:     "Object Manipulation",
			template: "{{ user.firstName }} {{ user.lastName }} is {{ user.age }} years old.",
			bindings: map[string]interface{}{
				"user": map[string]interface{}{"firstName": "John", "lastName": "Doe", "age": 28},
			},
			want: "John Doe is 28 years old.",
		},
		{
			name:     "Logical Operations",
			template: "You are {{ age >= 18 ? 'an adult' : 'a minor' }}.",
			bindings: map[string]interface{}{
				"age": 20,
			},
			want: "You are an adult.",
		},
		{
			name:     "String Concatenation",
			template: "{{ 'Hello, ' + user.name + '!'}}",
			bindings: map[string]interface{}{
				"user": map[string]interface{}{"name": "Jane"},
			},
			want: "Hello, Jane!",
		},
		{
			name:     "Using JavaScript Functions",
			template: "Your score is {{ Math.min(score, 100) }}.",
			bindings: map[string]interface{}{
				"score": 105,
			},
			want: "Your score is 100.",
		},
		{
			name:     "Nested Object Access",
			template: "Project {{ project.details.name }} is due on {{ project.details.dueDate }}.",
			bindings: map[string]interface{}{
				"project": map[string]interface{}{
					"details": map[string]interface{}{"name": "Apollo", "dueDate": "2024-03-01"},
				},
			},
			want: "Project Apollo is due on 2024-03-01.",
		},
		{
			name:     "Complex Expressions",
			template: "{{ user.isActive ? user.name + ' is active and has ' + user.roles.length + ' roles' : user.name + ' is not active' }}.",
			bindings: map[string]interface{}{
				"user": map[string]interface{}{"name": "Eve", "isActive": true, "roles": []string{"admin", "editor"}},
			},
			want: "Eve is active and has 2 roles.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveTemplate(tt.bindings, tt.template)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
