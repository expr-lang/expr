package ast_test

import (
	"github.com/antonmedv/expr/ast"
	"github.com/stretchr/testify/assert"
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
	node := &ast.BinaryNode{
		Operator: "+",
		Left:     &ast.IdentifierNode{Value: "foo"},
		Right:    &ast.IdentifierNode{Value: "bar"},
	}

	visitor := &visitor{}
	ast.Walk(node, visitor)
	assert.Equal(t, []string{"foo", "bar"}, visitor.identifiers)
}
