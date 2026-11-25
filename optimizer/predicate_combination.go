package optimizer

import (
	"math"
	"reflect"

	. "github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser/operator"
)

/*
predicateCombination is a visitor that combines multiple predicate calls into a single call.
For example, the following expression:

	all(x, x > 1) && all(x, x < 10) -> all(x, x > 1 && x < 10)
	any(x, x > 1) || any(x, x < 10) -> any(x, x > 1 || x < 10)
	none(x, x > 1) && none(x, x < 10) -> none(x, x > 1 || x < 10)
*/
type predicateCombination struct{}

func (v *predicateCombination) Visit(node *Node) {
	if op, ok := (*node).(*BinaryNode); ok && operator.IsBoolean(op.Operator) {
		if left, ok := op.Left.(*BuiltinNode); ok {
			if combinedOp, ok := combinedOperator(left.Name, op.Operator); ok {
				if right, ok := op.Right.(*BuiltinNode); ok && right.Name == left.Name {
					if left.Arguments[0].Type() == right.Arguments[0].Type() && sameNode(left.Arguments[0], right.Arguments[0]) {
						predicate := &PredicateNode{
							Node: &BinaryNode{
								Operator: combinedOp,
								Left:     left.Arguments[1].(*PredicateNode).Node,
								Right:    right.Arguments[1].(*PredicateNode).Node,
							},
						}
						v.Visit(&predicate.Node)
						patchCopyType(node, &BuiltinNode{
							Name: left.Name,
							Arguments: []Node{
								left.Arguments[0],
								predicate,
							},
						})
					}
				}
			}
		}
	}
}

func combinedOperator(fn, op string) (string, bool) {
	switch {
	case fn == "all" && (op == "and" || op == "&&"):
		return op, true
	case fn == "any" && (op == "or" || op == "||"):
		return op, true
	case fn == "none" && (op == "and" || op == "&&"):
		switch op {
		case "and":
			return "or", true
		case "&&":
			return "||", true
		}
	}
	return "", false
}

// sameNode reports whether a and b are the same expression. It compares the
// trees field by field instead of comparing String() output, because different
// trees can print the same: a BuiltinNode's Map field is not printed, and a
// folded ConstantNode prints 2.0 as 2.
func sameNode(a, b Node) bool {
	return sameValue(reflect.ValueOf(a), reflect.ValueOf(b), -1)
}

var constantNodeType = reflect.TypeOf(ConstantNode{})

// maxConstantDepth bounds how deep sameValue walks into a constant's value,
// so a value that refers to itself cannot recurse forever.
const maxConstantDepth = 100

// sameValue compares a and b by type and value. depth is -1 while walking
// nodes, where it skips unexported fields: they hold the location and type,
// which differ between two copies of one expression. Inside a constant's
// value, depth counts the nesting, and every field is compared. Floats
// compare by bits, so 0.0 and -0.0 differ. Values it cannot compare, and
// constants nested deeper than maxConstantDepth, are never the same.
func sameValue(a, b reflect.Value, depth int) bool {
	if depth >= 0 {
		depth++
		if depth > maxConstantDepth {
			return false
		}
	}
	if !a.IsValid() || !b.IsValid() {
		return a.IsValid() == b.IsValid()
	}
	if a.Type() != b.Type() {
		return false
	}
	switch a.Kind() {
	case reflect.Interface, reflect.Ptr:
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		return sameValue(a.Elem(), b.Elem(), depth)
	case reflect.Struct:
		if depth < 0 && a.Type() == constantNodeType {
			return sameValue(a.FieldByName("Value"), b.FieldByName("Value"), 0)
		}
		for i := 0; i < a.NumField(); i++ {
			if depth < 0 && !a.Type().Field(i).IsExported() {
				continue
			}
			if !sameValue(a.Field(i), b.Field(i), depth) {
				return false
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			if !sameValue(a.Index(i), b.Index(i), depth) {
				return false
			}
		}
		return true
	case reflect.Map:
		if a.IsNil() != b.IsNil() || a.Len() != b.Len() {
			return false
		}
		iter := a.MapRange()
		for iter.Next() {
			bv := b.MapIndex(iter.Key())
			if !bv.IsValid() || !sameValue(iter.Value(), bv, depth) {
				return false
			}
		}
		return true
	case reflect.Bool:
		return a.Bool() == b.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() == b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() == b.Uint()
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(a.Float()) == math.Float64bits(b.Float())
	case reflect.Complex64, reflect.Complex128:
		ac, bc := a.Complex(), b.Complex()
		return math.Float64bits(real(ac)) == math.Float64bits(real(bc)) &&
			math.Float64bits(imag(ac)) == math.Float64bits(imag(bc))
	case reflect.String:
		return a.String() == b.String()
	default:
		return false
	}
}
