package runtime

//go:generate sh -c "go run ./helpers > ./helpers[generated].go"

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"

	"github.com/expr-lang/expr/internal/deref"
)

type contextKey struct{}

// New instantiates a new context with the provided tag.
func New(tag string) *Context {
	return &Context{tag: tag}
}

// FromContext retrieves the *Context stored by NewContext, or nil if none is present.
func FromContext(ctx context.Context) *Context {
	c, _ := ctx.Value(contextKey{}).(*Context)
	return c
}

// fieldCacheKey is used inside Context to memoize struct field index lookups.
// It combines the struct type with the requested field name (which may be a
// tag value).
type fieldCacheKey struct {
	t reflect.Type
	f string
}

// Context holds the struct-tag key and a per-program field index cache for
// runtime struct field resolution. It must not be copied after first use
// (sync.Map is embedded). Embed it in vm.Program so each compiled program
// carries its own isolated context.
type Context struct {
	tag   string
	cache sync.Map
}

// With returns a new Go context with the runtime instance incorporated
// inside it.
func (c *Context) With(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}

// Tag returns the runtime context's tag used in struct fields.
func (c *Context) Tag() string {
	return c.tag
}

// TagName resolves the display name for a struct field under c's configured tag.
// Returns ("", false) if the field is excluded by a "-" tag value.
// Returns ("", true) if no tag is set — caller should fall back to the Go field name.
func (c *Context) TagName(tag reflect.StructTag) (string, bool) {
	if c == nil {
		return "", true
	}
	tagVal := tag.Get(c.tag)
	if i := strings.IndexByte(tagVal, ','); i >= 0 {
		tagVal = tagVal[:i]
	}
	if tagVal == "-" {
		return "", false
	}
	return tagVal, true
}

// Fetch retrieves the value addressed by i from from. For structs, it uses
// c.tag to map i to a Go field name and caches the resolved index in c.cache.
func (c *Context) Fetch(from, i any) any {
	v := reflect.ValueOf(from)
	if v.Kind() == reflect.Invalid {
		panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
	}

	// Methods can be defined on any type.
	if v.NumMethod() > 0 {
		if methodName, ok := i.(string); ok {
			method := v.MethodByName(methodName)
			if method.IsValid() {
				return method.Interface()
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
		fieldName := i.(string)
		t := v.Type()
		key := fieldCacheKey{t: t, f: fieldName}
		if cv, ok := c.cache.Load(key); ok {
			return v.FieldByIndex(cv.([]int)).Interface()
		}
		field, ok := t.FieldByNameFunc(func(name string) bool {
			f, _ := t.FieldByName(name)
			tagName, ok := c.TagName(f.Tag)
			if !ok {
				return false
			}
			if tagName != "" {
				return tagName == fieldName
			}
			return name == fieldName
		})
		if ok && field.IsExported() {
			value := v.FieldByIndex(field.Index)
			if value.IsValid() {
				c.cache.Store(key, field.Index)
				return value.Interface()
			}
		}
	}
	panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
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
					panic(fmt.Sprintf("cannot get %v from %v", field.Path[i], field.Path[i-1]))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}

type Method struct {
	Index int
	Name  string
}

func FetchMethod(from any, method *Method) any {
	v := reflect.ValueOf(from)
	kind := v.Kind()
	if kind != reflect.Invalid {
		// Methods can be defined on any type, no need to dereference.
		method := v.Method(method.Index)
		if method.IsValid() {
			return method.Interface()
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

// In reports whether needle is in array. For structs it checks whether a
// field whose name (or c.tag tag value) equals needle exists and is exported.
func (c *Context) In(needle any, array any) bool {
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
		fieldName := n.String()
		t := v.Type()
		field, ok := t.FieldByNameFunc(func(name string) bool {
			f, _ := t.FieldByName(name)
			tagName, ok := c.TagName(f.Tag)
			if !ok {
				return false
			}
			if tagName != "" {
				return tagName == fieldName
			}
			return name == fieldName
		})
		if !ok || !field.IsExported() {
			return false
		}
		value := v.FieldByIndex(field.Index)
		if value.IsValid() {
			return true
		}
		return false

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() {
			return c.In(needle, value.Interface())
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
