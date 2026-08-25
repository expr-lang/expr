package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

// TestIssue825 verifies that passing a pointer to a map as the environment is
// rejected with a clear message instead of panicking deep inside reflect.
//
// Supporting *map by dereferencing it was declined by the maintainer (#825 is
// labeled wontfix); the agreed direction was to reject *map with an error
// message. Previously conf.EnvWithCache selected the map branch on the
// dereferenced kind but then read the map keys/length from the original
// (pointer) value, panicking with the opaque:
//
//	reflect: call of reflect.Value.Len on ptr to non-array Value
func TestIssue825(t *testing.T) {
	m := map[string]any{"foo": 42}

	assert.PanicsWithValue(t,
		"environment must be a map, not a pointer to a map: *map[string]interface {}",
		func() {
			_, _ = expr.Compile("foo > 0", expr.Env(&m))
		},
	)
}

// TestIssue825_NilPointer verifies that a nil pointer to a map is rejected with
// the same clear message as a non-nil one, rather than falling through to the
// generic "unknown type" panic (the pointer-to-map is detected by type).
func TestIssue825_NilPointer(t *testing.T) {
	var m *map[string]any

	assert.PanicsWithValue(t,
		"environment must be a map, not a pointer to a map: *map[string]interface {}",
		func() {
			_, _ = expr.Compile("foo > 0", expr.Env(m))
		},
	)
}

// TestIssue825_MapStillWorks guards the common case: a map passed by value is
// unaffected and continues to work exactly as before.
func TestIssue825_MapStillWorks(t *testing.T) {
	m := map[string]any{"foo": 42}

	program, err := expr.Compile("foo + 1", expr.Env(m))
	require.NoError(t, err)

	out, err := expr.Run(program, m)
	require.NoError(t, err)
	assert.Equal(t, 43, out)
}
