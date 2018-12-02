# Expr [![Build Status](https://travis-ci.org/antonmedv/expr.svg?branch=master)](https://travis-ci.org/antonmedv/expr) [![Go Report Card](https://goreportcard.com/badge/github.com/antonmedv/expr)](https://goreportcard.com/report/github.com/antonmedv/expr) [![Code Coverage](https://scrutinizer-ci.com/g/antonmedv/expr/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/antonmedv/expr/?branch=master) <a href="https://stars.medv.io/antonmedv/expr"><img src="https://stars.medv.io/antonmedv/expr.svg" alt="Sparkline" height="24"></a>

Expr is an engine that can evaluate expressions. 

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

Inspired by 
* Symfony's [The ExpressionLanguage](https://github.com/symfony/expression-language) component,
* Rob Pike's talk [Lexical Scanning in Go](https://talks.golang.org/2011/lex.slide).

## Features

* Works with any valid Go object (structs, maps, etc)
* Static and dynamic typing ([example](https://godoc.org/github.com/antonmedv/expr#example-Define))
  ```go
  code := "groups[0].Title + user.Age"
  p, err := expr.Parse(code, expr.Define("groups", []Group{}), expr.Define("user", User{}))
  // err: invalid operation: groups[0].Name + user.Age (mismatched types string and int)
  ```
* User-friendly error messages
  ```
  unclosed "("
  (boo + bar]
  ----------^
  ```
* Reasonable set of basic operators
* Fast (faster otto and goja, see [bench](https://github.com/antonmedv/expr/wiki/Benchmarks))

## Install

```
go get -u github.com/antonmedv/expr
```

<a href="https://www.patreon.com/antonmedv">
	<img src="https://c5.patreon.com/external/logo/become_a_patron_button@2x.png" width="160">
</a>

## Documentation

* See [![GoDoc](https://godoc.org/github.com/antonmedv/expr?status.svg)](https://godoc.org/github.com/antonmedv/expr) for developer documentation,
* See [The Expression Syntax](https://github.com/antonmedv/expr/wiki/The-Expression-Syntax) page to learn the syntax of the Expr expressions.

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
type env struct {
	Foo int
	Bar bar
}

type bar struct {
	Value int
}

p, err := expr.Parse("Foo + Bar.Value", expr.Env(env{}))

out, err := expr.Run(p, env{1, bar{2}})
```

Using env's methods as functions inside expressions.

```go
type env struct {
	Name string
}

func (e env) Title() string {
	return strings.Title(e.Name)
}


p, err := expr.Parse("'Hello ' ~ Title()", expr.Env(env{}))

out, err := expr.Run(p, env{"world"})
```

Using embedded structs to construct env.

```go
type env struct {
	helpers
	Name string
}

type helpers struct{}

func (h helpers) Title(s string) string {
	return strings.Title(s)
}


p, err := expr.Parse("'Hello ' ~ Title(Name)", expr.Env(env{}))

out, err := expr.Run(p, env{"world"})
```

## License

MIT
