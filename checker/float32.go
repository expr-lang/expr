package checker

import (
	"reflect"

	"github.com/expr-lang/expr/ast"
)

var float32Type = reflect.TypeOf(float32(0))

type float32ComparisonLiterals struct{}

func (float32ComparisonLiterals) Visit(node *ast.Node) {
	binary, ok := (*node).(*ast.BinaryNode)
	if !ok || !isComparisonOperator(binary.Operator) {
		return
	}

	patchFloat32Literal(&binary.Right, binary.Left.Type())
	patchFloat32Literal(&binary.Left, binary.Right.Type())
}

func patchFloat32Literal(node *ast.Node, target reflect.Type) {
	if target != float32Type {
		return
	}

	switch literal := (*node).(type) {
	case *ast.FloatNode:
		ast.Patch(node, newFloat32Constant(float32(literal.Value)))
	case *ast.UnaryNode:
		if literal.Operator == "+" || literal.Operator == "-" {
			if value, ok := literal.Node.(*ast.FloatNode); ok {
				floatValue := float32(value.Value)
				if literal.Operator == "-" {
					floatValue = -floatValue
				}
				ast.Patch(node, newFloat32Constant(floatValue))
			}
		}
	}
}

func newFloat32Constant(value float32) *ast.ConstantNode {
	node := &ast.ConstantNode{Value: value}
	node.SetType(float32Type)
	return node
}

func isComparisonOperator(operator string) bool {
	switch operator {
	case "==", "!=", "<", ">", "<=", ">=":
		return true
	default:
		return false
	}
}
