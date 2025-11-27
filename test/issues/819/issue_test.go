package issue_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestIssue819(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		program, err := expr.Compile(`
			let a = [1];
			let b = type(a[0]) == 'array' ? a : [a];
			b[0][0]
		`)
		require.NoError(t, err)

		out, err := expr.Run(program, nil)
		require.NoError(t, err)
		require.Equal(t, 1, out)
	})

	t.Run("case 2", func(t *testing.T) {
		program, err := expr.Compile(`
			let range = [1,1000];
			let arr = false ? range : [range];
			map(arr, {len(#)})
		`)
		require.NoError(t, err)

		out, err := expr.Run(program, nil)
		require.NoError(t, err)
		require.Equal(t, []interface{}{2}, out)
	})
}
