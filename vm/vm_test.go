package vm_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/checker"
	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm"
)

func TestRun_NilProgram(t *testing.T) {
	_, err := vm.Run(nil, nil)
	require.Error(t, err)
}

func TestRun_ReuseVM(t *testing.T) {
	node, err := parser.Parse(`map(1..2, {#})`)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	reuse := vm.VM{}
	_, err = reuse.Run(program, nil)
	require.NoError(t, err)
	_, err = reuse.Run(program, nil)
	require.NoError(t, err)
}

func TestRun_ReuseVM_for_different_variables(t *testing.T) {
	v := vm.VM{}

	program, err := expr.Compile(`let a = 1; a + 1`)
	require.NoError(t, err)
	out, err := v.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, 2, out)

	program, err = expr.Compile(`let a = 2; a + 1`)
	require.NoError(t, err)
	out, err = v.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, 3, out)

	program, err = expr.Compile(`let a = 2; let b = 2; a + b`)
	require.NoError(t, err)
	out, err = v.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, 4, out)
}

func TestRun_Cast(t *testing.T) {
	tests := []struct {
		input  string
		expect reflect.Kind
		want   any
	}{
		{
			input:  `1`,
			expect: reflect.Float64,
			want:   float64(1),
		},
		{
			input:  `1`,
			expect: reflect.Int,
			want:   int(1),
		},
		{
			input:  `1`,
			expect: reflect.Int64,
			want:   int64(1),
		},
		{
			input:  `true`,
			expect: reflect.Bool,
			want:   true,
		},
		{
			input:  `false`,
			expect: reflect.Bool,
			want:   false,
		},
		{
			input:  `nil`,
			expect: reflect.Bool,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v %v", tt.expect, tt.input), func(t *testing.T) {
			tree, err := parser.Parse(tt.input)
			require.NoError(t, err)

			program, err := compiler.Compile(tree, &conf.Config{Expect: tt.expect})
			require.NoError(t, err)

			out, err := vm.Run(program, nil)
			require.NoError(t, err)

			require.Equal(t, tt.want, out)
		})
	}
}

func TestRun_Helpers(t *testing.T) {
	values := []any{
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		1,
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		float32(1),
		float64(1),
	}
	ops := []string{"+", "-", "*", "/", "%", "==", ">=", "<=", "<", ">"}

	for _, a := range values {
		for _, b := range values {
			for _, op := range ops {

				if op == "%" {
					switch a.(type) {
					case float32, float64:
						continue
					}
					switch b.(type) {
					case float32, float64:
						continue
					}
				}

				input := fmt.Sprintf("a %v b", op)
				env := map[string]any{
					"a": a,
					"b": b,
				}

				config := conf.CreateNew()

				tree, err := parser.Parse(input)
				require.NoError(t, err)

				_, err = checker.Check(tree, config)
				require.NoError(t, err)

				program, err := compiler.Compile(tree, config)
				require.NoError(t, err)

				_, err = vm.Run(program, env)
				require.NoError(t, err)
			}
		}
	}
}

type ErrorEnv struct {
	InnerEnv InnerEnv
}
type InnerEnv struct{}

func (ErrorEnv) WillError(param string) (bool, error) {
	if param == "yes" {
		return false, errors.New("error")
	}
	return true, nil
}

func (InnerEnv) WillError(param string) (bool, error) {
	if param == "yes" {
		return false, errors.New("inner error")
	}
	return true, nil
}

func TestRun_MethodWithError(t *testing.T) {
	input := `WillError("yes")`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := ErrorEnv{}
	funcConf := conf.New(env)
	_, err = checker.Check(tree, funcConf)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.EqualError(t, err, "error (1:1)\n | WillError(\"yes\")\n | ^")
	require.Equal(t, nil, out)

	selfErr := errors.Unwrap(err)
	require.NotNil(t, err)
	require.Equal(t, "error", selfErr.Error())
}

func TestRun_FastMethods(t *testing.T) {
	input := `hello() + world()`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := map[string]any{
		"hello": func(...any) any { return "hello " },
		"world": func(...any) any { return "world" },
	}
	funcConf := conf.New(env)
	_, err = checker.Check(tree, funcConf)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)

	require.Equal(t, "hello world", out)
}

