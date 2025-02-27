package checker_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"expr/internal/testify/assert"
	"expr/internal/testify/require"
	"expr/types"

	"expr"
	"expr/ast"
	"expr/checker"
	"expr/conf"
	"expr/parser"
)

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

	config := conf.CreateNew()
	expr.Env(struct {
		x struct {
			y bool `expr:"bar"`
		} `expr:"foo"`
	}{})(config)
	expr.AsBool()(config)

	_, err = checker.Check(tree, config)
	assert.NoError(t, err)
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

func TestCheck_types(t *testing.T) {
	env := types.Map{
		"foo": types.Map{
			"bar": types.Map{
				"baz":       types.String,
				types.Extra: types.String,
			},
		},
		"arr": types.Array(types.Map{
			"value": types.String,
		}),
		types.Extra: types.Any,
	}

	noerr := "no error"
	tests := []struct {
		code string
		err  string
	}{
		{`unknown`, noerr},
		{`[unknown + 42, another_unknown + "foo"]`, noerr},
		{`foo.bar.baz > 0`, `invalid operation: > (mismatched types string and int)`},
		{`foo.unknown.baz`, `unknown field unknown (1:5)`},
		{`foo.bar.unknown`, noerr},
		{`foo.bar.unknown + 42`, `invalid operation: + (mismatched types string and int)`},
		{`[foo] | map(.unknown)`, `unknown field unknown`},
		{`[foo] | map(.bar) | filter(.baz)`, `predicate should return boolean (got string)`},
		{`arr | filter(.value > 0)`, `invalid operation: > (mismatched types string and int)`},
		{`arr | filter(.value contains "a") | filter(.value == 0)`, `invalid operation: == (mismatched types string and int)`},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			tree, err := parser.Parse(test.code)
			require.NoError(t, err)

			config := conf.New(env)
			_, err = checker.Check(tree, config)
			if test.err == noerr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
