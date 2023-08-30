package builtin_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/test/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuiltin(t *testing.T) {
	env := map[string]any{
		"ArrayOfString": []string{"foo", "bar", "baz"},
		"ArrayOfInt":    []int{1, 2, 3},
		"ArrayOfAny":    []any{1, "2", true},
		"ArrayOfFoo":    []mock.Foo{{Value: "a"}, {Value: "b"}, {Value: "c"}},
	}

	var tests = []struct {
		input string
		want  any
	}{
		{`len(1..10)`, 10},
		{`len({foo: 1, bar: 2})`, 2},
		{`len("hello")`, 5},
		{`abs(-5)`, 5},
		{`abs(.5)`, .5},
		{`abs(-.5)`, .5},
		{`int(5.5)`, 5},
		{`int(5)`, 5},
		{`int("5")`, 5},
		{`float(5)`, 5.0},
		{`float(5.5)`, 5.5},
		{`float("5.5")`, 5.5},
		{`string(5)`, "5"},
		{`string(5.5)`, "5.5"},
		{`string("5.5")`, "5.5"},
		{`trim("  foo  ")`, "foo"},
		{`trim("__foo___", "_")`, "foo"},
		{`trimPrefix("prefix_foo", "prefix_")`, "foo"},
		{`trimSuffix("foo_suffix", "_suffix")`, "foo"},
		{`upper("foo")`, "FOO"},
		{`lower("FOO")`, "foo"},
		{`split("foo,bar,baz", ",")`, []string{"foo", "bar", "baz"}},
		{`split("foo,bar,baz", ",", 2)`, []string{"foo", "bar,baz"}},
		{`splitAfter("foo,bar,baz", ",")`, []string{"foo,", "bar,", "baz"}},
		{`splitAfter("foo,bar,baz", ",", 2)`, []string{"foo,", "bar,baz"}},
		{`replace("foo,bar,baz", ",", ";")`, "foo;bar;baz"},
		{`replace("foo,bar,baz,goo", ",", ";", 2)`, "foo;bar;baz,goo"},
		{`repeat("foo", 3)`, "foofoofoo"},
		{`join(ArrayOfString, ",")`, "foo,bar,baz"},
		{`join(ArrayOfString)`, "foobarbaz"},
		{`join(["foo", "bar", "baz"], ",")`, "foo,bar,baz"},
		{`join(["foo", "bar", "baz"])`, "foobarbaz"},
		{`indexOf("foo,bar,baz", ",")`, 3},
		{`lastIndexOf("foo,bar,baz", ",")`, 7},
		{`hasPrefix("foo,bar,baz", "foo")`, true},
		{`hasSuffix("foo,bar,baz", "baz")`, true},
		{`max(1, 2, 3)`, 3},
		{`max(1.5, 2.5, 3.5)`, 3.5},
		{`min(1, 2, 3)`, 1},
		{`min(1.5, 2.5, 3.5)`, 1.5},
		{`sum(1..9)`, 45},
		{`sum([.5, 1.5, 2.5])`, 4.5},
		{`sum([])`, 0},
		{`sum([1, 2, 3.0, 4])`, 10.0},
		{`mean(1..9)`, 5.0},
		{`mean([.5, 1.5, 2.5])`, 1.5},
		{`mean([])`, 0.0},
		{`mean([1, 2, 3.0, 4])`, 2.5},
		{`median(1..9)`, 5.0},
		{`median([.5, 1.5, 2.5])`, 1.5},
		{`median([])`, 0.0},
		{`median([1, 2, 3])`, 2.0},
		{`median([1, 2, 3, 4])`, 2.5},
		{`toJSON({foo: 1, bar: 2})`, "{\n  \"bar\": 2,\n  \"foo\": 1\n}"},
		{`fromJSON("[1, 2, 3]")`, []any{1.0, 2.0, 3.0}},
		{`toBase64("hello")`, "aGVsbG8="},
		{`fromBase64("aGVsbG8=")`, "hello"},
		{`now().Format("2006-01-02T15:04Z")`, time.Now().Format("2006-01-02T15:04Z")},
		{`duration("1h")`, time.Hour},
		{`date("2006-01-02T15:04:05Z")`, time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)},
		{`date("2006.01.02", "2006.01.02")`, time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC)},
		{`first(ArrayOfString)`, "foo"},
		{`first(ArrayOfInt)`, 1},
		{`first(ArrayOfAny)`, 1},
		{`first([])`, nil},
		{`last(ArrayOfString)`, "baz"},
		{`last(ArrayOfInt)`, 3},
		{`last(ArrayOfAny)`, true},
		{`last([])`, nil},
		{`get(ArrayOfString, 1)`, "bar"},
		{`get(ArrayOfString, 99)`, nil},
		{`get(ArrayOfInt, 1)`, 2},
		{`get(ArrayOfInt, -1)`, 3},
		{`get(ArrayOfAny, 1)`, "2"},
		{`get({foo: 1, bar: 2}, "foo")`, 1},
		{`get({foo: 1, bar: 2}, "unknown")`, nil},
		{`take(ArrayOfString, 2)`, []string{"foo", "bar"}},
		{`take(ArrayOfString, 99)`, []string{"foo", "bar", "baz"}},
		{`"foo" in keys({foo: 1, bar: 2})`, true},
		{`1 in values({foo: 1, bar: 2})`, true},
		{`len(toPairs({foo: 1, bar: 2}))`, 2},
		{`len(toPairs({}))`, 0},
		{`fromPairs([["foo", 1], ["bar", 2]])`, map[any]any{"foo": 1, "bar": 2}},
		{`fromPairs(toPairs({foo: 1, bar: 2}))`, map[any]any{"foo": 1, "bar": 2}},
		{`groupBy(1..9, # % 2)`, map[any][]any{0: {2, 4, 6, 8}, 1: {1, 3, 5, 7, 9}}},
		{`groupBy(1..9, # % 2)[0]`, []any{2, 4, 6, 8}},
		{`groupBy(1..3, # > 1)[true]`, []any{2, 3}},
		{`groupBy(1..3, # > 1 ? nil : "")[nil]`, []any{2, 3}},
		{`groupBy(ArrayOfFoo, .Value).a`, []any{mock.Foo{Value: "a"}}},
		{`reduce(1..9, # + #acc, 0)`, 45},
		{`reduce(1..9, # + #acc)`, 45},
		{`reduce([.5, 1.5, 2.5], # + #acc, 0)`, 4.5},
		{`reduce([], 5, 0)`, 0},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program, err := expr.Compile(test.input, expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, test.want, out)
		})
	}
}

