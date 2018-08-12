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
	return fmt.Sprintf("(%v%v)", n.operator, n.node)
}

func (n binaryNode) String() string {
	var leftOp, rightOp *info
	op := binaryOperators[n.operator]

	switch n.left.(type) {
	case binaryNode:
		v := binaryOperators[n.left.(binaryNode).operator]
		leftOp = &v
	}
	switch n.right.(type) {
	case binaryNode:
		v := binaryOperators[n.right.(binaryNode).operator]
		rightOp = &v
	}

	l, r := fmt.Sprintf("%v", n.left), fmt.Sprintf("%v", n.right)

	if leftOp != nil {
		if leftOp.precedence < op.precedence && op.associativity == left {
			l = fmt.Sprintf("(%v)", n.left)
		} else if leftOp.precedence >= op.precedence && op.associativity == right {
			l = fmt.Sprintf("(%v)", n.left)
		}
	}

	if rightOp != nil {
		if rightOp.precedence < op.precedence && op.associativity == left {
			r = fmt.Sprintf("(%v)", n.right)
		} else if rightOp.precedence >= op.precedence && op.associativity == right {
			r = fmt.Sprintf("(%v)", n.right)
		}
	}

	return fmt.Sprintf("%v %v %v", l, n.operator, r)
}

func (n matchesNode) String() string {
	return fmt.Sprintf("(%v matches %v)", n.left, n.right)
}

func (n propertyNode) String() string {
	return fmt.Sprintf("%v.%v", n.node, n.property)
}

func (n indexNode) String() string {
	return fmt.Sprintf("%v[%v]", n.node, n.index)
}

func (n methodNode) String() string {
	s := fmt.Sprintf("%v.%v(", n.node, n.method)
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
	case unaryNode:
		return fmt.Sprintf("%v: %v", n.key, n.value)
	case binaryNode:
		return fmt.Sprintf("(%v): %v", n.key, n.value)
	}
	return fmt.Sprintf("%q: %v", n.key, n.value)
}
