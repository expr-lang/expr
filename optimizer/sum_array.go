package optimizer

import (
	. "github.com/expr-lang/expr/ast"
)

type sumArray struct{}

func (*sumArray) Visit(node *Node) {
	if sumBuiltin, ok := (*node).(*BuiltinNode); ok &&
		sumBuiltin.Name == "sum" &&
		len(sumBuiltin.Arguments) == 1 {
		if array, ok := sumBuiltin.Arguments[0].(*ArrayNode); ok &&
			len(array.Nodes) >= 2 {
			patchCopyType(node, sumArrayFold(array.Nodes))
		}
	}
}

// sumArrayFold builds ((a + b) + c) + ..., adding the elements left to right
// like sum() does at runtime. Float addition is not associative, so any other
// grouping can give a different result than the unoptimized program.
func sumArrayFold(nodes []Node) Node {
	result := nodes[0]
	for _, n := range nodes[1:] {
		result = &BinaryNode{
			Operator: "+",
			Left:     result,
			Right:    n,
		}
	}
	return result
}
