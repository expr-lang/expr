package optimizer

import (
	. "github.com/antonmedv/expr/ast"
	"math"
)

type fold struct {
	applied bool
}
type inRange struct{}
type constRange struct{}

func patch(node *Node, newNode Node) {
	newNode.SetType((*node).GetType())
	newNode.SetLocation((*node).GetLocation())
	*node = newNode
}

func (*fold) Enter(node *Node) {}
func (fold *fold) Exit(node *Node) {
	patch := func(newNode Node) {
		fold.applied = true
		patch(node, newNode)
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

func (*inRange) Enter(node *Node) {}
func (*inRange) Exit(node *Node) {
	switch n := (*node).(type) {
	case *BinaryNode:
		if n.Operator == "in" || n.Operator == "not in" {
			if rng, ok := n.Right.(*BinaryNode); ok && rng.Operator == ".." {
				if from, ok := rng.Left.(*IntegerNode); ok {
					if to, ok := rng.Right.(*IntegerNode); ok {
						patch(node, &BinaryNode{
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
						})
						if n.Operator == "not in" {
							patch(node, &UnaryNode{
								Operator: "not",
								Node:     *node,
							})
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
					patch(node, &ConstantNode{
						Value: value,
					})
				}
			}
		}
	}
}

func Optimize(node *Node) {
	limit := 1000
	for {
		fold := &fold{}
		Walk(node, fold)
		limit--
		if !fold.applied || limit == 0 {
			break
		}
	}
	Walk(node, &inRange{})
	Walk(node, &constRange{})
}
