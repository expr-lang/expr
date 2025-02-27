package checker_test

import (
	"reflect"
	"testing"
	"time"

	"expr/internal/testify/require"

	"expr/checker"
)

func TestTypedFuncIndex(t *testing.T) {
	fn := func() time.Duration {
		return 1 * time.Second
	}
	index, ok := checker.TypedFuncIndex(reflect.TypeOf(fn), false)
	require.True(t, ok)
	require.Equal(t, 1, index)
}
