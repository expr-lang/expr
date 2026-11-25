package builtin

import (
	"fmt"
	"reflect"
	"time"

	"github.com/expr-lang/expr/internal/deref"
	"github.com/expr-lang/expr/vm/runtime"
)

var (
	anyType      = reflect.TypeOf(new(any)).Elem()
	integerType  = reflect.TypeOf(0)
	floatType    = reflect.TypeOf(float64(0))
	arrayType    = reflect.TypeOf([]any{})
	mapType      = reflect.TypeOf(map[any]any{})
	timeType     = reflect.TypeOf(new(time.Time)).Elem()
	locationType = reflect.TypeOf(new(time.Location))
)

func kind(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	t = deref.Type(t)
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

func containsEqual(values []any, value any) bool {
	for _, v := range values {
		if runtime.Equal(value, v) {
			return true
		}
	}
	return false
}

// uniqIndex finds values equal under runtime.Equal without a linear scan for
// strings and numbers. Strings only equal strings. Integers equal integers
// via int(x) and floats via float64(x), so integers are also keyed as
// float64 for float lookups. Everything else falls back to a scan.
type uniqIndex struct {
	strings   map[string]struct{}
	ints      map[int]struct{}
	intFloats map[float64]struct{}
	floats    map[float64]struct{}
	other     []any
}

func newUniqIndex(values []any) *uniqIndex {
	u := &uniqIndex{
		strings:   map[string]struct{}{},
		ints:      map[int]struct{}{},
		intFloats: map[float64]struct{}{},
		floats:    map[float64]struct{}{},
	}
	for _, v := range values {
		u.add(v)
	}
	return u
}

func (u *uniqIndex) has(value any) bool {
	if s, ok := value.(string); ok {
		_, found := u.strings[s]
		return found
	}
	if i, f, ok := intKeys(value); ok {
		_, found := u.ints[i]
		if !found {
			_, found = u.floats[f]
		}
		return found
	}
	if f, ok := floatKey(value); ok {
		_, found := u.floats[f]
		if !found {
			_, found = u.intFloats[f]
		}
		return found
	}
	return containsEqual(u.other, value)
}

func (u *uniqIndex) add(value any) {
	if s, ok := value.(string); ok {
		u.strings[s] = struct{}{}
	} else if i, f, ok := intKeys(value); ok {
		u.ints[i] = struct{}{}
		u.intFloats[f] = struct{}{}
	} else if f, ok := floatKey(value); ok {
		u.floats[f] = struct{}{}
	} else {
		u.other = append(u.other, value)
	}
}

// intKeys mirrors the conversions runtime.Equal applies to integers.
func intKeys(value any) (int, float64, bool) {
	switch x := value.(type) {
	case int:
		return int(x), float64(x), true
	case int8:
		return int(x), float64(x), true
	case int16:
		return int(x), float64(x), true
	case int32:
		return int(x), float64(x), true
	case int64:
		return int(x), float64(x), true
	case uint:
		return int(x), float64(x), true
	case uint8:
		return int(x), float64(x), true
	case uint16:
		return int(x), float64(x), true
	case uint32:
		return int(x), float64(x), true
	case uint64:
		return int(x), float64(x), true
	}
	return 0, 0, false
}

func floatKey(value any) (float64, bool) {
	switch x := value.(type) {
	case float32:
		return float64(x), true
	case float64:
		return x, true
	}
	return 0, false
}
