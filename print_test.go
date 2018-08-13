package expr

import (
	"fmt"
	"reflect"
	"testing"
)

type printTest struct {
	input    Node
	expected string
}

var printTests = []printTest{
	{
		methodNode{nameNode{"foo"}, "bar", []Node{textNode{"arg1"}, numberNode{2}, boolNode{true}}},
		`foo.bar("arg1", 2, true)`,
	},
	{
		indexNode{propertyNode{methodNode{methodNode{nameNode{"foo"}, "bar", []Node{}}, "foo", []Node{}}, "baz"}, numberNode{33}},
		"foo.bar().foo().baz[33]",
	},
	{
		mapNode{[]pairNode{{identifierNode{"foo"}, numberNode{1}}, {binaryNode{"+", numberNode{1}, numberNode{2}}, numberNode{2}}}},
		`{"foo": 1, (1 + 2): 2}`,
	},
	{
		functionNode{"call", []Node{propertyNode{arrayNode{[]Node{numberNode{1}, unaryNode{"not", boolNode{true}}}}, "foo"}}},
		"call([1, not true].foo)",
	},
	{
		builtinNode{"len", []Node{nameNode{"array"}}},
		"len(array)",
	},
	{
		binaryNode{"or", binaryNode{"or", nameNode{"a"}, nameNode{"b"}}, nameNode{"c"}},
		"((a or b) or c)",
	},
	{
		binaryNode{"or", nameNode{"a"}, binaryNode{"and", nameNode{"b"}, nameNode{"c"}}},
		"(a or (b and c))",
	},
	{
		binaryNode{"and", binaryNode{"or", nameNode{"a"}, nameNode{"b"}}, nameNode{"c"}},
		"((a or b) and c)",
	},
	{
		conditionalNode{nameNode{"a"}, nameNode{"a"}, nameNode{"b"}},
		"a ? a : b",
	},
}

func TestPrint(t *testing.T) {
	for _, test := range printTests {
		actual := fmt.Sprintf("%v", test.input)
		if actual != test.expected {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.expected, actual, test.expected)
		}
		// Parse again and check if ast same as before.
		ast, err := Parse(actual)
		if err != nil {
			t.Errorf("%s: can't parse printed expression", actual)
		}
		if !reflect.DeepEqual(ast, test.input) {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.expected, ast, test.input)
		}
	}
}

func TestPrint_matches(t *testing.T) {
	input := matchesNode{left: nameNode{"foo"}, right: textNode{"foobar"}}
	expected := "(foo matches \"foobar\")"

	actual := fmt.Sprintf("%v", input)
	if actual != expected {
		t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", expected, actual, expected)
	}
	// Parse again and check if ast same as before.
	ast, err := Parse(actual)
	if err != nil {
		t.Errorf("%s: can't parse printed expression", actual)
	}

	// Clear parsed regexp.
	m := ast.(matchesNode)
	m.r = nil

	if !reflect.DeepEqual(m, input) {
		t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", expected, m, input)
	}
}
