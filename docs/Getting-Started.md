# Getting Started

**Expr** provides a package for evaluating arbitrary expressions as well as type checking of such expression.
 For demonstration purpose, let's assume that we already have some types and structs describing our data. And we want to implement filtering of data flow, via providing our users to configure such filters via expressions.
 
```go
type Requset struct {
	Location string
	Date     time.Time
	Ticket   Ticket
}

type Ticket struct {
	Segments []Segment
}

type Segment struct {
	Origin, Destination string
	Date                time.Time
} 
```

First what we need to do is to create a way users will be creating, editing (and maybe deleted) filter rules, but this is out of scope for this article. Let's assume that we have some kind of Web interface where users can do all of this task. On creation or deletion of rules, we much check if rules are written correctly. And we want to give the user access to `Request` fields and only these fields. 

```go
import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

var program *vm.Program

program, err = expr.Compile(rule, expr.Env(&Request{}))
```

By passing `&Request{}` as env into compile method we turned on type checking, now if users make an error in field or compare strings with integers, he will get an error. Now users can save expressions.

```coffeescript
all(Ticket.Segments, {.Origin == Location}) and Date.Before(Ticket.Segments[0].Date)
``` 

Now we need to implement the execution of our compiled program. On each request, we have some kind of filter loop.

```go
output, err := expr.Run(program, requset)
if err != nil {
	return err
}

if !output.(bool) {
	continue
}
```

Now let's some add a function for repetitive tasks. 

```go
func (r *Requset) SameLocation() bool {
	same := false
	for _, s := range r.Ticket.Segments {
		same = same && s.Origin == r.Location
	}
	return same
}
```

Now users can use functions inside expressions.

```coffeescript
SameLocation() and Date.Before(Ticket.Segments[0].Date)
```

Much better. But using time's package methods isn't pretty. What if we can override operators? And we can! Let's describe another function.

```go
func (*Requset) Before(a, b time.Time) bool {
	return a.Before(b)
}
```

Now, on compile time, override `<` operator.

```go
program, err = expr.Compile(rule, expr.Env(&Request{}), expr.Operator("<", "Before"))
```

That's it! Now users can write expressions in a more pleasant way. Other operators `>`, `===`, etc can be overridden if similar way.

```coffeescript
SameLocation() and Date < Ticket.Segments[0].Date
```

Next
* [Language Definition](Language-Definition.md)
* [Usage](Usage.md)
