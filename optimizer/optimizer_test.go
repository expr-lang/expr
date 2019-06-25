package optimizer_test

import (
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptimize_constant_folding(t *testing.T) {
	tree, err := parser.Parse(`v in [1,2,3,4,5]`)
	require.NoError(t, err)

	optimizer.Optimize(&tree.Node)

	expected := &ast.BinaryNode{
		Operator: "in",
		Left:     &ast.IdentifierNode{Value: "v"},
		Right:    &ast.ConstantNode{Value: []int{1, 2, 3, 4, 5}},
	}

	assert.Equal(t, litter.Sdump(expected), litter.Sdump(tree.Node))
}

func TestOptimize_in_range(t *testing.T) {
	tree, err := parser.Parse(`age in 18..31`)
	require.NoError(t, err)

	optimizer.Optimize(&tree.Node)

	expected := &ast.BinaryNode{
		Operator: "in",
		Left: &ast.IdentifierNode{
			Value: "age",
		},
		Right: &ast.BinaryNode{
			Operator: "..",
			Left: &ast.IntegerNode{
				Value: 18,
			},
			Right: &ast.IntegerNode{
				Value: 31,
			},
		},
	}

	assert.Equal(t, litter.Sdump(expected), litter.Sdump(tree.Node))
}
