package main

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/vm"
	"github.com/chzyer/readline"
)

func main() {
	extra := []string{
		"exit",
		"opcodes",
		"map",
		"filter",
		"all",
		"any",
		"none",
		"one",
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: completer{append(builtin.Names, extra...)},
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

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

		program, err = expr.Compile(line, expr.Env(nil))
		if err != nil {
			fmt.Printf("compile error: %s\n", err)
			continue
		}
		output, err := expr.Run(program, nil)
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
