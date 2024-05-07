package issues584_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"

	"github.com/expr-lang/expr"
)

type Env struct{}

type Program struct {
}

func (p *Program) Foo() Value {
	return func(e *Env) float64 {
		return 5
	}
}

func (p *Program) Bar() Value {
	return func(e *Env) float64 {
		return 100
	}
}

func (p *Program) AndCondition(a, b Condition) Conditions {
	return Conditions{a, b}
}

func (p *Program) AndConditions(a Conditions, b Condition) Conditions {
	return append(a, b)
}

func (p *Program) ValueGreaterThan_float(v Value, i float64) Condition {
	return func(e *Env) bool {
		realized := v(e)
		return realized > i
	}
}

func (p *Program) ValueLessThan_float(v Value, i float64) Condition {
	return func(e *Env) bool {
		realized := v(e)
		return realized < i
	}
}

type Condition func(e *Env) bool
type Conditions []Condition

type Value func(e *Env) float64

func TestIssue584(t *testing.T) {
	code := `Foo() > 1.5 and Bar() < 200.0`

	p := &Program{}

	opt := []expr.Option{
		expr.Env(p),
		expr.Operator("and", "AndCondition", "AndConditions"),
		expr.Operator(">", "ValueGreaterThan_float"),
		expr.Operator("<", "ValueLessThan_float"),
	}

	program, err := expr.Compile(code, opt...)
	assert.Nil(t, err)

	state, err := expr.Run(program, p)
	assert.Nil(t, err)
	assert.NotNil(t, state)
}
