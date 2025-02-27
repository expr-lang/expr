package optimizer_test

import (
	"testing"

	"expr/internal/testify/require"

	"expr"
	"expr/vm"
)

func BenchmarkSumArray(b *testing.B) {
	env := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	program, err := expr.Compile(`sum([a, b, c, d])`, expr.Env(env))
	require.NoError(b, err)

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 10, out)
}
