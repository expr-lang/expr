package builtin

import (
	"context"
	"reflect"
)

type Function struct {
	Name            string
	Fast            func(arg any) any
	Func            func(args ...any) (any, error)
	Safe            func(args ...any) (any, uint, error)
	FuncWithContext func(ctx context.Context, args ...any) (any, error)
	Types           []reflect.Type
	Validate        func(args []reflect.Type) (reflect.Type, error)
	Deref           func(i int, arg reflect.Type) bool
	Predicate       bool
}

func (f *Function) Type() reflect.Type {
	if len(f.Types) > 0 {
		return f.Types[0]
	}
	return reflect.TypeOf(f.Func)
}
