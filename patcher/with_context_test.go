package patcher_test

import (
	"context"
)

type testEnvContext struct {
	Context context.Context `expr:"ctx"`
}

func (testEnvContext) Fn(ctx context.Context, a int) int {
	return ctx.Value("value").(int) + a
}

type TestFoo struct {
	contextValue int
}

func (f *TestFoo) GetValue(a int) int64 {
	return int64(f.contextValue + a)
}
