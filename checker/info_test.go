package checker_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/expr-lang/expr/checker"
	"github.com/expr-lang/expr/test/mock"
)

func TestTypedFuncIndex(t *testing.T) {
	fn := func([]any, string) string {
		return "foo"
	}
	index, ok := checker.TypedFuncIndex(reflect.TypeOf(fn), false)
	require.True(t, ok)
	require.Equal(t, 22, index)
}

func TestTypedFuncIndex_excludes_named_functions(t *testing.T) {
	var fn mock.MyFunc

	_, ok := checker.TypedFuncIndex(reflect.TypeOf(fn), false)
	require.False(t, ok)
}
