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
		w.visitor.Exit(node)
	case *IdentifierNode:
		w.visitor.Exit(node)
	case *IntegerNode:
		w.visitor.Exit(node)
	case *FloatNode:
		w.visitor.Exit(node)
	case *BoolNode:
		w.visitor.Exit(node)
	case *StringNode:
		w.visitor.Exit(node)
	case *UnaryNode:
		w.walk(&n.Node)
		w.visitor.Exit(node)
	case *BinaryNode:
		w.walk(&n.Left)
		w.walk(&n.Right)
		w.visitor.Exit(node)
	case *MatchesNode:
		w.walk(&n.Left)
		w.walk(&n.Right)
		w.visitor.Exit(node)
	case *PropertyNode:
		w.walk(&n.Node)
		w.visitor.Exit(node)
	case *IndexNode:
		w.walk(&n.Node)
		w.walk(&n.Index)
		w.visitor.Exit(node)
	case *SliceNode:
		w.walk(&n.From)
		w.walk(&n.To)
		w.visitor.Exit(node)
	case *MethodNode:
		w.walk(&n.Node)
		for _, arg := range n.Arguments {
			w.walk(&arg)
		}
		w.visitor.Exit(node)
	case *FunctionNode:
		for _, arg := range n.Arguments {
			w.walk(&arg)
		}
		w.visitor.Exit(node)
	case *BuiltinNode:
		for _, arg := range n.Arguments {
			w.walk(&arg)
		}
		w.visitor.Exit(node)
	case *ClosureNode:
		w.walk(&n.Node)
		w.visitor.Exit(node)
	case *PointerNode:
		w.visitor.Exit(node)
	case *ConditionalNode:
		w.walk(&n.Cond)
		w.walk(&n.Exp1)
		w.walk(&n.Exp2)
		w.visitor.Exit(node)
	case *ArrayNode:
		for _, node := range n.Nodes {
			w.walk(&node)
		}
		w.visitor.Exit(node)
	case *MapNode:
		var pair Node
		for _, pair = range n.Pairs {
			w.walk(&pair)
		}
		w.visitor.Exit(node)
	case *PairNode:
		w.walk(&n.Value)
		w.visitor.Exit(node)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
}
