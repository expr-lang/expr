package checker_test

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/checker/mock"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

//func TestVisitor_FunctionNode(t *testing.T) {
//	var err error
//
//	env := &mockEnv{}
//	input := `Set(1, "tag") + Add(2) + Get() + Sub(3) + Any()`
//
//	tree, err := parser.Parse(input)
//	assert.NoError(t, err)
//
//	out, err := checker.Check(tree, conf.New(env))
//	assert.NoError(t, err)
//
//	if err == nil {
//		assert.Equal(t, "int64", out.Name())
//	}
//}

//func TestVisitor_MethodNode(t *testing.T) {
//	var err error
//
//	env := &mockEnv{}
//	input := `Var.Set(1, 0.5)
//				+ Var.Add(2)
//				+ Var.Any(true)
//				+ Var.Get()
//				+ Var.Sub(3)
//				+ (Duration.String() == "" ? 1 : 0)
//				+ Interface.Method(0)
//				+ Tickets[0].Method(0)`
//
//	tree, err := parser.Parse(input)
//	assert.NoError(t, err)
//
//	out, err := checker.Check(tree, conf.New(env))
//	assert.NoError(t, err)
//
//	if err == nil {
//		assert.Equal(t, "int64", out.Name())
//	}
//}

//func TestVisitor_BuiltinNode(t *testing.T) {
//	typeTests := []string{
//		`all(Tickets, {.Price > 0}) && any(map(Tickets, {.Price}), {# < 1000})`,
//		`filter(map(Tickets, {.Origin}), {len(#) != 3})[0]`,
//		`none(Any, {#.Any < 1})`,
//		`none(Any, {.Thing != "awesome"})`,
//	}
//
//	for _, input := range typeTests {
//		tree, err := parser.Parse(input)
//		assert.NoError(t, err)
//
//		_, err = checker.Check(tree, conf.New(&mockEnv{}))
//		assert.NoError(t, err)
//	}
//}

var successTests = []string{
	"nil == nil",
	"nil == IntPtr",
	"nil == nil",
	"nil in ArrayOfFoo",
	"!Bool",
	"!BoolPtr == Bool",
	"'a' == 'b' + 'c'",
	"'foo' contains 'bar'",
	"'foo' endsWith 'bar'",
	"'foo' startsWith 'bar'",
	"(1 == 1) || (String matches Any)",
	"Int % Int > 1",
	"Int + Int + Int > 0",
	"Int == Any",
	"Int in Int..Int",
	"IntPtrPtr + 1 > 0",
	"1 + 2 + Int64 > 0",
	"Int64 % 1 > 0",
	"IntPtr == Int",
	"FloatPtr == 1 + 2.",
	"1 + 2 + Float + 3 + 4 < 0",
	"1 + Int + Float == 0.5",
	"-1 + +1 == 0",
	"1 / 2 == 0",
	"2**3 + 1 != 0",
	"Float == 1",
	"Float < 1.0",
	"Float <= 1.0",
	"Float > 1.0",
	"Float >= 1.0",
	`"a" < "b"`,
	"String + (true ? String : String) == ''",
	"(Any ? nil : '') == ''",
	"(Any ? 0 : nil) == 0",
	"(Any ? nil : nil) == nil",
	"!(Any ? Foo : Foo.Bar).Anything",
	"String in ArrayOfFoo",
	"String in Foo",
	"String in MapOfFoo",
	"String matches 'ok'",
	"String matches Any",
	"String not in ArrayOfFoo",
	"StringPtr == nil",
	"[1, 2, 3] == []",
	"len([]) > 0",
	"Any matches Any",
	"!Any.Things.Contains.Any",
	"!ArrayOfAny['string'].next.goes['any thing']",
	"ArrayOfFoo[0].Bar.Baz == ''",
	"ArrayOfFoo[0:10][0].Bar.Baz == ''",
	"Bool && Any",
	"FuncParam(true, 1, 'str')",
	"FuncParamAny(nil)",
	"!Fast(Any, String)",
	"Foo.Method().Baz == ''",
	"Foo.Bar == MapOfAny.id.Bar",
	"Foo.Bar.Baz == ''",
	"MapOfFoo['any'].Bar.Baz == ''",
	"Func(Foo) > 1",
	"Any() > 0",
	"Embed.EmbedString == ''",
	"EmbedString == ''",
	"EmbedMethod(0) == ''",
	"Embed.EmbedMethod(0) == ''",
	"Embed.EmbedString == ''",
	"EmbedString == ''",
	"{id: Foo.Bar.Baz, 'str': String} == {}",
	"Variadic(0, 1, 2) || Variadic(0)",
	"count(1..30, {# % 3 == 0}) > 0",
	"map(1..3, {#}) == [1,2,3]",
	"map(filter(ArrayOfFoo, {.Bar.Baz != ''}), {.Bar}) == []",
	"filter(Any, {.AnyMethod()})[0] == ''",
	"Time == Time",
	"Any == Time",
	"Any != Time",
	"Any > Time",
	"Any >= Time",
	"Any < Time",
	"Any <= Time",
	"Any - Time > Duration",
	"Any == Any",
	"Any != Any",
	"Any > Any",
	"Any >= Any",
	"Any < Any",
	"Any <= Any",
	"Any - Any < Duration",
	"Time == Any",
	"Time != Any",
	"Time > Any",
	"Time >= Any",
	"Time < Any",
	"Time <= Any",
	"Time - Any == Duration",
	"Time + Duration == Time",
	"Duration + Time == Time",
	// TODO: Next check not work as expected, we need to refactor binary node
	// for working well with any (interface{}) types.
	// "Duration + Any == Time",
	// "Any + Duration == Time",
}

