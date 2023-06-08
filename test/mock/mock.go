package mock

import (
	"context"
	"time"
)

type Env struct {
	Embed
	Ambiguous          string
	Any                interface{}
	Bool               bool
	Float              float64
	Int64              int64
	Int32              int32
	Int                int
	Uint32             uint32
	String             string
	BoolPtr            *bool
	FloatPtr           *float64
	IntPtr             *int
	IntPtrPtr          **int
	StringPtr          *string
	Foo                Foo
	Abstract           Abstract
	ArrayOfAny         []interface{}
	ArrayOfInt         []int
	ArrayOfFoo         []Foo
	MapOfFoo           map[string]Foo
	MapOfAny           map[string]interface{}
	FuncCtx            func(context.Context, bool, int, string) bool
	FuncParam          func(_ bool, _ int, _ string) bool
	FuncParamAny       func(_ interface{}) bool
	FuncTooManyReturns func() (int, int, error)
	FuncNamed          MyFunc
	NilFn              func()
	Variadic           func(_ int, _ ...int) bool
	Fast               func(...interface{}) interface{}
	Time               time.Time
	Duration           time.Duration
}

func (p Env) FuncFoo(_ Foo) int {
	return 0
}

func (p Env) Func() int {
	return 0
}

func (p Env) FuncTyped(_ string) int {
	return 2023
}

type Embed struct {
	EmbedEmbed
	EmbedString string
}

func (p Embed) EmbedMethod(_ int) string {
	return ""
}

type EmbedEmbed struct {
	EmbedEmbedString string
}

type Foo struct {
	Bar Bar
}

func (Foo) Method() Bar {
	return Bar{}
}

type Bar struct {
	Baz string
}

type Abstract interface {
	Method(int) int
}

type MyFunc func(string) int
