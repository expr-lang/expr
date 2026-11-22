package issue_test

import (
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

// Checking nested calls of overloaded functions used to take time
// exponential in the nesting depth. At depth 100 it would never finish.
func TestIssue1008(t *testing.T) {
	const depth = 100

	replaceAny := expr.Function("replaceAny",
		func(params ...any) (any, error) { return params[0], nil },
		new(func(any, string) any),
		new(func(any, string, string) any),
	)

	tests := []struct {
		name string
		code string
		opts []expr.Option
		want any
	}{
		{
			name: "replace",
			code: strings.Repeat("replace(", depth) + `"a"` + strings.Repeat(`, "a", "a")`, depth),
			want: "a",
		},
		{
			name: "join",
			code: strings.Repeat("join([", depth) + `"a"` + strings.Repeat(`], ",")`, depth),
			want: "a",
		},
		{
			name: "trim with binary argument",
			code: strings.Repeat(`trim("" + `, depth) + `"a"` + strings.Repeat(`, " ")`, depth),
			want: "a",
		},
		{
			name: "function returning any",
			code: strings.Repeat("replaceAny(", depth) + `"a"` + strings.Repeat(`, "a", "a")`, depth),
			opts: []expr.Option{replaceAny},
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.code, tt.opts...)
			require.NoError(t, err)
			out, err := expr.Run(program, nil)
			require.NoError(t, err)
			require.Equal(t, tt.want, out)
		})
	}
}
