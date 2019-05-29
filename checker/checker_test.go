package checker_test

import (
	"fmt"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/parser"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestVisitor_FunctionNode(t *testing.T) {
	var err error

	env := &mockEnv{}
	input := `Set(1, "tag") + Add(2) + Get() + Sub(3) + Any()`

	tree, err := parser.Parse(input)
	assert.NoError(t, err)

	out, err := checker.Check(tree, checker.Env(env))
	assert.NoError(t, err)

	if err == nil {
		assert.Equal(t, out.Name(), "int64")
	}
}

func TestVisitor_MethodNode(t *testing.T) {
	var err error

	env := &mockEnv{}
	input := `Var.Set(1, 0.5) 
				+ Var.Add(2) 
				+ Var.Any(true) 
				+ Var.Get() 
				+ Var.Sub(3)
				+ (Duration.String() == "" ? 1 : 0)`

	tree, err := parser.Parse(input)
	assert.NoError(t, err)

	out, err := checker.Check(tree, checker.Env(env))
	assert.NoError(t, err)

	if err == nil {
		assert.Equal(t, out.Name(), "int64")
	}
}

func TestVisitor_BuiltinNode(t *testing.T) {
	var typeTests = []string{
		`all(Tickets, {.Price > 0}) && any(map(Tickets, {.Price}), {# < 1000})`,
		`filter(map(Tickets, {.Origin}), {len(#) != 3})[0]`,
		`none(Any, {#.Any < 1})`,
		`none(Any, {.Thing != "awesome"})`,
	}

	for _, input := range typeTests {
		tree, err := parser.Parse(input)
		assert.NoError(t, err)

		_, err = checker.Check(tree, checker.Env(&mockEnv{}))
		assert.NoError(t, err)
	}
}

// Helper types and declarations.

type mockEnv struct {
	*mockEmbed
	Add      func(int64) int64
	Any      interface{}
	Var      *mockVar
	Tickets  []mockTicket
	Duration time.Duration
}

func (f *mockEnv) Set(v int64, any interface{}) int64 {
	return v
}

type mockEmbed struct {
	EmbedVar int64
	Sub      func(int64) int64
}

func (f *mockEmbed) Get() int64 {
	return 0
}

type mockVar struct {
	*mockEmbed
	Add func(int64) int64
	Any interface{}
}

func (*mockVar) Set(v int64, f float64) int64 {
	return 0
}

type mockTicket struct {
	Price  int
	Origin string
}

// Other tests.

func TestCheck(t *testing.T) {
	var typeTests = []string{
		"Foo.Bar.Baz",
		"Arr[0].Bar.Baz",
		"Map['string'].Bar.Baz",
		"Map.id.Bar.Baz",
		"Any.Thing.Is.Ok",
		"Irr['string'].next.goes['any thing']",
		"Fn(true, 1, 'str', Any)",
		"Foo.Fn()",
		"true ? Any : Any",
		"Str + (true ? Str : Str)",
		"Ok && Any",
		"Str matches 'ok'",
		"Str matches Any",
		"Any matches Any",
		"'foo' contains 'bar'",
		"'foo' startsWith 'bar'",
		"'foo' endsWith 'bar'",
		"len([])",
		"true == false",
		"nil",
		"!Ok",
		"[1,2,3]",
		"{id: Foo.Bar.Baz, 'str': Ok}",
		"Abc()",
		"Foo.Abc()",
		"'a' == 'b' + 'c'",
		"Num == 1",
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
		"1 < Num",
		"1 > Num",
		"1 >= Num",
		"1 <= Num",
		"Int + Int + Int",
		"Int % Int > 1",
		"Int in Int..Int",
		"EmbStr == ''",
		"Embedded.EmbStr",
		"EmbPtrStr == ''",
		"EmbeddedPtr.EmbPtrStr + Str",
		"SubStr + ''",
		"SubEmbedded.SubStr",
		"OkFn() and OkFn()",
		"Foo.Fn() or Foo.Fn()",
		"Method(Foo.Bar) > 1",
		"Embedded.Method() + Str",
		`"a" < "b"`,
	}
	for _, test := range typeTests {
		var err error

		tree, err := parser.Parse(test)
		assert.NoError(t, err, test)

		_, err = checker.Check(tree, checker.Env(Env{}))
		assert.NoError(t, err, test)
	}
}

