package interface_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
)

type Bar interface {
	Bar() int
}

type FooImpl struct{}

func (f FooImpl) Foo() Bar {
	return BarImpl{}
}

type BarImpl struct{}

// Aba is a special method that is not part of the Bar interface,
// but is used to test that the correct method is called. "Aba" name
// is chosen to be before "Bar" in the alphabet.
func (b BarImpl) Aba() bool {
	return true
}

func (b BarImpl) Bar() int {
	return 42
}

func TestInterfaceMethod(t *testing.T) {
	require.True(t, BarImpl{}.Aba())
	require.True(t, BarImpl{}.Bar() == 42)

	env := map[string]any{
		"var": FooImpl{},
	}
	p, err := expr.Compile(`var.Foo().Bar()`, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(p, env)
	assert.NoError(t, err)
	assert.Equal(t, 42, out)
}
