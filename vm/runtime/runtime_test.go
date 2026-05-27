package runtime_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"

	"github.com/expr-lang/expr/vm/runtime"
)

// TestIn_TypedSlices exercises the typed-slice fast paths in runtime.In to
// guarantee they preserve the semantics of the reflect-based fallback.
func TestIn_TypedSlices(t *testing.T) {
	cases := []struct {
		name   string
		needle any
		array  any
		want   bool
	}{
		// []string fast path
		{"string in []string (hit)", "b", []string{"a", "b", "c"}, true},
		{"string in []string (miss)", "z", []string{"a", "b", "c"}, false},
		{"string in empty []string", "x", []string{}, false},

		// []float64 fast path
		{"float64 in []float64 (hit)", 2.5, []float64{1.0, 2.5, 3.0}, true},
		{"float64 in []float64 (miss)", 9.9, []float64{1.0, 2.5, 3.0}, false},

		// []int64 fast path
		{"int64 in []int64 (hit)", int64(2), []int64{1, 2, 3}, true},
		{"int64 in []int64 (miss)", int64(9), []int64{1, 2, 3}, false},

		// []int fast path
		{"int in []int (hit)", 2, []int{1, 2, 3}, true},
		{"int in []int (miss)", 9, []int{1, 2, 3}, false},

		// []bool fast path
		{"true in []bool (hit)", true, []bool{false, true, false}, true},
		{"false in []bool (hit)", false, []bool{true, true, false}, true},
		{"true in []bool (miss all-false)", true, []bool{false, false}, false},

		// Type-mismatched needles must fall through to the reflect path so
		// Equal()'s cross-type semantics are preserved. e.g. an int needle
		// against a []float64 should still match via numeric promotion.
		{"int needle in []float64 (promoted hit)", 2, []float64{1.0, 2.0, 3.0}, true},
		{"int needle in []float64 (promoted miss)", 9, []float64{1.0, 2.0, 3.0}, false},
		{"int needle in []int64 (promoted hit)", 2, []int64{1, 2, 3}, true},

		// []any keeps using the reflect path (unchanged).
		{"string in []any (hit)", "b", []any{"a", "b", "c"}, true},
		{"int in []any (hit)", 2, []any{1, 2, 3}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, runtime.In(tc.needle, tc.array))
		})
	}
}

// TestIn_NilArray ensures the early-return for a nil right-hand side is
// preserved (it lives above the typed-slice fast paths).
func TestIn_NilArray(t *testing.T) {
	assert.False(t, runtime.In("x", nil))
}
