package expr

import (
	"fmt"
	"reflect"
)

func isBool(val interface{}) bool {
	return val != nil && reflect.TypeOf(val).Kind() == reflect.Bool
}

func toBool(val interface{}) bool {
	return reflect.ValueOf(val).Bool()
}

func isText(val interface{}) bool {
	return val != nil && reflect.TypeOf(val).Kind() == reflect.String
}

func toText(val interface{}) string {
	return reflect.ValueOf(val).String()
}

func equal(left, right interface{}) bool {
	if isNumber(left) && canBeNumber(right) {
		right, _ := cast(right)
		return left == right
	} else if canBeNumber(left) && isNumber(right) {
		left, _ := cast(left)
		return left == right
	} else {
		return reflect.DeepEqual(left, right)
	}
}

func isNumber(val interface{}) bool {
	return val != nil && reflect.TypeOf(val).Kind() == reflect.Float64
}

func cast(v interface{}) (float64, error) {
	if v != nil {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Float32, reflect.Float64:
			return v.(float64), nil

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(reflect.ValueOf(v).Int()), nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(reflect.ValueOf(v).Uint()), nil // TODO: Check if uint64 fits into float64.
		}
	}
	return 0, fmt.Errorf("can't cast %T to float64", v)
}

func canBeNumber(v interface{}) bool {
	if v != nil {
		return isNumberType(reflect.TypeOf(v))
	}
	return false
}

func extract(from interface{}, it interface{}) (interface{}, error) {
	if from != nil {
		switch reflect.TypeOf(from).Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			i, err := cast(it)
			if err != nil {
				return nil, err
			}

			value := reflect.ValueOf(from).Index(int(i))
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
		case reflect.Map:
			value := reflect.ValueOf(from).MapIndex(reflect.ValueOf(it))
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
		case reflect.Struct:
			value := reflect.ValueOf(from).FieldByName(reflect.ValueOf(it).String())
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
		case reflect.Ptr:
			value := reflect.ValueOf(from).Elem()
			if value.IsValid() && value.CanInterface() {
				return extract(value.Interface(), it)
			}
		}
	}
	return nil, fmt.Errorf("can't get %q from %T", it, from)
}

func contains(needle interface{}, array interface{}) (bool, error) {
	if array != nil {
		value := reflect.ValueOf(array)
		switch reflect.TypeOf(array).Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				value := value.Index(i)
				if value.IsValid() && value.CanInterface() {
					if equal(value.Interface(), needle) {
						return true, nil
					}
				}
			}
			return false, nil
		}
		return false, fmt.Errorf("operator in not defined on %T", array)
	}
	return false, nil
}

func count(node Node, array interface{}) (float64, error) {
	if array != nil {
		value := reflect.ValueOf(array)
		switch reflect.TypeOf(array).Kind() {
		case reflect.Array, reflect.Slice:
			return float64(value.Len()), nil
		case reflect.String:
			return float64(value.Len()), nil
		}
		return 0, fmt.Errorf("invalid argument %v (type %T) for len", node, array)
	}

	return 0, nil
}

func call(name string, fn interface{}, arguments []Node, env interface{}) (interface{}, error) {
	in := make([]reflect.Value, 0)
	for _, arg := range arguments {
		a, err := Run(arg, env)
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(a))
	}

	out := reflect.ValueOf(fn).Call(in)
	if len(out) == 0 {
		return nil, nil
	} else if len(out) > 1 {
		return nil, fmt.Errorf("func %q must return only one value", name)
	}

	if out[0].IsValid() && out[0].CanInterface() {
		return out[0].Interface(), nil
	}
	return nil, nil
}
