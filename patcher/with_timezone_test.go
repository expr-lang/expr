package patcher_test

import (
	"testing"
	"time"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
)

func TestWithTimezone_date(t *testing.T) {
	program, err := expr.Compile(`date("2024-05-07 23:00:00")`, expr.Timezone("Europe/Zurich"))
	require.NoError(t, err)

	out, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, "2024-05-07T23:00:00+02:00", out.(time.Time).Format(time.RFC3339))
}

func TestWithTimezone_now(t *testing.T) {
	program, err := expr.Compile(`now()`, expr.Timezone("Asia/Kamchatka"))
	require.NoError(t, err)

	out, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, "Asia/Kamchatka", out.(time.Time).Location().String())
}
