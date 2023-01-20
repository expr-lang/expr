package expr_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleEval() {
	output, err := expr.Eval("greet + name", map[string]interface{}{
		"greet": "Hello, ",
		"name":  "world!",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: Hello, world!
}

func ExampleEval_runtime_error() {
	_, err := expr.Eval(`map(1..3, {1 % (# - 3)})`, nil)
	fmt.Print(err)

	// Output: runtime error: integer divide by zero (1:14)
	//  | map(1..3, {1 % (# - 3)})
	//  | .............^
}

func ExampleCompile() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 99,
	}

	program, err := expr.Compile("foo in 1..99 and bar in 1..99", expr.Env(env))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEnv() {
	type Segment struct {
		Origin string
	}
	type Passengers struct {
		Adults int
	}
	type Meta struct {
		Tags map[string]string
	}
	type Env struct {
		Meta
		Segments   []*Segment
		Passengers *Passengers
		Marker     string
	}

	code := `all(Segments, {.Origin == "MOW"}) && Passengers.Adults > 0 && Tags["foo"] startsWith "bar"`

	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := Env{
		Meta: Meta{
			Tags: map[string]string{
				"foo": "bar",
			},
		},
		Segments: []*Segment{
			{Origin: "MOW"},
		},
		Passengers: &Passengers{
			Adults: 2,
		},
		Marker: "test",
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEnv_tagged_field_names() {
	env := struct {
		FirstWord  string
		Separator  string `expr:"Space"`
		SecondWord string `expr:"second_word"`
	}{
		FirstWord:  "Hello",
		Separator:  " ",
		SecondWord: "World",
	}

	output, err := expr.Eval(`FirstWord + Space + second_word`, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output : Hello World
}

func ExampleAsKind() {
	program, err := expr.Compile("{a: 1, b: 2}", expr.AsKind(reflect.Map))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: map[a:1 b:2]
}

func ExampleAsBool() {
	env := map[string]int{
		"foo": 0,
	}

	program, err := expr.Compile("foo >= 0", expr.Env(env), expr.AsBool())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output.(bool))

	// Output: true
}

func ExampleAsBool_error() {
	env := map[string]interface{}{
		"foo": 0,
	}

	_, err := expr.Compile("foo + 42", expr.Env(env), expr.AsBool())

	fmt.Printf("%v", err)

	// Output: expected bool, but got int
}

func ExampleAsInt() {
	program, err := expr.Compile("42", expr.AsInt())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%T(%v)", output, output)

	// Output: int(42)
}

func ExampleAsInt64() {
	env := map[string]interface{}{
		"rating": 5.5,
	}

	program, err := expr.Compile("rating", expr.Env(env), expr.AsInt64())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output.(int64))

	// Output: 5
}

func ExampleAsFloat64() {
	program, err := expr.Compile("42", expr.AsFloat64())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output.(float64))

	// Output: 42
}

func ExampleAsFloat64_error() {
	_, err := expr.Compile(`!!true`, expr.AsFloat64())

	fmt.Printf("%v", err)

	// Output: expected float64, but got bool
}

func ExampleOperator() {
	code := `
		Now() > CreatedAt &&
		(Now() - CreatedAt).Hours() > 24
	`

	type Env struct {
		CreatedAt time.Time
		Now       func() time.Time
		Sub       func(a, b time.Time) time.Duration
		After     func(a, b time.Time) bool
	}

	options := []expr.Option{
		expr.Env(Env{}),
		expr.Operator(">", "After"),
		expr.Operator("-", "Sub"),
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := Env{
		CreatedAt: time.Date(2018, 7, 14, 0, 0, 0, 0, time.UTC),
		Now:       func() time.Time { return time.Now() },
		Sub:       func(a, b time.Time) time.Duration { return a.Sub(b) },
		After:     func(a, b time.Time) bool { return a.After(b) },
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func ExampleConstExpr() {
	code := `[fib(5), fib(3+3), fib(dyn)]`

	env := map[string]interface{}{
		"fib": fib,
		"dyn": 0,
	}

	options := []expr.Option{
		expr.Env(env),
		expr.ConstExpr("fib"), // Mark fib func as constant expression.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// Only fib(5) and fib(6) calculated on Compile, fib(dyn) can be called at runtime.
	env["dyn"] = 7

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v\n", output)

	// Output: [5 8 13]
}

func ExampleAllowUndefinedVariables() {
	code := `name == nil ? "Hello, world!" : sprintf("Hello, %v!", name)`

	env := map[string]interface{}{
		"sprintf": fmt.Sprintf,
	}

	options := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(), // Allow to use undefined variables.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v\n", output)

	env["name"] = "you" // Define variables later on.

	output, err = expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v\n", output)

	// Output: Hello, world!
	// Hello, you!
}

func ExampleAllowUndefinedVariables_zero_value() {
	code := `name == "" ? foo + bar : foo + name`

	// If environment has different zero values, then undefined variables
	// will have it as default value.
	env := map[string]string{}

	options := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(), // Allow to use undefined variables.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env = map[string]string{
		"foo": "Hello, ",
		"bar": "world!",
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", output)

	// Output: Hello, world!
}

func ExampleAllowUndefinedVariables_zero_value_functions() {
	code := `words == "" ? Split("foo,bar", ",") : Split(words, ",")`

	// Env is map[string]string type on which methods are defined.
	env := mockMapStringStringEnv{}

	options := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(), // Allow to use undefined variables.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", output)

	// Output: [foo bar]
}

func ExamplePatch() {
	/*
		type patcher struct{}

		func (p *patcher) Visit(node *ast.Node) {
			switch n := (*node).(type) {
			case *ast.MemberNode:
				ast.Patch(node, &ast.CallNode{
					Callee:    &ast.IdentifierNode{Value: "get"},
					Arguments: []ast.Node{n.Node, n.Property},
				})
			}
		}
	*/

	program, err := expr.Compile(
		`greet.you.world + "!"`,
		expr.Patch(&patcher{}),
	)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := map[string]interface{}{
		"greet": "Hello",
		"get": func(a, b string) string {
			return a + ", " + b
		},
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", output)

	// Output : Hello, you, world!
}

func TestOperator_struct(t *testing.T) {
	env := &mockEnv{
		BirthDay: time.Date(2017, time.October, 23, 18, 30, 0, 0, time.UTC),
	}

	code := `BirthDay == "2017-10-23"`

	program, err := expr.Compile(code, expr.Env(&mockEnv{}), expr.Operator("==", "DateEqual"))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestOperator_options_another_order(t *testing.T) {
	code := `BirthDay == "2017-10-23"`
	_, err := expr.Compile(code, expr.Operator("==", "DateEqual"), expr.Env(&mockEnv{}))
	require.NoError(t, err)
}

func TestOperator_no_env(t *testing.T) {
	code := `BirthDay == "2017-10-23"`
	require.Panics(t, func() {
		_, _ = expr.Compile(code, expr.Operator("==", "DateEqual"))
	})
}

func TestOperator_interface(t *testing.T) {
	env := &mockEnv{
		Ticket: &ticket{Price: 100},
	}

	code := `Ticket == "$100" && "$100" == Ticket && Now != Ticket && Now == Now`

	program, err := expr.Compile(
		code,
		expr.Env(&mockEnv{}),
		expr.Operator("==", "StringerStringEqual", "StringStringerEqual", "StringerStringerEqual"),
		expr.Operator("!=", "NotStringerStringEqual", "NotStringStringerEqual", "NotStringerStringerEqual"),
	)
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestExpr_readme_example(t *testing.T) {
	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf,
	}

	code := `sprintf(greet, names[0])`

	program, err := expr.Compile(code, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)

	require.Equal(t, "Hello, world!", output)
}

func TestExpr(t *testing.T) {
	date := time.Date(2017, time.October, 23, 18, 30, 0, 0, time.UTC)
	timeNow := time.Now()
	oneDay, _ := time.ParseDuration("24h")
	timeNowPlusOneDay := timeNow.Add(oneDay)

	env := &mockEnv{
		Any:     "any",
		Int:     0,
		Int32:   0,
		Int64:   0,
		Uint64:  0,
		Float64: 0,
		Bool:    true,
		String:  "string",
		Array:   []int{1, 2, 3, 4, 5},
		Ticket: &ticket{
			Price: 100,
		},
		Passengers: &passengers{
			Adults: 1,
		},
		Segments: []*segment{
			{Origin: "MOW", Destination: "LED"},
			{Origin: "LED", Destination: "MOW"},
		},
		BirthDay:       date,
		Now:            timeNow,
		NowPlusOne:     timeNowPlusOneDay,
		OneDayDuration: oneDay,
		One:            1,
		Two:            2,
		Three:          3,
		MultiDimArray:  [][]int{{1, 2, 3}, {1, 2, 3}},
		Sum: func(list []int) int {
			var ret int
			for _, el := range list {
				ret += el
			}
			return ret
		},
		Inc:       func(a int) int { return a + 1 },
		Nil:       nil,
		Tweets:    []tweet{{"Oh My God!", date}, {"How you doin?", date}, {"Could I be wearing any more clothes?", date}},
		Lowercase: "lowercase",
	}

	tests := []struct {
		code string
		want interface{}
	}{
		{
			`1`,
			1,
		},
		{
			`-.5`,
			-.5,
		},
		{
			`true && false || false`,
			false,
		},
		{
			`Int == 0 && Int32 == 0 && Int64 == 0 && Float64 == 0 && Bool && String == "string"`,
			true,
		},
		{
			`-Int64 == 0`,
			true,
		},
		{
			`"a" != "b"`,
			true,
		},
		{
			`"a" != "b" || 1 == 2`,
			true,
		},
		{
			`Int + 0`,
			0,
		},
		{
			`Uint64 + 0`,
			0,
		},
		{
			`Uint64 + Int64`,
			0,
		},
		{
			`Int32 + Int64`,
			0,
		},
		{
			`Float64 + 0`,
			float64(0),
		},
		{
			`0 + Float64`,
			float64(0),
		},
		{
			`0 <= Float64`,
			true,
		},
		{
			`Float64 < 1`,
			true,
		},
		{
			`Int < 1`,
			true,
		},
		{
			`2 + 2 == 4`,
			true,
		},
		{
			`8 % 3`,
			2,
		},
		{
			`2 ** 8`,
			float64(256),
		},
		{
			`2 ^ 8`,
			float64(256),
		},
		{
			`-(2-5)**3-2/(+4-3)+-2`,
			float64(23),
		},
		{
			`"hello" + " " + "world"`,
			"hello world",
		},
		{
			`0 in -1..1 and 1 in 1..1`,
			true,
		},
		{
			`Int in 0..1`,
			true,
		},
		{
			`Int32 in 0..1`,
			true,
		},
		{
			`Int64 in 0..1`,
			true,
		},
		{
			`1 in [1, 2, 3] && "foo" in {foo: 0, bar: 1} && "Price" in Ticket`,
			true,
		},
		{
			`1 in [1.5] || 1 not in [1]`,
			false,
		},
		{
			`One in 0..1 && Two not in 0..1`,
			true,
		},
		{
			`Two not in 0..1`,
			true,
		},
		{
			`Two not    in 0..1`,
			true,
		},
		{
			`Int32 in [10, 20]`,
			false,
		},
		{
			`String matches "s.+"`,
			true,
		},
		{
			`String matches ("^" + String + "$")`,
			true,
		},
		{
			`"foobar" contains "bar"`,
			true,
		},
		{
			`"foobar" startsWith "foo"`,
			true,
		},
		{
			`"foobar" endsWith "bar"`,
			true,
		},
		{
			`(0..10)[5]`,
			5,
		},
		{
			`Ticket.Price`,
			100,
		},
		{
			`Add(10, 5) + GetInt()`,
			15,
		},
		{
			`Ticket.String()`,
			`$100`,
		},
		{
			`Ticket.PriceDiv(25)`,
			4,
		},
		{
			`len([1, 2, 3])`,
			3,
		},
		{
			`len([1, Two, 3])`,
			3,
		},
		{
			`len(["hello", "world"])`,
			2,
		},
		{
			`len("hello, world")`,
			12,
		},
		{
			`len(Array)`,
			5,
		},
		{
			`len({a: 1, b: 2, c: 2})`,
			3,
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]interface{}{"foo": 0, "bar": 1},
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]interface{}{"foo": 0, "bar": 1},
		},
		{
			`(true ? 0+1 : 2+3) + (false ? -1 : -2)`,
			-1,
		},
		{
			`filter(1..9, {# > 7})`,
			[]interface{}{8, 9},
		},
		{
			`map(1..3, {# * #})`,
			[]interface{}{1, 4, 9},
		},
		{
			`all(1..3, {# > 0})`,
			true,
		},
		{
			`none(1..3, {# == 0})`,
			true,
		},
		{
			`any([1,1,0,1], {# == 0})`,
			true,
		},
		{
			`one([1,1,0,1], {# == 0}) and not one([1,0,0,1], {# == 0})`,
			true,
		},
		{
			`count(1..30, {# % 3 == 0})`,
			10,
		},
		{
			`Now.After(BirthDay)`,
			true,
		},
		{
			`"a" < "b"`,
			true,
		},
		{
			`Now.Sub(Now).String() == Duration("0s").String()`,
			true,
		},
		{
			`8.5 * Passengers.Adults * len(Segments)`,
			float64(17),
		},
		{
			`1 + 1`,
			2,
		},
		{
			`(One * Two) * Three == One * (Two * Three)`,
			true,
		},
		{
			`Array[0]`,
			1,
		},
		{
			`Sum(Array)`,
			15,
		},
		{
			`Array[0] < Array[1]`,
			true,
		},
		{
			`Sum(MultiDimArray[0])`,
			6,
		},
		{
			`Sum(MultiDimArray[0]) + Sum(MultiDimArray[1])`,
			12,
		},
		{
			`Inc(Array[0] + Array[1])`,
			4,
		},
		{
			`Array[0] + Array[1]`,
			3,
		},
		{
			`Array[1:2]`,
			[]int{2},
		},
		{
			`Array[1:4]`,
			[]int{2, 3, 4},
		},
		{
			`Array[:3]`,
			[]int{1, 2, 3},
		},
		{
			`Array[3:]`,
			[]int{4, 5},
		},
		{
			`Array[0:5] == Array`,
			true,
		},
		{
			`Array[0:] == Array`,
			true,
		},
		{
			`Array[:5] == Array`,
			true,
		},
		{
			`Array[:] == Array`,
			true,
		},
		{
			`1 + 2 + Three`,
			6,
		},
		{
			`4 in 5..1`,
			false,
		},
		{
			`4..0`,
			[]int{},
		},
		{
			`MapArg({foo: "bar"})`,
			"bar",
		},
		{
			`NilStruct`,
			(*time.Time)(nil),
		},
		{
			`Nil == nil && nil == Nil && nil == nil && Nil == Nil && NilInt == nil && NilSlice == nil && NilStruct == nil`,
			true,
		},
		{
			`0 == nil || "str" == nil || true == nil`,
			false,
		},
		{
			`Variadic("head", 1, 2, 3)`,
			[]int{1, 2, 3},
		},
		{
			`Variadic("empty")`,
			[]int{},
		},
		{
			`String[:]`,
			"string",
		},
		{
			`String[:3]`,
			"str",
		},
		{
			`String[:9]`,
			"string",
		},
		{
			`String[3:9]`,
			"ing",
		},
		{
			`String[7:9]`,
			"",
		},
		{
			`Float(0)`,
			float64(0),
		},
		{
			`map(filter(Tweets, {len(.Text) > 10}), {Format(.Date)})`,
			[]interface{}{"23 Oct 17 18:30 UTC", "23 Oct 17 18:30 UTC"},
		},
		{
			`Concat("a", 1, [])`,
			`a1[]`,
		},
		{
			`Tweets[0].Date < Now`,
			true,
		},
		{
			`Now > Tweets[0].Date`,
			true,
		},
		{
			`Now == Now`,
			true,
		},
		{
			`Now >= Now`,
			true,
		},
		{
			`Now <= Now`,
			true,
		},
		{
			`Now == NowPlusOne`,
			false,
		},
		{
			`Now != Now`,
			false,
		},
		{
			`Now != NowPlusOne`,
			true,
		},
		{
			`NowPlusOne - Now`,
			oneDay,
		},
		{
			`Now + OneDayDuration`,
			timeNowPlusOneDay,
		},
		{
			`OneDayDuration + Now`,
			timeNowPlusOneDay,
		},
		{
			`lowercase`,
			"lowercase",
		},
	}

	for _, tt := range tests {
		program, err := expr.Compile(tt.code, expr.Env(&mockEnv{}))
		require.NoError(t, err, "compile error")

		got, err := expr.Run(program, env)
		require.NoError(t, err, "execution error")

		assert.Equal(t, tt.want, got, tt.code)
	}

	for _, tt := range tests {
		if tt.code == `-Int64 == 0` {
			program, err := expr.Compile(tt.code, expr.Optimize(false))
			require.NoError(t, err, "compile error")

			got, err := expr.Run(program, env)
			require.NoError(t, err, "run error")
			assert.Equal(t, tt.want, got, "unoptimized: "+tt.code)
		}
	}

	for _, tt := range tests {
		got, err := expr.Eval(tt.code, env)
		require.NoError(t, err, "eval error: "+tt.code)

		assert.Equal(t, tt.want, got, "eval: "+tt.code)
	}
}

func TestExpr_optional_chaining(t *testing.T) {
	env := map[string]interface{}{}
	program, err := expr.Compile("foo?.bar.baz", expr.Env(env), expr.AllowUndefinedVariables())
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_optional_chaining_property(t *testing.T) {
	env := map[string]interface{}{
		"foo": map[string]interface{}{},
	}
	program, err := expr.Compile("foo.bar?.baz", expr.Env(env))
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_optional_chaining_nested_chains(t *testing.T) {
	env := map[string]interface{}{
		"foo": map[string]interface{}{
			"id": 1,
			"bar": []map[string]interface{}{
				1: {
					"baz": "baz",
				},
			},
		},
	}
	program, err := expr.Compile("foo?.bar[foo?.id]?.baz", expr.Env(env))
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "baz", got)
}

func TestExpr_eval_with_env(t *testing.T) {
	_, err := expr.Eval("true", expr.Env(map[string]interface{}{}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "misused")
}

func TestExpr_fetch_from_func(t *testing.T) {
	_, err := expr.Eval("foo.Value", map[string]interface{}{
		"foo": func() {},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot fetch Value from func()")
}

func TestExpr_map_default_values(t *testing.T) {
	env := map[string]interface{}{
		"foo": map[string]string{},
		"bar": map[string]*string{},
	}

	input := `foo['missing'] == '' && bar['missing'] == nil`

	program, err := expr.Compile(input, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestExpr_map_default_values_compile_check(t *testing.T) {
	tests := []struct {
		env   interface{}
		input string
	}{
		{
			mockMapStringStringEnv{"foo": "bar"},
			`Split(foo, sep)`,
		},
		{
			mockMapStringIntEnv{"foo": 1},
			`foo / bar`,
		},
	}
	for _, tt := range tests {
		_, err := expr.Compile(tt.input, expr.Env(tt.env), expr.AllowUndefinedVariables())
		require.NoError(t, err)
	}
}

func TestExpr_calls_with_nil(t *testing.T) {
	env := map[string]interface{}{
		"equals": func(a, b interface{}) interface{} {
			assert.Nil(t, a, "a is not nil")
			assert.Nil(t, b, "b is not nil")
			return a == b
		},
		"is": is{},
	}

	p, err := expr.Compile(
		"a == nil && equals(b, nil) && is.Nil(c)",
		expr.Env(env),
		expr.Operator("==", "equals"),
		expr.AllowUndefinedVariables(),
	)
	require.NoError(t, err)

	out, err := expr.Run(p, env)
	require.NoError(t, err)
	require.Equal(t, true, out)
}

func TestExpr_call_floatarg_func_with_int(t *testing.T) {
	env := map[string]interface{}{
		"cnv": func(f float64) interface{} {
			return f
		},
	}
	for _, each := range []struct {
		input    string
		expected float64
	}{
		{"-1", -1.0},
		{"1+1", 2.0},
		{"+1", 1.0},
		{"1-1", 0.0},
		{"1/1", 1.0},
		{"1*1", 1.0},
	} {
		p, err := expr.Compile(
			fmt.Sprintf("cnv(%s)", each.input),
			expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(p, env)
		require.NoError(t, err)
		require.Equal(t, each.expected, out)
	}
}

func TestConstExpr_error_panic(t *testing.T) {
	env := map[string]interface{}{
		"divide": func(a, b int) int { return a / b },
	}

	_, err := expr.Compile(
		`1 + divide(1, 0)`,
		expr.Env(env),
		expr.ConstExpr("divide"),
	)
	require.Error(t, err)
	require.Equal(t, "compile error: integer divide by zero (1:5)\n | 1 + divide(1, 0)\n | ....^", err.Error())
}

type divideError struct{ Message string }

func (e divideError) Error() string {
	return e.Message
}

func TestConstExpr_error_as_error(t *testing.T) {
	env := map[string]interface{}{
		"divide": func(a, b int) (int, error) {
			if b == 0 {
				return 0, divideError{"integer divide by zero"}
			}
			return a / b, nil
		},
	}

	_, err := expr.Compile(
		`1 + divide(1, 0)`,
		expr.Env(env),
		expr.ConstExpr("divide"),
	)
	require.Error(t, err)
	require.Equal(t, "integer divide by zero", err.Error())
	require.IsType(t, divideError{}, err)
}

func TestConstExpr_error_wrong_type(t *testing.T) {
	env := map[string]interface{}{
		"divide": 0,
	}
	assert.Panics(t, func() {
		_, _ = expr.Compile(
			`1 + divide(1, 0)`,
			expr.Env(env),
			expr.ConstExpr("divide"),
		)
	})
}

func TestConstExpr_error_no_env(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = expr.Compile(
			`1 + divide(1, 0)`,
			expr.ConstExpr("divide"),
		)
	})
}

var stringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

type stringerPatcher struct{}

func (p *stringerPatcher) Visit(node *ast.Node) {
	t := (*node).Type()
	if t == nil {
		return
	}
	if t.Implements(stringer) {
		ast.Patch(node, &ast.CallNode{
			Callee: &ast.MemberNode{
				Node:     *node,
				Property: &ast.StringNode{Value: "String"},
			},
		})
	}
}

func TestPatch(t *testing.T) {
	program, err := expr.Compile(
		`Ticket == "$100" and "$90" != Ticket + "0"`,
		expr.Env(mockEnv{}),
		expr.Patch(&stringerPatcher{}),
	)
	require.NoError(t, err)

	env := mockEnv{
		Ticket: &ticket{Price: 100},
	}
	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

type lengthPatcher struct{}

func (p *lengthPatcher) Visit(node *ast.Node) {
	switch n := (*node).(type) {
	case *ast.MemberNode:
		if prop, ok := n.Property.(*ast.StringNode); ok && prop.Value == "length" {
			ast.Patch(node, &ast.BuiltinNode{
				Name:      "len",
				Arguments: []ast.Node{n.Node},
			})
		}
	}
}

func TestPatch_length(t *testing.T) {
	program, err := expr.Compile(
		`String.length == 5`,
		expr.Env(mockEnv{}),
		expr.Patch(&lengthPatcher{}),
	)
	require.NoError(t, err)

	env := mockEnv{String: "hello"}
	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestCompile_exposed_error(t *testing.T) {
	_, err := expr.Compile(`1 == true`)
	require.Error(t, err)

	fileError, ok := err.(*file.Error)
	require.True(t, ok, "error should be of type *file.Error")
	require.Equal(t, "invalid operation: == (mismatched types int and bool) (1:3)\n | 1 == true\n | ..^", fileError.Error())
	require.Equal(t, 2, fileError.Column)
	require.Equal(t, 1, fileError.Line)

	b, err := json.Marshal(err)
	require.NoError(t, err)
	require.Equal(t, `{"Line":1,"Column":2,"Message":"invalid operation: == (mismatched types int and bool)","Snippet":"\n | 1 == true\n | ..^"}`, string(b))
}

func TestCompile_deref(t *testing.T) {
	i := 1
	env := map[string]interface{}{
		"i": &i,
		"map": map[string]interface{}{
			"i": &i,
		},
	}
	{
		// With specified env, OpDeref added and == works as expected.
		program, err := expr.Compile(`i == 1 && map.i == 1`, expr.Env(env))
		require.NoError(t, err)

		env["any"] = &i
		out, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, true, out)
	}
	{
		// Compile without expr.Env() also works as expected,
		// and should add OpDeref automatically.
		program, err := expr.Compile(`i == 1 && map.i == 1`)
		require.NoError(t, err)

		out, err := expr.Run(program, env)
		require.NoError(t, err)
		require.Equal(t, true, out)
	}
}

func TestEval_deref(t *testing.T) {
	i := 1
	env := map[string]interface{}{
		"i": &i,
		"map": map[string]interface{}{
			"i": &i,
		},
	}

	out, err := expr.Eval(`i == 1 && map.i == 1`, env)
	require.NoError(t, err)
	require.Equal(t, true, out)
}

func TestAsBool_exposed_error(t *testing.T) {
	_, err := expr.Compile(`42`, expr.AsBool())
	require.Error(t, err)

	_, ok := err.(*file.Error)
	require.False(t, ok, "error must not be of type *file.Error")
	require.Equal(t, "expected bool, but got int", err.Error())
}

func TestEval_exposed_error(t *testing.T) {
	_, err := expr.Eval(`1 % 0`, nil)
	require.Error(t, err)

	fileError, ok := err.(*file.Error)
	require.True(t, ok, "error should be of type *file.Error")
	require.Equal(t, "runtime error: integer divide by zero (1:3)\n | 1 % 0\n | ..^", fileError.Error())
	require.Equal(t, 2, fileError.Column)
	require.Equal(t, 1, fileError.Line)
}

func TestIssue105(t *testing.T) {
	type A struct {
		Field string
	}
	type B struct {
		Field int
	}
	type C struct {
		A
		B
	}
	type Env struct {
		C
	}

	code := `
		A.Field == '' &&
		C.A.Field == '' &&
		B.Field == 0 &&
		C.B.Field == 0
	`

	_, err := expr.Compile(code, expr.Env(Env{}))
	require.NoError(t, err)

	_, err = expr.Compile(`Field == ''`, expr.Env(Env{}))
	require.Error(t, err)
	require.Contains(t, err.Error(), "ambiguous identifier Field")
}

func TestIssue_nested_closures(t *testing.T) {
	code := `all(1..3, { all(1..3, { # > 0 }) and # > 0 })`

	program, err := expr.Compile(code)
	require.NoError(t, err)

	output, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.True(t, output.(bool))
}

func TestIssue138(t *testing.T) {
	env := map[string]interface{}{}

	_, err := expr.Compile(`1 / (1 - 1)`, expr.Env(env))
	require.NoError(t, err)

	_, err = expr.Compile(`1 % 0`, expr.Env(env))
	require.Error(t, err)
	require.Equal(t, "integer divide by zero (1:3)\n | 1 % 0\n | ..^", err.Error())
}

func TestIssue154(t *testing.T) {
	type Data struct {
		Array  *[2]interface{}
		Slice  *[]interface{}
		Map    *map[string]interface{}
		String *string
	}

	type Env struct {
		Data *Data
	}

	b := true
	i := 10
	s := "value"

	Array := [2]interface{}{
		&b,
		&i,
	}

	Slice := []interface{}{
		&b,
		&i,
	}

	Map := map[string]interface{}{
		"Bool": &b,
		"Int":  &i,
	}

	env := Env{
		Data: &Data{
			Array:  &Array,
			Slice:  &Slice,
			Map:    &Map,
			String: &s,
		},
	}

	tests := []string{
		`Data.Array[0] == true`,
		`Data.Array[1] == 10`,
		`Data.Slice[0] == true`,
		`Data.Slice[1] == 10`,
		`Data.Map["Bool"] == true`,
		`Data.Map["Int"] == 10`,
		`Data.String == "value"`,
	}

	for _, input := range tests {
		program, err := expr.Compile(input, expr.Env(env))
		assert.NoError(t, err, input)

		output, err := expr.Run(program, env)
		assert.NoError(t, err)
		assert.True(t, output.(bool), input)
	}
}

func TestIssue270(t *testing.T) {
	env := map[string]interface{}{
		"int8":     int8(1),
		"int16":    int16(3),
		"int32":    int32(5),
		"int64":    int64(7),
		"uint8":    uint8(11),
		"uint16":   uint16(13),
		"uint32":   uint32(17),
		"uint64":   uint64(19),
		"int8a":    uint(23),
		"int8b":    uint(29),
		"int16a":   uint(31),
		"int16b":   uint(37),
		"int32a":   uint(41),
		"int32b":   uint(43),
		"int64a":   uint(47),
		"int64b":   uint(53),
		"uint8a":   uint(59),
		"uint8b":   uint(61),
		"uint16a":  uint(67),
		"uint16b":  uint(71),
		"uint32a":  uint(73),
		"uint32b":  uint(79),
		"uint64a":  uint(83),
		"uint64b":  uint(89),
		"float32a": float32(97),
		"float32b": float32(101),
		"float64a": float64(103),
		"float64b": float64(107),
	}
	for _, each := range []struct {
		input string
	}{
		{"int8 / int16"},
		{"int32 / int64"},
		{"uint8 / uint16"},
		{"uint32 / uint64"},
		{"int8 / uint64"},
		{"int64 / uint8"},
		{"int8a / int8b"},
		{"int16a / int16b"},
		{"int32a / int32b"},
		{"int64a / int64b"},
		{"uint8a / uint8b"},
		{"uint16a / uint16b"},
		{"uint32a / uint32b"},
		{"uint64a / uint64b"},
		{"float32a / float32b"},
		{"float64a / float64b"},
	} {
		p, err := expr.Compile(each.input, expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(p, env)
		require.NoError(t, err)
		require.IsType(t, float64(0), out)
	}
}

func TestIssue271(t *testing.T) {
	type BarArray []float64

	type Foo struct {
		Bar BarArray
		Baz int
	}

	type Env struct {
		Foo Foo
	}

	code := `Foo.Bar[0]`

	program, err := expr.Compile(code, expr.Env(Env{}))
	require.NoError(t, err)

	output, err := expr.Run(program, Env{
		Foo: Foo{
			Bar: BarArray{1.0, 2.0, 3.0},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 1.0, output)
}

// Mock types
type mockEnv struct {
	Any                  interface{}
	Int, One, Two, Three int
	Int32                int32
	Int64                int64
	Uint64               uint64
	Float64              float64
	Bool                 bool
	String               string
	Array                []int
	MultiDimArray        [][]int
	Sum                  func(list []int) int
	Inc                  func(int) int
	Ticket               *ticket
	Passengers           *passengers
	Segments             []*segment
	BirthDay             time.Time
	Now                  time.Time
	NowPlusOne           time.Time
	OneDayDuration       time.Duration
	Nil                  interface{}
	NilStruct            *time.Time
	NilInt               *int
	NilSlice             []ticket
	Tweets               []tweet
	Lowercase            string `expr:"lowercase"`
}

func (e *mockEnv) GetInt() int {
	return e.Int
}

func (*mockEnv) Add(a, b int) int {
	return a + b
}

func (*mockEnv) Duration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func (*mockEnv) MapArg(m map[string]interface{}) string {
	return m["foo"].(string)
}

func (*mockEnv) DateEqual(date time.Time, s string) bool {
	return date.Format("2006-01-02") == s
}

func (*mockEnv) StringerStringEqual(f fmt.Stringer, s string) bool {
	return f.String() == s
}

func (*mockEnv) StringStringerEqual(s string, f fmt.Stringer) bool {
	return s == f.String()
}

func (*mockEnv) StringerStringerEqual(f fmt.Stringer, g fmt.Stringer) bool {
	return f.String() == g.String()
}

func (*mockEnv) NotStringerStringEqual(f fmt.Stringer, s string) bool {
	return f.String() != s
}

func (*mockEnv) NotStringStringerEqual(s string, f fmt.Stringer) bool {
	return s != f.String()
}

func (*mockEnv) NotStringerStringerEqual(f fmt.Stringer, g fmt.Stringer) bool {
	return f.String() != g.String()
}

func (*mockEnv) Variadic(x string, xs ...int) []int {
	return xs
}

func (*mockEnv) Concat(list ...interface{}) string {
	out := ""
	for _, e := range list {
		out += fmt.Sprintf("%v", e)
	}
	return out
}

func (*mockEnv) Float(i interface{}) float64 {
	switch t := i.(type) {
	case int:
		return float64(t)
	case float64:
		return t
	default:
		panic("unexpected type")
	}
}

func (*mockEnv) Format(t time.Time) string {
	return t.Format(time.RFC822)
}

type ticket struct {
	Price int
}

func (t *ticket) PriceDiv(p int) int {
	return t.Price / p
}

func (t *ticket) String() string {
	return fmt.Sprintf("$%v", t.Price)
}

type passengers struct {
	Adults   uint32
	Children uint32
	Infants  uint32
}

type segment struct {
	Origin      string
	Destination string
	Date        time.Time
}

type tweet struct {
	Text string
	Date time.Time
}

type mockMapStringStringEnv map[string]string

func (m mockMapStringStringEnv) Split(s, sep string) []string {
	return strings.Split(s, sep)
}

type mockMapStringIntEnv map[string]int

type is struct{}

func (is) Nil(a interface{}) bool {
	return a == nil
}

type patcher struct{}

func (p *patcher) Visit(node *ast.Node) {
	switch n := (*node).(type) {
	case *ast.MemberNode:
		ast.Patch(node, &ast.CallNode{
			Callee:    &ast.IdentifierNode{Value: "get"},
			Arguments: []ast.Node{n.Node, n.Property},
		})
	}
}
