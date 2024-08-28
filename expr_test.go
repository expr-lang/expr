package expr_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/types"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/test/mock"
)

func ExampleEval() {
	output, err := expr.Eval("greet + name", map[string]any{
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
	env := map[string]any{
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
	env := map[string]any{
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
	env := map[string]any{
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

func ExampleWarnOnAny() {
	// Arrays always have []any type. The expression return type is any.
	// AsInt() instructs compiler to expect int or any, and cast to int,
	// if possible. WarnOnAny() instructs to return an error on any type.
	_, err := expr.Compile(`[42, true, "yes"][0]`, expr.AsInt(), expr.WarnOnAny())

	fmt.Printf("%v", err)

	// Output: expected int, but got interface {}
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

func ExampleOperator_Decimal() {
	type Decimal struct{ N float64 }
	code := `A + B - C`

	type Env struct {
		A, B, C Decimal
		Sub     func(a, b Decimal) Decimal
		Add     func(a, b Decimal) Decimal
	}

	options := []expr.Option{
		expr.Env(Env{}),
		expr.Operator("+", "Add"),
		expr.Operator("-", "Sub"),
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("Compile error: %v", err)
		return
	}

	env := Env{
		A:   Decimal{3},
		B:   Decimal{2},
		C:   Decimal{1},
		Sub: func(a, b Decimal) Decimal { return Decimal{a.N - b.N} },
		Add: func(a, b Decimal) Decimal { return Decimal{a.N + b.N} },
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: {4}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func ExampleConstExpr() {
	code := `[fib(5), fib(3+3), fib(dyn)]`

	env := map[string]any{
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

	env := map[string]any{
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
	env := mock.MapStringStringEnv{}

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

	env := map[string]any{
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

func ExampleWithContext() {
	env := map[string]any{
		"fn": func(ctx context.Context, _, _ int) int {
			// An infinite loop that can be canceled by context.
			for {
				select {
				case <-ctx.Done():
					return 42
				}
			}
		},
		"ctx": context.TODO(), // Context should be passed as a variable.
	}

	program, err := expr.Compile(`fn(1, 2)`,
		expr.Env(env),
		expr.WithContext("ctx"), // Pass context variable name.
	)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// Cancel context after 100 milliseconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	// After program is compiled, context can be passed to Run.
	env["ctx"] = ctx

	// Run will return 42 after 100 milliseconds.
	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)
	// Output: 42
}

func ExampleWithTimezone() {
	program, err := expr.Compile(`now().Location().String()`, expr.Timezone("Asia/Kamchatka"))
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
	// Output: Asia/Kamchatka
}

func TestExpr_readme_example(t *testing.T) {
	env := map[string]any{
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
	oneDay, _ := time.ParseDuration("24h")
	timeNowPlusOneDay := date.Add(oneDay)

	env := mock.Env{
		Embed:     mock.Embed{},
		Ambiguous: "",
		Any:       nil,
		Bool:      true,
		Float:     0,
		Int64:     0,
		Int32:     0,
		Int:       0,
		One:       1,
		Two:       2,
		Uint32:    0,
		String:    "string",
		BoolPtr:   nil,
		FloatPtr:  nil,
		IntPtr:    nil,
		IntPtrPtr: nil,
		StringPtr: nil,
		Foo: mock.Foo{
			Value: "foo",
			Bar: mock.Bar{
				Baz: "baz",
			},
		},
		Abstract:           nil,
		ArrayOfAny:         nil,
		ArrayOfInt:         []int{1, 2, 3, 4, 5},
		ArrayOfFoo:         []*mock.Foo{{Value: "foo"}, {Value: "bar"}, {Value: "baz"}},
		MapOfFoo:           nil,
		MapOfAny:           nil,
		FuncParam:          nil,
		FuncParamAny:       nil,
		FuncTooManyReturns: nil,
		FuncNamed:          nil,
		NilAny:             nil,
		NilFn:              nil,
		NilStruct:          nil,
		Variadic: func(head int, xs ...int) bool {
			sum := 0
			for _, x := range xs {
				sum += x
			}
			return head == sum
		},
		Fast:        nil,
		Time:        date,
		TimePlusDay: timeNowPlusOneDay,
		Duration:    oneDay,
	}

	tests := []struct {
		code string
		want any
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
			`1 in [1, 2, 3] && "foo" in {foo: 0, bar: 1} && "Bar" in Foo`,
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
			`-1 not in [1]`,
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
			`'foo' + 'bar' not matches 'foobar'`,
			false,
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
			`Foo.Bar.Baz`,
			"baz",
		},
		{
			`Add(10, 5) + GetInt()`,
			15,
		},
		{
			`Foo.Method().Baz`,
			`baz (from Foo.Method)`,
		},
		{
			`Foo.MethodWithArgs("prefix ")`,
			"prefix foo",
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
			`len(ArrayOfInt)`,
			5,
		},
		{
			`len({a: 1, b: 2, c: 2})`,
			3,
		},
		{
			`max([1, 2, 3])`,
			3,
		},
		{
			`max(1, 2, 3)`,
			3,
		},
		{
			`min([1, 2, 3])`,
			1,
		},
		{
			`min(1, 2, 3)`,
			1,
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]any{"foo": 0, "bar": 1},
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]any{"foo": 0, "bar": 1},
		},
		{
			`(true ? 0+1 : 2+3) + (false ? -1 : -2)`,
			-1,
		},
		{
			`filter(1..9, {# > 7})`,
			[]any{8, 9},
		},
		{
			`map(1..3, {# * #})`,
			[]any{1, 4, 9},
		},
		{
			`all(1..3, {# > 0})`,
			true,
		},
		{
			`count(1..30, {# % 3 == 0})`,
			10,
		},
		{
			`count([true, true, false])`,
			2,
		},
		{
			`"a" < "b"`,
			true,
		},
		{
			`Time.Sub(Time).String() == "0s"`,
			true,
		},
		{
			`1 + 1`,
			2,
		},
		{
			`(One * Two) * 3 == One * (Two * 3)`,
			true,
		},
		{
			`ArrayOfInt[1]`,
			2,
		},
		{
			`ArrayOfInt[0] < ArrayOfInt[1]`,
			true,
		},
		{
			`ArrayOfInt[-1]`,
			5,
		},
		{
			`ArrayOfInt[1:2]`,
			[]int{2},
		},
		{
			`ArrayOfInt[1:4]`,
			[]int{2, 3, 4},
		},
		{
			`ArrayOfInt[-4:-1]`,
			[]int{2, 3, 4},
		},
		{
			`ArrayOfInt[:3]`,
			[]int{1, 2, 3},
		},
		{
			`ArrayOfInt[3:]`,
			[]int{4, 5},
		},
		{
			`ArrayOfInt[0:5] == ArrayOfInt`,
			true,
		},
		{
			`ArrayOfInt[0:] == ArrayOfInt`,
			true,
		},
		{
			`ArrayOfInt[:5] == ArrayOfInt`,
			true,
		},
		{
			`ArrayOfInt[:] == ArrayOfInt`,
			true,
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
			`NilStruct`,
			(*mock.Foo)(nil),
		},
		{
			`NilAny == nil && nil == NilAny && nil == nil && NilAny == NilAny && NilInt == nil && NilSlice == nil && NilStruct == nil`,
			true,
		},
		{
			`0 == nil || "str" == nil || true == nil`,
			false,
		},
		{
			`Variadic(6, 1, 2, 3)`,
			true,
		},
		{
			`Variadic(0)`,
			true,
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
			`map(filter(ArrayOfInt, # >= 3), # + 1)`,
			[]any{4, 5, 6},
		},
		{
			`Time < Time + Duration`,
			true,
		},
		{
			`Time + Duration > Time`,
			true,
		},
		{
			`Time == Time`,
			true,
		},
		{
			`Time >= Time`,
			true,
		},
		{
			`Time <= Time`,
			true,
		},
		{
			`Time == Time + Duration`,
			false,
		},
		{
			`Time != Time`,
			false,
		},
		{
			`TimePlusDay - Duration`,
			date,
		},
		{
			`duration("1h") == duration("1h")`,
			true,
		},
		{
			`TimePlusDay - Time >= duration("24h")`,
			true,
		},
		{
			`duration("1h") > duration("1m")`,
			true,
		},
		{
			`duration("1h") < duration("1m")`,
			false,
		},
		{
			`duration("1h") >= duration("1m")`,
			true,
		},
		{
			`duration("1h") <= duration("1m")`,
			false,
		},
		{
			`duration("1h") > duration("1m")`,
			true,
		},
		{
			`duration("1h") + duration("1m")`,
			time.Hour + time.Minute,
		},
		{
			`7 * duration("1h")`,
			7 * time.Hour,
		},
		{
			`duration("1h") * 7`,
			7 * time.Hour,
		},
		{
			`duration("1s") * .5`,
			5e8,
		},
		{
			`1 /* one */ + 2 // two`,
			3,
		},
		{
			`let x = 1; x + 2`,
			3,
		},
		{
			`map(1..3, let x = #; let y = x * x; y * y)`,
			[]any{1, 16, 81},
		},
		{
			`map(1..2, let x = #; map(2..3, let y = #; x + y))`,
			[]any{[]any{3, 4}, []any{4, 5}},
		},
		{
			`len(filter(1..99, # % 7 == 0))`,
			14,
		},
		{
			`find(ArrayOfFoo, .Value == "baz")`,
			env.ArrayOfFoo[2],
		},
		{
			`findIndex(ArrayOfFoo, .Value == "baz")`,
			2,
		},
		{
			`filter(ArrayOfFoo, .Value == "baz")[0]`,
			env.ArrayOfFoo[2],
		},
		{
			`first(filter(ArrayOfFoo, .Value == "baz"))`,
			env.ArrayOfFoo[2],
		},
		{
			`first(filter(ArrayOfFoo, false))`,
			nil,
		},
		{
			`findLast(1..9, # % 2 == 0)`,
			8,
		},
		{
			`findLastIndex(1..9, # % 2 == 0)`,
			7,
		},
		{
			`filter(1..9, # % 2 == 0)[-1]`,
			8,
		},
		{
			`last(filter(1..9, # % 2 == 0))`,
			8,
		},
		{
			`map(filter(1..9, # % 2 == 0), # * 2)`,
			[]any{4, 8, 12, 16},
		},
		{
			`map(map(filter(1..9, # % 2 == 0), # * 2), # * 2)`,
			[]any{8, 16, 24, 32},
		},
		{
			`first(map(filter(1..9, # % 2 == 0), # * 2))`,
			4,
		},
		{
			`map(filter(1..9, # % 2 == 0), # * 2)[-1]`,
			16,
		},
		{
			`len(map(filter(1..9, # % 2 == 0), # * 2))`,
			4,
		},
		{
			`len(filter(map(1..9, # * 2), # % 2 == 0))`,
			9,
		},
		{
			`first(filter(map(1..9, # * 2), # % 2 == 0))`,
			2,
		},
		{
			`first(map(filter(1..9, # % 2 == 0), # * 2))`,
			4,
		},
		{
			`2^3 == 8`,
			true,
		},
		{
			`4/2 == 2`,
			true,
		},
		{
			`.5 in 0..1`,
			false,
		},
		{
			`.5 in ArrayOfInt`,
			false,
		},
		{
			`bitnot(10)`,
			-11,
		},
		{
			`bitxor(15, 32)`,
			47,
		},
		{
			`bitand(90, 34)`,
			2,
		},
		{
			`bitnand(35, 9)`,
			34,
		},
		{
			`bitor(10, 5)`,
			15,
		},
		{
			`bitshr(7, 2)`,
			1,
		},
		{
			`bitshl(7, 2)`,
			28,
		},
		{
			`bitushr(-100, 5)`,
			576460752303423484,
		},
		{
			`"hello"[1:3]`,
			"el",
		},
		{
			`[1, 2, 3]?.[0]`,
			1,
		},
		{
			`[[1, 2], 3, 4]?.[0]?.[1]`,
			2,
		},
		{
			`[nil, 3, 4]?.[0]?.[1]`,
			nil,
		},
		{
			`1 > 2 < 3`,
			false,
		},
		{
			`1 < 2 < 3`,
			true,
		},
		{
			`1 < 2 < 3 > 4`,
			false,
		},
		{
			`1 < 2 < 3 > 2`,
			true,
		},
		{
			`1 < 2 < 3 == true`,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			{
				program, err := expr.Compile(tt.code, expr.Env(mock.Env{}))
				require.NoError(t, err, "compile error")

				got, err := expr.Run(program, env)
				require.NoError(t, err, "run error")
				assert.Equal(t, tt.want, got)
			}
			{
				program, err := expr.Compile(tt.code, expr.Optimize(false))
				require.NoError(t, err, "unoptimized")

				got, err := expr.Run(program, env)
				require.NoError(t, err, "unoptimized")
				assert.Equal(t, tt.want, got, "unoptimized")
			}
			{
				got, err := expr.Eval(tt.code, env)
				require.NoError(t, err, "eval")
				assert.Equal(t, tt.want, got, "eval")
			}
			{
				program, err := expr.Compile(tt.code, expr.Env(mock.Env{}), expr.Optimize(false))
				require.NoError(t, err)

				code := program.Node().String()
				got, err := expr.Eval(code, env)
				require.NoError(t, err, code)
				assert.Equal(t, tt.want, got, code)
			}
		})
	}
}

func TestExpr_error(t *testing.T) {
	env := mock.Env{
		ArrayOfAny: []any{1, "2", 3, true},
	}

	tests := []struct {
		code string
		want string
	}{
		{
			`filter(1..9, # > 9)[0]`,
			`reflect: slice index out of range (1:20)
 | filter(1..9, # > 9)[0]
 | ...................^`,
		},
		{
			`ArrayOfAny[-7]`,
			`index out of range: -3 (array length is 4) (1:11)
 | ArrayOfAny[-7]
 | ..........^`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code, expr.Env(mock.Env{}))
			require.NoError(t, err)

			_, err = expr.Run(program, env)
			require.Error(t, err)
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestExpr_optional_chaining(t *testing.T) {
	env := map[string]any{}
	program, err := expr.Compile("foo?.bar.baz", expr.Env(env), expr.AllowUndefinedVariables())
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_optional_chaining_property(t *testing.T) {
	env := map[string]any{
		"foo": map[string]any{},
	}
	program, err := expr.Compile("foo.bar?.baz", expr.Env(env))
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_optional_chaining_nested_chains(t *testing.T) {
	env := map[string]any{
		"foo": map[string]any{
			"id": 1,
			"bar": []map[string]any{
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

func TestExpr_optional_chaining_array(t *testing.T) {
	env := map[string]any{}
	program, err := expr.Compile("foo?.[1]?.[2]?.[3]", expr.Env(env), expr.AllowUndefinedVariables())
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_eval_with_env(t *testing.T) {
	_, err := expr.Eval("true", expr.Env(map[string]any{}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "misused")
}

func TestExpr_fetch_from_func(t *testing.T) {
	_, err := expr.Eval("foo.Value", map[string]any{
		"foo": func() {},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot fetch Value from func()")
}

func TestExpr_map_default_values(t *testing.T) {
	env := map[string]any{
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
		env   any
		input string
	}{
		{
			mock.MapStringStringEnv{"foo": "bar"},
			`Split(foo, sep)`,
		},
		{
			mock.MapStringIntEnv{"foo": 1},
			`foo / bar`,
		},
	}
	for _, tt := range tests {
		_, err := expr.Compile(tt.input, expr.Env(tt.env), expr.AllowUndefinedVariables())
		require.NoError(t, err)
	}
}

func TestExpr_calls_with_nil(t *testing.T) {
	env := map[string]any{
		"equals": func(a, b any) any {
			assert.Nil(t, a, "a is not nil")
			assert.Nil(t, b, "b is not nil")
			return a == b
		},
		"is": mock.Is{},
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

func TestExpr_call_float_arg_func_with_int(t *testing.T) {
	env := map[string]any{
		"cnv": func(f float64) any {
			return f
		},
	}
	tests := []struct {
		input    string
		expected float64
	}{
		{"-1", -1.0},
		{"1+1", 2.0},
		{"+1", 1.0},
		{"1-1", 0.0},
		{"1/1", 1.0},
		{"1*1", 1.0},
		{"1^1", 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := expr.Compile(fmt.Sprintf("cnv(%s)", tt.input), expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(p, env)
			require.NoError(t, err)
			require.Equal(t, tt.expected, out)
		})
	}
}

func TestConstExpr_error_panic(t *testing.T) {
	env := map[string]any{
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
	env := map[string]any{
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
	env := map[string]any{
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
		`Foo == "Foo.String"`,
		expr.Env(mock.Env{}),
		expr.Patch(&mock.StringerPatcher{}),
	)
	require.NoError(t, err)

	output, err := expr.Run(program, mock.Env{})
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
	require.Equal(t,
		`{"from":2,"to":4,"line":1,"column":2,"message":"invalid operation: == (mismatched types int and bool)","snippet":"\n | 1 == true\n | ..^","prev":null}`,
		string(b),
	)
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
	require.Equal(t, "integer divide by zero (1:3)\n | 1 % 0\n | ..^", fileError.Error())
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
	env := map[string]any{}

	_, err := expr.Compile(`1 / (1 - 1)`, expr.Env(env))
	require.NoError(t, err)

	_, err = expr.Compile(`1 % 0`, expr.Env(env))
	require.Error(t, err)
	require.Equal(t, "integer divide by zero (1:3)\n | 1 % 0\n | ..^", err.Error())
}

func TestIssue154(t *testing.T) {
	type Data struct {
		Array  *[2]any
		Slice  *[]any
		Map    *map[string]any
		String *string
	}

	type Env struct {
		Data *Data
	}

	b := true
	i := 10
	s := "value"

	Array := [2]any{
		&b,
		&i,
	}

	Slice := []any{
		&b,
		&i,
	}

	Map := map[string]any{
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
		require.NoError(t, err, input)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.True(t, output.(bool), input)
	}
}

func TestIssue270(t *testing.T) {
	env := map[string]any{
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

type Issue346Array []Issue346Type

type Issue346Type struct {
	Bar string
}

func (i Issue346Array) Len() int {
	return len(i)
}

func TestIssue346(t *testing.T) {
	code := `Foo[0].Bar`

	env := map[string]any{
		"Foo": Issue346Array{
			{Bar: "bar"},
		},
	}
	program, err := expr.Compile(code, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, "bar", output)
}

func TestCompile_allow_to_use_interface_to_get_an_element_from_map(t *testing.T) {
	code := `{"value": "ok"}[vars.key]`
	env := map[string]any{
		"vars": map[string]any{
			"key": "value",
		},
	}

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, "ok", out)

	t.Run("with allow undefined variables", func(t *testing.T) {
		code := `{'key': 'value'}[Key]`
		env := mock.MapStringStringEnv{}
		options := []expr.Option{
			expr.AllowUndefinedVariables(),
		}

		program, err := expr.Compile(code, options...)
		assert.NoError(t, err)

		out, err := expr.Run(program, env)
		assert.NoError(t, err)
		assert.Equal(t, nil, out)
	})
}

func TestFastCall(t *testing.T) {
	env := map[string]any{
		"func": func(in any) float64 {
			return 8
		},
	}
	code := `func("8")`

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, float64(8), out)
}

func TestFastCall_OpCallFastErr(t *testing.T) {
	env := map[string]any{
		"func": func(...any) (any, error) {
			return 8, nil
		},
	}
	code := `func("8")`

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, 8, out)
}

func TestRun_custom_func_returns_an_error_as_second_arg(t *testing.T) {
	env := map[string]any{
		"semver": func(value string, cmp string) (bool, error) { return true, nil },
	}

	p, err := expr.Compile(`semver("1.2.3", "= 1.2.3")`, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(p, env)
	assert.NoError(t, err)
	assert.Equal(t, true, out)
}

func TestFunction(t *testing.T) {
	add := expr.Function(
		"add",
		func(p ...any) (any, error) {
			out := 0
			for _, each := range p {
				out += each.(int)
			}
			return out, nil
		},
		new(func(...int) int),
	)

	p, err := expr.Compile(`add() + add(1) + add(1, 2) + add(1, 2, 3) + add(1, 2, 3, 4)`, add)
	assert.NoError(t, err)

	out, err := expr.Run(p, nil)
	assert.NoError(t, err)
	assert.Equal(t, 20, out)
}

// Nil coalescing operator
func TestRun_NilCoalescingOperator(t *testing.T) {
	env := map[string]any{
		"foo": map[string]any{
			"bar": "value",
		},
	}

	t.Run("value", func(t *testing.T) {
		p, err := expr.Compile(`foo.bar ?? "default"`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, "value", out)
	})

	t.Run("default", func(t *testing.T) {
		p, err := expr.Compile(`foo.baz ?? "default"`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, "default", out)
	})

	t.Run("default with chain", func(t *testing.T) {
		p, err := expr.Compile(`foo?.bar ?? "default"`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, map[string]any{})
		assert.NoError(t, err)
		assert.Equal(t, "default", out)
	})
}

func TestEval_nil_in_maps(t *testing.T) {
	env := map[string]any{
		"m":     map[any]any{nil: "bar"},
		"empty": map[any]any{},
	}
	t.Run("nil key exists", func(t *testing.T) {
		p, err := expr.Compile(`m[nil]`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, "bar", out)
	})
	t.Run("no nil key", func(t *testing.T) {
		p, err := expr.Compile(`empty[nil]`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, nil, out)
	})
	t.Run("nil in m", func(t *testing.T) {
		p, err := expr.Compile(`nil in m`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, true, out)
	})
	t.Run("nil in empty", func(t *testing.T) {
		p, err := expr.Compile(`nil in empty`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, false, out)
	})
}

// Test the use of env keyword.  Forms env[] and env[”] are valid.
// The enclosed identifier must be in the expression env.
func TestEnv_keyword(t *testing.T) {
	env := map[string]any{
		"space test":                       "ok",
		"space_test":                       "not ok", // Seems to be some underscore substituting happening, check that.
		"Section 1-2a":                     "ok",
		`c:\ndrive\2015 Information Table`: "ok",
		"%*worst function name ever!!": func() string {
			return "ok"
		}(),
		"1":      "o",
		"2":      "k",
		"num":    10,
		"mylist": []int{1, 2, 3, 4, 5},
		"MIN": func(a, b int) int {
			if a < b {
				return a
			} else {
				return b
			}
		},
		"red":   "n",
		"irect": "um",
		"String Map": map[string]string{
			"one":   "two",
			"three": "four",
		},
		"OtherMap": map[string]string{
			"a": "b",
			"c": "d",
		},
	}

	// No error cases
	var tests = []struct {
		code string
		want any
	}{
		{"$env['space test']", "ok"},
		{"$env['Section 1-2a']", "ok"},
		{`$env["c:\\ndrive\\2015 Information Table"]`, "ok"},
		{"$env['%*worst function name ever!!']", "ok"},
		{"$env['String Map'].one", "two"},
		{"$env['1'] + $env['2']", "ok"},
		{"1 + $env['num'] + $env['num']", 21},
		{"MIN($env['num'],0)", 0},
		{"$env['nu' + 'm']", 10},
		{"$env[red + irect]", 10},
		{"$env['String Map']?.five", ""},
		{"$env.red", "n"},
		{"$env?.unknown", nil},
		{"$env.mylist[1]", 2},
		{"$env?.OtherMap?.a", "b"},
		{"$env?.OtherMap?.d", ""},
		{"'num' in $env", true},
		{"get($env, 'num')", 10},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {

			program, err := expr.Compile(tt.code, expr.Env(env))
			require.NoError(t, err, "compile error")

			got, err := expr.Run(program, env)
			require.NoError(t, err, "execution error")

			assert.Equal(t, tt.want, got, tt.code)
		})
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, env)
			require.NoError(t, err, "eval error: "+tt.code)

			assert.Equal(t, tt.want, got, "eval: "+tt.code)
		})
	}

	// error cases
	tests = []struct {
		code string
		want any
	}{
		{"env()", "bad"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			_, err := expr.Eval(tt.code, expr.Env(env))
			require.Error(t, err, "compile error")

		})
	}
}

func TestEnv_keyword_with_custom_functions(t *testing.T) {
	fn := expr.Function("fn", func(params ...any) (any, error) {
		return "ok", nil
	})

	var tests = []struct {
		code  string
		error bool
	}{
		{`fn()`, false},
		{`$env.fn()`, true},
		{`$env["fn"]`, true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			_, err := expr.Compile(tt.code, expr.Env(mock.Env{}), fn)
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIssue401(t *testing.T) {
	program, err := expr.Compile("(a - b + c) / d", expr.AllowUndefinedVariables())
	require.NoError(t, err, "compile error")

	output, err := expr.Run(program, map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	})
	require.NoError(t, err, "run error")
	require.Equal(t, 0.5, output)
}

func TestEval_slices_out_of_bound(t *testing.T) {
	tests := []struct {
		code string
		want any
	}{
		{"[1, 2, 3][:99]", []any{1, 2, 3}},
		{"[1, 2, 3][99:]", []any{}},
		{"[1, 2, 3][:-99]", []any{}},
		{"[1, 2, 3][-99:]", []any{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, nil)
			require.NoError(t, err, "eval error: "+tt.code)
			assert.Equal(t, tt.want, got, "eval: "+tt.code)
		})
	}
}

func TestMemoryBudget(t *testing.T) {
	tests := []struct {
		code string
	}{
		{`map(1..100, {map(1..100, {map(1..100, {0})})})`},
		{`len(1..10000000)`},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code)
			require.NoError(t, err, "compile error")

			_, err = expr.Run(program, nil)
			assert.Error(t, err, "run error")
			assert.Contains(t, err.Error(), "memory budget exceeded")
		})
	}
}

func TestExpr_custom_tests(t *testing.T) {
	f, err := os.Open("custom_tests.json")
	if os.IsNotExist(err) {
		t.Skip("no custom tests")
		return
	}

	require.NoError(t, err, "open file error")
	defer f.Close()

	var tests []string
	err = json.NewDecoder(f).Decode(&tests)
	require.NoError(t, err, "decode json error")

	for id, tt := range tests {
		t.Run(fmt.Sprintf("line %v", id+2), func(t *testing.T) {
			program, err := expr.Compile(tt)
			require.NoError(t, err)

			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(time.Second)
				timeout <- true
			}()

			done := make(chan bool, 1)
			go func() {
				out, err := expr.Run(program, nil)
				// Make sure out is used.
				_ = fmt.Sprintf("%v", out)
				assert.Error(t, err)
				done <- true
			}()

			select {
			case <-done:
				// Success.
			case <-timeout:
				t.Fatal("timeout")
			}
		})
	}
}

func TestIssue432(t *testing.T) {
	env := map[string]any{
		"func": func(
			paramUint32 uint32,
			paramUint16 uint16,
			paramUint8 uint8,
			paramUint uint,
			paramInt32 int32,
			paramInt16 int16,
			paramInt8 int8,
			paramInt int,
			paramFloat64 float64,
			paramFloat32 float32,
		) float64 {
			return float64(paramUint32) + float64(paramUint16) + float64(paramUint8) + float64(paramUint) +
				float64(paramInt32) + float64(paramInt16) + float64(paramInt8) + float64(paramInt) +
				float64(paramFloat64) + float64(paramFloat32)
		},
	}
	code := `func(1,1,1,1,1,1,1,1,1,1)`

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), out)
}

func TestIssue461(t *testing.T) {
	type EnvStr string
	type EnvField struct {
		S   EnvStr
		Str string
	}
	type Env struct {
		S        EnvStr
		Str      string
		EnvField EnvField
	}
	var tests = []struct {
		input string
		env   Env
		want  bool
	}{
		{
			input: "Str == S",
			env:   Env{S: "string", Str: "string"},
			want:  false,
		},
		{
			input: "Str == Str",
			env:   Env{Str: "string"},
			want:  true,
		},
		{
			input: "S == S",
			env:   Env{Str: "string"},
			want:  true,
		},
		{
			input: `Str == "string"`,
			env:   Env{Str: "string"},
			want:  true,
		},
		{
			input: `S == "string"`,
			env:   Env{Str: "string"},
			want:  false,
		},
		{
			input: "EnvField.Str == EnvField.S",
			env:   Env{EnvField: EnvField{S: "string", Str: "string"}},
			want:  false,
		},
		{
			input: "EnvField.Str == EnvField.Str",
			env:   Env{EnvField: EnvField{Str: "string"}},
			want:  true,
		},
		{
			input: "EnvField.S == EnvField.S",
			env:   Env{EnvField: EnvField{Str: "string"}},
			want:  true,
		},
		{
			input: `EnvField.Str == "string"`,
			env:   Env{EnvField: EnvField{Str: "string"}},
			want:  true,
		},
		{
			input: `EnvField.S == "string"`,
			env:   Env{EnvField: EnvField{Str: "string"}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := expr.Compile(tt.input, expr.Env(tt.env), expr.AsBool())

			out, err := expr.Run(program, tt.env)
			require.NoError(t, err)

			require.Equal(t, tt.want, out)
		})
	}
}

func TestIssue462(t *testing.T) {
	env := map[string]any{
		"foo": func() (string, error) {
			return "bar", nil
		},
	}
	_, err := expr.Compile(`$env.unknown(int())`, expr.Env(env))
	require.Error(t, err)
}

func TestIssue_embedded_pointer_struct(t *testing.T) {
	var tests = []struct {
		input string
		env   mock.Env
		want  any
	}{
		{
			input: "EmbedPointerEmbedInt > 0",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 123,
					},
				},
			},
			want: true,
		},
		{
			input: "(Embed).EmbedPointerEmbedInt > 0",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 123,
					},
				},
			},
			want: true,
		},
		{
			input: "(Embed).EmbedPointerEmbedInt > 0",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 0,
					},
				},
			},
			want: false,
		},
		{
			input: "(Embed).EmbedPointerEmbedMethod(0)",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 0,
					},
				},
			},
			want: "",
		},
		{
			input: "(Embed).EmbedPointerEmbedPointerReceiverMethod(0)",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: nil,
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := expr.Compile(tt.input, expr.Env(tt.env))
			require.NoError(t, err)

			out, err := expr.Run(program, tt.env)
			require.NoError(t, err)

			require.Equal(t, tt.want, out)
		})
	}
}

func TestIssue474(t *testing.T) {
	testCases := []struct {
		code string
		fail bool
	}{
		{
			code: `func("invalid")`,
			fail: true,
		},
		{
			code: `func(true)`,
			fail: true,
		},
		{
			code: `func([])`,
			fail: true,
		},
		{
			code: `func({})`,
			fail: true,
		},
		{
			code: `func(1)`,
			fail: false,
		},
		{
			code: `func(1.5)`,
			fail: false,
		},
	}

	for _, tc := range testCases {
		ltc := tc
		t.Run(ltc.code, func(t *testing.T) {
			t.Parallel()
			function := expr.Function("func", func(params ...any) (any, error) {
				return true, nil
			}, new(func(float64) bool))
			_, err := expr.Compile(ltc.code, function)
			if ltc.fail {
				if err == nil {
					t.Error("expected an error, but it was nil")
					t.FailNow()
				}
			} else {
				if err != nil {
					t.Errorf("expected nil, but it was %v", err)
					t.FailNow()
				}
			}
		})
	}
}

func TestRaceCondition_variables(t *testing.T) {
	program, err := expr.Compile(`let foo = 1; foo + 1`, expr.Env(mock.Env{}))
	require.NoError(t, err)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := expr.Run(program, mock.Env{})
			require.NoError(t, err)
			require.Equal(t, 2, out)
		}()
	}

	wg.Wait()
}

