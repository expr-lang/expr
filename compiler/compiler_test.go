package compiler_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"

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
					vm.OpJumpIfNotNil,
					vm.OpPop,
					vm.OpNil,
				},
				Arguments: []int{0, 1, 1, 2, 0, 0},
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
					vm.OpJumpIfNotNil,
					vm.OpPop,
					vm.OpNil,
				},
				Arguments: []int{0, 1, 1, 2, 0, 0},
			},
		},
		{
			`A?.B`,
			vm.Program{
				Constants: []any{
					&runtime.Field{
						Index: []int{0},
						Path:  []string{"A"},
					},
					&runtime.Field{
						Index: []int{1},
						Path:  []string{"B"},
					},
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNil,
					vm.OpFetchField,
					vm.OpJumpIfNotNil,
					vm.OpPop,
					vm.OpNil,
				},
				Arguments: []int{0, 1, 1, 2, 0, 0},
			},
		},
		{
			`A?.B ?? 42`,
			vm.Program{
				Constants: []any{
					&runtime.Field{
						Index: []int{0},
						Path:  []string{"A"},
					},
					&runtime.Field{
						Index: []int{1},
						Path:  []string{"B"},
					},
					42,
				},
				Bytecode: []vm.Opcode{
					vm.OpLoadField,
					vm.OpJumpIfNil,
					vm.OpFetchField,
					vm.OpJumpIfNotNil,
					vm.OpPop,
					vm.OpPush,
				},
				Arguments: []int{0, 1, 1, 2, 0, 2},
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
				Arguments: []int{0, 0, 1, 12},
			},
		},
		{
			`1; 2; 3`,
			vm.Program{
				Constants: []any{
					1,
					2,
					3,
				},
				Bytecode: []vm.Opcode{
					vm.OpPush,
					vm.OpPop,
					vm.OpPush,
					vm.OpPop,
					vm.OpPush,
				},
				Arguments: []int{0, 0, 1, 0, 2},
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
	require.Equal(t, 32, program.Arguments[3])
}

func TestCompile_FuncTypes_with_Method(t *testing.T) {
	env := mock.Env{}
	program, err := expr.Compile("FuncTyped('bar')", expr.Env(env))
	require.NoError(t, err)
	require.Equal(t, vm.OpCallTyped, program.Bytecode[2])
	require.Equal(t, 76, program.Arguments[2])
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

func TestCompile_optimizes_jumps(t *testing.T) {
	env := map[string]any{
		"a":   true,
		"b":   true,
		"c":   true,
		"d":   true,
		"i64": int64(1),
	}
	tests := []struct {
		code string
		want string
	}{
		{
			`let foo = true; let bar = false; let baz = true; foo || bar || baz`,
			`0   OpTrue
1   OpStore  <0>  foo
2   OpFalse
3   OpStore  <1>  bar
4   OpTrue
5   OpStore       <2>  baz
6   OpLoadVar     <0>  foo
7   OpJumpIfTrue  <5>  (13)
8   OpPop
9   OpLoadVar     <1>  bar
10  OpJumpIfTrue  <2>  (13)
11  OpPop
12  OpLoadVar  <2>  baz
`,
		},
		{
			`a && b && c`,
			`0  OpLoadFast     <0>  a
1  OpJumpIfFalse  <5>  (7)
2  OpPop
3  OpLoadFast     <1>  b
4  OpJumpIfFalse  <2>  (7)
5  OpPop
6  OpLoadFast  <2>  c
`,
		},
		{
			`a && b || c && d`,
			`0  OpLoadFast     <0>  a
1  OpJumpIfFalse  <2>  (4)
2  OpPop
3  OpLoadFast    <1>  b
4  OpJumpIfTrue  <5>  (10)
5  OpPop
6  OpLoadFast     <2>  c
7  OpJumpIfFalse  <2>  (10)
8  OpPop
9  OpLoadFast  <3>  d
`,
		},
		{
			`filter([1, 2, 3, 4, 5], # > 3 && # != 4 && # != 5)`,
			`0   OpPush  <0>  [1 2 3 4 5]
1   OpBegin
2   OpJumpIfEnd  <23>  (26)
3   OpPointer
4   OpPush  <1>  3
5   OpMore
6   OpJumpIfFalse  <16>  (23)
7   OpPop
8   OpPointer
9   OpPush  <2>  4
10  OpEqualInt
11  OpNot
12  OpJumpIfFalse  <10>  (23)
13  OpPop
14  OpPointer
15  OpPush  <3>  5
16  OpEqualInt
17  OpNot
18  OpJumpIfFalse  <4>  (23)
19  OpPop
20  OpIncrementCount
21  OpPointer
22  OpJump  <1>  (24)
23  OpPop
24  OpIncrementIndex
25  OpJumpBackward  <24>  (2)
26  OpGetCount
27  OpEnd
28  OpArray
`,
		},
		{
			`let foo = true; let bar = false; let baz = true; foo && bar || baz`,
			`0   OpTrue
1   OpStore  <0>  foo
2   OpFalse
3   OpStore  <1>  bar
4   OpTrue
5   OpStore        <2>  baz
6   OpLoadVar      <0>  foo
7   OpJumpIfFalse  <2>  (10)
8   OpPop
9   OpLoadVar     <1>  bar
10  OpJumpIfTrue  <2>  (13)
11  OpPop
12  OpLoadVar  <2>  baz
`,
		},
		{
			`true ?? nil ?? nil ?? nil`,
			`0  OpTrue
1  OpJumpIfNotNil  <8>  (10)
2  OpPop
3  OpNil
4  OpJumpIfNotNil  <5>  (10)
5  OpPop
6  OpNil
7  OpJumpIfNotNil  <2>  (10)
8  OpPop
9  OpNil
`,
		},
		{
			`let m = {"a": {"b": {"c": 1}}}; m?.a?.b?.c`,
			`0   OpPush  <0>  a
1   OpPush  <1>  b
2   OpPush  <2>  c
3   OpPush  <3>  1
4   OpPush  <3>  1
5   OpMap
6   OpPush  <3>  1
7   OpMap
8   OpPush  <3>  1
9   OpMap
10  OpStore      <0>  m
11  OpLoadVar    <0>  m
12  OpJumpIfNil  <8>  (21)
13  OpPush       <0>  a
14  OpFetch
15  OpJumpIfNil  <5>  (21)
16  OpPush       <1>  b
17  OpFetch
18  OpJumpIfNil  <2>  (21)
19  OpPush       <2>  c
20  OpFetch
21  OpJumpIfNotNil  <2>  (24)
22  OpPop
23  OpNil
`,
		},
		{
			`-1 not in [1, 2, 5]`,
			`0  OpPush  <0>  -1
1  OpPush  <1>  map[1:{} 2:{} 5:{}]
2  OpIn
3  OpNot
`,
		},
		{
			`1 + 8 not in [1, 2, 5]`,
			`0  OpPush  <0>  9
1  OpPush  <1>  map[1:{} 2:{} 5:{}]
2  OpIn
3  OpNot
`,
		},
		{
			`true ? false : 8 not in [1, 2, 5]`,
			`0  OpTrue
1  OpJumpIfFalse  <3>  (5)
2  OpPop
3  OpFalse
4  OpJump  <5>  (10)
5  OpPop
6  OpPush  <0>  8
7  OpPush  <1>  map[1:{} 2:{} 5:{}]
8  OpIn
9  OpNot
`,
		},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			program, err := expr.Compile(test.code, expr.Env(env))
			require.NoError(t, err)
			require.Equal(t, test.want, program.Disassemble())
		})
	}
}

