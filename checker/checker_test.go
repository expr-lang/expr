package checker_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/checker"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/test/mock"
)

func TestCheck(t *testing.T) {
	var tests = []struct {
		input string
	}{
		{"nil == nil"},
		{"nil == IntPtr"},
		{"nil == nil"},
		{"nil in ArrayOfFoo"},
		{"!Bool"},
		{"!BoolPtr == Bool"},
		{"'a' == 'b' + 'c'"},
		{"'foo' contains 'bar'"},
		{"'foo' endsWith 'bar'"},
		{"'foo' startsWith 'bar'"},
		{"(1 == 1) || (String matches Any)"},
		{"Int % Int > 1"},
		{"Int + Int + Int > 0"},
		{"Int == Any"},
		{"Int in Int..Int"},
		{"IntPtrPtr + 1 > 0"},
		{"1 + 2 + Int64 > 0"},
		{"Int64 % 1 > 0"},
		{"IntPtr == Int"},
		{"FloatPtr == 1 + 2."},
		{"1 + 2 + Float + 3 + 4 < 0"},
		{"1 + Int + Float == 0.5"},
		{"-1 + +1 == 0"},
		{"1 / 2 == 0"},
		{"2**3 + 1 != 0"},
		{"2^3 + 1 != 0"},
		{"Float == 1"},
		{"Float < 1.0"},
		{"Float <= 1.0"},
		{"Float > 1.0"},
		{"Float >= 1.0"},
		{`"a" < "b"`},
		{"String + (true ? String : String) == ''"},
		{"(Any ? nil : '') == ''"},
		{"(Any ? 0 : nil) == 0"},
		{"(Any ? nil : nil) == nil"},
		{"!(Any ? Foo : Foo.Bar).Anything"},
		{"Int in ArrayOfInt"},
		{"Int not in ArrayOfAny"},
		{"String in ArrayOfAny"},
		{"String in ArrayOfString"},
		{"String in Foo"},
		{"String in MapOfFoo"},
		{"String matches 'ok'"},
		{"String matches Any"},
		{"String not matches Any"},
		{"StringPtr == nil"},
		{"[1, 2, 3] == []"},
		{"len([]) > 0"},
		{"Any matches Any"},
		{"!Any.Things.Contains.Any"},
		{"!ArrayOfAny[0].next.goes['any thing']"},
		{"ArrayOfFoo[0].Bar.Baz == ''"},
		{"ArrayOfFoo[0:10][0].Bar.Baz == ''"},
		{"!ArrayOfAny[Any]"},
		{"Bool && Any"},
		{"FuncParam(true, 1, 'str')"},
		{"FuncParamAny(nil)"},
		{"!Fast(Any, String)"},
		{"Foo.Method().Baz == ''"},
		{"Foo.Bar == MapOfAny.id.Bar"},
		{"Foo.Bar.Baz == ''"},
		{"MapOfFoo['any'].Bar.Baz == ''"},
		{"Func() == 0"},
		{"FuncFoo(Foo) > 1"},
		{"Any() > 0"},
		{"Embed.EmbedString == ''"},
		{"EmbedString == ''"},
		{"EmbedMethod(0) == ''"},
		{"Embed.EmbedMethod(0) == ''"},
		{"Embed.EmbedString == ''"},
		{"EmbedString == ''"},
		{"{id: Foo.Bar.Baz, 'str': String} == {}"},
		{"Variadic(0, 1, 2) || Variadic(0)"},
		{"count(1..30, {# % 3 == 0}) > 0"},
		{"map(1..3, {#}) == [1,2,3]"},
		{"map(1..3, #index) == [0,1,2]"},
		{"map(filter(ArrayOfFoo, {.Bar.Baz != ''}), {.Bar}) == []"},
		{"filter(Any, {.AnyMethod()})[0] == ''"},
		{"Time == Time"},
		{"Any == Time"},
		{"Any != Time"},
		{"Any > Time"},
		{"Any >= Time"},
		{"Any < Time"},
		{"Any <= Time"},
		{"Any - Time > Duration"},
		{"Any == Any"},
		{"Any != Any"},
		{"Any > Any"},
		{"Any >= Any"},
		{"Any < Any"},
		{"Any <= Any"},
		{"Any - Any < Duration"},
		{"Time == Any"},
		{"Time != Any"},
		{"Time > Any"},
		{"Time >= Any"},
		{"Time < Any"},
		{"Time <= Any"},
		{"Time - Any == Duration"},
		{"Time + Duration == Time"},
		{"Duration + Time == Time"},
		{"Duration + Any == Time"},
		{"Any + Duration == Time"},
		{"Any.A?.B == nil"},
		{"(Any.Bool ?? Bool) > 0"},
		{"Bool ?? Bool"},
		{"let foo = 1; foo == 1"},
		{"(Embed).EmbedPointerEmbedInt > 0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var err error
			tree, err := parser.Parse(tt.input)
			require.NoError(t, err)

			config := conf.New(mock.Env{})
			expr.AsBool()(config)

			_, err = checker.Check(tree, config)
			assert.NoError(t, err)
		})
	}
}

