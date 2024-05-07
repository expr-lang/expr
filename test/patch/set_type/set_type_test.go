package set_type_test

import (
	"reflect"
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

func TestPatch_SetType(t *testing.T) {
	_, err := expr.Compile(
		`Value + "string"`,
		expr.Env(Env{}),
		expr.Function(
			"getValue",
			func(params ...any) (any, error) {
				return params[0].(Value).Int, nil
			},
			// We can set function type right here,
			// but we want to check what SetType in
			// getValuePatcher will take an effect.
		),
		expr.Patch(getValuePatcher{}),
	)
	require.Error(t, err)
}

type Value struct {
	Int int
}

type Env struct {
	Value Value
}

var valueType = reflect.TypeOf((*Value)(nil)).Elem()

type getValuePatcher struct{}

func (getValuePatcher) Visit(node *ast.Node) {
	id, ok := (*node).(*ast.IdentifierNode)
	if !ok {
		return
	}
	if id.Type() == valueType {
		newNode := &ast.CallNode{
			Callee:    &ast.IdentifierNode{Value: "getValue"},
			Arguments: []ast.Node{id},
		}
		newNode.SetType(reflect.TypeOf(0))
		ast.Patch(node, newNode)
	}
}
