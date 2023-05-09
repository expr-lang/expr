# Visitor and Patch

The [ast](https://pkg.go.dev/github.com/antonmedv/expr/ast?tab=doc) package 
provides the `ast.Visitor` interface and the `ast.Walk` function. It can be
used to customize the AST before compiling.

For example, to collect all variable names:

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

func (v *visitor) Visit(node *ast.Node) {
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

Specify a visitor to modify the AST with `expr.Patch` function.  

```go
program, err := expr.Compile(code, expr.Patch(&visitor{}))
```
 
For example, we are going to replace the expression `list[-1]` with the `list[len(list)-1]`.

```go
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

func (p *patcher) Visit(node *ast.Node) {
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

Type information is also available. In the following example, any struct
implementing the `fmt.Stringer` interface is automatically converted to `string` type.

```go
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

func (p *stringerPatcher) Visit(node *ast.Node) {
	t := (*node).Type()
	if t == nil {
		return
	}
	if t.Implements(stringer) {
		ast.Patch(node, &ast.CallNode{
			Callee: &ast.MemberNode{
				Node:  *node,
				Field: "String",
				Property: &ast.StringNode{Value: "String"},
			},
		})
	}

}
```
