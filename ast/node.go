package ast

import (
	"gopkg.in/antonmedv/expr.v2/internal/file"
	"reflect"
	"regexp"
)

// Node represents items of abstract syntax tree.
type Node interface {
	GetLocation() file.Location
	SetLocation(file.Location)
	GetType() reflect.Type
	SetType(reflect.Type)
}

type NilNode struct {
	l file.Location
	t reflect.Type
}

type IdentifierNode struct {
	l file.Location
	t reflect.Type

	Value string
}

type IntegerNode struct {
	l file.Location
	t reflect.Type

	Value int
}

type FloatNode struct {
	l file.Location
	t reflect.Type

	Value float64
}

type BoolNode struct {
	l file.Location
	t reflect.Type

	Value bool
}

type StringNode struct {
	l file.Location
	t reflect.Type

	Value string
}

type UnaryNode struct {
	l file.Location
	t reflect.Type

	Operator string
	Node     Node
}

type BinaryNode struct {
	l file.Location
	t reflect.Type

	Operator string
	Left     Node
	Right    Node
}

type MatchesNode struct {
	l file.Location
	t reflect.Type

	Regexp *regexp.Regexp
	Left   Node
	Right  Node
}

type PropertyNode struct {
	l file.Location
	t reflect.Type

	Node     Node
	Property string
}

type IndexNode struct {
	l file.Location
	t reflect.Type

	Node  Node
	Index Node
}

type MethodNode struct {
	l file.Location
	t reflect.Type

	Node      Node
	Method    string
	Arguments []Node
}

type FunctionNode struct {
	l file.Location
	t reflect.Type

	Name      string
	Arguments []Node
}

type BuiltinNode struct {
	l file.Location
	t reflect.Type

	Name      string
	Arguments []Node
}

type ClosureNode struct {
	l file.Location
	t reflect.Type

	Node Node
}

type PointerNode struct {
	l file.Location
	t reflect.Type
}

type ConditionalNode struct {
	l file.Location
	t reflect.Type

	Cond Node
	Exp1 Node
	Exp2 Node
}

type ArrayNode struct {
	l file.Location
	t reflect.Type

	Nodes []Node
}

type MapNode struct {
	l file.Location
	t reflect.Type

	Pairs []*PairNode
}

type PairNode struct {
	l file.Location
	t reflect.Type

	Key   Node
	Value Node
}
