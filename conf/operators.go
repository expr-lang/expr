package conf

import (
	"reflect"

	"github.com/expr-lang/expr/ast"
)

// OperatorsTable maps binary operators to corresponding list of functions.
// Functions should be provided in the environment to allow operator overloading.
type OperatorsTable map[string][]string

func FindSuitableOperatorOverload(fns []string, types TypesTable, funcs FunctionTable, l, r reflect.Type) (reflect.Type, string, bool) {
	t, fn, ok := FindSuitableOperatorOverloadInFunctions(fns, funcs, l, r)
	if !ok {
		t, fn, ok = FindSuitableOperatorOverloadInTypes(fns, types, l, r)
	}
	return t, fn, ok
}

func FindSuitableOperatorOverloadInTypes(fns []string, types TypesTable, l, r reflect.Type) (reflect.Type, string, bool) {
	for _, fn := range fns {
		fnType, ok := types[fn]
		if !ok {
			continue
		}
		firstInIndex := 0
		if fnType.Method {
			firstInIndex = 1 // As first argument to method is receiver.
		}
		ret, done := checkTypeSuits(fnType.Type, l, r, firstInIndex)
		if done {
			return ret, fn, true
		}
	}
	return nil, "", false
}

func FindSuitableOperatorOverloadInFunctions(fns []string, funcs FunctionTable, l, r reflect.Type) (reflect.Type, string, bool) {
	for _, fn := range fns {
		fnType, ok := funcs[fn]
		if !ok {
			continue
		}
		firstInIndex := 0
		for _, overload := range fnType.Types {
			ret, done := checkTypeSuits(overload, l, r, firstInIndex)
			if done {
				return ret, fn, true
			}
		}
	}
	return nil, "", false
}

func checkTypeSuits(t reflect.Type, l reflect.Type, r reflect.Type, firstInIndex int) (reflect.Type, bool) {
	firstArgType := t.In(firstInIndex)
	secondArgType := t.In(firstInIndex + 1)

	firstArgumentFit := l == firstArgType || (firstArgType.Kind() == reflect.Interface && (l == nil || l.Implements(firstArgType)))
	secondArgumentFit := r == secondArgType || (secondArgType.Kind() == reflect.Interface && (r == nil || r.Implements(secondArgType)))
	if firstArgumentFit && secondArgumentFit {
		return t.Out(0), true
	}
	return nil, false
}

type OperatorPatcher struct {
	Operators OperatorsTable
	Types     TypesTable
	Functions FunctionTable
}

func (p *OperatorPatcher) Visit(node *ast.Node) {
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

	ret, fn, ok := FindSuitableOperatorOverload(fns, p.Types, p.Functions, leftType, rightType)
	if ok {
		newNode := &ast.CallNode{
			Callee:    &ast.IdentifierNode{Value: fn},
			Arguments: []ast.Node{binaryNode.Left, binaryNode.Right},
		}
		newNode.SetType(ret)
		ast.Patch(node, newNode)
	}
}