func TestBuiltin_works_with_any(t *testing.T) {
	config := map[string]struct {
		arity int
	}{
		"get":    {2},
		"take":   {2},
		"sortBy": {2},
	}

	for _, b := range builtin.Builtins {
		if b.Predicate {
			continue
		}
		t.Run(b.Name, func(t *testing.T) {
			arity := 1
			if c, ok := config[b.Name]; ok {
				arity = c.arity
			}
			if len(b.Types) > 0 {
				arity = b.Types[0].NumIn()
			}
			args := make([]string, arity)
			for i := 1; i <= arity; i++ {
				args[i-1] = fmt.Sprintf("arg%d", i)
			}
			_, err := expr.Compile(fmt.Sprintf(`%s(%s)`, b.Name, strings.Join(args, ", "))) // expr.Env(env) is not needed
			assert.NoError(t, err)
		})
	}
}

func TestBuiltin_errors(t *testing.T) {
	var errorTests = []struct {
		input string
		err   string
	}{
		{`len()`, `invalid number of arguments (expected 1, got 0)`},
		{`len(1)`, `invalid argument for len (type int)`},
		{`abs()`, `invalid number of arguments (expected 1, got 0)`},
		{`abs(1, 2)`, `invalid number of arguments (expected 1, got 2)`},
		{`abs("foo")`, `invalid argument for abs (type string)`},
		{`int()`, `invalid number of arguments (expected 1, got 0)`},
		{`int(1, 2)`, `invalid number of arguments (expected 1, got 2)`},
		{`float()`, `invalid number of arguments (expected 1, got 0)`},
		{`float(1, 2)`, `invalid number of arguments (expected 1, got 2)`},
		{`string(1, 2)`, `too many arguments to call string`},
		{`trim()`, `not enough arguments to call trim`},
		{`max()`, `not enough arguments to call max`},
		{`max(1, "2")`, `invalid argument for max (type string)`},
		{`min()`, `not enough arguments to call min`},
		{`min(1, "2")`, `invalid argument for min (type string)`},
		{`duration("error")`, `invalid duration`},
		{`date("error")`, `invalid date`},
		{`get()`, `invalid number of arguments (expected 2, got 0)`},
		{`get(1, 2)`, `type int does not support indexing`},
	}
	for _, test := range errorTests {
		t.Run(test.input, func(t *testing.T) {
			_, err := expr.Eval(test.input, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.err)
		})
	}
}

func TestBuiltin_types(t *testing.T) {
	env := map[string]any{
		"num":           42,
		"str":           "foo",
		"ArrayOfString": []string{"foo", "bar", "baz"},
		"ArrayOfInt":    []int{1, 2, 3},
	}

	tests := []struct {
		input string
		want  reflect.Kind
	}{
		{`get(ArrayOfString, 0)`, reflect.String},
		{`get(ArrayOfInt, 0)`, reflect.Int},
		{`first(ArrayOfString)`, reflect.String},
		{`first(ArrayOfInt)`, reflect.Int},
		{`last(ArrayOfString)`, reflect.String},
		{`last(ArrayOfInt)`, reflect.Int},
		{`get($env, 'str')`, reflect.String},
		{`get($env, 'num')`, reflect.Int},
		{`get($env, 'ArrayOfString')`, reflect.Slice},
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

func TestBuiltin_memory_limits(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`repeat("\xc4<\xc4\xc4\xc4",10009999990)`},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			_, err := expr.Eval(test.input, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "memory budget exceeded")
		})
	}
}

