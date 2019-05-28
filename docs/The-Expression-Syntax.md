# The Expression Syntax

**Expr** package uses a specific syntax. In this document, you can find all supported
syntaxes.

## Supported Literals

The package supports:

* **strings** - single and double quotes (e.g. `"hello"`, `'hello'`)
* **numbers** - e.g. `103`, `2.5`, `1e+6`
* **arrays** - e.g. `[1, 2]`
* **maps** - e.g. `{foo: "bar"}`
* **booleans** - `true` and `false`
* **nil** - `nil`

## Working with Structs

When passing structs or maps into an expression, you can use different syntaxes to
access properties.

Public properties on structs can be accessed by using the `.` syntax:

```coffeescript
Foo.Bar.Baz
```

## Working with Functions

You can also use functions by calling them using C syntax.

```coffeescript
Upper("text")
```

Result will be set to `HELLO`. 

> Note: `Upper` function doesn't present in package by default. 

## Working with Arrays

Use the `[]` syntax to access values:

```coffeescript
array[2]
```

## Supported Operators

The package comes with a lot of operators:

### Arithmetic Operators

* `+` (addition)
* `-` (subtraction)
* `*` (multiplication)
* `/` (division)
* `%` (modulus)
* `**` (pow)

Example:

```coffeescript
life + universe + everything
``` 

### Comparison Operators

* `==` (equal)
* `!=` (not equal)
* `<` (less than)
* `>` (greater than)
* `<=` (less than or equal to)
* `>=` (greater than or equal to)

### Logical Operators

* `not` or `!`
* `and` or `&&`
* `or` or `||`

Example:

```
life < universe or life < everything
```

### String Operators

* `+` (concatenation)
* `matches` (regex match)
* `contains` (string contains)

To test if a string does *not* match a regex, use the logical `not` operator in combination with the `matches` operator:

```coffeescript
not ("foo" matches "^b.+")
```

You must use parenthesis because the unary operator `not` has precedence over the binary operator `matches`.

Example:

```coffeescript
'Arthur' + ' ' + 'Dent'
```

Result will be set to `Arthur Dent`.

### Membership Operators

* `in` (contain)
* `not in` (does not contain)

Example:

```coffeescript
User.Group in ["human_resources", "marketing"]
```

```coffeescript
"foo" in {foo: 1, bar: 2}
```

### Numeric Operators

* `..` (range)

Example:

```coffeescript
User.Age in 18..45
```

The range is inclusive:

```coffeescript
1..3 == [1, 2, 3]
```

### Ternary Operators

* `foo ? 'yes' : 'no'`

## Builtin functions

* `len` (length of array or string)
* `all` (will return `true` if all element satisfies the predicate)
* `none` (will return `true` if all element does NOT satisfies the predicate)
* `any` (will return `true` if any element satisfies the predicate)
* `one` (will return `true` if exactly ONE element satisfies the predicate)
* `filter` (filter array by the predicate)
* `map` (map all items with the closure)

Example:

```go
// Ensure all tweets are less than 140 chars.
all(Tweets, {.Size < 140})
```

## Closure

* `{...}` (closure)

**Expr** support closures with builtin functions. To access current item use `#` symbol.

```go
map(0..9, {# + 1})
```

If the item of array is struct, it's possible to access fields of struct with omitted `#` symbol (`#.Value` becomes `.Value`).

```go
filter(Tweets, {.Size > 140})
```
