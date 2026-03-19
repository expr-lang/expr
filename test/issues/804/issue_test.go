package issue_test

import (
	"math"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

// TestIssue804 verifies that arithmetic operations on int64 and uint64
// operands preserve 64-bit precision instead of truncating to int
// (which is 32-bit on 32-bit architectures).
func TestIssue804(t *testing.T) {
	env := map[string]any{
		"a": int64(math.MaxInt32),
		"b": int64(1),
	}

	// int64 + int64 must produce int64, not int.
	out, err := expr.Eval("a + b", env)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt32)+1, out)

	// int64 - int64
	out, err = expr.Eval("a - b", env)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt32)-1, out)

	// int64 * int64
	out, err = expr.Eval("a * b", env)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt32), out)
}

func TestIssue804_uint64(t *testing.T) {
	env := map[string]any{
		"a": uint64(math.MaxUint32),
		"b": uint64(1),
	}

	// uint64 + uint64 must produce uint64, not int.
	out, err := expr.Eval("a + b", env)
	require.NoError(t, err)
	require.Equal(t, uint64(math.MaxUint32)+1, out)
}

func TestIssue804_mixed(t *testing.T) {
	env := map[string]any{
		"a": int32(100),
		"b": int64(math.MaxInt32),
	}

	// int32 + int64 must produce int64.
	out, err := expr.Eval("a + b", env)
	require.NoError(t, err)
	require.Equal(t, int64(100)+int64(math.MaxInt32), out)
}
