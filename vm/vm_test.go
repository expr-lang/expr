package vm_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/require"
)

func TestRun_NilProgram(t *testing.T) {
	_, err := vm.Run(nil, nil)
	require.Error(t, err)
}

func TestRun_Debugger(t *testing.T) {
	input := `[1, 2]`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	debug := vm.Debug()
	go func() {
		debug.Step()
		debug.Step()
		debug.Step()
		debug.Step()
	}()
	go func() {
		for range debug.Position() {
		}
	}()

	_, err = debug.Run(program, nil)
	require.NoError(t, err)
	require.Len(t, debug.Stack(), 0)
	require.Nil(t, debug.Scope())
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

func TestRun_Cast(t *testing.T) {
	input := `1`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, &conf.Config{Expect: reflect.Float64})
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)

	require.Equal(t, float64(1), out)
}

func TestRun_Helpers(t *testing.T) {
	values := []interface{}{
		uint(3),
		uint8(3),
		uint16(3),
		uint32(3),
		uint64(3),
		3,
		int8(3),
		int16(3),
		int32(3),
		int64(3),
		float32(3),
		float64(3),
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
				t.Run(fmt.Sprintf("%v %v %s %v %v", reflect.TypeOf(a), a, op, reflect.TypeOf(b), b), func(t *testing.T) {
					input := fmt.Sprintf("a %v b", op)
					env := map[string]interface{}{
						"a": a,
						"b": b,
					}
					// Run.
					tree, err := parser.Parse(input)
					require.NoError(t, err, "parsing input should not fail")
					_, err = checker.Check(tree, nil)
					require.NoError(t, err, "parsed tree should be valid")
					program, err := compiler.Compile(tree, nil)
					require.NoError(t, err, "should compile")
					result, err := vm.Run(program, env)
					require.NoError(t, err, "run should not fail")
					// Parse and check result.
					var resultStr string
					switch result := result.(type) {
					case time.Duration:
						resultStr = fmt.Sprintf("%v", int64(result))
					default:
						resultStr = fmt.Sprintf("%v", result)
					}
					switch op {
					case "+":
						assert.Equal(t, "6", resultStr, "should return correct result")
					case "-":
						assert.Equal(t, "0", resultStr, "should return correct result")
					case "*":
						assert.Equal(t, "9", resultStr, "should return correct result")
					case "/":
						assert.Equal(t, "1", resultStr, "should return correct result")
					case "%":
						assert.Equal(t, "0", resultStr, "should return correct result")
					case "==":
						assert.Equal(t, "true", resultStr, "should return correct result")
					case ">=":
						assert.Equal(t, "true", resultStr, "should return correct result")
					case "<=":
						assert.Equal(t, "true", resultStr, "should return correct result")
					case "<":
						assert.Equal(t, "false", resultStr, "should return correct result")
					case ">":
						assert.Equal(t, "false", resultStr, "should return correct result")
					default:
						t.Error("unhandled operation in test")
					}
				})
			}
		}
	}
}

func TestRun_Helpers_Time(t *testing.T) {
	testTime := time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	testDuration := time.Duration(1)

	tests := []struct {
		a       interface{}
		b       interface{}
		op      string
		want    interface{}
		wantErr bool
	}{
		{a: testTime, b: testTime, op: "<", wantErr: false, want: false},
		{a: testTime, b: testTime, op: ">", wantErr: false, want: false},
		{a: testTime, b: testTime, op: "<=", wantErr: false, want: true},
		{a: testTime, b: testTime, op: ">=", wantErr: false, want: true},
		{a: testTime, b: testTime, op: "==", wantErr: false, want: true},
		{a: testTime, b: testTime, op: "!=", wantErr: false, want: false},
		{a: testTime, b: testTime, op: "-", wantErr: false},
		{a: testTime, b: testDuration, op: "+", wantErr: false},

		// error cases
		{a: testTime, b: int64(1), op: "<", wantErr: true},
		{a: testTime, b: float64(1), op: "<", wantErr: true},
		{a: testTime, b: testDuration, op: "<", wantErr: true},

		{a: testTime, b: int64(1), op: ">", wantErr: true},
		{a: testTime, b: float64(1), op: ">", wantErr: true},
		{a: testTime, b: testDuration, op: ">", wantErr: true},

		{a: testTime, b: int64(1), op: "<=", wantErr: true},
		{a: testTime, b: float64(1), op: "<=", wantErr: true},
		{a: testTime, b: testDuration, op: "<=", wantErr: true},

		{a: testTime, b: int64(1), op: ">=", wantErr: true},
		{a: testTime, b: float64(1), op: ">=", wantErr: true},
		{a: testTime, b: testDuration, op: ">=", wantErr: true},

		{a: testTime, b: int64(1), op: "==", wantErr: false, want: false},
		{a: testTime, b: float64(1), op: "==", wantErr: false, want: false},
		{a: testTime, b: testDuration, op: "==", wantErr: false, want: false},

		{a: testTime, b: int64(1), op: "!=", wantErr: false, want: true},
		{a: testTime, b: float64(1), op: "!=", wantErr: false, want: true},
		{a: testTime, b: testDuration, op: "!=", wantErr: false, want: true},

		{a: testTime, b: int64(1), op: "-", wantErr: true},
		{a: testTime, b: float64(1), op: "-", wantErr: true},
		{a: testTime, b: testDuration, op: "-", wantErr: true},

		{a: testTime, b: testTime, op: "+", wantErr: true},
		{a: testTime, b: int64(1), op: "+", wantErr: true},
		{a: testTime, b: float64(1), op: "+", wantErr: true},
		{a: testDuration, b: testTime, op: "+", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("time helper test `%T %s %T`", tt.a, tt.op, tt.b), func(t *testing.T) {
			input := fmt.Sprintf("a %v b", tt.op)
			env := map[string]interface{}{
				"a": tt.a,
				"b": tt.b,
			}

			tree, err := parser.Parse(input)
			require.NoError(t, err)

			_, err = checker.Check(tree, nil)
			require.NoError(t, err)

			program, err := compiler.Compile(tree, nil)
			require.NoError(t, err)

			got, err := vm.Run(program, env)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.want != nil {
					require.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func TestRun_MemoryBudget(t *testing.T) {
	input := `map(1..100, {map(1..100, {map(1..100, {0})})})`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, nil)
	require.NoError(t, err)

	_, err = vm.Run(program, nil)
	require.Error(t, err)
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

func (ErrorEnv) FastError(...interface{}) interface{} {
	return true
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
}

func TestRun_FastMethods(t *testing.T) {
	input := `hello() + world()`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := map[string]interface{}{
		"hello": func(...interface{}) interface{} { return "hello " },
		"world": func(...interface{}) interface{} { return "world" },
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

func TestRun_FastMethodWithError(t *testing.T) {
	input := `FastError()`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := ErrorEnv{}
	funcConf := conf.New(env)
	_, err = checker.Check(tree, funcConf)
	require.NoError(t, err)
	require.True(t, tree.Node.(*ast.CallNode).Fast, "method must be fast")

	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)

	require.Equal(t, true, out)
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
