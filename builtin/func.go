package builtin

import (
	"fmt"
	"reflect"
	"strconv"
)

func Len(x interface{}) interface{} {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len()
	default:
		panic(fmt.Sprintf("invalid argument for len (type %T)", x))
	}
}

func Abs(x interface{}) interface{} {
	switch x.(type) {
	case float32:
		if x.(float32) < 0 {
			return -x.(float32)
		} else {
			return x
		}
	case float64:
		if x.(float64) < 0 {
			return -x.(float64)
		} else {
			return x
		}
	case int:
		if x.(int) < 0 {
			return -x.(int)
		} else {
			return x
		}
	case int8:
		if x.(int8) < 0 {
			return -x.(int8)
		} else {
			return x
		}
	case int16:
		if x.(int16) < 0 {
			return -x.(int16)
		} else {
			return x
		}
	case int32:
		if x.(int32) < 0 {
			return -x.(int32)
		} else {
			return x
		}
	case int64:
		if x.(int64) < 0 {
			return -x.(int64)
		} else {
			return x
		}
	case uint:
		if x.(uint) < 0 {
			return -x.(uint)
		} else {
			return x
		}
	case uint8:
		if x.(uint8) < 0 {
			return -x.(uint8)
		} else {
			return x
		}
	case uint16:
		if x.(uint16) < 0 {
			return -x.(uint16)
		} else {
			return x
		}
	case uint32:
		if x.(uint32) < 0 {
			return -x.(uint32)
		} else {
			return x
		}
	case uint64:
		if x.(uint64) < 0 {
			return -x.(uint64)
		} else {
			return x
		}
	}
	panic(fmt.Sprintf("invalid argument for abs (type %T)", x))
}

func Int(x interface{}) interface{} {
	switch x := x.(type) {
	case float32:
		return int(x)
	case float64:
		return int(x)
	case int:
		return x
	case int8:
		return int(x)
	case int16:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)
	case uint:
		return int(x)
	case uint8:
		return int(x)
	case uint16:
		return int(x)
	case uint32:
		return int(x)
	case uint64:
		return int(x)
	case string:
		i, err := strconv.Atoi(x)
		if err != nil {
			panic(fmt.Sprintf("invalid operation: int(%s)", x))
		}
		return i
	default:
		panic(fmt.Sprintf("invalid operation: int(%T)", x))
	}
}

func Float(x interface{}) interface{} {
	switch x := x.(type) {
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
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid operation: float(%s)", x))
		}
		return f
	default:
		panic(fmt.Sprintf("invalid operation: float(%T)", x))
	}
}

func String(arg interface{}) interface{} {
	return fmt.Sprintf("%v", arg)
}
