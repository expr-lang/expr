package value

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

func Benchmark_valueAdd(b *testing.B) {
	env := make(map[string]any)
	env["ValueOne"] = &customInt{1}
	env["ValueTwo"] = &customInt{2}

	program, err := expr.Compile("ValueOne + ValueTwo", expr.Env(env), ValueGetter)
	require.NoError(b, err)

	var out any
	v := vm.VM{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = v.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 3, out.(int))
}

func Benchmark_valueUntypedAdd(b *testing.B) {
	env := make(map[string]any)
	env["ValueOne"] = &customUntypedInt{1}
	env["ValueTwo"] = &customUntypedInt{2}

	program, err := expr.Compile("ValueOne + ValueTwo", expr.Env(env), ValueGetter)
	require.NoError(b, err)

	var out any
	v := vm.VM{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = v.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 3, out.(int))
}

func Benchmark_valueTypedAdd(b *testing.B) {
	env := make(map[string]any)
	env["ValueOne"] = &customTypedInt{1}
	env["ValueTwo"] = &customTypedInt{2}

	program, err := expr.Compile("ValueOne + ValueTwo", expr.Env(env), ValueGetter)
	require.NoError(b, err)

	var out any
	v := vm.VM{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = v.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 3, out.(int))
}
