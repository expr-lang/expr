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
//
// A false result only skips an optimization, so the comparison is
// conservative: anything it cannot compare simply and quickly is different.
func sameNode(a, b Node) bool {
	c := comparer{steps: maxCompareSteps}
	return c.same(reflect.ValueOf(a), reflect.ValueOf(b), false)
}

var constantNodeType = reflect.TypeOf(ConstantNode{})

// maxCompareSteps bounds the work of one comparison. Constant values can
// refer to themselves or reuse one value in many places, and walking them
// could otherwise take forever.
const maxCompareSteps = 100_000

type comparer struct {
	steps int
}

// same compares a and b by type and value. While walking nodes, it skips
// unexported fields: they hold the location and type, which differ between
// two copies of one expression. Inside a constant's value it compares every
// field, and allows only plain data: no pointers, funcs or channels, and only
// map keys that == compares exactly. Floats compare by bits, so 0.0 and -0.0
// differ.
func (c *comparer) same(a, b reflect.Value, constant bool) bool {
	c.steps--
	if c.steps < 0 {
		return false
	}
	if !a.IsValid() || !b.IsValid() {
		return a.IsValid() == b.IsValid()
	}
	if a.Type() != b.Type() {
		return false
	}
	switch a.Kind() {
	case reflect.Interface:
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		return c.same(a.Elem(), b.Elem(), constant)
	case reflect.Ptr:
		if constant {
			return false
		}
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		return c.same(a.Elem(), b.Elem(), constant)
	case reflect.Struct:
		if !constant && a.Type() == constantNodeType {
			return c.same(a.FieldByName("Value"), b.FieldByName("Value"), true)
		}
		for i := 0; i < a.NumField(); i++ {
			if !constant && !a.Type().Field(i).IsExported() {
				continue
			}
			if !c.same(a.Field(i), b.Field(i), constant) {
				return false
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		if a.Kind() == reflect.Slice && a.IsNil() != b.IsNil() {
			return false
		}
		if a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			if !c.same(a.Index(i), b.Index(i), constant) {
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
			if !exactKey(iter.Key()) {
				return false
			}
			bv := b.MapIndex(iter.Key())
			if !bv.IsValid() || !c.same(iter.Value(), bv, constant) {
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

// exactKey reports whether == on a map key means the same value: true for
// bools, integers and strings, false for floats, pointers and anything else.
// With exact keys, the same key in two maps is found by a plain lookup.
func exactKey(key reflect.Value) bool {
	switch key.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	case reflect.Interface:
		return !key.IsNil() && exactKey(key.Elem())
	default:
		return false
	}
}
