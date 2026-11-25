package optimizer_test

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

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

func TestOptimize_in_array(t *testing.T) {
	config := conf.New(map[string]int{"v": 0})

	tree, err := parser.Parse(`v in [1,2,3]`)
	require.NoError(t, err)

	_, err = checker.Check(tree, config)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.BinaryNode{
		Operator: "in",
		Left:     &ast.IdentifierNode{Value: "v"},
		Right:    &ast.ConstantNode{Value: map[int]struct{}{1: {}, 2: {}, 3: {}}},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}

func TestOptimize_in_range(t *testing.T) {
	tree, err := parser.Parse(`age in 18..31`)
	require.NoError(t, err)

	config := conf.New(map[string]int{"age": 30})
	_, err = checker.Check(tree, config)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	left := &ast.IdentifierNode{
		Value: "age",
	}
	expected := &ast.BinaryNode{
		Operator: "and",
		Left: &ast.BinaryNode{
			Operator: ">=",
			Left:     left,
			Right: &ast.IntegerNode{
				Value: 18,
			},
		},
		Right: &ast.BinaryNode{
			Operator: "<=",
			Left:     left,
			Right: &ast.IntegerNode{
				Value: 31,
			},
		},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
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
			&ast.PredicateNode{
				Node: &ast.BinaryNode{
					Operator: "==",
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Right: &ast.StringNode{Value: "Bob"},
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
			&ast.PredicateNode{
				Node: &ast.BinaryNode{
					Operator: "==",
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Right: &ast.StringNode{Value: "Bob"},
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
			&ast.PredicateNode{
				Node: &ast.BinaryNode{
					Operator: "==",
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Right: &ast.StringNode{Value: "Bob"},
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
			&ast.PredicateNode{
				Node: &ast.BinaryNode{
					Operator: "==",
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Right: &ast.StringNode{Value: "Bob"},
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
			&ast.PredicateNode{
				Node: &ast.BinaryNode{
					Operator: "==",
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Right: &ast.StringNode{Value: "Bob"},
				},
			},
		},
		Throws: false,
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
			&ast.PredicateNode{
				Node: &ast.BinaryNode{
					Operator: "==",
					Left: &ast.MemberNode{
						Node:     &ast.PointerNode{},
						Property: &ast.StringNode{Value: "Name"},
					},
					Right: &ast.StringNode{Value: "Bob"},
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
					&ast.PredicateNode{
						Node: &ast.BinaryNode{
							Operator: tt.wantOp,
							Left: &ast.BinaryNode{
								Operator: "and",
								Left: &ast.BinaryNode{
									Operator: ">",
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Age"},
									},
									Right: &ast.IntegerNode{Value: 18},
								},
								Right: &ast.BinaryNode{
									Operator: "!=",
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Name"},
									},
									Right: &ast.StringNode{Value: "Bob"},
								},
							},
							Right: &ast.BinaryNode{
								Operator: "<",
								Left: &ast.MemberNode{
									Node:     &ast.PointerNode{},
									Property: &ast.StringNode{Value: "Age"},
								},
								Right: &ast.IntegerNode{Value: 30},
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
			&ast.PredicateNode{
				Node: &ast.BuiltinNode{
					Name: "all",
					Arguments: []ast.Node{
						&ast.MemberNode{
							Node:     &ast.PointerNode{},
							Property: &ast.StringNode{Value: "Friends"},
						},
						&ast.PredicateNode{
							Node: &ast.BinaryNode{
								Operator: "&&",
								Left: &ast.BinaryNode{
									Operator: "==",
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Age"},
									},
									Right: &ast.IntegerNode{Value: 18},
								},
								Right: &ast.BinaryNode{
									Operator: "!=",
									Left: &ast.MemberNode{
										Node:     &ast.PointerNode{},
										Property: &ast.StringNode{Value: "Name"},
									},
									Right: &ast.StringNode{Value: "Bob"},
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

func TestOptimize_predicate_combination_different_collections(t *testing.T) {
	env := map[string]any{
		"cond": true,
		"a":    []int{-1, 5},
		"b":    []int{7, 8},
	}
	for code, want := range map[string]bool{
		`all((cond ? a : b)[1:], # > 0) && all(cond ? a : b[1:], # > 0)`:     false,
		`any(cond ? a : b[1:], # == -5) || any((cond ? b : a)[1:], # == 8)`:  true,
		`any((cond ? a : b)[1:], # == -1) || any(cond ? a : b[1:], # == -1)`: true,
		`none((cond ? a : b)[1:], # < 0) && none(cond ? a : b[1:], # < 0)`:   false,
	} {
		t.Run(code, func(t *testing.T) {
			for _, optimize := range []bool{false, true} {
				program, err := expr.Compile(code, expr.Env(env), expr.Optimize(optimize))
				require.NoError(t, err)

				got, err := expr.Run(program, env)
				require.NoError(t, err)
				require.Equal(t, want, got)
			}
		})
	}
}

func TestOptimize_predicate_combination_with_non_json_float(t *testing.T) {
	// Constant folding turns [0/0] into a ConstantNode holding []any{NaN}.
	// predicateCombination compares the collection arguments via String(),
	// which must not panic on values json.Marshal rejects (NaN, ±Inf).
	for _, code := range []string{
		`any([0/0], # > 0) || any([0/0], # > 0)`,
		`all([1/0], # > 0) && all([1/0], # > 0)`,
		`none([0/0], # > 0) && none([0/0], # > 0)`,
	} {
		_, err := expr.Compile(code)
		require.NoError(t, err, code)
	}
}

func TestOptimize_predicate_combination_same_string(t *testing.T) {
	// predicateCombination must merge only calls on the same collection.
	// These collections differ but print the same via String().
	env := map[string]any{"a": []int{1, 2, 3}, "n": 2, "NaN": 7}
	for _, code := range []string{
		// filterMap sets BuiltinNode.Map, which String() does not print.
		`all(map(filter(a, # > 0), # * 10), # > 5) && all(filter(a, # > 0), # < 5)`,
		`not (any(map(filter(a, # > 0), # * 10), # == 1) || any(filter(a, # > 0), # == 30))`,
		// Folded (-2) prints as -2, so (-2) ** n printed as -2 ** n.
		`all([(-2) ** n], # > 0) && all([-2 ** n], # < 0)`,
		`all([(-2.0) ** n], # > 0) && all([-2.0 ** n], # < 0)`,
		// Folded constants print 2.0 as 2, and NaN like an identifier.
		`all([2.0], type(#) == "float") && all([2], type(#) == "int")`,
		`all([0/0], # != #) && all([NaN], # == 7)`,
		// 0.0 and -0.0 are equal, but they are not the same constant.
		`all([-0.0], 1 / # < 0) && all([0.0], 1 / # > 0)`,
	} {
		t.Run(code, func(t *testing.T) {
			unoptimized, err := expr.Compile(code, expr.Env(env), expr.Optimize(false))
			require.NoError(t, err)
			want, err := expr.Run(unoptimized, env)
			require.NoError(t, err)
			require.Equal(t, true, want)

			optimized, err := expr.Compile(code, expr.Env(env))
			require.NoError(t, err)
			got, err := expr.Run(optimized, env)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}

func TestOptimize_predicate_combination_same_constant(t *testing.T) {
	// Collections holding folded constants still merge when they are the same.
	env := map[string]any{"a": []int{1, 2, 3}}
	for _, code := range []string{
		`all([1, 2, 3], # > 0) && all([1, 2, 3], # < 5)`,
		`all(filter(a, # in [1, 2]), # > 0) && all(filter(a, # in [1, 2]), # < 5)`,
		`all([0/0], # != #) && all([0/0], # != 0)`,
	} {
		t.Run(code, func(t *testing.T) {
			tree, err := parser.Parse(code)
			require.NoError(t, err)
			config := conf.New(env)
			_, err = checker.Check(tree, config)
			require.NoError(t, err)
			err = optimizer.Optimize(&tree.Node, config)
			require.NoError(t, err)

			builtin, ok := tree.Node.(*ast.BuiltinNode)
			require.True(t, ok, "not merged: %s", tree.Node)
			assert.Equal(t, "all", builtin.Name)
		})
	}
}

func TestOptimize_predicate_combination_constant_values(t *testing.T) {
	// Constants holding values that the comparison cannot check exactly must
	// not be merged when they differ. The expression is true unless c2's
	// predicate wrongly runs on c1's value.
	code := `all([c1()], same(#, c1())) && all([c2()], same(#, c2()))`
	x0, x5, y0, y5 := 0, 5, 0, 5
	for _, tt := range []struct {
		name   string
		c1, c2 any
	}{
		{"pointer keys", []any{map[*int]int{&x0: 1, &x5: 2}, &x0}, []any{map[*int]int{&y0: 1, &y5: 2}, &y5}},
		{"NaN keys", map[float64]int{math.NaN(): 1, math.NaN(): 1}, map[float64]int{math.NaN(): 1, math.NaN(): 2}},
		{"signed zero keys", map[float64]int{math.Copysign(0, -1): 1}, map[float64]int{0: 1}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			env := map[string]any{
				"c1": func() any { return tt.c1 },
				"c2": func() any { return tt.c2 },
				// c1 and c2 return the same map or slice each time, so
				// compare identity.
				"same": func(a, b any) bool {
					return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
				},
			}
			program, err := expr.Compile(code, expr.Env(env), expr.ConstExpr("c1"), expr.ConstExpr("c2"))
			require.NoError(t, err)
			out, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, true, out)
		})
	}
}

func TestOptimize_predicate_combination_same_pointer(t *testing.T) {
	// A constant holding pointers merges when both sides hold the same data.
	type user struct{ Age int }
	users := []*user{{Age: 20}, {Age: 30}}
	env := map[string]any{"users": func() any { return users }}
	tree, err := parser.Parse(`all(users(), .Age > 18) && all(users(), .Age < 40)`)
	require.NoError(t, err)
	config := conf.New(env)
	config.ConstExpr("users")
	_, err = checker.Check(tree, config)
	require.NoError(t, err)
	err = optimizer.Optimize(&tree.Node, config)
	require.NoError(t, err)

	builtin, ok := tree.Node.(*ast.BuiltinNode)
	require.True(t, ok, "not merged: %T", tree.Node)
	assert.Equal(t, "all", builtin.Name)
}

func TestOptimize_predicate_combination_large_constant(t *testing.T) {
	// Constants that refer to themselves, or reuse one value in many places,
	// must not hang the optimizer. Each case builds two equal values that do
	// not share data, so the comparison has to walk them.
	cycle := func() any {
		m := map[string]any{}
		m["a"] = m
		m["b"] = m
		return m
	}
	slices := func() any {
		var v any = 1
		for i := 0; i < 40; i++ {
			v = []any{v, v}
		}
		return v
	}
	arrays := func() any {
		var v any = 1
		for i := 0; i < 40; i++ {
			v = [2]any{v, v}
		}
		return v
	}
	for name, build := range map[string]func() any{"cycle": cycle, "shared slices": slices, "shared arrays": arrays} {
		c1, c2 := build(), build()
		t.Run(name, func(t *testing.T) {
			env := map[string]any{
				"c1": func() any { return c1 },
				"c2": func() any { return c2 },
			}
			tree, err := parser.Parse(`all([c1()], # != nil) && all([c2()], # != 1)`)
			require.NoError(t, err)
			config := conf.New(env)
			config.ConstExpr("c1")
			config.ConstExpr("c2")
			_, err = checker.Check(tree, config)
			require.NoError(t, err)

			done := make(chan error, 1)
			go func() { done <- optimizer.Optimize(&tree.Node, config) }()
			select {
			case err := <-done:
				require.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("optimizer did not finish")
			}
		})
	}
}
