//go:build expr_noreflectmethod

package assert

import "testing"

// SkipNoReflectMethod skips the current test under the expr_noreflectmethod
// build tag, where reflect-based env method dispatch is unavailable.
func SkipNoReflectMethod(t *testing.T) {
	t.Helper()
	t.Skip("requires reflect-based method dispatch")
}
