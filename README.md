[![build and test](https://github.com/mustafaakin/templ8go/actions/workflows/go.yml/badge.svg)](https://github.com/mustafaakin/templ8go/actions/workflows/go.yml)
[![golang ci Lint](https://github.com/mustafaakin/templ8go/actions/workflows/go-lint.yml/badge.svg)](https://github.com/mustafaakin/templ8go/actions/workflows/go-lint.yml)

# templ8go

`templ8go` is a small library designed to resolve template strings like `Hello
{{ user.name }}` utilizing the V8 JavaScript engine. The play on words in the
name stems from blending "template" with "V8" (the JavaScript engine), and "8"
phonetically resembling "ate". Sorry for the dad joke, ChatGPT come up with it.

---

## Why?

While traditional templating engines are great for generating large blocks of
text and support a variety of built-in functions and pipe syntax, `templ8go`
takes a different approach. It leverages the power and flexibility of
JavaScript for manipulating and interpolating short strings. This allows for
more dynamic and programmable template resolution without leaving the comfort
of Go.

Consider the following examples where `templ8go` shines:

- `Hello {{ user.name }}` with `{user:{name: Mustafa}}` resolves to `Hello Mustafa`.
- `Your balance is {{ account.balance.toFixed(2) }}` with `{account:{balance: 1234.567}}` becomes `Your balance is 1234.57`.

These examples showcase the simplicity and power of using JavaScript
expressions within templates.

---

## Features

- **Dynamic Expression Evaluation**: Use JavaScript expressions right within
  your template strings.
- **Bindings Support**: Seamlessly pass Go variables into the JavaScript
  execution context to be used within your template expressions.
- **Easy Integration**: Designed to be easily integrated into any Go project
  that needs flexible string interpolation.
- **Security**: It leverages V8, the same Javascript engine that runs
  Cloudflare workers and Chrome. Though we can add even more hardening.

---

## Getting Started

### Installation

First, ensure you have Go installed on your machine (version 1.21 or newer is
recommended). Then, install `templ8go` using `go get`:

```sh
go get -u github.com/mustafaakin/templ8go
```

### Usage

Here's a quick example to get you started:

```go
package main

import (
    "fmt"
    "github.com/mustafaakin/templ8go"
)

func main() {
    template := "Hello {{ user.name }}"
    bindings := map[string]interface{}{
        "user": map[string]interface{}{
            "name": "Mustafa",
        },
    }

    result, err := templ8go.ResolveTemplate(bindings, template)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println(result) // Output: Hello Mustafa
}
```

Or you can use the Javascript resolve directly to get the result of object without string interpolation which 
can be useful in some environments.

```go
package main

import (
    "fmt"
    "github.com/mustafaakin/templ8go"
)

func main() {
    bindings := map[string]interface{}{
        "user": map[string]interface{}{
            "name": "Mustafa",
        },
    }
	
    result, err := templ8go.ResolveJSExpression(bindings, "user.name")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println(result) // Output: Mustafa
}
```

You can change the default execution timeout via calling;

```go
// now you have 200 Milliseconds for execution
SetDefaultExecutionTimeout(200 * time.Millisecond)

// rest is the same...
result, err := templ8go.ResolveJSExpression(bindings, "user.name")
...
```

---

## Supported Expressions

Since we use V8 engine underneath, many things are possible.

- **Simple Arithmetic**:
    - Template: `The sum of 5 and 3 is {{ 5 + 3 }}.`
    - Bindings: `{}`
    - Output: `The sum of 5 and 3 is 8.`

- **Conditional Greetings**:
    - Template: `Good {{ hour < 12 ? 'morning' : 'afternoon' }}, {{ user.name }}!`
    - Bindings: `{ "hour": 9, "user": {"name": "Alice"} }`
    - Output: `Good morning, Alice!`

- **Array Operations**:
    - Template: `Users list: {{ users.map(user => user.name).join(', ') }}`
    - Bindings `{users: [{name: 'Alice'}, {name: 'Bob'}, {name: 'Charlie'}]}`
    - Output: `Users list Alice, Bob, Charlie`

- **Object Manipulation**:
    - Template: `{{ user.firstName }} {{ user.lastName }} is {{ user.age }} years old.`
    - Bindings: `{ "user": {"firstName": "John", "lastName": "Doe", "age": 28} }`
    - Output: `John Doe is 28 years old.`

- **Logical Operations**:
    - Template: `You are {{ age >= 18 ? 'an adult' : 'a minor' }}.`
    - Bindings: `{ "age": 20 }`
    - Output: `You are an adult.`

- **String Concatenation**:
    - Template: `{{ 'Hello, ' + user.name + '!'}}`
    - Bindings: `{ "user": {"name": "Jane"} }`
    - Output: `Hello, Jane!`

- **Using JavaScript Functions**:
    - Template: `Your score is {{ Math.min(score, 100) }}.`
    - Bindings: `{ "score": 105 }`
    - Output: `Your score is 100.`

- **Nested Object Access**:
    - Template: `Project {{ project.details.name }} is due on {{ project.details.dueDate }}.`
    - Bindings: `{ "project": {"details": {"name": "Apollo", "dueDate": "2024-03-01"}} }`
    - Output: `Project Apollo is due on 2024-03-01.`

- **Complex Expressions**:
    - Template: `{{ user.isActive ? user.name + ' is active and has ' + user.roles.length + ' roles' : user.name + ' is not active' }}.`
    - Bindings: `{ "user": {"name": "Eve", "isActive": true, "roles": ["admin", "editor"]} }`
    - Output: `Eve is active and has 2 roles.`

---

## Development

To contribute or modify `templ8go`, clone the repository and ensure all dependencies are properly set up. Run tests with:

```sh
go test ./...
```

Feel free to submit pull requests or create issues for bugs, features, or suggestions.

---

## Contributing

Contributions are more than welcome! If you have an idea for an improvement or
find a bug, please feel free to fork the repository, make your changes, and
submit a pull request.

## License

`templ8go` is made available under the MIT License. For more details, see the LICENSE file in the repository.
