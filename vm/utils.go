package vm

import (
	"reflect"
	"time"

	"github.com/expr-lang/expr/vm/runtime"
)

type (
	Function     = func(params ...any) (any, error)
	SafeFunction = func(params ...any) (any, uint, error)
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

type uniqBy struct {
	Keys     []any
	Items    []any
	Hashable map[any]struct{}
}

func newUniqBy(size int) *uniqBy {
	return &uniqBy{
		Keys:     make([]any, 0, size),
		Items:    make([]any, 0, size),
		Hashable: make(map[any]struct{}, size),
	}
}

func (u *uniqBy) Add(key, item any) {
	if hash, ok := uniqByHash(key); ok {
		if _, exists := u.Hashable[hash]; exists {
			return
		}
		u.Hashable[hash] = struct{}{}
		u.Keys = append(u.Keys, key)
		u.Items = append(u.Items, item)
		return
	}

	for _, seen := range u.Keys {
		if runtime.Equal(key, seen) {
			return
		}
	}
	u.Keys = append(u.Keys, key)
	u.Items = append(u.Items, item)
}

func uniqByHash(key any) (any, bool) {
	if runtime.IsNil(key) {
		return nil, true
	}
	switch key := key.(type) {
	case string, bool, time.Duration:
		return key, true
	case time.Time:
		return key.UTC(), true
	default:
		return nil, false
	}
}

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
