package issue_test

import (
	"context"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type X struct{}

func (x *X) HelloCtx(ctx context.Context, text string) error {
	return nil
}

func TestIssue756(t *testing.T) {
	env := map[string]any{
		"_goctx_": context.TODO(),
		"_g_": map[string]*X{
			"rpc": {},
		},
		"text": "еуче",
	}
	exprStr := `let v = _g_.rpc.HelloCtx(text); v`
	program, err := expr.Compile(exprStr, expr.Env(env), expr.WithContext("_goctx_"))
	require.NoError(t, err)

	_, err = expr.Run(program, env)
	require.NoError(t, err)
}
