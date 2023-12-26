package checker

import (
	"reflect"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/conf"
)

func FieldIndex(types conf.TypesTable, node ast.Node) (bool, []int, string) {
	switch n := node.(type) {
	case *ast.IdentifierNode:
		if t, ok := types[n.Value]; ok && len(t.FieldIndex) > 0 {
			return true, t.FieldIndex, n.Value
		}
	case *ast.MemberNode:
		base := n.Node.Type()
		if kind(base) == reflect.Ptr {
			base = base.Elem()
		}
		if kind(base) == reflect.Struct {
			if prop, ok := n.Property.(*ast.StringNode); ok {
				name := prop.Value
				if field, ok := fetchField(base, name); ok {
					return true, field.Index, name
				}
			}
		}
	}
	return false, nil, ""
}

func MethodIndex(types conf.TypesTable, node ast.Node) (bool, int, string) {
	switch n := node.(type) {
	case *ast.IdentifierNode:
		if t, ok := types[n.Value]; ok {
			return t.Method, t.MethodIndex, n.Value
		}
	case *ast.MemberNode:
		if name, ok := n.Property.(*ast.StringNode); ok {
			base := n.Node.Type()
			if base != nil && base.Kind() != reflect.Interface {
				if m, ok := base.MethodByName(name.Value); ok {
					return true, m.Index, name.Value
				}
			}
		}
	}
	return false, 0, ""
}
