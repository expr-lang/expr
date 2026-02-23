package issue934

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type Env struct {
	Exported
	unexported
}

type Exported struct {
	Str string
	str string
}

type unexported struct {
	Integer int
	integer int
}

// TestIssue934 tests that accessing unexported fields on values whose type
// is unknown at compile time (e.g., from a ternary with mixed types) yields
// a descriptive error.
//
// OSS-Fuzz issue #486370271.
func TestIssue934(t *testing.T) {
	env := map[string]any{
		"v": Env{},
	}

	// Accessing unexported fields like time.Time.loc
        // should produce a proper error.
	_, err := expr.Eval(`(true ? v : "string").str`, env)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot fetch str")

	// Exported field access should still work.
	_, err = expr.Eval(`(true ? v : "string").Str`, env)
	require.NoError(t, err)

	// Access unexported field inherited from unexported embedded struct.
	_, err = expr.Eval(`(true ? v : "string").integer`, env)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot fetch integer")

	// Access exported field inherited from unexported embedded struct should work.
	_, err = expr.Eval(`(true ? v : "string").Integer`, env)
	require.NoError(t, err)
}
