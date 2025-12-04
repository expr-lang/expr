package main

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue857(t *testing.T) {
	foo := map[string]any{
		"entry": map[string]any{
			"alpha": "x",
			"beta":  1,
		},
	}
	bar := map[string]any{
		"entry": map[string]any{
			"alpha": "x",
			"beta":  1,
		},
	}

	env := map[string]any{
		"foo": foo,
		"bar": bar,
	}

	code := `
		foo
		| keys()
		| filter(# in bar)
		| filter(foo[#].alpha == bar[#].alpha)
		| filter(foo[#].beta == bar[#].beta)
	`

	_, err := expr.Compile(code, expr.Env(env))
	require.NoError(t, err)
}
