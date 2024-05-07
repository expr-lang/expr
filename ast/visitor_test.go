package ast_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"

	"github.com/expr-lang/expr/ast"
)

type visitor struct {
	identifiers []string
}

func (v *visitor) Visit(node *ast.Node) {
	if n, ok := (*node).(*ast.IdentifierNode); ok {
		v.identifiers = append(v.identifiers, n.Value)
	}
}

func TestWalk(t *testing.T) {
	var node ast.Node
	node = &ast.BinaryNode{
		Operator: "+",
		Left:     &ast.IdentifierNode{Value: "foo"},
		Right:    &ast.IdentifierNode{Value: "bar"},
	}

	visitor := &visitor{}
	ast.Walk(&node, visitor)
	assert.Equal(t, []string{"foo", "bar"}, visitor.identifiers)
}

type patcher struct{}

func (p *patcher) Visit(node *ast.Node) {
	if _, ok := (*node).(*ast.IdentifierNode); ok {
		*node = &ast.NilNode{}
	}
}

func TestWalk_patch(t *testing.T) {
	var node ast.Node
	node = &ast.BinaryNode{
		Operator: "+",
		Left:     &ast.IdentifierNode{Value: "foo"},
		Right:    &ast.IdentifierNode{Value: "bar"},
	}

	patcher := &patcher{}
	ast.Walk(&node, patcher)
	assert.IsType(t, &ast.NilNode{}, node.(*ast.BinaryNode).Left)
	assert.IsType(t, &ast.NilNode{}, node.(*ast.BinaryNode).Right)
}
