package vm

import (
	"reflect"
	"time"
)

type (
	Function = func(params ...any) (any, error)
	// SafeFunction is a function that reports the memory it allocates.
	// It receives the remaining memory budget, so it can check the size of
	// an allocation before making it.
	SafeFunction = func(budget uint, params ...any) (any, uint, error)
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

type Scope struct {
	Array reflect.Value
	Index int
	Len   int
	Count int
	Acc   any
	// Fast paths
	Ints    []int
	Floats  []float64
	Strings []string
	Anys    []any
}

// Item returns the current element from the scope using fast paths when available.
func (s *Scope) Item() any {
	if s.Ints != nil {
		return s.Ints[s.Index]
	}
	if s.Floats != nil {
		return s.Floats[s.Index]
	}
	if s.Strings != nil {
		return s.Strings[s.Index]
	}
	if s.Anys != nil {
		return s.Anys[s.Index]
	}
	return s.Array.Index(s.Index).Interface()
}

type groupBy = map[any][]any

type Span struct {
	Name       string  `json:"name"`
	Expression string  `json:"expression"`
	Duration   int64   `json:"duration"`
	Children   []*Span `json:"children"`
	start      time.Time
}

func GetSpan(program *Program) *Span {
	return program.span
}

// EncodeBuiltinSafe packs a builtin index and its number of arguments into
// the argument of OpCallBuiltinSafe. It reports false if they do not fit.
func EncodeBuiltinSafe(id, argc int) (int, bool) {
	if id < 0 || id >= 1<<15 || argc < 0 || argc >= 1<<16 {
		return 0, false
	}
	return id<<16 | argc, true
}

// DecodeBuiltinSafe unpacks the argument of OpCallBuiltinSafe.
func DecodeBuiltinSafe(arg int) (id, argc int) {
	return arg >> 16, arg & (1<<16 - 1)
}
