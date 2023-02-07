# Operator Overloading

**Expr** supports operator overloading. For example, you may rewrite such expression:

```
Now().Sub(CreatedAt) 
```

To use `-` operator:
 
```
Now() - CreatedAt
```

To overload the operator use [Operator](https://pkg.go.dev/github.com/antonmedv/expr?tab=doc#Operator) option:

```go
func main() {
	code := `Now() - CreatedAt`

	options := []expr.Option{
		expr.Env(Env{}),
		expr.Operator("-", "Sub"), // Replace `-` operator with function `Sub`.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		panic(err)
	}

	env := Env{
		CreatedAt: time.Date(1987, time.November, 24, 20, 0, 0, 0, time.UTC),
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

type Env struct {
	CreatedAt time.Time
}

func (Env) Now() time.Time                   { return time.Now() }
func (Env) Sub(a, b time.Time) time.Duration { return a.Sub(b) }
```

**Expr** uses functions from `Env` for operator overloading. If types of 
operands match types of a function, the operator will be replaced with a 
function call.
