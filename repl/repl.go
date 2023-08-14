package main

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: "> ",
		AutoComplete: readline.NewPrefixCompleter(
			readline.PcItem("exit"),
			readline.PcItem("opcodes"),
			readline.PcItem("len"),
			readline.PcItem("abs"),
			readline.PcItem("int"),
			readline.PcItem("float"),
			readline.PcItem("map"),
			readline.PcItem("filter"),
			readline.PcItem("all"),
			readline.PcItem("any"),
			readline.PcItem("none"),
			readline.PcItem("one"),
			readline.PcItem("count"),
			readline.PcItem("sum"),
		),
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
