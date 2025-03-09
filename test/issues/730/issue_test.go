package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type ModeEnum int

const (
	ModeEnumA ModeEnum = 1
)

type Env struct {
	Mode *ModeEnum
}

func TestIssue730(t *testing.T) {
	code := `int(Mode) == 1`

	tmp := ModeEnumA

	env := map[string]any{
		"Mode": &tmp,
	}

	program, err := expr.Compile(code, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.True(t, output.(bool))
}

func TestIssue730_warn_about_different_types(t *testing.T) {
	code := `Mode == 1`

	_, err := expr.Compile(code, expr.Env(Env{}))
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid operation: == (mismatched types issue_test.ModeEnum and int)")
}

func TestIssue730_eval(t *testing.T) {
	code := `Mode == 1`

	tmp := ModeEnumA

	env := map[string]any{
		"Mode": &tmp,
	}

	// Golang also does not allow this:
	// _ = ModeEnumA == int(1) // will not compile

	out, err := expr.Eval(code, env)
	require.NoError(t, err)
	require.False(t, out.(bool))
}
