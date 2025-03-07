package optimizer_test

import (
	"testing"

	. "github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/optimizer"
	"github.com/expr-lang/expr/parser"
)

func TestOptimize_filter_map(t *testing.T) {
	tree, err := parser.Parse(`map(filter(users, .Name == "Bob"), .Age)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &BuiltinNode{
		Name: "filter",
		Arguments: []Node{
			&IdentifierNode{Value: "users"},
			&PredicateNode{
				Node: &BinaryNode{
					Operator: "==",
					Left: &MemberNode{
						Node:     &PointerNode{},
						Property: &StringNode{Value: "Name"},
					},
					Right: &StringNode{Value: "Bob"},
				},
			},
		},
		Map: &MemberNode{
			Node:     &PointerNode{},
			Property: &StringNode{Value: "Age"},
		},
	}

	assert.Equal(t, Dump(expected), Dump(tree.Node))
}

func TestOptimize_filter_map_with_index_pointer(t *testing.T) {
	tree, err := parser.Parse(`map(filter(users, true), #index)`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &BuiltinNode{
		Name: "map",
		Arguments: []Node{
			&BuiltinNode{
				Name: "filter",
				Arguments: []Node{
					&IdentifierNode{Value: "users"},
					&PredicateNode{
						Node: &BoolNode{Value: true},
					},
				},
				Throws: false,
				Map:    nil,
			},
			&PredicateNode{
				Node: &PointerNode{Name: "index"},
			},
		},
		Throws: false,
		Map:    nil,
	}

	assert.Equal(t, Dump(expected), Dump(tree.Node))
}

func TestOptimize_filter_map_with_index_pointer_with_index_pointer_in_first_argument(t *testing.T) {
	tree, err := parser.Parse(`1..2 | map(map(filter([#index], true), 42))`)
	require.NoError(t, err)

	err = optimizer.Optimize(&tree.Node, nil)
	require.NoError(t, err)

	expected := &BuiltinNode{
		Name: "map",
		Arguments: []Node{
			&BinaryNode{
				Operator: "..",
				Left:     &IntegerNode{Value: 1},
				Right:    &IntegerNode{Value: 2},
			},
			&PredicateNode{
				Node: &BuiltinNode{
					Name: "filter",
					Arguments: []Node{
						&ArrayNode{
							Nodes: []Node{
								&PointerNode{Name: "index"},
							},
						},
						&PredicateNode{
							Node: &BoolNode{Value: true},
						},
					},
					Throws: false,
					Map:    &IntegerNode{Value: 42},
				},
			},
		},
		Throws: false,
		Map:    nil,
	}

	assert.Equal(t, Dump(expected), Dump(tree.Node))
}
