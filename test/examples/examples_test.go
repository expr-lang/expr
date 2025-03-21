package main

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestExamples(t *testing.T) {
	flag.Parse()

	b, err := os.ReadFile("../../testdata/examples.md")
	require.NoError(t, err)
	examples := extractCodeBlocks(string(b))

	for _, line := range examples {
		line := line
		t.Run(line, func(t *testing.T) {
			program, err := expr.Compile(line, expr.Env(nil))
			require.NoError(t, err)

			_, err = expr.Run(program, nil)
			require.NoError(t, err)
		})
	}
}

func extractCodeBlocks(markdown string) []string {
	var blocks []string
	var currentBlock []string
	inBlock := false

	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			if inBlock {
				blocks = append(blocks, strings.Join(currentBlock, "\n"))
				currentBlock = nil
				inBlock = false
			} else {
				inBlock = true
			}
			continue
		}
		if inBlock {
			currentBlock = append(currentBlock, line)
		}
	}
	return blocks
}
