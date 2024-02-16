package builtin

import (
	"fmt"
	"reflect"
)

var (
	anyType     = reflect.TypeOf(new(any)).Elem()
	integerType = reflect.TypeOf(0)
	floatType   = reflect.TypeOf(float64(0))
	arrayType   = reflect.TypeOf([]any{})
	mapType     = reflect.TypeOf(map[any]any{})
)

func kind(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}

func types(types ...any) []reflect.Type {
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

func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}

loop:
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = reflect.Indirect(v)
		switch v.Kind() {
		case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
			break loop
		default:
			v = v.Elem()
		}
	}

	if v.IsValid() {
		return v
	}

	panic(fmt.Sprintf("cannot deref %s", v))
}

func toInt(val any) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("cannot use %T as argument (type int)", val)
	}
}

func bitFunc(name string, fn func(x, y int) (any, error)) *Function {
	return &Function{
		Name: name,
		Func: func(args ...any) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("invalid number of arguments for %s (expected 2, got %d)", name, len(args))
			}
			x, err := toInt(args[0])
			if err != nil {
				return nil, fmt.Errorf("%v to call %s", err, name)
			}
			y, err := toInt(args[1])
			if err != nil {
				return nil, fmt.Errorf("%v to call %s", err, name)
			}
			return fn(x, y)
		},
		Types: types(new(func(int, int) int)),
	}
}

func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
