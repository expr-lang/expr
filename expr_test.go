package expr_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleEval() {
	output, err := expr.Eval("'hello world'", nil)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: hello world
}

func ExampleEval_map() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": []string{"zero", "hello world"},
		"swipe": func(in string) string {
			return strings.Replace(in, "world", "user", 1)
		},
	}

	output, err := expr.Eval("swipe(bar[foo])", env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: hello user
}

type mockMapEnv map[string]interface{}

func (mockMapEnv) Swipe(in string) string {
	return strings.Replace(in, "world", "user", 1)
}

func ExampleEval_map_method() {
	env := mockMapEnv{
		"foo": 1,
		"bar": []string{"zero", "hello world"},
	}

	program, err := expr.Compile("Swipe(bar[foo])", expr.Env(env))
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

	// Output: hello user
}

func ExampleEval_struct() {
	type C struct{ C int }
	type B struct{ B *C }
	type A struct{ A B }

	env := A{B{&C{42}}}

	output, err := expr.Eval("A.B.C", env)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: 42
}

func ExampleEval_error() {
	output, err := expr.Eval("(boo + bar]", nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: syntax error: mismatched input ']' expecting ')' (1:11)
	//  | (boo + bar]
	//  | ..........^
}

func ExampleEval_matches() {
	output, err := expr.Eval(`"a" matches "a("`, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: error parsing regexp: missing closing ): `a(` (1:13)
	//  | "a" matches "a("
	//  | ............^
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
	type Request struct {
		Segments   []*Segment
		Passengers *Passengers
		Marker     string
		Meta       map[string]interface{}
	}

	code := `Segments[0].Origin == "MOW" && Passengers.Adults == 2 && Marker == "test" && Meta["accept"]`

	program, err := expr.Compile(code, expr.Env(&Request{}))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	request := &Request{
		Segments: []*Segment{
			{Origin: "MOW"},
		},
		Passengers: &Passengers{
			Adults: 2,
		},
		Marker: "test",
		Meta:   map[string]interface{}{"accept": true},
	}

	output, err := expr.Run(program, request)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEnv_with_undefined_variables() {
	env := map[string]interface{}{
		"foo": 0,
		"bar": 0,
	}

	program, err := expr.Compile(`foo + (bar != nil ? bar : 2)`, expr.Env(env))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	request := map[string]interface{}{
		"foo": 3,
	}

	output, err := expr.Run(program, request)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: 5
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

func ExampleAsInt64() {
	env := map[string]float64{
		"foo": 3,
	}

	program, err := expr.Compile("foo + 2", expr.Env(env), expr.AsInt64())
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

func ExampleOperator() {
	type Place struct {
		Code string
	}
	type Segment struct {
		Origin Place
	}
	type Helpers struct {
		PlaceEq func(p Place, s string) bool
	}
	type Request struct {
		Segments []*Segment
		Helpers
	}

	code := `Segments[0].Origin == "MOW" && PlaceEq(Segments[0].Origin, "MOW")`

	program, err := expr.Compile(code, expr.Env(&Request{}), expr.Operator("==", "PlaceEq"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	request := &Request{
		Segments: []*Segment{
			{Origin: Place{Code: "MOW"}},
		},
		Helpers: Helpers{PlaceEq: func(p Place, s string) bool {
			return p.Code == s
		}},
	}

	output, err := expr.Run(program, request)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleOperator_time() {
	type Segment struct {
		Date time.Time
	}
	type Request struct {
		Segments []Segment
		Before   func(a, b time.Time) bool
		Date     func(s string) time.Time
	}

	code := `Date("2001-01-01") < Segments[0].Date`

	program, err := expr.Compile(code, expr.Env(&Request{}), expr.Operator("<", "Before"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	request := &Request{
		Segments: []Segment{
			{Date: time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC)},
		},
		Before: func(a, b time.Time) bool {
			return a.Before(b)
		},
		Date: func(s string) time.Time {
			date, err := time.Parse("2006-01-02", s)
			if err != nil {
				panic(err)
			}
			return date
		},
	}

	output, err := expr.Run(program, request)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEval_marshal() {
	env := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	program, err := expr.Compile("(foo + bar) in [1, 2, 3]", expr.Env(env))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	b, err := json.Marshal(program)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	unmarshaledProgram := &vm.Program{}
	err = json.Unmarshal(b, unmarshaledProgram)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(unmarshaledProgram, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
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

func TestExpr(t *testing.T) {
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
		BirthDay:      time.Date(2017, time.October, 23, 18, 30, 0, 0, time.UTC),
		Now:           time.Now(),
		One:           1,
		Two:           2,
		Three:         3,
		MultiDimArray: [][]int{{1, 2, 3}, {1, 2, 3}},
		Sum: func(list []int) int {
			var ret int
			for _, el := range list {
				ret += el
			}
			return ret
		},
		Inc: func(a int) int { return a + 1 },
		Nil: nil,
	}

	tests := []struct {
		code string
		want interface{}
	}{
		{
			`1`,
			int(1),
		},
		{
			`-.5`,
			float64(-.5),
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
			int(0),
		},
		{
			`Uint64 + Int64`,
			int64(0),
		},
		{
			`Int32 + Int64`,
			int64(0),
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
			`len(Array)`,
			5,
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
	}

	for _, tt := range tests {
		program, err := expr.Compile(tt.code, expr.Env(&mockEnv{}))
		require.NoError(t, err, tt.code)

		got, err := expr.Run(program, env)
		require.NoError(t, err, tt.code)

		assert.Equal(t, tt.want, got, tt.code)
	}

	for _, tt := range tests {
		program, err := expr.Compile(tt.code, expr.Optimize(false))
		require.NoError(t, err, tt.code)

		got, err := expr.Run(program, env)
		require.NoError(t, err, tt.code)

		assert.Equal(t, tt.want, got, "unoptimized: "+tt.code)
	}

	for _, tt := range tests {
		got, err := expr.Eval(tt.code, env)
		require.NoError(t, err, tt.code)

		assert.Equal(t, tt.want, got, "eval: "+tt.code)
	}
}

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
	Nil                  interface{}
	NilStruct            *time.Time
	NilInt               *int
	NilSlice             []ticket
}

func (e *mockEnv) GetInt() int {
	return e.Int
}

func (*mockEnv) Add(a, b int) int {
	return int(a + b)
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

func TestExpr_eval_with_env(t *testing.T) {
	_, err := expr.Eval("true", expr.Env(map[string]interface{}{}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "misused")
}
