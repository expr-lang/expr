# Getting Started

**Expr** is a simple, fast and extensible expression language for Go. It is
designed to be easy to use and integrate into your application.


## Evaluate

For simple use cases with one execution of an expression, you can use the 
`expr.Eval` function. It takes an expression and a map of variables and returns
the result of the expression.

```go
package main

import (
	"fmt"
	"github.com/antonmedv/expr"
)

func main() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	out, err := expr.Eval("foo + bar", env)

	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

## Compile

Usually, we want to compile, type check and verify that the expression returns a 
boolean (or another type). For example, if a user saves an expression from a
[web UI](https://antonmedv.github.io/expr/).

```go
	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf, // Use any functions.
	}

	code := `sprintf(greet, names[0])`

	// Compile the code into a bytecode. This step can be done only once and
	// the program may be reused. Specify an environment for the type check.
	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
```

An environment can be a struct and structs methods can be used as
functions. Expr supports embedded structs and methods defined on them too.

The struct fields can be renamed by adding struct tags such as `expr:"name"`.

```go
type Env struct {
	Messages []Message `expr:"messages"`
}

// Methods defined on the struct become functions.
func (Env) Format(t time.Time) string { return t.Format(time.RFC822) }

type Message struct {
	Text string
	Date time.Time
}

func main() {
	code := `map(filter(messages, len(.Text) > 0), Format(.Date) + '\t' + .Text + '\n')`

	// We can use an empty instance of the struct as an environment.
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		panic(err)
	}

	env := Env{
		Messages: []Message{
			{"Oh My God!", time.Now()}, 
			{"How you doin?", time.Now()}, 
			{"Could I be wearing any more clothes?", time.Now()},
		},
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
}
```

## Configuration

Expr can be configured to do more things. For example, with [AllowUndefinedVariables](https://pkg.go.dev/github.com/antonmedv/expr#AllowUndefinedVariables) or [AsBool](https://pkg.go.dev/github.com/antonmedv/expr#AsBool) to expect the boolean result from the expression.

```go
program, err := expr.Compile(code, expr.Env(Env{}), expr.AllowUndefinedVariables(), expr.AsBool())
```

## Functions

Expr supports any Go functions. For example, you can use `fmt.Sprintf` or methods of your structs. 
Also, Expr supports functions configured via [`expr.Function(name, fn[, ...types])`](https://pkg.go.dev/github.com/antonmedv/expr#Function) option.

```go
	atoi := expr.Function(
		"atoi",
		func(params ...any) (any, error) {
			return strconv.Atoi(params[0].(string))
		},
	)

	program, err := expr.Compile(`atoi("42")`, atoi)
```

Expr sees the `atoi` function as a function with a variadic number of arguments of type `any` and returns a value of type `any`. But, we can specify the types of arguments and the return value by adding the correct function
signature or multiple signatures.

```go
	atoi := expr.Function(
		"atoi",
		func(params ...any) (any, error) {
			return strconv.Atoi(params[0].(string))
		},
		new(func(string) int),
	)
```

Or we can simply reuse the `strconv.Atoi` function.

```go
	atoi := expr.Function(
		"atoi",
		func(params ...any) (any, error) {
			return strconv.Atoi(params[0].(string))
		},
		strconv.Atoi,
	)
```

Here is another example with a few function signatures:

```go
	toInt := expr.Function(
		"toInt",
		func(params ...any) (any, error) {
			switch params[0].(type) {
			case float64:
				return int(params[0].(float64)), nil
			case string:
				return strconv.Atoi(params[0].(string))
			}
			return nil, fmt.Errorf("invalid type")
		},
		new(func(float64) int),
		new(func(string) int),
	)
```

Functions taking `context.Context` as first argument are also accepted:

```go
	username := expr.FunctionWithContext(
		"username",
		func(params ...any) (any, error) {
			ctx := params[0].(context.Context)
			return ctx.Value("user"), nil
		},
		new(func() string),
	)
```

Program shall be run using `expr.RunWithContext`.
