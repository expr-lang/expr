package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue461(t *testing.T) {
	type EnvStr string
	type EnvField struct {
		S   EnvStr
		Str string
	}
	type Env struct {
		S        EnvStr
		Str      string
		EnvField EnvField
	}
	var tests = []struct {
		input string
		env   Env
		want  bool
		err   string
	}{
		{
			input: "Str == S",
			env:   Env{S: "string", Str: "string"},
			err:   "invalid operation: == (mismatched types string and issue_test.EnvStr)",
		},
		{
			input: "Str == Str",
			env:   Env{Str: "string"},
			want:  true,
		},
		{
			input: "S == S",
			env:   Env{Str: "string"},
			want:  true,
		},
		{
			input: `Str == "string"`,
			env:   Env{Str: "string"},
			want:  true,
		},
		{
			input: `S == "string"`,
			env:   Env{Str: "string"},
			err:   "invalid operation: == (mismatched types issue_test.EnvStr and string)",
		},
		{
			input: "EnvField.Str == EnvField.S",
			env:   Env{EnvField: EnvField{S: "string", Str: "string"}},
			err:   "invalid operation: == (mismatched types string and issue_test.EnvStr)",
		},
		{
			input: "EnvField.Str == EnvField.Str",
			env:   Env{EnvField: EnvField{Str: "string"}},
			want:  true,
		},
		{
			input: "EnvField.S == EnvField.S",
			env:   Env{EnvField: EnvField{Str: "string"}},
			want:  true,
		},
		{
			input: `EnvField.Str == "string"`,
			env:   Env{EnvField: EnvField{Str: "string"}},
			want:  true,
		},
		{
			input: `EnvField.S == "string"`,
			env:   Env{EnvField: EnvField{Str: "string"}},
			err:   "invalid operation: == (mismatched types issue_test.EnvStr and string)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := expr.Compile(tt.input, expr.Env(tt.env), expr.AsBool())

			if tt.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				out, err := expr.Run(program, tt.env)
				require.NoError(t, err)
				require.Equal(t, tt.want, out)
			}
		})
	}
}
