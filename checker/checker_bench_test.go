package checker_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/checker"
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
	expr.Function("func", f, f)
	expr.ConstExpr("func")

	for _, c := range cases {
		b.Run("name="+c.name, func(b *testing.B) {
			tree, err := parser.ParseWithConfig(c.input, config)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err = checker.Check(tree, config)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}

}

type recursive struct {
	Inner *recursive `expr:"a"`
}
