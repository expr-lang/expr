// Package expr (imported as github.com/expr-lang/expr/static) is a
// dead-code-elimination-friendly entry point for expr.
//
// It exposes the same API as the parent expr, without the `Eval()` function. It
// does not trigger the Go linker's REFLECTMETHOD analysis and methods of every
// reachable type can be eliminated by dead-code elimination. Calling a method
// defined on the Env type does NOT work. Compilation fails with "unknown name
// <method>". To expose a function, use the Function option.
package expr

//go:generate cp ../expr.go expr.go
