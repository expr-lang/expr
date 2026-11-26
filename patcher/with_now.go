package patcher

import (
	"time"

	"github.com/expr-lang/expr/ast"
)

// WithNow replaces now() calls with a constant time. This is useful to get
// deterministic results when testing expressions that depend on the current
// time. The replacement happens while patching the AST, so the compiled
// program always returns the same time from now().
type WithNow struct {
	Now time.Time
}

func (p WithNow) Visit(node *ast.Node) {
	btin, ok := (*node).(*ast.BuiltinNode)
	if !ok || btin.Name != "now" {
		return
	}
	now := p.Now
	// Preserve the location injected by WithTimezone, if any, so that
	// Timezone and MockNow can be combined.
	for _, arg := range btin.Arguments {
		c, ok := arg.(*ast.ConstantNode)
		if !ok {
			continue
		}
		if loc, ok := c.Value.(*time.Location); ok {
			now = now.In(loc)
		}
	}
	ast.Patch(node, &ast.ConstantNode{Value: now})
}
