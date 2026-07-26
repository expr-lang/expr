package expr_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestFloat32LiteralComparisons(t *testing.T) {
	env := struct {
		Positive float32
		Negative float32
		Float64  float64
	}{
		Positive: 12.34,
		Negative: -12.34,
		Float64:  12.34,
	}
	tests := []struct {
		code string
		want bool
	}{
		{"Positive == 12.34", true},
		{"12.34 == Positive", true},
		{"Negative == -12.34", true},
		{"Positive <= 12.34", true},
		{"Positive != 12.34", false},
		{"Float64 == 12.34", true},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			for _, optimize := range []bool{true, false} {
				program, err := expr.Compile(test.code, expr.Env(env), expr.Optimize(optimize))
				require.NoError(t, err)

				output, err := expr.Run(program, env)
				require.NoError(t, err)
				assert.Equal(t, test.want, output)
			}
		})
	}
}
