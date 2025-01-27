package deref

import (
	"fmt"
	"reflect"
)

func Deref(p any) any {
	if p == nil {
		return nil
	}

	v := reflect.ValueOf(p)

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.IsValid() {
		return v.Interface()
	}

	panic(fmt.Sprintf("cannot dereference %v", p))
}

func Type(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}

	// Preserve interface types immediately to maintain type information
	// This handles both empty (interface{}) and non-empty (e.g., io.Reader) interfaces
	if t.Kind() == reflect.Interface {
		return t
	}

	// Iteratively unwrap pointer types until we reach a non-pointer
	// or encounter an interface type that needs preservation
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t == nil {
			return nil
		}
		// Stop unwrapping if we hit an interface type to preserve its type information
		// This ensures interface method sets are not lost
		if t.Kind() == reflect.Interface {
			return t
		}
	}

	// Return the final unwrapped type, which could be any non-pointer, non-interface type
	return t
}

func Value(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
}
