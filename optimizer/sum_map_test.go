package optimizer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/optimizer"
	"github.com/expr-lang/expr/parser"
)

func TestOptimize_sum_map(t *testing.T) {
	tree, err := parser.Parse(`sum(map(users, {.Age}))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "sum",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.MemberNode{
					Node:     &ast.PointerNode{},
					Property: &ast.StringNode{Value: "Age"},
				},
			},
		},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}
