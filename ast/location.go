package ast

import (
	"github.com/antonmedv/expr/internal/helper"
)

func (n *NilNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *NilNode) GetLocation() helper.Location {
	return n.l
}

func (n *IdentifierNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *IdentifierNode) GetLocation() helper.Location {
	return n.l
}

func (n *IntegerNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *IntegerNode) GetLocation() helper.Location {
	return n.l
}

func (n *FloatNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *FloatNode) GetLocation() helper.Location {
	return n.l
}

func (n *BoolNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *BoolNode) GetLocation() helper.Location {
	return n.l
}

func (n *StringNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *StringNode) GetLocation() helper.Location {
	return n.l
}

func (n *UnaryNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *UnaryNode) GetLocation() helper.Location {
	return n.l
}

func (n *BinaryNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *BinaryNode) GetLocation() helper.Location {
	return n.l
}

func (n *MatchesNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *MatchesNode) GetLocation() helper.Location {
	return n.l
}

func (n *PropertyNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *PropertyNode) GetLocation() helper.Location {
	return n.l
}

func (n *IndexNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *IndexNode) GetLocation() helper.Location {
	return n.l
}

func (n *MethodNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *MethodNode) GetLocation() helper.Location {
	return n.l
}

func (n *FunctionNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *FunctionNode) GetLocation() helper.Location {
	return n.l
}

func (n *BuiltinNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *BuiltinNode) GetLocation() helper.Location {
	return n.l
}

func (n *ClosureNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *ClosureNode) GetLocation() helper.Location {
	return n.l
}

func (n *PointerNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *PointerNode) GetLocation() helper.Location {
	return n.l
}

func (n *ConditionalNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *ConditionalNode) GetLocation() helper.Location {
	return n.l
}

func (n *ArrayNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *ArrayNode) GetLocation() helper.Location {
	return n.l
}

func (n *MapNode) SetLocation(l helper.Location) {
	n.l = l
}

func (n *MapNode) GetLocation() helper.Location {
	return n.l
}

func (n *PairNode) GetLocation() helper.Location {
	return n.l
}

func (n *PairNode) SetLocation(l helper.Location) {
	n.l = l
}
