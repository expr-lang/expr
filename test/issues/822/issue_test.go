package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue822(t *testing.T) {
	var tests = []struct {
		input string
		want  any
		err   string
	}{
		{
			input: `let x = nil; x?.[0:1]`,
		},
		{
			input: `let x = nil; x?.[:1]`,
		},
		{
			input: `let x = nil; x?.[1:]`,
		},
		{
			input: `let x = [1, 2, 3]; x?.[0:2]`,
			want:  []any{1, 2},
		},
		{
			input: `let x = "test"; x?.[5:10]`,
			want:  "",
		},
		{
			input: `let x = 1; x?.[0:1]`,
			err:   "cannot slice int",
		},
		{
			input: `let x = 1.5; x?.[0:1]`,
			err:   "cannot slice float64",
		},
		{
			input: `let x = true; x?.[0:1]`,
			err:   "cannot slice bool",
		},
		{
			input: `let x = {a: 1}; x?.[0:1]`,
			err:   "cannot slice map[string]interface {}",
		},
		{
			input: `let x = nil; x?.[true:false]`,
			err:   "non-integer slice index bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := expr.Compile(tt.input)

			if tt.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)
				out, err := expr.Run(program, nil)
				require.NoError(t, err)
				require.Equal(t, tt.want, out)
			}
		})
	}
}

func TestIssue822_nil_from_env(t *testing.T) {
	env := map[string]any{"a": nil}

	program, err := expr.Compile(`a?.[0:1]`, expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Nil(t, out)
}
