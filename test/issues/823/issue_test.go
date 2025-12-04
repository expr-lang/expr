package issue_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type env struct {
	Ctx context.Context `expr:"ctx"`
}

func TestIssue823(t *testing.T) {
	p, err := expr.Compile(
		"now2().After(date2())",
		expr.Env(env{}),
		expr.WithContext("ctx"),
		expr.Function(
			"now2",
			func(params ...any) (any, error) { return time.Now(), nil },
			new(func(context.Context) time.Time),
		),
		expr.Function(
			"date2",
			func(params ...any) (any, error) { return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil },
			new(func(context.Context) time.Time),
		),
	)
	fmt.Printf("Compile result err: %v\n", err)
	require.NoError(t, err)

	r, err := expr.Run(p, &env{Ctx: context.Background()})
	require.NoError(t, err)
	require.True(t, r.(bool))
}
