package operator_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/test/mock"
)

func TestOperator_struct(t *testing.T) {
	env := mock.Env{
		Time: time.Date(2017, time.October, 23, 18, 30, 0, 0, time.UTC),
	}

	code := `Time == "2017-10-23"`

	program, err := expr.Compile(code, expr.Env(mock.Env{}), expr.Operator("==", "TimeEqualString"))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestOperator_options_another_order(t *testing.T) {
	code := `Time == "2017-10-23"`
	_, err := expr.Compile(code, expr.Operator("==", "TimeEqualString"), expr.Env(mock.Env{}))
	require.NoError(t, err)
}

func TestOperator_no_env(t *testing.T) {
	code := `Time == "2017-10-23"`
	require.Panics(t, func() {
		_, _ = expr.Compile(code, expr.Operator("==", "TimeEqualString"))
	})
}

func TestOperator_interface(t *testing.T) {
	env := mock.Env{}

	code := `Foo == "Foo.String" && "Foo.String" == Foo && Time != Foo && Time == Time`

	program, err := expr.Compile(
		code,
		expr.Env(mock.Env{}),
		expr.Operator("==", "StringerStringEqual", "StringStringerEqual", "StringerStringerEqual"),
		expr.Operator("!=", "NotStringerStringEqual", "NotStringStringerEqual", "NotStringerStringerEqual"),
	)
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

type Value struct {
	Int int
}

func TestOperator_Function(t *testing.T) {
	env := map[string]interface{}{
		"foo": Value{1},
		"bar": Value{2},
	}

	tests := []struct {
		input string
		want  int
	}{
		{
			input: `foo + bar`,
			want:  3,
		},
		{
			input: `2 + 4`,
			want:  6,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(`opertor function helper test %s`, tt.input), func(t *testing.T) {
			program, err := expr.Compile(
				tt.input,
				expr.Env(env),
				expr.Operator("+", "Add", "AddInt"),
				expr.Function("Add", func(args ...interface{}) (interface{}, error) {
					return args[0].(Value).Int + args[1].(Value).Int, nil
				},
					new(func(_ Value, __ Value) int),
				),
				expr.Function("AddInt", func(args ...interface{}) (interface{}, error) {
					return args[0].(int) + args[1].(int), nil
				},
					new(func(_ int, __ int) int),
				),
			)
			require.NoError(t, err)

			output, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, tt.want, output)
		})
	}

}

func TestOperator_Function_WithTypes(t *testing.T) {
	env := map[string]interface{}{
		"foo": Value{1},
		"bar": Value{2},
	}

	require.PanicsWithError(t, `function Add for + operator misses types`, func() {
		_, _ = expr.Compile(
			`foo + bar`,
			expr.Env(env),
			expr.Operator("+", "Add", "AddInt"),
			expr.Function("Add", func(args ...interface{}) (interface{}, error) {
				return args[0].(Value).Int + args[1].(Value).Int, nil
			}),
		)
	})

	require.PanicsWithError(t, `function Add for + operator does not have a correct signature`, func() {
		_, _ = expr.Compile(
			`foo + bar`,
			expr.Env(env),
			expr.Operator("+", "Add", "AddInt"),
			expr.Function("Add", func(args ...interface{}) (interface{}, error) {
				return args[0].(Value).Int + args[1].(Value).Int, nil
			},
				new(func(_ Value) int),
			),
		)
	})

}

func TestOperator_FunctionOverTypesPrecedence(t *testing.T) {
	env := struct {
		Add func(a, b int) int
	}{
		Add: func(a, b int) int {
			return a + b
		},
	}

	program, err := expr.Compile(
		`1 + 2`,
		expr.Env(env),
		expr.Operator("+", "Add"),
		expr.Function("Add", func(args ...interface{}) (interface{}, error) {
			// Wierd function that returns 100 + a + b in testing purposes.
			return args[0].(int) + args[1].(int) + 100, nil
		},
			new(func(_ int, __ int) int),
		),
	)
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 103, output)
}

func TestOperator_CanBeDefinedEitherInTypesOrInFunctions(t *testing.T) {
	env := struct {
		Add func(a, b int) int
	}{
		Add: func(a, b int) int {
			return a + b
		},
	}

	program, err := expr.Compile(
		`1 + 2`,
		expr.Env(env),
		expr.Operator("+", "Add", "AddValues"),
		expr.Function("AddValues", func(args ...interface{}) (interface{}, error) {
			return args[0].(Value).Int + args[1].(Value).Int, nil
		},
			new(func(_ Value, __ Value) int),
		),
	)
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 3, output)
}
