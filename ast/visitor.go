package ast

import "fmt"

type Visitor interface {
	Enter(node *Node)
	Exit(node *Node)
}

type walker struct {
	visitor Visitor
}

func Walk(node *Node, visitor Visitor) {
	w := walker{
		visitor: visitor,
	}
	w.walk(node)
}

func (w *walker) walk(node *Node) {
	w.visitor.Enter(node)

	switch n := (*node).(type) {
	case *NilNode:
	case *IdentifierNode:
	case *IntegerNode:
	case *FloatNode:
	case *BoolNode:
	case *StringNode:
	case *ConstantNode:
	case *UnaryNode:
		w.walk(&n.Node)
	case *BinaryNode:
		w.walk(&n.Left)
		w.walk(&n.Right)
	case *MatchesNode:
		w.walk(&n.Left)
		w.walk(&n.Right)
	case *ChainNode:
		w.walk(&n.Node)
	case *MemberNode:
		w.walk(&n.Node)
		w.walk(&n.Property)
	case *SliceNode:
		w.walk(&n.Node)
		if n.From != nil {
			w.walk(&n.From)
		}
		if n.To != nil {
			w.walk(&n.To)
		}
	case *CallNode:
		w.walk(&n.Callee)
		for i := range n.Arguments {
			w.walk(&n.Arguments[i])
		}
	case *BuiltinNode:
		for i := range n.Arguments {
			w.walk(&n.Arguments[i])
		}
	case *ClosureNode:
		w.walk(&n.Node)
	case *PointerNode:
	case *ConditionalNode:
		w.walk(&n.Cond)
		w.walk(&n.Exp1)
		w.walk(&n.Exp2)
	case *ArrayNode:
		for i := range n.Nodes {
			w.walk(&n.Nodes[i])
		}
	case *MapNode:
		for i := range n.Pairs {
			w.walk(&n.Pairs[i])
		}
	case *PairNode:
		w.walk(&n.Key)
		w.walk(&n.Value)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}

	w.visitor.Exit(node)
}