func TestRun_InnerMethodWithError(t *testing.T) {
	input := `InnerEnv.WillError("yes")`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := ErrorEnv{}
	funcConf := conf.New(env)
	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.EqualError(t, err, "inner error (1:10)\n | InnerEnv.WillError(\"yes\")\n | .........^")
	require.Equal(t, nil, out)
}

func TestRun_InnerMethodWithError_NilSafe(t *testing.T) {
	input := `InnerEnv?.WillError("yes")`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := ErrorEnv{}
	funcConf := conf.New(env)
	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.EqualError(t, err, "inner error (1:11)\n | InnerEnv?.WillError(\"yes\")\n | ..........^")
	require.Equal(t, nil, out)
}

func TestRun_TaggedFieldName(t *testing.T) {
	input := `value`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := struct {
		V string `expr:"value"`
	}{
		V: "hello world",
	}

	funcConf := conf.New(env)
	_, err = checker.Check(tree, funcConf)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)

	require.Equal(t, "hello world", out)
}

func TestRun_OpInvalid(t *testing.T) {
	program := &vm.Program{
		Bytecode:  []vm.Opcode{vm.OpInvalid},
		Arguments: []int{0},
	}

	_, err := vm.Run(program, nil)
	require.EqualError(t, err, "invalid opcode")
}

