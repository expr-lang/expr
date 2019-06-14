package parser_test

import (
	"fmt"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
	"gopkg.in/antonmedv/expr.v2/ast"
	"gopkg.in/antonmedv/expr.v2/parser"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	parseTests := []struct {
		input    string
		expected ast.Node
	}{
		{
			"a",
			&ast.IdentifierNode{Value: "a"},
		},
		{
			`"\"double\""`,
			&ast.StringNode{Value: "\"double\""},
		},
		{
			`'\'single\\ \''`,
			&ast.StringNode{Value: "'single\\ '"},
		},
		{
			"3",
			&ast.IntegerNode{Value: 3},
		},
		{
			"0xFF",
			&ast.IntegerNode{Value: 255},
		},
		{
			"10_000_000",
			&ast.IntegerNode{Value: 10000000},
		},
		{
			"2.5",
			&ast.FloatNode{Value: 2.5},
		},
		{
			"true",
			&ast.BoolNode{Value: true},
		},
		{
			"false",
			&ast.BoolNode{},
		},
		{
			"nil",
			&ast.NilNode{},
		},
		{
			"-3",
			&ast.UnaryNode{Operator: "-", Node: &ast.IntegerNode{Value: 3}},
		},
		{
			"1 - 2",
			&ast.BinaryNode{Operator: "-", Left: &ast.IntegerNode{Value: 1}, Right: &ast.IntegerNode{Value: 2}},
		},
		{
			"(1 - 2) * 3",
			&ast.BinaryNode{Operator: "*",
				Left: &ast.BinaryNode{Operator: "-", Left: &ast.IntegerNode{Value: 1},
					Right: &ast.IntegerNode{Value: 2}}, Right: &ast.IntegerNode{Value: 3},
			},
		},
		{
			"a or b or c",
			&ast.BinaryNode{Operator: "or", Left: &ast.BinaryNode{Operator: "or", Left: &ast.IdentifierNode{Value: "a"}, Right: &ast.IdentifierNode{Value: "b"}}, Right: &ast.IdentifierNode{Value: "c"}},
		},
		{
			"a or b and c",
			&ast.BinaryNode{Operator: "or", Left: &ast.IdentifierNode{Value: "a"}, Right: &ast.BinaryNode{Operator: "and", Left: &ast.IdentifierNode{Value: "b"}, Right: &ast.IdentifierNode{Value: "c"}}},
		},
		{
			"(a or b) and c",
			&ast.BinaryNode{Operator: "and", Left: &ast.BinaryNode{Operator: "or", Left: &ast.IdentifierNode{Value: "a"}, Right: &ast.IdentifierNode{Value: "b"}}, Right: &ast.IdentifierNode{Value: "c"}},
		},
		{
			"2**4-1",
			&ast.BinaryNode{Operator: "-", Left: &ast.BinaryNode{Operator: "**", Left: &ast.IntegerNode{Value: 2}, Right: &ast.IntegerNode{Value: 4}}, Right: &ast.IntegerNode{Value: 1}},
		},
		{
			"foo(bar())",
			&ast.FunctionNode{Name: "foo", Arguments: []ast.Node{&ast.FunctionNode{Name: "bar", Arguments: []ast.Node{}}}},
		},
		{
			`foo("arg1", 2, true)`,
			&ast.FunctionNode{Name: "foo", Arguments: []ast.Node{&ast.StringNode{Value: "arg1"}, &ast.IntegerNode{Value: 2}, &ast.BoolNode{Value: true}}},
		},
		{
			"foo.bar",
			&ast.PropertyNode{Node: &ast.IdentifierNode{Value: "foo"}, Property: "bar"},
		},
		{
			"foo['all']",
			&ast.IndexNode{Node: &ast.IdentifierNode{Value: "foo"}, Index: &ast.StringNode{Value: "all"}},
		},
		{
			"foo.bar()",
			&ast.MethodNode{Node: &ast.IdentifierNode{Value: "foo"}, Method: "bar", Arguments: []ast.Node{}},
		},
		{
			`foo.bar("arg1", 2, true)`,
			&ast.MethodNode{Node: &ast.IdentifierNode{Value: "foo"}, Method: "bar", Arguments: []ast.Node{&ast.StringNode{Value: "arg1"}, &ast.IntegerNode{Value: 2}, &ast.BoolNode{Value: true}}},
		},
		{
			"foo[3]",
			&ast.IndexNode{Node: &ast.IdentifierNode{Value: "foo"}, Index: &ast.IntegerNode{Value: 3}},
		},
		{
			"true ? true : false",
			&ast.ConditionalNode{Cond: &ast.BoolNode{Value: true}, Exp1: &ast.BoolNode{Value: true}, Exp2: &ast.BoolNode{}},
		},
		{
			"foo.bar().foo().baz[33]",
			&ast.IndexNode{
				Node: &ast.PropertyNode{Node: &ast.MethodNode{Node: &ast.MethodNode{
					Node: &ast.IdentifierNode{Value: "foo"}, Method: "bar", Arguments: []ast.Node{}}, Method: "foo", Arguments: []ast.Node{}}, Property: "baz"},
				Index: &ast.IntegerNode{Value: 33},
			},
		},
		{
			"'a' == 'b'",
			&ast.BinaryNode{Operator: "==", Left: &ast.StringNode{Value: "a"}, Right: &ast.StringNode{Value: "b"}},
		},
		{
			"+0 != -0",
			&ast.BinaryNode{Operator: "!=", Left: &ast.UnaryNode{Operator: "+", Node: &ast.IntegerNode{}}, Right: &ast.UnaryNode{Operator: "-", Node: &ast.IntegerNode{}}},
		},
		{
			"[a, b, c]",
			&ast.ArrayNode{Nodes: []ast.Node{&ast.IdentifierNode{Value: "a"}, &ast.IdentifierNode{Value: "b"}, &ast.IdentifierNode{Value: "c"}}},
		},
		{
			"{foo:1, bar:2}",
			&ast.MapNode{Pairs: []*ast.PairNode{{Key: &ast.StringNode{Value: "foo"}, Value: &ast.IntegerNode{Value: 1}}, {Key: &ast.StringNode{Value: "bar"}, Value: &ast.IntegerNode{Value: 2}}}},
		},
		{
			`{"a": 1, 'b': 2}`,
			&ast.MapNode{Pairs: []*ast.PairNode{{Key: &ast.StringNode{Value: "a"}, Value: &ast.IntegerNode{Value: 1}}, {Key: &ast.StringNode{Value: "b"}, Value: &ast.IntegerNode{Value: 2}}}},
		},
		{
			"[1].foo",
			&ast.PropertyNode{Node: &ast.ArrayNode{Nodes: []ast.Node{&ast.IntegerNode{Value: 1}}}, Property: "foo"},
		},
		{
			"{foo:1}.bar",
			&ast.PropertyNode{Node: &ast.MapNode{Pairs: []*ast.PairNode{{Key: &ast.StringNode{Value: "foo"}, Value: &ast.IntegerNode{Value: 1}}}}, Property: "bar"},
		},
		{
			"len(foo)",
			&ast.BuiltinNode{Name: "len", Arguments: []ast.Node{&ast.IdentifierNode{Value: "foo"}}},
		},
		{
			`foo matches "foo"`,
			&ast.MatchesNode{Left: &ast.IdentifierNode{Value: "foo"}, Right: &ast.StringNode{Value: "foo"}},
		},
		{
			`foo matches regex`,
			&ast.MatchesNode{Left: &ast.IdentifierNode{Value: "foo"}, Right: &ast.IdentifierNode{Value: "regex"}},
		},
		{
			`foo contains "foo"`,
			&ast.BinaryNode{Operator: "contains", Left: &ast.IdentifierNode{Value: "foo"}, Right: &ast.StringNode{Value: "foo"}},
		},
		{
			`foo startsWith "foo"`,
			&ast.BinaryNode{Operator: "startsWith", Left: &ast.IdentifierNode{Value: "foo"}, Right: &ast.StringNode{Value: "foo"}},
		},
		{
			`foo endsWith "foo"`,
			&ast.BinaryNode{Operator: "endsWith", Left: &ast.IdentifierNode{Value: "foo"}, Right: &ast.StringNode{Value: "foo"}},
		},
		{
			"1..9",
			&ast.BinaryNode{Operator: "..", Left: &ast.IntegerNode{Value: 1}, Right: &ast.IntegerNode{Value: 9}},
		},
		{
			"0 in []",
			&ast.BinaryNode{Operator: "in", Left: &ast.IntegerNode{}, Right: &ast.ArrayNode{Nodes: []ast.Node{}}},
		},
		{
			"all(Tickets, {.Price > 0})",
			&ast.BuiltinNode{Name: "all", Arguments: []ast.Node{&ast.IdentifierNode{Value: "Tickets"}, &ast.ClosureNode{Node: &ast.BinaryNode{Operator: ">", Left: &ast.PropertyNode{Node: &ast.PointerNode{}, Property: "Price"}, Right: &ast.IntegerNode{Value: 0}}}}},
		},
	}
	for _, test := range parseTests {
		actual, err := parser.Parse(test.input)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		if m, ok := (actual.Node).(*ast.MatchesNode); ok {
			m.Regexp = nil
			actual.Node = m
		}
		assert.Equal(t, litter.Sdump(test.expected), litter.Sdump(actual.Node), test.input)
	}
}

