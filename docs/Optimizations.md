# Optimizations

Expr has a bunch of optimization which will produce more optimal program during compile step.

## In array

```js
value in ['foo', 'bar', 'baz']
```

If expr finds an `in` or `not in` expression with an array, it will be transformed into:

```js
value in {"foo": true, "bar": true, "baz": true}
```

## Constant folding

Arithmetic expressions with constants is computed on compile step and replaced with result.

```js
-(2-5)**3-2/(+4-3)+-2
```

Will be compiled to just single number:

```js
23
```

So in expressions it's safe to use some arithmetics for better readability:

```js
percentage > 0.3 * 100
```

As it will be simpified to:

```js
percentage > 30
```

## In range

```js
user.Age in 18..32
```

Will be replaced with binary operator:

```js
18 <= user.Age && user.Age <= 32
```

`not in` operator will also work.

## Const range

```js
1..10_000
```

Ranges computed on compile stage, repleced with preallocated slices.

## Const expr

If some function marked as constant expression with `expr.ConstExpr`. It will be replaced with result
of call, if all arguments are constants.

```go
expr.ConstExpt("fib")
```

```js
fib(42)
``` 

Will be replaced with result of `fib(42)` on compile step. No need to calculate it during runtime.

[ConstExpr Example](https://pkg.go.dev/github.com/antonmedv/expr?tab=doc#ConstExpr)

## Reuse VM

It is possible to reuse a virtual machine between re-runs on the program. 
This adds a small increase in performance (from 4% to 40% depending on a program).

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

## Reduced use of reflect

To fetch fields from struct, values from map, get by indexes expr uses reflect package. 
Envs can implement vm.Fetcher interface, to avoid use reflect:
```go
type Fetcher interface {
	Fetch(interface{}) interface{}
}
```
When you need to fetch a field, the method will be used instead reflect functions.
If the field is not found, Fetch must return nil.
To generate Fetch for your types, use [Exprgen](Exprgen.md).


* [Contents](README.md)
