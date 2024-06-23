package runtime_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"

	"github.com/expr-lang/expr/vm/runtime"
)

var tests = []struct {
	name string
	a, b any
	want bool
}{
	{"int == int", 42, 42, true},
	{"int != int", 42, 33, false},
	{"int == int8", 42, int8(42), true},
	{"int == int16", 42, int16(42), true},
	{"int == int32", 42, int32(42), true},
	{"int == int64", 42, int64(42), true},
	{"float == float", 42.0, 42.0, true},
	{"float != float", 42.0, 33.0, false},
	{"float == int", 42.0, 42, true},
	{"float != int", 42.0, 33, false},
	{"string == string", "foo", "foo", true},
	{"string != string", "foo", "bar", false},
	{"bool == bool", true, true, true},
	{"bool != bool", true, false, false},
	{"[]any == []int", []any{1, 2, 3}, []int{1, 2, 3}, true},
	{"[]any != []int", []any{1, 2, 3}, []int{1, 2, 99}, false},
	{"deep []any == []any", []any{[]int{1}, 2, []any{"3"}}, []any{[]any{1}, 2, []string{"3"}}, true},
	{"deep []any != []any", []any{[]int{1}, 2, []any{"3", "42"}}, []any{[]any{1}, 2, []string{"3"}}, false},
	{"map[string]any == map[string]any", map[string]any{"a": 1}, map[string]any{"a": 1}, true},
	{"map[string]any != map[string]any", map[string]any{"a": 1}, map[string]any{"a": 1, "b": 2}, false},
}

func TestEqual(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runtime.Equal(tt.a, tt.b)
			assert.Equal(t, tt.want, got, "Equal(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
			got = runtime.Equal(tt.b, tt.a)
			assert.Equal(t, tt.want, got, "Equal(%v, %v) = %v; want %v", tt.b, tt.a, got, tt.want)
		})
	}

}

func BenchmarkEqual(b *testing.B) {
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runtime.Equal(tt.a, tt.b)
			}
		})
	}
}
