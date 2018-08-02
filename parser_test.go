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
		"2**4-1",
		binaryNode{"-", binaryNode{"**", numberNode{2}, numberNode{4}}, numberNode{1}},
	},
	{
		"foo(bar())",
		functionNode{"foo", []Node{functionNode{"bar", []Node{}}}},
	},
	{
		"foo.bar",
		propertyNode{nameNode{"foo"}, identifierNode{"bar"}},
	},
	{
		"foo.not",
		propertyNode{nameNode{"foo"}, identifierNode{"not"}},
	},
	{
		"foo.bar()",
		methodNode{nameNode{"foo"}, identifierNode{"bar"}, []Node{}},
	},
	{
		"foo.not()",
		methodNode{nameNode{"foo"}, identifierNode{"not"}, []Node{}},
	},
	{
		`foo.bar("arg1", 2, true)`,
		methodNode{nameNode{"foo"}, identifierNode{"bar"}, []Node{textNode{"arg1"}, numberNode{2}, boolNode{true}}},
	},
	{
		"foo[3]",
		propertyNode{nameNode{"foo"}, numberNode{3}},
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
		`"foo" matches "/foo/"`,
		binaryNode{"matches", textNode{"foo"}, textNode{"/foo/"}},
	},
	{
		"foo.bar().foo().baz[33]",
		propertyNode{propertyNode{methodNode{methodNode{nameNode{"foo"}, identifierNode{"bar"}, []Node{}}, identifierNode{"foo"}, []Node{}}, identifierNode{"baz"}}, numberNode{33}},
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
		propertyNode{arrayNode{[]Node{numberNode{1}}}, identifierNode{"foo"}},
	},
	{
		"{foo:1}.bar",
		propertyNode{mapNode{[]pairNode{{identifierNode{"foo"}, numberNode{1}}}}, identifierNode{"bar"}},
	},
	{
		"len(foo)",
		builtinNode{"len", []Node{nameNode{"foo"}}},
	},
}

type parseErrorTest struct {
	input string
	err   string
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
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		actual, err := Parse(test.input)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.input, actual, test.expected)
		}
	}
}

func TestParseError(t *testing.T) {
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

func TestParseErrorPosition(t *testing.T) {
	_, err := Parse("foo() + bar(**)")
	if err == nil {
		err = fmt.Errorf("<nil>")
	}

	expected := "unexpected token operator(**)\nfoo() + bar(**)\n------------^"
	if err.Error() != expected {
		t.Errorf("\ngot\n\t%+v\nexpected\n\t%v", err.Error(), expected)
	}
}
