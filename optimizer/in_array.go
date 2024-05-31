package optimizer

import (
	"reflect"

	. "github.com/expr-lang/expr/ast"
)

type inArray struct{}

func (*inArray) Visit(node *Node) {
	switch n := (*node).(type) {
	case *CompareNode:
		for i := 0; i < len(n.Operators); i++ {
			op := n.Operators[i]
			negate := op == "not"
			if negate {
				i++
				op = n.Operators[i]
			}

			if op == "in" {
				comparatorIdx := i
				if negate {
					comparatorIdx = i - 1
				}
				if array, ok := n.Comparators[comparatorIdx].(*ArrayNode); ok && len(array.Nodes) > 0 {
					var lType reflect.Type
					if comparatorIdx == 0 {
						lType = n.Left.Type()
					} else {
						lType = n.Comparators[comparatorIdx-1].Type()
					}
					if lType == nil || lType.Kind() != reflect.Int {
						goto string
					}

					if !allIntegerNodes(array.Nodes) {
						goto string
					}
					{
						value := make(map[int]struct{})
						for _, a := range array.Nodes {
							value[a.(*IntegerNode).Value] = struct{}{}
						}
						n.Comparators[comparatorIdx] = &ConstantNode{Value: value}
					}

				string:
					if !allStringNodes(array.Nodes) {
						continue
					}
					{
						value := make(map[string]struct{})
						for _, a := range array.Nodes {
							value[a.(*StringNode).Value] = struct{}{}
						}
						n.Comparators[comparatorIdx] = &ConstantNode{Value: value}
					}

				}
			}
		}
	}
}

func allIntegerNodes(nodes []Node) bool {
	for _, n := range nodes {
		if _, ok := n.(*IntegerNode); !ok {
			return false
		}
	}
	return true
}

func allStringNodes(nodes []Node) bool {
	for _, n := range nodes {
		if _, ok := n.(*StringNode); !ok {
			return false
		}
	}
	return true
}
