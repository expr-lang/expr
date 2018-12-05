package expr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/etcinitd/expr"
)

type typeTest string

type typeErrorTest struct {
	input string
	err   string
}

var typeTests = []typeTest{
	"Foo.Bar.Baz",
	"Arr[0].Bar.Baz",
	"Map['string'].Bar.Baz",
	"Map.id.Bar.Baz",
	"Any.Thing.Is.Ok",
	"Irr['string'].next.goes['any thing']",
	"Fn(Any)",
	"Foo.Fn()",
	"true ? Any : Any",
	"Str ~ (true ? Str : Str)",
	"Ok && Any",
	"Str matches 'ok'",
	"Str matches Any",
	"Any matches Any",
	"len([])",
	"true == false",
	"nil",
	"!Ok",
	"[1,2,3]",
	"{id: Foo.Bar.Baz, (1+1): Ok}",
	"Abc()",
	"Foo.Abc()",
	"'a' == 'b' ~ 'c'",
	"Num == 1",
	"Int == Num",
	"Num == Abc",
	"Abc == Num",
	"1 == 2 and true or Ok",
	"Int == Any",
	"IntPtr == Int",
	"!OkPtr == Ok",
	"1 == NumPtr",
	"Foo.Bar == Map.id.Bar",
	"StrPtr == nil",
	"nil == nil",
	"nil == IntPtr",
	"Foo2p.Bar.Baz",
	"Str in Foo",
	"Str in Arr",
	"nil in Arr",
	"Str not in Foo2p",
	"Int | Num",
	"Int ^ Num",
	"Int & Num",
	"Int < Num",
	"Int > Num",
	"Int >= Num",
	"Int <= Num",
	"Int + Num",
	"Int - Num",
	"Int * Num",
	"Int / Num",
	"Int % Num",
	"Int ** Num",
	"Int .. Num",
	"Int + Int + Int",
	"Int % Int > 1",
	"Int in Int..Int",
	"EmbStr == ''",
	"Embedded.EmbStr",
	"EmbPtrStr == ''",
	"EmbeddedPtr.EmbPtrStr ~ Str",
	"SubStr ~ ''",
	"SubEmbedded.SubStr",
	"OkFn() and OkFn()",
	"Foo.Fn() or Foo.Fn()",
	"Method() > 1",
	"Embedded.Method() ~ Str",
}

