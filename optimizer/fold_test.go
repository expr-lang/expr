package optimizer_test

import (
	"reflect"
	"testing"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/optimizer"
	"github.com/expr-lang/expr/parser"
)

func TestOptimize_constant_folding(t *testing.T) {
	tree, err := parser.Parse(`[1,2,3][5*5-25]`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.MemberNode{
		Node:     &ast.ConstantNode{Value: []any{1, 2, 3}},
		Property: &ast.IntegerNode{Value: 0},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_constant_folding_with_floats(t *testing.T) {
	tree, err := parser.Parse(`1 + 2.0 * ((1.0 * 2) / 2) - 0`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.FloatNode{Value: 3.0}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
	assert.Equal(t, reflect.Float64, tree.Node.Type().Kind())
}

func TestOptimize_constant_folding_with_bools(t *testing.T) {
	tree, err := parser.Parse(`(true and false) or (true or false) or (false and false) or (true and (true == false))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BoolNode{Value: true}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_constant_folding_filter_filter(t *testing.T) {
	tree, err := parser.Parse(`filter(filter(1..2, true), true)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "filter",
		Arguments: []ast.Node{
			&ast.BinaryNode{
				Operator: "..",
				Left: &ast.IntegerNode{
					Value: 1,
				},
				Right: &ast.IntegerNode{
					Value: 2,
				},
			},
			&ast.PredicateNode{
				Node: &ast.BoolNode{
					Value: true,
				},
			},
		},
		Throws: false,
		Map:    nil,
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}
