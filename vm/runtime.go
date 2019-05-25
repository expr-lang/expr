package vm

import (
	"fmt"
	"math"
	"reflect"
)

func negate(i interface{}) interface{} {
	switch v := i.(type) {
	case float32:
		return -v
	case float64:
		return -v

	case int:
		return -v
	case int8:
		return -v
	case int16:
		return -v
	case int32:
		return -v
	case int64:
		return -v

	case uint:
		return -v
	case uint8:
		return -v
	case uint16:
		return -v
	case uint32:
		return -v
	case uint64:
		return -v

	default:
		panic(fmt.Sprintf("invalid operation: - %T", v))
	}
}

func equal(a, b interface{}) bool {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x == float32(y)
		default:
			return x == b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x == float64(y)
		default:
			return x == b.(float64)
		}

	case int:
		return x == b.(int)
	case int8:
		return x == b.(int8)
	case int16:
		return x == b.(int16)
	case int32:
		return x == b.(int32)
	case int64:
		return x == b.(int64)

	case uint:
		return x == b.(uint)
	case uint8:
		return x == b.(uint8)
	case uint16:
		return x == b.(uint16)
	case uint32:
		return x == b.(uint32)
	case uint64:
		return x == b.(uint64)

	default:
		return reflect.DeepEqual(a, b)
	}
}

func fetch(from interface{}, i interface{}) interface{} {
	v := reflect.ValueOf(from)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		value := v.Index(i.(int))
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}
	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(i))
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}
	case reflect.Struct:
		value := v.FieldByName(reflect.ValueOf(i).String())
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}
	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() && value.CanInterface() {
			return fetch(value.Interface(), i)
		}
	}
	panic(fmt.Sprintf("%v doesn't contains %v", from, i))
}

func add(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x + float32(y)
		default:
			return x + b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x + float64(y)
		default:
			return x + b.(float64)
		}

	case int:
		return x + b.(int)
	case int8:
		return x + b.(int8)
	case int16:
		return x + b.(int16)
	case int32:
		return x + b.(int32)
	case int64:
		return x + b.(int64)

	case uint:
		return x + b.(uint)
	case uint8:
		return x + b.(uint8)
	case uint16:
		return x + b.(uint16)
	case uint32:
		return x + b.(uint32)
	case uint64:
		return x + b.(uint64)

	case string:
		return x + b.(string)

	default:
		panic(fmt.Sprintf("invalid operation: %T + %T", a, b))
	}
}

func subtract(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x - float32(y)
		default:
			return x - b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x - float64(y)
		default:
			return x - b.(float64)
		}

	case int:
		return x - b.(int)
	case int8:
		return x - b.(int8)
	case int16:
		return x - b.(int16)
	case int32:
		return x - b.(int32)
	case int64:
		return x - b.(int64)

	case uint:
		return x - b.(uint)
	case uint8:
		return x - b.(uint8)
	case uint16:
		return x - b.(uint16)
	case uint32:
		return x - b.(uint32)
	case uint64:
		return x - b.(uint64)

	default:
		panic(fmt.Sprintf("invalid operation: %T - %T", a, b))
	}
}

func multiply(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x * float32(y)
		default:
			return x * b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x * float64(y)
		default:
			return x * b.(float64)
		}

	case int:
		return x * b.(int)
	case int8:
		return x * b.(int8)
	case int16:
		return x * b.(int16)
	case int32:
		return x * b.(int32)
	case int64:
		return x * b.(int64)

	case uint:
		return x * b.(uint)
	case uint8:
		return x * b.(uint8)
	case uint16:
		return x * b.(uint16)
	case uint32:
		return x * b.(uint32)
	case uint64:
		return x * b.(uint64)

	default:
		panic(fmt.Sprintf("invalid operation: %T * %T", a, b))
	}
}

func divide(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x / float32(y)
		default:
			return x / b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x / float64(y)
		default:
			return x / b.(float64)
		}

	case int:
		return x / b.(int)
	case int8:
		return x / b.(int8)
	case int16:
		return x / b.(int16)
	case int32:
		return x / b.(int32)
	case int64:
		return x / b.(int64)

	case uint:
		return x / b.(uint)
	case uint8:
		return x / b.(uint8)
	case uint16:
		return x / b.(uint16)
	case uint32:
		return x / b.(uint32)
	case uint64:
		return x / b.(uint64)

	default:
		panic(fmt.Sprintf("invalid operation: %T / %T", a, b))
	}
}

func modulo(a, b interface{}) interface{} {
	switch x := a.(type) {
	case int:
		return x % b.(int)
	case int8:
		return x % b.(int8)
	case int16:
		return x % b.(int16)
	case int32:
		return x % b.(int32)
	case int64:
		return x % b.(int64)

	case uint:
		return x % b.(uint)
	case uint8:
		return x % b.(uint8)
	case uint16:
		return x % b.(uint16)
	case uint32:
		return x % b.(uint32)
	case uint64:
		return x % b.(uint64)

	default:
		panic(fmt.Sprintf("invalid operation: %T %v %T", a, "%", b))
	}
}

func exponent(a, b interface{}) float64 {
	return math.Pow(toFloat64(a), toFloat64(b))
}

func toFloat64(a interface{}) float64 {
	switch x := a.(type) {
	case float32:
		return float64(x)
	case float64:
		return x

	case int:
		return float64(x)
	case int8:
		return float64(x)
	case int16:
		return float64(x)
	case int32:
		return float64(x)
	case int64:
		return float64(x)

	case uint:
		return float64(x)
	case uint8:
		return float64(x)
	case uint16:
		return float64(x)
	case uint32:
		return float64(x)
	case uint64:
		return float64(x)

	default:
		panic(fmt.Sprintf("invalid operation: float64(%T)", x))
	}
}