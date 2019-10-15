package main

import (
	"fmt"
	"time"

	"github.com/antonmedv/expr"
)

var expressions = []string{
	"foo > 0",
	"bar.Value in ['a', 'b', 'c']",
	"name matches '^hello.+$'",
	"now().Sub(startedAt).String()",
	"all(tweets, {.Size <= 280}) ? '👍' : '👎'",
}

var environment = map[string]interface{}{
	"foo":       1,
	"bar":       struct{ Value string }{"c"},
	"name":      "hello world",
	"startedAt": time.Now(),
	"now":       func() time.Time { return time.Now() },
	"tweets":    []tweet{},
}

type tweet struct {
	Message string
	Size    int
}

func main() {
	for _, input := range expressions {
		program, err := expr.Compile(input, expr.Env(environment))
		if err != nil {
			panic(err)
		}

		output, err := expr.Run(program, environment)
		if err != nil {
			panic(err)
		}

		fmt.Println(output)
	}
}
