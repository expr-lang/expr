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
		Left: &ast.BinaryNode{
			Operator: "+",
			Left:     &ast.IdentifierNode{Value: "a"},
			Right:    &ast.IdentifierNode{Value: "b"},
		},
		Right: &ast.IdentifierNode{Value: "c"},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_sum_array_float_order(t *testing.T) {
	// sum() adds left to right, so the optimized program must too:
	// 1e16 + -1e16 + 1 is 1, but 1e16 + (-1e16 + 1) is 0.
	x, y, z := 0.1, 0.2, 0.3
	env := map[string]any{
		"a": 1e16,
		"b": -1e16,
		"c": 1.0,
		"x": x,
		"y": y,
		"z": z,
	}

	tests := []struct {
		code string
		want any
	}{
		{`sum([a, b, c])`, 1.0},
		{`sum([x, y, z])`, x + y + z},
		{`sum([x, y, z]) == x + y + z`, true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code, expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, out)

			unoptimized, err := expr.Compile(tt.code, expr.Env(env), expr.Optimize(false))
			require.NoError(t, err)

			want, err := expr.Run(unoptimized, env)
			require.NoError(t, err)
			assert.Equal(t, want, out)
		})
	}
}
