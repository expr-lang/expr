package expr_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestFloat32LiteralComparisons(t *testing.T) {
	f := float32(12.34)
	env := struct {
		Positive float32
		Negative float32
		Float64  float64
		Pointer  *float32
		Array    []float32
	}{
		Positive: 12.34,
		Negative: -12.34,
		Float64:  12.34,
		Pointer:  &f,
		Array:    []float32{1.5, 12.34},
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
		{"Pointer == 12.34", true},
		{"Array[1] == 12.34", true},
		{"Positive in [1.5, 12.34]", true},
		{"Positive not in [1.5, 12.34]", false},
		{"12.34 in Array", true},
		{"-12.34 in Array", false},
		{"Float64 in [12.34]", true},
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
