package patcher

import (
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/conf"
)

type Operator struct {
	Operators conf.OperatorsTable
	Types     conf.TypesTable
	Functions conf.FunctionTable
}

func (p *Operator) Visit(node *ast.Node) {
	binaryNode, ok := (*node).(*ast.BinaryNode)
	if !ok {
		return
	}

	fns, ok := p.Operators[binaryNode.Operator]
	if !ok {
		return
	}

	leftType := binaryNode.Left.Type()
	rightType := binaryNode.Right.Type()

	ret, fn, ok := conf.FindSuitableOperatorOverload(fns, p.Types, p.Functions, leftType, rightType)
	if ok {
		newNode := &ast.CallNode{
			Callee:    &ast.IdentifierNode{Value: fn},
			Arguments: []ast.Node{binaryNode.Left, binaryNode.Right},
		}
		newNode.SetType(ret)
		ast.Patch(node, newNode)
	}
}
