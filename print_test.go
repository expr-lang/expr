package expr

import (
	"fmt"
	"testing"
)

type printTest struct {
	input    Node
	expected string
}

var printTests = []printTest{
	{
		methodNode{nameNode{"foo"}, identifierNode{"bar"}, []Node{textNode{"arg1"}, numberNode{2}, boolNode{true}}},
		`foo.bar("arg1", 2, true)`,
	},
	{
		propertyNode{propertyNode{methodNode{methodNode{nameNode{"foo"}, identifierNode{"bar"}, []Node{}}, identifierNode{"foo"}, []Node{}}, identifierNode{"baz"}}, numberNode{33}},
		"foo.bar().foo().baz[33]",
	},
	{
		mapNode{[]pairNode{{identifierNode{"foo"}, numberNode{1}}, {binaryNode{"+", numberNode{1}, numberNode{2}}, numberNode{2}}}},
		`{"foo": 1, (1 + 2): 2}`,
	},
	{
		functionNode{"call", []Node{propertyNode{arrayNode{[]Node{numberNode{1}, unaryNode{"not", boolNode{true}}}}, identifierNode{"foo"}}}},
		"call([1, not true].foo)",
	},
	{
		builtinNode{"len", []Node{identifierNode{"array"}}},
		"len(array)",
	},
}

func TestPrint(t *testing.T) {
	for _, test := range printTests {
		actual := fmt.Sprintf("%v", test.input)
		if actual != test.expected {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.expected, actual, test.expected)
		}
	}
}
