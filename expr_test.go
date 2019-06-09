package expr_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
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
	tests := []struct {
		name    string
		code    string
		env     interface{}
		request interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "+ operator",
			code: "1 + 1",
			want: 2,
		},
		{
			name:    "associativity",
			code:    "(A * B) * C == A * (B * C)",
			env:     struct{ A, B, C int }{},
			request: struct{ A, B, C int }{A: 1, B: 2, C: 3},
			want:    true,
		},
		{
			name:    "indexing",
			code:    "A[0]",
			env:     struct{ A []int }{},
			request: struct{ A []int }{A: []int{1}},
			want:    1,
		},
		{
			name: "helpers",
			code: "Sum(A)",
			env: struct {
				A   []int
				Sum func(list []int) int
			}{},
			request: struct {
				A   []int
				Sum func(list []int) int
			}{
				A: []int{1, 2, 3},
				Sum: func(list []int) int {
					var ret int
					for _, el := range list {
						ret += el
					}
					return ret
				},
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.code, expr.Env(tt.env))
			if err != nil {
				t.Errorf("Compile() error = %v", err)
				return
			}

			got, err := expr.Run(program, tt.request)
			if err != nil {
				t.Errorf("Run() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
