package compiler_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompile_debug(t *testing.T) {
	input := `false && true && true`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	_, err = compiler.Compile(tree, nil)
	require.NoError(t, err)
}

func TestCompile(t *testing.T) {
	type test struct {
		input   string
		program vm.Program
	}
	var tests = []test{
		{
			`65535`,
			vm.Program{
				Constants: []interface{}{
					int(math.MaxUint16),
				},
				Bytecode: []byte{
					vm.OpPush, 0, 0,
				},
			},
		},
		{
			`.5`,
			vm.Program{
				Constants: []interface{}{
					float64(.5),
				},
				Bytecode: []byte{
					vm.OpPush, 0, 0,
				},
			},
		},
		{
			`true`,
			vm.Program{
				Bytecode: []byte{
					vm.OpTrue,
				},
			},
		},
		{
			`Name`,
			vm.Program{
				Constants: []interface{}{
					"Name",
				},
				Bytecode: []byte{
					vm.OpFetch, 0, 0,
				},
			},
		},
		{
			`"string"`,
			vm.Program{
				Constants: []interface{}{
					"string",
				},
				Bytecode: []byte{
					vm.OpPush, 0, 0,
				},
			},
		},
		{
			`"string" == "string"`,
			vm.Program{
				Constants: []interface{}{
					"string",
				},
				Bytecode: []byte{
					vm.OpPush, 0, 0,
					vm.OpPush, 0, 0,
					vm.OpEqual,
				},
			},
		},
		{
			`1000000 == 1000000`,
			vm.Program{
				Constants: []interface{}{
					int64(1000000),
				},
				Bytecode: []byte{
					vm.OpPush, 0, 0,
					vm.OpPush, 0, 0,
					vm.OpEqual,
				},
			},
		},
		{
			`-1`,
			vm.Program{
				Constants: []interface{}{1},
				Bytecode: []byte{
					vm.OpPush, 0, 0,
					vm.OpNegate,
				},
			},
		},
		{
			`true && true || true`,
			vm.Program{
				Bytecode: []byte{
					vm.OpTrue,
					vm.OpJumpIfFalse, 2, 0,
					vm.OpPop,
					vm.OpTrue,
					vm.OpJumpIfTrue, 2, 0,
					vm.OpPop,
					vm.OpTrue,
				},
			},
		},
	}

	for _, test := range tests {
		node, err := parser.Parse(test.input)
		require.NoError(t, err)

		program, err := compiler.Compile(node, nil)
		require.NoError(t, err, test.input)

		assert.Equal(t, test.program.Disassemble(), program.Disassemble(), test.input)
	}
}

func TestCompile_cast(t *testing.T) {
	input := `1`
	expected := &vm.Program{
		Constants: []interface{}{
			1,
		},
		Bytecode: []byte{
			vm.OpPush, 0, 0,
			vm.OpCast, 1, 0,
		},
	}

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, &conf.Config{Expect: reflect.Float64})
	require.NoError(t, err)

	assert.Equal(t, expected.Disassemble(), program.Disassemble())
}
