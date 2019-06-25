package optimizer

import (
	. "github.com/antonmedv/expr/ast"
)

type fold struct{}
type inRange struct{}

func (*fold) Enter(node *Node) {}
func (*fold) Exit(node *Node) {
	switch n := (*node).(type) {
	case *ArrayNode:
		if len(n.Nodes) > 0 {

			for _, a := range n.Nodes {
				if _, ok := a.(*IntegerNode); !ok {
					goto string
				}
			}
			{
				value := make([]int, len(n.Nodes))
				for i, a := range n.Nodes {
					value[i] = a.(*IntegerNode).Value
				}
				*node = &ConstantNode{
					Value: value,
				}
			}

		string:
			for _, a := range n.Nodes {
				if _, ok := a.(*StringNode); !ok {
					return
				}
			}
			{
				value := make([]string, len(n.Nodes))
				for i, a := range n.Nodes {
					value[i] = a.(*StringNode).Value
				}
				*node = &ConstantNode{
					Value: value,
				}
			}

		}
	}
}

func (*inRange) Enter(node *Node) {}
func (*inRange) Exit(node *Node) {
	switch n := (*node).(type) {
	case *BinaryNode:
		if n.Operator == "in" || n.Operator == "not in" {
			if rng, ok := n.Right.(*BinaryNode); ok && rng.Operator == ".." {
				if from, ok := n.Left.(*IntegerNode); ok {
					if to, ok := n.Right.(*IntegerNode); ok {
						*node = &BinaryNode{
							Operator: "and",
							Left: &BinaryNode{
								Operator: ">=",
								Left:     n.Left,
								Right:    from,
							},
							Right: &BinaryNode{
								Operator: "<=",
								Left:     n.Left,
								Right:    to,
							},
						}
						if n.Operator == "not in" {
							*node = &UnaryNode{
								Operator: "not",
								Node:     *node,
							}
						}
					}
				}
			}
		}
	}
}

func Optimize(node *Node) {
	Walk(node, &fold{})
	Walk(node, &inRange{})
}
