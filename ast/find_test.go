package ast_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr/ast"
)

func TestFind(t *testing.T) {
	left := &ast.IdentifierNode{
		Value: "a",
	}
	var root ast.Node = &ast.BinaryNode{
		Operator: "+",
		Left:     left,
		Right: &ast.IdentifierNode{
			Value: "b",
		},
	}

	x := ast.Find(root, func(node ast.Node) bool {
		if n, ok := node.(*ast.IdentifierNode); ok {
			return n.Value == "a"
		}
		return false
	})

	require.Equal(t, left, x)
}
