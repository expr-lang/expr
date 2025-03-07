//go:build expr_debug
// +build expr_debug

package vm_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm"
)

func TestDebugger(t *testing.T) {
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
	require.Len(t, debug.Stack, 0)
	require.Nil(t, debug.Scopes)
}
