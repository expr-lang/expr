package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/expr-lang/expr/parser/operator"
	"github.com/expr-lang/expr/parser/utils"
)

// wrapIf returns the node as a string, in parentheses if cond is true.
func wrapIf(cond bool, node Node) string {
	if cond {
		return fmt.Sprintf("(%s)", node.String())
	}
	return node.String()
}

// isSequence reports whether the node is a list of expressions. It can only be
// written at the top level, in parentheses or in braces.
func isSequence(node Node) bool {
	_, ok := node.(*SequenceNode)
	return ok
}

// isStatement reports whether the node is a list of expressions, or a variable
// declaration, which takes everything that follows it.
func isStatement(node Node) bool {
	_, ok := node.(*VariableDeclaratorNode)
	return ok || isSequence(node)
}

// isOpen reports whether the node ends where an operator starts, or the other
// way around, so it can not be an operand without parentheses.
func isOpen(node Node) bool {
	_, ok := node.(*ConditionalNode)
	return ok || isStatement(node)
}

// wrapBase returns the string of a node to which a member, a slice or a call is
// added.
func wrapBase(node Node) string {
	switch node.(type) {
	case *UnaryNode, *BinaryNode, *NilNode, *BoolNode, *IntegerNode, *FloatNode:
		return wrapIf(true, node)
	}
	return wrapIf(isOpen(node), node)
}

// list returns the nodes as a list of elements or arguments.
func list(nodes []Node) string {
	items := make([]string, len(nodes))
	for i, node := range nodes {
		items[i] = wrapIf(isSequence(node), node)
	}
	return strings.Join(items, ", ")
}

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
	str := fmt.Sprintf("%v", n.Value)
	if !strings.ContainsAny(str, ".eIN") {
		str += ".0" // Keep it a float: 2.0 is not 2.
	}
	return str
}

func (n *BoolNode) String() string {
	return fmt.Sprintf("%t", n.Value)
}

func (n *StringNode) String() string {
	return fmt.Sprintf("%q", n.Value)
}

func (n *BytesNode) String() string {
	return fmt.Sprintf("b%q", n.Value)
}

func (n *ConstantNode) String() string {
	if n.Value == nil {
		return "nil"
	}
	b, err := json.Marshal(n.Value)
	if err != nil {
		// json.Marshal rejects values such as NaN and ±Inf, which can
		// reach here after constant folding (e.g. an array literal like
		// [0/0]). Fall back to a non-panicking representation instead of
		// crashing the caller (the optimizer runs this during Compile).
		return fmt.Sprintf("%v", n.Value)
	}
	return string(b)
}

func (n *UnaryNode) String() string {
	op := n.Operator
	if n.Operator == "not" {
		op = fmt.Sprintf("%s ", n.Operator)
	}
	wrap := isOpen(n.Node)
	switch b := n.Node.(type) {
	case *BinaryNode:
		wrap = operator.Binary[b.Operator].Precedence <
			operator.Unary[n.Operator].Precedence
	case *UnaryNode:
		wrap = operator.Unary[b.Operator].Precedence <
			operator.Unary[n.Operator].Precedence
	}
	return op + wrapIf(wrap, n.Node)
}

func (n *BinaryNode) String() string {
	op := operator.Binary[n.Operator]

	lwrap, rwrap := isOpen(n.Left), isOpen(n.Right)

	switch l := n.Left.(type) {
	case *UnaryNode:
		lwrap = operator.Unary[l.Operator].Precedence < op.Precedence
	case *BinaryNode:
		lop := operator.Binary[l.Operator]
		lwrap = operator.Less(l.Operator, n.Operator) ||
			lop.Precedence == op.Precedence && op.Associativity == operator.Right ||
			l.Operator == "??" ||
			operator.IsBoolean(l.Operator) && n.Operator != l.Operator ||
			// a < b < c is a < b && b < c
			operator.IsComparison(l.Operator) && operator.IsComparison(n.Operator)
	}

	switch r := n.Right.(type) {
	case *UnaryNode:
		// Without parentheses, "not" would also take the operators that follow
		// and bind tighter than it does: (a * (not b)) * c.
		rwrap = operator.Unary[r.Operator].Precedence < op.Precedence &&
			operator.Unary[r.Operator].Precedence <= operator.Binary["*"].Precedence
	case *BinaryNode:
		rop := operator.Binary[r.Operator]
		rwrap = operator.Less(r.Operator, n.Operator) ||
			rop.Precedence == op.Precedence && op.Associativity == operator.Left ||
			operator.IsBoolean(r.Operator) && n.Operator != r.Operator
	}

	lhs, rhs := wrapIf(lwrap, n.Left), wrapIf(rwrap, n.Right)

	if n.Operator == ".." {
		if strings.HasPrefix(rhs, ".") {
			rhs = fmt.Sprintf("(%s)", rhs) // 1..(.foo)
		}
		return fmt.Sprintf("%s..%s", lhs, rhs)
	}
	return fmt.Sprintf("%s %s %s", lhs, n.Operator, rhs)
}

