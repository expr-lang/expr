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

func TestCheck_debug(t *testing.T) {
	input := `2**3 + 1`

	tree, err := parser.Parse(input)
	assert.NoError(t, err)

	out, err := checker.Check(tree, checker.Env(&mockEnv{}))
	assert.NoError(t, err)

	if err == nil {
		assert.Equal(t, "float64", out.Name())
	}
}

func TestVisitor_FunctionNode(t *testing.T) {
	var err error

	env := &mockEnv{}
	input := `Set(1, "tag") + Add(2) + Get() + Sub(3) + Any()`

	tree, err := parser.Parse(input)
	assert.NoError(t, err)

	out, err := checker.Check(tree, checker.Env(env))
	assert.NoError(t, err)

	if err == nil {
		assert.Equal(t, "int64", out.Name())
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
		assert.Equal(t, "int64", out.Name())
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
		"!Bool",
		"!BoolPtr == Bool",
		"'a' == 'b' + 'c'",
		"'foo' contains 'bar'",
		"'foo' endsWith 'bar'",
		"'foo' startsWith 'bar'",
		"(1 == 1) || (String matches Any)",
		"1 + 2 + Int64",
		"1 + 2 == FloatPtr",
		"1 + 2 + Float + 3 + 4",
		"1 < Float",
		"1 <= Float",
		"1 == 2 and true or Bool",
		"1 > Float",
		"1 >= Float",
		"2**3 + 1",
		"[1,2,3]",
		"Abc == Float",
		"Abc()",
		"Any matches Any",
		"Any.Thing.Is.Bool",
		"ArrayOfAny['string'].next.goes['any thing']",
		"ArrayOfFoo[0].Bar.Baz",
		"ArrayOfFoo[1].Int64 + 1",
		"Bool && Any",
		"BoolFn() and BoolFn()",
		"EmbedPtr.EmbPtrStr + String",
		"EmbPtrStr == ''",
		"Float == 1",
		"Float == Abc",
		"Fn(true, 1, 'str', Any)",
		"Foo.Abc()",
		"Foo.Bar == Map.id.Bar",
		"Foo.Bar.Baz",
		"Foo.Fn() or Foo.Fn()",
		"Foo.Fn()",
		"Foo2p.Bar.Baz",
		"Int % Int > 1",
		"Int + Int + Int",
		"Int == Any",
		"Int in Int..Int",
		"Int64 % 1",
		"IntPtr == Int",
		"len([])",
		"Map.id.Bar.Baz",
		"Map['string'].Bar.Baz",
		"Method(Foo.Bar) > 1",
		"nil == IntPtr",
		"nil == nil",
		"nil in ArrayOfFoo",
		"nil",
		"String + (true ? String : String)",
		"String in ArrayOfFoo",
		"String in Foo",
		"String matches 'ok'",
		"String matches Any",
		"String not in Foo2p",
		"StringPtr == nil",
		"Sub.Method(0) + String",
		"Sub.SubString",
		"SubStr + ''",
		"SubString == ''",
		"SubSub.SubStr",
		"true == false",
		"true ? Any : Any",
		"{id: Foo.Bar.Baz, 'str': Bool}",
		`"a" < "b"`,
	}
	for _, test := range typeTests {
		var err error

		tree, err := parser.Parse(test)
		assert.NoError(t, err, test)

		_, err = checker.Check(tree, checker.Env(mockEnv2{}))
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
			"ArrayOfFoo[0].Not",
			"type *checker_test.foo has no field Not",
		},
		{
			"ArrayOfFoo[Not]",
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
			"ArrayOfFoo.Not",
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
			"Bool && IntPtr",
			"invalid operation: && (mismatched types bool and *int)",
		},
		{
			"No ? Any.Bool : Any.Not",
			"unknown name No",
		},
		{
			"Any.Cond ? No : Any.Not",
			"unknown name No",
		},
		{
			"Any.Cond ? Any.Bool : No",
			"unknown name No",
		},
		{
			"ManOfAny ? Any : Any",
			"non-bool expression (type map[string]interface {}) used as condition",
		},
		{
			"String matches Int",
			"invalid operation: matches (mismatched types string and int)",
		},
		{
			"Int matches String",
			"invalid operation: matches (mismatched types int and string)",
		},
		{
			"String contains Int",
			"invalid operation: contains (mismatched types string and int)",
		},
		{
			"Int contains String",
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
			"invalid operation: and (mismatched types int and bool)",
		},
		{
			"true or 0",
			"invalid operation: or (mismatched types bool and int)",
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
			"Int < Bool",
			"invalid operation: < (mismatched types int and bool)",
		},
		{
			"Int > Bool",
			"invalid operation: > (mismatched types int and bool)",
		},
		{
			"Int >= Bool",
			"invalid operation: >= (mismatched types int and bool)",
		},
		{
			"Int <= Bool",
			"invalid operation: <= (mismatched types int and bool)",
		},
		{
			"Int + Bool",
			"invalid operation: + (mismatched types int and bool)",
		},
		{
			"Int - Bool",
			"invalid operation: - (mismatched types int and bool)",
		},
		{
			"Int * Bool",
			"invalid operation: * (mismatched types int and bool)",
		},
		{
			"Int / Bool",
			"invalid operation: / (mismatched types int and bool)",
		},
		{
			"Int % Bool",
			"invalid operation: % (mismatched types int and bool)",
		},
		{
			"Int ** Bool",
			"invalid operation: ** (mismatched types int and bool)",
		},
		{
			"Int .. Bool",
			"invalid operation: .. (mismatched types int and bool)",
		},
		{
			"NilFn() and BoolFn()",
			"func NilFn doesn't return value",
		},
		{
			"'str' in String",
			`invalid operation: in (mismatched types string and string)`,
		},
		{
			"1 in Foo",
			"invalid operation: in (mismatched types int and *checker_test.foo)",
		},
		{
			"1 + ''",
			`invalid operation: + (mismatched types int and string)`,
		},
		{
			`all(ArrayOfFoo, {#.Fn() < 0})`,
			`invalid operation: < (mismatched types bool and int)`,
		},
		{
			`map(Any, {0})[0] + "str"`,
			`invalid operation: + (mismatched types int and string)`,
		},
	}

	re, _ := regexp.Compile(`\s*\(\d+:\d+\)\s*`)

	for _, test := range typeErrorTests {

		tree, err := parser.Parse(test.input)
		assert.NoError(t, err)

		_, err = checker.Check(tree, checker.Env(mockEnv2{}))
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
	Int64 int64
	Bar   bar
	Fn    func() bool
	Abc   abc
}

type SubSub struct {
	SubStr string
}

type Sub struct {
	SubSub
	SubString string
}

func (p Sub) Method(i int) string {
	return ""
}

type EmbedPtr struct {
	EmbPtrStr string
}

type mockEnv2 struct {
	Sub
	*EmbedPtr
	Abc        abc
	Foo        *foo
	ArrayOfFoo []*foo
	Map        map[string]*foo
	Any        interface{}
	ArrayOfAny []interface{}
	ManOfAny   map[string]interface{}
	Fn         func(bool, int, string, interface{}) string
	Bool       bool
	Float      float64
	Int64      int64
	Int        int
	String     string
	BoolPtr    *bool
	FloatPtr   *float64
	IntPtr     *int
	StringPtr  *string
	Foo2p      **foo
	BoolFn     func() bool
	NilFn      func()
}

func (p mockEnv2) Method(_ bar) int {
	return 0
}
