package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestEnvFieldMethods(t *testing.T) {
	program, err := expr.Compile(`Func(0)`, expr.Env(&Env{}))
	require.NoError(t, err)

	env := &Env{}
	env.Func = func() int {
		return 42
	}

	out, err := expr.Run(program, Env{})
	require.NoError(t, err)

	require.Equal(t, 42, out)
}

type Env struct {
	EmbeddedEnv
	Func func() int
}

type EmbeddedEnv struct{}
