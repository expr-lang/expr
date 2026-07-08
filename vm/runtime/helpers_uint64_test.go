package runtime_test

import (
	"math"
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/vm/runtime"
)

// TestEqual_Uint64Large tests runtime.Equal with large uint64 values
// that exceed int64 range. These values would produce incorrect results
// if erroneously cast to int for comparison.
func TestEqual_Uint64Large(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want bool
	}{
		// Same large uint64 values should be equal
		{"uint64 large == same", uint64(16141183638984196173), uint64(16141183638984196173), true},
		{"uint64 maxuint64 == same", uint64(math.MaxUint64), uint64(math.MaxUint64), true},
		{"uint64 maxint64+1 == same", uint64(math.MaxInt64) + 1, uint64(math.MaxInt64) + 1, true},

		// Different large uint64 values should not be equal
		{"uint64 large != different", uint64(16141183638984196173), uint64(16141183638984196174), false},
		{"uint64 maxuint64 != maxuint64-1", uint64(math.MaxUint64), uint64(math.MaxUint64) - 1, false},

		// uint64 vs int comparisons (cross-type)
		{"uint64 large != int 0", uint64(16141183638984196173), 0, false},
		{"uint64 large != int 1", uint64(16141183638984196173), 1, false},
		{"uint64(5) == int(5)", uint64(5), 5, true},
		{"uint64(5) == int64(5)", uint64(5), int64(5), true},

		// uint64 boundary with int64
		{"uint64(maxint64) == int64(maxint64)", uint64(math.MaxInt64), int64(math.MaxInt64), true},
		{"uint64(maxint64+1) != int64(maxint64)", uint64(math.MaxInt64) + 1, int64(math.MaxInt64), false},

		// uint64 vs float64
		{"uint64(5) == float64(5)", uint64(5), float64(5), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runtime.Equal(tt.a, tt.b)
			assert.Equal(t, tt.want, got, "Equal(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
			// Also test reverse order (commutativity)
			got = runtime.Equal(tt.b, tt.a)
			assert.Equal(t, tt.want, got, "Equal(%v, %v) = %v; want %v", tt.b, tt.a, got, tt.want)
		})
	}
}

// TestLess_Uint64Large tests runtime.Less with large uint64 values.
func TestLess_Uint64Large(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want bool
	}{
		// uint64 vs uint64
		{"uint64(maxint64+1) < uint64(maxint64+2)", uint64(math.MaxInt64) + 1, uint64(math.MaxInt64) + 2, true},
		{"uint64(maxint64+2) < uint64(maxint64+1)", uint64(math.MaxInt64) + 2, uint64(math.MaxInt64) + 1, false},
		{"uint64(maxuint64-1) < uint64(maxuint64)", uint64(math.MaxUint64) - 1, uint64(math.MaxUint64), true},
		{"uint64(same) not less", uint64(16141183638984196173), uint64(16141183638984196173), false},

		// uint64 vs int (cross-type)
		{"uint64(maxint64+1) > int(0)", uint64(math.MaxInt64) + 1, 0, false}, // a < b is false
		{"int(0) < uint64(maxint64+1)", 0, uint64(math.MaxInt64) + 1, true},

		// Small values should work normally
		{"uint64(1) < uint64(2)", uint64(1), uint64(2), true},
		{"uint64(2) < uint64(1)", uint64(2), uint64(1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runtime.Less(tt.a, tt.b)
			assert.Equal(t, tt.want, got, "Less(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		})
	}
}

// TestMore_Uint64Large tests runtime.More with large uint64 values.
func TestMore_Uint64Large(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want bool
	}{
		// uint64 vs uint64
		{"uint64(maxint64+2) > uint64(maxint64+1)", uint64(math.MaxInt64) + 2, uint64(math.MaxInt64) + 1, true},
		{"uint64(maxint64+1) > uint64(maxint64+2)", uint64(math.MaxInt64) + 1, uint64(math.MaxInt64) + 2, false},
		{"uint64(maxuint64) > uint64(maxuint64-1)", uint64(math.MaxUint64), uint64(math.MaxUint64) - 1, true},
		{"uint64(same) not more", uint64(16141183638984196173), uint64(16141183638984196173), false},

		// uint64 vs int (cross-type)
		{"uint64(maxint64+1) > int(0)", uint64(math.MaxInt64) + 1, 0, true},
		{"int(0) > uint64(maxint64+1)", 0, uint64(math.MaxInt64) + 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runtime.More(tt.a, tt.b)
			assert.Equal(t, tt.want, got, "More(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		})
	}
}

// TestLessOrEqual_Uint64Large tests runtime.LessOrEqual with large uint64 values.
func TestLessOrEqual_Uint64Large(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want bool
	}{
		{"uint64 same <= same", uint64(16141183638984196173), uint64(16141183638984196173), true},
		{"uint64 less <= more", uint64(math.MaxInt64) + 1, uint64(math.MaxInt64) + 2, true},
		{"uint64 more <= less", uint64(math.MaxInt64) + 2, uint64(math.MaxInt64) + 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runtime.LessOrEqual(tt.a, tt.b)
			assert.Equal(t, tt.want, got, "LessOrEqual(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		})
	}
}

// TestMoreOrEqual_Uint64Large tests runtime.MoreOrEqual with large uint64 values.
func TestMoreOrEqual_Uint64Large(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want bool
	}{
		{"uint64 same >= same", uint64(16141183638984196173), uint64(16141183638984196173), true},
		{"uint64 more >= less", uint64(math.MaxInt64) + 2, uint64(math.MaxInt64) + 1, true},
		{"uint64 less >= more", uint64(math.MaxInt64) + 1, uint64(math.MaxInt64) + 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runtime.MoreOrEqual(tt.a, tt.b)
			assert.Equal(t, tt.want, got, "MoreOrEqual(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		})
	}
}
