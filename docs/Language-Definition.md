# Language Definition

**Expr** is an expression evaluation language for Go.

## Supported Literals

The package supports:

* **strings** - single and double quotes (e.g. `"hello"`, `'hello'`)
* **numbers** - e.g. `103`, `2.5`, `.5`
* **arrays** - e.g. `[1, 2, 3]`
* **maps** - e.g. `{foo: "bar"}`
* **booleans** - `true` and `false`
* **nil** - `nil`

## Digit separators

Integer literals may contain digit separators to allow digit grouping into more legible forms.

Example:

```
10_000_000_000
```

## Fields

Struct fields and map elements can be accessed by using the `.` or the `[]` syntax.

```
foo.Field
bar["some-key"]
```

## Functions

Functions may be called using the `()` syntax.

```
foo.Method()
```

## Operators

### Arithmetic Operators

* `+` (addition)
* `-` (subtraction)
* `*` (multiplication)
* `/` (division)
* `%` (modulus)
* `^` or `**` (exponent)

Example:

```
x^2 + y^2
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
life < universe || life < everything
```

### String Operators

* `+` (concatenation)
* `matches` (regex match)
* `contains` (string contains)
* `startsWith` (has prefix)
* `endsWith` (has suffix)

Example:

```
"hello" matches "h.*"
```

### Membership Operators

* `in` (contain)
* `not in` (does not contain)

Example:

```
user.Group in ["human_resources", "marketing"]
```

```
"foo" in {foo: 1, bar: 2}
```

### Numeric Operators

* `..` (range)

Example:

```
user.Age in 18..45
```

The range is inclusive:

```
1..3 == [1, 2, 3]
```

### Ternary Operators

* `foo ? 'yes' : 'no'`

Example:

```
user.Age > 30 ? "mature" : "immature"
```

## Builtin functions

* `len` (length of array, map or string)
* `all` (will return `true` if all element satisfies the predicate)
* `none` (will return `true` if all element does NOT satisfy the predicate)
* `any` (will return `true` if any element satisfies the predicate)
* `one` (will return `true` if exactly ONE element satisfies the predicate)
* `filter` (filter array by the predicate)
* `map` (map all items with the closure)
* `count` (returns number of elements what satisfies the predicate)

Examples:

Ensure all tweets are less than 280 chars.

```
all(Tweets, {.Size < 280})
```

Ensure there is exactly one winner.

```
one(Participants, {.Winner})
```

## Closures

The closure is an expression that accepts a single argument. To access 
the argument use the `#` symbol.

```
map(0..9, {# / 2})
```

If the item of array is struct, it is possible to access fields of struct with 
omitted `#` symbol (`#.Value` becomes `.Value`).

```
filter(Tweets, {len(.Value) > 280})
```

## Slices

* `array[:]` (slice)

Slices can work with arrays or strings.

Example:

Variable `array` is `[1,2,3,4,5]`.

```
array[1:5] == [2,3,4] 
array[3:] == [4,5]
array[:4] == [1,2,3]
array[:] == array
```
