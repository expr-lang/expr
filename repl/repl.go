package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/debug"
	"github.com/antonmedv/expr/vm"
	"github.com/chzyer/readline"
)

var keywords = []string{
	// Commands:
	"exit", "opcodes", "debug",

	// Operators:
	"and", "or", "in", "not", "not in",
	"contains", "matches", "startsWith", "endsWith",
}

func main() {
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

	env := map[string]interface{}{
		"ENV": os.Environ(),
	}
	var program *vm.Program

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF when Ctrl-D is pressed
			break
		}
		line = strings.TrimSpace(line)

		if line == "exit" {
			break
		}

		if line == "opcodes" {
			if program == nil {
				fmt.Println("no program")
				continue
			}
			fmt.Println(program.Disassemble())
			continue
		}

		if line == "debug" {
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
		output, err := expr.Run(program, env)
		if err != nil {
			fmt.Printf("runtime error: %s\n", err)
			continue
		}
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
