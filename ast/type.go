package ast

import "reflect"

func (n *NilNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *NilNode) GetType() reflect.Type {
	return n.t
}

func (n *IdentifierNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *IdentifierNode) GetType() reflect.Type {
	return n.t
}

func (n *IntegerNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *IntegerNode) GetType() reflect.Type {
	return n.t
}

func (n *FloatNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *FloatNode) GetType() reflect.Type {
	return n.t
}

func (n *BoolNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *BoolNode) GetType() reflect.Type {
	return n.t
}

func (n *StringNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *StringNode) GetType() reflect.Type {
	return n.t
}

func (n *ConstantNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *ConstantNode) GetType() reflect.Type {
	return n.t
}

func (n *UnaryNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *UnaryNode) GetType() reflect.Type {
	return n.t
}

func (n *BinaryNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *BinaryNode) GetType() reflect.Type {
	return n.t
}

func (n *MatchesNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *MatchesNode) GetType() reflect.Type {
	return n.t
}

func (n *PropertyNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *PropertyNode) GetType() reflect.Type {
	return n.t
}

func (n *IndexNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *IndexNode) GetType() reflect.Type {
	return n.t
}

func (n *SliceNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *SliceNode) GetType() reflect.Type {
	return n.t
}

func (n *MethodNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *MethodNode) GetType() reflect.Type {
	return n.t
}

func (n *FunctionNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *FunctionNode) GetType() reflect.Type {
	return n.t
}

func (n *BuiltinNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *BuiltinNode) GetType() reflect.Type {
	return n.t
}

func (n *ClosureNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *ClosureNode) GetType() reflect.Type {
	return n.t
}

func (n *PointerNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *PointerNode) GetType() reflect.Type {
	return n.t
}

func (n *ConditionalNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *ConditionalNode) GetType() reflect.Type {
	return n.t
}

func (n *ArrayNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *ArrayNode) GetType() reflect.Type {
	return n.t
}

func (n *MapNode) SetType(t reflect.Type) {
	n.t = t
}

func (n *MapNode) GetType() reflect.Type {
	return n.t
}

func (n *PairNode) GetType() reflect.Type {
	return n.t
}

func (n *PairNode) SetType(t reflect.Type) {
	n.t = t
}
