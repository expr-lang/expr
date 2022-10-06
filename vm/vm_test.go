package vm_test

import (
	"errors"
	"fmt"
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

func TestRun_debug(t *testing.T) {
	input := `[1, 2, 3]`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	_, err = vm.Run(program, nil)
	require.NoError(t, err)
}

func TestRun_cast(t *testing.T) {
	input := `1`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, &conf.Config{Expect: reflect.Float64})
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)

	require.Equal(t, float64(1), out)
}

func TestRun_helpers(t *testing.T) {
	values := []interface{}{
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		int(1),
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
				env := map[string]interface{}{
					"a": a,
					"b": b,
				}

				tree, err := parser.Parse(input)
				require.NoError(t, err)

				_, err = checker.Check(tree, nil)
				require.NoError(t, err)

				program, err := compiler.Compile(tree, nil)
				require.NoError(t, err)

				_, err = vm.Run(program, env)
				require.NoError(t, err)
			}
		}
	}
}

func TestRun_helpers_time(t *testing.T) {
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

func TestRun_memory_budget(t *testing.T) {
	input := `map(1..100, {map(1..100, {map(1..100, {0})})})`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, nil)
	require.NoError(t, err)

	_, err = vm.Run(program, nil)
	require.Error(t, err)
}

func TestRun_fast_function_with_error(t *testing.T) {
	input := `WillError()`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := map[string]interface{}{
		"WillError": func(...interface{}) (interface{}, error) { return 1, errors.New("error") },
	}
	funcConf := conf.New(env)
	_, err = checker.Check(tree, funcConf)
	require.NoError(t, err)

	require.True(t, tree.Node.(*ast.FunctionNode).Fast, "function must be fast")
	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.EqualError(t, err, "error")

	require.Equal(t, nil, out)
}

type ErrorEnv struct {
	InnerEnv InnerEnv
}
type InnerEnv struct{}

func (ErrorEnv) WillError() (bool, error) {
	return false, errors.New("method error")
}

func (ErrorEnv) FastError(...interface{}) (interface{}, error) {
	return true, nil
}

func (InnerEnv) WillError() (bool, error) {
	return false, errors.New("inner error")
}

func TestRun_method_with_error(t *testing.T) {
	input := `WillError()`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := ErrorEnv{}
	funcConf := conf.New(env)
	_, err = checker.Check(tree, funcConf)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.EqualError(t, err, "method error")

	require.Equal(t, nil, out)
}

func TestRun_fast_methods(t *testing.T) {
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

func TestRun_fast_method_with_error(t *testing.T) {
	input := `FastError()`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := ErrorEnv{}
	funcConf := conf.New(env)
	_, err = checker.Check(tree, funcConf)
	require.NoError(t, err)
	require.True(t, tree.Node.(*ast.FunctionNode).Fast, "method must be fast")

	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)

	require.Equal(t, true, out)
}

func TestRun_inner_method_with_error(t *testing.T) {
	input := `InnerEnv.WillError()`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	env := ErrorEnv{}
	funcConf := conf.New(env)
	program, err := compiler.Compile(tree, funcConf)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.EqualError(t, err, "inner error")

	require.Equal(t, nil, out)
}

func TestRun_tagged_field_name(t *testing.T) {
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
