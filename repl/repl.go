package main

import (
	"fmt"
	"github.com/antonmedv/expr/test/fuzz"
	"os"
	"runtime"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/debug"
	"github.com/antonmedv/expr/vm"
	"github.com/bettercap/readline"
)

var keywords = []string{
	// Commands:
	"exit", "opcodes", "debug", "mem",

	// Operators:
	"and", "or", "in", "not", "not in",
	"contains", "matches", "startsWith", "endsWith",
}

func main() {
	env := fuzz.NewEnv()
	for name := range env {
		keywords = append(keywords, name)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "❯ ",
		AutoComplete: completer{append(builtin.Names, keywords...)},
		HistoryFile:  home + "/.expr_history",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var memUsage uint64
	var program *vm.Program

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF when Ctrl-D is pressed
			break
		}
		line = strings.TrimSpace(line)

		switch line {
		case "":
			continue

		case "exit":
			return

		case "mem":
			fmt.Printf("memory usage: %s\n", humanizeBytes(memUsage))
			continue

		case "opcodes":
			if program == nil {
				fmt.Println("no program")
				continue
			}
			fmt.Println(program.Disassemble())
			continue

		case "debug":
			if program == nil {
				fmt.Println("no program")
				continue
			}
			debug.StartDebugger(program, env)
			continue
		}

		program, err = expr.Compile(line, expr.Env(env))
		if err != nil {
			fmt.Printf("compile error: %s\n", err)
			continue
		}

		start := memoryUsage()
		output, err := expr.Run(program, env)
		if err != nil {
			fmt.Printf("runtime error: %s\n", err)
			continue
		}
		memUsage = memoryUsage() - start

		fmt.Println(output)
	}
}

type completer struct {
	words []string
}

func (c completer) Do(line []rune, pos int) ([][]rune, int) {
	var lastWord string
	for i := pos - 1; i >= 0; i-- {
		if line[i] == ' ' {
			break
		}
		lastWord = string(line[i]) + lastWord
	}

	var words [][]rune
	for _, word := range c.words {
		if strings.HasPrefix(word, lastWord) {
			words = append(words, []rune(strings.TrimPrefix(word, lastWord)))
		}
	}

	return words, len(lastWord)
}

func memoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func humanizeBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
