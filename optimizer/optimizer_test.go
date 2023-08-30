package optimizer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestOptimize_const_range(t *testing.T) {
	tree, err := parser.Parse(`-1..1`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &ast.ConstantNode{
		Value: []int{-1, 0, 1},
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
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
			&ast.ClosureNode{
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
			&ast.ClosureNode{
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
			&ast.ClosureNode{
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
			&ast.ClosureNode{
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
