package main

import (
	"os"
	"strings"
	"testing"

	"expr/internal/testify/require"

	"expr"
)

func TestGenerated(t *testing.T) {
	b, err := os.ReadFile("../../testdata/examples.txt")
	require.NoError(t, err)

	examples := strings.TrimSpace(string(b))
	for _, line := range strings.Split(examples, "\n") {
		t.Run(line, func(t *testing.T) {
			program, err := expr.Compile(line, expr.Env(env))
			require.NoError(t, err)

			_, err = expr.Run(program, env)
			require.NoError(t, err)
		})
	}
}
