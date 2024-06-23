# Environment

The environment is a map or a struct that contains the variables and functions that the expression can access.

## Struct as Environment

Let's consider the following example:

```go
type Env struct {
    UpdatedAt time.Time
    Posts     []Post
    Map       map[string]string `expr:"tags"`
}
```

The `Env` struct contains 3 variables that the expression can access: `UpdatedAt`, `Posts`, and `tags`.

:::info
The `expr` tag is used to rename the `Map` field to `tags` variable in the expression.
:::

The `Env` struct can also contain methods. The methods defined on the struct become functions that the expression can
call.

```go
func (Env) Format(t time.Time) string {
    return t.Format(time.RFC822)
}
```

:::tip
Methods defined on embedded structs are also accessible.
```go
type Env struct {
    Helpers
}

type Helpers struct{}

func (Helpers) Format(t time.Time) string {
    return t.Format(time.RFC822)
}
```
:::

We can use an empty struct `Env{}` to with [expr.Env](https://pkg.go.dev/github.com/expr-lang/expr#Env) to create an environment. Expr will use reflection to find 
the fields and methods of the struct.

```go
program, err := expr.Compile(code, expr.Env(Env{}))
```

Compiler will type check the expression against the environment. After the compilation, we can run the program with the environment.
You should use the same type of environment that you passed to the `expr.Env` function.

```go
output, err := expr.Run(program, Env{
    UpdatedAt: time.Now(),
    Posts:     []Post{{Title: "Hello, World!"}},
    Map:       map[string]string{"tag1": "value1"},
})
```

## Map as Environment

You can also use a map as an environment.

```go
env := map[string]any{
    "UpdatedAt": time.Time{},
    "Posts":     []Post{},
    "tags":      map[string]string{},
    "sprintf":   fmt.Sprintf,
}

program, err := expr.Compile(code, expr.Env(env))
```

A map defines variables and functions that the expression can access. The key is the variable name, and the type
is the value's type.

```go
env := map[string]any{
    "object": map[string]any{
        "field": 42,
    }, 
}
```

Expr will infer the type of the `object` variable as `map[string]any`.

By default, Expr will return an error if unknown variables are used in the expression.

You can disable this behavior by passing [`AllowUndefinedVariables`](https://pkg.go.dev/github.com/expr-lang/expr#AllowUndefinedVariables) option to the compiler.
