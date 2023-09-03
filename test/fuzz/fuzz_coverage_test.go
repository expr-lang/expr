package fuzz_test

import (
	_ "embed"
	"github.com/antonmedv/expr/test/fuzz"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/antonmedv/expr"
)

//go:embed fuzz_corpus.txt
var fuzzCorpus string

func TestFuzzExpr_Coverage(t *testing.T) {
	inputs := strings.Split(strings.TrimSpace(fuzzCorpus), "\n")

	var env = fuzz.NewEnv()

	for _, code := range inputs {
		t.Run(code, func(t *testing.T) {
			program, err := expr.Compile(code, expr.Env(env))
			require.NoError(t, err)

			_, err = expr.Run(program, env)
			require.NoError(t, err)
		})
	}
}
