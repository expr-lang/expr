package optimizer_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/checker"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/optimizer"
	"github.com/expr-lang/expr/parser"
)

func TestOptimize(t *testing.T) {
	env := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	tests := []struct {
		expr string
		want any
	}{
		{`1 + 2`, 3},
		{`sum([])`, 0},
		{`sum([a])`, 1},
		{`sum([a, b])`, 3},
		{`sum([a, b, c])`, 6},
		{`sum([a, b, c, 4])`, 10},
		{`sum(1..10, # * 1000)`, 55000},
		{`sum(map(1..10, # * 1000), # / 1000)`, float64(55)},
		{`all(1..3, {# > 0}) && all(1..3, {# < 4})`, true},
		{`all(1..3, {# > 2}) && all(1..3, {# < 4})`, false},
		{`all(1..3, {# > 0}) && all(1..3, {# < 2})`, false},
		{`all(1..3, {# > 2}) && all(1..3, {# < 2})`, false},
		{`all(1..3, {# > 0}) || all(1..3, {# < 4})`, true},
		{`all(1..3, {# > 0}) || all(1..3, {# != 2})`, true},
		{`all(1..3, {# != 3}) || all(1..3, {# < 4})`, true},
		{`all(1..3, {# != 3}) || all(1..3, {# != 2})`, false},
		{`none(1..3, {# == 0})`, true},
		{`none(1..3, {# == 0}) && none(1..3, {# == 4})`, true},
		{`none(1..3, {# == 0}) && none(1..3, {# == 3})`, false},
		{`none(1..3, {# == 1}) && none(1..3, {# == 4})`, false},
		{`none(1..3, {# == 1}) && none(1..3, {# == 3})`, false},
		{`none(1..3, {# == 0}) || none(1..3, {# == 4})`, true},
		{`none(1..3, {# == 0}) || none(1..3, {# == 3})`, true},
		{`none(1..3, {# == 1}) || none(1..3, {# == 4})`, true},
		{`none(1..3, {# == 1}) || none(1..3, {# == 3})`, false},
		{`any([1, 1, 0, 1], {# == 0})`, true},
		{`any(1..3, {# == 1}) && any(1..3, {# == 2})`, true},
		{`any(1..3, {# == 0}) && any(1..3, {# == 2})`, false},
		{`any(1..3, {# == 1}) && any(1..3, {# == 4})`, false},
		{`any(1..3, {# == 0}) && any(1..3, {# == 4})`, false},
		{`any(1..3, {# == 1}) || any(1..3, {# == 2})`, true},
		{`any(1..3, {# == 0}) || any(1..3, {# == 2})`, true},
		{`any(1..3, {# == 1}) || any(1..3, {# == 4})`, true},
		{`any(1..3, {# == 0}) || any(1..3, {# == 4})`, false},
		{`one([1, 1, 0, 1], {# == 0}) and not one([1, 0, 0, 1], {# == 0})`, true},
		{`one(1..3, {# == 1}) and one(1..3, {# == 2})`, true},
		{`one(1..3, {# == 1 || # == 2}) and one(1..3, {# == 2})`, false},
		{`one(1..3, {# == 1}) and one(1..3, {# == 2 || # == 3})`, false},
		{`one(1..3, {# == 1 || # == 2}) and one(1..3, {# == 2 || # == 3})`, false},
		{`one(1..3, {# == 1}) or one(1..3, {# == 2})`, true},
		{`one(1..3, {# == 1 || # == 2}) or one(1..3, {# == 2})`, true},
		{`one(1..3, {# == 1}) or one(1..3, {# == 2 || # == 3})`, true},
		{`one(1..3, {# == 1 || # == 2}) or one(1..3, {# == 2 || # == 3})`, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.Env(env))
			require.NoError(t, err)

			output, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)

			unoptimizedProgram, err := expr.Compile(tt.expr, expr.Env(env), expr.Optimize(false))
			require.NoError(t, err)

			unoptimizedOutput, err := expr.Run(unoptimizedProgram, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, unoptimizedOutput)
		})
	}
}

