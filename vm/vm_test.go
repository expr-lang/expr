package vm_test

import (
	"fmt"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/internal/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestRun_debug(t *testing.T) {
	var input = `[1, 2, 3]`

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
