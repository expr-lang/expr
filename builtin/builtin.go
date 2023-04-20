package builtin

import (
	"fmt"
	"reflect"
)

var (
	anyType     = reflect.TypeOf(new(interface{})).Elem()
	integerType = reflect.TypeOf(0)
	floatType   = reflect.TypeOf(float64(0))
	stringType  = reflect.TypeOf("")
)

var (
	stringKindMap = map[reflect.Kind]struct{}{reflect.String: {}}
	numberKindMap = map[reflect.Kind]struct{}{
		reflect.Float32: {},
		reflect.Float64: {},
		reflect.Int:     {},
		reflect.Int8:    {},
		reflect.Int16:   {},
		reflect.Int32:   {},
		reflect.Int64:   {},
		reflect.Uint:    {},
		reflect.Uint8:   {},
		reflect.Uint16:  {},
		reflect.Uint32:  {},
		reflect.Uint64:  {},
	}
)

type Function struct {
	Name     string
	Func     func(args ...interface{}) (interface{}, error)
	Opcode   int
	Types    []reflect.Type
	Validate func(args []reflect.Type) (reflect.Type, error)
}

const (
	Len = iota + 1
	Abs
	Int
	Float
	Upper
	Lower
	Left
	Right
	LPad
	RPad
	Substring
	Reverse
)

var Builtins = map[int]*Function{
	Len: {
		Name:   "len",
		Opcode: Len,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for len (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Array, reflect.Map, reflect.Slice, reflect.String, reflect.Interface:
				return integerType, nil
			}
			return anyType, fmt.Errorf("invalid argument for len (type %s)", args[0])
		},
	},
	Abs: {
		Name:   "abs",
		Opcode: Abs,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for abs (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Interface:
				return args[0], nil
			}
			return anyType, fmt.Errorf("invalid argument for abs (type %s)", args[0])
		},
	},
	Int: {
		Name:   "int",
		Opcode: Int,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for int (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return integerType, nil
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return integerType, nil
			case reflect.String:
				return integerType, nil
			}
			return anyType, fmt.Errorf("invalid argument for int (type %s)", args[0])
		},
	},
	Float: {
		Name:   "float",
		Opcode: Float,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for float (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return floatType, nil
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return floatType, nil
			case reflect.String:
				return floatType, nil
			}
			return anyType, fmt.Errorf("invalid argument for float (type %s)", args[0])
		},
	},
	Upper: {
		Name:   "upper",
		Opcode: Upper,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for upper (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.String:
				return stringType, nil
			}
			return anyType, fmt.Errorf("invalid argument for upper (type %s)", args[0])
		},
	},
	Lower: {
		Name:   "lower",
		Opcode: Lower,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for lower (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.String:
				return stringType, nil
			default:
				return anyType, fmt.Errorf("invalid argument no. 1 for lower (type %s)", args[0])
			}
		},
	},
	Left: {
		Name:   "left",
		Opcode: Left,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if err := validateArgs("left", args, 2, []map[reflect.Kind]struct{}{
				stringKindMap, numberKindMap,
			}); err != nil {
				return anyType, err
			}
			return stringType, nil
		},
	},
	Right: {
		Name:   "right",
		Opcode: Right,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if err := validateArgs("right", args, 2, []map[reflect.Kind]struct{}{
				stringKindMap, numberKindMap,
			}); err != nil {
				return anyType, err
			}
			return stringType, nil
		},
	},
	LPad: {
		Name:   "lpad",
		Opcode: LPad,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if err := validateArgs("lpad", args, 3, []map[reflect.Kind]struct{}{
				stringKindMap, stringKindMap, numberKindMap,
			}); err != nil {
				return anyType, err
			}
			return stringType, nil
		},
	},
	RPad: {
		Name:   "rpad",
		Opcode: RPad,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if err := validateArgs("rpad", args, 3, []map[reflect.Kind]struct{}{
				stringKindMap, stringKindMap, numberKindMap,
			}); err != nil {
				return anyType, err
			}
			return stringType, nil
		},
	},
	Substring: {
		Name:   "substr",
		Opcode: Substring,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if err := validateArgs("substr", args, 3, []map[reflect.Kind]struct{}{
				stringKindMap, numberKindMap, numberKindMap,
			}); err != nil {
				return anyType, err
			}
			return stringType, nil
		},
	},
	Reverse: {
		Name:   "reverse",
		Opcode: Reverse,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments for reverse (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.String:
				return stringType, nil
			}
			return anyType, fmt.Errorf("invalid argument for reverse (type %s)", args[0])
		},
	},
}

func validateArgs(funcName string, args []reflect.Type, noArgs int, kindsForArgs []map[reflect.Kind]struct{}) error {
	if len(args) != noArgs {
		return fmt.Errorf("invalid number of arguments for %s (expected %d, got %d)", funcName, noArgs, len(args))
	}

	for argNo, argType := range args {
		if err := validateArg(funcName, argNo+1, argType, kindsForArgs[argNo]); err != nil {
			return err
		}
	}
	return nil
}

func validateArg(funcName string, argNo int, argType reflect.Type, kinds map[reflect.Kind]struct{}) error {
	if _, ok := kinds[kind(argType)]; !ok {
		return fmt.Errorf("invalid argument no. %d for %s (type %s)", argNo, funcName, argType)
	}
	return nil
}

func kind(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}
