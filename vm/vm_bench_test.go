package vm_test

import (
	"runtime"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/checker"
	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/vm"
)

func BenchmarkVM(b *testing.B) {
	cases := []struct {
		name, input string
	}{
		{"function calls", `
func(
	func(
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
	),
	func(
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
	),
	func(
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
		func(func(a, 'a', 1, nil), func(a, 'a', 1, nil), func(a, 'a', 1, nil)),
	)
)
		`},
	}

	a := new(recursive)
	for i, b := 0, a; i < 40*4; i++ {
		b.Inner = new(recursive)
		b = b.Inner
	}

	f := func(params ...any) (any, error) { return nil, nil }
	env := map[string]any{
		"a":    a,
		"b":    true,
		"func": f,
	}
	config := conf.New(env)
	expr.Function("func", f, f)(config)
	config.Check()

	for _, c := range cases {
		tree, err := checker.ParseCheck(c.input, config)
		if err != nil {
			b.Fatal(c.input, "parse and check", err)
		}
		prog, err := compiler.Compile(tree, config)
		if err != nil {
			b.Fatal(c.input, "compile", err)
		}
		//b.Logf("disassembled:\n%s", prog.Disassemble())
		//b.FailNow()
		runtime.GC()

		var vm vm.VM
		b.Run("name="+c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err = vm.Run(prog, env)
			}
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

type recursive struct {
	Inner *recursive `expr:"a"`
}
