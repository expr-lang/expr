package expr_test

import (
	"gopkg.in/antonmedv/expr.v2"
	"gopkg.in/antonmedv/expr.v2/vm"
	"testing"
)

func Benchmark_expr(b *testing.B) {
	params := make(map[string]interface{})
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = 1
	params["Value"] = 100

	program, err := expr.Compile(`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`, expr.Env(params))
	if err != nil {
		b.Fatal(err)
	}

	var out interface{}

	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, params)
	}

	if err != nil {
		b.Fatal(err)
	}
	if !out.(bool) {
		b.Fail()
	}
}

func Benchmark_filter(b *testing.B) {
	params := make(map[string]interface{})
	params["max"] = 50

	program, err := expr.Compile(`filter(1..100, {# > max})`, expr.Env(params))
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		_, err = vm.Run(program, params)
	}

	if err != nil {
		b.Fatal(err)
	}
}

func Benchmark_access(b *testing.B) {
	type Price struct {
		Value int
	}
	type Env struct {
		Price Price
	}

	program, err := expr.Compile(`Price.Value > 0`, expr.Env(Env{}))
	if err != nil {
		b.Fatal(err)
	}

	env := Env{Price: Price{Value: 1}}

	for n := 0; n < b.N; n++ {
		_, err = vm.Run(program, env)
	}

	if err != nil {
		b.Fatal(err)
	}
}
