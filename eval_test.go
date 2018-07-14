package expr

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type evalTest struct {
	input    string
	env      interface{}
	expected interface{}
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
		"foo[2]",
		map[string]interface{}{"foo": []rune{'a', 'b', 'c'}},
		'c',
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
		`not ("seafood" matches "[0-9]+") ? "a" : "b"`,
		nil,
		"a",
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
}

type evalErrorTest struct {
	input string
	env   interface{}
	err   string
}

var evalErrorTests = []evalErrorTest{
	{
		"bar",
		map[string]int{"foo": 1},
		`can't get "bar"`,
	},
	{
		`"foo" ~ nil`,
		nil,
		`operator ~ not defined on (string, <nil>)`,
	},
	{
		"foo['bar'].baz",
		map[string]interface{}{"foo": nil},
		`can't get "bar" from <nil>`,
	},
	{
		`"seafood" matches "a(b"`,
		nil,
		`error parsing regexp:`,
	},
	{
		`0 ? 1 : 2`,
		nil,
		`non-bool value used in cond`,
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
}

func TestEval(t *testing.T) {
	for _, test := range evalTests {
		actual, err := Eval(test.input, test.env)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.input, actual, test.expected)
		}
	}
}

func TestEvalError(t *testing.T) {
	for _, test := range evalErrorTests {
		_, err := Eval(test.input, test.env)
		if err == nil {
			err = fmt.Errorf("<nil>")
		}
		if !strings.HasPrefix(err.Error(), test.err) {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, err.Error(), test.err)
		}
	}
}

func TestEvalComplex(t *testing.T) {
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
	node, err := Parse(input, Names("Request"), Funcs("Values"))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Run(node, p)
	if err != nil {
		t.Fatal(err)
	}

	expected := true
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("TestEvalComplex:\ngot\n\t%#v\nexpected:\n\t%#v", actual, expected)
	}
}
