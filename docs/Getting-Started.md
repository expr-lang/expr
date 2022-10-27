# Getting Started

**Expr** provides a package for evaluating arbitrary expressions as well as type
checking of such expression.

## Evaluate

For simple use cases with one execution of an expression, you can use the 
`expr.Eval` function. It takes an expression and a map of variables and returns
the result of the expression.

```go
package main

import (
	"fmt"
	"github.com/antonmedv/expr"
)

func main() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	out, err := expr.Eval("foo + bar", env)

	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

## Compile

Usually we want to compile, type check and verify what the expression returns a 
boolean (or another type). 

For example, if a user saves an expression from a
[web UI](https://antonmedv.github.io/expr/).

```go
package main

import (
	"fmt"

	"github.com/antonmedv/expr"
)

func main() {
	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf, // Use any functions.
	}

	code := `sprintf(greet, names[0])`

	// Compile the code into a bytecode. This step can be done once 
	// and program may be reused. Specify an environment for type check.
	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
}
```

You may use existing types. 

For example, an environment can be a struct. And structs methods can be used as
functions. Expr supports embedded structs and methods defines on them too.

The struct fields can be renamed by adding struct tags such as `expr:"name"`.

```go
package main

import (
	"fmt"
	"time"

	"github.com/antonmedv/expr"
)

type Env struct {
	Tweets []Tweet
}

// Methods defined on the struct become functions.
func (Env) Format(t time.Time) string { return t.Format(time.RFC822) }

type Tweet struct {
	Text string
	Date time.Time `expr:"timestamp"`
}

func main() {
	code := `map(filter(Tweets, {len(.Text) > 0}), {.Text + Format(.timestamp)})`

	// We can use an empty instance of the struct as an environment.
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		panic(err)
	}

	env := Env{
		Tweets: []Tweet{
			{"Oh My God!", time.Now()}, 
			{"How you doin?", time.Now()}, 
			{"Could I be wearing any more clothes?", time.Now()},
		},
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
}
```

* Next: [Operator Overloading](Operator-Overloading.md)
