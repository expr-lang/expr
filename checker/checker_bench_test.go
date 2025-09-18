package checker_test

import (
	"runtime"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/checker"
	"github.com/expr-lang/expr/checker/nature"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/parser"
)

func BenchmarkChecker(b *testing.B) {
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
		{"unary and binary operations", `
!b && !b || !b == !b && !b != !b || 1 < 1.0 && 0.1 > 1 || 0 <= 1.0 && 0.1 >= 1 &&
!b && !b || !b == !b && !b != !b || 1 < 1.0 && 0.1 > 1 || 0 <= 1.0 && 0.1 >= 1 &&
!b && !b || !b == !b && !b != !b || 1 < 1.0 && 0.1 > 1 || 0 <= 1.0 && 0.1 >= 1 &&
!b && !b || !b == !b && !b != !b || 1 < 1.0 && 0.1 > 1 || 0 <= 1.0 && 0.1 >= 1 &&
!b && !b || !b == !b && !b != !b || 1 < 1.0 && 0.1 > 1 || 0 <= 1.0 && 0.1 >= 1 &&
!b && !b || !b == !b && !b != !b || 1 < 1.0 && 0.1 > 1 || 0 <= 1.0 && 0.1 >= 1
		`},
		{"deep struct access", `
a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.
a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.
a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.
a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a
		`},
	}

	f := func(params ...any) (any, error) { return nil, nil }
	env := map[string]any{
		"a":    new(recursive),
		"b":    true,
		"func": f,
	}
	config := conf.New(env)
	expr.Function("func", f, f)(config)
	expr.ConstExpr("func")(config)

	for _, c := range cases {
		batchSize := 100_000
		if batchSize > b.N {
			batchSize = b.N
		}
		trees := make([]*parser.Tree, 0, batchSize)
		for i := 0; i < batchSize; i++ {
			tree, err := parser.ParseWithConfig(c.input, config)
			if err != nil {
				b.Fatal(err)
			}
			trees = append(trees, tree)
		}
		runtime.GC() // try to cleanup the mess from the initialization

		b.Run("name="+c.name, func(b *testing.B) {
			var err error
			for i := 0; i < b.N; i++ {
				j := i
				if j < 0 || j >= len(trees) {
					b.StopTimer()
					invalidateTrees(trees...)
					j = 0
					b.StartTimer()
				}

				_, err = checker.Check(trees[j], config)
			}
			b.StopTimer()
			if err != nil {
				b.Fatal(err)
			}
		})
	}
}

type visitorFunc func(*ast.Node)

func (f visitorFunc) Visit(node *ast.Node) { f(node) }

func invalidateTrees(trees ...*parser.Tree) {
	for _, tree := range trees {
		ast.Walk(&tree.Node, visitorFunc(func(node *ast.Node) {
			(*node).SetNature(nature.Nature{})
		}))
	}
}

type recursive struct {
	Inner *recursive `expr:"a"`
}
