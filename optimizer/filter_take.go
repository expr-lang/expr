package optimizer

import (
	. "github.com/expr-lang/expr/ast"
)

type filterTake struct{}

func (*filterTake) Visit(node *Node) {
	take, ok := (*node).(*BuiltinNode)
	if !ok || take.Name != "take" || len(take.Arguments) != 2 {
		return
	}
	filter, ok := take.Arguments[0].(*BuiltinNode)
	if !ok || filter.Name != "filter" || len(filter.Arguments) != 2 {
		return
	}
	limit, ok := take.Arguments[1].(*IntegerNode)
	if !ok || limit.Value <= 0 {
		return
	}

	threshold := limit.Value
	if filter.Threshold != nil && *filter.Threshold < threshold {
		threshold = *filter.Threshold
	}
	patchCopyType(node, &BuiltinNode{
		Name:      "filter",
		Arguments: filter.Arguments,
		Map:       filter.Map,
		Threshold: &threshold,
	})
}
