package checker

import (
	"fmt"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/internal/helper"
	"reflect"
)

type TypesTable map[string]reflect.Type

func Check(node ast.Node, table TypesTable, source *helper.Source) (reflect.Type, error) {
	v := &visitor{
		types:  table,
		errors: helper.NewErrors(source),
	}
	ast.Walk(node, v)

	if v.errors.HasError() {
		return nil, v.errors
	}
	if len(v.stack) == 0 {
		return nil, fmt.Errorf("empty stack")
	}
	if len(v.stack) > 1 {
		return nil, fmt.Errorf("too long stack")
	}
	return v.stack[0], nil
}

type visitor struct {
	types  TypesTable
	stack  []reflect.Type
	errors *helper.Errors
}

func (v *visitor) push(node reflect.Type) {
	v.stack = append(v.stack, node)
}

func (v *visitor) pop() reflect.Type {
	if len(v.stack) == 0 {
		return nilType
	}
	node := v.stack[len(v.stack)-1]
	v.stack = v.stack[:len(v.stack)-1]
	return node
}

func (v *visitor) reportError(node ast.Node, format string, args ...interface{}) {
	v.errors.ReportError(node.GetLocation(), format, args...)
}

func (v *visitor) NilNode(node *ast.NilNode) {
	v.push(nilType)
}

func (v *visitor) IdentifierNode(node *ast.IdentifierNode) {
	if t, ok := v.types[node.Value]; ok {
		v.push(t)
	} else {
		v.reportError(node, "unknown name %v", node.Value)
	}
}

func (v *visitor) IntegerNode(node *ast.IntegerNode) {

}

func (v *visitor) FloatNode(node *ast.FloatNode) {

}

func (v *visitor) BoolNode(node *ast.BoolNode) {

}

func (v *visitor) StringNode(node *ast.StringNode) {

}

func (v *visitor) UnaryNode(node *ast.UnaryNode) {

}

func (v *visitor) BinaryNode(node *ast.BinaryNode) {
	v.reportError(node.Right, "left is %v", node.Left)
}

func (v *visitor) MatchesNode(node *ast.MatchesNode) {

}

func (v *visitor) PropertyNode(node *ast.PropertyNode) {

}

func (v *visitor) IndexNode(node *ast.IndexNode) {

}

func (v *visitor) MethodNode(node *ast.MethodNode) {

}

func (v *visitor) FunctionNode(node *ast.FunctionNode) {

}

func (v *visitor) BuiltinNode(node *ast.BuiltinNode) {

}

func (v *visitor) ClosureNode(node *ast.ClosureNode) {

}

func (v *visitor) PointerNode(node *ast.PointerNode) {

}

func (v *visitor) ConditionalNode(node *ast.ConditionalNode) {

}

func (v *visitor) ArrayNode(node *ast.ArrayNode) {

}

func (v *visitor) MapNode(node *ast.MapNode) {

}

func (v *visitor) PairNode(node *ast.PairNode) {

}
