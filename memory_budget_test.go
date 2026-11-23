package expr_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/vm"
)

// totalAlloc returns bytes allocated while running fn.
func totalAlloc(fn func()) uint64 {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	fn()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

func TestMemoryBudget_exceeded(t *testing.T) {
	long := strings.Repeat("A", 1000)
	tests := []struct {
		name string
		code string
		env  map[string]any
	}{
		{"concat doubling", `reduce(1..34, #acc + #acc, "a")`, nil},
		{"repeat", `repeat("` + long + `", 100000)`, nil},
		{"repeat of env string", `repeat(s, 1000000)`, map[string]any{"s": long}},
		{"replace", `replace(s, "A", s)`, map[string]any{"s": long + long}},
		{"replace doubling", `reduce(1..40, replace(#acc, "a", "aa"), "a")`, nil},
		{"join of aliased strings", `reduce(1..40, join([#acc, #acc]), "a")`, nil},
		{"join with glue", `join(split(s, ""), s)`, map[string]any{"s": long}},
		{"string of array", `reduce(1..40, string([#acc, #acc]), "a")`, nil},
		{"toJSON of string", `reduce(1..40, toJSON(#acc), "\"")`, nil},
		{"toJSON of array", `reduce(1..40, toJSON([#acc, #acc]), "a")`, nil},
		{"toBase64", `reduce(1..100, toBase64(#acc), "a")`, nil},
		{"split", `split(s, "")`, map[string]any{"s": strings.Repeat("a", 2_000_000)}},
		{"splitAfter", `splitAfter(s, "")`, map[string]any{"s": strings.Repeat("a", 2_000_000)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.code)
			require.NoError(t, err)

			alloc := totalAlloc(func() {
				_, err = expr.Run(program, tt.env)
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), "memory budget exceeded")
			require.Less(t, alloc, uint64(100<<20), "allocated %d bytes before failing", alloc)
		})
	}
}

func TestMemoryBudget_small_budget(t *testing.T) {
	program, err := expr.Compile(`split(s, "a")`)
	require.NoError(t, err)

	_, err = (&vm.VM{MemoryBudget: 10}).Run(program, map[string]any{"s": strings.Repeat("a", 200)})
	require.Error(t, err)
	require.Contains(t, err.Error(), "memory budget exceeded")
}

// Strings from the environment are not charged, only what the expression adds to them.
func TestMemoryBudget_large_env_strings(t *testing.T) {
	big := strings.Repeat("a", 5_000_000)
	tests := []string{
		`"prefix: " + s + " :suffix"`,
		`replace(s, "a", "b")`,
		`replace(s, "x", "yy")`,
		`upper(s)`,
		`trim(s)`,
		`toJSON(s)`,
		`len(split(s, "aaaaaaaaaa"))`,
		`repeat(s, 0)`,
		`join([s, "!"])`,
		`string(s)`,
	}
	for _, code := range tests {
		t.Run(code, func(t *testing.T) {
			_, err := expr.Eval(code, map[string]any{"s": big})
			require.NoError(t, err)
		})
	}
}

func TestMemoryBudget_reverse_before_allocation(t *testing.T) {
	input := make([]int, 1<<20)
	program, err := expr.Compile(`reverse(items)`)
	require.NoError(t, err)
	env := map[string]any{"items": input}
	var out any
	alloc := totalAlloc(func() { out, err = (&vm.VM{MemoryBudget: 8}).Run(program, env) })
	require.Error(t, err)
	require.Contains(t, err.Error(), "memory budget exceeded")
	require.Nil(t, out)
	t.Logf("reverse rejected with %d bytes allocated", alloc)
	require.Less(t, alloc, uint64(1<<20), "allocated %d bytes before rejecting reverse", alloc)
}

func TestMemoryBudget_reverse_remaining(t *testing.T) {
	for _, tt := range []struct {
		code  string
		limit uint
		fail  bool
	}{
		{`reverse(items)`, 4, false},
		{`reverse(items)`, 3, true},
		{`reverse(reverse(items))`, 7, false},
		{`reverse(reverse(items))`, 6, true},
	} {
		t.Run(fmt.Sprintf("%s/%d", tt.code, tt.limit), func(t *testing.T) {
			input := []int{1, 2, 3}
			program, err := expr.Compile(tt.code)
			require.NoError(t, err)
			_, err = (&vm.VM{MemoryBudget: tt.limit}).Run(program, map[string]any{"items": input})
			if tt.fail {
				require.Error(t, err)
				require.Contains(t, err.Error(), "memory budget exceeded")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, []int{1, 2, 3}, input)
		})
	}
}

func TestMemoryBudget_concat_before_allocation(t *testing.T) {
	input := make([]int, 1<<20)
	program, err := expr.Compile(`concat(items)`)
	require.NoError(t, err)
	env := map[string]any{"items": input}
	var out any
	alloc := totalAlloc(func() { out, err = (&vm.VM{MemoryBudget: 8}).Run(program, env) })
	require.Error(t, err)
	require.Contains(t, err.Error(), "memory budget exceeded")
	require.Nil(t, out)
	t.Logf("concat rejected with %d bytes allocated", alloc)
	require.Less(t, alloc, uint64(1<<20), "allocated %d bytes before rejecting concat", alloc)
}

func TestMemoryBudget_concat_remaining(t *testing.T) {
	for _, tt := range []struct {
		code  string
		limit uint
		fail  bool
	}{
		{`concat(items)`, 4, false},
		{`concat(items)`, 3, true},
		{`concat(items, items)`, 7, false},
		{`concat(items, items)`, 6, true},
	} {
		t.Run(fmt.Sprintf("%s/%d", tt.code, tt.limit), func(t *testing.T) {
			input := []int{1, 2, 3}
			program, err := expr.Compile(tt.code)
			require.NoError(t, err)
			_, err = (&vm.VM{MemoryBudget: tt.limit}).Run(program, map[string]any{"items": input})
			if tt.fail {
				require.Error(t, err)
				require.Contains(t, err.Error(), "memory budget exceeded")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, []int{1, 2, 3}, input)
		})
	}
}
