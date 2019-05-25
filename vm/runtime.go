package vm

import (
	"fmt"
	"math"
	"reflect"
)

func fetch(from interface{}, i interface{}) interface{} {
	v := reflect.ValueOf(from)
	switch v.Kind() {

	case reflect.Array, reflect.Slice, reflect.String:
		index := int(cast(i))
		value := v.Index(index)
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

func in(needle interface{}, array interface{}) bool {
	if array == nil {
		return false
	}
	v := reflect.ValueOf(array)

	switch v.Kind() {

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			if value.IsValid() && value.CanInterface() {
				if equal(value.Interface(), needle) {
					return true
				}
			}
		}
		return false

	case reflect.Map:
		n := reflect.ValueOf(needle)
		if !n.IsValid() {
			panic(fmt.Sprintf("cannot use %T as index to %T", needle, array))
		}
		value := v.MapIndex(n)
		if value.IsValid() {
			return true
		}
		return false

	case reflect.Struct:
		n := reflect.ValueOf(needle)
		if !n.IsValid() || n.Kind() != reflect.String {
			panic(fmt.Sprintf("cannot use %T as field name of %T", needle, array))
		}
		value := v.FieldByName(n.String())
		if value.IsValid() {
			return true
		}
		return false

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() && value.CanInterface() {
			return in(needle, value.Interface())
		}
		return false
	}

	panic(fmt.Sprintf(`operator "in"" not defined on %T`, array))
}

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
		return x == int(cast(b))
	case int8:
		return x == int8(cast(b))
	case int16:
		return x == int16(cast(b))
	case int32:
		return x == int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) == float64(y)
		default:
			return x == cast(b)
		}

	case uint:
		return x == uint(cast(b))
	case uint8:
		return x == uint8(cast(b))
	case uint16:
		return x == uint16(cast(b))
	case uint32:
		return x == uint32(cast(b))
	case uint64:
		return x == uint64(cast(b))

	default:
		return reflect.DeepEqual(a, b)
	}
}

func less(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x < float32(y)
		default:
			return x < b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x < float64(y)
		default:
			return x < b.(float64)
		}

	case int:
		return x < int(cast(b))
	case int8:
		return x < int8(cast(b))
	case int16:
		return x < int16(cast(b))
	case int32:
		return x < int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) < float64(y)
		default:
			return x < cast(b)
		}

	case uint:
		return x < uint(cast(b))
	case uint8:
		return x < uint8(cast(b))
	case uint16:
		return x < uint16(cast(b))
	case uint32:
		return x < uint32(cast(b))
	case uint64:
		return x < uint64(cast(b))

	default:
		panic(fmt.Sprintf("invalid operation: %T < %T", a, b))
	}
}

func more(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x > float32(y)
		default:
			return x > b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x > float64(y)
		default:
			return x > b.(float64)
		}

	case int:
		return x > int(cast(b))
	case int8:
		return x > int8(cast(b))
	case int16:
		return x > int16(cast(b))
	case int32:
		return x > int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) > float64(y)
		default:
			return x > cast(b)
		}

	case uint:
		return x > uint(cast(b))
	case uint8:
		return x > uint8(cast(b))
	case uint16:
		return x > uint16(cast(b))
	case uint32:
		return x > uint32(cast(b))
	case uint64:
		return x > uint64(cast(b))

	default:
		panic(fmt.Sprintf("invalid operation: %T > %T", a, b))
	}
}

func lessOrEqual(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x <= float32(y)
		default:
			return x <= b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x <= float64(y)
		default:
			return x <= b.(float64)
		}

	case int:
		return x <= int(cast(b))
	case int8:
		return x <= int8(cast(b))
	case int16:
		return x <= int16(cast(b))
	case int32:
		return x <= int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) <= float64(y)
		default:
			return x <= cast(b)
		}

	case uint:
		return x <= uint(cast(b))
	case uint8:
		return x <= uint8(cast(b))
	case uint16:
		return x <= uint16(cast(b))
	case uint32:
		return x <= uint32(cast(b))
	case uint64:
		return x <= uint64(cast(b))

	default:
		panic(fmt.Sprintf("invalid operation: %T <= %T", a, b))
	}
}

