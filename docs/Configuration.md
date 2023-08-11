# Configuration

Expr can be configured with options. For example, you can pass the environment with variables and functions.

## AllowUndefinedVariables()

This option allows undefined variables in the expression. By default, Expr will return an error 
if the expression contains undefined variables.

```go
program, err := expr.Compile(`foo + bar`, expr.AllowUndefinedVariables())
```

## AsBool()

This option forces the expression to return a boolean value. If the expression returns a non-boolean value,
Expr will return an error.

```go
program, err := expr.Compile(`Title contains "Hello"`, expr.AsBool())
```

## AsFloat64()

This option forces the expression to return a float64 value. If the expression returns a non-float64 value,
Expr will return an error.

```go
program, err := expr.Compile(`42`, expr.AsFloat64())
```

:::note
If the expression returns integer value, Expr will convert it to float64.
:::

## AsInt()

This option forces the expression to return an int value. If the expression returns a non-int value,
Expr will return an error.

```go
program, err := expr.Compile(`42`, expr.AsInt())
```

:::note
If the expression returns a float value, Expr truncates it to int.
:::

## AsInt64()

Same as `AsInt()` but returns an int64 value.

```go
program, err := expr.Compile(`42`, expr.AsInt64())
```

## AsKind()

This option forces the expression to return a value of the specified kind. 
If the expression returns a value of a different kind, Expr will return an error.

```go
program, err := expr.Compile(`42`, expr.AsKind(reflect.String))
```

## ConstExpr()

This option tells Expr to treat specified functions as constant expressions. 
If all arguments of the function are constants, Expr will replace the function call with the result 
during the compile step.

```go
program, err := expr.Compile(`fib(42)`, expr.ConstExpr("fib"))
```

[ConstExpr Example](https://pkg.go.dev/github.com/antonmedv/expr?tab=doc#ConstExpr)

## Env()

This option passes the environment with variables and functions to the expression.

```go
program, err := expr.Compile(`foo + bar`, expr.Env(Env{}))
```

## Function()

This option adds a function to the expression.

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


## Operator()

This options defines an [operator overloading](Operator-Overloading.md).

## Optimize()

This option enables [optimizations](Internals.md). By default, Expr will optimize the expression.

## Patch()

This option allows you to [patch the expression](Visitor-and-Patch.md) before compilation.
