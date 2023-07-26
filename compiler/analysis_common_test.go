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
		{
			input: `true ? 1 : 2`,
			expr:  "true ? 1 : 2",
			program: &vm.Program{
				Bytecode: []vm.Opcode{
					vm.OpTrue,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpPush,
					vm.OpJump,
					vm.OpPop,
					vm.OpPush,
				},
				Arguments: []int{0, 3, 0, 0, 2, 0, 1},
			},
		},
		{
			input: `[true, true, false, 1.7]`,
			expr:  `[true,true,false,1.7000000000]`,
			program: &vm.Program{
				Constants: []interface{}{1.7, 4},
				Bytecode: []vm.Opcode{
					vm.OpTrue,
					vm.OpTrue,
					vm.OpFalse,
					vm.OpPush,
					vm.OpPush,
					vm.OpArray,
				},
				Arguments: []int{0, 0, 0, 0, 1, 0},
			},
		},
		{
			input: `{1: true, 2: true, 3: false, 4: 1.7}`,
			expr:  `{"1":true,"2":true,"3":false,"4":1.7000000000}`,
			program: &vm.Program{
				Constants: []interface{}{1.7, 4},
				Bytecode: []vm.Opcode{
					vm.OpPush,
					vm.OpTrue,
					vm.OpPush,
					vm.OpTrue,
					vm.OpPush,
					vm.OpFalse,
					vm.OpPush,
					vm.OpPush,
					vm.OpPush,
					vm.OpMap,
				},
				Arguments: []int{0, 0, 1, 0, 2, 0, 3, 4, 5, 0},
			},
		},
		{
			input: `all([true, false], { # })`,
			expr:  `all([true,false],#)`,
			program: &vm.Program{
				Constants: []interface{}{1.7, 4},
				Bytecode: []vm.Opcode{
					vm.OpTrue,
					vm.OpFalse,
					vm.OpPush,
					vm.OpArray,
					vm.OpBegin,
					vm.OpJumpIfEnd,
					vm.OpPointer,
					vm.OpJumpIfFalse,
					vm.OpPop,
					vm.OpIncrementIt,
					vm.OpJumpBackward,
					vm.OpTrue,
					vm.OpEnd,
				},
				Arguments: []int{0, 0, 0, 0, 0, 5, 0, 4, 0, 0, 6, 0, 0},
			},
		},
		{
			input: `any([1+1, 1+1], { # > 1 })`,
			expr:  `any([1 + 1,1 + 1],# > 1)`,
			program: &vm.Program{
				Constants: []interface{}{1.7, 4},
				Bytecode: []vm.Opcode{
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpAdd,
					vm.OpSaveCommon,
					vm.OpLoadCommon,
					vm.OpJumpIfSaveCommon,
					vm.OpPop,
					vm.OpPush,
					vm.OpPush,
					vm.OpAdd,
					vm.OpSaveCommon,
					vm.OpPush,
					vm.OpArray,
					vm.OpBegin,
					vm.OpJumpIfEnd,
					vm.OpPointer,
					vm.OpPush,
					vm.OpMore,
					vm.OpJumpIfTrue,
					vm.OpPop,
					vm.OpIncrementIt,
					vm.OpJumpBackward,
					vm.OpFalse,
					vm.OpEnd,
				},
				Arguments: []int{0, 5, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 1, 0, 0, 7, 0, 0, 0, 4, 0, 0, 8, 0, 0},
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
