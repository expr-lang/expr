package optimizer_test

import (
	"testing"

	"github.com/expr-lang/expr"
	. "github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/optimizer"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm"
)

func TestOptimize_count_threshold_gt(t *testing.T) {
	tree, err := parser.Parse(`count(items, .active) > 100`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	// Operator should remain >, but count should have threshold set
	binary, ok := tree.Node.(*BinaryNode)
	require.True(t, ok, "expected BinaryNode, got %T", tree.Node)
	assert.Equal(t, ">", binary.Operator)

	count, ok := binary.Left.(*BuiltinNode)
	require.True(t, ok, "expected BuiltinNode, got %T", binary.Left)
	assert.Equal(t, "count", count.Name)
	require.NotNil(t, count.Threshold)
	assert.Equal(t, 101, *count.Threshold) // threshold = N + 1 for > operator
}

func TestOptimize_count_threshold_gte(t *testing.T) {
	tree, err := parser.Parse(`count(items, .active) >= 50`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	// Operator should remain >=, but count should have threshold set
	binary, ok := tree.Node.(*BinaryNode)
	require.True(t, ok, "expected BinaryNode, got %T", tree.Node)
	assert.Equal(t, ">=", binary.Operator)

	count, ok := binary.Left.(*BuiltinNode)
	require.True(t, ok, "expected BuiltinNode, got %T", binary.Left)
	assert.Equal(t, "count", count.Name)
	require.NotNil(t, count.Threshold)
	assert.Equal(t, 50, *count.Threshold) // threshold = N for >= operator
}

func TestOptimize_count_threshold_lt(t *testing.T) {
	tree, err := parser.Parse(`count(items, .active) < 100`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	// Operator should remain <, but count should have threshold set
	binary, ok := tree.Node.(*BinaryNode)
	require.True(t, ok, "expected BinaryNode, got %T", tree.Node)
	assert.Equal(t, "<", binary.Operator)

	count, ok := binary.Left.(*BuiltinNode)
	require.True(t, ok, "expected BuiltinNode, got %T", binary.Left)
	assert.Equal(t, "count", count.Name)
	require.NotNil(t, count.Threshold)
	assert.Equal(t, 100, *count.Threshold) // threshold = N for < operator
}

func TestOptimize_count_threshold_lte(t *testing.T) {
	tree, err := parser.Parse(`count(items, .active) <= 50`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	// Operator should remain <=, but count should have threshold set
	binary, ok := tree.Node.(*BinaryNode)
	require.True(t, ok, "expected BinaryNode, got %T", tree.Node)
	assert.Equal(t, "<=", binary.Operator)

	count, ok := binary.Left.(*BuiltinNode)
	require.True(t, ok, "expected BuiltinNode, got %T", binary.Left)
	assert.Equal(t, "count", count.Name)
	require.NotNil(t, count.Threshold)
	assert.Equal(t, 51, *count.Threshold) // threshold = N + 1 for <= operator
}

func TestOptimize_count_threshold_correctness(t *testing.T) {
	tests := []struct {
		expr string
		want bool
	}{
		// count > N (threshold = N + 1)
		{`count(1..1000, # <= 100) > 50`, true},   // 100 matches > 50
		{`count(1..1000, # <= 100) > 100`, false}, // 100 matches not > 100
		{`count(1..1000, # <= 100) > 99`, true},   // 100 matches > 99
		{`count(1..100, # > 0) > 50`, true},       // 100 matches > 50
		{`count(1..100, # > 0) > 100`, false},     // 100 matches not > 100

		// count >= N (threshold = N)
		{`count(1..1000, # <= 100) >= 100`, true},  // 100 matches >= 100
		{`count(1..1000, # <= 100) >= 101`, false}, // 100 matches not >= 101
		{`count(1..100, # > 0) >= 50`, true},       // 100 matches >= 50
		{`count(1..100, # > 0) >= 100`, true},      // 100 matches >= 100

		// count < N (threshold = N)
		{`count(1..1000, # <= 100) < 101`, true},  // 100 matches < 101
		{`count(1..1000, # <= 100) < 100`, false}, // 100 matches not < 100
		{`count(1..1000, # <= 100) < 50`, false},  // 100 matches not < 50
		{`count(1..100, # > 0) < 101`, true},      // 100 matches < 101
		{`count(1..100, # > 0) < 100`, false},     // 100 matches not < 100

		// count <= N (threshold = N + 1)
		{`count(1..1000, # <= 100) <= 100`, true},  // 100 matches <= 100
		{`count(1..1000, # <= 100) <= 99`, false},  // 100 matches not <= 99
		{`count(1..1000, # <= 100) <= 50`, false},  // 100 matches not <= 50
		{`count(1..100, # > 0) <= 100`, true},      // 100 matches <= 100
		{`count(1..100, # > 0) <= 99`, false},      // 100 matches not <= 99
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

func TestOptimize_count_threshold_no_optimization(t *testing.T) {
	// These should NOT get a threshold (handled by count_any or not optimizable)
	tests := []struct {
		code      string
		threshold bool
	}{
		{`count(items, .active) > 0`, false},   // handled by count_any
		{`count(items, .active) >= 1`, false},  // handled by count_any
		{`count(items, .active) < 1`, false},   // threshold = 1, skipped
		{`count(items, .active) <= 0`, false},  // threshold = 1, skipped
		{`count(items, .active) == 10`, false}, // not supported
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			tree, err := parser.Parse(tt.code)
			require.NoError(t, err)

			err = optimizer.Optimize(&tree.Node, nil)
			require.NoError(t, err)

			// Check if count has threshold set
			var count *BuiltinNode
			if binary, ok := tree.Node.(*BinaryNode); ok {
				count, _ = binary.Left.(*BuiltinNode)
			} else if builtin, ok := tree.Node.(*BuiltinNode); ok {
				count = builtin
			}

			if count != nil && count.Name == "count" {
				if tt.threshold {
					assert.NotNil(t, count.Threshold, "expected threshold to be set")
				} else {
					assert.Nil(t, count.Threshold, "expected threshold to be nil")
				}
			}
		})
	}
}

// Benchmark: count > 100 with early match (element 101 matches early)
func BenchmarkCountThresholdEarlyMatch(b *testing.B) {
	// Array of 10000 elements, all match predicate, threshold is 101
	// Should exit after ~101 iterations
	program, _ := expr.Compile(`count(1..10000, # > 0) > 100`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmark: count >= 50 with early match
func BenchmarkCountThresholdGteEarlyMatch(b *testing.B) {
	// All elements match, threshold is 50
	// Should exit after ~50 iterations
	program, _ := expr.Compile(`count(1..10000, # > 0) >= 50`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmark: count > 100 with no early exit (not enough matches)
func BenchmarkCountThresholdNoEarlyExit(b *testing.B) {
	// Only 100 elements match (# <= 100), threshold is 101
	// Must scan entire array
	program, _ := expr.Compile(`count(1..10000, # <= 100) > 100`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmark: Large threshold with early match
func BenchmarkCountThresholdLargeEarlyMatch(b *testing.B) {
	// All 10000 match, threshold is 1000
	// Should exit after ~1000 iterations
	program, _ := expr.Compile(`count(1..10000, # > 0) > 999`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmark: count < N with early exit (result is false)
func BenchmarkCountThresholdLtEarlyExit(b *testing.B) {
	// All 10000 match, threshold is 100
	// Should exit after ~100 iterations with result = false
	program, _ := expr.Compile(`count(1..10000, # > 0) < 100`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmark: count <= N with early exit (result is false)
func BenchmarkCountThresholdLteEarlyExit(b *testing.B) {
	// All 10000 match, threshold is 51
	// Should exit after ~51 iterations with result = false
	program, _ := expr.Compile(`count(1..10000, # > 0) <= 50`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmark: count < N without early exit (result is true)
func BenchmarkCountThresholdLtNoEarlyExit(b *testing.B) {
	// Only 100 elements match (# <= 100), threshold is 200
	// Must scan entire array, result = true
	program, _ := expr.Compile(`count(1..10000, # <= 100) < 200`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmark: count <= N without early exit (result is true)
func BenchmarkCountThresholdLteNoEarlyExit(b *testing.B) {
	// Only 100 elements match (# <= 100), threshold is 101
	// Must scan entire array, result = true
	program, _ := expr.Compile(`count(1..10000, # <= 100) <= 100`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}
