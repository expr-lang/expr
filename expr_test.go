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

	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/types"
	"github.com/expr-lang/expr/vm"

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

	// Output: integer divide by zero (1:14)
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

func ExampleOperator_with_decimal() {
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

func ExampleTimezone() {
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
			`len('北京')`,
			2,
		},
		{
			`len('👍🏻')`, // one grapheme cluster, two code points
			2,
		},
		{
			`len('👍')`, // one grapheme cluster, one code point
			1,
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
			`duration("1h") - duration("1m")`,
			time.Hour - time.Minute,
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
		{
			`if 1 > 2 { 333 * 2 + 1 } else { 444 }`,
			444,
		},
		{
			`let a = 3;
			let b = 2;
			if a>b {let c = Add(a, b); c+1} else {Add(10, b)}
			`,
			6,
		},
		{
			`if "a" < "b" {let x = "a"; x} else {"abc"}`,
			"a",
		},
		{
			`1; 2; 3`,
			3,
		},
		{
			`let a = 1; Add(2, 2); let b = 2; a + b`,
			3,
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
	require.NoError(t, err)

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
	require.Equal(t, "runtime error: integer divide by zero (1:3)\n | 1 / (1 - 1)\n | ..^", err.Error())

	_, err = expr.Compile(`1 % 0`, expr.Env(env))
	require.Error(t, err)
	require.Equal(t, "runtime error: integer divide by zero (1:3)\n | 1 % 0\n | ..^", err.Error())
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

// Test the use of env keyword.  Forms env[] and env["] are valid.
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

func TestIssue758_filter_map_index(t *testing.T) {
	env := map[string]interface{}{}

	exprStr := `
        let a_map = 0..5 | filter(# % 2 == 0) | map(#index);
        let b_filter = 0..5 | filter(# % 2 == 0);
        let b_map = b_filter | map(#index);
        [a_map, b_map]
    `

	result, err := expr.Eval(exprStr, env)
	require.NoError(t, err)

	expected := []interface{}{
		[]interface{}{0, 1, 2},
		[]interface{}{0, 1, 2},
	}

	require.Equal(t, expected, result)
}

func TestExpr_wierd_cases(t *testing.T) {
	env := map[string]any{}

	_, err := expr.Compile(`A(A)`, expr.Env(env))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown name A")
}

func TestIssue785_get_nil(t *testing.T) {
	exprStrs := []string{
		`get(nil, "a")`,
		`get({}, "a")`,
		`get(nil, "a")`,
		`get({}, "a")`,
		`({} | get("a") | get("b"))`,
	}

	for _, exprStr := range exprStrs {
		t.Run("get returns nil", func(t *testing.T) {
			env := map[string]interface{}{}

			result, err := expr.Eval(exprStr, env)
			require.NoError(t, err)

			require.Equal(t, nil, result)
		})
	}
}

func TestMaxNodes(t *testing.T) {
	maxNodes := uint(100)

	code := ""
	for i := 0; i < int(maxNodes); i++ {
		code += "1; "
	}

	_, err := expr.Compile(code, expr.MaxNodes(maxNodes))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum allowed nodes")

	_, err = expr.Compile(code, expr.MaxNodes(maxNodes+1))
	require.NoError(t, err)
}

func TestMaxNodesDisabled(t *testing.T) {
	code := ""
	for i := 0; i < 2*int(conf.DefaultMaxNodes); i++ {
		code += "1; "
	}

	_, err := expr.Compile(code, expr.MaxNodes(0))
	require.NoError(t, err)
}

func TestMemoryBudget(t *testing.T) {
	tests := []struct {
		code string
		max  int
	}{
		{`map(1..100, {map(1..100, {map(1..100, {0})})})`, -1},
		{`len(1..10000000)`, -1},
		{`1..100`, 100},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code)
			require.NoError(t, err, "compile error")

			vm := vm.VM{}
			if tt.max > 0 {
				vm.MemoryBudget = uint(tt.max)
			}
			_, err = vm.Run(program, nil)
			require.Error(t, err, "run error")
			assert.Contains(t, err.Error(), "memory budget exceeded")
		})
	}
}

