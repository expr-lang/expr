# Functions

Expr comes with a set of [builtin](language-definition.md) functions, but you can also define your own functions.

The easiest way to define a custom function is to add it to the environment.

```go
env := map[string]any{
    "add": func(a, b int) int {
        return a + b
    },
}
```

Or you can use functions defined on a struct:

```go
type Env struct{}

func (Env) Add(a, b int) int {
    return a + b
}
```

:::info
If functions are marked with [`ConstExpr`](./configuration.md#constexpr) option, they will be evaluated at compile time.
:::

The best way to define a function from a performance perspective is to use a [`Function`](https://pkg.go.dev/github.com/expr-lang/expr#Function) option.

```go
atoi := expr.Function(
    "atoi",
    func(params ...any) (any, error) {
        return strconv.Atoi(params[0].(string))
    },
)

program, err := expr.Compile(`atoi("42")`, atoi)
```

Type checker sees the `atoi` function as a function with a variadic number of arguments of type `any`, and returns 
a value of type `any`. But, we can specify the types of arguments and the return value by adding the correct function 
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

Or we can simply reuse the strconv.Atoi function as a type:

```go
atoi := expr.Function(
    "atoi",
    func(params ...any) (any, error) {
        return strconv.Atoi(params[0].(string))
    },
    // highlight-next-line
    strconv.Atoi,
)
```

It is possible to define multiple signatures for a function:

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
