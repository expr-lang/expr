package vm

import (
	"reflect"
	"time"
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

type groupBy = map[any][]any

type Span struct {
	Name       string
	Expression string
	Start      time.Time
	Duration   []int64
	Children   []*Span
}
