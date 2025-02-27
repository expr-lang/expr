package optimizer_test

import (
	"strings"
	"testing"

	"expr/internal/testify/assert"
	"expr/internal/testify/require"

	"expr"
	"expr/ast"
	"expr/conf"
	"expr/optimizer"
	"expr/parser"
)

func TestOptimize(t *testing.T) {
	env := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	tests := []struct {
		expr string
		want any
	}{
		{`1 + 2`, 3},
		{`sum([])`, 0},
		{`sum([a])`, 1},
		{`sum([a, b])`, 3},
		{`sum([a, b, c])`, 6},
		{`sum([a, b, c, 4])`, 10},
		{`sum(1..10, # * 1000)`, 55000},
		{`sum(map(1..10, # * 1000), # / 1000)`, float64(55)},
		{`all(1..3, {# > 0}) && all(1..3, {# < 4})`, true},
		{`all(1..3, {# > 2}) && all(1..3, {# < 4})`, false},
		{`all(1..3, {# > 0}) && all(1..3, {# < 2})`, false},
		{`all(1..3, {# > 2}) && all(1..3, {# < 2})`, false},
		{`all(1..3, {# > 0}) || all(1..3, {# < 4})`, true},
		{`all(1..3, {# > 0}) || all(1..3, {# != 2})`, true},
		{`all(1..3, {# != 3}) || all(1..3, {# < 4})`, true},
		{`all(1..3, {# != 3}) || all(1..3, {# != 2})`, false},
		{`none(1..3, {# == 0})`, true},
		{`none(1..3, {# == 0}) && none(1..3, {# == 4})`, true},
		{`none(1..3, {# == 0}) && none(1..3, {# == 3})`, false},
		{`none(1..3, {# == 1}) && none(1..3, {# == 4})`, false},
		{`none(1..3, {# == 1}) && none(1..3, {# == 3})`, false},
		{`none(1..3, {# == 0}) || none(1..3, {# == 4})`, true},
		{`none(1..3, {# == 0}) || none(1..3, {# == 3})`, true},
		{`none(1..3, {# == 1}) || none(1..3, {# == 4})`, true},
		{`none(1..3, {# == 1}) || none(1..3, {# == 3})`, false},
		{`any([1, 1, 0, 1], {# == 0})`, true},
		{`any(1..3, {# == 1}) && any(1..3, {# == 2})`, true},
		{`any(1..3, {# == 0}) && any(1..3, {# == 2})`, false},
		{`any(1..3, {# == 1}) && any(1..3, {# == 4})`, false},
		{`any(1..3, {# == 0}) && any(1..3, {# == 4})`, false},
		{`any(1..3, {# == 1}) || any(1..3, {# == 2})`, true},
		{`any(1..3, {# == 0}) || any(1..3, {# == 2})`, true},
		{`any(1..3, {# == 1}) || any(1..3, {# == 4})`, true},
		{`any(1..3, {# == 0}) || any(1..3, {# == 4})`, false},
		{`one([1, 1, 0, 1], {# == 0}) and not one([1, 0, 0, 1], {# == 0})`, true},
		{`one(1..3, {# == 1}) and one(1..3, {# == 2})`, true},
		{`one(1..3, {# == 1 || # == 2}) and one(1..3, {# == 2})`, false},
		{`one(1..3, {# == 1}) and one(1..3, {# == 2 || # == 3})`, false},
		{`one(1..3, {# == 1 || # == 2}) and one(1..3, {# == 2 || # == 3})`, false},
		{`one(1..3, {# == 1}) or one(1..3, {# == 2})`, true},
		{`one(1..3, {# == 1 || # == 2}) or one(1..3, {# == 2})`, true},
		{`one(1..3, {# == 1}) or one(1..3, {# == 2 || # == 3})`, true},
		{`one(1..3, {# == 1 || # == 2}) or one(1..3, {# == 2 || # == 3})`, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.Env(env))
			require.NoError(t, err)

			output, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)

			unoptimizedProgram, err := expr.Compile(tt.expr, expr.Env(env), expr.Optimize(false))
			require.NoError(t, err)

			unoptimizedOutput, err := expr.Run(unoptimizedProgram, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, unoptimizedOutput)
		})
	}
}

func TestOptimize_in_range_with_floats(t *testing.T) {
	out, err := expr.Eval(`f in 1..3`, map[string]any{"f": 1.5})
	require.NoError(t, err)
	assert.Equal(t, false, out)
}

func TestOptimize_const_expr(t *testing.T) {
	tree, err := parser.Parse(`toUpper("hello")`)
	require.NoError(t, err)

	env := map[string]any{
		"toUpper": strings.ToUpper,
	}

	config := conf.New(env)
	config.ConstExpr("toUpper")

	err = optimizer.Optimize(&tree.Node, config)
	require.NoError(t, err)

	expected := &ast.ConstantNode{
		Value: "HELLO",
	}

	assert.Equal(t, ast.Dump(expected), ast.Dump(tree.Node))
}
