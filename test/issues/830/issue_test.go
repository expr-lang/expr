package issues

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue830(t *testing.T) {
	program, err := expr.Compile("varNotExist", expr.AllowUndefinedVariables(), expr.AsBool())
	require.NoError(t, err)

	output, err := expr.Run(program, map[string]interface{}{})
	require.NoError(t, err)

	// The user expects output to be false (bool), but gets nil.
	assert.Equal(t, false, output)
}
