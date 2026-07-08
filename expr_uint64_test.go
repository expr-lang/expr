package expr_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

// TestUint64Literal_Parse tests that large integer literals (exceeding int64 range)
// can be parsed correctly as uint64 values.
func TestUint64Literal_Parse(t *testing.T) {
	tests := []struct {
		code string
		want uint64
	}{
		// The original reported issue
		{`16141183638984196173`, 16141183638984196173},
		// MaxInt64 + 1 (first value that overflows int64)
		{`9223372036854775808`, uint64(math.MaxInt64) + 1},
		// MaxUint64
		{`18446744073709551615`, math.MaxUint64},
		// Another large value
		{`9228157111460438039`, 9228157111460438039},
		// Hex literal exceeding int64
		{`0xFFFFFFFFFFFFFFFF`, math.MaxUint64},
		{`0x8000000000000000`, uint64(math.MaxInt64) + 1},
		// Binary literal exceeding int64
		{`0b1000000000000000000000000000000000000000000000000000000000000000`, uint64(math.MaxInt64) + 1},
		// Octal literal exceeding int64
		{`0o1000000000000000000000`, 9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code)
			require.NoError(t, err, "failed to compile: %s", tt.code)

			got, err := expr.Run(program, nil)
			require.NoError(t, err, "failed to run: %s", tt.code)
			assert.Equal(t, tt.want, got, "unexpected result for: %s", tt.code)
		})
	}
}

// TestUint64Literal_Boundary tests boundary values around int64 max.
func TestUint64Literal_Boundary(t *testing.T) {
	tests := []struct {
		code string
		want any
	}{
		// MaxInt64 should still parse as int (existing behavior)
		{`9223372036854775807`, int(math.MaxInt64)},
		// MaxInt64 + 1 should parse as uint64
		{`9223372036854775808`, uint64(math.MaxInt64) + 1},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestUint64Literal_Comparison tests comparison operators with large uint64 values.
func TestUint64Literal_Comparison(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		// Equal
		{`16141183638984196173 == 16141183638984196173`, true},
		{`9223372036854775808 == 9223372036854775808`, true},
		{`18446744073709551615 == 18446744073709551615`, true},
		{`9223372036854775808 == 9223372036854775809`, false},

		// Not equal
		{`16141183638984196173 != 16141183638984196174`, true},
		{`16141183638984196173 != 16141183638984196173`, false},

		// Less than
		{`9223372036854775808 < 9223372036854775809`, true},
		{`9223372036854775809 < 9223372036854775808`, false},
		{`9223372036854775808 < 9223372036854775808`, false},

		// Greater than
		{`9223372036854775809 > 9223372036854775808`, true},
		{`9223372036854775808 > 9223372036854775809`, false},
		{`9223372036854775808 > 9223372036854775808`, false},

		// Less than or equal
		{`9223372036854775808 <= 9223372036854775808`, true},
		{`9223372036854775808 <= 9223372036854775809`, true},
		{`9223372036854775809 <= 9223372036854775808`, false},

		// Greater than or equal
		{`9223372036854775808 >= 9223372036854775808`, true},
		{`9223372036854775809 >= 9223372036854775808`, true},
		{`9223372036854775808 >= 9223372036854775809`, false},

		// Mixed uint64 literal and int literal comparisons
		{`9223372036854775808 > 0`, true},
		{`9223372036854775808 > 100`, true},
		{`0 < 9223372036854775808`, true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, nil)
			require.NoError(t, err, "eval error: %s", tt.code)
			assert.Equal(t, tt.want, got, "unexpected result for: %s", tt.code)
		})
	}
}

// TestUint64Literal_ComparisonWithEnv tests comparison with uint64 values from environment.
func TestUint64Literal_ComparisonWithEnv(t *testing.T) {
	env := map[string]any{
		"big":  uint64(16141183638984196173),
		"big2": uint64(16141183638984196173),
		"big3": uint64(16141183638984196174),
		"max":  uint64(math.MaxUint64),
	}

	tests := []struct {
		code string
		want any
	}{
		// Comparing env uint64 with env uint64
		{`big == big2`, true},
		{`big == big3`, false},
		{`big < big3`, true},
		{`big3 > big`, true},
		{`big != big3`, true},

		// Comparing env uint64 with literal
		{`big == 16141183638984196173`, true},
		{`big != 16141183638984196174`, true},
		{`big < 16141183638984196174`, true},

		// MaxUint64 comparisons
		{`max == 18446744073709551615`, true},
		{`max > 9223372036854775808`, true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, env)
			require.NoError(t, err, "eval error: %s", tt.code)
			assert.Equal(t, tt.want, got, "unexpected result for: %s", tt.code)
		})
	}
}

