package expr

import (
	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
	"testing"
)

type segment struct {
	Origin string
}
type passengers struct {
	Adults int
}
type request struct {
	Segments   []*segment
	Passengers *passengers
	Marker     string
}

func Benchmark_expr(b *testing.B) {
	r := &request{
		Segments: []*segment{
			{Origin: "MOW"},
		},
		Passengers: &passengers{
			Adults: 2,
		},
		Marker: "test",
	}

	script, err := Parse(`Segments[0].Origin == "MOW" && Passengers.Adults == 2 && Marker == "test"`)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		Run(script, r)
	}
}

func Benchmark_otto(b *testing.B) {
	r := &request{
		Segments: []*segment{
			{Origin: "MOW"},
		},
		Passengers: &passengers{
			Adults: 2,
		},
		Marker: "test",
	}

	vm := otto.New()

	script, err := vm.Compile("", `r.Segments[0].Origin == "MOW" && r.Passengers.Adults == 2 && r.Marker == "test"`)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		vm.Set("r", r)
		vm.Run(script)
	}
}

func Benchmark_goja(b *testing.B) {
	r := &request{
		Segments: []*segment{
			{Origin: "MOW"},
		},
		Passengers: &passengers{
			Adults: 2,
		},
		Marker: "test",
	}

	vm := goja.New()
	program, err := goja.Compile("", `r.Segments[0].Origin == "MOW" && r.Passengers.Adults == 2 && r.Marker == "test"`, false)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		vm.Set("r", r)
		vm.RunProgram(program)
	}
}
