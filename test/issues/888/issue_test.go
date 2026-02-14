package main

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type Container struct {
	ID   string
	List []string
}

func (c Container) IncludesAny(s ...string) bool {
	for _, l := range c.List {
		// Note: original issue used "slices.Contains" but
		// it is not available in the minimum Go version of expr (1.18).
		for _, v := range s {
			if v == l {
				return true
			}
		}
	}
	return false
}

func TestIssue888(t *testing.T) {
	env := map[string]any{
		"Container": Container{
			ID:   "id",
			List: []string{"foo", "bar", "baz"},
		},
	}

	code := `Container.IncludesAny("nope", "nope again", "bar")`

	program, err := expr.Compile(code, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}
