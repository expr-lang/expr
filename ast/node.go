package ast

import (
	"github.com/antonmedv/expr/internal/helper"
	"reflect"
	"regexp"
)

// Node represents items of abstract syntax tree.
type Node interface {
	GetLocation() helper.Location
	SetLocation(helper.Location)
	GetType() reflect.Type
	SetType(reflect.Type)
}

type NilNode struct {
	location helper.Location
	t        reflect.Type
}

type IdentifierNode struct {
	location helper.Location
	t        reflect.Type
	Value    string
}

type IntegerNode struct {
	location helper.Location
	t        reflect.Type
	Value    int64
}

type FloatNode struct {
	location helper.Location
	t        reflect.Type
	Value    float64
}

type BoolNode struct {
	location helper.Location
	t        reflect.Type
	Value    bool
}

type StringNode struct {
	location helper.Location
	t        reflect.Type
	Value    string
}

type UnaryNode struct {
	location helper.Location
	t        reflect.Type
	Operator string
	Node     Node
}

type BinaryNode struct {
	location helper.Location
	t        reflect.Type
	Operator string
	Left     Node
	Right    Node
}

type MatchesNode struct {
	location helper.Location
	t        reflect.Type
	Regexp   *regexp.Regexp
	Left     Node
	Right    Node
}

type PropertyNode struct {
	location helper.Location
	t        reflect.Type
	Node     Node
	Property string
}

type IndexNode struct {
	location helper.Location
	t        reflect.Type
	Node     Node
	Index    Node
}

type MethodNode struct {
	location  helper.Location
	t         reflect.Type
	Node      Node
	Method    string
	Arguments []Node
}

type FunctionNode struct {
	location  helper.Location
	t         reflect.Type
	Name      string
	Arguments []Node
}

type BuiltinNode struct {
	location  helper.Location
	t         reflect.Type
	Name      string
	Arguments []Node
}

type ClosureNode struct {
	location helper.Location
	t        reflect.Type
	Node     Node
}

type PointerNode struct {
	location helper.Location
	t        reflect.Type
}

type ConditionalNode struct {
	location helper.Location
	t        reflect.Type
	Cond     Node
	Exp1     Node
	Exp2     Node
}

type ArrayNode struct {
	location helper.Location
	t        reflect.Type
	Nodes    []Node
}

type MapNode struct {
	location helper.Location
	t        reflect.Type
	Pairs    []*PairNode
}

type PairNode struct {
	location helper.Location
	t        reflect.Type
	Key      Node
	Value    Node
}
