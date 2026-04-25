package expr_test

import (
	"testing"

	"github.com/expr-lang/expr/static"
)

type env struct {
	X, Y int
}

// Sum should not be callable.
func (e env) Sum() int { return e.X + e.Y }

func TestCompileRun(t *testing.T) {
	prog, err := expr.Compile(`X + Y`, expr.Env(env{}))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	out, err := expr.Run(prog, env{X: 3, Y: 4})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if out != 7 {
		t.Fatalf("X+Y = %v, want 7", out)
	}
}

func TestEnvMethodNotSurfaced(t *testing.T) {
	_, err := expr.Compile(`Sum()`, expr.Env(env{}))
	if err == nil {
		t.Fatal("expected Compile to fail for env-defined method, got nil")
	}
}

func TestFunctionOption(t *testing.T) {
	prog, err := expr.Compile(`Add(X, Y)`,
		expr.Env(env{}),
		expr.Function("Add",
			func(p ...any) (any, error) { return p[0].(int) + p[1].(int), nil },
			new(func(int, int) int)))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	out, err := expr.Run(prog, env{X: 10, Y: 20})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if out != 30 {
		t.Fatalf("Add(X,Y) = %v, want 30", out)
	}
}
