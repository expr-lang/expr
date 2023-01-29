package builtin

import (
	"fmt"
	"reflect"

	"github.com/antonmedv/expr/vm/runtime"
)

type Function struct {
	Name     string
	Func     func(params ...interface{}) (interface{}, error)
	Types    []reflect.Type
	Validate func(args []reflect.Type) error
}

var Builtins = []*Function{
	{
		Name: "len",
		Func: runtime.Len,
		Validate: func(args []reflect.Type) error {
			switch args[0].Kind() {
			case reflect.Array, reflect.Map, reflect.Slice, reflect.String, reflect.Interface:
				return nil
			}
			return fmt.Errorf("invalid argument for len (type %s)", args[0])
		},
	},
}
