package ast

import (
	"github.com/antonmedv/expr/internal/file"
)

func (n *NilNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *NilNode) GetLocation() file.Location {
	return n.l
}

func (n *IdentifierNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *IdentifierNode) GetLocation() file.Location {
	return n.l
}

func (n *IntegerNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *IntegerNode) GetLocation() file.Location {
	return n.l
}

func (n *FloatNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *FloatNode) GetLocation() file.Location {
	return n.l
}

func (n *BoolNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *BoolNode) GetLocation() file.Location {
	return n.l
}

func (n *StringNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *StringNode) GetLocation() file.Location {
	return n.l
}

func (n *ConstantNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *ConstantNode) GetLocation() file.Location {
	return n.l
}

func (n *UnaryNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *UnaryNode) GetLocation() file.Location {
	return n.l
}

func (n *BinaryNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *BinaryNode) GetLocation() file.Location {
	return n.l
}

func (n *MatchesNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *MatchesNode) GetLocation() file.Location {
	return n.l
}

func (n *PropertyNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *PropertyNode) GetLocation() file.Location {
	return n.l
}

func (n *IndexNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *IndexNode) GetLocation() file.Location {
	return n.l
}

func (n *SliceNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *SliceNode) GetLocation() file.Location {
	return n.l
}

func (n *MethodNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *MethodNode) GetLocation() file.Location {
	return n.l
}

func (n *FunctionNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *FunctionNode) GetLocation() file.Location {
	return n.l
}

func (n *BuiltinNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *BuiltinNode) GetLocation() file.Location {
	return n.l
}

func (n *ClosureNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *ClosureNode) GetLocation() file.Location {
	return n.l
}

func (n *PointerNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *PointerNode) GetLocation() file.Location {
	return n.l
}

func (n *ConditionalNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *ConditionalNode) GetLocation() file.Location {
	return n.l
}

func (n *ArrayNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *ArrayNode) GetLocation() file.Location {
	return n.l
}

func (n *MapNode) SetLocation(l file.Location) {
	n.l = l
}

func (n *MapNode) GetLocation() file.Location {
	return n.l
}

func (n *PairNode) GetLocation() file.Location {
	return n.l
}

func (n *PairNode) SetLocation(l file.Location) {
	n.l = l
}
