package optimizer

import (
	"math"

	. "github.com/antonmedv/expr/ast"
)

type fold struct {
	applied bool
}

func (*fold) Enter(*Node) {}
func (fold *fold) Exit(node *Node) {
	patch := func(newNode Node) {
		fold.applied = true
		Patch(node, newNode)
	}

	switch n := (*node).(type) {
	case *UnaryNode:
		switch n.Operator {
		case "-":
			if i, ok := n.Node.(*IntegerNode); ok {
				patch(&IntegerNode{Value: -i.Value})
			}
		case "+":
			if i, ok := n.Node.(*IntegerNode); ok {
				patch(&IntegerNode{Value: i.Value})
			}
		}

	case *BinaryNode:
		switch n.Operator {
		case "+":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					patch(&IntegerNode{Value: a.Value + b.Value})
				}
			}
			if a, ok := n.Left.(*StringNode); ok {
				if b, ok := n.Right.(*StringNode); ok {
					patch(&StringNode{Value: a.Value + b.Value})
				}
			}
		case "-":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					patch(&IntegerNode{Value: a.Value - b.Value})
				}
			}
		case "*":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					patch(&IntegerNode{Value: a.Value * b.Value})
				}
			}
		case "/":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					patch(&IntegerNode{Value: a.Value / b.Value})
				}
			}
		case "%":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					patch(&IntegerNode{Value: a.Value % b.Value})
				}
			}
		case "**":
			if a, ok := n.Left.(*IntegerNode); ok {
				if b, ok := n.Right.(*IntegerNode); ok {
					patch(&FloatNode{Value: math.Pow(float64(a.Value), float64(b.Value))})
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
				patch(&ConstantNode{Value: value})
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
				patch(&ConstantNode{Value: value})
			}

		}
	}
}
