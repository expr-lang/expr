package expr

import (
	"fmt"
	"strconv"
)

func (n nilNode) String() string {
	return "nil"
}

func (n identifierNode) String() string {
	return n.value
}

func (n numberNode) String() string {
	return fmt.Sprintf("%v", n.value)
}

func (n boolNode) String() string {
	if n.value {
		return "true"
	}
	return "false"
}

func (n textNode) String() string {
	return strconv.Quote(n.value)
}

func (n nameNode) String() string {
	return n.name
}

func (n unaryNode) String() string {
	switch n.operator {
	case "!", "not":
		return fmt.Sprintf("%v %v", n.operator, n.node)
	}
	return fmt.Sprintf("%v%v", n.operator, n.node)
}

func (n binaryNode) String() string {
	return fmt.Sprintf("%v %v %v", n.left, n.operator, n.right)
}

func (n propertyNode) String() string {
	switch n.property.(type) {
	case identifierNode:
		return fmt.Sprintf("%v.%v", n.node, n.property)
	}
	return fmt.Sprintf("%v[%v]", n.node, n.property)
}

func (n methodNode) String() string {
	s := fmt.Sprintf("%v.%v(", n.node, n.property)
	for i, a := range n.arguments {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", a)
	}
	return s + ")"
}

func (n builtinNode) String() string {
	s := fmt.Sprintf("%v(", n.name)
	for i, a := range n.arguments {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", a)
	}
	return s + ")"
}

func (n functionNode) String() string {
	s := fmt.Sprintf("%v(", n.name)
	for i, a := range n.arguments {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", a)
	}
	return s + ")"
}

func (n conditionalNode) String() string {
	return fmt.Sprintf("%v ? %v : %v", n.cond, n.exp1, n.exp2)
}

func (n arrayNode) String() string {
	s := "["
	for i, n := range n.nodes {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", n)
	}
	return s + "]"
}

func (n mapNode) String() string {
	s := "{"
	for i, p := range n.pairs {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", p)
	}
	return s + "}"
}

func (n pairNode) String() string {
	switch n.key.(type) {
	case binaryNode, unaryNode:
		return fmt.Sprintf("(%v): %v", n.key, n.value)
	}
	return fmt.Sprintf("%q: %v", n.key, n.value)
}
