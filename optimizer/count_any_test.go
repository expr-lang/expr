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

func TestOptimize_count_any(t *testing.T) {
	tree, err := parser.Parse(`count(items, .active) > 0`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &BuiltinNode{
		Name: "any",
		Arguments: []Node{
			&IdentifierNode{Value: "items"},
			&PredicateNode{
				Node: &MemberNode{
					Node:     &PointerNode{},
					Property: &StringNode{Value: "active"},
				},
			},
		},
	}

	assert.Equal(t, Dump(expected), Dump(tree.Node))
}

func TestOptimize_count_any_gte_one(t *testing.T) {
	tree, err := parser.Parse(`count(items, .valid) >= 1`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &BuiltinNode{
		Name: "any",
		Arguments: []Node{
			&IdentifierNode{Value: "items"},
			&PredicateNode{
				Node: &MemberNode{
					Node:     &PointerNode{},
					Property: &StringNode{Value: "valid"},
				},
			},
		},
	}

	assert.Equal(t, Dump(expected), Dump(tree.Node))
}

func TestOptimize_count_any_correctness(t *testing.T) {
	tests := []struct {
		expr string
		want bool
	}{
		// count > 0 → any
		{`count(1..100, # == 1) > 0`, true},
		{`count(1..100, # == 50) > 0`, true},
		{`count(1..100, # == 100) > 0`, true},
		{`count(1..100, # == 0) > 0`, false},

		// count >= 1 → any
		{`count(1..100, # % 10 == 0) >= 1`, true},
		{`count(1..100, # > 100) >= 1`, false},
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

func TestOptimize_count_no_optimization(t *testing.T) {
	// These should NOT be optimized
	tests := []string{
		`count(items, .active) > 1`,  // not > 0
		`count(items, .active) >= 2`, // not >= 1
		`count(items, .active) == 0`, // not optimized (none has overhead)
		`count(items, .active) == 1`, // not == 0
		`count(items, .active) < 1`,  // not optimized (none has overhead)
		`count(items, .active) <= 0`, // not optimized (none has overhead)
		`count(items, .active) != 0`, // different operator
	}

	for _, code := range tests {
		t.Run(code, func(t *testing.T) {
			tree, err := parser.Parse(code)
			require.NoError(t, err)

			err = optimizer.Optimize(&tree.Node, nil)
			require.NoError(t, err)

			// Should still be a BinaryNode (not optimized to any)
			_, ok := tree.Node.(*BinaryNode)
			assert.True(t, ok, "expected BinaryNode, got %T", tree.Node)
		})
	}
}

// Benchmarks for count > 0 → any
func BenchmarkCountGtZero(b *testing.B) {
	program, _ := expr.Compile(`count(1..1000, # == 1) > 0`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

func BenchmarkCountGtZeroLargeEarlyMatch(b *testing.B) {
	program, _ := expr.Compile(`count(1..10000, # == 1) > 0`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

func BenchmarkCountGtZeroNoMatch(b *testing.B) {
	program, _ := expr.Compile(`count(1..1000, # == 0) > 0`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

// Benchmarks for count >= 1 → any
func BenchmarkCountGteOneEarlyMatch(b *testing.B) {
	program, _ := expr.Compile(`count(1..1000, # == 1) >= 1`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}

func BenchmarkCountGteOneNoMatch(b *testing.B) {
	program, _ := expr.Compile(`count(1..1000, # == 0) >= 1`)
	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = vm.Run(program, nil)
	}
	_ = out
}
