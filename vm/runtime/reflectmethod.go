//go:build !expr_noreflectmethod

package runtime

import "reflect"

// MethodByName resolves a method by name on v at runtime via reflection.
func MethodByName(v reflect.Value, name string) (any, bool) {
	method := v.MethodByName(name)
	if method.IsValid() {
		return method.Interface(), true
	}
	return nil, false
}

// MethodByIndex dispatches a method by index on v at runtime via reflection.
func MethodByIndex(v reflect.Value, index int) (any, bool) {
	method := v.Method(index)
	if method.IsValid() {
		return method.Interface(), true
	}
	return nil, false
}
