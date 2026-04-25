package expr

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr/checker/nature"
	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm/runtime"
)

// This file installs the method-dispatch hooks that the shared packages
// (vm/runtime, checker/nature, builtin) consult to perform reflective method
// lookups. The hooks are installed only when the parent expr package is
// imported.
//
// All four reflect.* method-resolution call sites that the linker treats as
// REFLECTMETHOD live exclusively in this file (or in functions reachable
// only from this file's hooks):
//
//   - reflect.Value.MethodByName  in fetchMethodByName
//   - reflect.Value.Method        in fetchMethodIndexed
//   - reflect.Type.Method         in nature.LookupMethod (transitively)
//   - reflect.Type.MethodByName   not used

func init() {
	runtime.MethodByNameHook = fetchMethodByName
	runtime.MethodIndexedHook = fetchMethodIndexed
	nature.MethodByNameHook = nature.LookupMethod
}

func fetchMethodByName(v reflect.Value, name string) (any, bool) {
	method := v.MethodByName(name)
	if method.IsValid() {
		return method.Interface(), true
	}
	return nil, false
}

func fetchMethodIndexed(v reflect.Value, index int) (any, bool) {
	method := v.Method(index)
	if method.IsValid() {
		return method.Interface(), true
	}
	return nil, false
}

// Eval parses, compiles and runs given input.
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
