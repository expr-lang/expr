package checker

import "reflect"

var (
	nilType       = reflect.TypeOf(nil)
	boolType      = reflect.TypeOf(true)
	numberType    = reflect.TypeOf(float64(0))
	textType      = reflect.TypeOf("")
	arrayType     = reflect.TypeOf([]interface{}{})
	mapType       = reflect.TypeOf(map[interface{}]interface{}{})
	interfaceType = reflect.TypeOf(new(interface{})).Elem()
)
