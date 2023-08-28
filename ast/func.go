package ast

import "reflect"

type Function struct {
	Name      string
	Func      func(args ...interface{}) (interface{}, error)
	Fast      func(arg interface{}) interface{}
	Types     []reflect.Type
	Validate  func(args []reflect.Type) (reflect.Type, error)
	Predicate bool
}
