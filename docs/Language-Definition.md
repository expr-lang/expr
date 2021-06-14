# Language Definition

**Expr** package uses a specific syntax. In this document, you can find all supported
syntaxes.

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

## Accessing Public Properties

Public properties on structs can be accessed by using the `.` syntax. 
If you pass an array into an expression, use the `[]` syntax to access array keys.

```js
foo.Array[0].Value
```

## Functions and Methods

Functions may be called using `()` syntax. The `.` syntax can also be used to call methods on an struct.

```js
price.String()
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

```js
life + universe + everything
``` 

#### Broadcasting and Array Operators

The term broadcasting describes how expr treats arrays and scalars during arithmetic operations. Subject to certain constraints, the scalar is broadcast across the larger array so that they have compatible shapes.

Example:

```js
2 * [3, 4]
``` 

Output:

```js
[6, 8]
``` 

Example:

```js
[-1, 1] * [3, 4]
``` 

Output:

```js
[-3, 4]
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

To test if a string does *not* match a regex, use the logical `not` operator in combination with the `matches` operator:

```js
not ("foo" matches "^b.+")
```

You must use parenthesis because the unary operator `not` has precedence over the binary operator `matches`.

Example:

```js
'Arthur' + ' ' + 'Dent'
```

Result will be set to `Arthur Dent`.

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

### Numeric Operators

* `..` (range)

Example:

```js
user.Age in 18..45
```

The range is inclusive:

```js
1..3 == [1, 2, 3]
```

### Ternary Operators

* `foo ? 'yes' : 'no'`

Example:

```js
user.Age > 30 ? "mature" : "immature"
```

## Builtin functions

* `len` (length of array, map or string)
* `all` (will return `true` if all element satisfies the predicate)
* `none` (will return `true` if all element does NOT satisfies the predicate)
* `any` (will return `true` if any element satisfies the predicate)
* `one` (will return `true` if exactly ONE element satisfies the predicate)
* `filter` (filter array by the predicate)
* `map` (map all items with the closure)
* `count` (returns number of elements what satisfies the predicate)

Examples:

Ensure all tweets are less than 280 chars.

```js
all(Tweets, {.Size < 280})
```

Ensure there is exactly one winner.

```js
one(Participants, {.Winner})
```

## Vector Math functions

* `abs` (returns the absolute value of x, for each element x in the input array)
* `acos` (returns the arccosine, in radians, of x, for each element x in the input array)
* `acosh` (returns the inverse hyperbolic cosine of x, for each element x in the input array)
* `asin` (returns the arcsine, in radians, of x, for each element x in the input array)
* `asinh` (returns the inverse hyperbolic sine of x, for each element x in the input array)
* `atan` (returns the arctangent, in radians, of x, for each element x in the input array)
* `atanh` (returns the inverse hyperbolic tangent of x, for each element x in the input array)
* `cbrt` (returns the cube root of x, for each element x in the input array)
* `ceil` (returns the least integer value greater than or equal to x, for each element x in the input array)
* `cos` (returns the cosine of the radian element x, for each element x in the input array)
* `cosh` (returns the hyperbolic cosine of x, for each element x in the input array)
* `erf` (returns the error function of x, for each element x in the input array)
* `erfc` (returns the complementary error function of x, for each element x in the input array)
* `erfcinv` (returns the inverse of Erfc(x), for each element x in the input array)
* `erfinv` (returns the inverse error function of x, for each element x in the input array)
* `exp` (returns `e**x`, the base-e exponential of x, for each element x in the input array)
* `exp2` (returns `2**x`, the base-2 exponential of x, for each element x in the input array)
* `expm1` (returns `e**x - 1`, the base-e exponential of x minus 1, for each element x in the input array)
* `floor` (returns the greatest integer value less than or equal to x, for each element x in the input array)
* `gamma` (returns the Gamma function of x, for each element x in the input array)
* `j0` (returns the order-zero Bessel function of the first kind, for each element x in the input array)
* `j1` (returns the order-one Bessel function of the first kind, for each element x in the input array)
* `log` (returns the natural logarithm of x, for each element x in the input array)
* `log10` (returns the decimal logarithm of x, for each element x in the input array)
* `log1p` (returns the natural logarithm of 1 plus its argument x, for each element x in the input array)
* `log2` (returns the binary logarithm of x, for each element x in the input array)
* `logb` (returns the binary exponent of x, for each element x in the input array)
* `round` (returns the nearest integer, rounding half away from zero, for each element x in the input array)
* `roundtoeven` (returns the nearest integer, rounding ties to even, for each element x in the input array)
* `sin` (returns the sine of the radian argument x, for each element x in the input array)
* `sinh` (returns the hyperbolic sine of x, for each element x in the input array)
* `sqrt` (returns the square root of x, for each element x in the input array)
* `tan` (returns the tangent of the radian argument x, for each element x in the input array)
* `tanh` (returns the hyperbolic tangent of x, for each element x in the input array)
* `trunc` (returns the integer value of x, for each element x in the input array)
* `y0` (returns the order-zero Bessel function of the second kind, for each element x in the input array)
* `y1` (returns the order-one Bessel function of the second kind, for each element x in the input array)
* `maximum` (returns the larger of x or y, for each parallel element x and y in the input arrays)
* `minimum` (returns the smaller of x or y, for each parallel element x and y in the input arrays)
* `mod` (returns the floating-point remainder of x/y, for each parallel element x and y in the input arrays)
* `pow` (returns `x**y`, the base-x exponential of y, for each parallel element x and y in the input arrays)
* `remainder` (returns the IEEE 754 floating-point remainder of x/y, for each parallel element x and y in the input arrays)
* `nanmin` (returns the minimum of an array, ignoring any NaN)
* `nanmax` (returns the maximum of an array, ignoring any NaN)
* `nanmean` (returns the mean of an array, ignoring any NaN)
* `nanstd` (returns the standard deviation of an array, ignoring any NaN)
* `nansum` (returns the sum of an array, ignoring any NaN)
* `nanprod` (returns the product of array elements, ignoring any NaN)

## Closures

* `{...}` (closure)

Closures allowed only with builtin functions. To access current item use `#` symbol.

```js
map(0..9, {# / 2})
```

If the item of array is struct, it's possible to access fields of struct with omitted `#` symbol (`#.Value` becomes `.Value`).

```js
filter(Tweets, {len(.Value) > 280})
```

## Slices

* `array[:]` (slice)

Slices can work with arrays or strings.

Example:

Variable `array` is `[1,2,3,4,5]`.

```js
array[1:5] == [2,3,4] 
array[3:] == [4,5]
array[:4] == [1,2,3]
array[:] == array
```
