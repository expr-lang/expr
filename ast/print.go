package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/expr-lang/expr/parser/operator"
	"github.com/expr-lang/expr/parser/utils"
)

var EnableCache bool

func (n *NilNode) String() string {
	return "nil"
}

func (n *IdentifierNode) String() string {
	return n.Value
}

func (n *IntegerNode) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *FloatNode) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n *BoolNode) String() string {
	return fmt.Sprintf("%t", n.Value)
}

func (n *StringNode) String() string {
	return fmt.Sprintf("%q", n.Value)
}

func (n *ConstantNode) String() string {
	if n.Value == nil {
		return "nil"
	}
	b, err := json.Marshal(n.Value)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (n *UnaryNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	op := n.Operator
	if n.Operator == "not" {
		op = fmt.Sprintf("%s ", n.Operator)
	}
	if _, ok := n.Node.(*BinaryNode); ok {
		n.strCache = fmt.Sprintf("%s(%s)", op, n.Node.String())
	} else {
		n.strCache = fmt.Sprintf("%s%s", op, n.Node.String())
	}
	return n.strCache
}

func (n *BinaryNode) String() string {
	if EnableCache &&  n.strCache != "" {
		return n.strCache
	}
	if n.Operator == ".." {
		n.strCache = fmt.Sprintf("%s..%s", n.Left, n.Right)
		return n.strCache
	}

	var lhs, rhs string
	var lwrap, rwrap bool

	lb, ok := n.Left.(*BinaryNode)
	if ok {
		if operator.Less(lb.Operator, n.Operator) {
			lwrap = true
		}
		if lb.Operator == "??" {
			lwrap = true
		}
		if operator.IsBoolean(lb.Operator) && n.Operator != lb.Operator {
			lwrap = true
		}
	}

	rb, ok := n.Right.(*BinaryNode)
	if ok {
		if operator.Less(rb.Operator, n.Operator) {
			rwrap = true
		}
		if operator.IsBoolean(rb.Operator) && n.Operator != rb.Operator {
			rwrap = true
		}
	}

	if lwrap {
		lhs = fmt.Sprintf("(%s)", n.Left.String())
	} else {
		lhs = n.Left.String()
	}

	if rwrap {
		rhs = fmt.Sprintf("(%s)", n.Right.String())
	} else {
		rhs = n.Right.String()
	}

	n.strCache = fmt.Sprintf("%s %s %s", lhs, n.Operator, rhs)
	return n.strCache
}

func (n *ChainNode) String() string {
	return n.Node.String()
}

func (n *MemberNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	if n.Optional {
		if str, ok := n.Property.(*StringNode); ok && utils.IsValidIdentifier(str.Value) {
			n.strCache = fmt.Sprintf("%s?.%s", n.Node.String(), str.Value)
		} else {
			n.strCache = fmt.Sprintf("%s?.[%s]", n.Node.String(), n.Property.String())
		}
	} else if str, ok := n.Property.(*StringNode); ok && utils.IsValidIdentifier(str.Value) {
		if _, ok := n.Node.(*PointerNode); ok {
			n.strCache = fmt.Sprintf(".%s", str.Value)
		} else {
			n.strCache = fmt.Sprintf("%s.%s", n.Node.String(), str.Value)
		}
	} else {
		n.strCache = fmt.Sprintf("%s[%s]", n.Node.String(), n.Property.String())
	}
	return n.strCache
}

func (n *SliceNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	if n.From == nil && n.To == nil {
		n.strCache = fmt.Sprintf("%s[:]", n.Node.String())
	} else if n.From == nil {
		n.strCache = fmt.Sprintf("%s[:%s]", n.Node.String(), n.To.String())
	} else if n.To == nil {
		n.strCache = fmt.Sprintf("%s[%s:]", n.Node.String(), n.From.String())
	} else {
		n.strCache = fmt.Sprintf("%s[%s:%s]", n.Node.String(), n.From.String(), n.To.String())
	}
	return n.strCache
}

func (n *CallNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	arguments := make([]string, len(n.Arguments))
	for i, arg := range n.Arguments {
		arguments[i] = arg.String()
	}
	n.strCache = fmt.Sprintf("%s(%s)", n.Callee.String(), strings.Join(arguments, ", "))
	return n.strCache
}

func (n *BuiltinNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	arguments := make([]string, len(n.Arguments))
	for i, arg := range n.Arguments {
		arguments[i] = arg.String()
	}
	n.strCache = fmt.Sprintf("%s(%s)", n.Name, strings.Join(arguments, ", "))
	return n.strCache
}

func (n *ClosureNode) String() string {
	return n.Node.String()
}

func (n *PointerNode) String() string {
	return fmt.Sprintf("#%s", n.Name)
}

func (n *VariableDeclaratorNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	n.strCache = fmt.Sprintf("let %s = %s; %s", n.Name, n.Value.String(), n.Expr.String())
	return n.strCache
}

func (n *ConditionalNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	var cond, exp1, exp2 string
	if _, ok := n.Cond.(*ConditionalNode); ok {
		cond = fmt.Sprintf("(%s)", n.Cond.String())
	} else {
		cond = n.Cond.String()
	}
	if _, ok := n.Exp1.(*ConditionalNode); ok {
		exp1 = fmt.Sprintf("(%s)", n.Exp1.String())
	} else {
		exp1 = n.Exp1.String()
	}
	if _, ok := n.Exp2.(*ConditionalNode); ok {
		exp2 = fmt.Sprintf("(%s)", n.Exp2.String())
	} else {
		exp2 = n.Exp2.String()
	}
	n.strCache = fmt.Sprintf("%s ? %s : %s", cond, exp1, exp2)
	return n.strCache
}

func (n *ArrayNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	nodes := make([]string, len(n.Nodes))
	for i, node := range n.Nodes {
		nodes[i] = node.String()
	}
	n.strCache = fmt.Sprintf("[%s]", strings.Join(nodes, ", "))
	return n.strCache
}

func (n *MapNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	pairs := make([]string, len(n.Pairs))
	for i, pair := range n.Pairs {
		pairs[i] = pair.String()
	}
	n.strCache = fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
	return n.strCache
}

func (n *PairNode) String() string {
	if EnableCache && n.strCache != "" {
		return n.strCache
	}
	if str, ok := n.Key.(*StringNode); ok {
		if utils.IsValidIdentifier(str.Value) {
			n.strCache = fmt.Sprintf("%s: %s", str.Value, n.Value.String())
		} else {
			n.strCache = fmt.Sprintf("%q: %s", str.String(), n.Value.String())
		}
	} else {
		n.strCache = fmt.Sprintf("(%s): %s", n.Key.String(), n.Value.String())
	}
	return n.strCache
}
