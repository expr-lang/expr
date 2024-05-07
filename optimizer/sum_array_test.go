package optimizer_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/optimizer"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm"
)

func BenchmarkSumArray(b *testing.B) {
	env := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	program, err := expr.Compile(`sum([a, b, c, d])`, expr.Env(env))
	require.NoError(b, err)

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 10, out)
}

func TestOptimize_sum_array(t *testing.T) {
	tree, err := parser.Parse(`sum([a, b])`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BinaryNode{
		Operator: "+",
		Left:     &ast.IdentifierNode{Value: "a"},
		Right:    &ast.IdentifierNode{Value: "b"},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_sum_array_3(t *testing.T) {
	tree, err := parser.Parse(`sum([a, b, c])`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BinaryNode{
		Operator: "+",
		Left:     &ast.IdentifierNode{Value: "a"},
		Right: &ast.BinaryNode{
			Operator: "+",
			Left:     &ast.IdentifierNode{Value: "b"},
			Right:    &ast.IdentifierNode{Value: "c"},
		},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}
