package ast_test

import (
	"expr/ast"
)

type visitor struct {
	identifiers []string
}

func (v *visitor) Visit(node *ast.Node) {
	if n, ok := (*node).(*ast.IdentifierNode); ok {
		v.identifiers = append(v.identifiers, n.Value)
	}
}

type patcher struct{}

func (p *patcher) Visit(node *ast.Node) {
	if _, ok := (*node).(*ast.IdentifierNode); ok {
		*node = &ast.NilNode{}
	}
}
