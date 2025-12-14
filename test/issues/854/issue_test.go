package main

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/types"
)

func TestIssue854(t *testing.T) {
	envType := types.Map{
		"user": types.Map{
			// If we do not specify `Profile` here,
			// this is a correct behavior to throw
			// on a missing property.
		},
	}

	code := `user.Profile?.Address ?? "Unknown address"`

	_, err := expr.Compile(code, expr.Env(envType))
	require.Error(t, err)
}
