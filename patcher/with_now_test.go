package patcher_test

import (
	"testing"
	"time"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
)

func TestWithNow(t *testing.T) {
	now := time.Date(2024, 5, 7, 23, 0, 0, 0, time.UTC)

	program, err := expr.Compile(`now()`, expr.MockNow(now))
	require.NoError(t, err)

	out, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, now, out)
}

func TestWithNow_is_deterministic(t *testing.T) {
	now := time.Date(2024, 5, 7, 23, 0, 0, 0, time.UTC)

	program, err := expr.Compile(`now().Year()`, expr.MockNow(now))
	require.NoError(t, err)

	out, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, 2024, out)
}

func TestWithNow_combined_with_timezone(t *testing.T) {
	now := time.Date(2024, 5, 7, 23, 0, 0, 0, time.UTC)

	program, err := expr.Compile(`now()`, expr.Timezone("Asia/Kamchatka"), expr.MockNow(now))
	require.NoError(t, err)

	out, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, "Asia/Kamchatka", out.(time.Time).Location().String())
	require.True(t, out.(time.Time).Equal(now))
}
