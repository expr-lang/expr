package builtin

import (
	"reflect"
)

var (
	anyType     = reflect.TypeOf(new(interface{})).Elem()
	integerType = reflect.TypeOf(0)
	floatType   = reflect.TypeOf(float64(0))
	stringType  = reflect.TypeOf("")
	arrayType   = reflect.TypeOf([]interface{}{})
)

func kind(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}

func types(types ...interface{}) []reflect.Type {
	ts := make([]reflect.Type, len(types))
	for i, t := range types {
		t := reflect.TypeOf(t)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Func {
			panic("not a function")
		}
		ts[i] = t
	}
	return ts
}
