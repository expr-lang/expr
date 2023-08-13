package builtin

import "reflect"

var (
	anyType     = reflect.TypeOf(new(interface{})).Elem()
	integerType = reflect.TypeOf(0)
	floatType   = reflect.TypeOf(float64(0))
)

func kind(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}
