package patcher_test

import (
	"context"
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/patcher"
)

func TestWithContext(t *testing.T) {
	env := map[string]any{
		"fn": func(ctx context.Context, a int) int {
			return ctx.Value("value").(int) + a
		},
		"ctx": context.TODO(),
	}

	withContext := patcher.WithContext{Name: "ctx"}

	program, err := expr.Compile(`fn(40)`, expr.Env(env), expr.Patch(withContext))
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), "value", 2)
	env["ctx"] = ctx

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 42, output)
}

func TestWithContext_with_env_Function(t *testing.T) {
	env := map[string]any{
		"ctx": context.TODO(),
	}

	fn := expr.Function("fn",
		func(params ...any) (any, error) {
			ctx := params[0].(context.Context)
			a := params[1].(int)

			return ctx.Value("value").(int) + a, nil
		},
		new(func(context.Context, int) int),
	)

	program, err := expr.Compile(
		`fn(40)`,
		expr.Env(env),
		expr.WithContext("ctx"),
		fn,
	)
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), "value", 2)
	env["ctx"] = ctx

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 42, output)
}

type testEnvContext struct {
	Context context.Context `expr:"ctx"`
}

func (testEnvContext) Fn(ctx context.Context, a int) int {
	return ctx.Value("value").(int) + a
}

func TestWithContext_env_struct(t *testing.T) {
	withContext := patcher.WithContext{Name: "ctx"}

	program, err := expr.Compile(`Fn(40)`, expr.Env(testEnvContext{}), expr.Patch(withContext))
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), "value", 2)
	env := testEnvContext{
		Context: ctx,
	}

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 42, output)
}

type TestFoo struct {
	contextValue int
}

func (f *TestFoo) GetValue(a int) int64 {
	return int64(f.contextValue + a)
}

func TestWithContext_with_env_method_chain(t *testing.T) {
	env := map[string]any{
		"ctx": context.TODO(),
	}

	fn := expr.Function("fn",
		func(params ...any) (any, error) {
			ctx := params[0].(context.Context)

			contextValue := ctx.Value("value").(int)

			return &TestFoo{
				contextValue: contextValue,
			}, nil
		},
		new(func(context.Context) *TestFoo),
	)

	program, err := expr.Compile(
		`fn().GetValue(40)`,
		expr.Env(env),
		expr.WithContext("ctx"),
		fn,
		expr.AsInt64(),
	)
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), "value", 2)
	env["ctx"] = ctx

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, int64(42), output)
}

func TestWithContext_issue529(t *testing.T) {
	env := map[string]any{
		"ctx": context.Background(),
		"foo": func(ctx context.Context, n int) int {
			if ctx == nil {
				panic("wanted a context")
			}
			return n + 1
		},
	}
	options := []expr.Option{
		expr.Env(env),
		expr.WithContext("ctx"),
	}

	code := "foo(0) | foo()"
	program, err := expr.Compile(code, options...)
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 2, out)
}