func TestBuiltin_disallow_builtins_override(t *testing.T) {
	t.Run("via env", func(t *testing.T) {
		env := map[string]any{
			"len": func() int { return 42 },
			"repeat": func(a string) string {
				return a
			},
		}
		assert.Panics(t, func() {
			_, _ = expr.Compile(`string(len("foo")) + repeat("0", 2)`, expr.Env(env))
		})
	})
	t.Run("via expr.Function", func(t *testing.T) {
		length := expr.Function("len",
			func(params ...any) (any, error) {
				return 42, nil
			},
			new(func() int),
		)
		repeat := expr.Function("repeat",
			func(params ...any) (any, error) {
				return params[0], nil
			},
			new(func(string) string),
		)
		assert.Panics(t, func() {
			_, _ = expr.Compile(`string(len("foo")) + repeat("0", 2)`, length, repeat)
		})
	})
}

func TestBuiltin_DisableBuiltin(t *testing.T) {
	t.Run("via env", func(t *testing.T) {
		for _, b := range builtin.Builtins {
			if b.Predicate {
				continue // TODO: allow to disable predicates
			}
			t.Run(b.Name, func(t *testing.T) {
				env := map[string]any{
					b.Name: func() int { return 42 },
				}
				program, err := expr.Compile(b.Name+"()", expr.Env(env), expr.DisableBuiltin(b.Name))
				require.NoError(t, err)

				out, err := expr.Run(program, env)
				require.NoError(t, err)
				assert.Equal(t, 42, out)
			})
		}
	})
	t.Run("via expr.Function", func(t *testing.T) {
		for _, b := range builtin.Builtins {
			if b.Predicate {
				continue // TODO: allow to disable predicates
			}
			t.Run(b.Name, func(t *testing.T) {
				fn := expr.Function(b.Name,
					func(params ...any) (any, error) {
						return 42, nil
					},
					new(func() int),
				)
				program, err := expr.Compile(b.Name+"()", fn, expr.DisableBuiltin(b.Name))
				require.NoError(t, err)

				out, err := expr.Run(program, nil)
				require.NoError(t, err)
				assert.Equal(t, 42, out)
			})
		}
	})
}

func TestBuiltin_DisableAllBuiltins(t *testing.T) {
	_, err := expr.Compile(`len("foo")`, expr.Env(nil), expr.DisableAllBuiltins())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown name len")
}

func TestBuiltin_EnableBuiltin(t *testing.T) {
	t.Run("via env", func(t *testing.T) {
		env := map[string]any{
			"repeat": func() string { return "repeat" },
		}
		program, err := expr.Compile(`len(repeat())`, expr.Env(env), expr.DisableAllBuiltins(), expr.EnableBuiltin("len"))
		require.NoError(t, err)

		out, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, 6, out)
	})
	t.Run("via expr.Function", func(t *testing.T) {
		fn := expr.Function("repeat",
			func(params ...any) (any, error) {
				return "repeat", nil
			},
			new(func() string),
		)
		program, err := expr.Compile(`len(repeat())`, fn, expr.DisableAllBuiltins(), expr.EnableBuiltin("len"))
		require.NoError(t, err)

		out, err := expr.Run(program, nil)
		require.NoError(t, err)
		assert.Equal(t, 6, out)
	})
}

func TestBuiltin_type(t *testing.T) {
	type Foo struct{}
	var b any = 1
	var a any = &b
	tests := []struct {
		obj  any
		want string
	}{
		{nil, "nil"},
		{true, "bool"},
		{1, "int"},
		{int8(1), "int"},
		{uint(1), "uint"},
		{1.0, "float"},
		{float32(1.0), "float"},
		{"string", "string"},
		{[]string{"foo", "bar"}, "array"},
		{map[string]any{"foo": "bar"}, "map"},
		{func() {}, "func"},
		{time.Now(), "time.Time"},
		{time.Second, "time.Duration"},
		{Foo{}, "github.com/antonmedv/expr/builtin_test.Foo"},
		{struct{}{}, "struct"},
		{a, "int"},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			env := map[string]any{
				"obj": test.obj,
			}
			program, err := expr.Compile(`type(obj)`, expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, test.want, out)
		})
	}
}

func TestBuiltin_sort(t *testing.T) {
	env := map[string]any{
		"ArrayOfString": []string{"foo", "bar", "baz"},
		"ArrayOfInt":    []int{3, 2, 1},
		"ArrayOfFoo":    []mock.Foo{{Value: "c"}, {Value: "a"}, {Value: "b"}},
	}
	tests := []struct {
		input string
		want  any
	}{
		{`sort([])`, []any{}},
		{`sort(ArrayOfInt)`, []any{1, 2, 3}},
		{`sort(ArrayOfInt, 'desc')`, []any{3, 2, 1}},
		{`sortBy(ArrayOfFoo, 'Value')`, []any{mock.Foo{Value: "a"}, mock.Foo{Value: "b"}, mock.Foo{Value: "c"}}},
		{`sortBy([{id: "a"}, {id: "b"}], "id", "desc")`, []any{map[string]any{"id": "b"}, map[string]any{"id": "a"}}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program, err := expr.Compile(test.input, expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, test.want, out)
		})
	}
}
