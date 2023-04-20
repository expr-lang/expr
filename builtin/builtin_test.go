package builtin_test

import (
	"testing"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	input string
	want  interface{}
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
	{`upper("hello")`, "HELLO"},
	{`lower("HELLO")`, "hello"},
	{`left("foobar", 3)`, "foo"},
	{`left("foobar", -4)`, "fo"},
	{`right("foobar", 3)`, "bar"},
	{`right("foobar", -4)`, "ar"},
}

func TestBuiltin(t *testing.T) {
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			out, err := expr.Eval(test.input, nil)
			assert.NoError(t, err)
			assert.Equal(t, test.want, out)
		})
	}
}

var errorTests = []struct {
	input string
	err   string
}{
	{`len()`, `invalid number of arguments for len (expected 1, got 0)`},
	{`len(1)`, `invalid argument for len (type int)`},
	{`abs()`, `invalid number of arguments for abs (expected 1, got 0)`},
	{`abs(1, 2)`, `invalid number of arguments for abs (expected 1, got 2)`},
	{`abs("foo")`, `invalid argument for abs (type string)`},
	{`int()`, `invalid number of arguments for int (expected 1, got 0)`},
	{`int(1, 2)`, `invalid number of arguments for int (expected 1, got 2)`},
	{`float()`, `invalid number of arguments for float (expected 1, got 0)`},
	{`float(1, 2)`, `invalid number of arguments for float (expected 1, got 2)`},
	{`upper()`, `invalid number of arguments for upper (expected 1, got 0)`},
	{`upper("hello", "world")`, `invalid number of arguments for upper (expected 1, got 2)`},
	{`lower()`, `invalid number of arguments for lower (expected 1, got 0)`},
	{`lower("hello", "world")`, `invalid number of arguments for lower (expected 1, got 2)`},
	{`left()`, `invalid number of arguments for left (expected 2, got 0)`},
	{`left("foo", "bar")`, `invalid argument no. 2 for left (type string)`},
	{`left("foobar", 3, 4)`, `invalid number of arguments for left (expected 2, got 3)`},
	{`right()`, `invalid number of arguments for right (expected 2, got 0)`},
	{`right("foo", "bar")`, `invalid argument no. 2 for right (type string)`},
	{`right("foobar", 3, 4)`, `invalid number of arguments for right (expected 2, got 3)`},
}

func TestBuiltinErrors(t *testing.T) {
	for _, test := range errorTests {
		t.Run(test.input, func(t *testing.T) {
			_, err := expr.Eval(test.input, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.err)
		})
	}
}
