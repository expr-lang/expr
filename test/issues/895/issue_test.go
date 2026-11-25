package main

import (
	"math"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue895(t *testing.T) {
	env := map[string]any{
		"a": map[string]any{"a": 1, "b": 2},
		"b": map[string]any{"b": 3, "c": 4},
	}

	program, err := expr.Compile(`merge(a, b)`, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, map[any]any{"a": 1, "b": 3, "c": 4}, output)
}

func TestIssue895_does_not_modify_input(t *testing.T) {
	a := map[string]any{"a": 1}
	b := map[string]any{"b": 2}
	env := map[string]any{
		"a": a,
		"b": b,
	}

	program, err := expr.Compile(`merge(a, b)`, expr.Env(env))
	require.NoError(t, err)

	_, err = expr.Run(program, env)
	require.NoError(t, err)

	// Original maps must be unmodified.
	require.Equal(t, map[string]any{"a": 1}, a)
	require.Equal(t, map[string]any{"b": 2}, b)
}

func TestIssue895_preserves_key_types(t *testing.T) {
	env := map[string]any{
		"a": map[int]string{1: "int"},
		"b": map[string]string{"1": "string"},
	}

	program, err := expr.Compile(`merge(a, b)`, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, map[any]any{1: "int", "1": "string"}, output)
}

func TestIssue895_keeps_non_reflexive_keys(t *testing.T) {
	env := map[string]any{
		"a": map[float64]string{math.NaN(): "nan", 1: "one"},
		"b": map[float64]string{2: "two"},
	}

	program, err := expr.Compile(`merge(a, b)`, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)

	merged := output.(map[any]any)
	require.Len(t, merged, 3)
	require.Equal(t, "one", merged[1.0])
	require.Equal(t, "two", merged[2.0])
}
