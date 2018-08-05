package expr_test

import (
	"fmt"
	"github.com/antonmedv/expr"
	"strings"
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

func ExampleNames() {
	_, err := expr.Parse("foo + bar + baz", expr.Names("foo", "bar"))

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output: err: unknown name baz
	//foo + bar + baz
	//------------^
}

func ExampleFuncs() {
	_, err := expr.Parse("foo(bar(baz()))", expr.Funcs("foo", "bar"))

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output: err: unknown func baz
	//foo(bar(baz()))
	//--------^
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