func TestCheck(t *testing.T) {
	for _, input := range successTests {
		var err error
		tree, err := parser.Parse(input)
		require.NoError(t, err, input)

		config := conf.New(mock.Env{})
		expr.AsBool()(config)

		_, err = checker.Check(tree, config)
		assert.NoError(t, err, input)
	}
}

var nilsafeTests = []string{
	"any?.value == nil",
	"any?.Method() == nil",
}

func TestCheck_NilSafe(t *testing.T) {
	for _, input := range nilsafeTests {
		var err error
		tree, err := parser.Parse(input)
		require.NoError(t, err, input)

		config := conf.New(mock.Env{})
		expr.AsBool()(config)

		_, err = checker.Check(tree, config)
		assert.NoError(t, err, input)
	}
}

const errorTests = `
Foo.Bar.Not
type mock.Bar has no field Not (1:9)
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
invalid operation: type mock.Foo does not support indexing (1:4)
 | Foo['string']
 | ...^

Foo.Method(Not)
too many arguments to call Method (1:5)
 | Foo.Method(Not)
 | ....^

Foo.Bar()
type mock.Foo has no method Bar (1:5)
 | Foo.Bar()
 | ....^

Foo.Bar.Not()
type mock.Bar has no method Not (1:9)
 | Foo.Bar.Not()
 | ........^

ArrayOfFoo[0].Not
type mock.Foo has no field Not (1:15)
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
type []mock.Foo has no field Not (1:12)
 | ArrayOfFoo.Not
 | ...........^

FuncParam(Not)
not enough arguments to call FuncParam (1:1)
 | FuncParam(Not)
 | ^

MapOfFoo['str'].Not
type mock.Foo has no field Not (1:17)
 | MapOfFoo['str'].Not
 | ................^

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

MapOfAny ? Any : Any
non-bool expression (type map[string]interface {}) used as condition (1:1)
 | MapOfAny ? Any : Any
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
invalid operation: in (mismatched types int and mock.Foo) (1:3)
 | 1 in Foo
 | ..^

1 + ''
invalid operation: + (mismatched types int and string) (1:3)
 | 1 + ''
 | ..^

all(ArrayOfFoo, {#.Method() < 0})
invalid operation: < (mismatched types mock.Bar and int) (1:29)
 | all(ArrayOfFoo, {#.Method() < 0})
 | ............................^

map(Any, {0})[0] + "str"
invalid operation: + (mismatched types int and string) (1:18)
 | map(Any, {0})[0] + "str"
 | .................^

Variadic()
not enough arguments to call Variadic (1:1)
 | Variadic()
 | ^

Variadic(0, '')
cannot use string as argument (type int) to call Variadic  (1:13)
 | Variadic(0, '')
 | ............^

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

filter(ArrayOfFoo, {.Bar.Baz})
closure should return boolean (got string) (1:20)
 | filter(ArrayOfFoo, {.Bar.Baz})
 | ...................^

map(1, {2})
builtin map takes only array (got int) (1:5)
 | map(1, {2})
 | ....^

map(filter(ArrayOfFoo, {true}), {.Not})
type mock.Foo has no field Not (1:35)
 | map(filter(ArrayOfFoo, {true}), {.Not})
 | ..................................^

ArrayOfFoo[Foo]
invalid operation: cannot use mock.Foo as index to mock.Foo (1:12)
 | ArrayOfFoo[Foo]
 | ...........^

ArrayOfFoo[Bool:]
invalid operation: non-integer slice index bool (1:12)
 | ArrayOfFoo[Bool:]
 | ...........^

ArrayOfFoo[1:Bool]
invalid operation: non-integer slice index bool (1:14)
 | ArrayOfFoo[1:Bool]
 | .............^

Bool[:]
invalid operation: cannot slice bool (1:5)
 | Bool[:]
 | ....^

FuncTooManyReturns()
func FuncTooManyReturns returns more then two values (1:1)
 | FuncTooManyReturns()
 | ^

len(42)
invalid argument for len (type int) (1:1)
 | len(42)
 | ^

any(42, {#})
builtin any takes only array (got int) (1:5)
 | any(42, {#})
 | ....^

filter(42, {#})
builtin filter takes only array (got int) (1:8)
 | filter(42, {#})
 | .......^
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
		require.NoError(t, err)

		_, err = checker.Check(tree, conf.New(mock.Env{}))
		if err == nil {
			err = fmt.Errorf("<nil>")
		}

		assert.Equal(t, input[1], err.Error(), input[0])
	}
}

func TestVisitor_ConstantNode(t *testing.T) {
	tree, err := parser.Parse(`re("[a-z]")`)
	require.NoError(t, err)

	regexValue := regexp.MustCompile("[a-z]")
	constNode := &ast.ConstantNode{Value: regexValue}
	ast.Patch(&tree.Node, constNode)

	_, err = checker.Check(tree, nil)
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(regexValue), tree.Node.Type())
}

func TestCheck_AsBool(t *testing.T) {
	tree, err := parser.Parse(`1+2`)
	require.NoError(t, err)

	config := &conf.Config{}
	expr.AsBool()(config)

	_, err = checker.Check(tree, config)
	assert.Error(t, err)
	assert.Equal(t, "expected bool, but got int", err.Error())
}

func TestCheck_AsInt64(t *testing.T) {
	tree, err := parser.Parse(`true`)
	require.NoError(t, err)

	config := &conf.Config{}
	expr.AsInt64()(config)

	_, err = checker.Check(tree, config)
	assert.Error(t, err)
	assert.Equal(t, "expected int64, but got bool", err.Error())
}

func TestCheck_TaggedFieldName(t *testing.T) {
	tree, err := parser.Parse(`foo.bar`)
	require.NoError(t, err)

	config := &conf.Config{}
	expr.Env(struct {
		x struct {
			y bool `expr:"bar"`
		} `expr:"foo"`
	}{})(config)
	expr.AsBool()(config)

	_, err = checker.Check(tree, config)
	assert.NoError(t, err)
}

func TestCheck_Ambiguous(t *testing.T) {
	type A struct {
		Ambiguous bool
	}
	type B struct {
		Ambiguous int
	}
	type Env struct {
		A
		B
	}

	tree, err := parser.Parse(`Ambiguous == 1`)
	require.NoError(t, err)

	_, err = checker.Check(tree, conf.New(Env{}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ambiguous identifier Ambiguous")
}

func TestCheck_NoConfig(t *testing.T) {
	tree, err := parser.Parse(`any`)
	require.NoError(t, err)

	_, err = checker.Check(tree, nil)
	assert.NoError(t, err)
}

func TestCheck_AllowUndefinedVariables(t *testing.T) {
	type Env struct {
		A int
	}

	tree, err := parser.Parse(`any + fn()`)
	require.NoError(t, err)

	config := conf.New(Env{})
	expr.AllowUndefinedVariables()(config)

	_, err = checker.Check(tree, config)
	assert.NoError(t, err)
}

func TestCheck_AllowUndefinedVariables_DefaultType(t *testing.T) {
	env := map[string]bool{}

	tree, err := parser.Parse(`any`)
	require.NoError(t, err)

	config := conf.New(env)
	expr.AllowUndefinedVariables()(config)
	expr.AsBool()(config)

	_, err = checker.Check(tree, config)
	assert.NoError(t, err)
}

func TestCheck_OperatorOverload(t *testing.T) {
	type Date struct{}
	env := map[string]interface{}{
		"a": Date{},
		"b": Date{},
		"add": func(a, b Date) bool {
			return true
		},
	}
	tree, err := parser.Parse(`a + b`)
	require.NoError(t, err)

	config := conf.New(env)
	expr.AsBool()(config)

	_, err = checker.Check(tree, config)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid operation: + (mismatched types checker_test.Date and checker_test.Date)")

	expr.Operator("+", "add")(config)
	_, err = checker.Check(tree, config)
	require.NoError(t, err)
}

func TestCheck_PointerNode(t *testing.T) {
	_, err := checker.Check(&parser.Tree{Node: &ast.PointerNode{}}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot use pointer accessor outside closure")
}
