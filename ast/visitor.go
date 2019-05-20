package ast

import "fmt"

type Visitor interface {
	NilNode(node *NilNode)
	IdentifierNode(node *IdentifierNode)
	IntegerNode(node *IntegerNode)
	FloatNode(node *FloatNode)
	BoolNode(node *BoolNode)
	StringNode(node *StringNode)
	UnaryNode(node *UnaryNode)
	BinaryNode(node *BinaryNode)
	MatchesNode(node *MatchesNode)
	PropertyNode(node *PropertyNode)
	IndexNode(node *IndexNode)
	MethodNode(node *MethodNode)
	FunctionNode(node *FunctionNode)
	BuiltinNode(node *BuiltinNode)
	ClosureNode(node *ClosureNode)
	PointerNode(node *PointerNode)
	ConditionalNode(node *ConditionalNode)
	ArrayNode(node *ArrayNode)
	MapNode(node *MapNode)
	PairNode(node *PairNode)
}

type walker struct {
	visitor Visitor
}

func Walk(node Node, visitor Visitor) {
	w := walker{
		visitor: visitor,
	}
	w.walk(node)
}

func (w *walker) walk(node Node) {

	switch n := node.(type) {
	case *NilNode:
		w.visitor.NilNode(n)
	case *IdentifierNode:
		w.visitor.IdentifierNode(n)
	case *IntegerNode:
		w.visitor.IntegerNode(n)
	case *FloatNode:
		w.visitor.FloatNode(n)
	case *BoolNode:
		w.visitor.BoolNode(n)
	case *StringNode:
		w.visitor.StringNode(n)
	case *UnaryNode:
		w.walk(n.Node)
		w.visitor.UnaryNode(n)
	case *BinaryNode:
		w.walk(n.Left)
		w.walk(n.Right)
		w.visitor.BinaryNode(n)
	case *MatchesNode:
		w.walk(n.Left)
		w.walk(n.Right)
		w.visitor.MatchesNode(n)
	case *PropertyNode:
		w.walk(n.Node)
		w.visitor.PropertyNode(n)
	case *IndexNode:
		w.walk(n.Node)
		w.walk(n.Index)
		w.visitor.IndexNode(n)
	case *MethodNode:
		w.walk(n.Node)
		for _, arg := range n.Arguments {
			w.walk(arg)
		}
		w.visitor.MethodNode(n)
	case *FunctionNode:
		for _, arg := range n.Arguments {
			w.walk(arg)
		}
		w.visitor.FunctionNode(n)
	case *BuiltinNode:
		for _, arg := range n.Arguments {
			w.walk(arg)
		}
		w.visitor.BuiltinNode(n)
	case *ClosureNode:
		w.walk(n.Node)
		w.visitor.ClosureNode(n)
	case *PointerNode:
		w.visitor.PointerNode(n)
	case *ConditionalNode:
		w.walk(n.Cond)
		w.walk(n.Exp1)
		w.walk(n.Exp2)
		w.visitor.ConditionalNode(n)
	case *ArrayNode:
		for _, node := range n.Nodes {
			w.walk(node)
		}
		w.visitor.ArrayNode(n)
	case *MapNode:
		for _, pair := range n.Pairs {
			w.walk(pair)
		}
		w.visitor.MapNode(n)
	case *PairNode:
		w.walk(n.Value)
		w.visitor.PairNode(n)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
}
