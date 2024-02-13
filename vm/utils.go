package vm

import (
	"reflect"
)

type (
	Function     = func(params ...any) (any, error)
	SafeFunction = func(params ...any) (any, uint, error)
)

var (
	// MemoryBudget represents an upper limit of memory usage.
	MemoryBudget uint = 1e6

	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

type Scope struct {
	Array reflect.Value
	Index int
	Len   int
	Count int
	Acc   any
}

type GroupBy = map[any][]any
