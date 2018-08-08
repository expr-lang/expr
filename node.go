package expr

import (
	"regexp"
)

// Node represents items of abstract syntax tree.
type Node interface {
	Type(table typesTable) (Type, error)
	Eval(env interface{}) (interface{}, error)
}

type nilNode struct{}

type identifierNode struct {
	value string
}

type numberNode struct {
	value float64
}

type boolNode struct {
	value bool
}

type textNode struct {
	value string
}

type nameNode struct {
	name string
}

type unaryNode struct {
	operator string
	node     Node
}

type binaryNode struct {
	operator string
	left     Node
	right    Node
}

type matchesNode struct {
	r     *regexp.Regexp
	left  Node
	right Node
}

type propertyNode struct {
	node     Node
	property string
}

type indexNode struct {
	node  Node
	index Node
}

type methodNode struct {
	node      Node
	method    string
	arguments []Node
}

type builtinNode struct {
	name      string
	arguments []Node
}

type functionNode struct {
	name      string
	arguments []Node
}

type conditionalNode struct {
	cond Node
	exp1 Node
	exp2 Node
}

type arrayNode struct {
	nodes []Node
}

type mapNode struct {
	pairs []pairNode
}

type pairNode struct {
	key   Node
	value Node
}
