package patch_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/vm"
	"github.com/expr-lang/expr/vm/runtime"
)

func TestPatch_change_ident(t *testing.T) {
	program, err := expr.Compile(
		`foo`,
		expr.Env(Env{}),
		expr.Patch(changeIdent{}),
	)
	require.NoError(t, err)

	expected := &vm.Program{
		Bytecode: []vm.Opcode{
			vm.OpLoadField,
		},
		Arguments: []int{
			0,
		},
		Constants: []any{
			&runtime.Field{
				Path:  []string{"bar"},
				Index: []int{1},
			},
		},
	}

	require.Equal(t, expected.Disassemble(), program.Disassemble())
}

type Env struct {
	Foo int `expr:"foo"`
	Bar int `expr:"bar"`
}

type changeIdent struct{}

func (changeIdent) Visit(node *ast.Node) {
	id, ok := (*node).(*ast.IdentifierNode)
	if !ok {
		return
	}
	if id.Value == "foo" {
		// A correct way to patch the node:
		//
		//	newNode := &ast.IdentifierNode{Value: "bar"}
		//	ast.Patch(node, newNode)
		//
		// But we can do it in a wrong way:
		id.Value = "bar"
	}
}
