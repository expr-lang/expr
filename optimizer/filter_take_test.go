package optimizer_test

import (
	"fmt"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/vm"
)

func TestOptimize_filter_take(t *testing.T) {
	tests := []struct {
		input string
		want  any
		calls int
	}{
		{`take(filter(items, keep(#)), 2)`, []any{2, 4}, 4},
		{`take(filter(items, keep(#)), 1)`, []any{2}, 2},
		{`take(filter(items, keep(#)), 1 + 1)`, []any{2, 4}, 4},
		{`take(take(filter(items, keep(#)), 1), 3)`, []any{2}, 2},
		{`take(take(filter(items, keep(#)), 3), 1)`, []any{2}, 2},
		{`take(filter(items, keep(#)), 9)`, []any{2, 4, 6, 8}, 8},
		{`take(filter(items, keep(#) && # > 6), 1)`, []any{8}, 8},
		{`take(filter(items, keep(#) && # > 8), 2)`, []any{}, 8},
		{`take(filter([], keep(#)), 2)`, []any{}, 0},
		{`take(map(filter(items, keep(#)), # * 10), 2)`, []any{20, 40}, 4},
		{`map(take(filter(items, keep(#)), 2), #index)`, []any{0, 1}, 4},
		{`take(map(filter(items, keep(#)), #index), 2)`, []any{0, 1}, 8},
		{`[99, take(filter(items, keep(#)), 2), 42]`, []any{99, []any{2, 4}, 42}, 4},
		{`len(take(filter(items, keep(#)), 2))`, 2, 4},
		{`let a = take(filter(items, keep(#)), 2); [a, count(items, # > 4)]`, []any{[]any{2, 4}, 4}, 4},
		{`map([1, 2], [#index, take(filter(items, keep(#)), 2), #])`, []any{[]any{0, []any{2, 4}, 1}, []any{1, []any{2, 4}, 2}}, 8},
		{`take(filter(items, keep(#)), limit)`, []any{2, 4}, 8},
		{`take(filter(items, keep(#)), 0)`, []any{}, 8},
		{fmt.Sprintf(`take(filter(items, keep(#)), %d)`, int(^uint(0)>>1)), []any{2, 4, 6, 8}, 8},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			calls := 0
			env := map[string]any{
				"items": []int{1, 2, 3, 4, 5, 6, 7, 8},
				"limit": 2,
				"keep": func(value int) bool {
					calls++
					return value%2 == 0
				},
			}
			program, err := expr.Compile(tt.input, expr.Env(env))
			require.NoError(t, err)

			machine := vm.VM{}
			for run := 0; run < 2; run++ {
				calls = 0
				output, err := machine.Run(program, env)
				require.NoError(t, err)
				assert.Equal(t, tt.want, output)
				assert.Equal(t, tt.calls, calls)
				assert.Empty(t, machine.Stack)
				assert.Empty(t, machine.Scopes)
			}
		})
	}
}

func TestOptimize_filter_take_errors(t *testing.T) {
	for _, failAt := range []int{3, 5} {
		t.Run(fmt.Sprint(failAt), func(t *testing.T) {
			calls := 0
			env := map[string]any{
				"items": []int{1, 2, 3, 4, 5, 6},
				"keep": func(value int) (bool, error) {
					calls++
					if value == failAt {
						return false, fmt.Errorf("predicate failed")
					}
					return value%2 == 0, nil
				},
			}
			program, err := expr.Compile(`take(filter(items, keep(#)), 2)`, expr.Env(env))
			require.NoError(t, err)
			output, err := expr.Run(program, env)
			if failAt == 3 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "predicate failed")
				assert.Equal(t, 3, calls)
			} else {
				require.NoError(t, err)
				assert.Equal(t, []any{2, 4}, output)
				assert.Equal(t, 4, calls)
			}
		})
	}

	t.Run("fused map", func(t *testing.T) {
		calls := 0
		env := map[string]any{
			"items": []int{1, 2, 3, 4, 5, 6},
			"project": func(value int) (int, error) {
				calls++
				if value == 4 {
					return 0, fmt.Errorf("map failed")
				}
				return value * 10, nil
			},
		}
		program, err := expr.Compile(`take(map(filter(items, # % 2 == 0), project(#)), 2)`, expr.Env(env))
		require.NoError(t, err)
		_, err = expr.Run(program, env)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "map failed")
		assert.Equal(t, 2, calls)
	})
}

func TestOptimize_filter_take_unchanged(t *testing.T) {
	calls := 0
	env := map[string]any{
		"items": []int{1, 2, 3, 4},
		"keep": func(value int) bool {
			calls++
			return value%2 == 0
		},
	}
	program, err := expr.Compile(`take(filter(items, keep(#)), -1)`, expr.Env(env))
	require.NoError(t, err)
	_, err = expr.Run(program, env)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot take negative number of elements")
	assert.Equal(t, 4, calls)

	for _, code := range []string{
		`take(filter(items, keep(#)))`,
		`take(filter(items, keep(#)), 1.5)`,
		`take(filter(items, keep(#)), "2")`,
	} {
		_, err := expr.Compile(code, expr.Env(env))
		require.Error(t, err)
	}

	calls = 0
	program, err = expr.Compile(`take(filter(items, keep(#)), 1)`, expr.Env(env), expr.Optimize(false))
	require.NoError(t, err)
	output, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, []any{2}, output)
	assert.Equal(t, 4, calls)

	calls = 0
	takeCalls := 0
	program, err = expr.Compile(`take(filter(items, keep(#)), 1)`, expr.Env(env),
		expr.Function("take", func(args ...any) (any, error) {
			takeCalls++
			return args[0], nil
		}))
	require.NoError(t, err)
	output, err = expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, []any{2, 4}, output)
	assert.Equal(t, 4, calls)
	assert.Equal(t, 1, takeCalls)
}

func TestOptimize_filter_take_memory_budget(t *testing.T) {
	env := map[string]any{"items": []int{1, 2, 3, 4, 5, 6}}
	program, err := expr.Compile(`take(filter(items, # % 2 == 0), 2)`, expr.Env(env))
	require.NoError(t, err)
	machine := vm.VM{MemoryBudget: 2}
	_, err = machine.Run(program, env)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "memory budget exceeded")

	machine.MemoryBudget = 3
	output, err := machine.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, []any{2, 4}, output)
	assert.Empty(t, machine.Stack)
	assert.Empty(t, machine.Scopes)
}

func BenchmarkFilterTake(b *testing.B) {
	items := make([]int, 10000)
	for i := range items {
		items[i] = i
	}
	env := map[string]any{"items": items}
	for _, tt := range []struct {
		name  string
		input string
	}{
		{"early", `take(filter(items, # % 2 == 0), 10)`},
		{"insufficient", `take(filter(items, # < 10), 20)`},
		{"dense_insufficient", `take(filter(items, # % 2 == 0), 10000)`},
	} {
		b.Run(tt.name, func(b *testing.B) {
			program, err := expr.Compile(tt.input, expr.Env(env))
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := expr.Run(program, env)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
