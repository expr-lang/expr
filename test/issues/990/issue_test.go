package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type Outer struct {
	Inner  map[string]struct{} `expr:"bar"`
	Plain  int
	Hidden int `expr:"-"`
}

type Env struct {
	Outer Outer `expr:"foo"`
}

func TestInRespectsExprTags(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		{`"bar" in foo`, true},
		{`"Plain" in foo`, true},
		{`"Hidden" in foo`, false},
		{`"missing" in foo`, false},
		{`"foo" in $env`, true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code, expr.Env(Env{}))
			require.NoError(t, err)

			out, err := expr.Run(program, Env{})
			require.NoError(t, err)
			require.Equal(t, tt.want, out)

			out, err = expr.Eval(tt.code, Env{})
			require.NoError(t, err)
			require.Equal(t, tt.want, out)
		})
	}
}
