package patcher

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/builtin"
	"github.com/expr-lang/expr/conf"
)

type OperatorOverride struct {
	Operator  string              // Operator token to override.
	Overrides []string            // List of function names to override operator with.
	Types     conf.TypesTable     // Env types.
	Functions conf.FunctionsTable // Env functions.
}

func (p *OperatorOverride) Visit(node *ast.Node) {
	binaryNode, ok := (*node).(*ast.BinaryNode)
	if !ok {
		return
	}

	if binaryNode.Operator != p.Operator {
		return
	}

	leftType := binaryNode.Left.Type()
	rightType := binaryNode.Right.Type()

	ret, fn, ok := p.FindSuitableOperatorOverload(leftType, rightType)
	if ok {
		newNode := &ast.CallNode{
			Callee:    &ast.IdentifierNode{Value: fn},
			Arguments: []ast.Node{binaryNode.Left, binaryNode.Right},
		}
		newNode.SetType(ret)
		ast.Patch(node, newNode)
	}
}

func (p *OperatorOverride) FindSuitableOperatorOverload(l, r reflect.Type) (reflect.Type, string, bool) {
	t, fn, ok := p.findSuitableOperatorOverloadInFunctions(l, r)
	if !ok {
		t, fn, ok = p.findSuitableOperatorOverloadInTypes(l, r)
	}
	return t, fn, ok
}

func (p *OperatorOverride) findSuitableOperatorOverloadInTypes(l, r reflect.Type) (reflect.Type, string, bool) {
	for _, fn := range p.Overrides {
		fnType, ok := p.Types[fn]
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

func (p *OperatorOverride) findSuitableOperatorOverloadInFunctions(l, r reflect.Type) (reflect.Type, string, bool) {
	for _, fn := range p.Overrides {
		fnType, ok := p.Functions[fn]
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

func (p *OperatorOverride) Check() {
	for _, fn := range p.Overrides {
		fnType, foundType := p.Types[fn]
		fnFunc, foundFunc := p.Functions[fn]
		if !foundFunc && (!foundType || fnType.Type.Kind() != reflect.Func) {
			panic(fmt.Errorf("function %s for %s operator does not exist in the environment", fn, p.Operator))
		}

		if foundType {
			checkType(fnType, fn, p.Operator)
		}

		if foundFunc {
			checkFunc(fnFunc, fn, p.Operator)
		}
	}
}

func checkType(fnType conf.Tag, fn string, operator string) {
	requiredNumIn := 2
	if fnType.Method {
		requiredNumIn = 3 // As first argument of method is receiver.
	}
	if fnType.Type.NumIn() != requiredNumIn || fnType.Type.NumOut() != 1 {
		panic(fmt.Errorf("function %s for %s operator does not have a correct signature", fn, operator))
	}
}

func checkFunc(fn *builtin.Function, name string, operator string) {
	if len(fn.Types) == 0 {
		panic(fmt.Errorf("function %s for %s operator misses types", name, operator))
	}
	for _, t := range fn.Types {
		if t.NumIn() != 2 || t.NumOut() != 1 {
			panic(fmt.Errorf("function %s for %s operator does not have a correct signature", name, operator))
		}
	}
}
