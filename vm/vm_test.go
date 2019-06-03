package vm_test

import (
	"fmt"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRun_debug(t *testing.T) {
	env := &mockEnv{}

	var input = `Int64 == 0`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	_, err = checker.Check(node, checker.Env(&mockEnv{}))
	require.NoError(t, err)

	program, err := compiler.Compile(node)
	require.NoError(t, err)

	_, err = vm.Run(program, env)
	require.NoError(t, err)
}

func TestRun(t *testing.T) {
	type test struct {
		input  string
		output interface{}
	}
	var tests = []test{
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
			uint64(0),
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
			`2 ** 4`,
			float64(16),
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
			`[1, 2, 3]`,
			[]interface{}{1, 2, 3},
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]interface{}{"foo": 0, "bar": 1},
		},
		{
			`[1, 2, 3]`,
			[]interface{}{1, 2, 3},
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]interface{}{"foo": 0, "bar": 1},
		},
		{
			`1 in [1, 2, 3] && "foo" in {foo: 0, bar: 1} && "Price" in Ticket`,
			true,
		},
		{
			`1.5 in [1] && 1 in [1.5]`,
			false,
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
	}

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
		Ticket: &mockTicket{
			Price: 100,
		},
		BirthDay: time.Date(2017, time.October, 23, 18, 30, 0, 0, time.UTC),
		Now:      time.Now(),
	}

	for _, test := range tests {
		tree, err := parser.Parse(test.input)
		require.NoError(t, err, test.input)

		_, err = checker.Check(tree, checker.Env(&mockEnv{}))
		require.NoError(t, err, test.input)

		program, err := compiler.Compile(tree)
		require.NoError(t, err, test.input)

		output, err := vm.Run(program, env)
		require.NoError(t, err, test.input)

		assert.Equal(t, test.output, output, test.input)
	}
}

type mockEnv struct {
	Any      interface{}
	Int      int
	Int32    int32
	Int64    int64
	Uint64   uint64
	Float64  float64
	Bool     bool
	String   string
	Array    []int
	Ticket   *mockTicket
	BirthDay time.Time
	Now      time.Time
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

type mockTicket struct {
	Price int
}

func (t *mockTicket) PriceDiv(p int) int {
	return t.Price / p
}

func (t *mockTicket) String() string {
	return fmt.Sprintf("$%v", t.Price)
}
