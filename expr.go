package expr

import (
	"gopkg.in/antonmedv/expr.v2/checker"
	"gopkg.in/antonmedv/expr.v2/compiler"
	"gopkg.in/antonmedv/expr.v2/parser"
	"gopkg.in/antonmedv/expr.v2/vm"
	"reflect"
)

// Eval parses, compiles and runs given input.
func Eval(input string, env interface{}) (interface{}, error) {
	tree, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	program, err := compiler.Compile(tree)
	if err != nil {
		return nil, err
	}

	output, err := vm.Run(program, env)
	if err != nil {
		return nil, err
	}

	return output, nil
}

type config struct {
	mapEnv bool
	types  checker.TypesTable
}

// OptionFn for configuring expr.
type OptionFn func(c *config)

// Env specifies expected input of env for type checks.
// If struct is passed, all fields will be treated as variables,
// as well as all fields of embedded structs and struct itself.
// If map is passed, all items will be treated as variables.
// Methods defined on this type will be available as functions.
func Env(i interface{}) OptionFn {
	return func(c *config) {
		if _, ok := i.(map[string]interface{}); ok {
			c.mapEnv = true
		}
		c.types = checker.CreateTypesTable(i)
	}
}

// Compile parses and compiles given input expression to bytecode program.
func Compile(input string, ops ...OptionFn) (*vm.Program, reflect.Type, error) {
	c := &config{}

	for _, op := range ops {
		op(c)
	}

	tree, err := parser.Parse(input)
	if err != nil {
		return nil, nil, err
	}

	var t reflect.Type
	if c.types != nil {
		t, err = checker.Check(tree, c.types)
		if err != nil {
			return nil, nil, err
		}
	}

	compilerOps := make([]compiler.OptionFn, 0)
	if c.mapEnv {
		compilerOps = append(compilerOps, compiler.MapEnv())
	}

	program, err := compiler.Compile(tree, compilerOps...)
	if err != nil {
		return nil, nil, err
	}

	return program, t, nil
}

// Run evaluates given bytecode program.
func Run(program *vm.Program, env interface{}) (interface{}, error) {
	return vm.Run(program, env)
}
