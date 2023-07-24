package compiler_test

import (
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeCommonExpr(t *testing.T) {
	type testCase struct {
		input   string
		expr    string
		program *vm.Program
	}
	for _, tt := range []testCase{
		{
			input: `(1+1) == 2 and 2 == (1 +   1)`,
			expr:  `((1 + 1) == 2) and ((1 + 1) == 2)`,
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpAdd,
					vm.OpSaveCommon,
					vm.OpPush,
					vm.OpEqualInt,
					vm.OpSaveCommon,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpAdd,
					vm.OpSaveCommon,
					vm.OpEqualInt,
					vm.OpSaveCommon,
				},
				Arguments: []int{
					0, 11, 0, 1, 5, 0, 0, 0, 0, 1, 1, 0, 0,
					14, 0,
					0, 11, 0, 1, 1, 5, 0, 0, 0, 0, 1, 0, 0,
				},
			},
		},
		{
			input: `not ((+(1+1)) ** 2 == 4 or 4 == (-(1 +   1)) ^ 2)`,
			expr:  `not (((-(1 + 1) ** 2) == 4) or ((1 + 1) ** 2) == 4))`,
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpAdd,
					vm.OpSaveCommon,
					vm.OpPush,
					vm.OpExponent,
					vm.OpPush,
					vm.OpEqual,
					vm.OpJumpIfTrue,
					vm.OpPop,
					vm.OpPush,
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpAdd,
					vm.OpSaveCommon,
					vm.OpNegate,
					vm.OpPush,
					vm.OpExponent,
					vm.OpEqual,
					vm.OpNot,
				},
				Arguments: []int{
					0, 5, 0, 0, 0, 0, 0, 1, 0, 2, 0, 13, 0, 2,
					0, 5, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0,
				},
			},
		},
		{
			input: `nil ==    nil`,
			expr:  `nil == nil`,
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpNil,
					vm.OpNil,
					vm.OpEqual,
				},
				Arguments: []int{0, 0, 0},
			},
		},
		{
			input: `true && (true || false)`,
			expr:  `(false or true) and true`,
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpTrue,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpTrue,
					vm.OpJumpIfTrue,
					vm.OpPop,
					vm.OpFalse,
				},
				Arguments: []int{0, 5, 0, 0, 2, 0, 0},
			},
		},
		{
			input: `1 >= 2 and 2 < 1`,
			expr:  `(2 < 1) and (2 < 1)`,
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpMoreOrEqual,
					vm.OpSaveCommon,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpLess,
					vm.OpSaveCommon,
				},
				Arguments: []int{0, 5, 0, 0, 1, 0, 0, 8, 0, 0, 5, 0, 1, 0, 0, 0},
			},
		},
		{
			input: `1 <= 2 and 2 > 1`,
			expr:  `(2 > 1) and (2 > 1)`,
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpLessOrEqual,
					vm.OpSaveCommon,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpMore,
					vm.OpSaveCommon,
				},
				Arguments: []int{0, 5, 0, 0, 1, 0, 0, 8, 0, 0, 5, 0, 1, 0, 0, 0},
			},
		},
		{
			input: `A.B.C.D`,
			expr:  "A.B.C.D",
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
				},
				Arguments: []int{0},
			},
		},
		{
			input: `A?.B.C.D`,
			expr:  "A?.B.C.D",
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNil,
					vm.OpFetchField,
				},
				Arguments: []int{0, 1, 1},
			},
		},
		{
			input: `A.B?.C.D`,
			expr:  "A.B?.C.D",
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNil,
					vm.OpFetchField,
				},
				Arguments: []int{0, 1, 1},
			},
		},
		{
			input: `A.Map["B"].C.D`,
			expr:  "A.Map[\"B\"].C.D",
			program: &vm.Program{
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
			input: `A ?? 1`,
			expr:  "A ?? 1",
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNotNil,
					vm.OpPop,
					vm.OpPush,
				},
				Arguments: []int{0, 2, 0, 1},
			},
		},
		{
			input: `int("123") == int("123")`,
			expr:  `int("123") == int("123")`,
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpPush,
					vm.OpBuiltin,
					vm.OpPush,
					vm.OpBuiltin,
					vm.OpEqualInt,
				},
				Arguments: []int{0, 3, 0, 3, 0},
			},
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			program, err := expr.Compile(tt.input, expr.Env(Env{}), expr.Optimize(false), expr.AllowReuseCommon(true))
			if err != nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expr, program.Node.SubExpr())
				assert.Equal(t, tt.program.Bytecode, program.Bytecode)
				assert.Equal(t, tt.program.Arguments, program.Arguments)
			}

		})
	}

}
