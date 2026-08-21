package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

// TestIssue825 verifies that passing a pointer to a map as the environment is
// dereferenced instead of panicking.
//
// conf.EnvWithCache selected the map branch using the dereferenced value's
// kind, but then read the map keys/length from the original (pointer) value,
// panicking with:
//
//	reflect: call of reflect.Value.Len on ptr to non-array Value
func TestIssue825(t *testing.T) {
	m := map[string]any{"foo": 42}

	program, err := expr.Compile("foo + 1", expr.Env(&m))
	require.NoError(t, err)

	out, err := expr.Run(program, m)
	require.NoError(t, err)
	assert.Equal(t, 43, out)
}

// TestIssue825_Strict verifies that a pointer-to-map env keeps the strict-mode
// and element-type information of the dereferenced map, i.e. it behaves exactly
// like compiling with the map value itself.
func TestIssue825_Strict(t *testing.T) {
	m := map[string]int{"a": 1}

	// Unknown names are rejected (strict), just like a plain map env.
	_, err := expr.Compile("unknown + 1", expr.Env(&m))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown name unknown")

	// The element type (int) is inferred from the dereferenced map.
	program, err := expr.Compile("a + 1", expr.Env(&m))
	require.NoError(t, err)

	out, err := expr.Run(program, m)
	require.NoError(t, err)
	assert.Equal(t, 2, out)
}
