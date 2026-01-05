package optimizer_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/optimizer"
	"github.com/expr-lang/expr/parser"
)

func TestOptimize_sum_range(t *testing.T) {
	tree, err := parser.Parse(`sum(1..100)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.IntegerNode{Value: 5050}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_sum_range_different_values(t *testing.T) {
	tests := []struct {
		expr string
		want int
	}{
		{`sum(1..10)`, 55},
		{`sum(1..100)`, 5050},
		{`sum(5..10)`, 45},
		{`sum(0..100)`, 5050},
		{`sum(1..1)`, 1},
		{`sum(0..0)`, 0},
		{`sum(10..20)`, 165},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr)
			require.NoError(t, err)

			output, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)
		})
	}
}

func TestOptimize_sum_range_with_predicate(t *testing.T) {
	tests := []struct {
		expr string
		want int
	}{
		// # (identity) - same as sum(m..n)
		{`sum(1..10, #)`, 55},
		{`sum(1..100, #)`, 5050},

		// # * k (multiply by constant)
		{`sum(1..10, # * 2)`, 110},    // 2 * 55
		{`sum(1..100, # * 2)`, 10100}, // 2 * 5050
		{`sum(1..10, # * 0)`, 0},
		{`sum(1..10, # * 1)`, 55},

		// k * # (multiply by constant, reversed)
		{`sum(1..10, 2 * #)`, 110},
		{`sum(1..100, 3 * #)`, 15150}, // 3 * 5050

		// # + k (add constant to each element)
		{`sum(1..10, # + 1)`, 65},    // 55 + 10*1
		{`sum(1..100, # + 1)`, 5150}, // 5050 + 100*1
		{`sum(1..10, # + 0)`, 55},
		{`sum(1..10, # + 10)`, 155}, // 55 + 10*10

		// k + # (add constant, reversed)
		{`sum(1..10, 1 + #)`, 65},
		{`sum(1..100, 5 + #)`, 5550}, // 5050 + 100*5

		// # - k (subtract constant from each element)
		{`sum(1..10, # - 1)`, 45},    // 55 - 10*1
		{`sum(1..100, # - 1)`, 4950}, // 5050 - 100*1
		{`sum(1..10, # - 0)`, 55},

		// k - # (constant minus each element)
		{`sum(1..10, 10 - #)`, 45}, // 10*10 - 55
		{`sum(1..10, 0 - #)`, -55}, // 10*0 - 55
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr)
			require.NoError(t, err)

			output, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)
		})
	}
}

func TestOptimize_sum_range_with_predicate_ast(t *testing.T) {
	// Verify that sum(1..10, # * 2) is optimized to a constant
	tree, err := parser.Parse(`sum(1..10, # * 2)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.IntegerNode{Value: 110}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_reduce_range_sum(t *testing.T) {
	tree, err := parser.Parse(`reduce(1..100, # + #acc)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.IntegerNode{Value: 5050}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_reduce_range_sum_different_values(t *testing.T) {
	tests := []struct {
		expr string
		want int
	}{
		{`reduce(1..10, # + #acc)`, 55},
		{`reduce(1..100, # + #acc)`, 5050},
		{`reduce(5..10, # + #acc)`, 45},
		{`reduce(0..100, # + #acc)`, 5050},
		{`reduce(1..1, # + #acc)`, 1},
		{`reduce(10..20, # + #acc)`, 165},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr)
			require.NoError(t, err)

			output, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)
		})
	}
}

func TestOptimize_reduce_range_sum_reverse_order(t *testing.T) {
	// Test #acc + # (reverse order) - should also be optimized
	tree, err := parser.Parse(`reduce(1..100, #acc + #)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.IntegerNode{Value: 5050}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_reduce_range_sum_with_initial_value(t *testing.T) {
	// Test reduce with initialValue: reduce(1..100, # + #acc, 10) => 5050 + 10 = 5060
	tree, err := parser.Parse(`reduce(1..100, # + #acc, 10)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.IntegerNode{Value: 5060}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_reduce_range_sum_with_initial_value_different_values(t *testing.T) {
	tests := []struct {
		expr string
		want int
	}{
		{`reduce(1..10, # + #acc, 0)`, 55},
		{`reduce(1..10, # + #acc, 10)`, 65},
		{`reduce(1..100, # + #acc, 0)`, 5050},
		{`reduce(1..100, # + #acc, 100)`, 5150},
		{`reduce(5..10, # + #acc, 5)`, 50},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr)
			require.NoError(t, err)

			output, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)
		})
	}
}

func TestOptimize_sum_range_reversed(t *testing.T) {
	// When n < m (e.g., 10..1), the range is empty and sum should return 0.
	// The optimization should NOT apply (n >= m check), so runtime handles it.
	tests := []struct {
		expr string
		want int
	}{
		{`sum(10..1)`, 0},
		{`sum(5..3)`, 0},
		{`sum(100..1)`, 0},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr)
			require.NoError(t, err)

			output, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)
		})
	}
}

func TestOptimize_sum_range_reversed_not_optimized(t *testing.T) {
	// Verify that reversed ranges are NOT optimized (left as BuiltinNode)
	tree, err := parser.Parse(`sum(10..1)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	// Should still be a BuiltinNode, not an IntegerNode
	_, isBuiltin := tree.Node.(*ast.BuiltinNode)
	assert.True(t, isBuiltin, "reversed range should not be optimized")
}

func TestOptimize_reduce_range_reversed_errors(t *testing.T) {
	// reduce on empty range (reversed) should error at runtime
	program, err := expr.Compile(`reduce(10..1, # + #acc)`)
	require.NoError(t, err)

	_, err = expr.Run(program, nil)
	require.Error(t, err, "reduce on empty range should error")
}

func BenchmarkSumRange_Optimized(b *testing.B) {
	program, err := expr.Compile(`sum(1..100)`)
	require.NoError(b, err)

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = expr.Run(program, nil)
	}
	b.StopTimer()

	require.Equal(b, 5050, out)
}

func BenchmarkReduceRangeSum_Optimized(b *testing.B) {
	program, err := expr.Compile(`reduce(1..100, # + #acc)`)
	require.NoError(b, err)

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = expr.Run(program, nil)
	}
	b.StopTimer()

	require.Equal(b, 5050, out)
}

func BenchmarkSumRange_Unoptimized(b *testing.B) {
	program, err := expr.Compile(`sum(1..100)`, expr.Optimize(false))
	require.NoError(b, err)

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = expr.Run(program, nil)
	}
	b.StopTimer()

	require.Equal(b, 5050, out)
}
