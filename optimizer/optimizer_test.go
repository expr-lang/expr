package optimizer_test

import (
	"strings"
	"testing"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptimize_constant_folding(t *testing.T) {
	tree, err := parser.Parse(`[1,2,3][5*5-25]`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.IndexNode{
		Node:  &ast.ConstantNode{Value: []int{1, 2, 3}},
		Index: &ast.IntegerNode{Value: 0},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_in_array(t *testing.T) {
	config := conf.New(map[string]int{"v": 0})

	tree, err := parser.Parse(`v in [1,2,3]`)
	require.NoError(t, err)

	_, err = checker.Check(tree, config)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BinaryNode{
		Operator: "in",
		Left:     &ast.IdentifierNode{Value: "v"},
		Right:    &ast.ConstantNode{Value: map[int]struct{}{1: {}, 2: {}, 3: {}}},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_in_range(t *testing.T) {
	tree, err := parser.Parse(`age in 18..31`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	left := &ast.IdentifierNode{
		Value: "age",
	}
	expected := &ast.BinaryNode{
		Operator: "and",
		Left: &ast.BinaryNode{
			Operator: ">=",
			Left:     left,
			Right: &ast.IntegerNode{
				Value: 18,
			},
		},
		Right: &ast.BinaryNode{
			Operator: "<=",
			Left:     left,
			Right: &ast.IntegerNode{
				Value: 31,
			},
		},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_const_range(t *testing.T) {
	tree, err := parser.Parse(`-1..1`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.ConstantNode{
		Value: []int{-1, 0, 1},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_const_expr(t *testing.T) {
	tree, err := parser.Parse(`upper("hello")`)
	require.NoError(t, err)

	env := map[string]interface{}{
		"upper": strings.ToUpper,
	}

	config := conf.New(env)
	config.ConstExpr("upper")

	err = optimizer.Optimize(&tree.Node, config)
	require.NoError(t, err)

	expected := &ast.ConstantNode{
		Value: "HELLO",
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}