func (n *ChainNode) String() string {
	return n.Node.String()
}

func (n *MemberNode) String() string {
	node := wrapBase(n.Node)

	str, ok := n.Property.(*StringNode)
	ok = ok && utils.IsValidIdentifier(str.Value)
	if _, chain := n.Node.(*ChainNode); chain {
		ok = false // A name would be a part of the chain: a?.b["c"] is not a?.b.c.
	}

	if ok {
		if n.Optional {
			return fmt.Sprintf("%s?.%s", node, str.Value)
		}
		if pointer, ok := n.Node.(*PointerNode); ok && pointer.Name == "" {
			return fmt.Sprintf(".%s", str.Value)
		}
		return fmt.Sprintf("%s.%s", node, str.Value)
	}

	property := wrapIf(isSequence(n.Property), n.Property)
	if n.Optional {
		return fmt.Sprintf("%s?.[%s]", node, property)
	}
	return fmt.Sprintf("%s[%s]", node, property)
}

func (n *SliceNode) String() string {
	op := ""
	if n.Optional {
		op = "?."
	}
	var from, to string
	if n.From != nil {
		from = wrapIf(isOpen(n.From), n.From)
	}
	if n.To != nil {
		to = wrapIf(isOpen(n.To), n.To)
	}
	return fmt.Sprintf("%s%s[%s:%s]", wrapBase(n.Node), op, from, to)
}

func (n *CallNode) String() string {
	return fmt.Sprintf("%s(%s)", wrapBase(n.Callee), list(n.Arguments))
}

func (n *BuiltinNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Name, list(n.Arguments))
}

func (n *PredicateNode) String() string {
	str := n.Node.String()
	// The braces of a map would be taken as those of the predicate.
	if isSequence(n.Node) || strings.HasPrefix(str, "{") {
		return fmt.Sprintf("{ %s }", str)
	}
	return str
}

func (n *PointerNode) String() string {
	return fmt.Sprintf("#%s", n.Name)
}

func (n *VariableDeclaratorNode) String() string {
	return fmt.Sprintf("let %s = %s; %s", n.Name, wrapIf(isStatement(n.Value), n.Value), n.Expr.String())
}

func (n *SequenceNode) String() string {
	nodes := make([]string, len(n.Nodes))
	for i, node := range n.Nodes {
		last := i == len(n.Nodes)-1
		nodes[i] = wrapIf(isSequence(node) || isStatement(node) && !last, node)
	}
	return strings.Join(nodes, "; ")
}

func (n *ConditionalNode) String() string {
	if !n.Ternary {
		cond := wrapIf(isStatement(n.Cond), n.Cond)
		exp1 := n.Exp1.String()
		if c2, ok := n.Exp2.(*ConditionalNode); ok && !c2.Ternary {
			return fmt.Sprintf("if %s { %s } else %s", cond, exp1, c2.String())
		}
		exp2 := n.Exp2.String()
		return fmt.Sprintf("if %s { %s } else { %s }", cond, exp1, exp2)
	}

	return fmt.Sprintf("%s ? %s : %s",
		wrapIf(isOpen(n.Cond), n.Cond),
		wrapIf(isOpen(n.Exp1), n.Exp1),
		wrapIf(isOpen(n.Exp2), n.Exp2),
	)
}

func (n *ArrayNode) String() string {
	return fmt.Sprintf("[%s]", list(n.Nodes))
}

func (n *MapNode) String() string {
	pairs := make([]string, len(n.Pairs))
	for i, pair := range n.Pairs {
		pairs[i] = pair.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (n *PairNode) String() string {
	value := wrapIf(isSequence(n.Value), n.Value)
	if str, ok := n.Key.(*StringNode); ok {
		if utils.IsValidIdentifier(str.Value) && !isKeyword(str.Value) {
			return fmt.Sprintf("%s: %s", str.Value, value)
		}
		return fmt.Sprintf("%s: %s", str.String(), value)
	}
	return fmt.Sprintf("(%s): %s", n.Key.String(), value)
}

// isKeyword reports whether the word is not read as an identifier.
func isKeyword(word string) bool {
	switch word {
	case "not", "in", "or", "and", "matches", "contains", "startsWith", "endsWith", "let", "if", "else":
		return true
	}
	return false
}
