package expr_test

import (
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/require"
)

func Benchmark_expr(b *testing.B) {
	params := make(map[string]interface{})
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = 1
	params["Value"] = 100

	program, err := expr.Compile(`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`, expr.Env(params))
	require.NoError(b, err)

	var out interface{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, params)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

func Benchmark_expr_reuseVm(b *testing.B) {
	params := make(map[string]interface{})
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = 1
	params["Value"] = 100

	program, err := expr.Compile(`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`, expr.Env(params))
	require.NoError(b, err)

	var out interface{}
	v := vm.VM{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = v.Run(program, params)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

func Benchmark_len(b *testing.B) {
	env := map[string]interface{}{
		"arr": make([]int, 100),
	}

	program, err := expr.Compile(`len(arr)`, expr.Env(env))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 100, out)
}

func Benchmark_filter(b *testing.B) {
	type Env struct {
		Ints []int
	}
	env := Env{
		Ints: make([]int, 1000),
	}
	for i := 1; i <= len(env.Ints); i++ {
		env.Ints[i-1] = i
	}

	program, err := expr.Compile(`filter(Ints, # % 7 == 0)`, expr.Env(Env{}))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Len(b, out.([]interface{}), 142)
}

func Benchmark_filterLen(b *testing.B) {
	type Env struct {
		Ints []int
	}
	env := Env{
		Ints: make([]int, 1000),
	}
	for i := 1; i <= len(env.Ints); i++ {
		env.Ints[i-1] = i
	}

	program, err := expr.Compile(`len(filter(Ints, # % 7 == 0))`, expr.Env(Env{}))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 142, out)
}

func Benchmark_filterFirst(b *testing.B) {
	type Env struct {
		Ints []int
	}
	env := Env{
		Ints: make([]int, 1000),
	}
	for i := 1; i <= len(env.Ints); i++ {
		env.Ints[i-1] = i
	}

	program, err := expr.Compile(`filter(Ints, # % 7 == 0)[0]`, expr.Env(Env{}))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 7, out)
}

func Benchmark_filterLast(b *testing.B) {
	type Env struct {
		Ints []int
	}
	env := Env{
		Ints: make([]int, 1000),
	}
	for i := 1; i <= len(env.Ints); i++ {
		env.Ints[i-1] = i
	}

	program, err := expr.Compile(`filter(Ints, # % 7 == 0)[-1]`, expr.Env(Env{}))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}

	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 994, out)
}

func Benchmark_filterMap(b *testing.B) {
	type Env struct {
		Ints []int
	}
	env := Env{
		Ints: make([]int, 1000),
	}
	for i := 1; i <= len(env.Ints); i++ {
		env.Ints[i-1] = i
	}

	program, err := expr.Compile(`map(filter(Ints, # % 7 == 0), # * 2)`, expr.Env(Env{}))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Len(b, out.([]interface{}), 142)
	require.Equal(b, 14, out.([]interface{})[0])
}

func Benchmark_arrayIndex(b *testing.B) {
	env := map[string]interface{}{
		"arr": make([]int, 100),
	}
	for i := 0; i < 100; i++ {
		env["arr"].([]int)[i] = i
	}

	program, err := expr.Compile(`arr[50]`, expr.Env(env))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, 50, out)
}

func Benchmark_envStruct(b *testing.B) {
	type Price struct {
		Value int
	}
	type Env struct {
		Price Price
	}

	program, err := expr.Compile(`Price.Value > 0`, expr.Env(Env{}))
	require.NoError(b, err)

	env := Env{Price: Price{Value: 1}}

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

func Benchmark_envMap(b *testing.B) {
	type Price struct {
		Value int
	}
	env := map[string]interface{}{
		"price": Price{Value: 1},
	}

	program, err := expr.Compile(`price.Value > 0`, expr.Env(env))
	require.NoError(b, err)

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

type CallEnv struct {
	A      int
	B      int
	C      int
	Fn     func() bool
	FnFast func(...interface{}) interface{}
	Foo    CallFoo
}

func (CallEnv) Func() string {
	return "func"
}

type CallFoo struct {
	D int
	E int
	F int
}

func (CallFoo) Method() string {
	return "method"
}

func Benchmark_callFunc(b *testing.B) {
	program, err := expr.Compile(`Func()`, expr.Env(CallEnv{}))
	require.NoError(b, err)

	env := CallEnv{}

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, "func", out)
}

func Benchmark_callMethod(b *testing.B) {
	program, err := expr.Compile(`Foo.Method()`, expr.Env(CallEnv{}))
	require.NoError(b, err)

	env := CallEnv{}

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, "method", out)
}

func Benchmark_callField(b *testing.B) {
	program, err := expr.Compile(`Fn()`, expr.Env(CallEnv{}))
	require.NoError(b, err)

	env := CallEnv{
		Fn: func() bool {
			return true
		},
	}

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

func Benchmark_callFast(b *testing.B) {
	program, err := expr.Compile(`FnFast()`, expr.Env(CallEnv{}))
	if err != nil {
		b.Fatal(err)
	}

	env := CallEnv{
		FnFast: func(s ...interface{}) interface{} {
			return "fn_fast"
		},
	}

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, "fn_fast", out)
}

func Benchmark_callConstExpr(b *testing.B) {
	program, err := expr.Compile(`Func()`, expr.Env(CallEnv{}), expr.ConstExpr("Func"))
	require.NoError(b, err)

	env := CallEnv{}

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.Equal(b, "func", out)
}

func Benchmark_largeStructAccess(b *testing.B) {
	type Env struct {
		Data  [1024 * 1024 * 10]byte
		Field int
	}

	program, err := expr.Compile(`Field > 0 && Field > 1 && Field < 99`, expr.Env(Env{}))
	require.NoError(b, err)

	env := Env{Field: 21}

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, &env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

func Benchmark_largeNestedStructAccess(b *testing.B) {
	type Env struct {
		Inner struct {
			Data  [1024 * 1024 * 10]byte
			Field int
		}
	}

	program, err := expr.Compile(`Inner.Field > 0 && Inner.Field > 1 && Inner.Field < 99`, expr.Env(Env{}))
	require.NoError(b, err)

	env := Env{}
	env.Inner.Field = 21

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, &env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}

func Benchmark_largeNestedArrayAccess(b *testing.B) {
	type Env struct {
		Data [1][1024 * 1024 * 10]byte
	}

	program, err := expr.Compile(`Data[0][0] > 0`, expr.Env(Env{}))
	require.NoError(b, err)

	env := Env{}
	env.Data[0][0] = 1

	var out interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, &env)
	}
	b.StopTimer()

	require.NoError(b, err)
	require.True(b, out.(bool))
}
