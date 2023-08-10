package deref_test

import (
	"context"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/require"
)

func TestDeref_binary(t *testing.T) {
	i := 1
	env := map[string]interface{}{
		"i": &i,
		"map": map[string]interface{}{
			"i": &i,
		},
	}
	t.Run("==", func(t *testing.T) {
		// With specified env, OpDeref added and == works as expected.
		program, err := expr.Compile(`i == 1 && map.i == 1`, expr.Env(env))
		require.NoError(t, err)

		env["any"] = &i
		out, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, true, out)
	})
	t.Run(">", func(t *testing.T) {
		// With specified env, OpDeref added and == works as expected.
		program, err := expr.Compile(`i > 1 && map.i > 1`, expr.Env(env))
		require.NoError(t, err)

		env["any"] = &i
		out, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, true, out)
	})
}

func TestDeref_eval(t *testing.T) {
	i := 1
	env := map[string]interface{}{
		"i": &i,
		"map": map[string]interface{}{
			"i": &i,
		},
	}
	out, err := expr.Eval(`i == 1 && map.i == 1`, env)
	require.NoError(t, err)
	require.Equal(t, true, out)
}

func TestDeref_emptyCtx(t *testing.T) {
	program, err := expr.Compile(`ctx`)
	require.NoError(t, err)

	output, err := expr.Run(program, map[string]interface{}{
		"ctx": context.Background(),
	})
	require.NoError(t, err)
	require.Implements(t, new(context.Context), output)
}

func TestDeref_emptyCtx_Eval(t *testing.T) {
	output, err := expr.Eval(`ctx`, map[string]interface{}{
		"ctx": context.Background(),
	})
	require.NoError(t, err)
	require.Implements(t, new(context.Context), output)
}

func TestDeref_context_WithValue(t *testing.T) {
	program, err := expr.Compile(`ctxWithValue`)
	require.NoError(t, err)

	output, err := expr.Run(program, map[string]interface{}{
		"ctxWithValue": context.WithValue(context.Background(), "value", "test"),
	})
	require.NoError(t, err)
	require.Implements(t, new(context.Context), output)
}

func TestDeref_method_on_int_pointer(t *testing.T) {
	output, err := expr.Eval(`foo.Bar()`, map[string]interface{}{
		"foo": new(foo),
	})
	require.NoError(t, err)
	require.Equal(t, 42, output)
}

type foo int

func (f *foo) Bar() int {
	return 42
}
