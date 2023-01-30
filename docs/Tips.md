# Tips

## Reuse VM

It is possible to reuse a virtual machine between re-runs on the program.
In some cases, it can add a small increase in performance (~10%).

```go
package main

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

func main() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	program, err := expr.Compile("foo + bar", expr.Env(env))
	if err != nil {
		panic(err)
	}

	// Reuse this vm instance between runs
	v := vm.VM{}

	out, err := v.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Print(out)
}
```
