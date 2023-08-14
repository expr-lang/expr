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
	items := []readline.PrefixCompleterInterface{}
	for _, name := range builtin.Names {
		items = append(items, readline.PcItem(name))
	}
	items = append(items, readline.PcItem("exit"))
	items = append(items, readline.PcItem("opcodes"))
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: readline.NewPrefixCompleter(items...),
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