func TestVM_OpcodeOperations(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		env         map[string]any
		want        any
		expectError string
	}{
		// Arithmetic Operations
		{
			name: "basic addition",
			expr: "2 + 3",
			want: 5,
		},
		{
			name: "mixed type arithmetic",
			expr: "2.5 + 3",
			want: 5.5,
		},
		{
			name: "chained arithmetic",
			expr: "1 + 2 * 3 - 4 / 2",
			want: 5.0,
		},
		{
			name: "modulo operation",
			expr: "5 % 2",
			want: 1,
		},
		{
			name: "exponent operation",
			expr: "2 ^ 3",
			want: 8.0,
		},
		{
			name: "negation",
			expr: "-5",
			want: -5,
		},

		// String Operations
		{
			name: "string concatenation",
			expr: `"hello" + " " + "world"`,
			want: "hello world",
		},
		{
			name: "string starts with",
			expr: `"hello world" startsWith "hello"`,
			want: true,
		},
		{
			name: "string ends with",
			expr: `"hello world" endsWith "world"`,
			want: true,
		},
		{
			name: "string contains",
			expr: `"hello world" contains "lo wo"`,
			want: true,
		},
		{
			name: "string matches regex",
			expr: `"hello123" matches "^hello\\d+$"`,
			want: true,
		},
		{
			name: "byte slice matches regex",
			expr: `b matches "^hello\\d+$"`,
			env:  map[string]any{"b": []byte("hello123")},
			want: true,
		},
		{
			name: "byte slice matches dynamic regex",
			expr: `b matches pattern`,
			env:  map[string]any{"b": []byte("hello123"), "pattern": "^hello\\d+$"},
			want: true,
		},

		// Data Structure Operations
		{
			name: "array creation and access",
			expr: "[1, 2, 3][1]",
			want: 2,
		},
		{
			name: "map creation and access",
			expr: `{"a": 1, "b": 2}.b`,
			want: 2,
		},
		{
			name: "array length",
			expr: "len([1, 2, 3])",
			want: 3,
		},
		{
			name: "array slice",
			expr: "[1, 2, 3, 4][1:3]",
			want: []any{2, 3},
		},
		{
			name: "array range",
			expr: "1..5",
			want: []int{1, 2, 3, 4, 5},
		},

		// Error Cases
		{
			name:        "invalid array index",
			expr:        "[1,2,3][5]",
			expectError: "index out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.Env(tt.env))
			require.NoError(t, err)

			testVM := &vm.VM{}
			got, err := testVM.Run(program, tt.env)

			if tt.expectError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestVM_GroupAndSortOperations(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		env         map[string]any
		want        any
		expectError string
	}{
		{
			name: "group by single field",
			expr: `groupBy([{"id": 1, "type": "a"}, {"id": 2, "type": "b"}, {"id": 3, "type": "a"}], #.type)`,
			want: map[any][]any{
				"a": {
					map[string]any{"id": 1, "type": "a"},
					map[string]any{"id": 3, "type": "a"},
				},
				"b": {
					map[string]any{"id": 2, "type": "b"},
				},
			},
		},
		{
			name: "sort by field ascending",
			expr: `sortBy([{"id": 3}, {"id": 1}, {"id": 2}], #.id)`,
			want: []any{
				map[string]any{"id": 1},
				map[string]any{"id": 2},
				map[string]any{"id": 3},
			},
		},
		{
			name: "sort by field descending",
			expr: `sortBy([{"id": 3}, {"id": 1}, {"id": 2}], #.id, "desc")`,
			want: []any{
				map[string]any{"id": 3},
				map[string]any{"id": 2},
				map[string]any{"id": 1},
			},
		},
		{
			name: "sort by computed value",
			expr: `sortBy([1, 2, 3, 4], # % 2)`,
			want: []any{2, 4, 1, 3},
		},
		{
			name: "group by with complex key",
			expr: `groupBy([1, 2, 3, 4, 5, 6], # % 2 == 0 ? "even" : "odd")`,
			want: map[any][]any{
				"even": {2, 4, 6},
				"odd":  {1, 3, 5},
			},
		},
		{
			name:        "invalid sort order",
			expr:        `sortBy([1, 2, 3], #, "invalid")`,
			expectError: "unknown order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.Env(tt.env))
			require.NoError(t, err)

			testVM := &vm.VM{}
			got, err := testVM.Run(program, tt.env)

			if tt.expectError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// TestVM_ProfileOperations tests the profiling opcodes
func TestVM_ProfileOperations(t *testing.T) {
	program := &vm.Program{
		Bytecode: []vm.Opcode{
			vm.OpProfileStart,
			vm.OpPush,
			vm.OpCall,
			vm.OpProfileEnd,
		},
		Arguments: []int{0, 1, 0, 0},
		Constants: []any{
			&vm.Span{},
			func() (any, error) {
				time.Sleep(time.Millisecond * 10)
				return nil, nil
			},
		},
	}

	testVM := &vm.VM{}
	_, err := testVM.Run(program, nil)
	require.NoError(t, err)

	span := program.Constants[0].(*vm.Span)
	require.Greater(t, span.Duration, time.Millisecond)
}

// TestVM_IndexOperations tests the index manipulation opcodes
func TestVM_IndexOperations(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want any
	}{
		{
			name: "decrement index in loop",
			expr: "reduce([1,2,3], #acc + #, 0)",
			want: 6,
		},
		{
			name: "set index in loop",
			expr: "map([1,2,3], # * 2)",
			want: []any{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr)
			require.NoError(t, err)

			testVM := &vm.VM{}
			got, err := testVM.Run(program, nil)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestVM_DirectCallOpcodes tests the specialized call opcodes directly
func TestVM_DirectCallOpcodes(t *testing.T) {
	tests := []struct {
		name     string
		bytecode []vm.Opcode
		args     []int
		consts   []any
		funcs    []vm.Function
		want     any
		wantErr  bool
	}{
		{
			name:     "OpCall0",
			bytecode: []vm.Opcode{vm.OpCall0},
			args:     []int{0},
			funcs: []vm.Function{
				func(args ...any) (any, error) {
					return 42, nil
				},
			},
			want: 42,
		},
		{
			name: "OpCall1",
			bytecode: []vm.Opcode{
				vm.OpPush,
				vm.OpCall1,
			},
			args:   []int{0, 0},
			consts: []any{10},
			funcs: []vm.Function{
				func(args ...any) (any, error) {
					return args[0].(int) * 2, nil
				},
			},
			want: 20,
		},
		{
			name: "OpCall2",
			bytecode: []vm.Opcode{
				vm.OpPush,
				vm.OpPush,
				vm.OpCall2,
			},
			args:   []int{0, 1, 0},
			consts: []any{10, 5},
			funcs: []vm.Function{
				func(args ...any) (any, error) {
					return args[0].(int) + args[1].(int), nil
				},
			},
			want: 15,
		},
		{
			name: "OpCall3",
			bytecode: []vm.Opcode{
				vm.OpPush,
				vm.OpPush,
				vm.OpPush,
				vm.OpCall3,
			},
			args:   []int{0, 1, 2, 0},
			consts: []any{10, 5, 2},
			funcs: []vm.Function{
				func(args ...any) (any, error) {
					return args[0].(int) + args[1].(int) + args[2].(int), nil
				},
			},
			want: 17,
		},
		{
			name: "OpCallN with error",
			bytecode: []vm.Opcode{
				vm.OpLoadFunc,
				vm.OpCallN,
			},
			args: []int{0, 0}, // Function index, number of args (0)
			funcs: []vm.Function{
				func(args ...any) (any, error) {
					return nil, fmt.Errorf("test error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := vm.NewProgram(
				file.Source{}, // source
				nil,           // node
				nil,           // locations
				0,             // variables
				tt.consts,
				tt.bytecode,
				tt.args,
				tt.funcs,
				nil, // debugInfo
				nil, // span
			)
			vm := &vm.VM{}
			got, err := vm.Run(program, nil)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestVM_CallN(t *testing.T) {
	input := `fn(1, 2, 3)`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := map[string]any{
		"fn": func(args ...any) (any, error) {
			sum := 0
			for _, arg := range args {
				sum += arg.(int)
			}
			return sum, nil
		},
	}

	config := conf.New(env)
	program, err := compiler.Compile(tree, config)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, 6, out)
}

// TestVM_IndexAndCountOperations tests the index and count manipulation opcodes directly
func TestVM_IndexAndCountOperations(t *testing.T) {
	tests := []struct {
		name     string
		bytecode []vm.Opcode
		args     []int
		consts   []any
		want     any
		wantErr  bool
	}{
		{
			name: "GetIndex",
			bytecode: []vm.Opcode{
				vm.OpPush,     // Push array to stack
				vm.OpBegin,    // Start scope
				vm.OpGetIndex, // Get current index
			},
			args:   []int{0, 0, 0},
			consts: []any{[]any{1, 2, 3}}, // Array for scope
			want:   0,                     // Initial index is 0
		},
		{
			name: "DecrementIndex",
			bytecode: []vm.Opcode{
				vm.OpPush,           // Push array to stack
				vm.OpBegin,          // Start scope
				vm.OpDecrementIndex, // Decrement index
				vm.OpGetIndex,       // Get current index
			},
			args:   []int{0, 0, 0, 0},
			consts: []any{[]any{1, 2, 3}}, // Array for scope
			want:   -1,                    // After decrement
		},
		{
			name: "GetCount",
			bytecode: []vm.Opcode{
				vm.OpPush,     // Push array to stack
				vm.OpBegin,    // Start scope
				vm.OpGetCount, // Get current count
			},
			args:   []int{0, 0, 0},
			consts: []any{[]any{1, 2, 3}}, // Array for scope
			want:   0,                     // Initial count is 0
		},
		{
			name: "IncrementCount",
			bytecode: []vm.Opcode{
				vm.OpPush,           // Push array to stack
				vm.OpBegin,          // Start scope
				vm.OpIncrementCount, // Increment count
				vm.OpGetCount,       // Get current count
			},
			args:   []int{0, 0, 0, 0},
			consts: []any{[]any{1, 2, 3}}, // Array for scope
			want:   1,                     // After increment
		},
		{
			name: "Multiple operations",
			bytecode: []vm.Opcode{
				vm.OpPush,           // Push array to stack
				vm.OpBegin,          // Start scope
				vm.OpIncrementCount, // Count = 1
				vm.OpIncrementCount, // Count = 2
				vm.OpDecrementIndex, // Index = -1
				vm.OpDecrementIndex, // Index = -2
				vm.OpGetCount,       // Push count (2)
				vm.OpGetIndex,       // Push index (-2)
				vm.OpAdd,            // Add them together
			},
			args:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			consts: []any{[]any{1, 2, 3}}, // Array for scope
			want:   0,                     // 2 + (-2) = 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := vm.NewProgram(
				file.Source{}, // source
				nil,           // node
				nil,           // locations
				0,             // variables
				tt.consts,
				tt.bytecode,
				tt.args,
				nil, // functions
				nil, // debugInfo
				nil, // span
			)
			vm := &vm.VM{}
			got, err := vm.Run(program, nil)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// TestVM_DirectBasicOpcodes tests basic opcodes directly
func TestVM_DirectBasicOpcodes(t *testing.T) {
	tests := []struct {
		name     string
		bytecode []vm.Opcode
		args     []int
		consts   []any
		env      any
		want     any
		wantErr  bool
	}{
		{
			name: "OpLoadEnv",
			bytecode: []vm.Opcode{
				vm.OpLoadEnv, // Load entire environment
			},
			args: []int{0},
			env:  map[string]any{"key": "value"},
			want: map[string]any{"key": "value"},
		},
		{
			name: "OpTrue",
			bytecode: []vm.Opcode{
				vm.OpTrue,
			},
			args: []int{0},
			want: true,
		},
		{
			name: "OpFalse",
			bytecode: []vm.Opcode{
				vm.OpFalse,
			},
			args: []int{0},
			want: false,
		},
		{
			name: "OpNil",
			bytecode: []vm.Opcode{
				vm.OpNil,
			},
			args: []int{0},
			want: nil,
		},
		{
			name: "OpNegate int",
			bytecode: []vm.Opcode{
				vm.OpPush,   // Push number
				vm.OpNegate, // Negate it
			},
			args:   []int{0, 0},
			consts: []any{42},
			want:   -42,
		},
		{
			name: "OpNegate float",
			bytecode: []vm.Opcode{
				vm.OpPush,   // Push number
				vm.OpNegate, // Negate it
			},
			args:   []int{0, 0},
			consts: []any{42.5},
			want:   -42.5,
		},
		{
			name: "OpNot true",
			bytecode: []vm.Opcode{
				vm.OpTrue, // Push true
				vm.OpNot,  // Negate it
			},
			args: []int{0, 0},
			want: false,
		},
		{
			name: "OpNot false",
			bytecode: []vm.Opcode{
				vm.OpFalse, // Push false
				vm.OpNot,   // Negate it
			},
			args: []int{0, 0},
			want: true,
		},
		{
			name: "OpNot error",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push non-bool
				vm.OpNot,  // Try to negate it
			},
			args:    []int{0, 0},
			consts:  []any{"not a bool"},
			wantErr: true,
		},
		{
			name: "OpEqualString equal",
			bytecode: []vm.Opcode{
				vm.OpPush,        // Push first string
				vm.OpPush,        // Push second string
				vm.OpEqualString, // Compare strings
			},
			args:   []int{0, 1, 0},
			consts: []any{"hello", "hello"},
			want:   true,
		},
		{
			name: "OpEqualString not equal",
			bytecode: []vm.Opcode{
				vm.OpPush,        // Push first string
				vm.OpPush,        // Push second string
				vm.OpEqualString, // Compare strings
			},
			args:   []int{0, 1, 0},
			consts: []any{"hello", "world"},
			want:   false,
		},
		{
			name: "OpEqualString with empty strings",
			bytecode: []vm.Opcode{
				vm.OpPush,        // Push first string
				vm.OpPush,        // Push second string
				vm.OpEqualString, // Compare strings
			},
			args:   []int{0, 1, 0},
			consts: []any{"", ""},
			want:   true,
		},
		{
			name: "OpEqualString type error",
			bytecode: []vm.Opcode{
				vm.OpPush,        // Push non-string
				vm.OpPush,        // Push string
				vm.OpEqualString, // Try to compare
			},
			args:    []int{0, 1, 0},
			consts:  []any{42, "hello"},
			wantErr: true,
		},
		{
			name: "OpInt",
			bytecode: []vm.Opcode{
				vm.OpInt, // Push int directly from args
			},
			args:   []int{42}, // The value 42 is passed directly in args
			consts: []any{},   // No constants needed
			want:   42,
		},
		{
			name: "OpInt negative",
			bytecode: []vm.Opcode{
				vm.OpInt, // Push negative int directly from args
			},
			args:   []int{-42}, // The value -42 is passed directly in args
			consts: []any{},    // No constants needed
			want:   -42,
		},
		{
			name: "OpInt zero",
			bytecode: []vm.Opcode{
				vm.OpInt, // Push zero directly from args
			},
			args:   []int{0}, // The value 0 is passed directly in args
			consts: []any{},  // No constants needed
			want:   0,
		},
		{
			name: "OpIn array true",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push element
				vm.OpPush, // Push array
				vm.OpIn,   // Check if element is in array
			},
			args:   []int{0, 1, 0},
			consts: []any{2, []any{1, 2, 3}},
			want:   true,
		},
		{
			name: "OpIn array false",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push element
				vm.OpPush, // Push array
				vm.OpIn,   // Check if element is in array
			},
			args:   []int{0, 1, 0},
			consts: []any{4, []any{1, 2, 3}},
			want:   false,
		},
		{
			name: "OpIn map true",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push key
				vm.OpPush, // Push map
				vm.OpIn,   // Check if key is in map
			},
			args:   []int{0, 1, 0},
			consts: []any{"b", map[string]any{"a": 1, "b": 2}},
			want:   true,
		},
		{
			name: "OpIn map false",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push key
				vm.OpPush, // Push map
				vm.OpIn,   // Check if key is in map
			},
			args:   []int{0, 1, 0},
			consts: []any{"c", map[string]any{"a": 1, "b": 2}},
			want:   false,
		},
		{
			name: "OpExponent integers",
			bytecode: []vm.Opcode{
				vm.OpPush,     // Push base
				vm.OpPush,     // Push exponent
				vm.OpExponent, // Calculate power
			},
			args:   []int{0, 1, 0},
			consts: []any{2, 3},
			want:   8.0,
		},
		{
			name: "OpExponent floats",
			bytecode: []vm.Opcode{
				vm.OpPush,     // Push base
				vm.OpPush,     // Push exponent
				vm.OpExponent, // Calculate power
			},
			args:   []int{0, 1, 0},
			consts: []any{2.0, 3.0},
			want:   8.0,
		},
		{
			name: "OpExponent negative exponent",
			bytecode: []vm.Opcode{
				vm.OpPush,     // Push base
				vm.OpPush,     // Push exponent
				vm.OpExponent, // Calculate power
			},
			args:   []int{0, 1, 0},
			consts: []any{2.0, -2.0},
			want:   0.25,
		},
		{
			name: "OpMatches valid regex",
			bytecode: []vm.Opcode{
				vm.OpPush,    // Push string
				vm.OpPush,    // Push pattern
				vm.OpMatches, // Match string against pattern
			},
			args:   []int{0, 1, 0},
			consts: []any{"hello123", "^hello\\d+$"},
			want:   true,
		},
		{
			name: "OpMatches non-matching regex",
			bytecode: []vm.Opcode{
				vm.OpPush,    // Push string
				vm.OpPush,    // Push pattern
				vm.OpMatches, // Match string against pattern
			},
			args:   []int{0, 1, 0},
			consts: []any{"hello", "^\\d+$"},
			want:   false,
		},
		{
			name: "OpMatches invalid regex",
			bytecode: []vm.Opcode{
				vm.OpPush,    // Push string
				vm.OpPush,    // Push pattern
				vm.OpMatches, // Match string against pattern
			},
			args:    []int{0, 1, 0},
			consts:  []any{"hello", "[invalid"},
			wantErr: true,
		},
		{
			name: "OpMatches type error",
			bytecode: []vm.Opcode{
				vm.OpPush,    // Push non-string
				vm.OpPush,    // Push pattern
				vm.OpMatches, // Match against pattern
			},
			args:    []int{0, 1, 0},
			consts:  []any{42, "^\\d+$"},
			wantErr: true,
		},
		{
			name: "OpCast int to float64",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push int
				vm.OpCast, // Cast to float64
			},
			args:   []int{0, 2},
			consts: []any{42},
			want:   float64(42),
		},
		{
			name: "OpCast int32 to int64",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push int32
				vm.OpCast, // Cast to int64
			},
			args:   []int{0, 1},
			consts: []any{int32(42)},
			want:   int64(42),
		},
		{
			name: "OpCast bool to bool",
			bytecode: []vm.Opcode{
				vm.OpTrue, // Push true
				vm.OpCast, // Cast to bool
			},
			args: []int{0, 3},
			want: true,
		},
		{
			name: "OpCast nil to bool",
			bytecode: []vm.Opcode{
				vm.OpNil,  // Push nil
				vm.OpCast, // Cast to bool
			},
			args: []int{0, 3},
			want: false,
		},
		{
			name: "OpCast int to bool",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push int
				vm.OpCast, // Cast to bool
			},
			args:    []int{0, 3},
			consts:  []any{1},
			wantErr: true,
		},
		{
			name: "OpCast invalid type",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push string
				vm.OpCast, // Try to cast to float64
			},
			args:    []int{0, 0},
			consts:  []any{"not a number"},
			wantErr: true,
		},
		{
			name: "OpLen array",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push array
				vm.OpLen,  // Get length
			},
			args:   []int{0, 0},
			consts: []any{[]any{1, 2, 3}},
			want:   3,
		},
		{
			name: "OpLen empty array",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push empty array
				vm.OpLen,  // Get length
			},
			args:   []int{0, 0},
			consts: []any{[]any{}},
			want:   0,
		},
		{
			name: "OpLen string",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push string
				vm.OpLen,  // Get length
			},
			args:   []int{0, 0},
			consts: []any{"hello"},
			want:   5,
		},
		{
			name: "OpLen empty string",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push empty string
				vm.OpLen,  // Get length
			},
			args:   []int{0, 0},
			consts: []any{""},
			want:   0,
		},
		{
			name: "OpLen map",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push map
				vm.OpLen,  // Get length
			},
			args:   []int{0, 0},
			consts: []any{map[string]any{"a": 1, "b": 2, "c": 3}},
			want:   3,
		},
		{
			name: "OpLen empty map",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push empty map
				vm.OpLen,  // Get length
			},
			args:   []int{0, 0},
			consts: []any{map[string]any{}},
			want:   0,
		},
		{
			name: "OpLen invalid type",
			bytecode: []vm.Opcode{
				vm.OpPush, // Push number
				vm.OpLen,  // Try to get length
			},
			args:    []int{0, 0},
			consts:  []any{42},
			wantErr: true,
		},
		{
			name: "OpThrow with string",
			bytecode: []vm.Opcode{
				vm.OpPush,  // Push error message
				vm.OpThrow, // Throw error
			},
			args:    []int{0, 0},
			consts:  []any{"test error"},
			wantErr: true,
		},
		{
			name: "OpThrow with error",
			bytecode: []vm.Opcode{
				vm.OpPush,  // Push error
				vm.OpThrow, // Throw error
			},
			args:    []int{0, 0},
			consts:  []any{fmt.Errorf("test error")},
			wantErr: true,
		},
		{
			name: "OpDefault",
			bytecode: []vm.Opcode{
				vm.OpEnd + 1, // OpEnd is always last, this is anunknown opcode
			},
			args:    []int{0, 0},
			consts:  []any{fmt.Errorf("test error")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := vm.NewProgram(
				file.Source{}, // source
				nil,           // node
				nil,           // locations
				0,             // variables
				tt.consts,
				tt.bytecode,
				tt.args,
				nil, // functions
				nil, // debugInfo
				nil, // span
			)
			vm := &vm.VM{}
			got, err := vm.Run(program, tt.env)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestVM_MemoryBudget(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		memBudget   uint
		expectError string
	}{
		{
			name:      "under budget",
			expr:      "map(1..10, #)",
			memBudget: 100,
		},
		{
			name:        "exceeds budget",
			expr:        "map(1..1000, #)",
			memBudget:   10,
			expectError: "memory budget exceeded",
		},
		{
			name:      "zero budget uses default",
			expr:      "map(1..10, #)",
			memBudget: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parser.Parse(tt.expr)
			require.NoError(t, err)

			program, err := compiler.Compile(node, nil)
			require.NoError(t, err)

			vm := vm.VM{MemoryBudget: tt.memBudget}
			out, err := vm.Run(program, nil)

			if tt.expectError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, out)
			}
		})
	}
}

