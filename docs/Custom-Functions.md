# Custom functions

User can provide custom functions in environment. 
This functions can either be defined as functions or as methods.

Functions can be typed, in which case if any of the arguments passed to such function will not match required types - error will be returned.

By default, function need to return at least one value. 
Other return values will be silently skipped.

Only exception is if second return value is `error`, in which case returned error will be returned if it is non-nil.
Read about returning errors below.

## Functions

Simple example of custom functions would be to define function in map which will be used as environment:

```go
package main

import (
	"fmt"
	"github.com/antonmedv/expr"
)

func main() {
	env := map[string]interface{}{
		"foo": 1,
		"double": func(i int) int { return i * 2 },
	}

	out, err := expr.Eval("double(foo)", env)

	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

## Methods

Methods can be defined on type that is provided as environment.

Methods MUST be exported in order to be callable.

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

// Methods defined on such struct will be functions.
func (Env) Format(t time.Time) string { return t.Format(time.RFC822) }

type Tweet struct {
	Text string
	Date time.Time
}

func main() {
	code := `map(filter(Tweets, {len(.Text) > 0}), {.Text + Format(.Date)})`

	// We can use an empty instance of the struct as an environment.
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		panic(err)
	}

	env := Env{
		Tweets: []Tweet{{"Oh My God!", time.Now()}, {"How you doin?", time.Now()}, {"Could I be wearing any more clothes?", time.Now()}},
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
}
```

## Fast functions

Fast functions are functions that don't use reflection for calling them. 
This improves performance but drops ability to have typed arguments.

Such functions have strict signatures for them:
```go
func(...interface{}) interface{}
```
or
```go
func(...interface{}) (interface{}, error)
```

Methods can also be used as fast functions if they will have signature specified above.

Example:
```go
package main

import (
	"fmt"
	"github.com/antonmedv/expr"
)

type Env map[string]interface{}

func (Env) FastMethod(...interface{}) interface{} {
	return "Hello, "
}

func main() {
	env := Env{
		"fast_func": func(...interface{}) interface{} { return "world" },
	}

	out, err := expr.Eval("FastMethod() + fast_func()", env)

	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

## Returning errors

Both normal and fast functions can return `error`s as second return value.
In this case if function will return any value and non-nil error - such error will be returned to the caller.

```go
package main

import (
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
)

func main() {
	env := map[string]interface{}{
		"foo": -1,
		"double": func(i int) (int, error) {
			if i < 0 {
				return 0, errors.New("value cannot be less than zero")
			}
			return i * 2, nil
		},
	}

	out, err := expr.Eval("double(foo)", env)

	// This `err` will be the one returned from `double` function.
	// err.Error() == "value cannot be less than zero"
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

* [Contents](README.md)
* Next: [Operator Override](Operator-Override.md)
