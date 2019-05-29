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
)

func TestRun_debug(t *testing.T) {
	var test = struct {
		input  string
		output interface{}
	}{
		`filter([1,2,3], {# > 2})`,
		[]interface{}{int64(3)},
	}

	env := &mockEnv{}

	node, err := parser.Parse(test.input)
	require.NoError(t, err, test.input)

	_, err = checker.Check(node, checker.Env(&mockEnv{}))
	require.NoError(t, err, test.input)

	program, err := compiler.Compile(node)
	require.NoError(t, err, test.input)

	output, err := vm.Run(program, env)
	require.NoError(t, err, test.input)

	assert.Equal(t, test.output, output, test.input)
}

func TestRun(t *testing.T) {
	type test struct {
		input  string
		output interface{}
	}
	var tests = []test{
		{
			`1`,
			int64(1),
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
			`Int64 == 0 && Float64 == 0 && Bool && String == "string"`,
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
			int64(2),
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
			int64(5),
		},
		{
			`Ticket.Price`,
			int(100),
		},
		{
			`Add(10, 5) + GetInt()`,
			int(15),
		},
		{
			`Ticket.String()`,
			`$100`,
		},
		{
			`[1, 2, 3]`,
			[]interface{}{int64(1), int64(2), int64(3)},
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]interface{}{"foo": int64(0), "bar": int64(1)},
		},
		{
			`[1, 2, 3]`,
			[]interface{}{int64(1), int64(2), int64(3)},
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]interface{}{"foo": int64(0), "bar": int64(1)},
		},
		{
			`1 in [1, 2, 3] && "foo" in {foo: 0, bar: 1} && "Price" in Ticket`,
			true,
		},
		{
			`(true ? 0+1 : 2+3) + (false ? -1 : -2)`,
			int64(-1),
		},
		{
			`len(Array)`,
			int64(5),
		},
		{
			`filter(1..9, {# > 7})`,
			[]interface{}{int64(8), int64(9)},
		},
		{
			`map(1..3, {# * #})`,
			[]interface{}{int64(1), int64(4), int64(9)},
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
	Any     interface{}
	Int     int
	Int32   int32
	Int64   int64
	Uint64  uint64
	Float64 float64
	Bool    bool
	String  string
	Array   []int
	Ticket  *mockTicket
}

func (e *mockEnv) GetInt() int {
	return e.Int
}

func (*mockEnv) Add(a, b int64) int {
	return int(a + b)
}

type mockTicket struct {
	Price int
}

func (t *mockTicket) String() string {
	return fmt.Sprintf("$%v", t.Price)
}
