package main

import (
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/checker"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/parser"
)

func TestIssue844(t *testing.T) {
	testCases := []struct {
		name       string
		env        any
		expression string
		shouldFail bool
	}{
		{
			name:       "exported env, exported field",
			env:        ExportedEnv{},
			expression: `ExportedEmbedded`,
			shouldFail: false,
		},
		{
			name:       "exported env, unexported field",
			env:        ExportedEnv{},
			expression: `unexportedEmbedded`,
			shouldFail: true,
		},
		{
			name:       "exported env, exported field inherited from exported field",
			env:        ExportedEnv{},
			expression: `Str`,
			shouldFail: false,
		},
		{
			name:       "exported env, unexported field inherited from exported field",
			env:        ExportedEnv{},
			expression: `str`,
			shouldFail: true,
		},
		{
			name:       "exported env, exported field inherited from exported field",
			env:        ExportedEnv{},
			expression: `Integer`,
			shouldFail: false,
		},
		{
			name:       "exported env, unexported field inherited from exported field",
			env:        ExportedEnv{},
			expression: `integer`,
			shouldFail: true,
		},
		{
			name:       "exported env, exported field directly accessed from exported field",
			env:        ExportedEnv{},
			expression: `ExportedEmbedded.Str`,
			shouldFail: false,
		},
		{
			name:       "exported env, unexported field directly accessed from exported field",
			env:        ExportedEnv{},
			expression: `ExportedEmbedded.str`,
			shouldFail: true,
		},
		{
			name:       "exported env, exported field directly accessed from exported field",
			env:        ExportedEnv{},
			expression: `unexportedEmbedded.Integer`,
			shouldFail: true,
		},
		{
			name:       "exported env, unexported field directly accessed from exported field",
			env:        ExportedEnv{},
			expression: `unexportedEmbedded.integer`,
			shouldFail: true,
		},
		{
			name:       "unexported env, exported field",
			env:        unexportedEnv{},
			expression: `ExportedEmbedded`,
			shouldFail: false,
		},
		{
			name:       "unexported env, unexported field",
			env:        unexportedEnv{},
			expression: `unexportedEmbedded`,
			shouldFail: true,
		},
		{
			name:       "unexported env, exported field inherited from exported field",
			env:        unexportedEnv{},
			expression: `Str`,
			shouldFail: false,
		},
		{
			name:       "unexported env, unexported field inherited from exported field",
			env:        unexportedEnv{},
			expression: `str`,
			shouldFail: true,
		},
		{
			name:       "unexported env, exported field inherited from exported field",
			env:        unexportedEnv{},
			expression: `Integer`,
			shouldFail: false,
		},
		{
			name:       "unexported env, unexported field inherited from exported field",
			env:        unexportedEnv{},
			expression: `integer`,
			shouldFail: true,
		},
		{
			name:       "unexported env, exported field directly accessed from exported field",
			env:        unexportedEnv{},
			expression: `ExportedEmbedded.Str`,
			shouldFail: false,
		},
		{
			name:       "unexported env, unexported field directly accessed from exported field",
			env:        unexportedEnv{},
			expression: `ExportedEmbedded.str`,
			shouldFail: true,
		},
		{
			name:       "unexported env, exported field directly accessed from exported field",
			env:        unexportedEnv{},
			expression: `unexportedEmbedded.Integer`,
			shouldFail: true,
		},
		{
			name:       "unexported env, unexported field directly accessed from exported field",
			env:        unexportedEnv{},
			expression: `unexportedEmbedded.integer`,
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := conf.New(tc.env)
			tree, err := parser.ParseWithConfig(tc.expression, config)
			require.NoError(t, err)

			_, err = new(checker.Checker).PatchAndCheck(tree, config)
			if tc.shouldFail {
				require.Error(t, err)
				errStr := err.Error()
				if !strings.Contains(errStr, "unknown name") &&
					!strings.Contains(errStr, " has no field ") {
					t.Fatalf("expected a different error, got: %v", err)
				}
			} else {
				require.NoError(t, err)
			}

			// We add this because the issue was actually not catching something
			// that sometimes failed with the error:
			//	reflect.Value.Interface: cannot return value obtained from unexported field or method
			// This way, we test that everything we allow passing is also
			// allowed later
			_, err = expr.Eval(tc.expression, tc.env)
			if tc.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type ExportedEnv struct {
	ExportedEmbedded
	unexportedEmbedded
}

type unexportedEnv struct {
	ExportedEmbedded
	unexportedEmbedded
}

type ExportedEmbedded struct {
	Str string
	str string
}

type unexportedEmbedded struct {
	Integer int
	integer int
}
