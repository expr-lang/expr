package ast_test

import (
	"math"
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`nil`, `nil`},
		{`true`, `true`},
		{`false`, `false`},
		{`1`, `1`},
		{`1.1`, `1.1`},
		{`"a"`, `"a"`},
		{`'a'`, `"a"`},
		{`a`, `a`},
		{`a.b`, `a.b`},
		{`a[0]`, `a[0]`},
		{`a["the b"]`, `a["the b"]`},
		{`a.b[0]`, `a.b[0]`},
		{`a?.b`, `a?.b`},
		{`x[0][1]`, `x[0][1]`},
		{`x?.[0]?.[1]`, `x?.[0]?.[1]`},
		{`-a`, `-a`},
		{`!a`, `!a`},
		{`not a`, `not a`},
		{`a + b`, `a + b`},
		{`a + b * c`, `a + b * c`},
		{`(a + b) * c`, `(a + b) * c`},
		{`a * (b + c)`, `a * (b + c)`},
		{`-(a + b) * c`, `-(a + b) * c`},
		{`a == b`, `a == b`},
		{`a matches b`, `a matches b`},
		{`a in b`, `a in b`},
		{`a not in b`, `not (a in b)`},
		{`a and b`, `a and b`},
		{`a or b`, `a or b`},
		{`a or b and c`, `a or (b and c)`},
		{`a or (b and c)`, `a or (b and c)`},
		{`(a or b) and c`, `(a or b) and c`},
		{`a ? b : c`, `a ? b : c`},
		{`a ? b : c ? d : e`, `a ? b : (c ? d : e)`},
		{`(a ? b : c) ? d : e`, `(a ? b : c) ? d : e`},
		{`a ? (b ? c : d) : e`, `a ? (b ? c : d) : e`},
		{`func()`, `func()`},
		{`func(a)`, `func(a)`},
		{`func(a, b)`, `func(a, b)`},
		{`{}`, `{}`},
		{`{a: b}`, `{a: b}`},
		{`{a: b, c: d}`, `{a: b, c: d}`},
		{`{"a": b, 'c': d}`, `{a: b, c: d}`},
		{`{"a": b, c: d}`, `{a: b, c: d}`},
		{`{"a": b, 8: 8}`, `{a: b, "8": 8}`},
		{`{"9": 9, '8': 8, "foo": d}`, `{"9": 9, "8": 8, foo: d}`},
		{`[]`, `[]`},
		{`[a]`, `[a]`},
		{`[a, b]`, `[a, b]`},
		{`len(a)`, `len(a)`},
		{`map(a, # > 0)`, `map(a, # > 0)`},
		{`map(a, {# > 0})`, `map(a, # > 0)`},
		{`map(a, .b)`, `map(a, .b)`},
		{`a.b()`, `a.b()`},
		{`a.b(c)`, `a.b(c)`},
		{`a[1:-1]`, `a[1:-1]`},
		{`a[1:]`, `a[1:]`},
		{`a[1:]`, `a[1:]`},
		{`a[:]`, `a[:]`},
		{`a?.[1:-1]`, `a?.[1:-1]`},
		{`a?.[1:]`, `a?.[1:]`},
		{`a?.[:1]`, `a?.[:1]`},
		{`a?.[:]`, `a?.[:]`},
		{`(nil ?? 1) > 0`, `(nil ?? 1) > 0`},
		{`{("a" + "b"): 42}`, `{("a" + "b"): 42}`},
		{`(One == 1 ? true : false) && Two == 2`, `(One == 1 ? true : false) && Two == 2`},
		{`not (a == 1 ? b > 1 : b < 2)`, `not (a == 1 ? b > 1 : b < 2)`},
		{`(-(1+1)) ** 2`, `(-(1 + 1)) ** 2`},
		{`2 ** (-(1+1))`, `2 ** -(1 + 1)`},
		{`(2 ** 2) ** 3`, `(2 ** 2) ** 3`},
		{`(3 + 5) / (5 % 3)`, `(3 + 5) / (5 % 3)`},
		{`(-(1+1)) == 2`, `-(1 + 1) == 2`},
		{`if true { 1 } else { 2 }`, `if true { 1 } else { 2 }`},
		{`if true { 1 } else if false { 2 } else { 3 }`, `if true { 1 } else if false { 2 } else { 3 }`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tree, err := parser.Parse(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, tree.Node.String())
		})
	}
}