var typeErrorTests = []typeErrorTest{
	{
		"Foo.Bar.Not",
		"Foo.Bar.Not undefined (type expr_test.bar has no field Not)",
	},
	{
		"Noo",
		"unknown name Noo",
	},
	{
		"Noo()",
		"unknown func Noo()",
	},
	{
		"Foo()",
		"unknown func Foo()",
	},
	{
		"Foo['string']",
		`invalid operation: Foo["string"] (type *expr_test.foo does not support indexing)`,
	},
	{
		"Foo.Fn(Not)",
		"unknown name Not",
	},
	{
		"Foo.Bar()",
		"Foo.Bar() undefined (type *expr_test.foo has no method Bar)",
	},
	{
		"Foo.Bar.Not()",
		"Foo.Bar.Not() undefined (type expr_test.bar has no method Not)",
	},
	{
		"Arr[0].Not",
		"Arr[0].Not undefined (type *expr_test.foo has no field Not)",
	},
	{
		"Arr[Not]",
		"unknown name Not",
	},
	{
		"Not[0]",
		"unknown name Not",
	},
	{
		"Not.Bar",
		"unknown name Not",
	},
	{
		"Arr.Not",
		"Arr.Not undefined (type []*expr_test.foo has no field Not)",
	},
	{
		"Fn(Not)",
		"unknown name Not",
	},
	{
		"Map['str'].Not",
		`Map["str"].Not undefined (type *expr_test.foo has no field Not)`,
	},
	{
		"Ok && IntPtr",
		"invalid operation: Ok && IntPtr (mismatched types bool and *int)",
	},
	{
		"No ? Any.Ok : Any.Not",
		"unknown name No",
	},
	{
		"Any.Cond ? No : Any.Not",
		"unknown name No",
	},
	{
		"Any.Cond ? Any.Ok : No",
		"unknown name No",
	},
	{
		"Many ? Any : Any",
		"non-bool Many (type map[string]interface {}) used as condition",
	},
	{
		"Str matches Int",
		"invalid operation: (Str matches Int) (mismatched types string and int)",
	},
	{
		"Int matches Str",
		"invalid operation: (Int matches Str) (mismatched types int and string)",
	},
	{
		"!Not",
		"unknown name Not",
	},
	{
		"Not == Any",
		"unknown name Not",
	},
	{
		"[Not]",
		"unknown name Not",
	},
	{
		"{id: Not}",
		"unknown name Not",
	},
	{
		"{(Not): Any}",
		"unknown name Not",
	},
	{
		"(nil).Foo",
		"nil.Foo undefined (type <nil> has no field Foo)",
	},
	{
		"(nil)['Foo']",
		`invalid operation: nil["Foo"] (type <nil> does not support indexing)`,
	},
	{
		"1 and false",
		"invalid operation: 1 and false (mismatched types float64 and bool)",
	},
	{
		"true or 0",
		"invalid operation: true or 0 (mismatched types bool and float64)",
	},
	{
		"not IntPtr",
		"invalid operation: not IntPtr (mismatched type *int)",
	},
	{
		"len(Not)",
		"unknown name Not",
	},
	{
		"Int | Ok",
		"invalid operation: Int | Ok (mismatched types int and bool)",
	},
	{
		"Int ^ Ok",
		"invalid operation: Int ^ Ok (mismatched types int and bool)",
	},
	{
		"Int & Ok",
		"invalid operation: Int & Ok (mismatched types int and bool)",
	},
	{
		"Int < Ok",
		"invalid operation: Int < Ok (mismatched types int and bool)",
	},
	{
		"Int > Ok",
		"invalid operation: Int > Ok (mismatched types int and bool)",
	},
	{
		"Int >= Ok",
		"invalid operation: Int >= Ok (mismatched types int and bool)",
	},
	{
		"Int <= Ok",
		"invalid operation: Int <= Ok (mismatched types int and bool)",
	},
	{
		"Int + Ok",
		"invalid operation: Int + Ok (mismatched types int and bool)",
	},
	{
		"Int - Ok",
		"invalid operation: Int - Ok (mismatched types int and bool)",
	},
	{
		"Int * Ok",
		"invalid operation: Int * Ok (mismatched types int and bool)",
	},
	{
		"Int / Ok",
		"invalid operation: Int / Ok (mismatched types int and bool)",
	},
	{
		"Int % Ok",
		"invalid operation: Int % Ok (mismatched types int and bool)",
	},
	{
		"Int ** Ok",
		"invalid operation: Int ** Ok (mismatched types int and bool)",
	},
	{
		"Int .. Ok",
		"invalid operation: Int .. Ok (mismatched types int and bool)",
	},
	{
		"NilFn() and OkFn()",
		"invalid operation: NilFn() and OkFn() (mismatched types <nil> and bool)",
	},
	{
		"'str' in Str",
		`invalid operation: "str" in Str (mismatched types string and string)`,
	},
	{
		"1 in Foo",
		"invalid operation: 1 in Foo (mismatched types float64 and *expr_test.foo)",
	},
	{
		"1 ~ ''",
		`invalid operation: 1 ~ "" (mismatched types float64 and string)`,
	},
}

type abc interface {
	Abc()
}
type bar struct {
	Baz string
}
type foo struct {
	Bar bar
	Fn  func() bool
	Abc abc
}
type SubEmbedded struct {
	SubStr string
}
type Embedded struct {
	SubEmbedded
	EmbStr string
}
type EmbeddedPtr struct {
	EmbPtrStr string
}
type payload struct {
	Embedded
	*EmbeddedPtr
	Abc    abc
	Foo    *foo
	Arr    []*foo
	Map    map[string]*foo
	Any    interface{}
	Irr    []interface{}
	Many   map[string]interface{}
	Fn     func()
	Ok     bool
	Num    float64
	Int    int
	Str    string
	OkPtr  *bool
	NumPtr *float64
	IntPtr *int
	StrPtr *string
	Foo2p  **foo
	OkFn   func() bool
	NilFn  func()
}

func (p payload) Method() int {
	return 0
}

func (p Embedded) Method() string {
	return ""
}

func TestType(t *testing.T) {
	for _, test := range typeTests {
		_, err := expr.Parse(string(test), expr.Env(&payload{}))
		if err != nil {
			t.Errorf("%s:\n\t%+v", test, err.Error())
		}
	}
}

func TestType_error(t *testing.T) {
	for _, test := range typeErrorTests {
		_, err := expr.Parse(test.input, expr.Env(&payload{}))
		if err == nil {
			err = fmt.Errorf("<nil>")
		}
		if !strings.HasPrefix(err.Error(), test.err) || test.err == "" {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, err.Error(), test.err)
		}
	}
}
