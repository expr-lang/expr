package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"

	. "github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
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
			"-1 not in [1, 2, 3, 4]",
			&UnaryNode{Operator: "not",
				Node: &BinaryNode{Operator: "in",
					Left: &UnaryNode{Operator: "-", Node: &IntegerNode{Value: 1}},
					Right: &ArrayNode{Nodes: []Node{
						&IntegerNode{Value: 1},
						&IntegerNode{Value: 2},
						&IntegerNode{Value: 3},
						&IntegerNode{Value: 4},
					}}}},
		},
		{
			"1*8 not in [1, 2, 3, 4]",
			&UnaryNode{Operator: "not",
				Node: &BinaryNode{Operator: "in",
					Left: &BinaryNode{Operator: "*",
						Left:  &IntegerNode{Value: 1},
						Right: &IntegerNode{Value: 8},
					},
					Right: &ArrayNode{Nodes: []Node{
						&IntegerNode{Value: 1},
						&IntegerNode{Value: 2},
						&IntegerNode{Value: 3},
						&IntegerNode{Value: 4},
					}}}},
		},
		{
			"2==2 ? false : 3 not in [1, 2, 5]",
			&ConditionalNode{
				Cond: &BinaryNode{
					Operator: "==",
					Left:     &IntegerNode{Value: 2},
					Right:    &IntegerNode{Value: 2},
				},
				Exp1: &BoolNode{Value: false},
				Exp2: &UnaryNode{
					Operator: "not",
					Node: &BinaryNode{
						Operator: "in",
						Left:     &IntegerNode{Value: 3},
						Right: &ArrayNode{Nodes: []Node{
							&IntegerNode{Value: 1},
							&IntegerNode{Value: 2},
							&IntegerNode{Value: 5},
						}}}}},
		},
		{
			"'foo' + 'bar' not matches 'foobar'",
			&UnaryNode{Operator: "not",
				Node: &BinaryNode{Operator: "matches",
					Left: &BinaryNode{Operator: "+",
						Left:  &StringNode{Value: "foo"},
						Right: &StringNode{Value: "bar"}},
					Right: &StringNode{Value: "foobar"}}},
		},
		{
			"all(Tickets, #)",
			&BuiltinNode{
				Name: "all",
				Arguments: []Node{
					&IdentifierNode{Value: "Tickets"},
					&PredicateNode{
						Node: &PointerNode{},
					}}},
		},
		{
			"all(Tickets, {.Price > 0})",
			&BuiltinNode{
				Name: "all",
				Arguments: []Node{
					&IdentifierNode{Value: "Tickets"},
					&PredicateNode{
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
					&PredicateNode{
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
					&PredicateNode{Node: &BinaryNode{Operator: ">",
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
					&PredicateNode{
						Node: &PointerNode{Name: "index"},
					},
				},
			},
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
		{
			`1 < 2 > 3`,
			&BinaryNode{
				Operator: "&&",
				Left: &BinaryNode{
					Operator: "<",
					Left:     &IntegerNode{Value: 1},
					Right:    &IntegerNode{Value: 2},
				},
				Right: &BinaryNode{
					Operator: ">",
					Left:     &IntegerNode{Value: 2},
					Right:    &IntegerNode{Value: 3},
				},
			},
		},
		{
			`1 < 2 < 3 < 4`,
			&BinaryNode{
				Operator: "&&",
				Left: &BinaryNode{
					Operator: "&&",
					Left: &BinaryNode{
						Operator: "<",
						Left:     &IntegerNode{Value: 1},
						Right:    &IntegerNode{Value: 2},
					},
					Right: &BinaryNode{
						Operator: "<",
						Left:     &IntegerNode{Value: 2},
						Right:    &IntegerNode{Value: 3},
					},
				},
				Right: &BinaryNode{
					Operator: "<",
					Left:     &IntegerNode{Value: 3},
					Right:    &IntegerNode{Value: 4},
				},
			},
		},
		{
			`1 < 2 < 3 == true`,
			&BinaryNode{
				Operator: "==",
				Left: &BinaryNode{
					Operator: "&&",
					Left: &BinaryNode{
						Operator: "<",
						Left:     &IntegerNode{Value: 1},
						Right:    &IntegerNode{Value: 2},
					},
					Right: &BinaryNode{
						Operator: "<",
						Left:     &IntegerNode{Value: 2},
						Right:    &IntegerNode{Value: 3},
					},
				},
				Right: &BoolNode{Value: true},
			},
		},
		{
			"if a>b {true} else {x}",
			&ConditionalNode{
				Cond: &BinaryNode{
					Operator: ">",
					Left:     &IdentifierNode{Value: "a"},
					Right:    &IdentifierNode{Value: "b"},
				},
				Exp1: &BoolNode{Value: true},
				Exp2: &IdentifierNode{Value: "x"}},
		},
		{
			"1; 2; 3",
			&SequenceNode{
				Nodes: []Node{
					&IntegerNode{Value: 1},
					&IntegerNode{Value: 2},
					&IntegerNode{Value: 3},
				},
			},
		},
		{
			"1; (2; 3)",
			&SequenceNode{
				Nodes: []Node{
					&IntegerNode{Value: 1},
					&SequenceNode{
						Nodes: []Node{
							&IntegerNode{Value: 2},
							&IntegerNode{Value: 3}},
					},
				},
			},
		},
		{
			"true ? 1 : 2; 3 ; 4",
			&SequenceNode{
				Nodes: []Node{
					&ConditionalNode{
						Cond: &BoolNode{Value: true},
						Exp1: &IntegerNode{Value: 1},
						Exp2: &IntegerNode{Value: 2}},
					&IntegerNode{Value: 3},
					&IntegerNode{Value: 4},
				},
			},
		},
		{
			"true ? 1 : ( 2; 3; 4 )",
			&ConditionalNode{
				Cond: &BoolNode{Value: true},
				Exp1: &IntegerNode{Value: 1},
				Exp2: &SequenceNode{
					Nodes: []Node{
						&IntegerNode{Value: 2},
						&IntegerNode{Value: 3},
						&IntegerNode{Value: 4},
					},
				},
			},
		},
		{
			"true ?: 1; 2; 3",
			&SequenceNode{
				Nodes: []Node{
					&ConditionalNode{
						Cond: &BoolNode{Value: true},
						Exp1: &BoolNode{Value: true},
						Exp2: &IntegerNode{Value: 1}},
					&IntegerNode{Value: 2},
					&IntegerNode{Value: 3},
				},
			},
		},
		{
			`let x = true ? 1 : 2; x`,
			&VariableDeclaratorNode{
				Name: "x",
				Value: &ConditionalNode{
					Cond: &BoolNode{Value: true},
					Exp1: &IntegerNode{Value: 1},
					Exp2: &IntegerNode{Value: 2}},
				Expr: &IdentifierNode{Value: "x"}},
		},
		{
			"let x = true ? 1 : ( 2; 3; 4 ); x",
			&VariableDeclaratorNode{
				Name: "x",
				Value: &ConditionalNode{
					Cond: &BoolNode{Value: true},
					Exp1: &IntegerNode{Value: 1},
					Exp2: &SequenceNode{
						Nodes: []Node{
							&IntegerNode{Value: 2},
							&IntegerNode{Value: 3},
							&IntegerNode{Value: 4},
						},
					},
				},
				Expr: &IdentifierNode{Value: "x"}},
		},
		{
			"if true { 1; 2; 3 } else { 4; 5; 6 }",
			&ConditionalNode{
				Cond: &BoolNode{Value: true},
				Exp1: &SequenceNode{
					Nodes: []Node{
						&IntegerNode{Value: 1},
						&IntegerNode{Value: 2},
						&IntegerNode{Value: 3}}},
				Exp2: &SequenceNode{
					Nodes: []Node{
						&IntegerNode{Value: 4},
						&IntegerNode{Value: 5},
						&IntegerNode{Value: 6}}}},
		},
		{
			`all(ls, if true { 1 } else { 2 })`,
			&BuiltinNode{
				Name: "all",
				Arguments: []Node{
					&IdentifierNode{Value: "ls"},
					&PredicateNode{
						Node: &ConditionalNode{
							Cond: &BoolNode{Value: true},
							Exp1: &IntegerNode{Value: 1},
							Exp2: &IntegerNode{Value: 2}}}}},
		},
		{
			`let x = if true { 1 } else { 2 }; x`,
			&VariableDeclaratorNode{
				Name: "x",
				Value: &ConditionalNode{
					Cond: &BoolNode{Value: true},
					Exp1: &IntegerNode{Value: 1},
					Exp2: &IntegerNode{Value: 2}},
				Expr: &IdentifierNode{Value: "x"}},
		},
		{
			`call(if true { 1 } else { 2 })`,
			&CallNode{
				Callee: &IdentifierNode{Value: "call"},
				Arguments: []Node{
					&ConditionalNode{
						Cond: &BoolNode{Value: true},
						Exp1: &IntegerNode{Value: 1},
						Exp2: &IntegerNode{Value: 2}}}},
		},
		{
			`[if true { 1 } else { 2 }]`,
			&ArrayNode{
				Nodes: []Node{
					&ConditionalNode{
						Cond: &BoolNode{Value: true},
						Exp1: &IntegerNode{Value: 1},
						Exp2: &IntegerNode{Value: 2}}}},
		},
		{
			`map(ls, { 1; 2; 3 })`,
			&BuiltinNode{
				Name: "map",
				Arguments: []Node{
					&IdentifierNode{Value: "ls"},
					&PredicateNode{
						Node: &SequenceNode{
							Nodes: []Node{
								&IntegerNode{Value: 1},
								&IntegerNode{Value: 2},
								&IntegerNode{Value: 3},
							},
						},
					},
				}},
		},
		{
			`let x = 1; 2; 3 + x`,
			&VariableDeclaratorNode{
				Name:  "x",
				Value: &IntegerNode{Value: 1},
				Expr: &SequenceNode{
					Nodes: []Node{
						&IntegerNode{Value: 2},
						&BinaryNode{
							Operator: "+",
							Left:     &IntegerNode{Value: 3},
							Right:    &IdentifierNode{Value: "x"},
						},
					},
				},
			},
		},
		{
			`let x = 1; let y = 2; 3; 4; x + y`,
			&VariableDeclaratorNode{
				Name:  "x",
				Value: &IntegerNode{Value: 1},
				Expr: &VariableDeclaratorNode{
					Name:  "y",
					Value: &IntegerNode{Value: 2},
					Expr: &SequenceNode{
						Nodes: []Node{
							&IntegerNode{Value: 3},
							&IntegerNode{Value: 4},
							&BinaryNode{
								Operator: "+",
								Left:     &IdentifierNode{Value: "x"},
								Right:    &IdentifierNode{Value: "y"},
							},
						},
					}}},
		},
		{
			`let x = (1; 2; 3); x`,
			&VariableDeclaratorNode{
				Name: "x",
				Value: &SequenceNode{
					Nodes: []Node{
						&IntegerNode{Value: 1},
						&IntegerNode{Value: 2},
						&IntegerNode{Value: 3},
					},
				},
				Expr: &IdentifierNode{Value: "x"},
			},
		},
		{
			`all(
				[
				  true,
				  false,
				],
				#,
			)`,
			&BuiltinNode{
				Name: "all",
				Arguments: []Node{
					&ArrayNode{
						Nodes: []Node{
							&BoolNode{Value: true},
							&BoolNode{Value: false},
						},
					},
					&PredicateNode{
						Node: &PointerNode{},
					},
				},
			},
		},
		{
			`list | all(#,)`,
			&BuiltinNode{
				Name: "all",
				Arguments: []Node{
					&IdentifierNode{Value: "list"},
					&PredicateNode{
						Node: &PointerNode{},
					},
				},
			},
		},
		{
			`func(
				parameter1,
				parameter2,
			)`,
			&CallNode{
				Callee: &IdentifierNode{Value: "func"},
				Arguments: []Node{
					&IdentifierNode{Value: "parameter1"},
					&IdentifierNode{Value: "parameter2"},
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

func TestParse_error(t *testing.T) {
	var tests = []struct {
		input string
		err   string
	}{
		{`foo.`, `unexpected end of expression (1:4)
 | foo.
 | ...^`},
		{`a+`, `unexpected token EOF (1:2)
 | a+
 | .^`},
		{`a ? (1+2) c`, `unexpected token Identifier("c") (1:11)
 | a ? (1+2) c
 | ..........^`},
		{`[a b]`, `unexpected token Identifier("b") (1:4)
 | [a b]
 | ...^`},
		{`foo.bar(a b)`, `unexpected token Identifier("b") (1:11)
 | foo.bar(a b)
 | ..........^`},
		{`{-}`, `a map key must be a quoted string, a number, a identifier, or an expression enclosed in parentheses (unexpected token Operator("-")) (1:2)
 | {-}
 | .^`},
		{`foo({.bar})`, `a map key must be a quoted string, a number, a identifier, or an expression enclosed in parentheses (unexpected token Operator(".")) (1:6)
 | foo({.bar})
 | .....^`},
		{`[1, 2, 3,,]`, `unexpected token Operator(",") (1:10)
 | [1, 2, 3,,]
 | .........^`},
		{`[,]`, `unexpected token Operator(",") (1:2)
 | [,]
 | .^`},
		{`{,}`, `a map key must be a quoted string, a number, a identifier, or an expression enclosed in parentheses (unexpected token Operator(",")) (1:2)
 | {,}
 | .^`},
		{`{foo:1, bar:2, ,}`, `unexpected token Operator(",") (1:16)
 | {foo:1, bar:2, ,}
 | ...............^`},
		{`foo ?? bar || baz`, `Operator (||) and coalesce expressions (??) cannot be mixed. Wrap either by parentheses. (1:12)
 | foo ?? bar || baz
 | ...........^`},
		{`0b15`, `bad number syntax: "0b15" (1:4)
 | 0b15
 | ...^`},
		{`0X10G`, `bad number syntax: "0X10G" (1:5)
 | 0X10G
 | ....^`},
		{`0o1E`, `invalid float literal: strconv.ParseFloat: parsing "0o1E": invalid syntax (1:4)
 | 0o1E
 | ...^`},
		{`0b1E`, `invalid float literal: strconv.ParseFloat: parsing "0b1E": invalid syntax (1:4)
 | 0b1E
 | ...^`},
		{`0b1E+6`, `bad number syntax: "0b1E+6" (1:6)
 | 0b1E+6
 | .....^`},
		{`0b1E+1`, `invalid float literal: strconv.ParseFloat: parsing "0b1E+1": invalid syntax (1:6)
 | 0b1E+1
 | .....^`},
		{`0o1E+1`, `invalid float literal: strconv.ParseFloat: parsing "0o1E+1": invalid syntax (1:6)
 | 0o1E+1
 | .....^`},
		{`1E`, `invalid float literal: strconv.ParseFloat: parsing "1E": invalid syntax (1:2)
 | 1E
 | .^`},
		{`1 not == [1, 2, 5]`, `unexpected token Operator("==") (1:7)
 | 1 not == [1, 2, 5]
 | ......^`},
		{`foo(1; 2; 3)`, `unexpected token Operator(";") (1:6)
 | foo(1; 2; 3)
 | .....^`},
		{
			`map(ls, 1; 2; 3)`,
			`wrap predicate with brackets { and } (1:10)
 | map(ls, 1; 2; 3)
 | .........^`,
		},
		{
			`[1; 2; 3]`,
			`unexpected token Operator(";") (1:3)
 | [1; 2; 3]
 | ..^`,
		},
		{
			`1 + if true { 2 } else { 3 }`,
			`unexpected token Operator("if") (1:5)
 | 1 + if true { 2 } else { 3 }
 | ....^`,
		},
		{
			`list | all(#,,)`,
			`unexpected token Operator(",") (1:14)
 | list | all(#,,)
 | .............^`,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			_, err := parser.Parse(test.input)
			if err == nil {
				err = fmt.Errorf("<nil>")
			}
			assert.Equal(t, test.err, err.Error(), test.input)
		})
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
							&PredicateNode{
								Node: &MemberNode{
									Node:     &PointerNode{},
									Property: &StringNode{Value: "foo"},
								}}}}}}}}

	actual, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, Dump(expect), Dump(actual.Node))
}

func TestNodeBudget(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		maxNodes    uint
		shouldError bool
	}{
		{
			name:        "simple expression equal to limit",
			expr:        "a + b",
			maxNodes:    3,
			shouldError: false,
		},
		{
			name:        "medium expression under limit",
			expr:        "a + b * c / d",
			maxNodes:    20,
			shouldError: false,
		},
		{
			name:        "deeply nested expression over limit",
			expr:        "1 + (2 + (3 + (4 + (5 + (6 + (7 + 8))))))",
			maxNodes:    10,
			shouldError: true,
		},
		{
			name:        "array expression over limit",
			expr:        "[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]",
			maxNodes:    5,
			shouldError: true,
		},
		{
			name:        "disabled node budget",
			expr:        "1 + (2 + (3 + (4 + (5 + (6 + (7 + 8))))))",
			maxNodes:    0,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := conf.CreateNew()
			config.MaxNodes = tt.maxNodes
			config.Disabled = make(map[string]bool, 0)

			_, err := parser.ParseWithConfig(tt.expr, config)
			hasError := err != nil && strings.Contains(err.Error(), "exceeds maximum allowed nodes")

			if hasError != tt.shouldError {
				t.Errorf("ParseWithConfig(%q) error = %v, shouldError %v", tt.expr, err, tt.shouldError)
			}

			// Verify error message format when expected
			if tt.shouldError && err != nil {
				expected := "compilation failed: expression exceeds maximum allowed nodes"
				if !strings.Contains(err.Error(), expected) {
					t.Errorf("Expected error message to contain %q, got %q", expected, err.Error())
				}
			}
		})
	}
}

func TestNodeBudgetDisabled(t *testing.T) {
	config := conf.CreateNew()
	config.MaxNodes = 0 // Disable node budget

	expr := strings.Repeat("a + ", 1000) + "b"
	_, err := parser.ParseWithConfig(expr, config)

	if err != nil && strings.Contains(err.Error(), "exceeds maximum allowed nodes") {
		t.Error("Node budget check should be disabled when MaxNodes is 0")
	}
}