func TestPrint_parse(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`(a + b)[1:2]`, `(a + b)[1:2]`},
		{`(a ? b : c)[1:]`, `(a ? b : c)[1:]`},
		{`(a ?? b)[0]`, `(a ?? b)[0]`},
		{`(-a).b`, `(-a).b`},
		{`(a ? b : c).d`, `(a ? b : c).d`},
		{`(a ? b : c).d()`, `(a ? b : c).d()`},
		{`(1).a`, `(1).a`},
		{`-((1)[a])`, `-(1)[a]`},
		{`(nil)?.a`, `(nil)?.a`},
		{`(a?.b)["c"]`, `a?.b["c"]`},
		{`a?.b.c`, `a?.b.c`},
		{`(a?.["b"])?.["c"]`, `a?.b?.["c"]`},
		{`2.0 * a`, `2.0 * a`},
		{`1e21`, `1e+21`},
		{`(a < b) < c`, `(a < b) < c`},
		{`a < b < c`, `a < b && b < c`},
		{`(a > b)..c`, `(a > b)..c`},
		{`a..(b ? c : d)`, `a..(b ? c : d)`},
		{`map(a, 1..(.b))`, `map(a, 1..(.b))`},
		{`(a * not b) * c`, `a * (not b) * c`},
		{`-(not a) * b`, `-(not a) * b`},
		{`(a ** -b) * c`, `a ** -b * c`},
		{`(let x = a; x) + b`, `(let x = a; x) + b`},
		{`-(a; b)`, `-(a; b)`},
		{`f((a; b), c)`, `f((a; b), c)`},
		{`[(a; b), c]`, `[(a; b), c]`},
		{`{a: (b; c)}`, `{a: (b; c)}`},
		{`a[(b; c)]`, `a[(b; c)]`},
		{`a[(b; c):]`, `a[(b; c):]`},
		{`(a; b) ? (c; d) : (e; f)`, `(a; b) ? (c; d) : (e; f)`},
		{`if (a; b) { c; d } else { e; f }`, `if (a; b) { c; d } else { e; f }`},
		{`let x = (let y = a; y); x`, `let x = (let y = a; y); x`},
		{`(let x = a; x); b`, `(let x = a; x); b`},
		{`a; let x = b; x`, `a; let x = b; x`},
		{`(a; b); c`, `(a; b); c`},
		{`map(a, {b; #})`, `map(a, { b; # })`},
		{`map(a, {{b: #}})`, `map(a, { {b: #} })`},
		{`map(a, {{b: #}.b})`, `map(a, { {b: #}.b })`},
		{`reduce(a, #acc.b + #.c)`, `reduce(a, #acc.b + .c)`},
		{`{"not": a, "b c": d}`, `{"not": a, "b c": d}`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tree, err := parser.Parse(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, tree.Node.String())

			printed, err := parser.Parse(tree.Node.String())
			require.NoError(t, err)
			assert.Equal(t, ast.Dump(tree.Node), ast.Dump(printed.Node))
		})
	}
}

func TestPrint_MemberNode(t *testing.T) {
	node := &ast.MemberNode{
		Node: &ast.IdentifierNode{
			Value: "a",
		},
		Property: &ast.StringNode{Value: "b c"},
		Optional: true,
	}
	require.Equal(t, `a?.["b c"]`, node.String())
}

func TestPrint_ConstantNode(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{nil, `nil`},
		{true, `true`},
		{false, `false`},
		{1, `1`},
		{1.1, `1.1`},
		{"a", `"a"`},
		{[]int{1, 2, 3}, `[1,2,3]`},
		{map[string]int{"a": 1}, `{"a":1}`},
		{[]any{math.NaN()}, `[NaN]`},
		{math.Inf(1), `+Inf`},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			node := &ast.ConstantNode{
				Value: tt.input,
			}
			require.Equal(t, tt.want, node.String())
		})
	}
}
