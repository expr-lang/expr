package ast_test

import (
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
