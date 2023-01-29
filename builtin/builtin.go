package builtin

import (
	"fmt"
	"reflect"

	"github.com/antonmedv/expr/vm/runtime"
)

var (
	anyType     = reflect.TypeOf(new(interface{})).Elem()
	integerType = reflect.TypeOf(0)
)

type Function struct {
	Name     string
	Func     func(args ...interface{}) (interface{}, error)
	Types    []reflect.Type
	Validate func(args []reflect.Type) (reflect.Type, error)
}

var Builtins = []*Function{
	{
		Name: "len",
		Func: runtime.Len,
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
		Name: "abs",
		Func: func(args ...interface{}) (interface{}, error) {
			x := args[0]
			switch x.(type) {
			case float32:
				if x.(float32) < 0 {
					return -x.(float32), nil
				}
			case float64:
				if x.(float64) < 0 {
					return -x.(float64), nil
				}
			case int:
				if x.(int) < 0 {
					return -x.(int), nil
				}
			case int8:
				if x.(int8) < 0 {
					return -x.(int8), nil
				}
			case int16:
				if x.(int16) < 0 {
					return -x.(int16), nil
				}
			case int32:
				if x.(int32) < 0 {
					return -x.(int32), nil
				}
			case int64:
				if x.(int64) < 0 {
					return -x.(int64), nil
				}
			case uint:
				if x.(uint) < 0 {
					return -x.(uint), nil
				}
			case uint8:
				if x.(uint8) < 0 {
					return -x.(uint8), nil
				}
			case uint16:
				if x.(uint16) < 0 {
					return -x.(uint16), nil
				}
			case uint32:
				if x.(uint32) < 0 {
					return -x.(uint32), nil
				}
			case uint64:
				if x.(uint64) < 0 {
					return -x.(uint64), nil
				}
			}
			return nil, fmt.Errorf("invalid argument for abs (type %T)", x)
		},
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
