package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue710_disable_env_builtin(t *testing.T) {
	env := map[string]any{
		"foo": 42,
	}

	t.Run("$env disabled in strict mode", func(t *testing.T) {
		_, err := expr.Compile(`$env`, expr.Env(env), expr.DisableBuiltin("$env"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown name $env")
	})

	t.Run("$env.field disabled in strict mode", func(t *testing.T) {
		_, err := expr.Compile(`$env.foo`, expr.Env(env), expr.DisableBuiltin("$env"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown name $env")
	})

	t.Run("get($env, field) disabled in strict mode", func(t *testing.T) {
		_, err := expr.Compile(`get($env, "foo")`, expr.Env(env), expr.DisableBuiltin("$env"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown name $env")
	})

	t.Run("$env() disabled in strict mode", func(t *testing.T) {
		_, err := expr.Compile(`$env()`, expr.Env(env), expr.DisableBuiltin("$env"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown name $env")
	})

	t.Run("$env disabled without strict mode", func(t *testing.T) {
		// Without expr.Env(), the checker runs in non-strict mode.
		// Disabled $env falls through to normal identifier resolution,
		// which returns nil for unknown names in non-strict mode.
		program, err := expr.Compile(`$env`, expr.DisableBuiltin("$env"))
		require.NoError(t, err)

		out, err := expr.Run(program, map[string]any{})
		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("$env enabled by default", func(t *testing.T) {
		program, err := expr.Compile(`$env.foo`, expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, 42, out)
	})

	t.Run("$env standalone enabled by default", func(t *testing.T) {
		program, err := expr.Compile(`$env`, expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, env, out)
	})
}
