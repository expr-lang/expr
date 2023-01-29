package builtin

import (
	"fmt"
	"reflect"
)

var (
	anyType     = reflect.TypeOf(new(interface{})).Elem()
	integerType = reflect.TypeOf(0)
)

type Function struct {
	Name      string
	Func      func(args ...interface{}) (interface{}, error)
	Types     []reflect.Type
	Validate  func(args []reflect.Type) (reflect.Type, error)
	BuiltinId int
}

const (
	Len = iota + 1
	Abs
)

var Builtins = []*Function{
	{
		Name:      "len",
		BuiltinId: Len,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for len (expected 1, got %d)", len(args))
			}
			switch args[0].Kind() {
			case reflect.Array, reflect.Map, reflect.Slice, reflect.String, reflect.Interface:
				return integerType, nil
			}
			return anyType, fmt.Errorf("invalid argument for len (type %s)", args[0])
		},
	},
	{
		Name:      "abs",
		BuiltinId: Abs,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for abs (expected 1, got %d)", len(args))
			}
			switch args[0].Kind() {
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Interface:
				return args[0], nil
			}
			return anyType, fmt.Errorf("invalid argument for abs (type %s)", args[0])
		},
	},
}
