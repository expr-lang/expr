package vm

//go:generate go run ./generate

import (
	"fmt"
	"math"
	"reflect"
)

type Call struct {
	Name string
	Size int
}

type Scope map[string]interface{}

func fetch(from, i interface{}, nilsafe bool) interface{} {
	v := reflect.ValueOf(from)
	kind := v.Kind()

	// Structures can be access through a pointer or through a value, when they
	// are accessed through a pointer we don't want to copy them to a value.
	if kind == reflect.Ptr && reflect.Indirect(v).Kind() == reflect.Struct {
		v = reflect.Indirect(v)
		kind = v.Kind()
	}

	switch kind {

	case reflect.Array, reflect.Slice, reflect.String:
		value := v.Index(toInt(i))
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}

	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(i))
		if value.IsValid() {
			if value.CanInterface() {
				return value.Interface()
			}
		} else {
			elem := reflect.TypeOf(from).Elem()
			return reflect.Zero(elem).Interface()
		}

	case reflect.Struct:
		value := v.FieldByName(reflect.ValueOf(i).String())
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}
	}
	if !nilsafe {
		panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
	}
	return nil
}

func slice(array, from, to interface{}) interface{} {
	v := reflect.ValueOf(array)

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		length := v.Len()
		a, b := toInt(from), toInt(to)

		if b > length {
			b = length
		}
		if a > b {
			a = b
		}

		value := v.Slice(a, b)
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() && value.CanInterface() {
			return slice(value.Interface(), from, to)
		}

	}
	panic(fmt.Sprintf("cannot slice %v", from))
}

func toNpArray(in []interface{}) interface{} {
	if len(in) == 0 {
		out := make([]float64, 0)
		return out
	}
	kind := reflect.ValueOf(in[0]).Kind()
	switch kind {
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		out := make([]float64, len(in))
		for i := range in {
			out[i] = toFloat64(in[i])
		}
		return out
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		out := make([]int64, len(in))
		for i := range in {
			out[i] = toInt64(in[i])
		}
		return out
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uint:
		out := make([]uint64, len(in))
		for i := range in {
			out[i] = toUint64(in[i])
		}
		return out
	default:
		return nil
	}
}

func FetchFn(from interface{}, name string) reflect.Value {
	v := reflect.ValueOf(from)

	// Methods can be defined on any type.
	if v.NumMethod() > 0 {
		method := v.MethodByName(name)
		if method.IsValid() {
			return method
		}
	}

	d := v
	if v.Kind() == reflect.Ptr {
		d = v.Elem()
	}

	switch d.Kind() {
	case reflect.Map:
		value := d.MapIndex(reflect.ValueOf(name))
		if value.IsValid() && value.CanInterface() {
			return value.Elem()
		}
	case reflect.Struct:
		// If struct has not method, maybe it has func field.
		// To access this field we need dereference value.
		value := d.FieldByName(name)
		if value.IsValid() {
			return value
		}
	}
	panic(fmt.Sprintf(`cannot get "%v" from %T`, name, from))
}

func FetchFnNil(from interface{}, name string) reflect.Value {
	if v := reflect.ValueOf(from); !v.IsValid() {
		return v
	}
	return FetchFn(from, name)
}

func FetchNp(name string) reflect.Value {
	switch name {
	case "abs":
		return reflect.ValueOf(abs)
	case "acos":
		return reflect.ValueOf(acos)
	case "acosh":
		return reflect.ValueOf(acosh)
	case "asin":
		return reflect.ValueOf(asin)
	case "asinh":
		return reflect.ValueOf(asinh)
	case "atan":
		return reflect.ValueOf(atan)
	case "atanh":
		return reflect.ValueOf(atanh)
	case "cbrt":
		return reflect.ValueOf(cbrt)
	case "ceil":
		return reflect.ValueOf(ceil)
	case "cos":
		return reflect.ValueOf(cos)
	case "cosh":
		return reflect.ValueOf(cosh)
	case "erf":
		return reflect.ValueOf(erf)
	case "erfc":
		return reflect.ValueOf(erfc)
	case "erfcinv":
		return reflect.ValueOf(erfcinv)
	case "erfinv":
		return reflect.ValueOf(erfinv)
	case "exp":
		return reflect.ValueOf(exp)
	case "exp2":
		return reflect.ValueOf(exp2)
	case "expm1":
		return reflect.ValueOf(expm1)
	case "floor":
		return reflect.ValueOf(floor)
	case "gamma":
		return reflect.ValueOf(gamma)
	case "j0":
		return reflect.ValueOf(j0)
	case "j1":
		return reflect.ValueOf(j1)
	case "log":
		return reflect.ValueOf(log)
	case "log10":
		return reflect.ValueOf(log10)
	case "log1p":
		return reflect.ValueOf(log1p)
	case "log2":
		return reflect.ValueOf(log2)
	case "logb":
		return reflect.ValueOf(logb)
	case "round":
		return reflect.ValueOf(round)
	case "roundtoeven":
		return reflect.ValueOf(roundtoeven)
	case "sin":
		return reflect.ValueOf(sin)
	case "sinh":
		return reflect.ValueOf(sinh)
	case "sqrt":
		return reflect.ValueOf(sqrt)
	case "tan":
		return reflect.ValueOf(tan)
	case "tanh":
		return reflect.ValueOf(tanh)
	case "trunc":
		return reflect.ValueOf(trunc)
	case "y0":
		return reflect.ValueOf(y0)
	case "y1":
		return reflect.ValueOf(y1)
	case "maximum":
		return reflect.ValueOf(max)
	case "minimum":
		return reflect.ValueOf(min)
	case "mod":
		return reflect.ValueOf(mod)
	case "pow":
		return reflect.ValueOf(pow)
	case "remainder":
		return reflect.ValueOf(remainder)
	case "nanmin":
		return reflect.ValueOf(nanmin)
	case "nanmax":
		return reflect.ValueOf(nanmax)
	case "nanmean":
		return reflect.ValueOf(nanmean)
	case "nanstd":
		return reflect.ValueOf(nanstd)
	case "nansum":
		return reflect.ValueOf(nansum)
	case "nanprod":
		return reflect.ValueOf(nanprod)
	}

	panic(fmt.Sprintf(`cannot execute "%v" from builtins`, name))
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
				if equal(value.Interface(), needle).(bool) {
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

func length(a interface{}) int {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len()
	default:
		panic(fmt.Sprintf("invalid argument for len (type %T)", a))
	}
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

func exponent(a, b interface{}) float64 {
	return math.Pow(toFloat64(a), toFloat64(b))
}

func makeRange(min, max int) []int {
	size := max - min + 1
	if size <= 0 {
		return []int{}
	}
	rng := make([]int, size)
	for i := range rng {
		rng[i] = min + i
	}
	return rng
}

func toInt(a interface{}) int {
	switch x := a.(type) {
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

	default:
		panic(fmt.Sprintf("invalid operation: int(%T)", x))
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
		return x

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

func toUint64(a interface{}) uint64 {
	switch x := a.(type) {
	case float32:
		return uint64(x)
	case float64:
		return uint64(x)

	case int:
		return uint64(x)
	case int8:
		return uint64(x)
	case int16:
		return uint64(x)
	case int32:
		return uint64(x)
	case int64:
		return uint64(x)

	case uint:
		return uint64(x)
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x

	default:
		panic(fmt.Sprintf("invalid operation: uint64(%T)", x))
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

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	r := reflect.ValueOf(v)
	switch r.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return r.IsNil()
	default:
		return false
	}
}
