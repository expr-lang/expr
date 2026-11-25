// Package budget contains helpers for memory budget accounting.
package budget

import "errors"

// ErrExceeded is returned when an operation would exceed the memory budget.
var ErrExceeded = errors.New("memory budget exceeded")

// StringGrowth returns the memory cost of a string of length n produced from
// the input value in.
//
// Strings are immutable, so the inputs of an operation already exist in memory:
// they are either part of the environment or were already accounted for. The
// cost is therefore the number of bytes by which the result exceeds its input.
// Scalar inputs (numbers, bools, etc.) produce strings of bounded size and cost
// nothing. Any other input (arrays, maps, structs) may reference the same
// string many times, so the whole result is charged.
func StringGrowth(n int, in any) uint {
	switch v := in.(type) {
	case string:
		if n > len(v) {
			return uint(n - len(v))
		}
		return 0
	case nil, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return 0
	}
	return uint(n)
}

// Concat returns the memory cost of concatenating strings of lengths a and b.
// The result of a + b grows the largest of them by the size of the smallest.
func Concat(a, b int) uint {
	if a < b {
		return uint(a)
	}
	return uint(b)
}