// Add tests for all legal and illegal pairs for arithmetic and comparison operators
func TestArithmeticAndComparisonOperators(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Addition
		{" now() - duration('1h'``) +1  ", 3, false},
		{"1 + 2", 3, false},
		{"1.5 + 2.5", 4.0, false},
		{"1 + 2.5", 3.5, false},
		{"2.5 + 1", 3.5, false},
		{"'a' + 'b'", "ab", false},
		{"'a' + 1", "a1", false},
		{"1 + 'a'", "1a", false},
		// Subtraction
		{"5 - 2", 3, false},
		{"5.5 - 2.5", 3.0, false},
		{"5 - 2.5", 2.5, false},
		{"5.5 - 2", 3.5, false},
		// Multiplication
		{"2 * 3", 6, false},
		{"2.5 * 4", 10.0, false},
		{"2 * 4.5", 9.0, false},
		{"2.5 * 4.5", 11.25, false},
		// Division
		{"8 / 2", 4.0, false},
		{"8.0 / 2", 4.0, false},
		{"8 / 2.0", 4.0, false},
		{"8.0 / 2.0", 4.0, false},
		// Modulo
		{"5 % 2", 1, false},
		{"5.5 % 2.0", 1.5, false},
		{"5 % 2.5", 0.0, false},
		{"5.5 % 2", 1.5, false},
		{"5.5 % 2.5", 0.5, false},
		// Exponentiation
		{"2 ** 3", 8.0, false},
		{"2.0 ** 3", 8.0, false},
		{"2 ** 3.0", 8.0, false},
		{"2.0 ** 3.0", 8.0, false},
		// Comparison
		{"1 < 2", true, false},
		{"2 > 1", true, false},
		{"2 <= 2", true, false},
		{"2 >= 2", true, false},
		{"2 == 2", true, false},
		{"2 != 3", true, false},
		{"1.5 < 2.5", true, false},
		{"'a' < 'b'", true, false},
		{"true == true", true, false},
		// Illegal pairs
		{"true + 1", nil, true},
		{"1 + true", nil, true},
		{"true < 1", nil, true},
		{"1 < true", nil, true},
		{"'a' - 'b'", nil, true},
		{"'a' * 2", nil, true},
		{"2 * 'a'", nil, true},
		{"'a' / 2", nil, true},
		{"2 / 'a'", nil, true},
		{"'a' % 2", nil, true},
		{"2 % 'a'", nil, true},
		{"'a' ** 2", nil, true},
		{"2 ** 'a'", nil, true},
		{"true % false", nil, true},
		{"1 % 0", nil, true},
		{"1.0 % 0.0", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equal(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func equal(a, b any) bool {
	switch a := a.(type) {
	case float64:
		if b, ok := b.(float64); ok {
			return (a-b) < 1e-9 && (b-a) < 1e-9
		}
	}
	return a == b
}

func TestStringConcatenationWithAllTypes(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// String + numeric types
		{"'hello' + 42", "hello42", false},
		{"'hello' + 42.5", "hello42.5", false},
		{"'hello' + 0", "hello0", false},
		{"'hello' + (-5)", "hello-5", false},
		{"'hello' + 3.14159", "hello3.14159", false},

		// Numeric types + string
		{"42 + 'world'", "42world", false},
		{"42.5 + 'world'", "42.5world", false},
		{"0 + 'world'", "0world", false},
		{"(-5) + 'world'", "-5world", false},
		{"3.14159 + 'world'", "3.14159world", false},

		// String + boolean
		{"'hello' + true", "hellotrue", false},
		{"'hello' + false", "hellofalse", false},
		{"true + 'world'", "trueworld", false},
		{"false + 'world'", "falseworld", false},

		// String + nil (if supported)
		{"'hello' + nil", "hello", false},
		{"nil + 'world'", "world", false},

		// Empty string concatenations
		{"'' + 42", "42", false},
		{"42 + ''", "42", false},
		{"'' + ''", "", false},

		// Multiple concatenations
		{"'a' + 1 + 'b' + 2", "a1b2", false},
		{"'result: ' + (5 + 3)", "result: 8", false},
		{"'pi is approximately ' + 3.14", "pi is approximately 3.14", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalStrings(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNumericTypeArithmetic(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Int + Int variations
		{"1 + 2", 3, false},
		{"10 - 3", 7, false},
		{"4 * 5", 20, false},
		{"15 / 3", 5.0, false}, // Division always returns float
		{"17 % 5", 2, false},
		{"2 ** 3", 8.0, false}, // Power always returns float

		// Float + Float
		{"1.5 + 2.5", 4.0, false},
		{"5.5 - 2.5", 3.0, false},
		{"2.5 * 4.0", 10.0, false},
		{"7.5 / 2.5", 3.0, false},
		{"5.5 % 2.0", 1.5, false},
		{"2.0 ** 3.0", 8.0, false},

		// Int + Float combinations
		{"1 + 2.5", 3.5, false},
		{"5 - 2.5", 2.5, false},
		{"3 * 2.5", 7.5, false},
		{"7 / 2.0", 3.5, false},
		{"5 % 2.0", 1.0, false},
		{"2 ** 3.0", 8.0, false},

		// Float + Int combinations
		{"1.5 + 2", 3.5, false},
		{"5.5 - 2", 3.5, false},
		{"2.5 * 3", 7.5, false},
		{"7.5 / 2", 3.75, false},
		{"5.5 % 2", 1.5, false},
		{"2.0 ** 3", 8.0, false},

		// Zero operations
		{"0 + 5", 5, false},
		{"5 + 0", 5, false},
		{"0 - 5", -5, false},
		{"5 - 0", 5, false},
		{"0 * 5", 0, false},
		{"5 * 0", 0, false},
		{"0 / 5", 0.0, false},
		{"0 % 5", 0, false},
		{"0 ** 5", 0.0, false},

		// Negative numbers
		{"(-5) + 3", -2, false},
		{"3 + (-5)", -2, false},
		{"(-5) - 3", -8, false},
		{"3 - (-5)", 8, false},
		{"(-5) * 3", -15, false},
		{"3 * (-5)", -15, false},
		{"(-6) / 2", -3.0, false},
		{"(-7) % 3", -1, false},
		{"(-2) ** 3", -8.0, false},

		// Large numbers
		{"1000000 + 2000000", 3000000, false},
		{"1000000.5 + 2000000.5", 3000000.5 + 0.5, false},
		{"999999 * 2", 1999998, false},

		// Small decimal numbers
		{"0.1 + 0.2", 0.3, false},
		{"0.5 - 0.3", 0.2, false},
		{"0.1 * 0.2", 0.02, false},
		{"0.6 / 0.2", 3.0, false},

		// Error cases
		{"5 / 0", nil, true},     // Division by zero
		{"5.5 / 0.0", nil, true}, // Float division by zero
		{"5 % 0", nil, true},     // Modulo by zero
		{"5.5 % 0.0", nil, true}, // Float modulo by zero
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalNumbers(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestMixedTypeComparisons(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Int vs Float comparisons
		{"1 == 1.0", true, false},
		{"1 != 1.0", false, false},
		{"1 < 1.5", true, false},
		{"1.5 > 1", true, false},
		{"2 <= 2.0", true, false},
		{"2.0 >= 2", true, false},

		// String comparisons
		{"'abc' == 'abc'", true, false},
		{"'abc' != 'def'", true, false},
		{"'abc' < 'def'", true, false},
		{"'def' > 'abc'", true, false},
		{"'abc' <= 'abc'", true, false},
		{"'abc' >= 'abc'", true, false},

		// Zero comparisons
		{"0 == 0.0", true, false},
		{"0 != 0.0", false, false},
		{"0 < 0.1", true, false},
		{"0.0 <= 0", true, false},

		// Negative number comparisons
		{"(-1) == (-1.0)", true, false},
		{"(-1) < 0", true, false},
		{"(-1.5) < (-1)", true, false},
		{"0 > (-1)", true, false},

		// Boolean comparisons
		{"true == true", true, false},
		{"false == false", true, false},
		{"true != false", true, false},
		{"false != true", true, false},

		// Error cases - comparing incompatible types
		{"'abc' == 123", false, false}, // Should return false, not error
		{"true == 1", false, false},    // Should return false, not error
		{"'abc' < 123", nil, true},     // This should error
		{"true < 1", nil, true},        // This should error
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if result != tt.want {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Mixed arithmetic and string operations
		{"'Result: ' + (10 + 5)", "Result: 15", false},
		{"'Value: ' + (3.14 * 2)", "Value: 6.28", false},
		{"'Answer: ' + (100 / 4)", "Answer: 25", false},

		// Parentheses and order of operations
		{"(1 + 2) * 3", 9, false},
		{"1 + (2 * 3)", 7, false},
		{"(10 - 5) / (3 - 1)", 2.5, false},
		{"2 ** (3 + 1)", 16.0, false},

		// Multiple string concatenations with calculations
		{"'a' + 1 + 'b' + (2 * 3)", "a1b6", false},
		{"(5 + 3) + 'items'", "8items", false},

		// Complex comparisons
		{"(1 + 2) == 3", true, false},
		{"(1.5 * 2) > 2", true, false},
		{"'hello' + 'world' == 'helloworld'", true, false},

		// Nested operations
		{"((1 + 2) * 3) + 4", 13, false},
		{"1 + 2 * 3 + 4", 11, false}, // 1 + (2*3) + 4 = 1 + 6 + 4 = 11
		{"(1 + 2) * (3 + 4)", 21, false},

		// String with complex numeric expressions
		{"'Result: ' + ((10 + 5) * 2)", "Result: 30", false},
		{"'Pi doubled: ' + (3.14159 * 2)", "Pi doubled: 6.28318", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

// Helper functions for comparisons
func equalStrings(a, b any) bool {
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if aOk && bOk {
		return aStr == bStr
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func equalNumbers(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Convert both to float64 for comparison
	aFloat := toFloat64(a)
	bFloat := toFloat64(b)

	// Handle integer comparisons
	if isInteger(a) && isInteger(b) {
		return int64(aFloat) == int64(bFloat)
	}

	// Handle float comparisons with tolerance
	diff := aFloat - bFloat
	if diff < 0 {
		diff = -diff
	}
	return diff < 1e-9
}

func equalValues(a, b any) bool {
	// Try string comparison first
	if equalStrings(a, b) {
		return true
	}
	// Try numeric comparison
	if equalNumbers(a, b) {
		return true
	}
	// Fall back to direct comparison
	return a == b
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	default:
		return 0
	}
}

func isInteger(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

func TestNilArithmeticOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":    nil,
		"intValue":    42,
		"floatValue":  3.14,
		"stringValue": "hello",
		"boolValue":   true,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil + other types
		{"nilValue + 5", "5", false},                 // nil + int becomes string concatenation
		{"nilValue + 5.5", "5.5", false},             // nil + float becomes string concatenation
		{"nilValue + 'world'", "<nil>world", false},  // nil + string
		{"nilValue + true", "<nil>true", false},      // nil + bool
		{"nilValue + nilValue", "<nil><nil>", false}, // nil + nil

		// other types + nil
		{"5 + nilValue", "5", false},                // int + nil becomes string concatenation
		{"5.5 + nilValue", "5.5", false},            // float + nil becomes string concatenation
		{"'hello' + nilValue", "hello<nil>", false}, // string + nil
		{"true + nilValue", "true<nil>", false},     // bool + nil

		// nil - other types
		{"nilValue - 5", -5, false},       // nil - int = 0 - int
		{"nilValue - 5.5", -5.5, false},   // nil - float = 0 - float
		{"nilValue - nilValue", 0, false}, // nil - nil = 0

		// other types - nil
		{"5 - nilValue", 5, false},     // int - nil = int - 0
		{"5.5 - nilValue", 5.5, false}, // float - nil = float - 0

		// nil * other types
		{"nilValue * 5", 0, false}, // nil * anything = 0
		{"nilValue * 5.5", 0, false},
		{"nilValue * nilValue", 0, false},

		// other types * nil
		{"5 * nilValue", 0, false}, // anything * nil = 0
		{"5.5 * nilValue", 0, false},

		// nil / other types
		{"nilValue / 5", 0.0, false}, // nil / anything = 0
		{"nilValue / 5.5", 0.0, false},

		// other types / nil
		{"5 / nilValue", 0.0, false}, // anything / nil = 0 (treated as 0, not error)
		{"5.5 / nilValue", 0.0, false},

		// nil % other types
		{"nilValue % 5", 0, false}, // nil % anything = 0
		{"nilValue % 5.5", 0, false},

		// other types % nil
		{"5 % nilValue", 0, false}, // anything % nil = 0
		{"5.5 % nilValue", 0, false},

		// nil power operations
		{"nilValue ** 2", 0.0, false}, // nil ** anything = 0
		{"2 ** nilValue", 1.0, false}, // anything ** nil = anything ** 0 = 1
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestNilComparisonOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":    nil,
		"intValue":    42,
		"floatValue":  3.14,
		"stringValue": "hello",
		"boolValue":   true,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil == comparisons
		{"nilValue == nilValue", true, false}, // nil == nil
		{"nilValue == 0", false, false},       // nil != 0
		{"nilValue == 0.0", false, false},     // nil != 0.0
		{"nilValue == ''", false, false},      // nil != empty string
		{"nilValue == false", false, false},   // nil != false

		// other types == nil
		{"0 == nilValue", false, false},
		{"0.0 == nilValue", false, false},
		{"'' == nilValue", false, false},
		{"false == nilValue", false, false},

		// nil != comparisons
		{"nilValue != nilValue", false, false}, // nil is equal to nil
		{"nilValue != 0", true, false},         // nil is not equal to 0
		{"nilValue != 'hello'", true, false},   // nil is not equal to string
		{"nilValue != true", true, false},      // nil is not equal to true

		// nil < comparisons (nil is less than any non-nil value)
		{"nilValue < 5", true, false},
		{"nilValue < 0", true, false},
		{"nilValue < (-5)", true, false},
		{"nilValue < 'a'", true, false},
		{"nilValue < nilValue", false, false}, // nil is not less than nil

		// other types < nil
		{"5 < nilValue", false, false}, // non-nil is not less than nil
		{"0 < nilValue", false, false},
		{"'a' < nilValue", false, false},

		// nil > comparisons (nil is not greater than anything)
		{"nilValue > 5", false, false},
		{"nilValue > 0", false, false},
		{"nilValue > (-5)", false, false},
		{"nilValue > nilValue", false, false},

		// other types > nil (non-nil is greater than nil)
		{"5 > nilValue", true, false},
		{"0 > nilValue", true, false},
		{"'a' > nilValue", true, false},

		// nil <= comparisons
		{"nilValue <= nilValue", true, false}, // nil <= nil is true
		{"nilValue <= 5", true, false},        // nil <= anything is true
		{"nilValue <= 0", true, false},
		{"nilValue <= (-5)", true, false},

		// nil >= comparisons
		{"nilValue >= nilValue", true, false}, // nil >= nil is true
		{"nilValue >= 5", false, false},       // nil >= non-nil is false
		{"nilValue >= 0", false, false},

		// other types >= nil (non-nil >= nil is true)
		{"5 >= nilValue", true, false},
		{"0 >= nilValue", true, false},
		{"(-5) >= nilValue", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if result != tt.want {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNilLogicalOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":   nil,
		"trueValue":  true,
		"falseValue": false,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil && operations (nil is falsy)
		{"nilValue && true", nil, false},
		{"nilValue && false", nil, false},
		{"nilValue && nilValue", nil, false},

		// other types && nil
		{"true && nilValue", nil, false},    // true && nil = false (nil is falsy)
		{"false && nilValue", false, false}, // false && nil = false

		// nil || operations (nil is falsy)
		{"nilValue || true", true, false},    // nil || true = true
		{"nilValue || false", false, false},  // nil || false = false
		{"nilValue || nilValue", nil, false}, // nil || nil = false

		// other types || nil
		{"true || nilValue", true, false}, // true || nil = true
		{"false || nilValue", nil, false}, // false || nil = false

		// Complex logical operations with nil
		{"nilValue && true || false", false, false},
		{"true || nilValue && false", true, false},
		{"(nilValue || false) && true", false, false},
		{"(nilValue || true) && false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if result != tt.want {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNilInCollectionOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":        nil,
		"arrayWithNil":    []any{1, nil, "hello", nil, 5},
		"arrayWithoutNil": []any{1, 2, "hello", 4, 5},
		"mapWithNil": map[string]any{
			"key1": "value1",
			"key2": nil,
			"key3": 42,
		},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil in array operations
		{"nilValue in arrayWithNil", true, false},     // nil is in the array
		{"nilValue in arrayWithoutNil", false, false}, // nil is not in the array

		// nil in map operations (checking if nil is a value)
		{"nilValue in mapWithNil", false, false}, // nil as key doesn't exist

		// Array/slice operations with nil
		{"len(arrayWithNil)", 5, false}, // length includes nil elements
		{"arrayWithNil[1]", nil, false}, // accessing nil element

		// Map operations with nil values
		{"mapWithNil['key2']", nil, false},              // accessing nil value in map
		{"mapWithNil['key2'] == nilValue", true, false}, // comparing nil values
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNilUnaryOperations(t *testing.T) {
	env := map[string]any{
		"nilValue": nil,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Unary operations on nil
		{"!nilValue", true, false},    // !nil = true (nil is falsy)
		{"not nilValue", true, false}, // not nil = true
		{"-nilValue", 0.0, false},     // -nil = 0 (treated as numeric 0)
		{"+nilValue", nil, false},     // +nil = 0 (treated as numeric 0)
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestComplexNilExpressions(t *testing.T) {
	env := map[string]any{
		"nilValue": nil,
		"num":      10,
		"str":      "test",
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Complex expressions involving nil
		{"nilValue ?? 'default'", "default", false},               // Nil coalescing
		{"num ?? nilValue", 10, false},                            // Non-nil ?? nil
		{"nilValue ?? nilValue ?? 'fallback'", "fallback", false}, // Multiple nil coalescing

		// Conditional expressions with nil
		{"nilValue ? 'yes' : 'no'", "no", false},   // nil condition is falsy
		{"!nilValue ? 'yes' : 'no'", "yes", false}, // !nil is truthy

		// Parentheses with nil operations
		{"(nilValue + 5) * 2", "10", false}, // String result: "5" * 2 might not work as expected
		{"2 * (nilValue + 5)", "10", false}, // Check string concatenation behavior

		// Mixed nil and non-nil in complex expressions
		{"nilValue == nil && num > 5", true, false},
		{"nilValue != nil || str == 'test'", true, false},
		{"(nilValue || false) && (num > 0)", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

// Comprehensive test cases for operations between different data types
// Similar to JavaScript's type coercion behavior

func TestCrossTypeArithmeticOperations(t *testing.T) {
	env := map[string]any{
		"nil":    nil,
		"true":   true,
		"false":  false,
		"zero":   0,
		"one":    1,
		"neg":    -5,
		"float":  3.14,
		"empty":  "",
		"str":    "hello",
		"numStr": "42",
		"array":  []any{1, 2, 3},
		"obj":    map[string]any{"key": "value"},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ ADDITION (+) ============
		// JavaScript-like behavior: if either operand is string, concatenate; otherwise add numerically

		// nil + various types
		{"nil + nil", 0, false, "nil + nil should be 0"},
		{"nil + true", 1, false, "nil + true: nil becomes 0, true becomes 1"},
		{"nil + false", 0, false, "nil + false: both become 0"},
		{"nil + zero", 0, false, "nil + 0: both are 0"},
		{"nil + one", 1, false, "nil + 1: nil becomes 0"},
		{"nil + neg", -5, false, "nil + (-5): nil becomes 0"},
		{"nil + float", 3.14, false, "nil + 3.14: nil becomes 0"},
		{"nil + empty", "", false, "nil + empty string: nil becomes empty, string concatenation"},
		{"nil + str", "hello", false, "nil + string: nil becomes empty, string concatenation"},
		{"nil + numStr", "42", false, "nil + numeric string: string concatenation"},

		// boolean + various types
		{"true + true", 2, false, "true + true: both become 1"},
		{"true + false", 1, false, "true + false: 1 + 0"},
		{"false + false", 0, false, "false + false: 0 + 0"},
		{"true + zero", 1, false, "true + 0: true becomes 1"},
		{"true + one", 2, false, "true + 1: true becomes 1"},
		{"true + neg", -4, false, "true + (-5): 1 + (-5)"},
		{"true + float", 4.14, false, "true + 3.14: 1 + 3.14"},
		{"true + empty", "true", false, "true + empty string: string concatenation"},
		{"true + str", "truehello", false, "true + string: string concatenation"},
		{"true + numStr", "true42", false, "true + numeric string: string concatenation"},
		{"false + zero", 0, false, "false + 0: both are 0"},
		{"false + empty", "false", false, "false + empty string: string concatenation"},

		// number + various types
		{"zero + empty", "0", false, "0 + empty string: string concatenation"},
		{"zero + str", "0hello", false, "0 + string: string concatenation"},
		{"zero + numStr", "042", false, "0 + numeric string: string concatenation"},
		{"one + empty", "1", false, "1 + empty string: string concatenation"},
		{"one + str", "1hello", false, "1 + string: string concatenation"},
		{"one + numStr", "142", false, "1 + numeric string: string concatenation"},
		{"neg + empty", "-5", false, "(-5) + empty string: string concatenation"},
		{"neg + str", "-5hello", false, "(-5) + string: string concatenation"},
		{"float + empty", "3.14", false, "3.14 + empty string: string concatenation"},
		{"float + str", "3.14hello", false, "3.14 + string: string concatenation"},

		// string + various types (reverse of above)
		{"empty + nil", "", false, "empty string + nil: nil becomes empty"},
		{"empty + true", "true", false, "empty string + true: string concatenation"},
		{"empty + false", "false", false, "empty string + false: string concatenation"},
		{"empty + zero", "0", false, "empty string + 0: string concatenation"},
		{"empty + one", "1", false, "empty string + 1: string concatenation"},
		{"empty + neg", "-5", false, "empty string + (-5): string concatenation"},
		{"empty + float", "3.14", false, "empty string + 3.14: string concatenation"},
		{"empty + empty", "", false, "empty string + empty string: concatenation"},
		{"str + nil", "hello", false, "string + nil: nil becomes empty"},
		{"str + true", "hellotrue", false, "string + true: string concatenation"},
		{"str + false", "hellofalse", false, "string + false: string concatenation"},
		{"str + zero", "hello0", false, "string + 0: string concatenation"},
		{"str + one", "hello1", false, "string + 1: string concatenation"},
		{"str + neg", "hello-5", false, "string + (-5): string concatenation"},
		{"str + float", "hello3.14", false, "string + 3.14: string concatenation"},
		{"numStr + nil", "42", false, "numeric string + nil: nil becomes empty"},
		{"numStr + true", "42true", false, "numeric string + true: string concatenation"},
		{"numStr + zero", "420", false, "numeric string + 0: string concatenation"},
		{"numStr + one", "421", false, "numeric string + 1: string concatenation"},

		// ============ SUBTRACTION (-) ============
		// All operands should be converted to numbers

		{"nil - nil", 0, false, "nil - nil: 0 - 0"},
		{"nil - true", -1, false, "nil - true: 0 - 1"},
		{"nil - false", 0, false, "nil - false: 0 - 0"},
		{"nil - zero", 0, false, "nil - 0: 0 - 0"},
		{"nil - one", -1, false, "nil - 1: 0 - 1"},
		{"nil - neg", 5, false, "nil - (-5): 0 - (-5)"},
		{"nil - float", -3.14, false, "nil - 3.14: 0 - 3.14"},
		{"nil - empty", 0, false, "nil - empty string: 0 - 0 (empty string becomes 0)"},
		{"nil - numStr", -42, false, "nil - '42': 0 - 42"},

		{"true - nil", 1, false, "true - nil: 1 - 0"},
		{"true - true", 0, false, "true - true: 1 - 1"},
		{"true - false", 1, false, "true - false: 1 - 0"},
		{"true - zero", 1, false, "true - 0: 1 - 0"},
		{"true - one", 0, false, "true - 1: 1 - 1"},
		{"true - neg", 6, false, "true - (-5): 1 - (-5)"},
		{"true - float", -2.14, false, "true - 3.14: 1 - 3.14"},
		{"true - empty", 1, false, "true - empty string: 1 - 0"},
		{"true - numStr", -41, false, "true - '42': 1 - 42"},

		{"false - nil", 0, false, "false - nil: 0 - 0"},
		{"false - true", -1, false, "false - true: 0 - 1"},
		{"false - zero", 0, false, "false - 0: 0 - 0"},
		{"false - empty", 0, false, "false - empty string: 0 - 0"},

		{"zero - nil", 0, false, "0 - nil: 0 - 0"},
		{"zero - true", -1, false, "0 - true: 0 - 1"},
		{"zero - false", 0, false, "0 - false: 0 - 0"},
		{"zero - empty", 0, false, "0 - empty string: 0 - 0"},
		{"zero - numStr", -42, false, "0 - '42': 0 - 42"},

		{"one - nil", 1, false, "1 - nil: 1 - 0"},
		{"one - true", 0, false, "1 - true: 1 - 1"},
		{"one - false", 1, false, "1 - false: 1 - 0"},
		{"one - empty", 1, false, "1 - empty string: 1 - 0"},
		{"one - numStr", -41, false, "1 - '42': 1 - 42"},

		{"neg - nil", -5, false, "(-5) - nil: (-5) - 0"},
		{"neg - true", -6, false, "(-5) - true: (-5) - 1"},
		{"neg - false", -5, false, "(-5) - false: (-5) - 0"},
		{"neg - empty", -5, false, "(-5) - empty string: (-5) - 0"},
		{"neg - numStr", -47, false, "(-5) - '42': (-5) - 42"},

		{"float - nil", 3.14, false, "3.14 - nil: 3.14 - 0"},
		{"float - true", 2.14, false, "3.14 - true: 3.14 - 1"},
		{"float - false", 3.14, false, "3.14 - false: 3.14 - 0"},
		{"float - empty", 3.14, false, "3.14 - empty string: 3.14 - 0"},
		{"float - numStr", -38.86, false, "3.14 - '42': 3.14 - 42"},

		{"empty - nil", 0, false, "empty string - nil: 0 - 0"},
		{"empty - true", -1, false, "empty string - true: 0 - 1"},
		{"empty - false", 0, false, "empty string - false: 0 - 0"},
		{"empty - zero", 0, false, "empty string - 0: 0 - 0"},
		{"empty - one", -1, false, "empty string - 1: 0 - 1"},
		{"empty - empty", 0, false, "empty string - empty string: 0 - 0"},
		{"empty - numStr", -42, false, "empty string - '42': 0 - 42"},

		{"numStr - nil", 42, false, "'42' - nil: 42 - 0"},
		{"numStr - true", 41, false, "'42' - true: 42 - 1"},
		{"numStr - false", 42, false, "'42' - false: 42 - 0"},
		{"numStr - zero", 42, false, "'42' - 0: 42 - 0"},
		{"numStr - one", 41, false, "'42' - 1: 42 - 1"},
		{"numStr - empty", 42, false, "'42' - empty string: 42 - 0"},
		{"numStr - numStr", 0, false, "'42' - '42': 42 - 42"},

		// ============ MULTIPLICATION (*) ============
		// All operands should be converted to numbers

		{"nil * nil", 0, false, "nil * nil: 0 * 0"},
		{"nil * true", 0, false, "nil * true: 0 * 1"},
		{"nil * false", 0, false, "nil * false: 0 * 0"},
		{"nil * zero", 0, false, "nil * 0: 0 * 0"},
		{"nil * one", 0, false, "nil * 1: 0 * 1"},
		{"nil * neg", 0, false, "nil * (-5): 0 * (-5)"},
		{"nil * float", 0, false, "nil * 3.14: 0 * 3.14"},
		{"nil * empty", 0, false, "nil * empty string: 0 * 0"},
		{"nil * numStr", 0, false, "nil * '42': 0 * 42"},

		{"true * nil", 0, false, "true * nil: 1 * 0"},
		{"true * true", 1, false, "true * true: 1 * 1"},
		{"true * false", 0, false, "true * false: 1 * 0"},
		{"true * zero", 0, false, "true * 0: 1 * 0"},
		{"true * one", 1, false, "true * 1: 1 * 1"},
		{"true * neg", -5, false, "true * (-5): 1 * (-5)"},
		{"true * float", 3.14, false, "true * 3.14: 1 * 3.14"},
		{"true * empty", 0, false, "true * empty string: 1 * 0"},
		{"true * numStr", 42, false, "true * '42': 1 * 42"},

		{"false * nil", 0, false, "false * nil: 0 * 0"},
		{"false * true", 0, false, "false * true: 0 * 1"},
		{"false * zero", 0, false, "false * 0: 0 * 0"},
		{"false * one", 0, false, "false * 1: 0 * 1"},
		{"false * empty", 0, false, "false * empty string: 0 * 0"},
		{"false * numStr", 0, false, "false * '42': 0 * 42"},

		{"zero * nil", 0, false, "0 * nil: 0 * 0"},
		{"zero * true", 0, false, "0 * true: 0 * 1"},
		{"zero * empty", 0, false, "0 * empty string: 0 * 0"},
		{"zero * numStr", 0, false, "0 * '42': 0 * 42"},

		{"one * nil", 0, false, "1 * nil: 1 * 0"},
		{"one * true", 1, false, "1 * true: 1 * 1"},
		{"one * false", 0, false, "1 * false: 1 * 0"},
		{"one * empty", 0, false, "1 * empty string: 1 * 0"},
		{"one * numStr", 42, false, "1 * '42': 1 * 42"},

		{"neg * nil", 0, false, "(-5) * nil: (-5) * 0"},
		{"neg * true", -5, false, "(-5) * true: (-5) * 1"},
		{"neg * false", 0, false, "(-5) * false: (-5) * 0"},
		{"neg * empty", 0, false, "(-5) * empty string: (-5) * 0"},
		{"neg * numStr", -210, false, "(-5) * '42': (-5) * 42"},

		{"float * nil", 0, false, "3.14 * nil: 3.14 * 0"},
		{"float * true", 3.14, false, "3.14 * true: 3.14 * 1"},
		{"float * false", 0, false, "3.14 * false: 3.14 * 0"},
		{"float * empty", 0, false, "3.14 * empty string: 3.14 * 0"},
		{"float * numStr", 131.88, false, "3.14 * '42': 3.14 * 42"},

		{"empty * nil", 0, false, "empty string * nil: 0 * 0"},
		{"empty * true", 0, false, "empty string * true: 0 * 1"},
		{"empty * false", 0, false, "empty string * false: 0 * 0"},
		{"empty * zero", 0, false, "empty string * 0: 0 * 0"},
		{"empty * one", 0, false, "empty string * 1: 0 * 1"},
		{"empty * empty", 0, false, "empty string * empty string: 0 * 0"},
		{"empty * numStr", 0, false, "empty string * '42': 0 * 42"},

		{"numStr * nil", 0, false, "'42' * nil: 42 * 0"},
		{"numStr * true", 42, false, "'42' * true: 42 * 1"},
		{"numStr * false", 0, false, "'42' * false: 42 * 0"},
		{"numStr * zero", 0, false, "'42' * 0: 42 * 0"},
		{"numStr * one", 42, false, "'42' * 1: 42 * 1"},
		{"numStr * empty", 0, false, "'42' * empty string: 42 * 0"},
		{"numStr * numStr", 1764, false, "'42' * '42': 42 * 42"},

		// ============ DIVISION (/) ============
		// All operands should be converted to numbers

		{"nil / one", 0.0, false, "nil / 1: 0 / 1 = 0"},
		{"nil / neg", 0.0, false, "nil / (-5): 0 / (-5) = 0"},
		{"nil / float", 0.0, false, "nil / 3.14: 0 / 3.14 = 0"},
		{"nil / numStr", 0.0, false, "nil / '42': 0 / 42 = 0"},

		{"true / one", 1.0, false, "true / 1: 1 / 1 = 1"},
		{"true / neg", -0.2, false, "true / (-5): 1 / (-5) = -0.2"},
		{"true / float", 0.318, false, "true / 3.14: 1 / 3.14 ≈ 0.318"}, // Approximate
		{"true / numStr", 0.024, false, "true / '42': 1 / 42 ≈ 0.024"},  // Approximate

		{"false / one", 0.0, false, "false / 1: 0 / 1 = 0"},
		{"false / neg", 0.0, false, "false / (-5): 0 / (-5) = 0"},
		{"false / numStr", 0.0, false, "false / '42': 0 / 42 = 0"},

		{"zero / one", 0.0, false, "0 / 1: 0 / 1 = 0"},
		{"zero / neg", 0.0, false, "0 / (-5): 0 / (-5) = 0"},
		{"zero / numStr", 0.0, false, "0 / '42': 0 / 42 = 0"},

		{"one / true", 1.0, false, "1 / true: 1 / 1 = 1"},
		{"one / neg", -0.2, false, "1 / (-5): 1 / (-5) = -0.2"},
		{"one / numStr", 0.024, false, "1 / '42': 1 / 42 ≈ 0.024"}, // Approximate

		{"neg / true", -5.0, false, "(-5) / true: (-5) / 1 = -5"},
		{"neg / one", -5.0, false, "(-5) / 1: (-5) / 1 = -5"},
		{"neg / numStr", -0.119, false, "(-5) / '42': (-5) / 42 ≈ -0.119"}, // Approximate

		{"float / true", 3.14, false, "3.14 / true: 3.14 / 1 = 3.14"},
		{"float / one", 3.14, false, "3.14 / 1: 3.14 / 1 = 3.14"},
		{"float / numStr", 0.075, false, "3.14 / '42': 3.14 / 42 ≈ 0.075"}, // Approximate

		{"numStr / true", 42.0, false, "'42' / true: 42 / 1 = 42"},
		{"numStr / one", 42.0, false, "'42' / 1: 42 / 1 = 42"},
		{"numStr / neg", -8.4, false, "'42' / (-5): 42 / (-5) = -8.4"},
		{"numStr / float", 13.375, false, "'42' / 3.14: 42 / 3.14 ≈ 13.375"}, // Approximate

		// Division by zero cases (should error or return special values)
		{"nil / nil", 0.0, true, "nil / nil: 0 / 0 should error"},
		{"nil / false", 0.0, true, "nil / false: 0 / 0 should error"},
		{"nil / zero", 0.0, true, "nil / 0: 0 / 0 should error"},
		{"nil / empty", 0.0, true, "nil / empty string: 0 / 0 should error"},
		{"true / nil", 0.0, true, "true / nil: 1 / 0 should error"},
		{"true / false", 0.0, true, "true / false: 1 / 0 should error"},
		{"true / zero", 0.0, true, "true / 0: 1 / 0 should error"},
		{"true / empty", 0.0, true, "true / empty string: 1 / 0 should error"},
		{"one / nil", 0.0, true, "1 / nil: 1 / 0 should error"},
		{"one / false", 0.0, true, "1 / false: 1 / 0 should error"},
		{"one / zero", 0.0, true, "1 / 0: 1 / 0 should error"},
		{"one / empty", 0.0, true, "1 / empty string: 1 / 0 should error"},

		// ============ MODULO (%) ============
		// All operands should be converted to numbers

		{"nil % one", 0, false, "nil % 1: 0 % 1 = 0"},
		{"nil % neg", 0, false, "nil % (-5): 0 % (-5) = 0"},
		{"nil % numStr", 0, false, "nil % '42': 0 % 42 = 0"},

		{"true % one", 0, false, "true % 1: 1 % 1 = 0"},
		{"true % neg", 1, false, "true % (-5): 1 % (-5) = 1"},
		{"true % numStr", 1, false, "true % '42': 1 % 42 = 1"},

		{"false % one", 0, false, "false % 1: 0 % 1 = 0"},
		{"false % neg", 0, false, "false % (-5): 0 % (-5) = 0"},
		{"false % numStr", 0, false, "false % '42': 0 % 42 = 0"},

		{"zero % one", 0, false, "0 % 1: 0 % 1 = 0"},
		{"zero % neg", 0, false, "0 % (-5): 0 % (-5) = 0"},
		{"zero % numStr", 0, false, "0 % '42': 0 % 42 = 0"},

		{"one % true", 0, false, "1 % true: 1 % 1 = 0"},
		{"one % neg", 1, false, "1 % (-5): 1 % (-5) = 1"},
		{"one % numStr", 1, false, "1 % '42': 1 % 42 = 1"},

		{"neg % true", 0, false, "(-5) % true: (-5) % 1 = 0"},
		{"neg % one", 0, false, "(-5) % 1: (-5) % 1 = 0"},
		{"neg % numStr", -5, false, "(-5) % '42': (-5) % 42 = -5"},

		{"numStr % true", 0, false, "'42' % true: 42 % 1 = 0"},
		{"numStr % one", 0, false, "'42' % 1: 42 % 1 = 0"},
		{"numStr % neg", 2, false, "'42' % (-5): 42 % (-5) = 2"},

		// Modulo by zero cases (should error)
		{"nil % nil", 0, true, "nil % nil: 0 % 0 should error"},
		{"nil % false", 0, true, "nil % false: 0 % 0 should error"},
		{"nil % zero", 0, true, "nil % 0: 0 % 0 should error"},
		{"nil % empty", 0, true, "nil % empty string: 0 % 0 should error"},
		{"true % nil", 0, true, "true % nil: 1 % 0 should error"},
		{"true % false", 0, true, "true % false: 1 % 0 should error"},
		{"true % zero", 0, true, "true % 0: 1 % 0 should error"},
		{"true % empty", 0, true, "true % empty string: 1 % 0 should error"},

		// ============ POWER (**) ============
		// All operands should be converted to numbers

		{"nil ** nil", 1.0, false, "nil ** nil: 0 ** 0 = 1"},
		{"nil ** true", 0.0, false, "nil ** true: 0 ** 1 = 0"},
		{"nil ** false", 1.0, false, "nil ** false: 0 ** 0 = 1"},
		{"nil ** zero", 1.0, false, "nil ** 0: 0 ** 0 = 1"},
		{"nil ** one", 0.0, false, "nil ** 1: 0 ** 1 = 0"},
		{"nil ** neg", 0.0, false, "nil ** (-5): 0 ** (-5) should be infinity or error"},
		{"nil ** numStr", 0.0, false, "nil ** '42': 0 ** 42 = 0"},

		{"true ** nil", 1.0, false, "true ** nil: 1 ** 0 = 1"},
		{"true ** true", 1.0, false, "true ** true: 1 ** 1 = 1"},
		{"true ** false", 1.0, false, "true ** false: 1 ** 0 = 1"},
		{"true ** zero", 1.0, false, "true ** 0: 1 ** 0 = 1"},
		{"true ** one", 1.0, false, "true ** 1: 1 ** 1 = 1"},
		{"true ** neg", 1.0, false, "true ** (-5): 1 ** (-5) = 1"},
		{"true ** numStr", 1.0, false, "true ** '42': 1 ** 42 = 1"},

		{"false ** nil", 1.0, false, "false ** nil: 0 ** 0 = 1"},
		{"false ** true", 0.0, false, "false ** true: 0 ** 1 = 0"},
		{"false ** zero", 1.0, false, "false ** 0: 0 ** 0 = 1"},
		{"false ** one", 0.0, false, "false ** 1: 0 ** 1 = 0"},
		{"false ** neg", 0.0, false, "false ** (-5): 0 ** (-5) should be infinity or error"},
		{"false ** numStr", 0.0, false, "false ** '42': 0 ** 42 = 0"},

		{"zero ** nil", 1.0, false, "0 ** nil: 0 ** 0 = 1"},
		{"zero ** true", 0.0, false, "0 ** true: 0 ** 1 = 0"},
		{"zero ** false", 1.0, false, "0 ** false: 0 ** 0 = 1"},
		{"zero ** one", 0.0, false, "0 ** 1: 0 ** 1 = 0"},
		{"zero ** neg", 0.0, false, "0 ** (-5): 0 ** (-5) should be infinity or error"},
		{"zero ** numStr", 0.0, false, "0 ** '42': 0 ** 42 = 0"},

		{"one ** nil", 1.0, false, "1 ** nil: 1 ** 0 = 1"},
		{"one ** true", 1.0, false, "1 ** true: 1 ** 1 = 1"},
		{"one ** false", 1.0, false, "1 ** false: 1 ** 0 = 1"},
		{"one ** zero", 1.0, false, "1 ** 0: 1 ** 0 = 1"},
		{"one ** neg", 1.0, false, "1 ** (-5): 1 ** (-5) = 1"},
		{"one ** numStr", 1.0, false, "1 ** '42': 1 ** 42 = 1"},

		{"neg ** nil", 1.0, false, "(-5) ** nil: (-5) ** 0 = 1"},
		{"neg ** true", -5.0, false, "(-5) ** true: (-5) ** 1 = -5"},
		{"neg ** false", 1.0, false, "(-5) ** false: (-5) ** 0 = 1"},
		{"neg ** zero", 1.0, false, "(-5) ** 0: (-5) ** 0 = 1"},
		{"neg ** one", -5.0, false, "(-5) ** 1: (-5) ** 1 = -5"},

		{"numStr ** nil", 1.0, false, "'42' ** nil: 42 ** 0 = 1"},
		{"numStr ** true", 42.0, false, "'42' ** true: 42 ** 1 = 42"},
		{"numStr ** false", 1.0, false, "'42' ** false: 42 ** 0 = 1"},
		{"numStr ** zero", 1.0, false, "'42' ** 0: 42 ** 0 = 1"},
		{"numStr ** one", 42.0, false, "'42' ** 1: 42 ** 1 = 42"},

		// Error cases with non-numeric strings
		{"str - nil", 0, true, "non-numeric string - nil should error"},
		{"str * nil", 0, true, "non-numeric string * nil should error"},
		{"str / nil", 0.0, true, "non-numeric string / nil should error"},
		{"str % nil", 0, true, "non-numeric string % nil should error"},
		{"str ** nil", 1.0, true, "non-numeric string ** nil should error"},
		{"nil - str", 0, true, "nil - non-numeric string should error"},
		{"nil * str", 0, true, "nil * non-numeric string should error"},
		{"nil / str", 0.0, true, "nil / non-numeric string should error"},
		{"nil % str", 0, true, "nil % non-numeric string should error"},
		{"nil ** str", 1.0, true, "nil ** non-numeric string should error"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if !approximatelyEqual(result, tt.want) {
					t.Errorf("%q (%s): got %v, want %v", tt.expr, tt.description, result, tt.want)
				}
			}
		})
	}
}

func TestCrossTypeComparisonOperations(t *testing.T) {
	env := map[string]any{
		"nil":    nil,
		"true":   true,
		"false":  false,
		"zero":   0,
		"one":    1,
		"neg":    -5,
		"float":  3.14,
		"empty":  "",
		"str":    "hello",
		"numStr": "42",
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ EQUALITY (==) ============
		// JavaScript-like type coercion for equality

		// nil equality
		{"nil == nil", true, false, "nil == nil should be true"},
		{"nil == false", false, false, "nil == false should be false (different types)"},
		{"nil == zero", false, false, "nil == 0 should be false (different types)"},
		{"nil == empty", false, false, "nil == empty string should be false (different types)"},

		// boolean equality
		{"true == true", true, false, "true == true should be true"},
		{"false == false", true, false, "false == false should be true"},
		{"true == false", false, false, "true == false should be false"},
		{"true == one", true, false, "true == 1 should be true (boolean to number coercion)"},
		{"false == zero", true, false, "false == 0 should be true (boolean to number coercion)"},
		{"true == numStr", false, false, "true == '42' should be false (1 != 42)"},
		{"false == empty", true, false, "false == empty string should be true (both falsy)"},

		// number equality
		{"zero == zero", true, false, "0 == 0 should be true"},
		{"zero == false", true, false, "0 == false should be true (number to boolean coercion)"},
		{"one == true", true, false, "1 == true should be true (number to boolean coercion)"},
		{"one == numStr", false, false, "1 == '42' should be false (1 != 42)"},
		{"zero == empty", true, false, "0 == empty string should be true (both convert to 0)"},

		// string equality
		{"empty == empty", true, false, "empty string == empty string should be true"},
		{"empty == false", true, false, "empty string == false should be true (both falsy)"},
		{"empty == zero", true, false, "empty string == 0 should be true (both convert to 0)"},
		{"numStr == numStr", true, false, "'42' == '42' should be true"},
		{"str == str", true, false, "'hello' == 'hello' should be true"},
		{"str == empty", false, false, "'hello' == empty string should be false"},

		// ============ INEQUALITY (!=) ============
		// Opposite of equality

		{"nil != nil", false, false, "nil != nil should be false"},
		{"nil != false", true, false, "nil != false should be true"},
		{"nil != zero", true, false, "nil != 0 should be true"},
		{"nil != empty", true, false, "nil != empty string should be true"},
		{"true != false", true, false, "true != false should be true"},
		{"true != one", false, false, "true != 1 should be false (they're equal)"},
		{"false != zero", false, false, "false != 0 should be false (they're equal)"},
		{"one != numStr", true, false, "1 != '42' should be true"},
		{"empty != zero", false, false, "empty string != 0 should be false (they're equal)"},

		// ============ LESS THAN (<) ============
		// Convert to comparable types

		{"nil < nil", false, false, "nil < nil should be false"},
		{"nil < true", true, false, "nil < true: 0 < 1 should be true"},
		{"nil < false", false, false, "nil < false: 0 < 0 should be false"},
		{"nil < zero", false, false, "nil < 0: 0 < 0 should be false"},
		{"nil < one", true, false, "nil < 1: 0 < 1 should be true"},
		{"nil < neg", false, false, "nil < (-5): 0 < (-5) should be false"},
		{"nil < float", true, false, "nil < 3.14: 0 < 3.14 should be true"},
		{"nil < empty", false, false, "nil < empty string: 0 < 0 should be false"},
		{"nil < numStr", true, false, "nil < '42': 0 < 42 should be true"},

		{"true < nil", false, false, "true < nil: 1 < 0 should be false"},
		{"true < true", false, false, "true < true: 1 < 1 should be false"},
		{"true < false", false, false, "true < false: 1 < 0 should be false"},
		{"true < zero", false, false, "true < 0: 1 < 0 should be false"},
		{"true < one", false, false, "true < 1: 1 < 1 should be false"},
		{"true < neg", false, false, "true < (-5): 1 < (-5) should be false"},
		{"true < float", true, false, "true < 3.14: 1 < 3.14 should be true"},
		{"true < empty", false, false, "true < empty string: 1 < 0 should be false"},
		{"true < numStr", true, false, "true < '42': 1 < 42 should be true"},

		{"false < nil", false, false, "false < nil: 0 < 0 should be false"},
		{"false < true", true, false, "false < true: 0 < 1 should be true"},
		{"false < zero", false, false, "false < 0: 0 < 0 should be false"},
		{"false < one", true, false, "false < 1: 0 < 1 should be true"},
		{"false < neg", false, false, "false < (-5): 0 < (-5) should be false"},
		{"false < float", true, false, "false < 3.14: 0 < 3.14 should be true"},
		{"false < empty", false, false, "false < empty string: 0 < 0 should be false"},
		{"false < numStr", true, false, "false < '42': 0 < 42 should be true"},

		{"zero < nil", false, false, "0 < nil: 0 < 0 should be false"},
		{"zero < true", true, false, "0 < true: 0 < 1 should be true"},
		{"zero < false", false, false, "0 < false: 0 < 0 should be false"},
		{"zero < one", true, false, "0 < 1: 0 < 1 should be true"},
		{"zero < neg", false, false, "0 < (-5): 0 < (-5) should be false"},
		{"zero < float", true, false, "0 < 3.14: 0 < 3.14 should be true"},
		{"zero < empty", false, false, "0 < empty string: 0 < 0 should be false"},
		{"zero < numStr", true, false, "0 < '42': 0 < 42 should be true"},

		{"one < nil", false, false, "1 < nil: 1 < 0 should be false"},
		{"one < true", false, false, "1 < true: 1 < 1 should be false"},
		{"one < false", false, false, "1 < false: 1 < 0 should be false"},
		{"one < zero", false, false, "1 < 0: 1 < 0 should be false"},
		{"one < neg", false, false, "1 < (-5): 1 < (-5) should be false"},
		{"one < float", true, false, "1 < 3.14: 1 < 3.14 should be true"},
		{"one < empty", false, false, "1 < empty string: 1 < 0 should be false"},
		{"one < numStr", true, false, "1 < '42': 1 < 42 should be true"},

		{"neg < nil", true, false, "(-5) < nil: (-5) < 0 should be true"},
		{"neg < true", true, false, "(-5) < true: (-5) < 1 should be true"},
		{"neg < false", true, false, "(-5) < false: (-5) < 0 should be true"},
		{"neg < zero", true, false, "(-5) < 0: (-5) < 0 should be true"},
		{"neg < one", true, false, "(-5) < 1: (-5) < 1 should be true"},
		{"neg < float", true, false, "(-5) < 3.14: (-5) < 3.14 should be true"},
		{"neg < empty", true, false, "(-5) < empty string: (-5) < 0 should be true"},
		{"neg < numStr", true, false, "(-5) < '42': (-5) < 42 should be true"},

		{"float < nil", false, false, "3.14 < nil: 3.14 < 0 should be false"},
		{"float < true", false, false, "3.14 < true: 3.14 < 1 should be false"},
		{"float < false", false, false, "3.14 < false: 3.14 < 0 should be false"},
		{"float < zero", false, false, "3.14 < 0: 3.14 < 0 should be false"},
		{"float < one", false, false, "3.14 < 1: 3.14 < 1 should be false"},
		{"float < neg", false, false, "3.14 < (-5): 3.14 < (-5) should be false"},
		{"float < empty", false, false, "3.14 < empty string: 3.14 < 0 should be false"},
		{"float < numStr", true, false, "3.14 < '42': 3.14 < 42 should be true"},

		{"empty < nil", false, false, "empty string < nil: 0 < 0 should be false"},
		{"empty < true", true, false, "empty string < true: 0 < 1 should be true"},
		{"empty < false", false, false, "empty string < false: 0 < 0 should be false"},
		{"empty < zero", false, false, "empty string < 0: 0 < 0 should be false"},
		{"empty < one", true, false, "empty string < 1: 0 < 1 should be true"},
		{"empty < neg", false, false, "empty string < (-5): 0 < (-5) should be false"},
		{"empty < float", true, false, "empty string < 3.14: 0 < 3.14 should be true"},
		{"empty < numStr", true, false, "empty string < '42': 0 < 42 should be true"},

		{"numStr < nil", false, false, "'42' < nil: 42 < 0 should be false"},
		{"numStr < true", false, false, "'42' < true: 42 < 1 should be false"},
		{"numStr < false", false, false, "'42' < false: 42 < 0 should be false"},
		{"numStr < zero", false, false, "'42' < 0: 42 < 0 should be false"},
		{"numStr < one", false, false, "'42' < 1: 42 < 1 should be false"},
		{"numStr < neg", false, false, "'42' < (-5): 42 < (-5) should be false"},
		{"numStr < float", false, false, "'42' < 3.14: 42 < 3.14 should be false"},
		{"numStr < empty", false, false, "'42' < empty string: 42 < 0 should be false"},

		// String comparisons (when both operands are strings, use lexicographical comparison)
		{"str < str", false, false, "'hello' < 'hello' should be false"},
		{"empty < str", true, false, "empty string < 'hello' should be true"},
		{"str < empty", false, false, "'hello' < empty string should be false"},

		// ============ GREATER THAN (>) ============
		// Opposite of less than

		{"nil > nil", false, false, "nil > nil should be false"},
		{"true > nil", true, false, "true > nil: 1 > 0 should be true"},
		{"false > nil", false, false, "false > nil: 0 > 0 should be false"},
		{"one > nil", true, false, "1 > nil: 1 > 0 should be true"},
		{"neg > nil", false, false, "(-5) > nil: (-5) > 0 should be false"},
		{"float > nil", true, false, "3.14 > nil: 3.14 > 0 should be true"},
		{"empty > nil", false, false, "empty string > nil: 0 > 0 should be false"},
		{"numStr > nil", true, false, "'42' > nil: 42 > 0 should be true"},

		{"nil > true", false, false, "nil > true: 0 > 1 should be false"},
		{"nil > one", false, false, "nil > 1: 0 > 1 should be false"},
		{"nil > neg", true, false, "nil > (-5): 0 > (-5) should be true"},
		{"nil > float", false, false, "nil > 3.14: 0 > 3.14 should be false"},
		{"nil > numStr", false, false, "nil > '42': 0 > 42 should be false"},

		// ============ LESS THAN OR EQUAL (<=) ============
		// Less than OR equal

		{"nil <= nil", true, false, "nil <= nil should be true"},
		{"nil <= true", true, false, "nil <= true: 0 <= 1 should be true"},
		{"nil <= zero", true, false, "nil <= 0: 0 <= 0 should be true"},
		{"nil <= one", true, false, "nil <= 1: 0 <= 1 should be true"},
		{"true <= one", true, false, "true <= 1: 1 <= 1 should be true"},
		{"false <= zero", true, false, "false <= 0: 0 <= 0 should be true"},
		{"empty <= zero", true, false, "empty string <= 0: 0 <= 0 should be true"},

		// ============ GREATER THAN OR EQUAL (>=) ============
		// Greater than OR equal

		{"nil >= nil", true, false, "nil >= nil should be true"},
		{"true >= nil", true, false, "true >= nil: 1 >= 0 should be true"},
		{"zero >= nil", true, false, "0 >= nil: 0 >= 0 should be true"},
		{"one >= nil", true, false, "1 >= nil: 1 >= 0 should be true"},
		{"one >= true", true, false, "1 >= true: 1 >= 1 should be true"},
		{"zero >= false", true, false, "0 >= false: 0 >= 0 should be true"},
		{"zero >= empty", true, false, "0 >= empty string: 0 >= 0 should be true"},

		// Error cases with non-comparable types
		{"str < nil", false, true, "non-numeric string < nil should error or use special rules"},
		{"str > nil", false, true, "non-numeric string > nil should error or use special rules"},
		{"nil < str", false, true, "nil < non-numeric string should error or use special rules"},
		{"nil > str", false, true, "nil > non-numeric string should error or use special rules"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if result != tt.want {
					t.Errorf("%q (%s): got %v, want %v", tt.expr, tt.description, result, tt.want)
				}
			}
		})
	}
}

func TestCrossTypeLogicalOperations(t *testing.T) {
	env := map[string]any{
		"nil":        nil,
		"true":       true,
		"false":      false,
		"zero":       0,
		"one":        1,
		"neg":        -5,
		"float":      3.14,
		"empty":      "",
		"str":        "hello",
		"numStr":     "42",
		"array":      []any{1, 2, 3},
		"emptyArray": []any{},
		"obj":        map[string]any{"key": "value"},
		"emptyObj":   map[string]any{},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ LOGICAL AND (&&) ============
		// In JavaScript: if left is falsy, return left; otherwise return right
		// Falsy values: false, 0, "", null, undefined, NaN

		// nil && others (nil is falsy)
		{"nil && nil", nil, false, "nil && nil should return first nil"},
		{"nil && true", nil, false, "nil && true should return nil (first falsy)"},
		{"nil && false", nil, false, "nil && false should return nil (first falsy)"},
		{"nil && zero", nil, false, "nil && 0 should return nil (first falsy)"},
		{"nil && one", nil, false, "nil && 1 should return nil (first falsy)"},
		{"nil && neg", nil, false, "nil && (-5) should return nil (first falsy)"},
		{"nil && float", nil, false, "nil && 3.14 should return nil (first falsy)"},
		{"nil && empty", nil, false, "nil && empty string should return nil (first falsy)"},
		{"nil && str", nil, false, "nil && 'hello' should return nil (first falsy)"},
		{"nil && numStr", nil, false, "nil && '42' should return nil (first falsy)"},

		// false && others (false is falsy)
		{"false && nil", false, false, "false && nil should return false (first falsy)"},
		{"false && true", false, false, "false && true should return false (first falsy)"},
		{"false && false", false, false, "false && false should return false (first falsy)"},
		{"false && zero", false, false, "false && 0 should return false (first falsy)"},
		{"false && one", false, false, "false && 1 should return false (first falsy)"},
		{"false && empty", false, false, "false && empty string should return false (first falsy)"},
		{"false && str", false, false, "false && 'hello' should return false (first falsy)"},

		// zero && others (0 is falsy)
		{"zero && nil", 0, false, "0 && nil should return 0 (first falsy)"},
		{"zero && true", 0, false, "0 && true should return 0 (first falsy)"},
		{"zero && false", 0, false, "0 && false should return 0 (first falsy)"},
		{"zero && one", 0, false, "0 && 1 should return 0 (first falsy)"},
		{"zero && empty", 0, false, "0 && empty string should return 0 (first falsy)"},
		{"zero && str", 0, false, "0 && 'hello' should return 0 (first falsy)"},

		// empty && others (empty string is falsy)
		{"empty && nil", "", false, "empty string && nil should return empty string (first falsy)"},
		{"empty && true", "", false, "empty string && true should return empty string (first falsy)"},
		{"empty && false", "", false, "empty string && false should return empty string (first falsy)"},
		{"empty && zero", "", false, "empty string && 0 should return empty string (first falsy)"},
		{"empty && one", "", false, "empty string && 1 should return empty string (first falsy)"},
		{"empty && str", "", false, "empty string && 'hello' should return empty string (first falsy)"},

		// Truthy values && others (return second operand)
		{"true && nil", nil, false, "true && nil should return nil (second operand)"},
		{"true && true", true, false, "true && true should return true (second operand)"},
		{"true && false", false, false, "true && false should return false (second operand)"},
		{"true && zero", 0, false, "true && 0 should return 0 (second operand)"},
		{"true && one", 1, false, "true && 1 should return 1 (second operand)"},
		{"true && neg", -5, false, "true && (-5) should return -5 (second operand)"},
		{"true && float", 3.14, false, "true && 3.14 should return 3.14 (second operand)"},
		{"true && empty", "", false, "true && empty string should return empty string (second operand)"},
		{"true && str", "hello", false, "true && 'hello' should return 'hello' (second operand)"},
		{"true && numStr", "42", false, "true && '42' should return '42' (second operand)"},

		{"one && nil", nil, false, "1 && nil should return nil (second operand)"},
		{"one && true", true, false, "1 && true should return true (second operand)"},
		{"one && false", false, false, "1 && false should return false (second operand)"},
		{"one && zero", 0, false, "1 && 0 should return 0 (second operand)"},
		{"one && one", 1, false, "1 && 1 should return 1 (second operand)"},
		{"one && empty", "", false, "1 && empty string should return empty string (second operand)"},
		{"one && str", "hello", false, "1 && 'hello' should return 'hello' (second operand)"},

		{"neg && nil", nil, false, "(-5) && nil should return nil (second operand)"},
		{"neg && true", true, false, "(-5) && true should return true (second operand)"},
		{"neg && false", false, false, "(-5) && false should return false (second operand)"},
		{"neg && zero", 0, false, "(-5) && 0 should return 0 (second operand)"},
		{"neg && str", "hello", false, "(-5) && 'hello' should return 'hello' (second operand)"},

		{"float && nil", nil, false, "3.14 && nil should return nil (second operand)"},
		{"float && true", true, false, "3.14 && true should return true (second operand)"},
		{"float && zero", 0, false, "3.14 && 0 should return 0 (second operand)"},
		{"float && str", "hello", false, "3.14 && 'hello' should return 'hello' (second operand)"},

		{"str && nil", nil, false, "'hello' && nil should return nil (second operand)"},
		{"str && true", true, false, "'hello' && true should return true (second operand)"},
		{"str && false", false, false, "'hello' && false should return false (second operand)"},
		{"str && zero", 0, false, "'hello' && 0 should return 0 (second operand)"},
		{"str && one", 1, false, "'hello' && 1 should return 1 (second operand)"},
		{"str && empty", "", false, "'hello' && empty string should return empty string (second operand)"},
		{"str && str", "hello", false, "'hello' && 'hello' should return 'hello' (second operand)"},

		{"numStr && nil", nil, false, "'42' && nil should return nil (second operand)"},
		{"numStr && false", false, false, "'42' && false should return false (second operand)"},
		{"numStr && str", "hello", false, "'42' && 'hello' should return 'hello' (second operand)"},

		// ============ LOGICAL OR (||) ============
		// In JavaScript: if left is truthy, return left; otherwise return right

		// Falsy values || others (return second operand)
		{"nil || nil", nil, false, "nil || nil should return nil (second operand)"},
		{"nil || true", true, false, "nil || true should return true (second operand)"},
		{"nil || false", false, false, "nil || false should return false (second operand)"},
		{"nil || zero", 0, false, "nil || 0 should return 0 (second operand)"},
		{"nil || one", 1, false, "nil || 1 should return 1 (second operand)"},
		{"nil || neg", -5, false, "nil || (-5) should return -5 (second operand)"},
		{"nil || float", 3.14, false, "nil || 3.14 should return 3.14 (second operand)"},
		{"nil || empty", "", false, "nil || empty string should return empty string (second operand)"},
		{"nil || str", "hello", false, "nil || 'hello' should return 'hello' (second operand)"},
		{"nil || numStr", "42", false, "nil || '42' should return '42' (second operand)"},

		{"false || nil", nil, false, "false || nil should return nil (second operand)"},
		{"false || true", true, false, "false || true should return true (second operand)"},
		{"false || false", false, false, "false || false should return false (second operand)"},
		{"false || zero", 0, false, "false || 0 should return 0 (second operand)"},
		{"false || one", 1, false, "false || 1 should return 1 (second operand)"},
		{"false || empty", "", false, "false || empty string should return empty string (second operand)"},
		{"false || str", "hello", false, "false || 'hello' should return 'hello' (second operand)"},

		{"zero || nil", nil, false, "0 || nil should return nil (second operand)"},
		{"zero || true", true, false, "0 || true should return true (second operand)"},
		{"zero || false", false, false, "0 || false should return false (second operand)"},
		{"zero || one", 1, false, "0 || 1 should return 1 (second operand)"},
		{"zero || empty", "", false, "0 || empty string should return empty string (second operand)"},
		{"zero || str", "hello", false, "0 || 'hello' should return 'hello' (second operand)"},

		{"empty || nil", nil, false, "empty string || nil should return nil (second operand)"},
		{"empty || true", true, false, "empty string || true should return true (second operand)"},
		{"empty || false", false, false, "empty string || false should return false (second operand)"},
		{"empty || zero", 0, false, "empty string || 0 should return 0 (second operand)"},
		{"empty || one", 1, false, "empty string || 1 should return 1 (second operand)"},
		{"empty || str", "hello", false, "empty string || 'hello' should return 'hello' (second operand)"},

		// Truthy values || others (return first operand)
		{"true || nil", true, false, "true || nil should return true (first truthy)"},
		{"true || true", true, false, "true || true should return true (first truthy)"},
		{"true || false", true, false, "true || false should return true (first truthy)"},
		{"true || zero", true, false, "true || 0 should return true (first truthy)"},
		{"true || one", true, false, "true || 1 should return true (first truthy)"},
		{"true || empty", true, false, "true || empty string should return true (first truthy)"},
		{"true || str", true, false, "true || 'hello' should return true (first truthy)"},

		{"one || nil", 1, false, "1 || nil should return 1 (first truthy)"},
		{"one || true", 1, false, "1 || true should return 1 (first truthy)"},
		{"one || false", 1, false, "1 || false should return 1 (first truthy)"},
		{"one || zero", 1, false, "1 || 0 should return 1 (first truthy)"},
		{"one || empty", 1, false, "1 || empty string should return 1 (first truthy)"},
		{"one || str", 1, false, "1 || 'hello' should return 1 (first truthy)"},

		{"neg || nil", -5, false, "(-5) || nil should return -5 (first truthy)"},
		{"neg || true", -5, false, "(-5) || true should return -5 (first truthy)"},
		{"neg || false", -5, false, "(-5) || false should return -5 (first truthy)"},
		{"neg || zero", -5, false, "(-5) || 0 should return -5 (first truthy)"},
		{"neg || str", -5, false, "(-5) || 'hello' should return -5 (first truthy)"},

		{"float || nil", 3.14, false, "3.14 || nil should return 3.14 (first truthy)"},
		{"float || true", 3.14, false, "3.14 || true should return 3.14 (first truthy)"},
		{"float || zero", 3.14, false, "3.14 || 0 should return 3.14 (first truthy)"},
		{"float || str", 3.14, false, "3.14 || 'hello' should return 3.14 (first truthy)"},

		{"str || nil", "hello", false, "'hello' || nil should return 'hello' (first truthy)"},
		{"str || true", "hello", false, "'hello' || true should return 'hello' (first truthy)"},
		{"str || false", "hello", false, "'hello' || false should return 'hello' (first truthy)"},
		{"str || zero", "hello", false, "'hello' || 0 should return 'hello' (first truthy)"},
		{"str || empty", "hello", false, "'hello' || empty string should return 'hello' (first truthy)"},
		{"str || str", "hello", false, "'hello' || 'hello' should return 'hello' (first truthy)"},

		{"numStr || nil", "42", false, "'42' || nil should return '42' (first truthy)"},
		{"numStr || false", "42", false, "'42' || false should return '42' (first truthy)"},
		{"numStr || str", "42", false, "'42' || 'hello' should return '42' (first truthy)"},

		// Collections (arrays and objects are generally truthy unless empty)
		{"array && true", true, false, "non-empty array && true should return true (second operand)"},
		{"emptyArray && true", true, false, "empty array && true should return true (arrays are truthy even if empty)"},
		{"obj && true", true, false, "non-empty object && true should return true (second operand)"},
		{"emptyObj && true", true, false, "empty object && true should return true (objects are truthy even if empty)"},

		{"true || array", true, false, "true || array should return true (first truthy)"},
		{"false || array", []any{1, 2, 3}, false, "false || array should return array (second operand)"},
		{"nil || emptyArray", []any{}, false, "nil || empty array should return empty array (second operand)"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q (%s): got %v (type %T), want %v (type %T)", tt.expr, tt.description, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestCrossTypeSpecialOperations(t *testing.T) {
	env := map[string]any{
		"nil":    nil,
		"true":   true,
		"false":  false,
		"zero":   0,
		"one":    1,
		"str":    "hello",
		"numStr": "42",
		"array":  []any{1, nil, "hello", 0, false},
		"obj": map[string]any{
			"nil":  nil,
			"bool": true,
			"num":  42,
			"str":  "hello",
		},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ NIL COALESCING (??) ============
		// Return first non-nil value

		{"nil ?? nil", nil, false, "nil ?? nil should return nil"},
		{"nil ?? true", true, false, "nil ?? true should return true"},
		{"nil ?? false", false, false, "nil ?? false should return false"},
		{"nil ?? zero", 0, false, "nil ?? 0 should return 0"},
		{"nil ?? str", "hello", false, "nil ?? 'hello' should return 'hello'"},
		{"nil ?? numStr", "42", false, "nil ?? '42' should return '42'"},

		{"true ?? nil", true, false, "true ?? nil should return true (first non-nil)"},
		{"false ?? nil", false, false, "false ?? nil should return false (first non-nil)"},
		{"zero ?? nil", 0, false, "0 ?? nil should return 0 (first non-nil)"},
		{"str ?? nil", "hello", false, "'hello' ?? nil should return 'hello' (first non-nil)"},

		{"true ?? false", true, false, "true ?? false should return true (first non-nil)"},
		{"false ?? true", false, false, "false ?? true should return false (first non-nil)"},
		{"zero ?? one", 0, false, "0 ?? 1 should return 0 (first non-nil)"},
		{"str ?? numStr", "hello", false, "'hello' ?? '42' should return 'hello' (first non-nil)"},

		// Chained nil coalescing
		{"nil ?? nil ?? str", "hello", false, "nil ?? nil ?? 'hello' should return 'hello'"},
		{"nil ?? false ?? str", false, false, "nil ?? false ?? 'hello' should return false (first non-nil)"},

		// ============ IN OPERATOR ============
		// Check if value exists in collection

		{"nil in array", true, false, "nil should be found in array containing nil"},
		{"true in array", false, false, "true should not be found in array"},
		{"zero in array", true, false, "0 should be found in array containing 0"},
		{"false in array", true, false, "false should be found in array containing false"},
		{"str in array", true, false, "'hello' should be found in array containing 'hello'"},
		{"one in array", true, false, "1 should be found in array containing 1"},

		// Check if key exists in object (not value)
		{"'nil' in obj", true, false, "'nil' key should exist in object"},
		{"'bool' in obj", true, false, "'bool' key should exist in object"},
		{"'num' in obj", true, false, "'num' key should exist in object"},
		{"'str' in obj", true, false, "'str' key should exist in object"},
		{"'missing' in obj", false, false, "'missing' key should not exist in object"},

		// Type coercion in 'in' operator
		{"zero in [false, 0, '']", true, false, "0 should be found in array (exact match, no coercion for 'in')"},
		{"false in [0, false, '']", true, false, "false should be found in array (exact match)"},

		// ============ TERNARY OPERATOR (condition ? true : false) ============
		// Test truthiness of different types

		{"nil ? 'truthy' : 'falsy'", "falsy", false, "nil should be falsy"},
		{"true ? 'truthy' : 'falsy'", "truthy", false, "true should be truthy"},
		{"false ? 'truthy' : 'falsy'", "falsy", false, "false should be falsy"},
		{"zero ? 'truthy' : 'falsy'", "falsy", false, "0 should be falsy"},
		{"one ? 'truthy' : 'falsy'", "truthy", false, "1 should be truthy"},
		{"str ? 'truthy' : 'falsy'", "truthy", false, "non-empty string should be truthy"},
		{"'' ? 'truthy' : 'falsy'", "falsy", false, "empty string should be falsy"},
		{"numStr ? 'truthy' : 'falsy'", "truthy", false, "numeric string should be truthy"},
		{"array ? 'truthy' : 'falsy'", "truthy", false, "array should be truthy"},
		{"obj ? 'truthy' : 'falsy'", "truthy", false, "object should be truthy"},

		// Complex ternary expressions
		{"(nil || false) ? 'yes' : 'no'", "no", false, "(nil || false) should be falsy"},
		{"(nil || true) ? 'yes' : 'no'", "yes", false, "(nil || true) should be truthy"},
		{"(zero && one) ? 'yes' : 'no'", "no", false, "(0 && 1) should be falsy (returns 0)"},
		{"(one && str) ? 'yes' : 'no'", "yes", false, "(1 && 'hello') should be truthy (returns 'hello')"},

		// Nested ternary
		{"true ? (false ? 'inner-true' : 'inner-false') : 'outer-false'", "inner-false", false, "nested ternary"},
		{"false ? 'outer-true' : (true ? 'inner-true' : 'inner-false')", "inner-true", false, "nested ternary"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q (%s): got %v (type %T), want %v (type %T)", tt.expr, tt.description, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

// Helper functions for better comparisons

func approximatelyEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle float comparisons with tolerance
	aFloat, aIsFloat := toFloat64Safe(a)
	bFloat, bIsFloat := toFloat64Safe(b)
	if aIsFloat && bIsFloat {
		diff := aFloat - bFloat
		if diff < 0 {
			diff = -diff
		}
		return diff < 0.001 // Tolerance for floating point comparisons
	}

	// Fall back to direct comparison
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func toFloat64Safe(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

func deepEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle slices
	aSlice, aIsSlice := a.([]any)
	bSlice, bIsSlice := b.([]any)
	if aIsSlice && bIsSlice {
		if len(aSlice) != len(bSlice) {
			return false
		}
		for i := range aSlice {
			if !deepEqual(aSlice[i], bSlice[i]) {
				return false
			}
		}
		return true
	}

	// Handle maps
	aMap, aIsMap := a.(map[string]any)
	bMap, bIsMap := b.(map[string]any)
	if aIsMap && bIsMap {
		if len(aMap) != len(bMap) {
			return false
		}
		for k, v := range aMap {
			if bVal, exists := bMap[k]; !exists || !deepEqual(v, bVal) {
				return false
			}
		}
		return true
	}

	// For other types, use approximatelyEqual
	return approximatelyEqual(a, b)
}
