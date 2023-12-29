package compiler_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/test/mock"
	"github.com/expr-lang/expr/test/playground"
	"github.com/expr-lang/expr/vm"
	"github.com/expr-lang/expr/vm/runtime"
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

func (B) FuncInB() int {
	return 0
}

type Env struct {
	A struct {
		_   byte
		B   B
		Map map[string]B
		Ptr *int
	}
}

// AFunc is a method what goes before Func in the alphabet.
func (e Env) AFunc() int {
	return 0
}

func (e Env) Func() B {
	return B{}
}

func TestCompile(t *testing.T) {
	var tests = []struct {
		code string
		want vm.Program
	}{
		{
			`65535`,
			vm.Program{
				Constants: []any{
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
				Constants: []any{
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
				Constants: []any{
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
				Constants: []any{
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
				Constants: []any{
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
				Constants: []any{1},
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
				Constants: []any{
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
				Constants: []any{
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
				Constants: []any{
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
				Constants: []any{
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
				Constants: []any{
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
		{
			`A.Ptr + 1`,
			vm.Program{
				Constants: []any{
					&runtime.Field{
						Index: []int{0, 3},
						Path:  []string{"A", "Ptr"},
					},
					1,
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpDeref,
					vm.OpPush,
					vm.OpAdd,
				},
				Arguments: []int{0, 0, 1, 0},
			},
		},
		{
			`Func()`,
			vm.Program{
				Constants: []any{
					&runtime.Method{
						Index: 1,
						Name:  "Func",
					},
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadMethod,
					vm.OpCall,
				},
				Arguments: []int{0, 0},
			},
		},
		{
			`Func().FuncInB()`,
			vm.Program{
				Constants: []any{
					&runtime.Method{
						Index: 1,
						Name:  "Func",
					},
					&runtime.Method{
						Index: 0,
						Name:  "FuncInB",
					},
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadMethod,
					vm.OpCall,
					vm.OpMethod,
					vm.OpCallTyped,
				},
				Arguments: []int{0, 0, 1, 10},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			program, err := expr.Compile(test.code, expr.Env(Env{}), expr.Optimize(false))
			require.NoError(t, err)

			assert.Equal(t, test.want.Disassemble(), program.Disassemble())
		})
	}
}

func TestCompile_panic(t *testing.T) {
	tests := []string{
		`(TotalPosts.Profile[Authors > TotalPosts == get(nil, TotalLikes)] > Authors) ^ (TotalLikes / (Posts?.PublishDate[TotalPosts] < Posts))`,
		`one(Posts, nil)`,
		`trim(TotalViews, Posts) <= get(Authors, nil)`,
		`Authors.IsZero(nil * Authors) - (TotalViews && Posts ? nil : nil)[TotalViews.IsZero(false, " ").IsZero(Authors)]`,
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			_, err := expr.Compile(test, expr.Env(playground.Blog{}))
			require.Error(t, err)
		})
	}
}

func TestCompile_FuncTypes(t *testing.T) {
	env := map[string]any{
		"fn": func([]any, string) string {
			return "foo"
		},
	}
	program, err := expr.Compile("fn([1, 2], 'bar')", expr.Env(env))
	require.NoError(t, err)
	require.Equal(t, vm.OpCallTyped, program.Bytecode[3])
	require.Equal(t, 22, program.Arguments[3])
}

func TestCompile_FuncTypes_with_Method(t *testing.T) {
	env := mock.Env{}
	program, err := expr.Compile("FuncTyped('bar')", expr.Env(env))
	require.NoError(t, err)
	require.Equal(t, vm.OpCallTyped, program.Bytecode[2])
	require.Equal(t, 42, program.Arguments[2])
}

func TestCompile_FuncTypes_excludes_named_functions(t *testing.T) {
	env := mock.Env{}
	program, err := expr.Compile("FuncNamed('bar')", expr.Env(env))
	require.NoError(t, err)
	require.Equal(t, vm.OpCall, program.Bytecode[2])
	require.Equal(t, 1, program.Arguments[2])
}

func TestCompile_OpCallFast(t *testing.T) {
	env := mock.Env{}
	program, err := expr.Compile("Fast(3, 2, 1)", expr.Env(env))
	require.NoError(t, err)
	require.Equal(t, vm.OpCallFast, program.Bytecode[4])
	require.Equal(t, 3, program.Arguments[4])
}
