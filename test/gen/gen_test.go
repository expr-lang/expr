package main

import (
	"bufio"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

var updateFlag = flag.Bool("update", false, "Drop failing lines from examples.txt")

func TestGenerated(t *testing.T) {
	flag.Parse()

	b, err := os.ReadFile("../../testdata/generated.txt")
	require.NoError(t, err)

	examples := strings.TrimSpace(string(b))
	var validLines []string

	for _, line := range strings.Split(examples, "\n") {
		line := line
		t.Run(line, func(t *testing.T) {
			program, err := expr.Compile(line, expr.Env(Env))
			if err != nil {
				if !*updateFlag {
					t.Errorf("Compilation failed: %v", err)
				}
				return
			}

			_, err = expr.Run(program, Env)
			if err != nil {
				if !*updateFlag {
					t.Errorf("Execution failed: %v", err)
				}
				return
			}

			validLines = append(validLines, line)
		})
	}

	if *updateFlag {
		file, err := os.Create("../../testdata/examples.txt")
		if err != nil {
			t.Fatalf("Failed to update examples.txt: %v", err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		writer := bufio.NewWriter(file)
		for _, line := range validLines {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				t.Fatalf("Failed to write to examples.txt: %v", err)
			}
		}
		_ = writer.Flush()
	}
}
