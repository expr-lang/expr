package deref_test

import (
	"context"
	"testing"
	"time"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
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

func TestDeref_pointer_of_interface(t *testing.T) {
	v := 42
	a := &v
	b := any(a)
	c := any(&b)
	t.Run("returned as is", func(t *testing.T) {
		output, err := expr.Eval(`c`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, c, output)
		require.IsType(t, (*interface{})(nil), output)
	})
	t.Run("+ works", func(t *testing.T) {
		output, err := expr.Eval(`c+2`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, 44, output)
	})
}

func TestDeref_nil(t *testing.T) {
	var b *int = nil
	c := &b
	t.Run("returned as is", func(t *testing.T) {
		output, err := expr.Eval(`c`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, c, output)
		require.IsType(t, (**int)(nil), output)
	})
	t.Run("== nil works", func(t *testing.T) {
		output, err := expr.Eval(`c == nil`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, true, output)
	})
}

func TestDeref_nil_in_pointer_of_interface(t *testing.T) {
	var a *int32 = nil
	b := any(a)
	c := any(&b)
	t.Run("returned as is", func(t *testing.T) {
		output, err := expr.Eval(`c`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, c, output)
		require.IsType(t, (*interface{})(nil), output)
	})
	t.Run("== nil works", func(t *testing.T) {
		output, err := expr.Eval(`c == nil`, map[string]any{
			"c": c,
		})
		require.NoError(t, err)
		require.Equal(t, true, output)
	})
}

func TestDeref_commutative(t *testing.T) {
	a := "ok"
	b := "ok"

	type Env struct {
		A string
		B *string
	}

	env := Env{
		A: a,
		B: &b,
	}

	tests := []struct {
		code string
		want bool
	}{
		{`A == B`, true},
		{`B == A`, true},
		{`A != B`, false},
		{`B != A`, false},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			program, err := expr.Compile(test.code, expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, test.want, out)
		})
	}
}

func TestDeref_fetch_from_interface_mix_pointer(t *testing.T) {
	type FooBar struct {
		Value string
	}
	foobar := &FooBar{"waldo"}
	var foobarAny any = foobar
	var foobarPtrAny any = &foobarAny

	res, err := expr.Eval("foo.Value", map[string]any{
		"foo": foobarPtrAny,
	})
	assert.NoError(t, err)
	assert.Equal(t, "waldo", res)
}

func TestDeref_func_args(t *testing.T) {
	i := 20
	env := map[string]any{
		"var": &i,
		"fn": func(p int) int {
			return p + 1
		},
	}

	program, err := expr.Compile(`fn(var) + fn(var + 0)`, expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 42, out)
}

func TestDeref_struct_func_args(t *testing.T) {
	n, _ := time.Parse(time.RFC3339, "2024-05-12T18:30:00+00:00")
	duration := 30 * time.Minute
	env := map[string]any{
		"time":     n,
		"duration": &duration,
	}

	program, err := expr.Compile(`time.Add(duration).Format('2006-01-02T15:04:05Z07:00')`, expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, "2024-05-12T19:00:00Z", out)
}

func TestDeref_ignore_func_args(t *testing.T) {
	f := foo(1)
	env := map[string]any{
		"foo": &f,
		"fn": func(f *foo) int {
			return f.Bar()
		},
	}

	program, err := expr.Compile(`fn(foo)`, expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 42, out)
}

func TestDeref_ignore_struct_func_args(t *testing.T) {
	n := time.Now()
	location, _ := time.LoadLocation("UTC")
	env := map[string]any{
		"time":     n,
		"location": location,
	}

	program, err := expr.Compile(`time.In(location).Location().String()`, expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, "UTC", out)
}
