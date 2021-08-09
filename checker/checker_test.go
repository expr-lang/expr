package checker_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/stretchr/testify/assert"
)

func TestCheck_debug(t *testing.T) {
	input := `2**3 + 1`

	tree, err := parser.Parse(input)
	assert.NoError(t, err)

	out, err := checker.Check(tree, conf.New(&mockEnv{}))
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

	out, err := checker.Check(tree, conf.New(env))
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
				+ (Duration.String() == "" ? 1 : 0)
				+ Interface.Method(0)
				+ Tickets[0].Method(0)`

	tree, err := parser.Parse(input)
	assert.NoError(t, err)

	out, err := checker.Check(tree, conf.New(env))
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

		_, err = checker.Check(tree, conf.New(&mockEnv{}))
		assert.NoError(t, err)
	}
}

func TestVisitor_ConstantNode(t *testing.T) {
	tree, err := parser.Parse(`re("[a-z]")`)

	regexValue := regexp.MustCompile("[a-z]")
	constNode := &ast.ConstantNode{Value: regexValue}
	ast.Patch(&tree.Node, constNode)

	_, err = checker.Check(tree, conf.New(&mockEnv{}))
	assert.NoError(t, err)

	assert.Equal(t, reflect.TypeOf(regexValue), tree.Node.Type())
}

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
		"1 + Int + Float",
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
		"Variadic('', 1, 2) + Variadic('')",
		"Foo.Variadic('', 1, 2) + Foo.Variadic('')",
		"count(1..30, {# % 3 == 0}) > 0",
		"map(1..3, {#}) == [1,2,3]",
		"map(filter(ArrayOfFoo, {.Int64 > 0}), {.Bar})",
	}
	for _, test := range typeTests {
		var err error

		tree, err := parser.Parse(test)
		assert.NoError(t, err, test)

		_, err = checker.Check(tree, conf.New(mockEnv2{}))
		assert.NoError(t, err, test)
	}
}

const errorTests = `
Foo.Bar.Not
type checker_test.bar has no field Not (1:9)
 | Foo.Bar.Not
 | ........^

Noo
unknown name Noo (1:1)
 | Noo
 | ^

Foo()
unknown func Foo (1:1)
 | Foo()
 | ^

Foo['string']
invalid operation: type *checker_test.foo does not support indexing (1:4)
 | Foo['string']
 | ...^

Foo.Fn(Not)
too many arguments to call Fn (1:5)
 | Foo.Fn(Not)
 | ....^

Foo.Bar()
type *checker_test.foo has no method Bar (1:5)
 | Foo.Bar()
 | ....^

Foo.Bar.Not()
type checker_test.bar has no method Not (1:9)
 | Foo.Bar.Not()
 | ........^

ArrayOfFoo[0].Not
type *checker_test.foo has no field Not (1:15)
 | ArrayOfFoo[0].Not
 | ..............^

ArrayOfFoo[Not]
unknown name Not (1:12)
 | ArrayOfFoo[Not]
 | ...........^

Not[0]
unknown name Not (1:1)
 | Not[0]
 | ^

Not.Bar
unknown name Not (1:1)
 | Not.Bar
 | ^

ArrayOfFoo.Not
type []*checker_test.foo has no field Not (1:12)
 | ArrayOfFoo.Not
 | ...........^

Fn(Not)
not enough arguments to call Fn (1:1)
 | Fn(Not)
 | ^

Map['str'].Not
type *checker_test.foo has no field Not (1:12)
 | Map['str'].Not
 | ...........^

Bool && IntPtr
invalid operation: && (mismatched types bool and *int) (1:6)
 | Bool && IntPtr
 | .....^

No ? Any.Bool : Any.Not
unknown name No (1:1)
 | No ? Any.Bool : Any.Not
 | ^

Any.Cond ? No : Any.Not
unknown name No (1:12)
 | Any.Cond ? No : Any.Not
 | ...........^

Any.Cond ? Any.Bool : No
unknown name No (1:23)
 | Any.Cond ? Any.Bool : No
 | ......................^

ManOfAny ? Any : Any
non-bool expression (type map[string]interface {}) used as condition (1:1)
 | ManOfAny ? Any : Any
 | ^

String matches Int
invalid operation: matches (mismatched types string and int) (1:8)
 | String matches Int
 | .......^

Int matches String
invalid operation: matches (mismatched types int and string) (1:5)
 | Int matches String
 | ....^

String contains Int
invalid operation: contains (mismatched types string and int) (1:8)
 | String contains Int
 | .......^

Int contains String
invalid operation: contains (mismatched types int and string) (1:5)
 | Int contains String
 | ....^

!Not
unknown name Not (1:2)
 | !Not
 | .^

Not == Any
unknown name Not (1:1)
 | Not == Any
 | ^

[Not]
unknown name Not (1:2)
 | [Not]
 | .^

{id: Not}
unknown name Not (1:6)
 | {id: Not}
 | .....^

(nil).Foo
type <nil> has no field Foo (1:7)
 | (nil).Foo
 | ......^

(nil)['Foo']
invalid operation: type <nil> does not support indexing (1:6)
 | (nil)['Foo']
 | .....^

1 and false
invalid operation: and (mismatched types int and bool) (1:3)
 | 1 and false
 | ..^

true or 0
invalid operation: or (mismatched types bool and int) (1:6)
 | true or 0
 | .....^

not IntPtr
invalid operation: not (mismatched type *int) (1:1)
 | not IntPtr
 | ^

len(Not)
unknown name Not (1:5)
 | len(Not)
 | ....^

Int < Bool
invalid operation: < (mismatched types int and bool) (1:5)
 | Int < Bool
 | ....^

Int > Bool
invalid operation: > (mismatched types int and bool) (1:5)
 | Int > Bool
 | ....^

Int >= Bool
invalid operation: >= (mismatched types int and bool) (1:5)
 | Int >= Bool
 | ....^

Int <= Bool
invalid operation: <= (mismatched types int and bool) (1:5)
 | Int <= Bool
 | ....^

Int + Bool
invalid operation: + (mismatched types int and bool) (1:5)
 | Int + Bool
 | ....^

Int - Bool
invalid operation: - (mismatched types int and bool) (1:5)
 | Int - Bool
 | ....^

Int * Bool
invalid operation: * (mismatched types int and bool) (1:5)
 | Int * Bool
 | ....^

Int / Bool
invalid operation: / (mismatched types int and bool) (1:5)
 | Int / Bool
 | ....^

Int % Bool
invalid operation: % (mismatched types int and bool) (1:5)
 | Int % Bool
 | ....^

Int ** Bool
invalid operation: ** (mismatched types int and bool) (1:5)
 | Int ** Bool
 | ....^

Int .. Bool
invalid operation: .. (mismatched types int and bool) (1:5)
 | Int .. Bool
 | ....^

NilFn() and BoolFn()
func NilFn doesn't return value (1:1)
 | NilFn() and BoolFn()
 | ^

'str' in String
invalid operation: in (mismatched types string and string) (1:7)
 | 'str' in String
 | ......^

1 in Foo
invalid operation: in (mismatched types int and *checker_test.foo) (1:3)
 | 1 in Foo
 | ..^

1 + ''
invalid operation: + (mismatched types int and string) (1:3)
 | 1 + ''
 | ..^

all(ArrayOfFoo, {#.Fn() < 0})
invalid operation: < (mismatched types bool and int) (1:25)
 | all(ArrayOfFoo, {#.Fn() < 0})
 | ........................^

map(Any, {0})[0] + "str"
invalid operation: + (mismatched types int and string) (1:18)
 | map(Any, {0})[0] + "str"
 | .................^

Variadic()
not enough arguments to call Variadic (1:1)
 | Variadic()
 | ^

Variadic('', '')
cannot use string as argument (type int) to call Variadic  (1:14)
 | Variadic('', '')
 | .............^

Foo.Variadic()
not enough arguments to call Variadic (1:5)
 | Foo.Variadic()
 | ....^

Foo.Variadic('', '')
cannot use string as argument (type int) to call Variadic  (1:18)
 | Foo.Variadic('', '')
 | .................^

count(1, {#})
builtin count takes only array (got int) (1:7)
 | count(1, {#})
 | ......^

count(ArrayOfInt, {#})
closure should return boolean (got int) (1:19)
 | count(ArrayOfInt, {#})
 | ..................^

all(ArrayOfInt, {# + 1})
closure should return boolean (got int) (1:17)
 | all(ArrayOfInt, {# + 1})
 | ................^

filter(ArrayOfFoo, {.Int64})
closure should return boolean (got int64) (1:20)
 | filter(ArrayOfFoo, {.Int64})
 | ...................^

map(1, {2})
builtin map takes only array (got int) (1:5)
 | map(1, {2})
 | ....^

map(filter(ArrayOfFoo, {.Int64 > 0}), {.Var})
type *checker_test.foo has no field Var (1:41)
 | map(filter(ArrayOfFoo, {.Int64 > 0}), {.Var})
 | ........................................^
`

func TestCheck_error(t *testing.T) {
	tests := strings.Split(strings.Trim(errorTests, "\n"), "\n\n")

	for _, test := range tests {
		input := strings.SplitN(test, "\n", 2)
		if len(input) != 2 {
			t.Errorf("syntax error in test: %q", test)
			break
		}

		tree, err := parser.Parse(input[0])
		assert.NoError(t, err)

		_, err = checker.Check(tree, conf.New(mockEnv2{}))
		if err == nil {
			err = fmt.Errorf("<nil>")
		}

		assert.Equal(t, input[1], err.Error(), input[0])
	}
}

func TestCheck_AsBool(t *testing.T) {
	input := `1+2`

	tree, err := parser.Parse(input)
	assert.NoError(t, err)

	config := &conf.Config{}
	expr.AsBool()(config)

	_, err = checker.Check(tree, config)
	assert.Error(t, err)
	assert.Equal(t, "expected bool, but got int", err.Error())
}

//
// Mock types
//

type mockEnv struct {
	*mockEmbed
	Add       func(int64) int64
	Any       interface{}
	Var       *mockVar
	Tickets   []mockTicket
	Duration  time.Duration
	Interface mockInterface
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

type mockInterface interface {
	Method(int) int
}

type mockTicket struct {
	Price  int
	Origin string
}

func (t mockTicket) Method(int) int {
	return 0
}

type abc interface {
	Abc()
}

type bar struct {
	Baz string
}

type foo struct {
	Int64    int64
	Bar      bar
	Fn       func() bool
	Abc      abc
	Variadic func(head string, xs ...int) int
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
	ArrayOfInt []int
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
	Variadic   func(head string, xs ...int) int
}

func (p mockEnv2) Method(_ bar) int {
	return 0
}
