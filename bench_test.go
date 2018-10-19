package expr_test

import (
	"testing"

	"github.com/antonmedv/expr"
)

type Segment struct {
	Origin string
}
type Passengers struct {
	Adults int
}
type Env struct {
	Segments   []Segment
	Passengers Passengers
	Marker     string
}

func (e *Env) First(s []Segment) string {
	return s[0].Origin
}

var env = Env{
	Segments: []Segment{
		{Origin: "LED"},
		{Origin: "HKT"},
	},
	Passengers: Passengers{
		Adults: 2,
	},
	Marker: "test",
}

func Benchmark_struct(b *testing.B) {
	program, err := expr.Parse(
		`Segments[0].Origin == "LED" && Passengers.Adults == 2 && Marker == "test"`,
		expr.Env(Env{}),
	)
	if err != nil {
		b.Fatal(err)
	}
	out, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	if !out.(bool) {
		panic("unexpected result")
	}

	for n := 0; n < b.N; n++ {
		expr.Run(program, env)
	}
}

func Benchmark_map(b *testing.B) {
	env := map[string]interface{}{
		"segments":   env.Segments,
		"passengers": env.Passengers,
		"marker":     env.Marker,
	}

	program, err := expr.Parse(`segments[0].Origin == "LED" && passengers.Adults == 2 && marker == "test"`)
	if err != nil {
		b.Fatal(err)
	}
	out, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	if !out.(bool) {
		panic("unexpected result")
	}

	for n := 0; n < b.N; n++ {
		expr.Run(program, env)
	}
}

func Benchmark_func(b *testing.B) {
	program, err := expr.Parse(`First(Segments)`, expr.Env(&Env{}))
	if err != nil {
		b.Fatal(err)
	}
	out, err := expr.Run(program, &env)
	if err != nil {
		panic(err)
	}
	if out.(string) != "LED" {
		panic("unexpected result")
	}

	for n := 0; n < b.N; n++ {
		expr.Run(program, &env)
	}
}
