package ast

import (
	"reflect"
	"regexp"

	"github.com/antonmedv/expr/file"
)

// Node represents items of abstract syntax tree.
type Node interface {
	Location() file.Location
	SetLocation(file.Location)
	Type() reflect.Type
	SetType(reflect.Type)
}

type Base struct {
	loc      file.Location
	nodeType reflect.Type
}

func (n *Base) Location() file.Location {
	return n.loc
}

func (n *Base) SetLocation(loc file.Location) {
	n.loc = loc
}

func (n *Base) Type() reflect.Type {
	return n.nodeType
}

func (n *Base) SetType(t reflect.Type) {
	n.nodeType = t
}

func Loc(l file.Location) Base {
	return Base{loc: l}
}

type NilNode struct {
	Base
}

type IdentifierNode struct {
	Base
	Value string
}

type IntegerNode struct {
	Base
	Value int
}

type FloatNode struct {
	Base
	Value float64
}

type BoolNode struct {
	Base
	Value bool
}

type StringNode struct {
	Base
	Value string
}

type ConstantNode struct {
	Base
	Value interface{}
}

type UnaryNode struct {
	Base
	Operator string
	Node     Node
}

type BinaryNode struct {
	Base
	Operator string
	Left     Node
	Right    Node
}

type MatchesNode struct {
	Base
	Regexp *regexp.Regexp
	Left   Node
	Right  Node
}

type PropertyNode struct {
	Base
	Node     Node
	Property string
}

type IndexNode struct {
	Base
	Node  Node
	Index Node
}

type SliceNode struct {
	Base
	Node Node
	From Node
	To   Node
}

type MethodNode struct {
	Base
	Node      Node
	Method    string
	Arguments []Node
}

type FunctionNode struct {
	Base
	Name      string
	Arguments []Node
	Fast      bool
}

type BuiltinNode struct {
	Base
	Name      string
	Arguments []Node
}

type ClosureNode struct {
	Base
	Node Node
}

type PointerNode struct {
	Base
}

type ConditionalNode struct {
	Base
	Cond Node
	Exp1 Node
	Exp2 Node
}

type ArrayNode struct {
	Base
	Nodes []Node
}

type MapNode struct {
	Base
	Pairs []Node
}

type PairNode struct {
	Base
	Key   Node
	Value Node
}
