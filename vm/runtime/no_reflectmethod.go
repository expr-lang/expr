//go:build expr_noreflectmethod

package runtime

import "reflect"

// MethodByName is a no-op stub used when building with the
// expr_noreflectmethod tag. It avoids reaching reflect.Value.MethodByName so
// the Go linker can perform full method dead-code elimination on user types.
func MethodByName(v reflect.Value, name string) (any, bool) {
	return nil, false
}

// MethodByIndex is a no-op stub used when building with the
// expr_noreflectmethod tag. It avoids reaching reflect.Value.Method so the Go
// linker can perform full method dead-code elimination on user types.
func MethodByIndex(v reflect.Value, index int) (any, bool) {
	return nil, false
}
