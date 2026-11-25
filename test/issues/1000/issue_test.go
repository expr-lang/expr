package issue1000

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type Rec struct{ X int }
type Env struct{ L []Rec }

func TestIssue1000(t *testing.T) {
	env := Env{L: []Rec{{1}, {2}}}
	tests := []struct {
		code string
		want int
	}{
		{`filter(L, true)[0].X`, 1},
		{`filter(L, true)[-1].X`, 2},
		{`find(L, true).X`, 1},
		{`findLast(L, true).X`, 2},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code, expr.Env(Env{}))
			require.NoError(t, err)
			out, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, tt.want, out)
		})
	}
}
