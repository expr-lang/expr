package optimizer

import (
	. "github.com/antonmedv/expr/ast"
	"math"
)

type fold struct{}
type inRange struct{}
type constRange struct{}

func (*fold) Enter(node *Node) {}
func (*fold) Exit(node *Node) {
	switch n := (*node).(type) {
	case *UnaryNode:
		if n.Operator == "-" {
			if i, ok := n.Node.(*IntegerNode); ok {
				*node = &IntegerNode{
					Value: -i.Value,
				}
			}
		}

	case *BinaryNode:
		switch n.Operator {
		case "+":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					*node = &IntegerNode{
						Value: a.Value + b.Value,
					}
				}
			}
			if a, ok := n.Left.(*StringNode); ok {
				if b, ok := n.Right.(*StringNode); ok {
					*node = &StringNode{
						Value: a.Value + b.Value,
					}
				}
			}
		case "**":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					*node = &FloatNode{
						Value: math.Pow(float64(a.Value), float64(b.Value)),
					}
				}
			}
		}

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
				if from, ok := rng.Left.(*IntegerNode); ok {
					if to, ok := rng.Right.(*IntegerNode); ok {
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

func (*constRange) Enter(node *Node) {}
func (*constRange) Exit(node *Node) {
	switch n := (*node).(type) {
	case *BinaryNode:
		if n.Operator == ".." {
			if min, ok := n.Left.(*IntegerNode); ok {
				if max, ok := n.Right.(*IntegerNode); ok {
					size := max.Value - min.Value + 1
					value := make([]int, size)
					for i := range value {
						value[i] = min.Value + i
					}
					*node = &ConstantNode{
						Value: value,
					}
				}
			}
		}
	}
}

func Optimize(node *Node) {
	Walk(node, &fold{})
	Walk(node, &inRange{})
	Walk(node, &constRange{})
}