func TestParse_error(t *testing.T) {
	var parseErrorTests = []struct {
		input string
		err   string
	}{
		{
			"foo.",
			"syntax error: missing Identifier at",
		},
		{
			"a+",
			"syntax error: mismatched input '<EOF>'",
		},
		{
			"a ? (1+2) c",
			"syntax error: missing ':' at 'c'",
		},
		{
			"[a b]",
			"syntax error: extraneous input 'b' expecting {']', ','}",
		},
		{
			"foo.bar(a b)",
			"syntax error: extraneous input 'b' expecting ')'",
		},
		{
			"{-}",
			"syntax error: no viable alternative at input '{-'",
		},
		{
			"a matches 'a)(b'",
			"error parsing regexp: unexpected )",
		},
		{
			`a matches "*"`,
			"error parsing regexp: missing argument to repetition operator: `*` (1:11)\n | a matches \"*\"\n | ..........^",
		},
		{
			`.foo`,
			"parse error: dot property accessor can be only inside closure",
		},
		{
			`foo({.bar})`,
			"syntax error: no viable alternative at input '{.'",
		},
	}
	for _, test := range parseErrorTests {
		_, err := parser.Parse(test.input)
		if err == nil {
			err = fmt.Errorf("<nil>")
		}
		if !strings.Contains(err.Error(), test.err) || test.err == "" {
			assert.Equal(t, test.err, err.Error(), test.input)
		}
	}
}