func TestCompile_IntegerArgsFunc(t *testing.T) {
	env := mock.Env{}
	tests := []struct{ code string }{
		{"FuncInt(0)"},
		{"FuncInt8(0)"},
		{"FuncInt16(0)"},
		{"FuncInt32(0)"},
		{"FuncInt64(0)"},
		{"FuncUint(0)"},
		{"FuncUint8(0)"},
		{"FuncUint16(0)"},
		{"FuncUint32(0)"},
		{"FuncUint64(0)"},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			_, err := expr.Compile(tt.code, expr.Env(env))
			require.NoError(t, err)
		})
	}
}

func TestCompile_call_on_nil(t *testing.T) {
	env := map[string]any{
		"foo": nil,
	}
	_, err := expr.Compile(`foo()`, expr.Env(env))
	require.Error(t, err)
	require.Contains(t, err.Error(), "foo is nil; cannot call nil as function")
}

func TestCompile_Expect(t *testing.T) {
	tests := []struct {
		input  string
		option expr.Option
		op     vm.Opcode
		arg    int
	}{
		{
			input:  "1",
			option: expr.AsKind(reflect.Int),
			op:     vm.OpCast,
			arg:    0,
		},
		{
			input:  "1",
			option: expr.AsInt64(),
			op:     vm.OpCast,
			arg:    1,
		},
		{
			input:  "1",
			option: expr.AsFloat64(),
			op:     vm.OpCast,
			arg:    2,
		},
		{
			input:  "true",
			option: expr.AsBool(),
			op:     vm.OpCast,
			arg:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := expr.Compile(tt.input, tt.option)
			require.NoError(t, err)

			lastOp := program.Bytecode[len(program.Bytecode)-1]
			lastArg := program.Arguments[len(program.Arguments)-1]

			assert.Equal(t, tt.op, lastOp)
			assert.Equal(t, tt.arg, lastArg)
		})
	}
}