// Helper functions for creating deeply nested expressions
func createNestedArithmeticExpr(t *testing.T, depth int) string {
	t.Helper()
	if depth == 0 {
		return "a"
	}
	return fmt.Sprintf("(%s + %d)", createNestedArithmeticExpr(t, depth-1), depth)
}

func createNestedMapExpr(t *testing.T, depth int) string {
	t.Helper()
	if depth == 0 {
		return `{"value": 1}`
	}
	return fmt.Sprintf(`{"nested": %s}`, createNestedMapExpr(t, depth-1))
}

func TestVM_Limits(t *testing.T) {
	tests := []struct {
		name         string
		expr         string
		memoryBudget uint
		maxNodes     uint
		env          map[string]any
		expectError  string
	}{
		{
			name:         "nested arithmetic allowed with max nodes and memory budget",
			expr:         createNestedArithmeticExpr(t, 100),
			env:          map[string]any{"a": 1},
			maxNodes:     1000,
			memoryBudget: 1, // arithmetic expressions not counted towards memory budget
		},
		{
			name:         "nested arithmetic blocked by max nodes",
			expr:         createNestedArithmeticExpr(t, 10000),
			env:          map[string]any{"a": 1},
			maxNodes:     100,
			memoryBudget: 1, // arithmetic expressions not counted towards memory budget
			expectError:  "compilation failed: expression exceeds maximum allowed nodes",
		},
		{
			name:         "nested map blocked by memory budget",
			expr:         createNestedMapExpr(t, 100),
			env:          map[string]any{},
			maxNodes:     1000,
			memoryBudget: 10, // Small memory budget to trigger limit
			expectError:  "memory budget exceeded",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var options []expr.Option
			options = append(options, expr.Env(test.env))
			if test.maxNodes > 0 {
				options = append(options, func(c *conf.Config) {
					c.MaxNodes = test.maxNodes
				})
			}

			program, err := expr.Compile(test.expr, options...)
			if err != nil {
				if test.expectError != "" && strings.Contains(err.Error(), test.expectError) {
					return
				}
				t.Fatal(err)
			}

			testVM := &vm.VM{
				MemoryBudget: test.memoryBudget,
			}

			_, err = testVM.Run(program, test.env)

			if test.expectError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectError)
			}
		})
	}
}

