# Patch

Sometimes it may be necessary to modify an expression before the compilation.
For example, you may want to replace a variable with a constant, transform an expression into a function call, 
or even modify the expression to use a different operator.

## Simple example

Let's start with a simple example. We have an expression that uses a variable `foo`:

```go
program, err := expr.Compile(`foo + bar`)
```

We want to replace the `foo` variable with a constant `42`. First, we need to implement a [visitor](./visitor.md):

```go
type FooPatcher struct{}

func (FooPatcher) Visit(node *ast.Node) {
    if n, ok := (*node).(*ast.IdentifierNode); ok && n.Value == "foo" {
        // highlight-next-line
        ast.Patch(node, &ast.IntegerNode{Value: 42})
    }
}
```

We used the [ast.Patch](https://pkg.go.dev/github.com/expr-lang/expr/ast#Patch) function to replace the `foo` variable with an integer node.

Now we can use the `FooPatcher` to modify the expression on compilation via the [expr.Patch](https://pkg.go.dev/github.com/expr-lang/expr#Patch) option:

```go
program, err := expr.Compile(`foo + bar`, expr.Patch(FooPatcher{}))
```

## Advanced example

Let's consider a more complex example. We have an expression that uses variables `foo` and `bar` of type `Decimal`:

```go
type Decimal struct {
    Value int
}
```

And we want to transform the following expression:

```expr
a + b + c
```

Into functions calls that accept `Decimal` arguments:

```expr
add(add(a, b), c)
```

First, we need to implement a visitor that will transform the expression:

```go
type DecimalPatcher struct{}

var decimalType = reflect.TypeOf(Decimal{})

func (DecimalPatcher) Visit(node *ast.Node) {
    if n, ok := (*node).(*ast.BinaryNode); ok && n.Operator == "+" {
        
        if !n.Left.Type().AssignableTo(decimalType) {
            return // skip, left side is not a Decimal
        }
		
        if !n.Right.Type().AssignableTo(decimalType) {
            return // skip, right side is not a Decimal
        }
		
        // highlight-start
        callNode := &ast.CallNode{
            Callee:    &ast.IdentifierNode{Value: "add"},
            Arguments: []ast.Node{n.Left, n.Right},
        }
        ast.Patch(node, callNode)
        // highlight-end
		
        (*node).SetType(decimalType) // set the type, so the patcher can be applied recursively
    }
}
```

We used [Type()](https://pkg.go.dev/github.com/expr-lang/expr/ast#Node.Type) method to get the type of the expression node.
The `AssignableTo` method is used to check if the type is `Decimal`. If both sides are `Decimal`, we replace the expression with a function call.

The important part of this patcher is to set correct types for the nodes. As we constructed a new `CallNode`, it lacks the type information.
So after the first patcher run, if we want the patcher to be applied recursively, we need to set the type of the node.


Now we can use the `DecimalPatcher` to modify the expression:

```go
env := map[string]interface{}{
    "a": Decimal{1},
    "b": Decimal{2},
    "c": Decimal{3},
    "add": func(x, y Decimal) Decimal {
        return Decimal{x.Value + y.Value}
    },
}

code := `a + b + c`

// highlight-next-line
program, err := expr.Compile(code, expr.Env(env), expr.Patch(DecimalPatcher{}))
if err != nil {
    panic(err)
}

output, err := expr.Run(program, env)
if err != nil {
    panic(err)
}

fmt.Println(output) // Decimal{6}
```


:::info
Expr comes with already implemented patcher that simplifies operator overloading.

The `DecimalPatcher` can be replaced with the [Operator](https://pkg.go.dev/github.com/expr-lang/expr#Operator) option.

```go
program, err := expr.Compile(code, expr.Env(env), expr.Operator("+", "add"))
```

Operator overloading patcher will check if provided functions (`"add"`) satisfy the operator (`"+"`), and
replace the operator with the function call.
:::
