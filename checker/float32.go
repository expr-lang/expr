package checker

import (
	"reflect"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/internal/deref"
)

var float32Type = reflect.TypeOf(float32(0))

// float32ComparisonLiterals converts float literals compared with a float32
// operand to float32 constants, as Go converts an untyped constant. Otherwise
// the float64 literal is compared with the rounded float32 value, and
// `f32 == 12.34` is false.
type float32ComparisonLiterals struct{}

func (float32ComparisonLiterals) Visit(node *ast.Node) {
	binary, ok := (*node).(*ast.BinaryNode)
	if !ok {
		return
	}
	switch binary.Operator {
	case "==", "!=", "<", ">", "<=", ">=":
		patchFloat32Literal(&binary.Right, binary.Left.Type())
		patchFloat32Literal(&binary.Left, binary.Right.Type())
	case "in":
		if array, ok := binary.Right.(*ast.ArrayNode); ok {
			for i := range array.Nodes {
				patchFloat32Literal(&array.Nodes[i], binary.Left.Type())
			}
		}
		if t := deref.Type(binary.Right.Type()); t != nil && (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) {
			patchFloat32Literal(&binary.Left, t.Elem())
		}
	}
}

func patchFloat32Literal(node *ast.Node, target reflect.Type) {
	if deref.Type(target) != float32Type {
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
