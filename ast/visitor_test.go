package ast_test

import (
	"github.com/antonmedv/expr/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

type visitor struct {
	ast.BaseVisitor
	lastSeenInteger int64
}

func (v *visitor) IntegerNode(node *ast.IntegerNode) {
	v.lastSeenInteger = node.Value
}

func TestWalk(t *testing.T) {
	node := &ast.BinaryNode{
		Operator: "",
		Left:     &ast.IntegerNode{Value: 12},
		Right:    &ast.IntegerNode{Value: 42},
	}

	visitor := &visitor{}
	ast.Walk(node, visitor)
	assert.Equal(t, int64(42), visitor.lastSeenInteger)
}
