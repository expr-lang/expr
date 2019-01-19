package expr_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/antonmedv/expr"
)

type evalTest struct {
	input    string
	env      interface{}
	expected interface{}
}

type evalErrorTest struct {
	input string
	env   interface{}
	err   string
}

type evalParams map[string]interface{}

func (p evalParams) Max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

func (p evalParams) Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

var evalTests = []evalTest{
	{
		"foo",
		map[string]int{"foo": 33},
		33,
	},
	{
		"foo == bar",
		map[string]interface{}{"foo": 1, "bar": 1},
		true,
	},
	{
		"foo || (bar && !false && true)",
		map[string]interface{}{"foo": false, "bar": true},
		true,
	},
	{
		"foo && bar",
		map[string]interface{}{"foo": false, "bar": true},
		false,
	},
	{
		"!foo && bar",
		map[string]interface{}{"foo": false, "bar": true},
		true,
	},
	{
		"true || false",
		nil,
		true,
	},
	{
		"false && true",
		nil,
		false,
	},
	{
		"2+2==4",
		nil,
		true,
	},
	{
		"2+3",
		nil,
		float64(5),
	},
	{
		"5-2",
		nil,
		float64(3),
	},
	{
		"2*3",
		nil,
		float64(6),
	},
	{
		"6/2",
		nil,
		float64(3),
	},
	{
		"8%3",
		nil,
		float64(2),
	},
	{
		"2**4",
		nil,
		float64(16),
	},
	{
		"2**4",
		nil,
		float64(16),
	},
	{
		"-(2-5)**3-2/(+4-3)+-2",
		nil,
		float64(23),
	},
	{
		`"hello" ~ ' ' ~ "world"`,
		nil,
		"hello world",
	},
	{
		"+0 == -0",
		nil,
		true,
	},
	{
		"1 < 2 and 3 > 2",
		nil,
		true,
	},
	{
		"!(1 != 1) && 2 >= 2 && 3 <= 3",
		nil,
		true,
	},
	{
		"[1, 02, 1e3, 1.2e-4]",
		nil,
		[]interface{}{float64(1), float64(2), float64(1000), float64(0.00012)},
	},
	{
		"1..5",
		nil,
		[]float64{1, 2, 3, 4, 5},
	},
	{
		"{foo: 1}",
		nil,
		map[interface{}]interface{}{"foo": float64(1)},
	},
	{
		`{foo: "bar"}.foo`,
		nil,
		"bar",
	},
	{
		`Foo.Bar`,
		struct{ Foo struct{ Bar bool } }{Foo: struct{ Bar bool }{Bar: true}},
		true,
	},
	{
		`Foo.Bar`,
		&struct{ Foo *struct{ Bar bool } }{Foo: &struct{ Bar bool }{Bar: true}},
		true,
	},
	{
		"foo[2]",
		map[string]interface{}{"foo": []rune{'a', 'b', 'c'}},
		'c',
	},
	{
		"len(foo) == 3",
		map[string]interface{}{"foo": []rune{'a', 'b', 'c'}},
		true,
	},
	{
		`len(foo) == 6`,
		map[string]string{"foo": "foobar"},
		true,
	},
	{
		"[1, 2, 3][2/2]",
		nil,
		float64(2),
	},
	{
		`[true][A]`,
		struct{ A int }{0},
		true,
	},
	{
		`A-1`,
		struct{ A int }{1},
		float64(0),
	},
	{
		`A == 0`,
		struct{ A uint8 }{0},
		true,
	},
	{
		`A == B`,
		struct {
			A uint8
			B float64
		}{1, 1},
		true,
	},
	{
		`A == B`,
		struct {
			A float64
			B interface{}
		}{1, new(interface{})},
		false,
	},
	{
		`A == B`,
		struct {
			A interface{}
			B float64
		}{new(interface{}), 1},
		false,
	},
	{
		`[true][A]`,
		&struct{ A int }{0},
		true,
	},
	{
		`A-1`,
		&struct{ A int }{1},
		float64(0),
	},
	{
		`A == 0`,
		&struct{ A uint8 }{0},
		true,
	},
	{
		`A == B`,
		&struct {
			A uint8
			B float64
		}{1, 1},
		true,
	},
	{
		`5 in 0..9`,
		nil,
		true,
	},
	{
		`"1" in ["1", "2"]`,
		nil,
		true,
	},
	{
		`"0" not in ["1", "2"]`,
		nil,
		true,
	},
	{
		`"a" in {a:1, b:2}`,
		nil,
		true,
	},
	{
		`"Bar" in Foo`,
		struct{ Foo struct{ Bar bool } }{struct{ Bar bool }{true}},
		true,
	},
	{
		`"Bar" in Ptr`,
		struct{ Ptr *struct{ Bar bool } }{&struct{ Bar bool }{true}},
		true,
	},
	{
		`"Bar" in NilPtr`,
		struct{ NilPtr *bool }{nil},
		false,
	},
	{
		`0 in nil`,
		nil,
		false,
	},
	{
		`A == nil`,
		struct{ A interface{} }{nil},
		true,
	},
	{
		"foo['bar'].baz",
		map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": true}}},
		true,
	},
	{
		"foo.Bar['baz']",
		map[string]interface{}{"foo": &struct{ Bar map[string]interface{} }{Bar: map[string]interface{}{"baz": true}}},
		true,
	},
	{
		`60 & 13`,
		nil,
		12,
	},
	{
		`60 ^ 13`,
		nil,
		49,
	},
	{
		`60 | 13`,
		nil,
		61,
	},
	{
		`"seafood" matches "foo.*"`,
		nil,
		true,
	},
	{
		`"seafood" matches "sea" ~ "food"`,
		nil,
		true,
	},
	{
		`not ("seafood" matches "[0-9]+") ? "a" : "b"`,
		nil,
		"a",
	},
	{
		`false ? "a" : "b"`,
		nil,
		"b",
	},
	{
		`foo.bar("world")`,
		map[string]interface{}{"foo": map[string]interface{}{"bar": func(in string) string { return "hello " + in }}},
		"hello world",
	},
	{
		`foo.bar()`,
		map[string]interface{}{"foo": map[string]interface{}{"bar": func() {}}},
		nil,
	},
	{
		`foo("world")`,
		map[string]interface{}{"foo": func(in string) string { return "hello " + in }},
		"hello world",
	},
	{
		"Max(a, b)",
		evalParams{"a": 1.23, "b": 3.21},
		3.21,
	},
	{
		"Min(a, b)",
		evalParams{"a": 1.23, "b": 3.21},
		1.23,
	},
}

