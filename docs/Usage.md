# Usage

**Expr** package can compile and evaluate expressions.
A simple example of an expression is `1 + 2`. You can also use more complicated expressions, such as `foo[3].Method('bar')`. 
See [Language Definition](Language-Definition.md) to learn the syntax of the **expr** package.

The package provides 2 ways to work with expressions:

* **compile**: the expression is compiled to bytecode, so it can be cached and evaluated.
* **evaluation**: the expression is evaluated using our own virtual machine;

```go
import "gopkg.in/antonmedv/expr.v2"

program, err := expr.Compile(`1 + 2`)

output, err := expr.Run(program, nil) 

fmt.Println(out) // outputs 3
```

## Passing in Variables

You can also pass variables into the expression, which can be map or struct:

```go
env := map[string]interface{}{
	"Foo": ...
	"Bar": ...
}

// or
env := Env{
	Foo: ...
	Bar: ...
}

// Pass env option to compile for static type checking.
program, err := expr.Compile(`Foo == Bar`, expr.Env(env))

output, err := expr.Run(program, env) 
```

Expr uses reflection for accessing and iterating passed data. 
For example you can pass nested structures without any modification or preparation:

```go
type Cookie struct {
	Key   string
	Value string
}
type User struct {
	UserAgent string
	Cookies   []Cookie
}
type Request struct {
	User *user
}

req := Request{
	User: &User{
        Cookies:   []Cookie{{"origin", "www"}},
	    UserAgent: "Firefox",
    },
}

program, err := expr.Compile(`User.UserAgent matches "^Fire.+$" and User.Cookies[0].Value == "www"`, expr.Env(env))

output, err := expr.Run(program, env) 
```

## Passing in Functions

You can also pass functions into the expression:

```go
env := map[string]interface{}{
	"Request": req,
	"Values": func(xs []Cookie) []string {
		vs := make([]string, 0)
		for _, x := range xs {
			vs = append(vs, x.Value)
		}
		return vs
	},
}

program, err := expr.Compile(`"www" in Values(Request.User.Cookies)`, expr.Env(env))

output, err := expr.Run(program, env) 
```

### Struct's methods

All methods of passed struct also available as functions inside expr:

```go
type Env struct {
	value int
}

func (e *Env) Value() int {
	return e.value
}

program, err := expr.Compile(`Value()`, expr.Env(&Env{}))

output, err := expr.Run(program, &Env{1}) 
```

### Map types

As well as methods defined of map types.

```go
type Env map[string]interface{}

func (Env) Swipe(in string) string {
	return strings.Replace(in, "world", "user", 1)
}

env := Env{
	"greeting": "hello world",
}

program, err := expr.Compile(`Swipe(greeting)`, expr.Env(env))

output, err := expr.Run(program, env)

fmt.Println(out) // outputs "hello user"
```

### Embedded structs

If you have lots of different contexts for expressions, but want to separate functionality you can use embedded structs.

```go
type EnvContextOne struct {
    *Helpers
	Price int
}

type EnvContextTwo struct {
    *Helpers
	City *City
}

type Helpers struct {
	Value int
}

func (h *Helper) IsMore(i int) bool {
	return i > h.Value
}

program, err := expr.Compile(`IsMore(Price)`, expr.Env(&EnvContextOne{}))

output, err := expr.Run(program, &EnvContextOne{...})

// ...

program, err := expr.Compile(`IsMore(City.Population)`, expr.Env(&EnvContextTwo{}))

output, err := expr.Run(program, &EnvContextTwo{...})
```

## Marshaling program

Compiled program is possible to marshal and unmarshal before running.

```go
    env := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	program, err := expr.Compile("foo + bar", expr.Env(env))
	b, err := json.Marshal(program)

	unmarshaledProgram := &vm.Program{}
	err = json.Unmarshal(b, unmarshaledProgram)
	
	output, err := expr.Run(unmarshaledProgram, env)

	fmt.Printf("%v", output) // outputs 3
```

## Visitor

[ast](https://godoc.org/gopkg.in/antonmedv/expr.v2/ast) package provides `Visitor` interface and `BaseVisitor` implementation. 
You can use it for traveling ast tree of compiled program.

For example if you want to collect all variable names:

```go
import "gopkg.in/antonmedv/expr.v2/ast"

type visitor struct {
	ast.BaseVisitor
	identifiers []string
}

func (v *visitor) IdentifierNode(node *ast.IdentifierNode) {
	v.identifiers = append(v.identifiers, node.Value)
}

program, err := expr.Compile("foo + bar", expr.Env(env))

visitor := &visitor{}
ast.Walk(node, visitor)
	
fmt.Printf("%v", visitor.identifiers) // outputs [foo bar]

```
