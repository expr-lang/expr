package pipes_test

import (
	"fmt"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/require"
)

func TestPipes(t *testing.T) {
	env := map[string]interface{}{
		"sprintf": fmt.Sprintf,
	}

	program, err := expr.Compile(`"%s bar %d" | sprintf("foo", -42 | abs())`, expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, "foo bar 42", out)
}

func TestPipes_map_filter(t *testing.T) {
	program, err := expr.Compile(`1..9 | map(# + 1) | filter(# % 2 == 0)`)
	require.NoError(t, err)

	out, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, []interface{}{2, 4, 6, 8, 10}, out)
}