func TestOperatorDependsOnEnv(t *testing.T) {
	env := map[string]any{
		"plus": func(a, b int) int {
			return 42
		},
	}
	program, err := expr.Compile(`1 + 2`, expr.Operator("+", "plus"), expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, 42, out)
}

func TestIssue624(t *testing.T) {
	type tag struct {
		Name string
	}

	type item struct {
		Tags []tag
	}

	i := item{
		Tags: []tag{
			{Name: "one"},
			{Name: "two"},
		},
	}

	rule := `[
true && true, 
one(Tags, .Name in ["one"]), 
one(Tags, .Name in ["two"]), 
one(Tags, .Name in ["one"]) && one(Tags, .Name in ["two"])
]`
	resp, err := expr.Eval(rule, i)
	require.NoError(t, err)
	require.Equal(t, []interface{}{true, true, true, true}, resp)
}

func TestPredicateCombination(t *testing.T) {
	tests := []struct {
		code1 string
		code2 string
	}{
		{"all(1..3, {# > 0}) && all(1..3, {# < 4})", "all(1..3, {# > 0 && # < 4})"},
		{"all(1..3, {# > 1}) && all(1..3, {# < 4})", "all(1..3, {# > 1 && # < 4})"},
		{"all(1..3, {# > 0}) && all(1..3, {# < 2})", "all(1..3, {# > 0 && # < 2})"},
		{"all(1..3, {# > 1}) && all(1..3, {# < 2})", "all(1..3, {# > 1 && # < 2})"},

		{"any(1..3, {# > 0}) || any(1..3, {# < 4})", "any(1..3, {# > 0 || # < 4})"},
		{"any(1..3, {# > 1}) || any(1..3, {# < 4})", "any(1..3, {# > 1 || # < 4})"},
		{"any(1..3, {# > 0}) || any(1..3, {# < 2})", "any(1..3, {# > 0 || # < 2})"},
		{"any(1..3, {# > 1}) || any(1..3, {# < 2})", "any(1..3, {# > 1 || # < 2})"},

		{"none(1..3, {# > 0}) && none(1..3, {# < 4})", "none(1..3, {# > 0 || # < 4})"},
		{"none(1..3, {# > 1}) && none(1..3, {# < 4})", "none(1..3, {# > 1 || # < 4})"},
		{"none(1..3, {# > 0}) && none(1..3, {# < 2})", "none(1..3, {# > 0 || # < 2})"},
		{"none(1..3, {# > 1}) && none(1..3, {# < 2})", "none(1..3, {# > 1 || # < 2})"},
	}
	for _, tt := range tests {
		t.Run(tt.code1, func(t *testing.T) {
			out1, err := expr.Eval(tt.code1, nil)
			require.NoError(t, err)

			out2, err := expr.Eval(tt.code2, nil)
			require.NoError(t, err)

			require.Equal(t, out1, out2)
		})
	}
}

