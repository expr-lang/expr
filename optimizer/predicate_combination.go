package optimizer

import (
	. "github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser/operator"
)

type predicateCombination struct{}

func (v *predicateCombination) Visit(node *Node) {
	if op, ok := (*node).(*BinaryNode); ok && operator.IsBoolean(op.Operator) {
		if left, ok := op.Left.(*BuiltinNode); ok {
			if combinedOp, ok := combinedOperator(left.Name, op.Operator); ok {
				if right, ok := op.Right.(*BuiltinNode); ok && right.Name == left.Name {
					if left.Arguments[0].Type() == right.Arguments[0].Type() && left.Arguments[0].String() == right.Arguments[0].String() {
						closure := &ClosureNode{
							Node: &BinaryNode{
								Operator: combinedOp,
								Left:     left.Arguments[1].(*ClosureNode).Node,
								Right:    right.Arguments[1].(*ClosureNode).Node,
							},
						}
						v.Visit(&closure.Node)
						Patch(node, &BuiltinNode{
							Name: left.Name,
							Arguments: []Node{
								left.Arguments[0],
								closure,
							},
						})
					}
				}
			}
		}
	}
}

func combinedOperator(fn, op string) (string, bool) {
	switch fn {
	case "all", "any":
		return op, true
	case "one", "none":
		switch op {
		case "and":
			return "or", true
		case "&&":
			return "||", true
		}
	}
	return "", false
}