func moreOrEqual(a, b interface{}) interface{} {
	switch x := a.(type) {
	case float32:
		switch y := b.(type) {
		case int64:
			return x >= float32(y)
		default:
			return x >= b.(float32)
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x >= float64(y)
		default:
			return x >= b.(float64)
		}

	case int:
		return x >= int(cast(b))
	case int8:
		return x >= int8(cast(b))
	case int16:
		return x >= int16(cast(b))
	case int32:
		return x >= int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) >= float64(y)
		default:
			return x >= cast(b)
		}

	case uint:
		return x >= uint(cast(b))
	case uint8:
		return x >= uint8(cast(b))
	case uint16:
		return x >= uint16(cast(b))
	case uint32:
		return x >= uint32(cast(b))
	case uint64:
		return x >= uint64(cast(b))

	default:
		panic(fmt.Sprintf("invalid operation: %T >= %T", a, b))
	}
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
		return x + int(cast(b))
	case int8:
		return x + int8(cast(b))
	case int16:
		return x + int16(cast(b))
	case int32:
		return x + int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) + float64(y)
		default:
			return x + cast(b)
		}

	case uint:
		return x + uint(cast(b))
	case uint8:
		return x + uint8(cast(b))
	case uint16:
		return x + uint16(cast(b))
	case uint32:
		return x + uint32(cast(b))
	case uint64:
		return x + uint64(cast(b))

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
		return x - int(cast(b))
	case int8:
		return x - int8(cast(b))
	case int16:
		return x - int16(cast(b))
	case int32:
		return x - int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) - float64(y)
		default:
			return x - cast(b)
		}

	case uint:
		return x - uint(cast(b))
	case uint8:
		return x - uint8(cast(b))
	case uint16:
		return x - uint16(cast(b))
	case uint32:
		return x - uint32(cast(b))
	case uint64:
		return x - uint64(cast(b))

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
		return x * int(cast(b))
	case int8:
		return x * int8(cast(b))
	case int16:
		return x * int16(cast(b))
	case int32:
		return x * int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) * float64(y)
		default:
			return x * cast(b)
		}

	case uint:
		return x * uint(cast(b))
	case uint8:
		return x * uint8(cast(b))
	case uint16:
		return x * uint16(cast(b))
	case uint32:
		return x * uint32(cast(b))
	case uint64:
		return x * uint64(cast(b))

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
		return x / int(cast(b))
	case int8:
		return x / int8(cast(b))
	case int16:
		return x / int16(cast(b))
	case int32:
		return x / int32(cast(b))
	case int64:
		switch y := b.(type) {
		case float64:
			return float64(x) / float64(y)
		default:
			return x / cast(b)
		}

	case uint:
		return x / uint(cast(b))
	case uint8:
		return x / uint8(cast(b))
	case uint16:
		return x / uint16(cast(b))
	case uint32:
		return x / uint32(cast(b))
	case uint64:
		return x / uint64(cast(b))

	default:
		panic(fmt.Sprintf("invalid operation: %T / %T", a, b))
	}
}

func modulo(a, b interface{}) interface{} {
	switch x := a.(type) {
	case int:
		return x % int(cast(b))
	case int8:
		return x % int8(cast(b))
	case int16:
		return x % int16(cast(b))
	case int32:
		return x % int32(cast(b))
	case int64:
		return x % int64(cast(b))

	case uint:
		return x % uint(cast(b))
	case uint8:
		return x % uint8(cast(b))
	case uint16:
		return x % uint16(cast(b))
	case uint32:
		return x % uint32(cast(b))
	case uint64:
		return x % uint64(cast(b))

	default:
		panic(fmt.Sprintf("invalid operation: %T %v %T", a, "%", b))
	}
}

func exponent(a, b interface{}) float64 {
	return math.Pow(toFloat64(a), toFloat64(b))
}

func makeRange(a, b interface{}) []int64 {
	min := toInt64(a)
	max := toInt64(b)
	size := max - min + 1
	rng := make([]int64, size)
	for i := range rng {
		rng[i] = min + int64(i)
	}
	return rng
}

func cast(a interface{}) int64 {
	switch x := a.(type) {
	case int:
		return int64(x)
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return int64(x)

	case uint:
		return int64(x)
	case uint8:
		return int64(x)
	case uint16:
		return int64(x)
	case uint32:
		return int64(x)
	case uint64:
		return int64(x)

	default:
		panic(fmt.Sprintf("can't cast %T to int64", a))
	}
}

func toInt64(a interface{}) int64 {
	switch x := a.(type) {
	case float32:
		return int64(x)
	case float64:
		return int64(x)

	case int:
		return int64(x)
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return int64(x)

	case uint:
		return int64(x)
	case uint8:
		return int64(x)
	case uint16:
		return int64(x)
	case uint32:
		return int64(x)
	case uint64:
		return int64(x)

	default:
		panic(fmt.Sprintf("invalid operation: int64(%T)", x))
	}
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
