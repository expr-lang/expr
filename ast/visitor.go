package ast

import "fmt"

type Visitor interface {
	Visit(node *Node)
}

func Walk(node *Node, v Visitor) {
	if *node == nil {
		return
	}
	switch n := (*node).(type) {
	case *NilNode:
	case *IdentifierNode:
	case *IntegerNode:
	case *FloatNode:
	case *BoolNode:
	case *StringNode:
	case *ConstantNode:
	case *ChainNode:
		Walk(&n.Node, v)
	case *MemberNode:
		Walk(&n.Node, v)
		Walk(&n.Property, v)
	case *SliceNode:
		Walk(&n.Node, v)
		if n.From != nil {
			Walk(&n.From, v)
		}
		if n.To != nil {
			Walk(&n.To, v)
		}
	case *CallNode:
		Walk(&n.Callee, v)
		for i := range n.Arguments {
			Walk(&n.Arguments[i], v)
		}
	case *BuiltinNode:
		for i := range n.Arguments {
			Walk(&n.Arguments[i], v)
		}
	case *ArrayNode:
		for i := range n.Nodes {
			Walk(&n.Nodes[i], v)
		}
	case *MapNode:
		for i := range n.Pairs {
			Walk(&n.Pairs[i], v)
		}
	case *PairNode:
		Walk(&n.Key, v)
		Walk(&n.Value, v)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}

	v.Visit(node)
}
