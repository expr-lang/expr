package ast

import (
	"github.com/antonmedv/expr/internal/helper"
	"regexp"
)

// Node represents items of abstract syntax tree.
type Node interface {
	GetLocation() helper.Location
	SetLocation(helper.Location)
}

type NilNode struct {
	location helper.Location
}

type IdentifierNode struct {
	location helper.Location
	Value    string
}

type IntegerNode struct {
	location helper.Location
	Value    int64
}

type FloatNode struct {
	location helper.Location
	Value    float64
}

type BoolNode struct {
	location helper.Location
	Value    bool
}

type StringNode struct {
	location helper.Location
	Value    string
}

type UnaryNode struct {
	location helper.Location
	Operator string
	Node     Node
}

type BinaryNode struct {
	location helper.Location
	Operator string
	Left     Node
	Right    Node
}

type MatchesNode struct {
	location helper.Location
	Regexp   *regexp.Regexp
	Left     Node
	Right    Node
}

type PropertyNode struct {
	location helper.Location
	Node     Node
	Property string
}

type IndexNode struct {
	location helper.Location
	Node     Node
	Index    Node
}

type MethodNode struct {
	location  helper.Location
	Node      Node
	Method    string
	Arguments []Node
}

type FunctionNode struct {
	location  helper.Location
	Name      string
	Arguments []Node
}

type BuiltinNode struct {
	location  helper.Location
	Name      string
	Arguments []Node
}

type ClosureNode struct {
	location helper.Location
	Node     Node
}

type PointerNode struct {
	location helper.Location
}

type ConditionalNode struct {
	location helper.Location
	Cond     Node
	Exp1     Node
	Exp2     Node
}

type ArrayNode struct {
	location helper.Location
	Nodes    []Node
}

type MapNode struct {
	location helper.Location
	Pairs    []*PairNode
}

type PairNode struct {
	location helper.Location
	Key      Node
	Value    Node
}
