package runtime

//go:generate sh -c "go run ./helpers > ./helpers[generated].go"

import (
	"fmt"
	"math"
	"reflect"
	"sync"

	"github.com/expr-lang/expr/internal/deref"
)

// structFields caches, for each struct type, the index paths of its fields
// by their expr names (map[reflect.Type]map[string][]int).
var structFields sync.Map

func Fetch(from, i any) any {
	v := reflect.ValueOf(from)
	if v.Kind() == reflect.Invalid {
		panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
	}

	// Methods can be defined on any type.
	if v.NumMethod() > 0 {
		if methodName, ok := i.(string); ok {
			if m, ok := MethodByName(v, methodName); ok {
				return m
			}
		}
	}

	// Structs, maps, and slices can be access through a pointer or through
	// a value, when they are accessed through a pointer we don't want to
	// copy them to a value.
	// De-reference everything if necessary (interface and pointers)
	v = deref.Value(v)

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		if _, ok := i.(string); ok {
			panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
		}
		index := ToInt(i)
		l := v.Len()
		if index < 0 {
			index = l + index
		}
		if index < 0 || index >= l {
			panic(fmt.Sprintf("index out of range: %v (array length is %v)", index, l))
		}
		value := v.Index(index)
		if value.IsValid() {
			return value.Interface()
		}

	case reflect.Map:
		var value reflect.Value
		if i == nil {
			value = v.MapIndex(reflect.Zero(v.Type().Key()))
		} else {
			value = v.MapIndex(reflect.ValueOf(i))
		}
		if value.IsValid() {
			return value.Interface()
		} else {
			elem := reflect.TypeOf(from).Elem()
			return reflect.Zero(elem).Interface()
		}

	case reflect.Struct:
		if fieldName, ok := i.(string); ok {
			if value, ok := StructField(v, fieldName); ok {
				return value
			}
		}
	}
	panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
}

// StructField looks up a field of struct v by its expr name. It is the single
// place where struct field names are resolved at runtime, shared by field
// access, the "in" operator, and the builtin get(). It reports false if there
// is no such field, or if the field is behind a nil embedded pointer.
func StructField(v reflect.Value, fieldName string) (any, bool) {
	if index, ok := structFieldIndex(v.Type(), fieldName); ok {
		if value, ok := readField(v, index); ok {
			return value.Interface(), true
		}
		return nil, false
	}
	// Field isn't found via standard Go promotion. Try to find it
	// by traversing embedded interface values whose concrete types
	// may contain the requested field.
	return FetchFromEmbeddedInterfaces(v, fieldName)
}

// structFieldIndex returns the index path of the field of struct type t that
// has the given expr name.
func structFieldIndex(t reflect.Type, name string) ([]int, bool) {
	fields, ok := structFields.Load(t)
	if !ok {
		fields, _ = structFields.LoadOrStore(t, collectStructFields(t))
	}
	index, ok := fields.(map[string][]int)[name]
	return index, ok
}

// collectStructFields maps the expr names of the fields of struct type t to
// their index paths, using the same rules as the type checker (see
// checker/nature/utils.go):
//   - a field is named by its expr tag, or by its Go name if it has no tag;
//   - unexported fields and fields tagged expr:"-" have no name;
//   - a struct's own fields come before the fields of its embedded structs,
//     which are visited depth first in declaration order;
//   - the first field found for a name wins.
func collectStructFields(t reflect.Type) map[string][]int {
	fields := make(map[string][]int)
	visited := make(map[reflect.Type]bool)
	var collect func(t reflect.Type, index []int)
	collect = func(t reflect.Type, index []int) {
		// A type visited before can't add new names: they are all taken.
		if visited[t] {
			return
		}
		visited[t] = true
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			name := f.Name
			switch tag := f.Tag.Get("expr"); tag {
			case "-":
				continue
			case "":
			default:
				name = tag
			}
			if _, ok := fields[name]; !ok {
				fields[name] = append(index[:len(index):len(index)], i)
			}
		}
		// Unexported embedded structs are visited too: their exported
		// fields are promoted.
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if ft := deref.Type(f.Type); f.Anonymous && ft.Kind() == reflect.Struct {
				collect(ft, append(index[:len(index):len(index)], i))
			}
		}
	}
	collect(t, nil)
	return fields
}

// readField returns the field of struct v at the index path, or false if the
// path goes through a nil embedded pointer.
func readField(v reflect.Value, index []int) (reflect.Value, bool) {
	for i, x := range index {
		if i > 0 && v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}, false
			}
			v = v.Elem()
		}
		v = v.Field(x)
	}
	return v, true
}

type Field struct {
	Index []int
	Path  []string
}

func FetchField(from any, field *Field) any {
	v := reflect.ValueOf(from)
	if v.Kind() != reflect.Invalid {
		v = reflect.Indirect(v)

		// We can use v.FieldByIndex here, but it will panic if the field
		// is not exists. And we need to recover() to generate a more
		// user-friendly error message.
		// Also, our fieldByIndex() function is slightly faster than the
		// v.FieldByIndex() function as we don't need to verify what a field
		// is a struct as we already did it on compilation step.
		value := fieldByIndex(v, field)
		if value.IsValid() {
			return value.Interface()
		}
	}
	panic(fmt.Sprintf("cannot get %v from %T", field.Path[0], from))
}

