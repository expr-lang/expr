/*
Package expr is an engine that can evaluate expressions.

	// Evaluate expression on data.
	result, err := expr.Eval("expression", data)

	// Or precompile expression to ast first.
	node, err := expr.Parse("expression")

	// And run later.
	result, err := expr.Run(node, data)

*/
package expr
