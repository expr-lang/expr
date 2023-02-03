package compiler_test

import (
	"math"
	"testing"

	"github.com/antonmedv/expr"
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
					math.MaxUint16,
				},
				Bytecode: []vm.Opcode{
					vm.OpPush,
				},
				Arguments: []int{0},
			},
		},
		{
			`.5`,
			vm.Program{
				Constants: []interface{}{
					.5,
				},
				Bytecode: []vm.Opcode{
					vm.OpPush,
				},
				Arguments: []int{0},
			},
		},
		{
			`true`,
			vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpTrue,
				},
				Arguments: []int{0},
			},
		},
		{
			`"string"`,
			vm.Program{
				Constants: []interface{}{
					"string",
				},
				Bytecode: []vm.Opcode{
					vm.OpPush,
				},
				Arguments: []int{0},
			},
		},
		{
			`"string" == "string"`,
			vm.Program{
				Constants: []interface{}{
					"string",
				},
				Bytecode: []vm.Opcode{
					vm.OpPush,
					vm.OpPush,
					vm.OpEqualString,
				},
				Arguments: []int{0, 0, 0},
			},
		},
		{
			`1000000 == 1000000`,
			vm.Program{
				Constants: []interface{}{
					int64(1000000),
				},
				Bytecode: []vm.Opcode{
					vm.OpPush,
					vm.OpPush,
					vm.OpEqualInt,
				},
				Arguments: []int{0, 0, 0},
			},
		},
		{
			`-1`,
			vm.Program{
				Constants: []interface{}{1},
				Bytecode: []vm.Opcode{
					vm.OpPush,
					vm.OpNegate,
				},
				Arguments: []int{0, 0},
			},
		},
		{
			`true && true || true`,
			vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpTrue,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpTrue,
					vm.OpJumpIfTrue,
					vm.OpPop,
					vm.OpTrue,
				},
				Arguments: []int{0, 2, 0, 0, 2, 0, 0},
			},
		},
		{
			`true && (true || true)`,
			vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpTrue,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpTrue,
					vm.OpJumpIfTrue,
					vm.OpPop,
					vm.OpTrue,
				},
				Arguments: []int{0, 5, 0, 0, 2, 0, 0},
			},
		},
		{
			`A.B.C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0, 1, 2, 3},
						Path:  []string{"A", "B", "C", "D"},
					},
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
				},
				Arguments: []int{0},
			},
		},
		{
			`A?.B.C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0},
						Path:  []string{"A"},
					},
					&runtime.Field{
						Index: []int{1, 2, 3},
						Path:  []string{"B", "C", "D"},
					},
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNil,
					vm.OpFetchField,
				},
				Arguments: []int{0, 1, 1},
			},
		},
		{
			`A.B?.C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0, 1},
						Path:  []string{"A", "B"},
					},
					&runtime.Field{
						Index: []int{2, 3},
						Path:  []string{"C", "D"},
					},
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNil,
					vm.OpFetchField,
				},
				Arguments: []int{0, 1, 1},
			},
		},
		{
			`A.Map["B"].C.D`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0, 2},
						Path:  []string{"A", "Map"},
					},
					"B",
					&runtime.Field{
						Index: []int{2, 3},
						Path:  []string{"C", "D"},
					},
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpPush,
					vm.OpFetch,
					vm.OpFetchField,
				},
				Arguments: []int{0, 1, 0, 2},
			},
		},
		{
			`A ?? 1`,
			vm.Program{
				Constants: []interface{}{
					&runtime.Field{
						Index: []int{0},
						Path:  []string{"A"},
					},
					1,
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNotNil,
					vm.OpPop,
					vm.OpPush,
				},
				Arguments: []int{0, 2, 0, 1},
			},
		},
	}

	for _, test := range tests {
		program, err := expr.Compile(test.input, expr.Env(Env{}), expr.Optimize(false))
		require.NoError(t, err, test.input)

		assert.Equal(t, test.program.Disassemble(), program.Disassemble(), test.input)
	}
}
