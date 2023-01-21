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

## Built-in Functions

<table>
<tr>
  <td>
    <a href="#all">all()</a><br>
    <a href="#any">any()</a><br>
    <a href="#len">one()</a><br>
    <a href="#none">none()</a><br>
  </td>
  <td>
    <a href="#len">len()</a><br>
    <a href="#map">map()</a><br>
    <a href="#filter">filter()</a><br>
    <a href="#count">count()</a><br>    
  </td>
</tr>
</table>


#### `all(array, predicate)`

Returns **true** if all elements satisfies the predicate (or if the array is empty).

```
all(Tweets, {.Size < 280})
```

#### `any(array, predicate)`

Returns **true** if any elements satisfies the predicate. If the array is empty, returns **false**.


#### `one(array, predicate)`

Returns **true** if _exactly one_ element satisfies the predicate. If the array is empty, returns **false**.

```
one(Participants, {.Winner})
```

#### `none(array, predicate)`

Returns **true** if _all elements does not_ satisfy the predicate. If the array is empty, returns **true**.

#### `len(v)`

Returns the length of an array, a map or a string.

#### `map(array, closure)`

Returns new array by applying the closure to each element of the array.

#### `filter(array, predicate)`

Returns new array by filtering elements of the array by predicate.

#### `count(array, predicate)`

Returns the number of elements what satisfies the predicate. Equivalent to:

```
len(filter(array, predicate))
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
array[1:4] == [2,3,4]
array[:3] == [1,2,3]
array[3:] == [4,5]
array[:] == array
```
