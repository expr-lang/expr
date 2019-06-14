package ast_test

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/antonmedv/expr.v2/ast"
	"testing"
)

type visitor struct {
	ast.BaseVisitor
	identifiers []string
}

func (v *visitor) IdentifierNode(node *ast.IdentifierNode) {
	v.identifiers = append(v.identifiers, node.Value)
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

type patcher struct {
	ast.BaseVisitor
}

func (p *patcher) Node(node *ast.Node) {
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
