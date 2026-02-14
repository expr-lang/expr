package issue_test

import (
	"context"
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type env struct {
	Ctx context.Context `expr:"ctx"`
}

// TestIssue823 verifies that WithContext injects context into nested custom
// function calls. The bug was that date2() nested as an argument to After()
// didn't receive the context because its callee type was unknown.
func TestIssue823(t *testing.T) {
	now2Called := false
	date2Called := false

	p, err := expr.Compile(
		"now2().After(date2())",
		expr.Env(env{}),
		expr.WithContext("ctx"),
		expr.Function(
			"now2",
			func(params ...any) (any, error) {
				require.Len(t, params, 1, "now2 should receive context")
				_, ok := params[0].(context.Context)
				require.True(t, ok, "now2 first param should be context.Context")
				now2Called = true
				return time.Now(), nil
			},
			new(func(context.Context) time.Time),
		),
		expr.Function(
			"date2",
			func(params ...any) (any, error) {
				require.Len(t, params, 1, "date2 should receive context")
				_, ok := params[0].(context.Context)
				require.True(t, ok, "date2 first param should be context.Context")
				date2Called = true
				return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil
			},
			new(func(context.Context) time.Time),
		),
	)
	require.NoError(t, err)

	r, err := expr.Run(p, &env{Ctx: context.Background()})
	require.NoError(t, err)
	require.True(t, r.(bool))
	require.True(t, now2Called, "now2 should have been called")
	require.True(t, date2Called, "date2 should have been called")
}

// envWithMethods tests that Env methods with context.Context work correctly
// when nested in method chains (similar to TestIssue823 but with Env methods).
type envWithMethods struct {
	Ctx context.Context `expr:"ctx"`
}

func (e *envWithMethods) Now2(ctx context.Context) time.Time {
	return time.Now()
}

func (e *envWithMethods) Date2(ctx context.Context) time.Time {
	return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestIssue823_EnvMethods(t *testing.T) {
	p, err := expr.Compile(
		"Now2().After(Date2())",
		expr.Env(&envWithMethods{}),
		expr.WithContext("ctx"),
	)
	require.NoError(t, err)

	r, err := expr.Run(p, &envWithMethods{Ctx: context.Background()})
	require.NoError(t, err)
	require.True(t, r.(bool))
}