func TestOptimize_constant_folding(t *testing.T) {
	tree, err := parser.Parse(`[1,2,3][5*5-25]`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.MemberNode{
		Node:     &ast.ConstantNode{Value: []any{1, 2, 3}},
		Property: &ast.IntegerNode{Value: 0},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_constant_folding_with_floats(t *testing.T) {
	tree, err := parser.Parse(`1 + 2.0 * ((1.0 * 2) / 2) - 0`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.FloatNode{Value: 3.0}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
	assert.Equal(t, reflect.Float64, tree.Node.Type().Kind())
}

func TestOptimize_constant_folding_with_bools(t *testing.T) {
	tree, err := parser.Parse(`(true and false) or (true or false) or (false and false) or (true and (true == false))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BoolNode{Value: true}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_in_array(t *testing.T) {
	config := conf.New(map[string]int{"v": 0})

	tree, err := parser.Parse(`v in [1,2,3]`)
	require.NoError(t, err)

	_, err = checker.Check(tree, config)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.CompareNode{
		Left:        &ast.IdentifierNode{Value: "v"},
		Operators:   []string{"in"},
		Comparators: []ast.Node{&ast.ConstantNode{Value: map[int]struct{}{1: {}, 2: {}, 3: {}}}},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_in_range(t *testing.T) {
	tests := []struct {
		code     string
		expected ast.Node
	}{
		{`age in 18..31`, &ast.CompareNode{
			Left:        &ast.IntegerNode{Value: 18},
			Operators:   []string{"<=", "<="},
			Comparators: []ast.Node{&ast.IdentifierNode{Value: "age"}, &ast.IntegerNode{Value: 31}},
		}},
		{`age in 18..31 == true`, &ast.CompareNode{
			Left: &ast.CompareNode{
				Operators:   []string{"<=", "<="},
				Left:        &ast.IntegerNode{Value: 18},
				Comparators: []ast.Node{&ast.IdentifierNode{Value: "age"}, &ast.IntegerNode{Value: 31}},
			},
			Operators:   []string{"&&", "=="},
			Comparators: []ast.Node{&ast.BoolNode{Value: true}},
		}},
		{`age > 19 in 18..31`, &ast.CompareNode{
			Left:      &ast.IdentifierNode{Value: "age"},
			Operators: []string{">", "&&"},
			Comparators: []ast.Node{&ast.IntegerNode{Value: 19}, &ast.CompareNode{
				Operators:   []string{"<=", "<="},
				Left:        &ast.IntegerNode{Value: 18},
				Comparators: []ast.Node{&ast.IntegerNode{Value: 19}, &ast.IntegerNode{Value: 31}},
			}},
		}},
		{`age > 19 in 18..31 != true`, &ast.CompareNode{
			Left:      &ast.IdentifierNode{Value: "age"},
			Operators: []string{">", "&&", "!="},
			Comparators: []ast.Node{&ast.IntegerNode{Value: 19}, &ast.CompareNode{
				Operators:   []string{"<=", "<="},
				Left:        &ast.IntegerNode{Value: 18},
				Comparators: []ast.Node{&ast.IntegerNode{Value: 19}, &ast.IntegerNode{Value: 31}},
			}, &ast.BoolNode{Value: true}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			tree, err := parser.Parse(tt.code)
			require.NoError(t, err)

			config := conf.New(map[string]int{"age": 30})
			_, err = checker.Check(tree, config)

			err = optimizer.Optimize(&tree.Node, nil)
			require.NoError(t, err)
			assert.Equal(t, ast.Dump(tt.expected), ast.Dump(tree.Node))
		})
	}
}

func TestOptimize_in_range_with_floats(t *testing.T) {
	out, err := expr.Eval(`f in 1..3`, map[string]any{"f": 1.5})
	require.NoError(t, err)
	assert.Equal(t, false, out)
}

func TestOptimize_const_expr(t *testing.T) {
	tree, err := parser.Parse(`toUpper("hello")`)
	require.NoError(t, err)

	env := map[string]any{
		"toUpper": strings.ToUpper,
	}

	config := conf.New(env)
	config.ConstExpr("toUpper")

	err = optimizer.Optimize(&tree.Node, config)
	require.NoError(t, err)

	expected := &ast.ConstantNode{
		Value: "HELLO",
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_filter_len(t *testing.T) {
	tree, err := parser.Parse(`len(filter(users, .Name == "Bob"))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "count",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.CompareNode{
					Operators: []string{"=="},
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
				},
			},
		},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_filter_0(t *testing.T) {
	tree, err := parser.Parse(`filter(users, .Name == "Bob")[0]`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "find",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.CompareNode{
					Operators: []string{"=="},
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
				},
			},
		},
		Throws: true,
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_filter_first(t *testing.T) {
	tree, err := parser.Parse(`first(filter(users, .Name == "Bob"))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "find",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.CompareNode{
					Operators: []string{"=="},
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
				},
			},
		},
		Throws: false,
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_filter_minus_1(t *testing.T) {
	tree, err := parser.Parse(`filter(users, .Name == "Bob")[-1]`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "findLast",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.CompareNode{
					Operators: []string{"=="},
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
				},
			},
		},
		Throws: true,
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_filter_last(t *testing.T) {
	tree, err := parser.Parse(`last(filter(users, .Name == "Bob"))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "findLast",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.CompareNode{
					Operators: []string{"=="},
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
				},
			},
		},
		Throws: false,
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_filter_map(t *testing.T) {
	tree, err := parser.Parse(`map(filter(users, .Name == "Bob"), .Age)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "filter",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.CompareNode{
					Operators: []string{"=="},
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
				},
			},
		},
		Map: &ast.MemberNode{
			Node:     &ast.PointerNode{},
			Property: &ast.StringNode{Value: "Age"},
		},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_filter_map_first(t *testing.T) {
	tree, err := parser.Parse(`first(map(filter(users, .Name == "Bob"), .Age))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "find",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.CompareNode{
					Operators: []string{"=="},
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
				},
			},
		},
		Map: &ast.MemberNode{
			Node:     &ast.PointerNode{},
			Property: &ast.StringNode{Value: "Age"},
		},
		Throws: false,
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_predicate_combination(t *testing.T) {
	tests := []struct {
		op     string
		fn     string
		wantOp string
	}{
		{"and", "all", "and"},
		{"&&", "all", "&&"},
		{"or", "any", "or"},
		{"||", "any", "||"},
		{"and", "none", "or"},
		{"&&", "none", "||"},
	}

	for _, tt := range tests {
		rule := fmt.Sprintf(`%s(users, .Age > 18 and .Name != "Bob") %s %s(users, .Age < 30)`, tt.fn, tt.op, tt.fn)
		t.Run(rule, func(t *testing.T) {
			tree, err := parser.Parse(rule)
			require.NoError(t, err)

			err = optimizer.Optimize(&tree.Node, nil)
			require.NoError(t, err)

			expected := &ast.BuiltinNode{
				Name: tt.fn,
				Arguments: []ast.Node{
					&ast.IdentifierNode{Value: "users"},
					&ast.ClosureNode{
						Node: &ast.BinaryNode{
							Operator: tt.wantOp,
							Left: &ast.BinaryNode{
								Operator: "and",
								Left: &ast.CompareNode{
									Operators: []string{">"},
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Age"},
									},
									Comparators: []ast.Node{&ast.IntegerNode{Value: 18}},
								},
								Right: &ast.CompareNode{
									Operators: []string{"!="},
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Name"},
									},
									Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
								},
							},
							Right: &ast.CompareNode{
								Operators: []string{"<"},
								Left: &ast.MemberNode{
									Node:     &ast.PointerNode{},
									Property: &ast.StringNode{Value: "Age"},
								},
								Comparators: []ast.Node{&ast.IntegerNode{Value: 30}},
							},
						},
					},
				},
			}
			assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
		})
	}
}

func TestOptimize_predicate_combination_nested(t *testing.T) {
	tree, err := parser.Parse(`all(users, {all(.Friends, {.Age == 18 })}) && all(users, {all(.Friends, {.Name != "Bob" })})`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BuiltinNode{
		Name: "all",
		Arguments: []ast.Node{
			&ast.IdentifierNode{Value: "users"},
			&ast.ClosureNode{
				Node: &ast.BuiltinNode{
					Name: "all",
					Arguments: []ast.Node{
						&ast.MemberNode{
							Node:     &ast.PointerNode{},
							Property: &ast.StringNode{Value: "Friends"},
						},
						&ast.ClosureNode{
							Node: &ast.BinaryNode{
								Operator: "&&",
								Left: &ast.CompareNode{
									Operators: []string{"=="},
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Age"},
									},
									Comparators: []ast.Node{&ast.IntegerNode{Value: 18}},
								},
								Right: &ast.CompareNode{
									Operators: []string{"!="},
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Name"},
									},
									Comparators: []ast.Node{&ast.StringNode{Value: "Bob"}},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}