func TestCheck_error(t *testing.T) {
	type test struct {
		input string
		err   string
	}
	var typeErrorTests = []test{
		{
			"Foo.Bar.Not",
			"type checker_test.bar has no field Not",
		},
		{
			"Noo",
			"unknown name Noo",
		},
		{
			"Noo()",
			"unknown func Noo",
		},
		{
			"Foo()",
			"unknown func Foo",
		},
		{
			"Foo['string']",
			`invalid operation: type *checker_test.foo does not support indexing`,
		},
		{
			"Foo.Fn(Not)",
			"too many arguments to call Fn",
		},
		{
			"Foo.Bar()",
			"type *checker_test.foo has no method Bar",
		},
		{
			"Foo.Bar.Not()",
			"type checker_test.bar has no method Not",
		},
		{
			"Arr[0].Not",
			"type *checker_test.foo has no field Not",
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
			"type []*checker_test.foo has no field Not",
		},
		{
			"Fn(Not)",
			"not enough arguments to call Fn",
		},
		{
			"Map['str'].Not",
			`type *checker_test.foo has no field Not`,
		},
		{
			"Ok && IntPtr",
			"invalid operation: && (mismatched types bool and *int)",
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
			"non-bool expression (type map[string]interface {}) used as condition",
		},
		{
			"Str matches Int",
			"invalid operation: matches (mismatched types string and int)",
		},
		{
			"Int matches Str",
			"invalid operation: matches (mismatched types int and string)",
		},
		{
			"Str contains Int",
			"invalid operation: contains (mismatched types string and int)",
		},
		{
			"Int contains Str",
			"invalid operation: contains (mismatched types int and string)",
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
			"(nil).Foo",
			"type <nil> has no field Foo",
		},
		{
			"(nil)['Foo']",
			`invalid operation: type <nil> does not support indexing`,
		},
		{
			"1 and false",
			"invalid operation: and (mismatched types int64 and bool)",
		},
		{
			"true or 0",
			"invalid operation: or (mismatched types bool and int64)",
		},
		{
			"not IntPtr",
			"invalid operation: not (mismatched type *int)",
		},
		{
			"len(Not)",
			"unknown name Not",
		},
		{
			"Int < Ok",
			"invalid operation: < (mismatched types int and bool)",
		},
		{
			"Int > Ok",
			"invalid operation: > (mismatched types int and bool)",
		},
		{
			"Int >= Ok",
			"invalid operation: >= (mismatched types int and bool)",
		},
		{
			"Int <= Ok",
			"invalid operation: <= (mismatched types int and bool)",
		},
		{
			"Int + Ok",
			"invalid operation: + (mismatched types int and bool)",
		},
		{
			"Int - Ok",
			"invalid operation: - (mismatched types int and bool)",
		},
		{
			"Int * Ok",
			"invalid operation: * (mismatched types int and bool)",
		},
		{
			"Int / Ok",
			"invalid operation: / (mismatched types int and bool)",
		},
		{
			"Int % Ok",
			"invalid operation: % (mismatched types int and bool)",
		},
		{
			"Int ** Ok",
			"invalid operation: ** (mismatched types int and bool)",
		},
		{
			"Int .. Ok",
			"invalid operation: .. (mismatched types int and bool)",
		},
		{
			"NilFn() and OkFn()",
			"func NilFn doesn't return value",
		},
		{
			"'str' in Str",
			`invalid operation: in (mismatched types string and string)`,
		},
		{
			"1 in Foo",
			"invalid operation: in (mismatched types int64 and *checker_test.foo)",
		},
		{
			"1 + ''",
			`invalid operation: + (mismatched types int64 and string)`,
		},
		{
			`all(Arr, {#.Fn() < 0})`,
			`invalid operation: < (mismatched types bool and int64)`,
		},
		{
			`map(Any, {0})[0] + "str"`,
			`invalid operation: + (mismatched types int64 and string)`,
		},
	}

	re, _ := regexp.Compile(`\s*\(\d+:\d+\)\s*`)

	for _, test := range typeErrorTests {

		tree, err := parser.Parse(test.input)
		assert.NoError(t, err)

		_, err = checker.Check(tree, checker.Env(Env{}))
		if err == nil {
			err = fmt.Errorf("<nil>")
		}

		// Trim code snippet.
		lines := strings.Split(err.Error(), "\n")
		firstLine := string(re.ReplaceAll([]byte(lines[0]), []byte{}))

		assert.Equal(t, test.err, firstLine, test.input)
	}
}

// Other helper types.

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

type Env struct {
	Embedded
	*EmbeddedPtr
	Abc    abc
	Foo    *foo
	Arr    []*foo
	Map    map[string]*foo
	Any    interface{}
	Irr    []interface{}
	Many   map[string]interface{}
	Fn     func(bool, int64, string, interface{}) string
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

func (p Env) Method(_ bar) int {
	return 0
}

func (p Embedded) Method() string {
	return ""
}