var evalErrorTests = []evalErrorTest{
	{
		"bar",
		map[string]int{"foo": 1},
		`undefined: bar`,
	},
	{
		`"foo" ~ foo`,
		map[string]*int{"foo": nil},
		`interface conversion: interface {} is *int, not string`,
	},
	{
		"1 or 0",
		nil,
		"interface conversion: interface {} is float64, not bool",
	},
	{
		"not nil",
		nil,
		"interface conversion: interface {} is nil, not bool",
	},
	{
		"nil matches 'nil'",
		nil,
		"interface conversion: interface {} is nil, not string",
	},
	{
		"foo['bar'].baz",
		map[string]interface{}{"foo": nil},
		`cannot get "bar" from <nil>: foo["bar"]`,
	},
	{
		"foo.bar(abc)",
		map[string]interface{}{"foo": nil},
		`cannot get method bar from <nil>: foo.bar(abc)`,
	},
	{
		`"seafood" matches "a(b"`,
		nil,
		"error parsing regexp: missing closing ): `a(b`",
	},
	{
		`"seafood" matches "a" ~ ")b"`,
		nil,
		"error parsing regexp: unexpected ): `a)b`",
	},
	{
		`1 matches "1" ~ "2"`,
		nil,
		"interface conversion: interface {} is float64, not string",
	},
	{
		`1 matches "1"`,
		nil,
		"interface conversion: interface {} is float64, not string",
	},
	{
		`"1" matches 1`,
		nil,
		"interface conversion: interface {} is float64, not string",
	},
	{
		`foo ? 1 : 2`,
		map[string]interface{}{"foo": 0},
		`interface conversion: interface {} is int, not bool`,
	},
	{
		`foo()`,
		map[string]interface{}{"foo": func() (int, int) { return 0, 1 }},
		`func "foo" must return only one value`,
	},
	{
		`foo()`,
		map[string]interface{}{"foo": nil},
		`reflect: call of reflect.Value.Call on zero Value`,
	},
	{
		"1..1e6+1",
		nil,
		"range 1..1000001 exceeded max size of 1e6",
	},
	{
		"1/0",
		nil,
		"division by zero",
	},
	{
		"1%0",
		nil,
		"division by zero",
	},
	{
		"1 + 'a'",
		nil,
		`cannot convert "a" (type string) to type float64`,
	},
	{
		"'a' + 1",
		nil,
		`cannot convert "a" (type string) to type float64`,
	},
	{
		"[1, 2]['a']",
		nil,
		`cannot get "a" from []interface {}: [1, 2]["a"]`,
	},
	{
		`1 in "a"`,
		nil,
		`operator "in" not defined on string`,
	},
	{
		`nil in map`,
		map[string]interface{}{"map": map[string]interface{}{"true": "yes"}},
		`cannot use <nil> as index to map[string]interface {}`,
	},
	{
		`nil in foo`,
		map[string]interface{}{"foo": struct{ Bar bool }{true}},
		`cannot use <nil> as field name of struct { Bar bool }`,
	},
	{
		`true in foo`,
		map[string]interface{}{"foo": struct{ Bar bool }{true}},
		`cannot use bool as field name of struct { Bar bool }`,
	},
	{
		"len()",
		nil,
		"missing argument: len()",
	},
	{
		"len(1)",
		nil,
		"invalid argument len(1) (type float64)",
	},
	{
		"len(a, b)",
		nil,
		"too many arguments: len(a, b)",
	},
	{
		"Foo.Map",
		struct{ Foo map[string]int }{Foo: nil},
		"Foo is nil",
	},
	{
		"Foo.Bar",
		struct{ Foo *struct{ Bar bool } }{Foo: nil},
		"Foo is nil",
	},
	{
		"Foo.Panic",
		struct{ Foo interface{} }{Foo: nil},
		"Foo is nil",
	},
}

