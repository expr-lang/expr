package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type Foo interface {
	Add(a int, b *int) int
}

type FooImpl struct {
}

func (*FooImpl) Add(a int, b *int) int {
	return 0
}

type Env struct {
	Foo Foo `expr:"foo"`
}

func (Env) Any(x any) any {
	return x
}

func TestNoInterfaceMethodWithNil(t *testing.T) {
	program, err := expr.Compile(`foo.Add(1, nil)`)
	require.NoError(t, err)

	_, err = expr.Run(program, Env{Foo: &FooImpl{}})
	require.NoError(t, err)
}

func TestNoInterfaceMethodWithNil_with_env(t *testing.T) {
	program, err := expr.Compile(`foo.Add(1, nil)`, expr.Env(Env{}))
	require.NoError(t, err)

	_, err = expr.Run(program, Env{Foo: &FooImpl{}})
	require.NoError(t, err)
}

func TestNoInterfaceMethodWithNil_with_any(t *testing.T) {
	program, err := expr.Compile(`Any(nil)`, expr.Env(Env{}))
	require.NoError(t, err)

	out, err := expr.Run(program, Env{Foo: &FooImpl{}})
	require.NoError(t, err)
	require.Equal(t, nil, out)
}
