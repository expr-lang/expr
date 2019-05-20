package ast

import (
	"github.com/antonmedv/expr/internal/helper"
)

func (n *NilNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *NilNode) GetLocation() helper.Location {
	return n.location
}

func (n *IdentifierNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *IdentifierNode) GetLocation() helper.Location {
	return n.location
}

func (n *IntegerNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *IntegerNode) GetLocation() helper.Location {
	return n.location
}

func (n *FloatNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *FloatNode) GetLocation() helper.Location {
	return n.location
}

func (n *BoolNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *BoolNode) GetLocation() helper.Location {
	return n.location
}

func (n *StringNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *StringNode) GetLocation() helper.Location {
	return n.location
}

func (n *UnaryNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *UnaryNode) GetLocation() helper.Location {
	return n.location
}

func (n *BinaryNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *BinaryNode) GetLocation() helper.Location {
	return n.location
}

func (n *MatchesNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *MatchesNode) GetLocation() helper.Location {
	return n.location
}

func (n *PropertyNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *PropertyNode) GetLocation() helper.Location {
	return n.location
}

func (n *IndexNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *IndexNode) GetLocation() helper.Location {
	return n.location
}

func (n *MethodNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *MethodNode) GetLocation() helper.Location {
	return n.location
}

func (n *FunctionNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *FunctionNode) GetLocation() helper.Location {
	return n.location
}

func (n *BuiltinNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *BuiltinNode) GetLocation() helper.Location {
	return n.location
}

func (n *ClosureNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *ClosureNode) GetLocation() helper.Location {
	return n.location
}

func (n *PointerNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *PointerNode) GetLocation() helper.Location {
	return n.location
}

func (n *ConditionalNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *ConditionalNode) GetLocation() helper.Location {
	return n.location
}

func (n *ArrayNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *ArrayNode) GetLocation() helper.Location {
	return n.location
}

func (n *MapNode) SetLocation(l helper.Location) {
	n.location = l
}

func (n *MapNode) GetLocation() helper.Location {
	return n.location
}

func (n *PairNode) GetLocation() helper.Location {
	return n.location
}

func (n *PairNode) SetLocation(l helper.Location) {
	n.location = l
}
