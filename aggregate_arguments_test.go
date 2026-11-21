package expr_test

import (
	"fmt"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestCollectionBuiltinArgumentsScope(t *testing.T) {
	groups := []map[string]any{
		{"items": []map[string]any{{"value": 10, "initial": 1, "direction": "asc"}, {"value": 20, "initial": 2, "direction": "asc"}}, "initial": 100, "direction": "desc"},
		{"items": []map[string]any{{"value": 30, "initial": 3, "direction": "desc"}, {"value": 40, "initial": 4, "direction": "desc"}}, "initial": 200, "direction": "asc"},
	}
	tests := []struct {
		name string
		code string
		env  map[string]any
		want any
	}{
		{"reduce outer field", `map(groups, reduce(.items, #acc + .value, .initial))`, map[string]any{"groups": groups}, []any{130, 270}},
		{"sortBy outer field", `map(groups, map(sortBy(.items, .value, .direction), .value))`, map[string]any{"groups": groups}, []any{[]any{20, 10}, []any{30, 40}}},
		{"reduce outer index", `map(groups, reduce(#, #acc + #, #index))`, map[string]any{"groups": [][]int{{10, 20}, {30, 40}, {50, 60}}}, []any{30, 71, 112}},
		{"sortBy outer index", `map(groups, sortBy(#, #, #index == 0 ? "asc" : "desc"))`, map[string]any{"groups": [][]int{{2, 1}, {3, 4}, {5, 6}}}, []any{[]any{1, 2}, []any{4, 3}, []any{6, 5}}},
		{"reduce outer accumulator", `reduce([10, 20], reduce([1, 2], #acc + #, #acc), 100)`, nil, 106},
		{"reduce nested argument", `map(groups, reduce(#, #acc + #, reduce([1, 2], #acc + #, #index)))`, map[string]any{"groups": [][]int{{10, 20}, {30, 40}}}, []any{33, 74}},
		{"sortBy nested argument", `map(groups, sortBy(#, #, reduce([1], #acc, #index) == 0 ? "asc" : "desc"))`, map[string]any{"groups": [][]int{{2, 1}, {3, 4}}}, []any{[]any{1, 2}, []any{4, 3}}},
		{"reduce empty input", `map(groups, reduce(.items, #acc + #, .initial))`, map[string]any{"groups": []map[string]any{{"items": []int{}, "initial": 100}, {"items": []int{}, "initial": 200}}}, []any{100, 200}},
		{"sortBy empty input", `map(groups, sortBy(.items, #, .direction))`, map[string]any{"groups": []map[string]any{{"items": []int{}, "direction": "desc"}}}, []any{[]any{}}},
		{"reduce default", `map(groups, reduce(#, #acc + #))`, map[string]any{"groups": [][]int{{10, 20}, {30, 40}}}, []any{30, 70}},
		{"sortBy default", `map(groups, sortBy(#, #))`, map[string]any{"groups": [][]int{{2, 1}, {4, 3}}}, []any{[]any{1, 2}, []any{3, 4}}},
		{"reduce predicate index", `reduce([10, 20, 30], #acc + #index, 0)`, nil, 3},
		{"sortBy predicate value", `sortBy([10, 20, 30], -#)`, nil, []any{30, 20, 10}},
	}
	for _, tt := range tests {
		for _, optimize := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/optimize=%v", tt.name, optimize), func(t *testing.T) {
				program, err := expr.Compile(tt.code, expr.Env(tt.env), expr.Optimize(optimize))
				require.NoError(t, err)
				got, err := expr.Run(program, tt.env)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			})
		}
	}
}

func TestCollectionBuiltinArgumentsEvaluationOrder(t *testing.T) {
	for _, builtin := range []string{"reduce", "sortBy"} {
		for _, empty := range []bool{false, true} {
			for _, optimize := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/empty=%v/optimize=%v", builtin, empty, optimize), func(t *testing.T) {
					var calls []string
					items := [][]int{{10, 20}, {30, 40}}
					if empty {
						items = [][]int{{}, {}}
					}
					env := map[string]any{
						"groups": items,
						"collection": func(values []int, index int) []int {
							calls = append(calls, fmt.Sprintf("collection:%d", index))
							return values
						},
						"argument": func(index int) int {
							calls = append(calls, fmt.Sprintf("argument:%d", index))
							return index
						},
						"order": func(index int) string {
							calls = append(calls, fmt.Sprintf("argument:%d", index))
							return "desc"
						},
						"predicate": func(index int) int {
							calls = append(calls, fmt.Sprintf("predicate:%d", index))
							return index
						},
					}
					code := `map(groups, reduce(collection(#, #index), #acc + # + predicate(#index), argument(#index)))`
					var want any = []any{31, 72}
					if builtin == "sortBy" {
						code = `map(groups, sortBy(collection(#, #index), # + predicate(#), order(#index)))`
						want = []any{[]any{20, 10}, []any{40, 30}}
					}
					expectedCalls := []string{"collection:0", "argument:0", "predicate:0", "predicate:1", "collection:1", "argument:1", "predicate:0", "predicate:1"}
					if builtin == "sortBy" {
						expectedCalls = []string{"collection:0", "argument:0", "predicate:10", "predicate:20", "collection:1", "argument:1", "predicate:30", "predicate:40"}
					}
					if empty {
						expectedCalls = []string{"collection:0", "argument:0", "collection:1", "argument:1"}
						want = []any{0, 1}
						if builtin == "sortBy" {
							want = []any{[]any{}, []any{}}
						}
					}
					program, err := expr.Compile(code, expr.Env(env), expr.Optimize(optimize))
					require.NoError(t, err)
					got, err := expr.Run(program, env)
					require.NoError(t, err)
					require.Equal(t, expectedCalls, calls)
					require.Equal(t, want, got)
				})
			}
		}
	}
}
