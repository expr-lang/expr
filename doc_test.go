package expr_test

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
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
		fmt.Printf("err: %v", err)
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
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: 42
}

func ExampleEval_error() {
	output, err := expr.Eval("(boo + bar]", nil)

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: err: unclosed "("
	//(boo + bar]
	//----------^
}

func ExampleEval_matches() {
	output, err := expr.Eval(`"a" matches "a("`, nil)

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: err: error parsing regexp: missing closing ): `a(`
	//"a" matches "a("
	//----------------^
}

func ExampleParse() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 99,
	}

	ast, err := expr.Parse("foo in 1..99 and bar in 1..99")

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	output, err := expr.Run(ast, env)

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleRun() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 99,
	}

	ast, err := expr.Parse("foo + bar not in 99..100")

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	output, err := expr.Run(ast, env)

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: false
}

func ExampleDefine() {
	type Group struct {
		Name string
	}
	type User struct {
		Age int
	}

	_, err := expr.Parse("groups[0].Name + user.Age", expr.Define("groups", []Group{}), expr.Define("user", User{}))

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output: err: invalid operation: groups[0].Name + user.Age (mismatched types string and int)
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
	ast, err := expr.Parse(code, expr.Env(&Request{}))

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	r := &Request{
		Segments: []*Segment{
			{Origin: "MOW"},
		},
		Passengers: &Passengers{
			Adults: 2,
		},
		Marker: "test",
		Meta:   map[string]interface{}{"accept": true},
	}
	output, err := expr.Run(ast, r)

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleFuncs() {
	var foo, bar func()
	_, err := expr.Parse("foo(bar(baz()))", expr.Define("foo", foo), expr.Define("bar", bar))

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output: err: unknown func baz()
}

func ExampleNode() {
	node, err := expr.Parse("foo.bar")

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", node)

	// Output: foo.bar
}
