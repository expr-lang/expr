package issue_test

import (
	"fmt"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue817_1(t *testing.T) {
	out, err := expr.Eval(
		`sprintf("result: %v %v", 1, nil)`,
		map[string]any{
			"sprintf": fmt.Sprintf,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "result: 1 <nil>", out)
}

func TestIssue817_2(t *testing.T) {
	out, err := expr.Eval(
		`thing(nil)`,
		map[string]any{
			"thing": func(arg ...any) string {
				return fmt.Sprintf("result: (%T) %v", arg[0], arg[0])
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "result: (<nil>) <nil>", out)
}
