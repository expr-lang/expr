package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type InputStruct struct {
	Enabled *bool `json:"enabled"`
}

func TestIssue836(t *testing.T) {
	str := "foo"
	ptrStr := &str
	b := true
	ptrBool := &b
	i := 1
	ptrInt := &i

	env := map[string]interface{}{
		"ptrStr":  ptrStr,
		"ptrBool": ptrBool,
		"ptrInt":  ptrInt,
		"arr":     []int{1, 2, 3},
		"mapPtr":  map[*int]int{ptrInt: 42},
	}

	t.Run("map access with pointer key", func(t *testing.T) {
		program, err := expr.Compile(`{"foo": "bar"}[ptrStr]`, expr.Env(env))
		require.NoError(t, err)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, "bar", output)
	})

	t.Run("conditional with pointer condition", func(t *testing.T) {
		program, err := expr.Compile(`ptrBool ? 1 : 0`, expr.Env(env))
		require.NoError(t, err)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, 1, output)
	})

	t.Run("get() with pointer key", func(t *testing.T) {
		program, err := expr.Compile(`get({"foo": "bar"}, ptrStr)`, expr.Env(env))
		require.NoError(t, err)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, "bar", output)
	})

	t.Run("struct field pointer check in ternary", func(t *testing.T) {
		var v InputStruct
		// v.Enabled is nil

		env := map[string]any{
			"v": v,
		}

		code := `v.Enabled == nil ? 'default' : ( v.Enabled ? 'enabled' : 'disabled' )`

		program, err := expr.Compile(code, expr.Env(env))
		require.NoError(t, err)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, "default", output)
	})

	t.Run("struct field pointer check in ternary (enabled)", func(t *testing.T) {
		b := true
		v := InputStruct{Enabled: &b}

		env := map[string]any{
			"v": v,
		}

		code := `v.Enabled == nil ? 'default' : ( v.Enabled ? 'enabled' : 'disabled' )`

		program, err := expr.Compile(code, expr.Env(env))
		require.NoError(t, err)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, "enabled", output)
	})

	t.Run("slice with pointer indices", func(t *testing.T) {
		program, err := expr.Compile(`arr[ptrInt:ptrInt]`, expr.Env(env))
		require.NoError(t, err)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, []int{}, output)
	})
}
