# Expr [![Build Status](https://travis-ci.org/antonmedv/expr.svg?branch=master)](https://travis-ci.org/antonmedv/expr) [![Code Coverage](https://scrutinizer-ci.com/g/antonmedv/expr/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/antonmedv/expr/?branch=master) 

Expr is an engine that can evaluate expressions. 

The purpose of the package is to allow users to use expressions inside configuration for more complex logic. 
It is a perfect candidate for the foundation of a _business rule engine_. 
The idea is to let configure things in a dynamic way without recompile of a program:

```ruby
# Get the special price if
user.Group() in ["good_customers", "collaborator"]

# Promote article to the homepage when
article.CommentCount > 100 and article.Category not in ["misc"]

# Send an alert when
product.Stock < 15
```

Inspired by 
* Symfony's [The ExpressionLanguage](https://github.com/symfony/expression-language) component,
* Rob Pike's talk [Lexical Scanning in Go](https://talks.golang.org/2011/lex.slide).

## Install

```
go get -u github.com/antonmedv/expr
```

## Documentation

### Usage
```go
// Evaluate expression on data.
result, err := expr.Eval("expression", data)

// Or precompile expression to ast first.
node, err := expr.Parse("expression")

// And run later.
result, err := expr.Run(node, data)
```

### Expression Syntax
See [The Expression Syntax](https://github.com/antonmedv/expr/wiki/The-Expression-Syntax) to learn the syntax of the Expr expressions.

### Passing in Variables
You can pass variables into the expression, which can be of any valid Go type (including structs):
```go
// Maps
data := map[string]interface{}{
    "Foo": ...
    "Bar": ...
}

// Structs
data := Payload{
	Foo: ...
	Bar: ...
}

// Pass object
result, err := expr.Eval("Foo == Bar", data)
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
    User user
}

req := Request{User{
    Cookies:   []cookie{{"origin", "www"}},
    UserAgent: "Firefox",
}}

ok, err := expr.Eval(`User.UserAgent matches "Firefox" and User.Cookies[0].Value == "www"`, req)
``` 

### Passing in Functions
You can also pass functions into the expression:
```go
data := map[string]interface{}{
    "Request": req,
    "Values": func(xs []Cookie) []string {
        vs := make([]string, 0)
        for _, x := range xs {
            vs = append(vs, x.Value)
        }
        return vs
    },
}

ok, err := expr.Eval(`"www" in Values(Request.User.Cookies)`, data)
```

### Caching
If you planning to execute some expression lots times, it's good to compile it first and only one time: 

```go
// Precompile
node, err := expr.Parse(expression)

// Run
ok, err := expr.Run(node, data)
```

### Checking variables and functions
It is possible to check used variables and functions during parsing of the expression.

```go
expression := `Request.User.UserAgent matches "Firefox" && "www" in Values(Request.User.Cookies)`

node, err := expr.Parse(expression, expr.Names("Request"), expr.Funcs("Values"))
```

Only `Request` and `Values` will bbe available inside expression, otherwise parse error.

### Printing
Compiled ast can be compiled back to string expression using _String()_:

```go
node, err := expr.Parse(expression)
code := fmt.Sprintf("%v", node)
``` 

### Number type
Inside Expr engine there is no distinguish between int, uint and float types (as in JavaScript).
All numbers inside Expr engine represented as `float64`. 
You should remember about it if you use any of binary operators (`+`, `-`, `/`, `*`, etc).
Otherwise type remain unchanged.

```go
data := map[string]int{
    "Foo": 1,
    "Bar": 2,
}

out, err := expr.Eval(`Foo`, data) // int

out, err := expr.Eval(`Foo + Bar`, data) // float64
```

## License

MIT
