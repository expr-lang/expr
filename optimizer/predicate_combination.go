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
	c := comparer{visited: map[visit]bool{}}
	return c.same(reflect.ValueOf(a), reflect.ValueOf(b), -1)
}

var constantNodeType = reflect.TypeOf(ConstantNode{})

// maxConstantDepth bounds how deep the comparison walks into a constant's
// value, so a deeply nested value cannot overflow the stack.
const maxConstantDepth = 100

// visit is a pair of pointers, with their type, that the comparison has
// already reached.
type visit struct {
	a, b uintptr
	t    reflect.Type
}

type comparer struct {
	visited map[visit]bool
}

// same compares a and b by type and value. depth is -1 while walking nodes,
// where it skips unexported fields: they hold the location and type, which
// differ between two copies of one expression. Inside a constant's value,
// depth counts the nesting, and every field is compared. Floats compare by
// bits, so 0.0 and -0.0 differ. Values it cannot compare, and constants
// nested deeper than maxConstantDepth, are never the same.
//
// Like reflect.DeepEqual, it remembers each pair of pointers it reaches and
// treats a pair it meets again as the same. Any difference stops the whole
// comparison, so this ends cycles and keeps values that share parts from
// being walked once per path.
func (c *comparer) same(a, b reflect.Value, depth int) bool {
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
	case reflect.Ptr, reflect.Map, reflect.Slice:
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		if a.Kind() == reflect.Slice && a.Len() != b.Len() {
			return false
		}
		v := visit{a.Pointer(), b.Pointer(), a.Type()}
		if c.visited[v] {
			return true
		}
		c.visited[v] = true
	}
	switch a.Kind() {
	case reflect.Interface:
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		return c.same(a.Elem(), b.Elem(), depth)
	case reflect.Ptr:
		return c.same(a.Elem(), b.Elem(), depth)
	case reflect.Struct:
		if depth < 0 && a.Type() == constantNodeType {
			return c.same(a.FieldByName("Value"), b.FieldByName("Value"), 0)
		}
		for i := 0; i < a.NumField(); i++ {
			if depth < 0 && !a.Type().Field(i).IsExported() {
				continue
			}
			if !c.same(a.Field(i), b.Field(i), depth) {
				return false
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			if !c.same(a.Index(i), b.Index(i), depth) {
				return false
			}
		}
		return true
	case reflect.Map:
		if a.Len() != b.Len() {
			return false
		}
		iter := a.MapRange()
		for iter.Next() {
			bv, ok := c.mapIndex(b, iter.Key(), depth)
			if !ok || !c.same(iter.Value(), bv, depth) {
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

// mapIndex finds the value for key in m. Go map lookup uses ==, which finds
// key 0.0 for -0.0, so keys that are not compared exactly by == are matched
// by searching m's keys with same.
func (c *comparer) mapIndex(m, key reflect.Value, depth int) (reflect.Value, bool) {
	if exactKey(key) {
		v := m.MapIndex(key)
		return v, v.IsValid()
	}
	iter := m.MapRange()
	for iter.Next() {
		if c.same(key, iter.Key(), depth) {
			return iter.Value(), true
		}
	}
	return reflect.Value{}, false
}

// exactKey reports whether == on key means the same value: true for bools,
// integers and strings, false for floats and anything that may hold one.
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
