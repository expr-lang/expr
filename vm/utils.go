package vm

import (
	"reflect"
)

type (
	Function     = func(params ...any) (any, error)
	SafeFunction = func(params ...any) (any, uint, error)
)

// MemoryBudget represents an upper limit of memory usage.
var MemoryBudget uint = 1e6

var errorType = reflect.TypeOf((*error)(nil)).Elem()
