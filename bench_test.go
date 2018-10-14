package expr_test

import (
	"testing"

	"github.com/antonmedv/expr"
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

	script, err := expr.Parse(`Segments[0].Origin == "MOW" && Passengers.Adults == 2 && Marker == "test"`)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		expr.Run(script, r)
	}
}
