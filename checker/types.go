package checker

import (
	"reflect"
	"time"

	. "github.com/expr-lang/expr/checker/nature"
	"github.com/expr-lang/expr/conf"
)

var (
	unknown        = Nature{}
	nilNature      = Nature{Type: reflect.TypeOf(Nil{})}
	boolNature     = Nature{Type: reflect.TypeOf(true)}
	integerNature  = Nature{Type: reflect.TypeOf(0)}
	floatNature    = Nature{Type: reflect.TypeOf(float64(0))}
	stringNature   = Nature{Type: reflect.TypeOf("")}
	arrayNature    = Nature{Type: reflect.TypeOf([]any{})}
	mapNature      = Nature{Type: reflect.TypeOf(map[string]any{})}
	timeNature     = Nature{Type: reflect.TypeOf(time.Time{})}
	durationNature = Nature{Type: reflect.TypeOf(time.Duration(0))}
)

var (
	anyType      = reflect.TypeOf(new(any)).Elem()
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Duration(0))
	arrayType    = reflect.TypeOf([]any{})
)

// Nil is a special type to represent nil.
type Nil struct{}

func isNil(nt Nature) bool {
	if nt.Type == nil {
		return false
	}
	return nt.Type == nilNature.Type
}

func combined(l, r Nature) Nature {
	if isUnknown(l) || isUnknown(r) {
		return unknown
	}
	if isFloat(l) || isFloat(r) {
		return floatNature
	}
	return integerNature
}

func anyOf(nt Nature, fns ...func(Nature) bool) bool {
	for _, fn := range fns {
		if fn(nt) {
			return true
		}
	}
	return false
}

func or(l, r Nature, fns ...func(Nature) bool) bool {
	if isUnknown(l) && isUnknown(r) {
		return true
	}
	if isUnknown(l) && anyOf(r, fns...) {
		return true
	}
	if isUnknown(r) && anyOf(l, fns...) {
		return true
	}
	return false
}

func isUnknown(nt Nature) bool {
	switch {
	case nt.Type == nil:
		return true
	case nt.Kind() == reflect.Interface:
		return true
	}
	return false
}

func isInteger(nt Nature) bool {
	switch nt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

func isFloat(nt Nature) bool {
	switch nt.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func isNumber(nt Nature) bool {
	return isInteger(nt) || isFloat(nt)
}

func isTime(nt Nature) bool {
	switch nt.Type {
	case timeType:
		return true
	}
	return false
}

func isDuration(nt Nature) bool {
	switch nt.Type {
	case durationType:
		return true
	}
	return false
}

func isBool(nt Nature) bool {
	switch nt.Kind() {
	case reflect.Bool:
		return true
	}
	return false
}

func isString(nt Nature) bool {
	switch nt.Kind() {
	case reflect.String:
		return true
	}
	return false
}

func isArray(nt Nature) bool {
	switch nt.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

func isMap(nt Nature) bool {
	switch nt.Kind() {
	case reflect.Map:
		return true
	}
	return false
}

func isStruct(nt Nature) bool {
	switch nt.Kind() {
	case reflect.Struct:
		return true
	}
	return false
}

func isFunc(nt Nature) bool {
	switch nt.Kind() {
	case reflect.Func:
		return true
	}
	return false
}

func fetchField(nt Nature, name string) (reflect.StructField, bool) {
	// First check all structs fields.
	for i := 0; i < nt.NumField(); i++ {
		field := nt.Field(i)
		// Search all fields, even embedded structs.
		if conf.FieldName(field) == name {
			return field, true
		}
	}

	// Second check fields of embedded structs.
	for i := 0; i < nt.NumField(); i++ {
		anon := nt.Field(i)
		if anon.Anonymous {
			anonType := anon.Type
			if anonType.Kind() == reflect.Pointer {
				anonType = anonType.Elem()
			}
			if field, ok := fetchField(Nature{Type: anonType}, name); ok {
				field.Index = append(anon.Index, field.Index...)
				return field, true
			}
		}
	}

	return reflect.StructField{}, false
}

func kind(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}

func isComparable(l, r Nature) bool {
	switch {
	case l.Kind() == r.Kind():
		return true
	case isNumber(l) && isNumber(r):
		return true
	case isNil(l) || isNil(r):
		return true
	case isUnknown(l) || isUnknown(r):
		return true
	}
	return false
}
