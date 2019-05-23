package checker

import (
	"fmt"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/internal/helper"
	"reflect"
)

func Check(node ast.Node, source *helper.Source, ops ...OptionFn) (t reflect.Type, err error) {
	defer func() {
		if r := recover(); r != nil {
			if h, ok := r.(helper.Error); ok {
				err = fmt.Errorf("%v", h.Format(source))
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	v := &visitor{
		types: make(typesTable),
	}

	for _, op := range ops {
		op(v)
	}

	t = v.visit(node)
	return
}

type visitor struct {
	types typesTable
}

func (v *visitor) visit(node ast.Node) reflect.Type {
	switch n := node.(type) {
	case *ast.NilNode:
		return v.NilNode(n)
	case *ast.IdentifierNode:
		return v.IdentifierNode(n)
	case *ast.IntegerNode:
		return v.IntegerNode(n)
	case *ast.FloatNode:
		return v.FloatNode(n)
	case *ast.BoolNode:
		return v.BoolNode(n)
	case *ast.StringNode:
		return v.StringNode(n)
	case *ast.UnaryNode:
		return v.UnaryNode(n)
	case *ast.BinaryNode:
		return v.BinaryNode(n)
	case *ast.MatchesNode:
		return v.MatchesNode(n)
	case *ast.PropertyNode:
		return v.PropertyNode(n)
	case *ast.IndexNode:
		return v.IndexNode(n)
	case *ast.MethodNode:
		return v.MethodNode(n)
	case *ast.FunctionNode:
		return v.FunctionNode(n)
	case *ast.BuiltinNode:
		return v.BuiltinNode(n)
	case *ast.ClosureNode:
		return v.ClosureNode(n)
	case *ast.PointerNode:
		return v.PointerNode(n)
	case *ast.ConditionalNode:
		return v.ConditionalNode(n)
	case *ast.ArrayNode:
		return v.ArrayNode(n)
	case *ast.MapNode:
		return v.MapNode(n)
	case *ast.PairNode:
		return v.PairNode(n)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
}

func (v *visitor) error(node ast.Node, format string, args ...interface{}) helper.Error {
	return helper.Error{
		Location: node.GetLocation(),
		Message:  fmt.Sprintf(format, args...),
	}
}

func (v *visitor) NilNode(node *ast.NilNode) reflect.Type {
	return nilType
}

func (v *visitor) IdentifierNode(node *ast.IdentifierNode) reflect.Type {
	if t, ok := v.types[node.Value]; ok {
		return t
	}
	panic(v.error(node, "unknown name %v", node.Value))
}

func (v *visitor) IntegerNode(node *ast.IntegerNode) reflect.Type {
	return integerType
}

func (v *visitor) FloatNode(node *ast.FloatNode) reflect.Type {
	return floatType
}

func (v *visitor) BoolNode(node *ast.BoolNode) reflect.Type {
	return boolType
}

func (v *visitor) StringNode(node *ast.StringNode) reflect.Type {
	return stringType
}

func (v *visitor) UnaryNode(node *ast.UnaryNode) reflect.Type {
	t := v.visit(node.Node)

	switch node.Operator {

	case "!", "not":
		if isBool(t) || isInterface(t) {
			return boolType
		}

	case "~":
		if isInteger(t) || isInterface(t) {
			return integerType
		}

	case "+", "-":
		if isInteger(t) || isFloat(t) || isInterface(t) {
			return t
		}

	default:
		panic(v.error(node, "unknown operator (%v)", node.Operator))
	}

	panic(v.error(node, `invalid operation: %v (mismatched type %v)`, node.Operator, t))
}

func (v *visitor) BinaryNode(node *ast.BinaryNode) reflect.Type {
	l := v.visit(node.Left)
	r := v.visit(node.Right)

	switch node.Operator {
	case "==", "!=":
		if isComparable(l, r) {
			return boolType
		}

	case "or", "||", "and", "&&":
		if (isBool(l) || isInterface(l)) && (isBool(r) || isInterface(r)) {
			return boolType
		}

	case "in", "not in":
		if (isString(l) || isInterface(l)) && (isStruct(r) || isInterface(r)) {
			return boolType
		}
		if isArray(r) || isMap(r) || isInterface(r) {
			return boolType
		}

	case "<", ">", ">=", "<=":
		if (isInteger(l) || isInterface(l)) && (isInteger(r) || isInterface(r)) {
			return boolType
		}
		if (isFloat(l) || isInterface(l)) && (isFloat(r) || isInterface(r)) {
			return boolType
		}

	case "/", "-", "*", "**":
		if (isInteger(l) || isInterface(l)) && (isInteger(r) || isInterface(r)) {
			return integerType
		}
		if (isFloat(l) || isInterface(l)) && (isFloat(r) || isInterface(r)) {
			return floatType
		}

	case "|", "^", "&", "%":
		if (isInteger(l) || isInterface(l)) && (isInteger(r) || isInterface(r)) {
			return integerType
		}

	case "+":
		if (isInteger(l) || isInterface(l)) && (isInteger(r) || isInterface(r)) {
			return integerType
		}
		if (isFloat(l) || isInterface(l)) && (isFloat(r) || isInterface(r)) {
			return floatType
		}
		if (isString(l) || isInterface(l)) && (isString(r) || isInterface(r)) {
			return stringType
		}

	case "contains":
		if (isString(l) || isInterface(l)) && (isString(r) || isInterface(r)) {
			return boolType
		}

	case "..":
		if (isInteger(l) || isInterface(l)) && (isInteger(r) || isInterface(r)) {
			return arrayType
		}

	default:
		panic(v.error(node, "unknown operator (%v)", node.Operator))

	}

	panic(v.error(node, `invalid operation: %v (mismatched types %v and %v)`, node.Operator, l, r))
}

func (v *visitor) MatchesNode(node *ast.MatchesNode) reflect.Type {
	l := v.visit(node.Left)
	r := v.visit(node.Right)

	if (isString(l) || isInterface(l)) && (isString(r) || isInterface(r)) {
		return stringType
	}

	panic(v.error(node, `invalid operation: matches (mismatched types %v and %v)`, l, r))
}

func (v *visitor) PropertyNode(node *ast.PropertyNode) reflect.Type {
	t := v.visit(node.Node)

	if t, ok := fieldType(t, node.Property); ok {
		return t
	}

	panic(v.error(node, "type %v has no field %v", t, node.Property))
}

func (v *visitor) IndexNode(node *ast.IndexNode) reflect.Type {
	t := v.visit(node.Node)
	i := v.visit(node.Index)

	if t, ok := indexType(t); ok {
		if !isInteger(i) && !isString(i) {
			panic(v.error(node, "invalid operation: can't use %v as index to %v", i, t))
		}
		return t
	}

	panic(v.error(node, "invalid operation: type %v does not support indexing", t))
}

func (v *visitor) FunctionNode(node *ast.FunctionNode) reflect.Type {
	if f, ok := v.types[node.Name]; ok {
		if fn, ok := funcType(f); ok {
			if isInterface(fn) {
				return integerType
			}

			if fn.NumOut() != 1 {
				panic(v.error(node, "func %v returns more then one value", node.Name))
			}

			numIn := fn.NumIn()

			// If func is method on an env, first argument should be a receiver,
			// and actual arguments less then numIn by one.
			if f.method {
				numIn--
			}

			if len(node.Arguments) > numIn {
				panic(v.error(node, "too many arguments to call %v", node.Name))
			}
			if len(node.Arguments) < numIn {
				panic(v.error(node, "not enough arguments to call %v", node.Name))
			}

			n := 0

			// Skip first argument in case of the receiver.
			if f.method {
				n = 1
			}

			for _, arg := range node.Arguments {
				t := v.visit(arg)
				in := fn.In(n)
				if !t.AssignableTo(in) {
					panic(v.error(node, "can't use %v as argument (type %v) to call %v ", t, in, node.Name))
				}
				n++
			}

			return fn.Out(0)

		}
	}
	panic(v.error(node, "unknown func %v", node.Name))
}

func (v *visitor) MethodNode(node *ast.MethodNode) reflect.Type {
	panic("a")
}

func (v *visitor) BuiltinNode(node *ast.BuiltinNode) reflect.Type {
	panic("imme")
}

func (v *visitor) ClosureNode(node *ast.ClosureNode) reflect.Type {
	panic("imme")
}

func (v *visitor) PointerNode(node *ast.PointerNode) reflect.Type {
	panic("imme")
}

func (v *visitor) ConditionalNode(node *ast.ConditionalNode) reflect.Type {
	panic("imme")
}

func (v *visitor) ArrayNode(node *ast.ArrayNode) reflect.Type {
	panic("imme")
}

func (v *visitor) MapNode(node *ast.MapNode) reflect.Type {
	panic("imme")
}

func (v *visitor) PairNode(node *ast.PairNode) reflect.Type {
	panic("imme")
}
