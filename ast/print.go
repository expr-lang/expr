package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/expr-lang/expr/parser/utils"
)

func (n *NilNode) String() string {
	res := "nil"
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *IdentifierNode) String() string {
	res := n.Value
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *IntegerNode) String() string {
	res := fmt.Sprintf("%d", n.Value)
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *FloatNode) String() string {
	res := fmt.Sprintf("%v", n.Value)
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *BoolNode) String() string {
	res := fmt.Sprintf("%t", n.Value)
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *StringNode) String() string {
	res := fmt.Sprintf("%q", n.Value)
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *ConstantNode) String() string {
	res := func() string {
		if n.Value == nil {
			return "nil"
		}
		b, err := json.Marshal(n.Value)
		if err != nil {
			panic(err)
		}
		return string(b)
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *UnaryNode) String() string {
	res := func() string {
		op := n.Operator
		if n.Operator == "not" {
			op = fmt.Sprintf("%s ", n.Operator)
		}
		return fmt.Sprintf("%s%s", op, n.Node.String())
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *BinaryNode) String() string {
	res := func() string {
		if n.Operator == ".." {
			return fmt.Sprintf("%s..%s", n.Left, n.Right)
		}
		return fmt.Sprintf("%s %s %s", n.Left, n.Operator, n.Right)
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *ChainNode) String() string {
	res := n.Node.String()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *MemberNode) String() string {
	res := func() string {
		if n.Optional {
			if str, ok := n.Property.(*StringNode); ok && utils.IsValidIdentifier(str.Value) {
				return fmt.Sprintf("%s?.%s", n.Node.String(), str.Value)
			} else {
				return fmt.Sprintf("%s?.[%s]", n.Node.String(), n.Property.String())
			}
		}
		if str, ok := n.Property.(*StringNode); ok && utils.IsValidIdentifier(str.Value) {
			if _, ok := n.Node.(*PointerNode); ok {
				return fmt.Sprintf(".%s", str.Value)
			}
			return fmt.Sprintf("%s.%s", n.Node.String(), str.Value)
		}
		return fmt.Sprintf("%s[%s]", n.Node.String(), n.Property.String())
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *SliceNode) String() string {
	res := func() string {
		if n.From == nil && n.To == nil {
			return fmt.Sprintf("%s[:]", n.Node.String())
		}
		if n.From == nil {
			return fmt.Sprintf("%s[:%s]", n.Node.String(), n.To.String())
		}
		if n.To == nil {
			return fmt.Sprintf("%s[%s:]", n.Node.String(), n.From.String())
		}
		return fmt.Sprintf("%s[%s:%s]", n.Node.String(), n.From.String(), n.To.String())
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *CallNode) String() string {
	res := func() string {
		arguments := make([]string, len(n.Arguments))
		for i, arg := range n.Arguments {
			arguments[i] = arg.String()
		}
		return fmt.Sprintf("%s(%s)", n.Callee.String(), strings.Join(arguments, ", "))
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *BuiltinNode) String() string {
	res := func() string {
		arguments := make([]string, len(n.Arguments))
		for i, arg := range n.Arguments {
			arguments[i] = arg.String()
		}
		return fmt.Sprintf("%s(%s)", n.Name, strings.Join(arguments, ", "))
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *ClosureNode) String() string {
	res := n.Node.String()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *PointerNode) String() string {
	res := fmt.Sprintf("#%s", n.Name)
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *VariableDeclaratorNode) String() string {
	res := fmt.Sprintf("let %s = %s; %s", n.Name, n.Value.String(), n.Expr.String())
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *ConditionalNode) String() string {
	res := fmt.Sprintf("%s ? %s : %s", n.Cond, n.Exp1, n.Exp2)
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *ArrayNode) String() string {
	res := func() string {
		nodes := make([]string, len(n.Nodes))
		for i, node := range n.Nodes {
			nodes[i] = node.String()
		}
		return fmt.Sprintf("[%s]", strings.Join(nodes, ", "))
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *MapNode) String() string {
	res := func() string {
		pairs := make([]string, len(n.Pairs))
		for i, pair := range n.Pairs {
			pairs[i] = pair.String()
		}
		return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}

func (n *PairNode) String() string {
	res := func() string {
		if str, ok := n.Key.(*StringNode); ok {
			if utils.IsValidIdentifier(str.Value) {
				return fmt.Sprintf("%s: %s", str.Value, n.Value.String())
			}
			return fmt.Sprintf("%q: %s", str.String(), n.Value.String())
		}
		return fmt.Sprintf("%s: %s", n.Key, n.Value)
	}()
	if n.parenthesis() {
		res = fmt.Sprintf("(%s)", res)
	}
	return res
}
