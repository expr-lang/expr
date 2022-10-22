package compiler_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/antonmedv/expr/vm/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type B struct {
	_ byte
	_ byte
	C struct {
		_ byte
		_ byte
		_ byte
		D int
	}
}

type Env struct {
	A struct {
		_   byte
		B   B
		Map map[string]B
	}
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
					vm.OpEqualString,
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
					vm.OpEqualInt,
				},
			},
		},
		{
			`-1`,
			vm.Program{
				Constants: []interface{}{-1},
				Bytecode: []byte{
					vm.OpPush, 0, 0,
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
		{
			`A.B.C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0, 1, 2, 3},
						Path:  "A.B.C.D",
					},
				},
				Bytecode: []byte{
					vm.OpFetchEnvField, 0, 0,
				},
			},
		},
		{
			`A?.B.C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0},
						Path:  "A",
					},
					&runtime.Field{
						Index: []int{1, 2, 3},
						Path:  "B.C.D",
					},
				},
				Bytecode: []byte{
					vm.OpFetchEnvField, 0, 0,
					vm.OpJumpIfNil, 3, 0,
					vm.OpFetchField, 1, 0,
				},
			},
		},
		{
			`A.B?.C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0, 1},
						Path:  "A.B",
					},
					&runtime.Field{
						Index: []int{2, 3},
						Path:  "C.D",
					},
				},
				Bytecode: []byte{
					vm.OpFetchEnvField, 0, 0,
					vm.OpJumpIfNil, 3, 0,
					vm.OpFetchField, 1, 0,
				},
			},
		},
		{
			`A.Map["B"].C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0, 2},
						Path:  "A.Map",
					},
					"B",
					&runtime.Field{
						Index: []int{2, 3},
						Path:  "C.D",
					},
				},
				Bytecode: []byte{
					vm.OpFetchEnvField, 0, 0,
					vm.OpPush, 1, 0,
					vm.OpFetch,
					vm.OpFetchField, 2, 0,
				},
			},
		},
	}

	for _, test := range tests {
		program, err := expr.Compile(test.input, expr.Env(Env{}))
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
