package examples

import (
	"testing"
	"time"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/require"
)

func TestExamples_dates(t *testing.T) {
	code := `
		Now() > Date("2020-01-01") &&
		Now() - CreatedAt > Duration("24h")
	`

	options := []expr.Option{
		expr.Env(Env{}),

		// Operators override for date comprising.
		expr.Operator("==", "Equal"),
		expr.Operator("<", "Before"),
		expr.Operator("<=", "BeforeOrEqual"),
		expr.Operator(">", "After"),
		expr.Operator(">=", "AfterOrEqual"),

		// Time and duration manipulation.
		expr.Operator("+", "Add"),
		expr.Operator("-", "Sub"),

		// Operators override for duration comprising.
		expr.Operator("==", "EqualDuration"),
		expr.Operator("<", "BeforeDuration"),
		expr.Operator("<=", "BeforeOrEqualDuration"),
		expr.Operator(">", "AfterDuration"),
		expr.Operator(">=", "AfterOrEqualDuration"),
	}

	program, err := expr.Compile(code, options...)
	require.NoError(t, err)

	env := Env{
		CreatedAt: Env{}.Date("2018-07-14"), // first commit date
	}

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

type Env struct {
	datetime
	CreatedAt time.Time
}

type datetime struct{}

func (datetime) Date(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}
func (datetime) Duration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}
func (datetime) Now() time.Time                                { return time.Now() }
func (datetime) Equal(a, b time.Time) bool                     { return a.Equal(b) }
func (datetime) Before(a, b time.Time) bool                    { return a.Before(b) }
func (datetime) BeforeOrEqual(a, b time.Time) bool             { return a.Before(b) || a.Equal(b) }
func (datetime) After(a, b time.Time) bool                     { return a.After(b) }
func (datetime) AfterOrEqual(a, b time.Time) bool              { return a.After(b) || a.Equal(b) }
func (datetime) Add(a time.Time, b time.Duration) time.Time    { return a.Add(b) }
func (datetime) Sub(a, b time.Time) time.Duration              { return a.Sub(b) }
func (datetime) EqualDuration(a, b time.Duration) bool         { return a == b }
func (datetime) BeforeDuration(a, b time.Duration) bool        { return a < b }
func (datetime) BeforeOrEqualDuration(a, b time.Duration) bool { return a <= b }
func (datetime) AfterDuration(a, b time.Duration) bool         { return a > b }
func (datetime) AfterOrEqualDuration(a, b time.Duration) bool  { return a >= b }
