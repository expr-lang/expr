package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"expr/internal/testify/assert"
	"expr/internal/testify/require"

	. "expr/ast"
	"expr/parser"
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
			"`hello\nworld`",
			&StringNode{Value: `hello
world`},
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
			"0X63",
			&IntegerNode{Value: 99},
		},
		{
			"0o600",
			&IntegerNode{Value: 384},
		},
		{
			"0O45",
			&IntegerNode{Value: 37},
		},
		{
			"0b10",
			&IntegerNode{Value: 2},
		},
		{
			"0B101011",
			&IntegerNode{Value: 43},
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
				Property: &StringNode{Value: "bar"}, Method: true},
				Arguments: []Node{}},
		},
		{
			`foo.bar("arg1", 2, true)`,
			&CallNode{Callee: &MemberNode{Node: &IdentifierNode{Value: "foo"},
				Property: &StringNode{Value: "bar"}, Method: true},
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
									Method: true,
								},
								Arguments: []Node{},
							},
							Property: &StringNode{
								Value: "c",
							},
							Method: true,
						},
						Arguments: []Node{},
					},
					Property: &StringNode{
						Value: "d",
					},
				},
				Property: &IntegerNode{Value: 33}},
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
			`::split("a,b,c", ",")`,
			&BuiltinNode{
				Name: "split",
				Arguments: []Node{
					&StringNode{Value: "a,b,c"},
					&StringNode{Value: ","},
				},
			},
		},
		{
			`::split("a,b,c", ",")[0]`,
			&MemberNode{
				Node: &BuiltinNode{
					Name: "split",
					Arguments: []Node{
						&StringNode{Value: "a,b,c"},
						&StringNode{Value: ","},
					},
				},
				Property: &IntegerNode{Value: 0},
			},
		},
		{
			`"hello"[1:3]`,
			&SliceNode{
				Node: &StringNode{Value: "hello"},
				From: &IntegerNode{Value: 1},
				To:   &IntegerNode{Value: 3},
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
cannot use pointer accessor outside predicate (1:1)
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

0b15
bad number syntax: "0b15" (1:4)
 | 0b15
 | ...^

0X10G
bad number syntax: "0X10G" (1:5)
 | 0X10G
 | ....^

0o1E
invalid float literal: strconv.ParseFloat: parsing "0o1E": invalid syntax (1:4)
 | 0o1E
 | ...^

0b1E
invalid float literal: strconv.ParseFloat: parsing "0b1E": invalid syntax (1:4)
 | 0b1E
 | ...^

0b1E+6
bad number syntax: "0b1E+6" (1:6)
 | 0b1E+6
 | .....^

0b1E+1
invalid float literal: strconv.ParseFloat: parsing "0b1E+1": invalid syntax (1:6)
 | 0b1E+1
 | .....^

0o1E+1
invalid float literal: strconv.ParseFloat: parsing "0o1E+1": invalid syntax (1:6)
 | 0o1E+1
 | .....^

1E
invalid float literal: strconv.ParseFloat: parsing "1E": invalid syntax (1:2)
 | 1E
 | .^

1 not == [1, 2, 5]
unexpected token Operator("==") (1:7)
 | 1 not == [1, 2, 5]
 | ......^
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
		{
			"foo.bar?.[0]",
			&ChainNode{
				Node: &MemberNode{
					Node: &MemberNode{
						Node:     &IdentifierNode{Value: "foo"},
						Property: &StringNode{Value: "bar"},
					},
					Property: &IntegerNode{Value: 0},
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
