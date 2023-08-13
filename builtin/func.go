package builtin

import (
	"fmt"
	"reflect"
	"strconv"
)

func Len(xs ...interface{}) (interface{}, error) {
	x := xs[0]
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len(), nil
	default:
		return nil, fmt.Errorf("invalid argument for len (type %T)", x)
	}
}

func Abs(xs ...interface{}) (interface{}, error) {
	x := xs[0]
	switch x.(type) {
	case float32:
		if x.(float32) < 0 {
			return -x.(float32), nil
		} else {
			return x, nil
		}
	case float64:
		if x.(float64) < 0 {
			return -x.(float64), nil
		} else {
			return x, nil
		}
	case int:
		if x.(int) < 0 {
			return -x.(int), nil
		} else {
			return x, nil
		}
	case int8:
		if x.(int8) < 0 {
			return -x.(int8), nil
		} else {
			return x, nil
		}
	case int16:
		if x.(int16) < 0 {
			return -x.(int16), nil
		} else {
			return x, nil
		}
	case int32:
		if x.(int32) < 0 {
			return -x.(int32), nil
		} else {
			return x, nil
		}
	case int64:
		if x.(int64) < 0 {
			return -x.(int64), nil
		} else {
			return x, nil
		}
	case uint:
		if x.(uint) < 0 {
			return -x.(uint), nil
		} else {
			return x, nil
		}
	case uint8:
		if x.(uint8) < 0 {
			return -x.(uint8), nil
		} else {
			return x, nil
		}
	case uint16:
		if x.(uint16) < 0 {
			return -x.(uint16), nil
		} else {
			return x, nil
		}
	case uint32:
		if x.(uint32) < 0 {
			return -x.(uint32), nil
		} else {
			return x, nil
		}
	case uint64:
		if x.(uint64) < 0 {
			return -x.(uint64), nil
		} else {
			return x, nil
		}
	}
	return nil, fmt.Errorf("invalid argument for abs (type %T)", x)
}

func Int(xs ...interface{}) (interface{}, error) {
	x := xs[0]
	switch x := x.(type) {
	case float32:
		return int(x), nil
	case float64:
		return int(x), nil
	case int:
		return x, nil
	case int8:
		return int(x), nil
	case int16:
		return int(x), nil
	case int32:
		return int(x), nil
	case int64:
		return int(x), nil
	case uint:
		return int(x), nil
	case uint8:
		return int(x), nil
	case uint16:
		return int(x), nil
	case uint32:
		return int(x), nil
	case uint64:
		return int(x), nil
	case string:
		i, err := strconv.Atoi(x)
		if err != nil {
			return nil, fmt.Errorf("invalid operation: int(%s)", x)
		}
		return i, nil
	default:
		return nil, fmt.Errorf("invalid operation: int(%T)", x)
	}
}

func Float(xs ...interface{}) (interface{}, error) {
	x := xs[0]
	switch x := x.(type) {
	case float32:
		return float64(x), nil
	case float64:
		return x, nil
	case int:
		return float64(x), nil
	case int8:
		return float64(x), nil
	case int16:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case uint8:
		return float64(x), nil
	case uint16:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid operation: float(%s)", x)
		}
		return f, nil
	default:
		return nil, fmt.Errorf("invalid operation: float(%T)", x)
	}
}
