package optimizer

import (
	"github.com/antonmedv/expr/ast"
)

type optimizer struct {
}

func (*optimizer) Enter(node *ast.Node) {}

func (*optimizer) Exit(node *ast.Node) {
	switch n := (*node).(type) {
	case *ast.ArrayNode:
		if len(n.Nodes) > 0 {

			for _, a := range n.Nodes {
				if _, ok := a.(*ast.IntegerNode); !ok {
					goto string
				}
			}
			{
				value := make([]int, len(n.Nodes))
				for i, a := range n.Nodes {
					value[i] = a.(*ast.IntegerNode).Value
				}
				*node = &ast.ConstantNode{
					Value: value,
				}
			}

		string:
			for _, a := range n.Nodes {
				if _, ok := a.(*ast.StringNode); !ok {
					return
				}
			}
			{
				value := make([]string, len(n.Nodes))
				for i, a := range n.Nodes {
					value[i] = a.(*ast.StringNode).Value
				}
				*node = &ast.ConstantNode{
					Value: value,
				}
			}

		}
	}
}

func Optimize(node *ast.Node) {
	o := &optimizer{}
	ast.Walk(node, o)
}
