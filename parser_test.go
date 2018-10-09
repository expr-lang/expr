package expr

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type parseTest struct {
	input    string
	expected Node
}

type parseErrorTest struct {
	input string
	err   string
}

var parseTests = []parseTest{
	{
		"a",
		nameNode{"a"},
	},
	{
		`"a"`,
		textNode{"a"},
	},
	{
		"3",
		numberNode{3},
	},
	{
		"true",
		boolNode{true},
	},
	{
		"false",
		boolNode{false},
	},
	{
		"nil",
		nilNode{},
	},
	{
		"-3",
		unaryNode{"-", numberNode{3}},
	},
	{
		"1 - 2",
		binaryNode{"-", numberNode{1}, numberNode{2}},
	},
	{
		"(1 - 2) * 3",
		binaryNode{"*", binaryNode{"-", numberNode{1}, numberNode{2}}, numberNode{3}},
	},
	{
		"a or b or c",
		binaryNode{"or", binaryNode{"or", nameNode{"a"}, nameNode{"b"}}, nameNode{"c"}},
	},
	{
		"a or b and c",
		binaryNode{"or", nameNode{"a"}, binaryNode{"and", nameNode{"b"}, nameNode{"c"}}},
	},
	{
		"(a or b) and c",
		binaryNode{"and", binaryNode{"or", nameNode{"a"}, nameNode{"b"}}, nameNode{"c"}},
	},
	{
		"2**4-1",
		binaryNode{"-", binaryNode{"**", numberNode{2}, numberNode{4}}, numberNode{1}},
	},
	{
		"foo(bar())",
		functionNode{"foo", []Node{functionNode{"bar", []Node{}}}},
	},
	{
		"foo.bar",
		propertyNode{nameNode{"foo"}, "bar"},
	},
	{
		"foo.not",
		propertyNode{nameNode{"foo"}, "not"},
	},
	{
		"foo.bar()",
		methodNode{nameNode{"foo"}, "bar", []Node{}},
	},
	{
		"foo.not()",
		methodNode{nameNode{"foo"}, "not", []Node{}},
	},
	{
		`foo.bar("arg1", 2, true)`,
		methodNode{nameNode{"foo"}, "bar", []Node{textNode{"arg1"}, numberNode{2}, boolNode{true}}},
	},
	{
		"foo[3]",
		indexNode{nameNode{"foo"}, numberNode{3}},
	},
	{
		"true ? true : false",
		conditionalNode{boolNode{true}, boolNode{true}, boolNode{false}},
	},
	{
		"a ?: b",
		conditionalNode{nameNode{"a"}, nameNode{"a"}, nameNode{"b"}},
	},
	{
		"foo.bar().foo().baz[33]",
		indexNode{propertyNode{methodNode{methodNode{nameNode{"foo"}, "bar", []Node{}}, "foo", []Node{}}, "baz"}, numberNode{33}},
	},
	{
		"+0 != -0",
		binaryNode{"!=", unaryNode{"+", numberNode{0}}, unaryNode{"-", numberNode{0}}},
	},
	{
		"[a, b, c]",
		arrayNode{[]Node{nameNode{"a"}, nameNode{"b"}, nameNode{"c"}}},
	},
	{
		"{foo:1, bar:2}",
		mapNode{[]pairNode{{identifierNode{"foo"}, numberNode{1}}, {identifierNode{"bar"}, numberNode{2}}}},
	},
	{
		`{"foo":1, (1+2):2}`,
		mapNode{[]pairNode{{identifierNode{"foo"}, numberNode{1}}, {binaryNode{"+", numberNode{1}, numberNode{2}}, numberNode{2}}}},
	},
	{
		"[1].foo",
		propertyNode{arrayNode{[]Node{numberNode{1}}}, "foo"},
	},
	{
		"{foo:1}.bar",
		propertyNode{mapNode{[]pairNode{{identifierNode{"foo"}, numberNode{1}}}}, "bar"},
	},
	{
		"len(foo)",
		builtinNode{"len", []Node{nameNode{"foo"}}},
	},
	{
		`foo matches "foo"`,
		matchesNode{left: nameNode{"foo"}, right: textNode{"foo"}},
	},
	{
		`foo matches regex`,
		matchesNode{left: nameNode{"foo"}, right: nameNode{"regex"}},
	},
}

var parseErrorTests = []parseErrorTest{
	{
		"foo.",
		"unexpected end of expression",
	},
	{
		"a+",
		"unexpected token EOF",
	},
	{
		"a ? (1+2) c",
		"unexpected token name(c)",
	},
	{
		"[a b]",
		"array items must be separated by a comma",
	},
	{
		"foo.bar(a b)",
		"arguments must be separated by a comma",
	},
	{
		"{-}",
		"a map key must be a",
	},
	{
		"a matches 'a)(b'",
		"error parsing regexp: unexpected )",
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		actual, err := Parse(test.input)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		if m, ok := actual.(matchesNode); ok {
			m.r = nil
			actual = m
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.input, actual, test.expected)
		}
	}
}

func TestParse_error(t *testing.T) {
	for _, test := range parseErrorTests {
		_, err := Parse(test.input)
		if err == nil {
			err = fmt.Errorf("<nil>")
		}
		if !strings.HasPrefix(err.Error(), test.err) || test.err == "" {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, err.Error(), test.err)
		}
	}
}

func TestParser_createTypesTable(t *testing.T) {
	var intType = reflect.TypeOf(0)

	type (
		D struct {
			F2 int
		}

		C struct {
			F int
		}

		B struct {
			C
		}

		A struct {
			*D
			B
		}
	)

	p := parser{}
	types := p.createTypesTable(A{})

	if len(types) != 5 {
		t.Error("unexpected number of fields")
	}
	if types["F"] != intType {
		t.Error("expected embedded struct field 'F'")
	}
	if types["F2"] != intType {
		t.Error("expected embedded struct field 'F2'")
	}
}
