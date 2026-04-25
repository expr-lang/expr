//go:build expr_noreflectmethod

package nature

// MethodByName is a no-op stub used when building with the
// expr_noreflectmethod tag. It avoids reaching reflect.Type.Method via the
// methodset cache so the Go linker can perform full method dead-code
// elimination on user types.
func (n *Nature) MethodByName(c *Cache, name string) (Nature, bool) {
	return Nature{}, false
}
