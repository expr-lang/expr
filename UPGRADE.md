# Upgrade

This file contains information on how to upgrade between major versions.

## 0.x → 1.x

### Deprecated `expr.Names` and `expr.Funcs`

Use `expr.Define` or `expr.With` instead of.

In new version you can specify types of variables, if you don't need type info use `interface{}`:

```go
expr.Parse(code, expr.Names('Foo', 'Bar'))

// Change to

var foo, bar interface{}
expr.Parse(code, expr.Define('Foo', foo), expr.Define('Bar', bar))

// Or

type params struct {
    Foo interface{}
    Bar interface{}
}
expr.Parse(code, expr.With(params{}))
```
  
