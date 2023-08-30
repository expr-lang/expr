package vm_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/file"
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

func (ErrorEnv) FastError(...any) any {
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

func TestRun_OpInvalid(t *testing.T) {
	program := &vm.Program{
		Locations: []file.Location{{0, 0}},
		Bytecode:  []vm.Opcode{vm.OpInvalid},
		Arguments: []int{0},
	}

	_, err := vm.Run(program, nil)
	require.EqualError(t, err, "invalid opcode")
}