func TestEval(t *testing.T) {
	for _, test := range evalTests {
		actual, err := expr.Eval(test.input, test.env)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.input, actual, test.expected)
		}
	}
}

func TestEval_error(t *testing.T) {
	for _, test := range evalErrorTests {
		result, err := expr.Eval(test.input, test.env)
		if err == nil {
			err = fmt.Errorf("%v, <nil>", result)
		}
		if !strings.HasPrefix(err.Error(), test.err) || test.err == "" {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, err.Error(), test.err)
		}
	}
}

func TestEval_complex(t *testing.T) {
	type cookie struct {
		Key   string
		Value string
	}
	type user struct {
		UserAgent string
		Cookies   []cookie
	}
	type request struct {
		User user
	}

	p := map[string]interface{}{
		"Request": request{user{
			Cookies:   []cookie{{"origin", "www"}},
			UserAgent: "Mozilla 1",
		}},
		"Values": func(xs []cookie) []string {
			vs := make([]string, 0)
			for _, x := range xs {
				vs = append(vs, x.Value)
			}
			return vs
		},
	}

	input := `Request.User.UserAgent matches "Mozilla" && "www" in Values(Request.User.Cookies)`
	node, err := expr.Parse(input, expr.Env(p))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := expr.Run(node, p)
	if err != nil {
		t.Fatal(err)
	}

	expected := true
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("TestEval_complex:\ngot\n\t%#v\nexpected:\n\t%#v", actual, expected)
	}
}

func TestEval_panic(t *testing.T) {
	node, err := expr.Parse("foo()")
	if err != nil {
		t.Fatal(err)
	}

	_, err = expr.Run(node, map[string]interface{}{"foo": nil})
	if err == nil {
		err = fmt.Errorf("<nil>")
	}

	expected := "reflect: call of reflect.Value.Call on zero Value"
	if err.Error() != expected {
		t.Errorf("\ngot\n\t%+v\nexpected\n\t%v", err.Error(), expected)
	}
}

func TestEval_func(t *testing.T) {
	type testEnv struct {
		Func func() string
	}

	env := &testEnv{
		Func: func() string {
			return "func"
		},
	}

	input := `Func()`

	node, err := expr.Parse(input, expr.Env(&testEnv{}))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := expr.Run(node, env)
	if err != nil {
		t.Fatal(err)
	}

	expected := "func"
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("TestEval_method:\ngot\n\t%#v\nexpected:\n\t%#v", actual, expected)
	}
}

func TestEval_method(t *testing.T) {
	env := &testEnv{
		Hello: "hello",
		World: testWorld{
			name: []string{"w", "o", "r", "l", "d"},
		},
		testVersion: &testVersion{
			version: 2,
		},
	}

	input := `Title(Hello) ~ Space() ~ (CompareVersion(1, 3) ? World.String() : '')`

	node, err := expr.Parse(input, expr.Env(&testEnv{}))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := expr.Run(node, env)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello world"
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("TestEval_method:\ngot\n\t%#v\nexpected:\n\t%#v", actual, expected)
	}
}

type testVersion struct {
	version float64
}

func (c *testVersion) CompareVersion(min, max float64) bool {
	return min < c.version && c.version < max
}

type testWorld struct {
	name []string
}

func (w testWorld) String() string {
	return strings.Join(w.name, "")
}

type testEnv struct {
	*testVersion
	Hello string
	World testWorld
}

func (e *testEnv) Title(s string) string {
	return strings.Title(s)
}

func (e *testEnv) Space() string {
	return " "
}