func fieldByIndex(v reflect.Value, field *Field) reflect.Value {
	if len(field.Index) == 1 {
		return v.Field(field.Index[0])
	}
	for i, x := range field.Index {
		if i > 0 {
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					// Index can be longer than Path, as embedded structs
					// add steps to the index that have no name in Path.
					panic(fmt.Sprintf("cannot get %v from nil %v", field.Path[len(field.Path)-1], v.Type()))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}

// FetchFromEmbeddedInterfaces looks up fieldName on the concrete types held by
// anonymous (embedded) interface fields of the struct value v, recursing through
// embedded structs to reach further embedded interfaces. It is shared by Fetch
// and the builtin get() so that both resolve such fields consistently.
func FetchFromEmbeddedInterfaces(v reflect.Value, fieldName string) (any, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.Anonymous {
			continue
		}
		fv := deref.Value(v.Field(i))
		if fv.Kind() != reflect.Struct {
			continue
		}
		// Embedded interfaces need an explicit field lookup on the concrete
		// value. Embedded structs are already covered by structFieldIndex,
		// so we only recurse into them to find further embedded
		// interfaces.
		if deref.Type(f.Type).Kind() == reflect.Interface {
			if index, ok := structFieldIndex(fv.Type(), fieldName); ok {
				if value, ok := readField(fv, index); ok {
					return value.Interface(), true
				}
			}
		}
		if result, found := FetchFromEmbeddedInterfaces(fv, fieldName); found {
			return result, true
		}
	}
	return nil, false
}

type Method struct {
	Index int
	Name  string
}

func FetchMethod(from any, method *Method) any {
	v := reflect.ValueOf(from)
	if v.Kind() != reflect.Invalid {
		// Methods can be defined on any type, no need to dereference.
		if m, ok := MethodByIndex(v, method.Index); ok {
			return m
		}
	}
	panic(fmt.Sprintf("cannot fetch %v from %T", method.Name, from))
}

func Slice(array, from, to any) any {
	v := reflect.ValueOf(array)

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		length := v.Len()
		a, b := ToInt(from), ToInt(to)
		if a < 0 {
			a = length + a
		}
		if a < 0 {
			a = 0
		}
		if b < 0 {
			b = length + b
		}
		if b < 0 {
			b = 0
		}
		if b > length {
			b = length
		}
		if a > b {
			a = b
		}
		if v.Kind() == reflect.Array && !v.CanAddr() {
			newValue := reflect.New(v.Type()).Elem()
			newValue.Set(v)
			v = newValue
		}
		value := v.Slice(a, b)
		if value.IsValid() {
			return value.Interface()
		}

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() {
			return Slice(value.Interface(), from, to)
		}

	}
	panic(fmt.Sprintf("cannot slice %v", from))
}

func In(needle any, array any) bool {
	if array == nil {
		return false
	}
	v := reflect.ValueOf(array)

	switch v.Kind() {

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			if value.IsValid() {
				if Equal(value.Interface(), needle) {
					return true
				}
			}
		}
		return false

	case reflect.Map:
		var value reflect.Value
		if needle == nil {
			value = v.MapIndex(reflect.Zero(v.Type().Key()))
		} else {
			value = v.MapIndex(reflect.ValueOf(needle))
		}
		if value.IsValid() {
			return true
		}
		return false

	case reflect.Struct:
		n := reflect.ValueOf(needle)
		if !n.IsValid() || n.Kind() != reflect.String {
			panic(fmt.Sprintf("cannot use %T as field name of %T", needle, array))
		}
		// Like "in" on a map, this asks whether the field exists, so it
		// doesn't read the field (which may be behind a nil pointer).
		if _, ok := structFieldIndex(v.Type(), n.String()); ok {
			return true
		}
		_, ok := FetchFromEmbeddedInterfaces(v, n.String())
		return ok

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() {
			return In(needle, value.Interface())
		}
		return false
	}

	panic(fmt.Sprintf(`operator "in" not defined on %T`, array))
}

func Len(a any) int {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len()
	default:
		panic(fmt.Sprintf("invalid argument for len (type %T)", a))
	}
}

func Negate(i any) any {
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

func Exponent(a, b any) float64 {
	return math.Pow(ToFloat64(a), ToFloat64(b))
}

func MakeRange(min, max int) []int {
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

func ToInt(a any) int {
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

func ToInt64(a any) int64 {
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

func ToFloat64(a any) float64 {
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
		panic(fmt.Sprintf("invalid operation: float(%T)", x))
	}
}

func ToBool(a any) bool {
	if a == nil {
		return false
	}
	switch x := a.(type) {
	case bool:
		return x
	default:
		panic(fmt.Sprintf("invalid operation: bool(%T)", x))
	}
}

func IsNil(v any) bool {
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
