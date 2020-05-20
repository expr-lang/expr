# Visitor and Patch

[ast](https://pkg.go.dev/github.com/antonmedv/expr/ast?tab=doc) package provides `ast.Visitor` interface and `ast.Walk` function. 
You can use it for traveling ast tree of compiled program.

For example if you want to collect all variable names:

```go
package main

import (
	"fmt"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/parser"
)

type visitor struct {
	identifiers []string
}

func (v *visitor) Enter(node *ast.Node) {}
func (v *visitor) Exit(node *ast.Node) {
	if n, ok := (*node).(*ast.IdentifierNode); ok {
		v.identifiers = append(v.identifiers, n.Value)
	}
}

func main() {
	tree, err := parser.Parse("foo + bar")
	if err != nil {
		panic(err)
	}

	visitor := &visitor{}
	ast.Walk(&tree.Node, visitor)

	fmt.Printf("%v", visitor.identifiers) // outputs [foo bar]
}
```

## Patch

Implemented visitor can be applied before compiling AST to bytecode in `expr.Compile` function.

```go
program, err := expr.Compile(code, expr.Patch(&visitor{}))
```

This can be useful for some edge cases, there you want to extend functionality of **Expr** language. 
In next example we are going to replace expression `list[-1]` with `list[len(list)-1]`.

```go
package main

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
)

func main() {
	env := map[string]interface{}{
		"list": []int{1, 2, 3},
	}

	code := `list[-1]` // will output 3

	program, err := expr.Compile(code, expr.Env(env), expr.Patch(&patcher{}))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

type patcher struct{}

func (p *patcher) Enter(_ *ast.Node) {}
func (p *patcher) Exit(node *ast.Node) {
	n, ok := (*node).(*ast.IndexNode)
	if !ok {
		return
	}
	unary, ok := n.Index.(*ast.UnaryNode)
	if !ok {
		return
	}
	if unary.Operator == "-" {
		ast.Patch(&n.Index, &ast.BinaryNode{
			Operator: "-",
			Left:     &ast.BuiltinNode{Name: "len", Arguments: []ast.Node{n.Node}},
			Right:    unary.Node,
		})
	}

}
```

Type information is also available. Here is an example, there all `fmt.Stringer` interface automatically 
converted to `string` type.

```go
package main

import (
	"fmt"
	"reflect"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
)

func main() {
	code := `Price == "$100"`

	program, err := expr.Compile(code, expr.Env(Env{}), expr.Patch(&stringerPatcher{}))
	if err != nil {
		panic(err)
	}

	env := Env{100_00}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

type Env struct {
	Price Price
}

type Price int

func (p Price) String() string {
	return fmt.Sprintf("$%v", int(p)/100)
}

var stringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

type stringerPatcher struct{}

func (p *stringerPatcher) Enter(_ *ast.Node) {}
func (p *stringerPatcher) Exit(node *ast.Node) {
	t := (*node).Type()
	if t == nil {
		return
	}
	if t.Implements(stringer) {
		ast.Patch(node, &ast.MethodNode{
			Node:   *node,
			Method: "String",
		})
	}

}
```

* [Contents](README.md)
* Next: [Optimizations](Optimizations.md)
