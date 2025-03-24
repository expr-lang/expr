package bench_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/vm"
)

type Env struct {
	Fn func() bool
}

func BenchmarkCall_callTyped(b *testing.B) {
	code := `Fn()`

	p, err := expr.Compile(code, expr.Env(Env{}))
	require.NoError(b, err)
	require.Equal(b, p.Bytecode[1], vm.OpCallTyped)

	env := Env{
		Fn: func() bool {
			return true
		},
	}

	var out any

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		program, _ := expr.Compile(code, expr.Env(env))
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

func BenchmarkCall_eval(b *testing.B) {
	code := `Fn()`

	p, err := expr.Compile(code)
	require.NoError(b, err)
	require.Equal(b, p.Bytecode[1], vm.OpCall)

	env := Env{
		Fn: func() bool {
			return true
		},
	}

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = expr.Eval(code, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}
