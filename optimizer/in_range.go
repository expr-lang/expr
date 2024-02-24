package optimizer

import (
	"reflect"

	. "github.com/expr-lang/expr/ast"
)

type inRange struct{}

func (*inRange) Visit(node *Node) {
	switch n := (*node).(type) {
	case *BinaryNode:
		if n.Operator == "in" {
			t := n.Left.Type()
			if t == nil {
				return
			}
			if t.Kind() != reflect.Int {
				return
			}
			if rangeOp, ok := n.Right.(*BinaryNode); ok && rangeOp.Operator == ".." {
				if from, ok := rangeOp.Left.(*IntegerNode); ok {
					if to, ok := rangeOp.Right.(*IntegerNode); ok {
						Patch(node, &ComparisonNode{
							Left:        from,
							Comparators: []Node{n.Left, to},
							Operators:   []string{"<=", "<="},
						})
					}
				}
			}
		}
	}
}
