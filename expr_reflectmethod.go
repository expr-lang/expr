//go:build !expr_noreflectmethod

package expr

import (
	"fmt"

	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/parser"
)

// Eval parses, compiles and runs given input.
//
// Eval is excluded from the build under the expr_noreflectmethod tag because
// it relies on runtime dispatch on the env, which requires reflect-based
// method resolution.
func Eval(input string, env any) (any, error) {
	if _, ok := env.(Option); ok {
		return nil, fmt.Errorf("misused expr.Eval: second argument (env) should be passed without expr.Env")
	}

	tree, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	program, err := compiler.Compile(tree, nil)
	if err != nil {
		return nil, err
	}

	output, err := Run(program, env)
	if err != nil {
		return nil, err
	}

	return output, nil
}
