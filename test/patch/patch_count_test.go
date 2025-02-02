package patch_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/test/mock"
)

// This patcher tracks how many nodes it patches which can 
// be used to verify if it was run too many times or not at all
type countingPatcher struct {
	PatchCount int
}

func (c *countingPatcher) Visit(node *ast.Node) {
	switch (*node).(type) {
	case *ast.IntegerNode:
		c.PatchCount++
	}
}

// Test over a simple expression
func TestPatch_Count(t *testing.T) {
	patcher := countingPatcher{}

	_, err := expr.Compile(
		`5 + 5`,
		expr.Env(mock.Env{}),
		expr.Patch(&patcher),
	)
	require.NoError(t, err)

	require.Equal(t, 2, patcher.PatchCount, "Patcher run an unexpected number of times during compile")
}

// Test with operator overloading
func TestPatchOperator_Count(t *testing.T) {
	patcher := countingPatcher{}

	_, err := expr.Compile(
		`5 + 5`,
		expr.Env(mock.Env{}),
		expr.Patch(&patcher),
		expr.Operator("+", "_intAdd"),
		expr.Function(
			"_intAdd",
			func(params ...any) (any, error) {
				return params[0].(int) + params[1].(int), nil
			},
			new(func(int, int) int),
		),
	)

	require.NoError(t, err)

	require.Equal(t, 2, patcher.PatchCount, "Patcher run an unexpected number of times during compile")
}
