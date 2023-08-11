# Configuration

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
