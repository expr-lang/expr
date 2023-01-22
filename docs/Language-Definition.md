# Language Definition

<table>
  <tr>
    <th colspan="2">Built-in Functions</th>
    <th colspan="2">Operators</th>
  </tr>
  <tr>
    <td>
      <a href="#allarray-predicate">all()</a><br>
      <a href="#anyarray-predicate">any()</a><br>
      <a href="#lenarray-predicate">one()</a><br>
      <a href="#nonearray-predicate">none()</a><br>
    </td>
    <td>
      <a href="#lenv">len()</a><br>
      <a href="#maparray-predicate">map()</a><br>
      <a href="#filterarray-predicate">filter()</a><br>
      <a href="#countarray-predicate">count()</a><br>    
    </td>
    <td>
      <a href="#string-operators">matches</a><br>
      <a href="#string-operators">contains</a><br>
      <a href="#string-operators">startsWith</a><br>
      <a href="#string-operators">endsWith</a><br>    
    </td>
    <td>
      <a href="#membership-operators">in</a><br>
      <a href="#membership-operators">not in</a><br>
      <a href="#range-operator">x..y</a><br>
      <a href="#slice-operator">[x:y]</a><br>
    </td>
  </tr>
</table>

## Literals

* `true`
* `false`
* `nil`

### Strings

Single or double quotes. Unicode sequences (`\uXXXX`) are supported.

### Numbers

Integers and floats.

* `42`
* `3.14`
* `1e6`
* `0x2A`
* `1_000_000`

### Arrays

* `[1, 2, 3]`

Tailing commas are allowed.

### Maps

* `{foo: "bar"}`

Tailing commas are allowed.

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

* `.` (dot)
* `in` (contain)
* `not in` (does not contain)

Struct fields and map elements can be accessed by using the `.` or the `[]`
syntax.

Example:

```
user.Group in ["human_resources", "marketing"]
```

```
data["tag-name"] in {foo: 1, bar: 2}
```

### Range Operator

* `..` (range)

Example:

```
user.Age in 18..45
```

The range is inclusive:

```
1..3 == [1, 2, 3]
```

### Slice Operator

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

### Ternary Operator

* `foo ? 'yes' : 'no'`

Example:

```
user.Age > 30 ? "mature" : "immature"
```

## Built-in Functions

### `all(array, predicate)`

Returns **true** if all elements satisfies the [predicate](#predicate).
If the array is empty, returns **true**.

```
all(Tweets, {.Size < 280})
```

### `any(array, predicate)`

Returns **true** if any elements satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

### `one(array, predicate)`

Returns **true** if _exactly one_ element satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

```
one(Participants, {.Winner})
```

### `none(array, predicate)`

Returns **true** if _all elements does not_ satisfy the [predicate](#predicate).
If the array is empty, returns **true**.

### `len(v)`

Returns the length of an array, a map or a string.

### `map(array, predicate)`

Returns new array by applying the [predicate](#predicate) to each element of
the array.

### `filter(array, predicate)`

Returns new array by filtering elements of the array by [predicate](#predicate).

### `count(array, predicate)`

Returns the number of elements what satisfies the [predicate](#predicate).
Equivalent to:

```
len(filter(array, predicate))
```

## Predicate

The predicate is an expression that accepts a single argument. To access
the argument use the `#` symbol.

```
map(0..9, {# / 2})
```

If items of the array is a struct or a map, it is possible to access fields with
omitted `#` symbol (`#.Value` becomes `.Value`).

```
filter(Tweets, {len(.Value) > 280})
```
