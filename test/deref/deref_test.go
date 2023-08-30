package deref_test

import (
	"context"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/require"
)

func TestDeref_binary(t *testing.T) {
	i := 1
	env := map[string]any{
		"i": &i,
		"obj": map[string]any{
			"i": &i,
		},
	}
	t.Run("==", func(t *testing.T) {
		program, err := expr.Compile(`i == 1 && obj.i == 1`, expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, true, out)
	})
	t.Run("><", func(t *testing.T) {
		program, err := expr.Compile(`i > 0 && obj.i < 99`, expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, true, out)
	})
	t.Run("??+", func(t *testing.T) {
		program, err := expr.Compile(`(i ?? obj.i) + 1`, expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, 2, out)
	})
}

func TestDeref_unary(t *testing.T) {
	i := 1
	ok := true
	env := map[string]any{
		"i": &i,
		"obj": map[string]any{
			"ok": &ok,
		},
	}

	program, err := expr.Compile(`-i < 0 && !!obj.ok`, expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, out)
}

func TestDeref_eval(t *testing.T) {
	i := 1
	env := map[string]any{
		"i": &i,
		"obj": map[string]any{
			"i": &i,
		},
	}
	out, err := expr.Eval(`i == 1 && obj.i == 1`, env)
	require.NoError(t, err)
	require.Equal(t, true, out)
}

func TestDeref_emptyCtx(t *testing.T) {
	program, err := expr.Compile(`ctx`)
	require.NoError(t, err)

	output, err := expr.Run(program, map[string]any{
		"ctx": context.Background(),
	})
	require.NoError(t, err)
	require.Implements(t, new(context.Context), output)
}

func TestDeref_emptyCtx_Eval(t *testing.T) {
	output, err := expr.Eval(`ctx`, map[string]any{
		"ctx": context.Background(),
	})
	require.NoError(t, err)
	require.Implements(t, new(context.Context), output)
}

func TestDeref_context_WithValue(t *testing.T) {
	program, err := expr.Compile(`ctxWithValue`)
	require.NoError(t, err)

	output, err := expr.Run(program, map[string]any{
		"ctxWithValue": context.WithValue(context.Background(), "value", "test"),
	})
	require.NoError(t, err)
	require.Implements(t, new(context.Context), output)
}

func TestDeref_method_on_int_pointer(t *testing.T) {
	output, err := expr.Eval(`foo.Bar()`, map[string]any{
		"foo": new(foo),
	})
	require.NoError(t, err)
	require.Equal(t, 42, output)
}

type foo int

func (f *foo) Bar() int {
	return 42
}

func TestDeref_multiple_pointers(t *testing.T) {
	a := 42
	b := &a
	c := &b
	t.Run("returned as is", func(t *testing.T) {
		output, err := expr.Eval(`c`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, c, output)
		require.IsType(t, (**int)(nil), output)
	})
	t.Run("+ works", func(t *testing.T) {
		output, err := expr.Eval(`c+2`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, 44, output)
	})
}
