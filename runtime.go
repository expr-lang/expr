package expr

import (
	"fmt"
	"reflect"
)

func toNumber(n Node, val interface{}) float64 {
	v, ok := cast(val)
	if ok {
		return v
	}
	panic(fmt.Sprintf("cannot convert %v (type %T) to type float64", n, val))
}

func cast(val interface{}) (float64, bool) {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float(), true

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true // TODO: Check if uint64 fits into float64.
	}
	return 0, false
}

func isNumber(val interface{}) bool {
	return val != nil && reflect.TypeOf(val).Kind() == reflect.Float64
}

func canBeNumber(val interface{}) bool {
	if val != nil {
		return isNumberType(reflect.TypeOf(val))
	}
	return false
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

func extract(val interface{}, i interface{}) (interface{}, bool) {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		n, ok := cast(i)
		if !ok {
			break
		}

		value := v.Index(int(n))
		if value.IsValid() && value.CanInterface() {
			return value.Interface(), true
		}
	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(i))
		if value.IsValid() && value.CanInterface() {
			return value.Interface(), true
		}
	case reflect.Struct:
		value := v.FieldByName(reflect.ValueOf(i).String())
		if value.IsValid() && value.CanInterface() {
			return value.Interface(), true
		}
	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() && value.CanInterface() {
			return extract(value.Interface(), i)
		}
	}
	return nil, false
}

func getFunc(val interface{}, name string) (interface{}, bool) {
	v := reflect.ValueOf(val)
	d := v
	if v.Kind() == reflect.Ptr {
		d = v.Elem()
	}

	switch d.Kind() {
	case reflect.Map:
		value := d.MapIndex(reflect.ValueOf(name))
		if value.IsValid() && value.CanInterface() {
			return value.Interface(), true
		}
		// A map may have method too.
		if v.NumMethod() > 0 {
			method := v.MethodByName(name)
			if method.IsValid() && method.CanInterface() {
				return method.Interface(), true
			}
		}
	case reflect.Struct:
		method := v.MethodByName(name)
		if method.IsValid() && method.CanInterface() {
			return method.Interface(), true
		}

		// If struct has not method, maybe it has func field.
		// To access this field we need dereferenced value.
		value := d.FieldByName(name)
		if value.IsValid() && value.CanInterface() {
			return value.Interface(), true
		}
	}
	return nil, false
}

func contains(needle interface{}, array interface{}) (bool, error) {
	if array != nil {
		v := reflect.ValueOf(array)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				value := v.Index(i)
				if value.IsValid() && value.CanInterface() {
					if equal(value.Interface(), needle) {
						return true, nil
					}
				}
			}
			return false, nil
		case reflect.Map:
			n := reflect.ValueOf(needle)
			if !n.IsValid() {
				return false, fmt.Errorf("cannot use %T as index to %T", needle, array)
			}
			value := v.MapIndex(n)
			if value.IsValid() {
				return true, nil
			}
			return false, nil
		case reflect.Struct:
			n := reflect.ValueOf(needle)
			if !n.IsValid() || n.Kind() != reflect.String {
				return false, fmt.Errorf("cannot use %T as field name of %T", needle, array)
			}
			value := v.FieldByName(n.String())
			if value.IsValid() {
				return true, nil
			}
			return false, nil
		case reflect.Ptr:
			value := v.Elem()
			if value.IsValid() && value.CanInterface() {
				return contains(needle, value.Interface())
			}
			return false, nil
		}
		return false, fmt.Errorf("operator \"in\" not defined on %T", array)
	}
	return false, nil
}

func isNil(val interface{}) bool {
	v := reflect.ValueOf(val)
	return !v.IsValid() || v.IsNil()
}
