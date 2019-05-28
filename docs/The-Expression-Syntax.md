The Expr package uses a specific syntax. In this document, you can find all supported
syntaxes.

## Supported Literals

The package supports:

* **strings** - single and double quotes (e.g. `"hello"`, `'hello'`)
* **numbers** - e.g. `103`, `2.5`
* **arrays** - using JSON-like notation (e.g. `[1, 2]`)
* **maps** - using JSON-like notation (e.g. `{ foo: "bar" }`)
* **booleans** - `true` and `false`
* **nil** - `nil`

## Working with Structs

When passing structs or maps into an expression, you can use different syntaxes to
access properties.

Public properties on structs can be accessed by using the `.` syntax:

```
Foo.Bar.Baz
```

## Working with Functions

You can also use pass functions into the expression and call them using C syntax.

```
Upper("text")
```

Result will be set to `HELLO`. 

Note: `Upper` function doesn't present in package by default. 


## Working with Arrays

If you pass an array into an expression, use the `[]` syntax to access values:

```
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

```
life + universe + everything
```

> Note what result will be `float64` as all binary operators what works with numbers, cast to `float64` automatically. 

### Bitwise Operators

* `&` (and)
* `|` (or)
* `^` (xor)

### Comparison Operators

* `==` (equal)
* `!=` (not equal)
* `<` (less than)
* `>` (greater than)
* `<=` (less than or equal to)
* `>=` (greater than or equal to)
* `matches` (regex match)

To test if a string does *not* match a regex, use the logical `not` operator in combination with the `matches` operator:

```
not ("foo" matches "bar")
```

You must use parenthesis because the unary operator `not` has precedence over the binary operator `matches`.

### Logical Operators

* `not` or `!`
* `and` or `&&`
* `or` or `||`

Example:

```
life < universe or life < everything
```

### String Operators

* `~` (concatenation)

Example:

```go
'Arthur' ~ ' ' ~ 'Dent'
```

Result will be set to `Arthur Dent`.

### Membership Operators

* `in` (contain)
* `not in` (does not contain)

Example:

```
User.Group in ["human_resources", "marketing"]
```

```
"foo" in {foo: 1, bar: 2}
```

### Numeric Operators

* `..` (range)

Example:

```
User.Age in 18..45
```

The range is inclusive:

```
1..3 == [1, 2, 3]
```

### Ternary Operators

* `foo ? 'yes' : 'no'`
* `foo ?: 'no'` (equal to `foo ? foo : 'no'`)

## Builtin functions

* `len` (length of array or string)

Example:

```
len([1, 2, 3]) == len("foo")
```

Result will be set to `true`.