func TestVM_OpJump_NegativeOffset(t *testing.T) {
	program := vm.NewProgram(
		file.Source{},
		nil,
		nil,
		0,
		nil,
		[]vm.Opcode{
			vm.OpInt,
			vm.OpInt,
			vm.OpJump,
			vm.OpInt,
			vm.OpJump,
		},
		[]int{
			1,
			2,
			-2, // negative offset for a forward jump opcode
			3,
			-2,
		},
		nil,
		nil,
		nil,
	)

	_, err := vm.Run(program, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "negative jump offset is invalid")
}

func TestVM_StackUnderflow(t *testing.T) {
	tests := []struct {
		name        string
		bytecode    []vm.Opcode
		args        []int
		expectError string
	}{
		{
			name:     "pop after push",
			bytecode: []vm.Opcode{vm.OpInt, vm.OpPop},
			args:     []int{42, 0},
		},
		{
			name:        "underflow after valid operations",
			bytecode:    []vm.Opcode{vm.OpInt, vm.OpInt, vm.OpPop, vm.OpPop, vm.OpPop},
			args:        []int{1, 2, 0, 0, 0},
			expectError: "stack underflow",
		},
		{
			name:        "pop on empty stack",
			bytecode:    []vm.Opcode{vm.OpPop},
			args:        []int{0},
			expectError: "stack underflow",
		},
		{
			name:     "pop after push",
			bytecode: []vm.Opcode{vm.OpInt, vm.OpPop},
			args:     []int{123, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := &vm.Program{
				Bytecode:  tt.bytecode,
				Arguments: tt.args,
				Constants: []any{},
			}

			_, err := vm.Run(program, nil)
			if tt.expectError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