func TestArrayComparison(t *testing.T) {
	tests := []struct {
		env  any
		code string
	}{
		{[]string{"A", "B"}, "foo == ['A', 'B']"},
		{[]int{1, 2}, "foo == [1, 2]"},
		{[]uint8{1, 2}, "foo == [1, 2]"},
		{[]float64{1.1, 2.2}, "foo == [1.1, 2.2]"},
		{[]any{"A", 1, 1.1, true}, "foo == ['A', 1, 1.1, true]"},
		{[]string{"A", "B"}, "foo != [1, 2]"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			env := map[string]any{"foo": tt.env}
			program, err := expr.Compile(tt.code, expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, true, out)
		})
	}
}

func TestIssue_570(t *testing.T) {
	type Student struct {
		Name string
	}

	env := map[string]any{
		"student": (*Student)(nil),
	}

	program, err := expr.Compile("student?.Name", expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.IsType(t, nil, out)
}

func TestIssue_integer_truncated_by_compiler(t *testing.T) {
	env := map[string]any{
		"fn": func(x byte) byte {
			return x
		},
	}

	_, err := expr.Compile("fn(255)", expr.Env(env))
	require.NoError(t, err)

	_, err = expr.Compile("fn(256)", expr.Env(env))
	require.Error(t, err)
}

func TestExpr_crash(t *testing.T) {
	content, err := os.ReadFile("testdata/crash.txt")
	require.NoError(t, err)

	_, err = expr.Compile(string(content))
	require.Error(t, err)
}

func TestExpr_nil_op_str(t *testing.T) {
	// Let's test operators, which do `.(string)` in VM, also check for nil.

	var str *string = nil
	env := map[string]any{
		"nilString": str,
	}

	tests := []struct{ code string }{
		{`nilString == "str"`},
		{`nilString contains "str"`},
		{`nilString matches "str"`},
		{`nilString startsWith "str"`},
		{`nilString endsWith "str"`},

		{`"str" == nilString`},
		{`"str" contains nilString`},
		{`"str" matches nilString`},
		{`"str" startsWith nilString`},
		{`"str" endsWith nilString`},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code)
			require.NoError(t, err)

			output, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, false, output)
		})
	}
}

func TestExpr_env_types_map(t *testing.T) {
	envTypes := types.Map{
		"foo": types.Map{
			"bar": types.String,
		},
	}

	program, err := expr.Compile(`foo.bar`, expr.Env(envTypes))
	require.NoError(t, err)

	env := map[string]any{
		"foo": map[string]any{
			"bar": "value",
		},
	}

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, "value", output)
}

func TestExpr_env_types_map_error(t *testing.T) {
	envTypes := types.Map{
		"foo": types.Map{
			"bar": types.String,
		},
	}

	program, err := expr.Compile(`foo.bar`, expr.Env(envTypes))
	require.NoError(t, err)

	_, err = expr.Run(program, envTypes)
	require.Error(t, err)
}
