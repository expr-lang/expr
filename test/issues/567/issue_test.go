package expr_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue567(t *testing.T) {
	program, err := expr.Compile("concat(1..2, 3..4)")
	require.NoError(t, err)

	var buf bytes.Buffer
	program.DisassembleWriter(&buf)
	output := buf.String()

	// Check if "concat" is mentioned in the output
	require.True(t, strings.Contains(output, "concat"), "expected 'concat' in disassembly output")

	// It should appear as the called builtin
	require.True(t, strings.Contains(output, "OpCallBuiltinSafe\t"), "expected 'OpCallBuiltinSafe' in disassembly output")
	require.True(t, strings.Contains(output, "concat(2 args)"), "expected 'concat(2 args)' in disassembly output")
}
