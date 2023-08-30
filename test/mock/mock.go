package mock

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/antonmedv/expr/ast"
)

type Env struct {
	Embed
	Ambiguous          string
	Any                any
	Bool               bool
	Float              float64
	Int64              int64
	Int32              int32
	Int, One, Two      int
	Uint32             uint32
	Uint64             uint64
	Float32            float32
	Float64            float64
	String             string
	BoolPtr            *bool
	FloatPtr           *float64
	IntPtr             *int
	IntPtrPtr          **int
	StringPtr          *string
	Foo                Foo
	Abstract           Abstract
	ArrayOfAny         []any
	ArrayOfInt         []int
	ArrayOfString      []string
	ArrayOfFoo         []*Foo
	MapOfFoo           map[string]Foo
	MapOfAny           map[string]any
	MapIntAny          map[int]string
	FuncParam          func(_ bool, _ int, _ string) bool
	FuncParamAny       func(_ any) bool
	FuncTooManyReturns func() (int, int, error)
	FuncNamed          MyFunc
	NilAny             any
	NilInt             *int
	NilFn              func()
	NilStruct          *Foo
	NilSlice           []any
	Variadic           func(_ int, _ ...int) bool
	Fast               func(...any) any
	Time               time.Time
	TimePlusDay        time.Time
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

func (p Env) TimeEqualString(a time.Time, s string) bool {
	return a.Format("2006-01-02") == s
}

func (p Env) GetInt() int {
	return p.Int
}

func (Env) Add(a, b int) int {
	return a + b
}

func (Env) StringerStringEqual(f fmt.Stringer, s string) bool {
	return f.String() == s
}

func (Env) StringStringerEqual(s string, f fmt.Stringer) bool {
	return s == f.String()
}

func (Env) StringerStringerEqual(f fmt.Stringer, g fmt.Stringer) bool {
	return f.String() == g.String()
}

func (Env) NotStringerStringEqual(f fmt.Stringer, s string) bool {
	return f.String() != s
}

func (Env) NotStringStringerEqual(s string, f fmt.Stringer) bool {
	return s != f.String()
}

func (Env) NotStringerStringerEqual(f fmt.Stringer, g fmt.Stringer) bool {
	return f.String() != g.String()
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
	Value string
	Bar   Bar
}

func (Foo) Method() Bar {
	return Bar{
		Baz: "baz (from Foo.Method)",
	}
}

func (f Foo) MethodWithArgs(prefix string) string {
	return prefix + f.Value
}

func (Foo) String() string {
	return "Foo.String"
}

type Bar struct {
	Baz string
}

type Abstract interface {
	Method(int) int
}

type MyFunc func(string) int

var stringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

type StringerPatcher struct{}

func (*StringerPatcher) Visit(node *ast.Node) {
	t := (*node).Type()
	if t == nil {
		return
	}
	if t.Implements(stringer) {
		ast.Patch(node, &ast.CallNode{
			Callee: &ast.MemberNode{
				Node:     *node,
				Property: &ast.StringNode{Value: "String"},
			},
		})
	}
}

type MapStringStringEnv map[string]string

func (m MapStringStringEnv) Split(s, sep string) []string {
	return strings.Split(s, sep)
}

type MapStringIntEnv map[string]int

type Is struct{}

func (Is) Nil(a any) bool {
	return a == nil
}
