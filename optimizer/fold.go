package optimizer

import (
	. "expr/ast"
	"expr/file"
)

type fold struct {
	applied bool
	err     *file.Error
}

func (fold *fold) Visit(node *Node) {
	patch := func(newNode Node) {
		fold.applied = true
		patchWithType(node, newNode)
	}

	switch n := (*node).(type) {
	case *ArrayNode:
		if len(n.Nodes) > 0 {
			for _, a := range n.Nodes {
				switch a.(type) {
				case *IntegerNode, *FloatNode, *StringNode, *BoolNode:
					continue
				default:
					return
				}
			}
			value := make([]any, len(n.Nodes))
			for i, a := range n.Nodes {
				switch b := a.(type) {
				case *IntegerNode:
					value[i] = b.Value
				case *FloatNode:
					value[i] = b.Value
				case *StringNode:
					value[i] = b.Value
				case *BoolNode:
					value[i] = b.Value
				}
			}
			patch(&ConstantNode{Value: value})
		}
	}
}

func toString(n Node) *StringNode {
	switch a := n.(type) {
	case *StringNode:
		return a
	}
	return nil
}

func toInteger(n Node) *IntegerNode {
	switch a := n.(type) {
	case *IntegerNode:
		return a
	}
	return nil
}

func toFloat(n Node) *FloatNode {
	switch a := n.(type) {
	case *FloatNode:
		return a
	}
	return nil
}

func toBool(n Node) *BoolNode {
	switch a := n.(type) {
	case *BoolNode:
		return a
	}
	return nil
}
