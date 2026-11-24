//go:build !expr_noreflectmethod

package assert

import "testing"

// SkipNoReflectMethod is a no-op in builds without the expr_noreflectmethod
// tag.
func SkipNoReflectMethod(t *testing.T) {}
