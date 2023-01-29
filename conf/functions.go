package conf

import "reflect"

type Function struct {
	Name  string
	Func  func(params ...interface{}) (interface{}, error)
	Types []reflect.Type
}
