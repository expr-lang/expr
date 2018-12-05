package expr

import (
	"fmt"
	"regexp"
)

type Visitor interface {
	Nil() (interface{}, error)
	Value(interface{}) (interface{}, error)
	Name(string) (interface{}, error)
	Unary(string, Node) (interface{}, error)
	Binary(string, Node, Node) (interface{}, error)
	Matches(*regexp.Regexp, Node, Node) (interface{}, error)
	Property(Node, string) (interface{}, error)
	Index(Node, Node) (interface{}, error)
	Method(Node, string, []Node) (interface{}, error)
	Function(string, []Node) (interface{}, error)
	Len(Node) (interface{}, error)
	Conditional(Node, Node, Node) (interface{}, error)
	Array([]Node) (interface{}, error)
	Map([]PairNode) (interface{}, error)
}

// Visit traverses given ast.
func Visit(node Node, visitor Visitor) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	return node.Visit(visitor)
}

func (n nilNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Nil()
}

func (n identifierNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Value(n.value)
}

func (n numberNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Value(n.value)
}

func (n boolNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Value(n.value)
}

func (n textNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Value(n.value)
}

func (n nameNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Name(n.name)
}

func (n unaryNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Unary(n.operator, n.node)
}

func (n binaryNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Binary(n.operator, n.left, n.right)
}

func (n matchesNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Matches(n.r, n.left, n.right)
}

func (n propertyNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Property(n.node, n.property)
}

func (n indexNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Index(n.node, n.index)
}

func (n methodNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Method(n.node, n.method, n.arguments)
}

func (n builtinNode) Visit(visitor Visitor) (interface{}, error) {
	switch n.name {
	case "len":
		if len(n.arguments) == 0 {
			return nil, fmt.Errorf("missing argument: %v", n)
		}
		if len(n.arguments) > 1 {
			return nil, fmt.Errorf("too many arguments: %v", n)
		}
		return visitor.Len(n.arguments[0])
	}

	return nil, fmt.Errorf("unknown %q builtin", n.name)
}

func (n functionNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Function(n.name, n.arguments)
}

func (n conditionalNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Conditional(n.cond, n.exp1, n.exp2)
}

func (n arrayNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Array(n.nodes)
}

func (n mapNode) Visit(visitor Visitor) (interface{}, error) {
	return visitor.Map(n.pairs)
}
