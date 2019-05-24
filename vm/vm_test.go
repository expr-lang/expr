package vm_test

import (
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/internal/helper"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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
	}

	env := &mockEnv{
		Int64:   0,
		Uint64:  0,
		Float64: 0,
		Bool:    true,
		String:  "string",
	}

	for _, test := range tests {
		source := helper.NewSource(test.input)

		node, err := parser.ParseSource(source)
		require.NoError(t, err)

		_, err = checker.Check(node, source, checker.Env(&mockEnv{}))
		require.NoError(t, err)

		program, err := compiler.Compile(node)
		require.NoError(t, err)

		output, err := vm.Run(program, env)
		require.NoError(t, err)

		assert.Equal(t, test.output, output, test.input)
	}
}

type mockEnv struct {
	Int64   int64
	Uint64  uint64
	Float64 float64
	Bool    bool
	String  string
}
