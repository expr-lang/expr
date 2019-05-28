package expr_test

import (
	"encoding/json"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"strings"
	"testing"
)

func Benchmark_expr(b *testing.B) {
	params := make(map[string]interface{})
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = int64(1)
	params["Value"] = int64(100)

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

func ExampleEval() {
	output, err := expr.Eval("'hello world'", nil)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: hello world
}

func ExampleEval_map() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": []string{"zero", "hello world"},
		"swipe": func(in string) string {
			return strings.Replace(in, "world", "user", 1)
		},
	}

	output, err := expr.Eval("swipe(bar[foo])", env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: hello user
}

type mockMapEnv map[string]interface{}

func (mockMapEnv) Swipe(in string) string {
	return strings.Replace(in, "world", "user", 1)
}

func ExampleEval_map_method() {
	env := mockMapEnv{
		"foo": 1,
		"bar": []string{"zero", "hello world"},
	}

	program, err := expr.Compile("Swipe(bar[foo])", expr.Env(env))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: hello user
}

func ExampleEval_struct() {
	type C struct{ C int }
	type B struct{ B *C }
	type A struct{ A B }

	env := A{B{&C{42}}}

	output, err := expr.Eval("A.B.C", env)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: 42
}

func ExampleEval_error() {
	output, err := expr.Eval("(boo + bar]", nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: syntax error: mismatched input ']' expecting ')' (1:11)
	//  | (boo + bar]
	//  | ..........^
}

func ExampleEval_matches() {
	output, err := expr.Eval(`"a" matches "a("`, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: error parsing regexp: missing closing ): `a(` (1:13)
	//  | "a" matches "a("
	//  | ............^
}

func ExampleRun() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 99,
	}

	program, err := expr.Compile("foo in 1..99 and bar in 1..99", expr.Env(env))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEnv() {
	type Segment struct {
		Origin string
	}
	type Passengers struct {
		Adults int
	}
	type Request struct {
		Segments   []*Segment
		Passengers *Passengers
		Marker     string
		Meta       map[string]interface{}
	}

	code := `Segments[0].Origin == "MOW" && Passengers.Adults == 2 && Marker == "test" && Meta["accept"]`

	program, err := expr.Compile(code, expr.Env(&Request{}))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	request := &Request{
		Segments: []*Segment{
			{Origin: "MOW"},
		},
		Passengers: &Passengers{
			Adults: 2,
		},
		Marker: "test",
		Meta:   map[string]interface{}{"accept": true},
	}

	output, err := expr.Run(program, request)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEval_marshal() {
	env := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	program, err := expr.Compile("foo + bar", expr.Env(env))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	b, err := json.Marshal(program)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	unmarshaledProgram := &vm.Program{}
	err = json.Unmarshal(b, unmarshaledProgram)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(unmarshaledProgram, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: 3
}
