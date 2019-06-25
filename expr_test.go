package expr_test

import (
	"encoding/json"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

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

func ExampleOperator() {
	type Place struct {
		Code string
	}
	type Segment struct {
		Origin Place
	}
	type Helpers struct {
		PlaceEq func(p Place, s string) bool
	}
	type Request struct {
		Segments []*Segment
		Helpers
	}

	code := `Segments[0].Origin == "MOW" && PlaceEq(Segments[0].Origin, "MOW")`

	program, err := expr.Compile(code, expr.Env(&Request{}), expr.Operator("==", "PlaceEq"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	request := &Request{
		Segments: []*Segment{
			{Origin: Place{Code: "MOW"}},
		},
		Helpers: Helpers{PlaceEq: func(p Place, s string) bool {
			return p.Code == s
		}},
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

func TestExpr(t *testing.T) {
	type mockEnv struct {
		One, Two, Three int
		IntArray        []int
		MultiDimArray   [][]int
		Sum             func(list []int) int
		Inc             func(int) int
	}

	request := mockEnv{
		One:           1,
		Two:           2,
		Three:         3,
		IntArray:      []int{1, 2, 3},
		MultiDimArray: [][]int{{1, 2, 3}, {1, 2, 3}},
		Sum: func(list []int) int {
			var ret int
			for _, el := range list {
				ret += el
			}
			return ret
		},
		Inc: func(a int) int { return a + 1 },
	}

	tests := []struct {
		code string
		want interface{}
	}{
		{
			`1 + 1`,
			2,
		},
		{
			`(One * Two) * Three == One * (Two * Three)`,
			true,
		},
		{
			`IntArray[0]`,
			1,
		},
		{
			`Sum(IntArray)`,
			6,
		},
		{
			`IntArray[0] < IntArray[1]`,
			true,
		},
		{
			`Sum(MultiDimArray[0])`,
			6,
		},
		{
			`Sum(MultiDimArray[0]) + Sum(MultiDimArray[1])`,
			12,
		},
		{
			`Inc(IntArray[0] + IntArray[1])`,
			4,
		},
		{
			`IntArray[0] + IntArray[1]`,
			3,
		},
	}

	for _, tt := range tests {
		program, err := expr.Compile(tt.code, expr.Env(mockEnv{}))
		require.NoError(t, err, tt.code)

		got, err := expr.Run(program, request)
		require.NoError(t, err, tt.code)

		assert.Equal(t, tt.want, got, tt.code)
	}
}