// TestUint64Literal_Arithmetic tests arithmetic operations with large uint64 values.
func TestUint64Literal_Arithmetic(t *testing.T) {
	env := map[string]any{
		"a": uint64(9223372036854775808), // MaxInt64 + 1
		"b": uint64(1),
	}

	tests := []struct {
		code string
		want any
	}{
		// Addition
		{`9223372036854775808 + 1`, uint64(9223372036854775809)},
		{`a + b`, uint64(9223372036854775809)},

		// Subtraction
		{`9223372036854775809 - 1`, uint64(9223372036854775808)},
		{`9223372036854775808 - 9223372036854775808`, uint64(0)},

		// Multiplication
		{`9223372036854775808 * 2`, uint64(0)}, // wraps around on overflow
		{`a * 2`, uint64(0)},                   // 0x8000000000000000 * 2 overflows to 0

		// Modulo
		{`9223372036854775809 % 2`, uint64(1)},
		{`18446744073709551615 % 2`, uint64(1)},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, env)
			require.NoError(t, err, "eval error: %s", tt.code)
			assert.Equal(t, tt.want, got, "unexpected result for: %s", tt.code)
		})
	}
}

// TestUint64Literal_Overflow tests that values exceeding uint64 max still produce errors.
func TestUint64Literal_Overflow(t *testing.T) {
	overflowCases := []string{
		// MaxUint64 + 1
		`18446744073709551616`,
		// Much larger than uint64
		`99999999999999999999`,
	}

	for _, code := range overflowCases {
		t.Run(code, func(t *testing.T) {
			_, err := expr.Compile(code)
			require.Error(t, err, "expected error for overflow: %s", code)
		})
	}
}

// TestUint64Literal_InExpression tests uint64 literals used in more complex expressions.
func TestUint64Literal_InExpression(t *testing.T) {
	env := map[string]any{
		"values": []uint64{16141183638984196173, 9223372036854775808, 18446744073709551615},
		"target": uint64(16141183638984196173),
	}

	tests := []struct {
		code string
		want any
	}{
		// uint64 in ternary
		{`true ? 9223372036854775808 : 0`, uint64(9223372036854775808)},

		// uint64 in let expression
		{`let x = 9223372036854775808; x == 9223372036854775808`, true},

		// uint64 in array contains (from env)
		{`target in values`, true},

		// Negation should fail or be handled correctly
		// (can't negate uint64 > MaxInt64)
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, env)
			require.NoError(t, err, "eval error: %s", tt.code)
			assert.Equal(t, tt.want, got, "unexpected result for: %s", tt.code)
		})
	}
}

// TestUint64_RuntimeEqual tests the runtime Equal function with large uint64 values
// that would overflow if incorrectly cast to int.
func TestUint64_RuntimeEqual(t *testing.T) {
	// These values are > MaxInt64, they would produce wrong results
	// if cast to int for comparison.
	env := map[string]any{
		"a": uint64(16141183638984196173),
		"b": uint64(16141183638984196173),
		"c": uint64(16141183638984196174),
	}

	tests := []struct {
		code string
		want bool
	}{
		{`a == b`, true},
		{`a == c`, false},
		{`a != c`, true},
		{`a != b`, false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestUint64_WithUnderscore tests uint64 literals with underscore separators.
func TestUint64_WithUnderscore(t *testing.T) {
	tests := []struct {
		code string
		want uint64
	}{
		{`16_141_183_638_984_196_173`, 16141183638984196173},
		{`9_223_372_036_854_775_808`, uint64(math.MaxInt64) + 1},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestUint64_String tests that uint64 values can be used with string formatting.
func TestUint64Literal_String(t *testing.T) {
	env := map[string]any{
		"val":     uint64(16141183638984196173),
		"sprintf": fmt.Sprintf,
	}

	got, err := expr.Eval(`sprintf("%d", val)`, env)
	require.NoError(t, err)
	assert.Equal(t, "16141183638984196173", got)
}
