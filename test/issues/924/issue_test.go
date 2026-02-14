package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue924_allow_disabling_builtins_and_providing_fn_at_runtime(t *testing.T) {
	// We disable the builtin "upper", but do not env information,
	// but we can provide a function at runtime.
	program, err := expr.Compile(`upper(1)`, expr.DisableBuiltin("upper"))
	require.NoError(t, err)

	env := map[string]any{
		"upper": func(a int) int { return a },
	}

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 1, out)
}
