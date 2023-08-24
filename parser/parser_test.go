package parser_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  Node
	}{
		{
			"a",
			&IdentifierNode{Value: "a"},
		},
		{
			`"str"`,
			&StringNode{Value: "str"},
		},
		{
			"3",
			&IntegerNode{Value: 3},
		},
		{
			"0xFF",
			&IntegerNode{Value: 255},
		},
		{
			"0x6E",
			&IntegerNode{Value: 110},
		},
		{
			"10_000_000",
			&IntegerNode{Value: 10_000_000},
		},
		{
			"2.5",
			&FloatNode{Value: 2.5},
		},
		{
			"1e9",
			&FloatNode{Value: 1e9},
		},
		{
			"true",
			&BoolNode{Value: true},
		},
		{
			"false",
			&BoolNode{Value: false},
		},
		{
			"nil",
			&NilNode{},
		},
		{
			"-3",
			&UnaryNode{Operator: "-",
				Node: &IntegerNode{Value: 3}},
		},
		{
			"-2^2",
			&UnaryNode{
				Operator: "-",
				Node: &BinaryNode{
					Operator: "^",
					Left:     &IntegerNode{Value: 2},
					Right:    &IntegerNode{Value: 2},
				},
			},
		},
		{
			"1 - 2",
			&BinaryNode{Operator: "-",
				Left:  &IntegerNode{Value: 1},
				Right: &IntegerNode{Value: 2}},
		},
		{
			"(1 - 2) * 3",
			&BinaryNode{
				Operator: "*",
				Left: &BinaryNode{
					Operator: "-",
					Left:     &IntegerNode{Value: 1},
					Right:    &IntegerNode{Value: 2},
				},
				Right: &IntegerNode{Value: 3},
			},
		},
		{
			"a or b or c",
			&BinaryNode{Operator: "or",
				Left: &BinaryNode{Operator: "or",
					Left:  &IdentifierNode{Value: "a"},
					Right: &IdentifierNode{Value: "b"}},
				Right: &IdentifierNode{Value: "c"}},
		},
		{
			"a or b and c",
			&BinaryNode{Operator: "or",
				Left: &IdentifierNode{Value: "a"},
				Right: &BinaryNode{Operator: "and",
					Left:  &IdentifierNode{Value: "b"},
					Right: &IdentifierNode{Value: "c"}}},
		},
		{
			"(a or b) and c",
			&BinaryNode{Operator: "and",
				Left: &BinaryNode{Operator: "or",
					Left:  &IdentifierNode{Value: "a"},
					Right: &IdentifierNode{Value: "b"}},
				Right: &IdentifierNode{Value: "c"}},
		},
		{
			"2**4-1",
			&BinaryNode{Operator: "-",
				Left: &BinaryNode{Operator: "**",
					Left:  &IntegerNode{Value: 2},
					Right: &IntegerNode{Value: 4}},
				Right: &IntegerNode{Value: 1}},
		},
		{
			"foo(bar())",
			&CallNode{Callee: &IdentifierNode{Value: "foo"},
				Arguments: []Node{&CallNode{Callee: &IdentifierNode{Value: "bar"},
					Arguments: []Node{}}}},
		},
		{
			`foo("arg1", 2, true)`,
			&CallNode{Callee: &IdentifierNode{Value: "foo"},
				Arguments: []Node{&StringNode{Value: "arg1"},
					&IntegerNode{Value: 2},
					&BoolNode{Value: true}}},
		},
		{
			"foo.bar",
			&MemberNode{Node: &IdentifierNode{Value: "foo"},
				Property: &StringNode{Value: "bar"}},
		},
		{
			"foo['all']",
			&MemberNode{Node: &IdentifierNode{Value: "foo"},
				Property: &StringNode{Value: "all"}},
		},
		{
			"foo.bar()",
			&CallNode{Callee: &MemberNode{Node: &IdentifierNode{Value: "foo"},
				Property: &StringNode{Value: "bar"}},
				Arguments: []Node{}},
		},
		{
			`foo.bar("arg1", 2, true)`,
			&CallNode{Callee: &MemberNode{Node: &IdentifierNode{Value: "foo"},
				Property: &StringNode{Value: "bar"}},
				Arguments: []Node{&StringNode{Value: "arg1"},
					&IntegerNode{Value: 2},
					&BoolNode{Value: true}}},
		},
		{
			"foo[3]",
			&MemberNode{Node: &IdentifierNode{Value: "foo"},
				Property: &IntegerNode{Value: 3}},
		},
		{
			"true ? true : false",
			&ConditionalNode{Cond: &BoolNode{Value: true},
				Exp1: &BoolNode{Value: true},
				Exp2: &BoolNode{}},
		},
		{
			"a?[b]:c",
			&ConditionalNode{Cond: &IdentifierNode{Value: "a"},
				Exp1: &ArrayNode{Nodes: []Node{&IdentifierNode{Value: "b"}}},
				Exp2: &IdentifierNode{Value: "c"}},
		},
		{
			"a.b().c().d[33]",
			&MemberNode{
				Node: &MemberNode{
					Node: &CallNode{
						Callee: &MemberNode{
							Node: &CallNode{
								Callee: &MemberNode{
									Node: &IdentifierNode{
										Value: "a",
									},
									Property: &StringNode{
										Value: "b",
									},
								},
								Arguments: []Node{},
								Fast:      false,
							},
							Property: &StringNode{
								Value: "c",
							},
						},
						Arguments: []Node{},
						Fast:      false,
					},
					Property: &StringNode{
						Value: "d",
					},
				},
				Property: &IntegerNode{Value: 33}},
		},
		{
			"'a' == 'b'",
			&BinaryNode{Operator: "==",
				Left:  &StringNode{Value: "a"},
				Right: &StringNode{Value: "b"}},
		},
		{
			"+0 != -0",
			&BinaryNode{Operator: "!=",
				Left: &UnaryNode{Operator: "+",
					Node: &IntegerNode{}},
				Right: &UnaryNode{Operator: "-",
					Node: &IntegerNode{}}},
		},
		{
			"[a, b, c]",
			&ArrayNode{Nodes: []Node{&IdentifierNode{Value: "a"},
				&IdentifierNode{Value: "b"},
				&IdentifierNode{Value: "c"}}},
		},
		{
			"{foo:1, bar:2}",
			&MapNode{Pairs: []Node{&PairNode{Key: &StringNode{Value: "foo"},
				Value: &IntegerNode{Value: 1}},
				&PairNode{Key: &StringNode{Value: "bar"},
					Value: &IntegerNode{Value: 2}}}},
		},
		{
			"{foo:1, bar:2, }",
			&MapNode{Pairs: []Node{&PairNode{Key: &StringNode{Value: "foo"},
				Value: &IntegerNode{Value: 1}},
				&PairNode{Key: &StringNode{Value: "bar"},
					Value: &IntegerNode{Value: 2}}}},
		},
		{
			`{"a": 1, 'b': 2}`,
			&MapNode{Pairs: []Node{&PairNode{Key: &StringNode{Value: "a"},
				Value: &IntegerNode{Value: 1}},
				&PairNode{Key: &StringNode{Value: "b"},
					Value: &IntegerNode{Value: 2}}}},
		},
		{
			"[1].foo",
			&MemberNode{Node: &ArrayNode{Nodes: []Node{&IntegerNode{Value: 1}}},
				Property: &StringNode{Value: "foo"}},
		},
		{
			"{foo:1}.bar",
			&MemberNode{Node: &MapNode{Pairs: []Node{&PairNode{Key: &StringNode{Value: "foo"},
				Value: &IntegerNode{Value: 1}}}},
				Property: &StringNode{Value: "bar"}},
		},
		{
			"len(foo)",
			&BuiltinNode{
				Name: "len",
				Arguments: []Node{
					&IdentifierNode{Value: "foo"},
				},
			},
		},
		{
			`foo matches "foo"`,
			&BinaryNode{
				Operator: "matches",
				Left:     &IdentifierNode{Value: "foo"},
				Right:    &StringNode{Value: "foo"}},
		},
		{
			`foo not matches "foo"`,
			&UnaryNode{
				Operator: "not",
				Node: &BinaryNode{
					Operator: "matches",
					Left:     &IdentifierNode{Value: "foo"},
					Right:    &StringNode{Value: "foo"}}},
		},
		{
			`foo matches regex`,
			&BinaryNode{
				Operator: "matches",
				Left:     &IdentifierNode{Value: "foo"},
				Right:    &IdentifierNode{Value: "regex"}},
		},
		{
			`foo contains "foo"`,
			&BinaryNode{
				Operator: "contains",
				Left:     &IdentifierNode{Value: "foo"},
				Right:    &StringNode{Value: "foo"}},
		},
		{
			`foo not contains "foo"`,
			&UnaryNode{
				Operator: "not",
				Node: &BinaryNode{Operator: "contains",
					Left:  &IdentifierNode{Value: "foo"},
					Right: &StringNode{Value: "foo"}}},
		},
		{
			`foo startsWith "foo"`,
			&BinaryNode{Operator: "startsWith",
				Left:  &IdentifierNode{Value: "foo"},
				Right: &StringNode{Value: "foo"}},
		},
		{
			`foo endsWith "foo"`,
			&BinaryNode{Operator: "endsWith",
				Left:  &IdentifierNode{Value: "foo"},
				Right: &StringNode{Value: "foo"}},
		},
		{
			"1..9",
			&BinaryNode{Operator: "..",
				Left:  &IntegerNode{Value: 1},
				Right: &IntegerNode{Value: 9}},
		},
		{
			"0 in []",
			&BinaryNode{Operator: "in",
				Left:  &IntegerNode{},
				Right: &ArrayNode{Nodes: []Node{}}},
		},
		{
			"not in_var",
			&UnaryNode{Operator: "not",
				Node: &IdentifierNode{Value: "in_var"}},
		},
		{
			"all(Tickets, #)",
			&BuiltinNode{
				Name: "all",
				Arguments: []Node{
					&IdentifierNode{Value: "Tickets"},
					&ClosureNode{
						Node: &PointerNode{},
					}}},
		},
		{
			"all(Tickets, {.Price > 0})",
			&BuiltinNode{
				Name: "all",
				Arguments: []Node{
					&IdentifierNode{Value: "Tickets"},
					&ClosureNode{
						Node: &BinaryNode{
							Operator: ">",
							Left: &MemberNode{Node: &PointerNode{},
								Property: &StringNode{Value: "Price"}},
							Right: &IntegerNode{Value: 0}}}}},
		},
		{
			"one(Tickets, {#.Price > 0})",
			&BuiltinNode{
				Name: "one",
				Arguments: []Node{
					&IdentifierNode{Value: "Tickets"},
					&ClosureNode{
						Node: &BinaryNode{
							Operator: ">",
							Left: &MemberNode{
								Node:     &PointerNode{},
								Property: &StringNode{Value: "Price"},
							},
							Right: &IntegerNode{Value: 0}}}}},
		},
		{
			"filter(Prices, {# > 100})",
			&BuiltinNode{Name: "filter",
				Arguments: []Node{&IdentifierNode{Value: "Prices"},
					&ClosureNode{Node: &BinaryNode{Operator: ">",
						Left:  &PointerNode{},
						Right: &IntegerNode{Value: 100}}}}},
		},
		{
			"array[1:2]",
			&SliceNode{Node: &IdentifierNode{Value: "array"},
				From: &IntegerNode{Value: 1},
				To:   &IntegerNode{Value: 2}},
		},
		{
			"array[:2]",
			&SliceNode{Node: &IdentifierNode{Value: "array"},
				To: &IntegerNode{Value: 2}},
		},
		{
			"array[1:]",
			&SliceNode{Node: &IdentifierNode{Value: "array"},
				From: &IntegerNode{Value: 1}},
		},
		{
			"array[:]",
			&SliceNode{Node: &IdentifierNode{Value: "array"}},
		},
		{
			"[]",
			&ArrayNode{},
		},
		{
			"foo ?? bar",
			&BinaryNode{Operator: "??",
				Left:  &IdentifierNode{Value: "foo"},
				Right: &IdentifierNode{Value: "bar"}},
		},
		{
			"foo ?? bar ?? baz",
			&BinaryNode{Operator: "??",
				Left: &BinaryNode{Operator: "??",
					Left:  &IdentifierNode{Value: "foo"},
					Right: &IdentifierNode{Value: "bar"}},
				Right: &IdentifierNode{Value: "baz"}},
		},
		{
			"foo ?? (bar || baz)",
			&BinaryNode{Operator: "??",
				Left: &IdentifierNode{Value: "foo"},
				Right: &BinaryNode{Operator: "||",
					Left:  &IdentifierNode{Value: "bar"},
					Right: &IdentifierNode{Value: "baz"}}},
		},
		{
			"foo || bar ?? baz",
			&BinaryNode{Operator: "||",
				Left: &IdentifierNode{Value: "foo"},
				Right: &BinaryNode{Operator: "??",
					Left:  &IdentifierNode{Value: "bar"},
					Right: &IdentifierNode{Value: "baz"}}},
		},
		{
			"foo ?? bar()",
			&BinaryNode{Operator: "??",
				Left:  &IdentifierNode{Value: "foo"},
				Right: &CallNode{Callee: &IdentifierNode{Value: "bar"}}},
		},
		{
			"true | ok()",
			&CallNode{
				Callee: &IdentifierNode{Value: "ok"},
				Arguments: []Node{
					&BoolNode{Value: true}}}},
		{
			`let foo = a + b; foo + c`,
			&VariableDeclaratorNode{
				Name: "foo",
				Value: &BinaryNode{Operator: "+",
					Left:  &IdentifierNode{Value: "a"},
					Right: &IdentifierNode{Value: "b"}},
				Expr: &BinaryNode{Operator: "+",
					Left:  &IdentifierNode{Value: "foo"},
					Right: &IdentifierNode{Value: "c"}}},
		},
		{
			`map([], #index)`,
			&BuiltinNode{
				Name: "map",
				Arguments: []Node{
					&ArrayNode{},
					&ClosureNode{
						Node: &PointerNode{Name: "index"},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual, err := parser.Parse(test.input)
			require.NoError(t, err)
			assert.Equal(t, Dump(test.want), Dump(actual.Node))
		})
	}
}

const errorTests = `
foo.
unexpected end of expression (1:4)
 | foo.
 | ...^

a+
unexpected token EOF (1:2)
 | a+
 | .^

a ? (1+2) c
unexpected token Identifier("c") (1:11)
 | a ? (1+2) c
 | ..........^

[a b]
unexpected token Identifier("b") (1:4)
 | [a b]
 | ...^

foo.bar(a b)
unexpected token Identifier("b") (1:11)
 | foo.bar(a b)
 | ..........^

{-}
a map key must be a quoted string, a number, a identifier, or an expression enclosed in parentheses (unexpected token Operator("-")) (1:2)
 | {-}
 | .^

foo({.bar})
a map key must be a quoted string, a number, a identifier, or an expression enclosed in parentheses (unexpected token Operator(".")) (1:6)
 | foo({.bar})
 | .....^

.foo
cannot use pointer accessor outside closure (1:1)
 | .foo
 | ^

[1, 2, 3,,]
unexpected token Operator(",") (1:10)
 | [1, 2, 3,,]
 | .........^

[,]
unexpected token Operator(",") (1:2)
 | [,]
 | .^

{,}
a map key must be a quoted string, a number, a identifier, or an expression enclosed in parentheses (unexpected token Operator(",")) (1:2)
 | {,}
 | .^

{foo:1, bar:2, ,}
unexpected token Operator(",") (1:16)
 | {foo:1, bar:2, ,}
 | ...............^

foo ?? bar || baz
Operator (||) and coalesce expressions (??) cannot be mixed. Wrap either by parentheses. (1:12)
 | foo ?? bar || baz
 | ...........^
`

func TestParse_error(t *testing.T) {
	tests := strings.Split(strings.Trim(errorTests, "\n"), "\n\n")
	for _, test := range tests {
		input := strings.SplitN(test, "\n", 2)
		if len(input) != 2 {
			t.Errorf("syntax error in test: %q", test)
			break
		}
		_, err := parser.Parse(input[0])
		if err == nil {
			err = fmt.Errorf("<nil>")
		}
		assert.Equal(t, input[1], err.Error(), input[0])
	}
}

func TestParse_optional_chaining(t *testing.T) {
	parseTests := []struct {
		input    string
		expected Node
	}{
		{
			"foo?.bar.baz",
			&ChainNode{
				Node: &MemberNode{
					Node: &MemberNode{
						Node:     &IdentifierNode{Value: "foo"},
						Property: &StringNode{Value: "bar"},
						Optional: true,
					},
					Property: &StringNode{Value: "baz"},
				},
			},
		},
		{
			"foo.bar?.baz",
			&ChainNode{
				Node: &MemberNode{
					Node: &MemberNode{
						Node:     &IdentifierNode{Value: "foo"},
						Property: &StringNode{Value: "bar"},
					},
					Property: &StringNode{Value: "baz"},
					Optional: true,
				},
			},
		},
		{
			"foo?.bar?.baz",
			&ChainNode{
				Node: &MemberNode{
					Node: &MemberNode{
						Node:     &IdentifierNode{Value: "foo"},
						Property: &StringNode{Value: "bar"},
						Optional: true,
					},
					Property: &StringNode{Value: "baz"},
					Optional: true,
				},
			},
		},
		{
			"!foo?.bar.baz",
			&UnaryNode{
				Operator: "!",
				Node: &ChainNode{
					Node: &MemberNode{
						Node: &MemberNode{
							Node:     &IdentifierNode{Value: "foo"},
							Property: &StringNode{Value: "bar"},
							Optional: true,
						},
						Property: &StringNode{Value: "baz"},
					},
				},
			},
		},
		{
			"foo.bar[a?.b]?.baz",
			&ChainNode{
				Node: &MemberNode{
					Node: &MemberNode{
						Node: &MemberNode{
							Node:     &IdentifierNode{Value: "foo"},
							Property: &StringNode{Value: "bar"},
						},
						Property: &ChainNode{
							Node: &MemberNode{
								Node:     &IdentifierNode{Value: "a"},
								Property: &StringNode{Value: "b"},
								Optional: true,
							},
						},
					},
					Property: &StringNode{Value: "baz"},
					Optional: true,
				},
			},
		},
	}
	for _, test := range parseTests {
		actual, err := parser.Parse(test.input)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		assert.Equal(t, Dump(test.expected), Dump(actual.Node), test.input)
	}
}

func TestParse_pipe_operator(t *testing.T) {
	input := "arr | map(.foo) | len() | Foo()"
	expect := &CallNode{
		Callee: &IdentifierNode{Value: "Foo"},
		Arguments: []Node{
			&BuiltinNode{
				Name: "len",
				Arguments: []Node{
					&BuiltinNode{
						Name: "map",
						Arguments: []Node{
							&IdentifierNode{Value: "arr"},
							&ClosureNode{
								Node: &MemberNode{
									Node:     &PointerNode{},
									Property: &StringNode{Value: "foo"},
								}}}}}}}}

	actual, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, Dump(expect), Dump(actual.Node))
}
