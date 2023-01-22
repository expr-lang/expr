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

## Supported Literals

The package supports:

* **strings** - single and double quotes (e.g. `"hello"`, `'hello'`)
* **numbers** - e.g. `103`, `2.5`, `.5`
* **arrays** - e.g. `[1, 2, 3]`
* **maps** - e.g. `{foo: "bar"}`
* **booleans** - `true` and `false`
* **nil** - `nil`

## Digit separators

Integer literals may contain digit separators to allow digit grouping into more
legible forms.

Example:

```js
10_000_000_000
```

## Fields

Struct fields and map elements can be accessed by using the `.` or the `[]`
syntax.

```js
foo.Field
bar["some-key"]
```

## Functions

Functions may be called using the `()` syntax.

```js
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

```js
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

```js
life < universe || life < everything
```

### String Operators

* `+` (concatenation)
* `matches` (regex match)
* `contains` (string contains)
* `startsWith` (has prefix)
* `endsWith` (has suffix)

Example:

```js
"hello" matches "h.*"
```

### Membership Operators

* `in` (contain)
* `not in` (does not contain)

Example:

```js
user.Group in ["human_resources", "marketing"]
```

```js
"foo" in {foo: 1, bar: 2}
```

### Range Operator

* `..` (range)

Example:

```js
user.Age in 18..45
```

The range is inclusive:

```js
1..3 == [1, 2, 3]
```

### Slice Operator

* `array[:]` (slice)

Slices can work with arrays or strings.

Example:

Variable `array` is `[1,2,3,4,5]`.

```js
array[1:4] == [2,3,4]
array[:3] == [1,2,3]
array[3:] == [4,5]
array[:] == array
```

### Ternary Operator

* `foo ? 'yes' : 'no'`

Example:

```js
user.Age > 30 ? "mature" : "immature"
```

## Built-in Functions

### `all(array, predicate)`

Returns **true** if all elements satisfies the [predicate](#predicate).
If the array is empty, returns **true**.

```js
all(Tweets, {.Size < 280})
```

### `any(array, predicate)`

Returns **true** if any elements satisfies the [predicate](#predicate).
If the array is empty, returns **false**.


### `one(array, predicate)`

Returns **true** if _exactly one_ element satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

```js
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

```js
len(filter(array, predicate))
```

## Predicate

The predicate is an expression that accepts a single argument. To access 
the argument use the `#` symbol.

```js
map(0..9, {# / 2})
```

If items of the array is a struct or a map, it is possible to access fields with 
omitted `#` symbol (`#.Value` becomes `.Value`).

```js
filter(Tweets, {len(.Value) > 280})
```
