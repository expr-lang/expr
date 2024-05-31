package optimizer

import (
	"reflect"

	. "github.com/expr-lang/expr/ast"
)

type inRange struct{}

func (*inRange) Visit(node *Node) {
	switch n := (*node).(type) {
	case *CompareNode:
		opSize := len(n.Operators)
		for i := 0; i < opSize; i++ {
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
				if rangeOp, ok := n.Comparators[comparatorIdx].(*BinaryNode); ok && rangeOp.Operator == ".." {
					if from, ok := rangeOp.Left.(*IntegerNode); ok {
						if to, ok := rangeOp.Right.(*IntegerNode); ok {
							var lNode Node
							if comparatorIdx == 0 {
								lNode = n.Left
							} else {
								lNode = n.Comparators[comparatorIdx-1]
							}
							if lType := lNode.Type(); lType != nil && lType.Kind() == reflect.Int {
								if comparatorIdx == 0 {
									if len(n.Comparators) == 1 {
										n.Operators = []string{"<=", "<="}
										n.Comparators = []Node{lNode, to}
										n.Left = from
										if negate {
											Patch(node, &UnaryNode{
												Operator: "not",
												Node:     n,
											})
										}
										return
									} else {
										n.Operators[comparatorIdx] = "&&"
										n.Left = &CompareNode{
											Left:        from,
											Operators:   []string{"<=", "<="},
											Comparators: []Node{lNode, to},
										}
										n.Comparators = n.Comparators[comparatorIdx+1:]
										if negate {
											n.Left = &UnaryNode{
												Operator: "not",
												Node:     n.Left,
											}
											n.Operators = append(n.Operators[:comparatorIdx+1], n.Operators[comparatorIdx+2:]...)
											opSize--
										}
									}
								} else {
									n.Operators[comparatorIdx] = "&&"
									var comparator Node = &CompareNode{
										Left:        from,
										Operators:   []string{"<=", "<="},
										Comparators: []Node{lNode, to},
									}
									if negate {
										comparator = &UnaryNode{
											Operator: "not",
											Node:     comparator,
										}
										n.Operators = append(n.Operators[:comparatorIdx+1], n.Operators[comparatorIdx+2:]...)
										opSize--
									}
									n.Comparators[comparatorIdx] = comparator
								}
							}
						}
					}
				}
			}
		}
	}
}
