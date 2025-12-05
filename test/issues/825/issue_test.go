package main

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue825(t *testing.T) {
	m := map[string]any{"foo": 42}
	env := &m

	prog, err := expr.Compile("foo > 0", expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(prog, env)
	require.NoError(t, err)
	require.Equal(t, true, out)
}
