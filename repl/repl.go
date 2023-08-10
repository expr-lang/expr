package main

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF when Ctrl-D is pressed
			break
		}

		if line == "exit" {
			break
		}

		result := evaluate(line)
		fmt.Println(result)
	}
}

func evaluate(input string) string {
	result, err := expr.Eval(input, nil)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v", result)
}
