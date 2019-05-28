# Expr 
[![Build Status](https://travis-ci.org/antonmedv/expr.svg?branch=master)](https://travis-ci.org/antonmedv/expr) 
[![Go Report Card](https://goreportcard.com/badge/github.com/antonmedv/expr)](https://goreportcard.com/report/github.com/antonmedv/expr) 
[![Code Coverage](https://scrutinizer-ci.com/g/antonmedv/expr/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/antonmedv/expr/?branch=master) 
[![GoDoc](https://godoc.org/github.com/antonmedv/expr?status.svg)](https://godoc.org/github.com/antonmedv/expr)

**Expr** package provides an engine that can compile and evaluate expressions. 
An expression is a one-liner that returns a value (mostly, but not limited to, booleans).
It is designed for simplicity, speed and safety.

The purpose of the package is to allow users to use expressions inside configuration for more complex logic. 
It is a perfect candidate for the foundation of a _business rule engine_. 
The idea is to let configure things in a dynamic way without recompile of a program:

```coffeescript
# Get the special price if
user.Group in ["good_customers", "collaborator"]

# Promote article to the homepage when
len(article.Comments) > 100 and article.Category not in ["misc"]

# Send an alert when
product.Stock < 15
```

## Features

* Seamless integration with Go.
* Static typing ([example](https://godoc.org/github.com/antonmedv/expr#example-Env)).
  ```go
  out, err := expr.Eval("'hello' + 10")
  // err: invalid operation + (mismatched types string and int64)
  // | 'hello' + 10
  // | ........^
  ```
* User-friendly error messages.
* Reasonable set of basic operators.
* Builtins `all`, `none`, `any`, `one`, `filter`, `map`.
  ```coffeescript
  all(Tweets, {.Size < 140})
  ```
* Fast ([benchmarks](https://github.com/antonmedv/golang-expression-evaluation-comparison)).

## Install

```
go get github.com/antonmedv/expr
```

<a href="https://www.patreon.com/antonmedv">
	<img src="https://c5.patreon.com/external/logo/become_a_patron_button@2x.png" width="160">
</a>

## Documentation

* See [docs](docs) page for developer documentation.
* See [The Expression Syntax](docs/The-Expression-Syntax.md) page to learn the syntax.

## Examples

Executing arbitrary expressions.

```go
env := map[string]interface{}{
    "foo": 1,
    "bar": struct{Value int}{1},
}

out, err := expr.Eval("foo + bar.Value", env)
```

Static type checker with struct as environment.

```go
type Env struct {
	Foo int
	Bar bar
}

type Bar struct {
	Value int
}

program, err := expr.Compile("Foo + Bar.Value", expr.Env(Env{}))

out, err := expr.Run(program, Env{1, Bar{2}})
```

Using env's methods as functions inside expressions.

```go
type Env struct {
	Name string
}

func (e *Env) Title() string {
	return strings.Title(e.Name)
}

program, err := expr.Compile(`"Hello " + Title()`, expr.Env(&Env{}))

out, err := expr.Run(program, &Env{"world"})
```

Using embedded structs to construct env.

```go
type Env struct {
	Helpers
	Name string
}

type Helpers struct{}

func (Helpers) Title(s string) string {
	return strings.Title(s)
}


program, err := expr.Compile(`"Hello " + Title(Name)`, expr.Env(Env{}))

out, err := expr.Run(program, Env{"world"})
```

## Who is using Expr?

* [Aviasales](https://aviasales.ru) are actively using Expr for different parts of the search engine.

## License

[MIT](LICENSE)
