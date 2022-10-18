package checker

import (
	"github.com/antonmedv/expr/conf"
	"reflect"
	"time"

	"github.com/antonmedv/expr/ast"
)

var (
	nilType      = reflect.TypeOf(nil)
	boolType     = reflect.TypeOf(true)
	integerType  = reflect.TypeOf(0)
	floatType    = reflect.TypeOf(float64(0))
	stringType   = reflect.TypeOf("")
	arrayType    = reflect.TypeOf([]interface{}{})
	mapType      = reflect.TypeOf(map[string]interface{}{})
	anyType      = reflect.TypeOf(new(interface{})).Elem()
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Duration(0))
	errorType    = reflect.TypeOf((*error)(nil)).Elem()
)

func typeWeight(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Uint:
		return 1
	case reflect.Uint8:
		return 2
	case reflect.Uint16:
		return 3
	case reflect.Uint32:
		return 4
	case reflect.Uint64:
		return 5
	case reflect.Int:
		return 6
	case reflect.Int8:
		return 7
	case reflect.Int16:
		return 8
	case reflect.Int32:
		return 9
	case reflect.Int64:
		return 10
	case reflect.Float32:
		return 11
	case reflect.Float64:
		return 12
	default:
		return 0
	}
}

func combined(a, b reflect.Type) reflect.Type {
	if typeWeight(a) > typeWeight(b) {
		return a
	} else {
		return b
	}
}

func anyOf(t reflect.Type, fns ...func(reflect.Type) bool) bool {
	for _, fn := range fns {
		if fn(t) {
			return true
		}
	}
	return false
}

func or(l, r reflect.Type, fns ...func(reflect.Type) bool) bool {
	if isAny(l) && isAny(r) {
		return true
	}
	if isAny(l) && anyOf(r, fns...) {
		return true
	}
	if isAny(r) && anyOf(l, fns...) {
		return true
	}
	return false
}

func isAny(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Interface:
			return true
		}
	}
	return false
}

func isInteger(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		}
	}
	return false
}

func isFloat(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Float32, reflect.Float64:
			return true
		}
	}
	return false
}

func isNumber(t reflect.Type) bool {
	return isInteger(t) || isFloat(t)
}

func isTime(t reflect.Type) bool {
	if t != nil {
		switch t {
		case timeType:
			return true
		}
	}
	return isAny(t)
}

func isDuration(t reflect.Type) bool {
	if t != nil {
		switch t {
		case durationType:
			return true
		}
	}
	return false
}

func isBool(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Bool:
			return true
		}
	}
	return false
}

func isString(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.String:
			return true
		}
	}
	return false
}

func isArray(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return isArray(t.Elem())
		case reflect.Slice, reflect.Array:
			return true
		}
	}
	return false
}

func isMap(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return isMap(t.Elem())
		case reflect.Map:
			return true
		}
	}
	return false
}

func isStruct(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return isStruct(t.Elem())
		case reflect.Struct:
			return true
		}
	}
	return false
}

func isFunc(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return isFunc(t.Elem())
		case reflect.Func:
			return true
		}
	}
	return false
}

func fetchType(t reflect.Type, name string) (reflect.Type, bool) {
	if t != nil {
		switch t.Kind() {
		case reflect.Interface:
			return anyType, true
		case reflect.Struct:
			// First check all structs fields.
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				// Search all fields, even embedded structs.
				if conf.FieldName(f) == name {
					return f.Type, true
				}
			}

			// Second check fields of embedded structs.
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Anonymous {
					if t, ok := fetchType(f.Type, name); ok {
						return t, true
					}
				}
			}
		}
	}

	return nil, false
}

func deref(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() == reflect.Interface {
		return t, true
	}
	found := false
	for t != nil && t.Kind() == reflect.Ptr {
		e := t.Elem()
		switch e.Kind() {
		case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
			return t, false
		default:
			found = true
			t = e
		}
	}
	return t, found
}

func isIntegerOrArithmeticOperation(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.IntegerNode:
		return true
	case *ast.UnaryNode:
		switch n.Operator {
		case "+", "-":
			return true
		}
	case *ast.BinaryNode:
		switch n.Operator {
		case "+", "/", "-", "*":
			return true
		}
	}
	return false
}

func setTypeForIntegers(node ast.Node, t reflect.Type) {
	switch n := node.(type) {
	case *ast.IntegerNode:
		n.SetType(t)
	case *ast.UnaryNode:
		switch n.Operator {
		case "+", "-":
			setTypeForIntegers(n.Node, t)
		}
	case *ast.BinaryNode:
		switch n.Operator {
		case "+", "/", "-", "*":
			setTypeForIntegers(n.Left, t)
			setTypeForIntegers(n.Right, t)
		}
	}
}