const errorTests = `
Any > Foo
invalid operation: > (mismatched types interface {} and mock.Foo) (1:5)
 | Any > Foo
 | ....^
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

		_, err = checker.Check(tree, conf.New(mock.Env{}))
		if err == nil {
			err = fmt.Errorf("<nil>")
		}

		assert.Equal(t, input[1], err.Error(), input[0])
	}
}

func TestCheck_FloatVsInt(t *testing.T) {
	tree, err := parser.Parse(`Int + Float`)
	require.NoError(t, err)

	typ, err := checker.Check(tree, conf.New(mock.Env{}))
	assert.NoError(t, err)
	assert.Equal(t, typ.Kind(), reflect.Float64)
}

func TestCheck_IntSums(t *testing.T) {
	tree, err := parser.Parse(`Uint32 + Int32`)
	require.NoError(t, err)

	typ, err := checker.Check(tree, conf.New(mock.Env{}))
	assert.NoError(t, err)
	assert.Equal(t, typ.Kind(), reflect.Int)
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

	_, err = checker.Check(tree, conf.CreateNew())
	assert.NoError(t, err)
}

func TestCheck_AllowUndefinedVariables(t *testing.T) {
	type Env struct {
		A int
	}

	tree, err := parser.Parse(`Any + fn()`)
	require.NoError(t, err)

	config := conf.New(Env{})
	expr.AllowUndefinedVariables()(config)

	_, err = checker.Check(tree, config)
	assert.NoError(t, err)
}

func TestCheck_AllowUndefinedVariables_DefaultType(t *testing.T) {
	env := map[string]bool{}

	tree, err := parser.Parse(`Any`)
	require.NoError(t, err)

	config := conf.New(env)
	expr.AllowUndefinedVariables()(config)
	expr.AsBool()(config)

	_, err = checker.Check(tree, config)
	assert.NoError(t, err)
}

func TestCheck_AllowUndefinedVariables_OptionalChaining(t *testing.T) {
	type Env struct{}

	tree, err := parser.Parse("Not?.A.B == nil")
	require.NoError(t, err)

	config := conf.New(Env{})
	expr.AllowUndefinedVariables()(config)

	_, err = checker.Check(tree, config)
	assert.NoError(t, err)
}

func TestCheck_PointerNode(t *testing.T) {
	_, err := checker.Check(&parser.Tree{Node: &ast.PointerNode{}}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot use pointer accessor outside closure")
}

func TestCheck_TypeWeights(t *testing.T) {
	types := map[string]any{
		"Uint":    uint(1),
		"Uint8":   uint8(2),
		"Uint16":  uint16(3),
		"Uint32":  uint32(4),
		"Uint64":  uint64(5),
		"Int":     6,
		"Int8":    int8(7),
		"Int16":   int16(8),
		"Int32":   int32(9),
		"Int64":   int64(10),
		"Float32": float32(11),
		"Float64": float64(12),
	}
	for a := range types {
		for b := range types {
			tree, err := parser.Parse(fmt.Sprintf("%s + %s", a, b))
			require.NoError(t, err)

			config := conf.New(types)

			_, err = checker.Check(tree, config)
			require.NoError(t, err)
		}
	}
}

func TestCheck_works_with_nil_types(t *testing.T) {
	env := map[string]any{
		"null": nil,
	}

	tree, err := parser.Parse("null")
	require.NoError(t, err)

	_, err = checker.Check(tree, conf.New(env))
	require.NoError(t, err)
}

func TestCheck_cast_to_expected_works_with_interface(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		type Env struct {
			Any any
		}

		tree, err := parser.Parse("Any")
		require.NoError(t, err)

		config := conf.New(Env{})
		expr.AsFloat64()(config)
		expr.AsAny()(config)

		_, err = checker.Check(tree, config)
		require.NoError(t, err)
	})

	t.Run("kind", func(t *testing.T) {
		env := map[string]any{
			"Any": any("foo"),
		}

		tree, err := parser.Parse("Any")
		require.NoError(t, err)

		config := conf.New(env)
		expr.AsKind(reflect.String)(config)

		_, err = checker.Check(tree, config)
		require.NoError(t, err)
	})
}

func TestCheck_operator_in_works_with_interfaces(t *testing.T) {
	tree, err := parser.Parse(`'Tom' in names`)
	require.NoError(t, err)

	config := conf.New(nil)
	expr.AllowUndefinedVariables()(config)

	_, err = checker.Check(tree, config)
	require.NoError(t, err)
}

func TestCheck_Function_types_are_checked(t *testing.T) {
	add := expr.Function(
		"add",
		func(p ...any) (any, error) {
			out := 0
			for _, each := range p {
				out += each.(int)
			}
			return out, nil
		},
		new(func(int) int),
		new(func(int, int) int),
		new(func(int, int, int) int),
		new(func(...int) int),
	)

	config := conf.CreateNew()
	add(config)

	tests := []string{
		"add(1)",
		"add(1, 2)",
		"add(1, 2, 3)",
		"add(1, 2, 3, 4)",
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			tree, err := parser.Parse(test)
			require.NoError(t, err)

			_, err = checker.Check(tree, config)
			require.NoError(t, err)
			require.Equal(t, reflect.Int, tree.Node.Type().Kind())
		})
	}

	t.Run("errors", func(t *testing.T) {
		tree, err := parser.Parse("add(1, '2')")
		require.NoError(t, err)

		_, err = checker.Check(tree, config)
		require.Error(t, err)
		require.Equal(t, "cannot use string as argument (type int) to call add  (1:8)\n | add(1, '2')\n | .......^", err.Error())
	})
}

func TestCheck_Function_without_types(t *testing.T) {
	add := expr.Function(
		"add",
		func(p ...any) (any, error) {
			out := 0
			for _, each := range p {
				out += each.(int)
			}
			return out, nil
		},
	)

	tree, err := parser.Parse("add(1, 2, 3)")
	require.NoError(t, err)

	config := conf.CreateNew()
	add(config)

	_, err = checker.Check(tree, config)
	require.NoError(t, err)
	require.Equal(t, reflect.Interface, tree.Node.Type().Kind())
}

func TestCheck_dont_panic_on_nil_arguments_for_builtins(t *testing.T) {
	tests := []string{
		"len(nil)",
		"abs(nil)",
		"int(nil)",
		"float(nil)",
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			tree, err := parser.Parse(test)
			require.NoError(t, err)

			_, err = checker.Check(tree, conf.New(nil))
			require.Error(t, err)
		})
	}
}

func TestCheck_do_not_override_params_for_functions(t *testing.T) {
	env := map[string]any{
		"foo": func(p string) string {
			return "foo"
		},
	}
	config := conf.New(env)
	expr.Function(
		"bar",
		func(p ...any) (any, error) {
			return p[0].(string), nil
		},
		new(func(string) string),
	)(config)
	config.Check()

	t.Run("func from env", func(t *testing.T) {
		tree, err := parser.Parse("foo(1)")
		require.NoError(t, err)

		_, err = checker.Check(tree, config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot use int as argument")
	})

	t.Run("func from function", func(t *testing.T) {
		tree, err := parser.Parse("bar(1)")
		require.NoError(t, err)

		_, err = checker.Check(tree, config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot use int as argument")
	})
}

func TestCheck_env_keyword(t *testing.T) {
	env := map[string]any{
		"num":  42,
		"str":  "foo",
		"name": "str",
	}

	tests := []struct {
		input string
		want  reflect.Kind
	}{
		{`$env['str']`, reflect.String},
		{`$env['num']`, reflect.Int},
		{`$env[name]`, reflect.Interface},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tree, err := parser.Parse(test.input)
			require.NoError(t, err)

			rtype, err := checker.Check(tree, conf.New(env))
			require.NoError(t, err)
			require.True(t, rtype.Kind() == test.want, fmt.Sprintf("expected %s, got %s", test.want, rtype.Kind()))
		})
	}
}

func TestCheck_builtin_without_call(t *testing.T) {
	tests := []struct {
		input string
		err   string
	}{
		{`len + 1`, "invalid operation: + (mismatched types func(...interface {}) (interface {}, error) and int) (1:5)\n | len + 1\n | ....^"},
		{`string.A`, "type func(interface {}) string[string] is undefined (1:8)\n | string.A\n | .......^"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tree, err := parser.Parse(test.input)
			require.NoError(t, err)

			_, err = checker.Check(tree, conf.New(nil))
			require.Error(t, err)
			require.Equal(t, test.err, err.Error())
		})
	}
}
