# Language Definition

## Literals

<table>
    <tr>
        <td>Comment</td>
        <td>
             <code>/* */</code> or <code>//</code>
        </td>
    </tr>
    <tr>
        <td>Boolean</td>
        <td>
            <code>true</code>, <code>false</code>
        </td>
    </tr>
    <tr>
        <td>Integer</td>
        <td>
            <code>42</code>, <code>0x2A</code>
        </td>
    </tr>
    <tr>
        <td>Float</td>
        <td>
            <code>0.5</code>, <code>.5</code>
        </td>
    </tr>
    <tr>
        <td>String</td>
        <td>
            <code>"foo"</code>, <code>'bar'</code>
        </td>
    </tr>
    <tr>
        <td>Duration</td>
        <td>
            <code>1h16m7ms</code>
        </td>
    </tr>
    <tr>
        <td>Array</td>
        <td>
            <code>[1, 2, 3]</code>
        </td>
    </tr>
    <tr>
        <td>Map</td>
        <td>
            <code>&#123;a: 1, b: 2, c: 3&#125;</code>
        </td>
    </tr>
    <tr>
        <td>Nil</td>
        <td>
            <code>nil</code>
        </td>
    </tr>
</table>

## Operators

<table>
    <tr>
        <td>Arithmetic</td>
        <td>
            <code>+</code>, <code>-</code>, <code>*</code>, <code>/</code>, <code>%</code> (modulus), <code>^</code> or <code>**</code> (exponent)
        </td>
    </tr>
    <tr>
        <td>Comparison</td>
        <td>
            <code>==</code>, <code>!=</code>, <code>&lt;</code>, <code>&gt;</code>, <code>&lt;=</code>, <code>&gt;=</code>
        </td>
    </tr>
    <tr>
        <td>Logical</td>
        <td>
            <code>not</code> or <code>!</code>, <code>and</code> or <code>&amp;&amp;</code>, <code>or</code> or <code>||</code>
        </td>
    </tr>
    <tr>
        <td>Conditional</td>
        <td>
            <code>?:</code> (ternary), <code>??</code> (nil coalescing)
        </td>
    </tr>
    <tr>
        <td>Membership</td>
        <td>
            <code>[]</code>, <code>.</code>, <code>?.</code>, <code>in</code>
        </td>
    </tr>
    <tr>
        <td>String</td>
        <td>
            <code>+</code> (concatenation), <code>contains</code>, <code>startsWith</code>, <code>endsWith</code>
        </td>
    </tr>
    <tr>
        <td>Regex</td>
        <td>
            <code>matches</code>
        </td>
    </tr>
    <tr>
        <td>Range</td>
        <td>
            <code>..</code>
        </td>
    </tr>
    <tr>
        <td>Slice</td>
        <td>
            <code>[:]</code>
        </td>
    </tr>
    <tr>
        <td>Pipe</td>
        <td>
            <code>|</code>
        </td>
    </tr>
</table>

Examples:

```expr
user.Age in 18..45 and user.Name not in ["admin", "root"]
```

```expr
foo matches "^[A-Z].*"
```

```expr
tweets | filter(.Size < 280) | map(.Content) | join(" -- ")
```

```expr
filter(posts, {now() - .CreatedAt >= 7 * duration("24h")})
```

### Membership Operator

Fields of structs and items of maps can be accessed with `.` operator
or `[]` operator. Elements of arrays and slices can be accessed with
`[]` operator. Negative indices are supported with `-1` being
the last element.

The `in` operator can be used to check if an item is in an array or a map.

```expr
user.Name in list["available-names"]
```

#### Optional chaining

The `?.` operator can be used to access a field of a struct or an item of a map
without checking if the struct or the map is `nil`. If the struct or the map is
`nil`, the result of the expression is `nil`.

```expr
author?.User?.Name
```

#### Nil coalescing

The `??` operator can be used to return the left-hand side if it is not `nil`,
otherwise the right-hand side is returned.

```expr
author?.User?.Name ?? "Anonymous"
```

### Slice Operator

The slice operator `[:]` can be used to access a slice of an array.

For example, variable `array` is `[1, 2, 3, 4, 5]`:

```expr
array[1:4] == [2, 3, 4]
array[1:-1] == [2, 3, 4]
array[:3] == [1, 2, 3]
array[3:] == [4, 5]
array[:] == array
```

### Pipe Operator

The pipe operator `|` can be used to pass the result of the left-hand side
expression as the first argument of the right-hand side expression.

For example, expression `split(lower(user.Name), " ")` can be written as:

```expr
user.Name | lower() | split(" ")
```

## String Functions

### trim(str[, chars])

Removes white spaces from both ends of a string `str`.
If the optional `chars` argument is given, it is a string specifying the set of characters to be removed.

```expr
trim("  Hello  ") == "Hello"
trim("__Hello__", "_") == "Hello"
```

### trimPrefix(str, prefix)

Removes the specified prefix from the string `str` if it starts with that prefix.

```expr
trimPrefix("HelloWorld", "Hello") == "World"
```

### trimSuffix(str, suffix)

Removes the specified suffix from the string `str` if it ends with that suffix.

```expr
trimSuffix("HelloWorld", "World") == "Hello"
```

### upper(str)

Converts all the characters in string `str` to uppercase.

```expr
upper("hello") == "HELLO"
```

### lower(str)

Converts all the characters in string `str` to lowercase.

```expr
lower("HELLO") == "hello"
```

### split(str, delimiter[, n])

Splits the string `str` at each instance of the delimiter and returns an array of substrings.

```expr
split("apple,orange,grape", ",") == ["apple", "orange", "grape"]
split("apple,orange,grape", ",", 2) == ["apple", "orange,grape"]
```

### splitAfter(str, delimiter[, n])

Splits the string `str` after each instance of the delimiter.

```expr
splitAfter("apple,orange,grape", ",") == ["apple,", "orange,", "grape"]
splitAfter("apple,orange,grape", ",", 2) == ["apple,", "orange,grape"]
```

### replace(str, old, new)

Replaces all occurrences of `old` in string `str` with `new`.

```expr
replace("Hello World", "World", "Universe") == "Hello Universe"
```

### repeat(str, n)

Repeats the string `str` `n` times.

```expr
repeat("Hi", 3) == "HiHiHi"
```

### indexOf(str, substring)

Returns the index of the first occurrence of the substring in string `str` or -1 if not found.

```expr
indexOf("apple pie", "pie") == 6
```

### lastIndexOf(str, substring)

Returns the index of the last occurrence of the substring in string `str` or -1 if not found.

```expr
lastIndexOf("apple pie apple", "apple") == 10
```

### hasPrefix(str, prefix)

Returns `true` if string `str` starts with the given prefix.

```expr
hasPrefix("HelloWorld", "Hello") == true
```

### hasSuffix(str, suffix)

Returns `true` if string `str` ends with the given suffix.

```expr
hasSuffix("HelloWorld", "World") == true
```

## Date Functions

The following operators can be used to manipulate dates:

```expr
date("2023-08-14") + duration("1h")
date("2023-08-14") - duration("1h")
date("2023-08-14") - date("2023-08-13") == duration("24h")
```

### now()

Returns the current date and time.

```expr
createdAt > now() - duration(1h)
```

### duration(str)

Returns [time.Duration](https://pkg.go.dev/time#Duration) value of the given string `str`.

Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

```expr
duration("1h").Seconds() == 3600
```

### date(str[, format[, timezone]])

Converts the given string `str` into a date representation.

If the optional `format` argument is given, it is a string specifying the format of the date.
The format string uses the same formatting rules as the standard
Go [time package](https://pkg.go.dev/time#pkg-constants).

If the optional `timezone` argument is given, it is a string specifying the timezone of the date.

If the `format` argument is not given, the `v` argument must be in one of the following formats:

- 2006-01-02
- 15:04:05
- 2006-01-02 15:04:05
- RFC3339
- RFC822,
- RFC850,
- RFC1123,

```expr
date("2023-08-14")
date("15:04:05")
date("2023-08-14T00:00:00Z")
date("2023-08-14 00:00:00", "2006-01-02 15:04:05", "Europe/Zurich")
```

## Duration Functions

The following operators can be used to manipulate durations:

```expr
-1h == -1 * 1h
+1h == 1h
2 * 1h == 2h
1h + 1m == 1h1m
1h - 1m == 59m
1h / 10m == 6
1h / 2 == 30m
```

Some number functions (max, min and abs) are compatible with durations as well.

### round(d1, d2)

Returns the result of rounding d1 to the nearest multiple of d2.

```expr
round(24h2m, 1h) == 24h
```

## Number Functions

### max(n1, n2)

Returns the maximum of the two numbers `n1` and `n2`.

```expr
max(5, 7) == 7
```

### min(n1, n2)

Returns the minimum of the two numbers `n1` and `n2`.

```expr
min(5, 7) == 5
```

### abs(n)

Returns the absolute value of a number.

### ceil(n)

Returns the least integer value greater than or equal to x.

```expr
ceil(1.5) == 2.0
```

### floor(n)

Returns the greatest integer value less than or equal to x.

```expr
floor(1.5) == 1.0
```

### round(n)

Returns the nearest integer, rounding half away from zero.

```expr
round(1.5) == 2.0
```

## Array Functions

### all(array, predicate)

Returns **true** if all elements satisfies the [predicate](#predicate).
If the array is empty, returns **true**.

```expr
all(tweets, {.Size < 280})
```

### any(array, predicate)

Returns **true** if any elements satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

### one(array, predicate)

Returns **true** if _exactly one_ element satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

```expr
one(participants, {.Winner})
```

### none(array, predicate)

Returns **true** if _all elements does not_ satisfy the [predicate](#predicate).
If the array is empty, returns **true**.

### map(array, predicate)

Returns new array by applying the [predicate](#predicate) to each element of
the array.

```expr
map(tweets, {.Size})
```

### filter(array, predicate)

Returns new array by filtering elements of the array by [predicate](#predicate).

```expr
filter(users, .Name startsWith "J")
```

### find(array, predicate)

Finds the first element in an array that satisfies the [predicate](#predicate).

```expr
find([1, 2, 3, 4], # > 2) == 3
```

### findIndex(array, predicate)

Finds the index of the first element in an array that satisfies the [predicate](#predicate).

```expr
findIndex([1, 2, 3, 4], # > 2) == 2
```

### findLast(array, predicate)

Finds the last element in an array that satisfies the [predicate](#predicate).

```expr
findLast([1, 2, 3, 4], # > 2) == 4
```

### findLastIndex(array, predicate)

Finds the index of the last element in an array that satisfies the [predicate](#predicate).

```expr
findLastIndex([1, 2, 3, 4], # > 2) == 3
```

### groupBy(array, predicate)

Groups the elements of an array by the result of the [predicate](#predicate).

```expr
groupBy(users, .Age)
```

### count(array, predicate)

Returns the number of elements what satisfies the [predicate](#predicate).

Equivalent to:

```expr
len(filter(array, predicate))
```

### join(array[, delimiter])

Joins an array of strings into a single string with the given delimiter.
If no delimiter is given, an empty string is used.

```expr
join(["apple", "orange", "grape"], ",") == "apple,orange,grape"
join(["apple", "orange", "grape"]) == "appleorangegrape"
```

### reduce(array, predicate[, initialValue])

Applies a predicate to each element in the array, reducing the array to a single value.
Optional `initialValue` argument can be used to specify the initial value of the accumulator.
If `initialValue` is not given, the first element of the array is used as the initial value.

Following variables are available in the predicate:

- `#` - the current element
- `#acc` - the accumulator
- `#index` - the index of the current element

```expr
reduce(1..9, #acc + #)
reduce(1..9, #acc + #, 0)
```

### sum(array)

Returns the sum of all numbers in the array.

```expr
sum([1, 2, 3]) == 6
```

### mean(array)

Returns the average of all numbers in the array.

```expr
mean([1, 2, 3]) == 2.0
```

### median(array)

Returns the median of all numbers in the array.

```expr
median([1, 2, 3]) == 2.0
```

### first(array)

Returns the first element from an array. If the array is empty, returns `nil`.

```expr
first([1, 2, 3]) == 1
```

### last(array)

Returns the last element from an array. If the array is empty, returns `nil`.

```expr
last([1, 2, 3]) == 3
```

### take(array, n)

Returns the first `n` elements from an array. If the array has fewer than `n` elements, returns the whole array.

```expr
take([1, 2, 3, 4], 2) == [1, 2]
```

### sort(array[, order])

Sorts an array in ascending order. Optional `order` argument can be used to specify the order of sorting: `asc`
or `desc`.

```expr
sort([3, 1, 4]) == [1, 3, 4]
sort([3, 1, 4], "desc") == [4, 3, 1]
```

### sortBy(array, key[, order])

Sorts an array of maps by a specific key in ascending order. Optional `order` argument can be used to specify the order
of sorting: `asc` or `desc`.

```expr
sortBy(users, "Age")
sortBy(users, "Age", "desc")
```

## Map Functions

### keys(map)

Returns an array containing the keys of the map.

```expr
keys({"name": "John", "age": 30}) == ["name", "age"]
```

### values(map)

Returns an array containing the values of the map.

```expr
values({"name": "John", "age": 30}) == ["John", 30]
```

## Type Conversion Functions

### type(v)

Returns the type of the given value `v`.
Returns on of the following types: `nil`, `bool`, `int`, `uint`, `float`, `string`, `array`, `map`.
For named types and structs, the type name is returned.

```expr
type(42) == "int"
type("hello") == "string"
type(now()) == "time.Time"
```

### int(v)

Returns the integer value of a number or a string.

```expr
int("123") == 123
```

### float(v)

Returns the float value of a number or a string.

### string(v)

Converts the given value `v` into a string representation.

```expr
string(123) == "123"
```

### toJSON(v)

Converts the given value `v` to its JSON string representation.

```expr
toJSON({"name": "John", "age": 30})
```

### fromJSON(v)

Parses the given JSON string `v` and returns the corresponding value.

```expr
fromJSON('{"name": "John", "age": 30}')
```

### toBase64(v)

Encodes the string `v` into Base64 format.

```expr
toBase64("Hello World") == "SGVsbG8gV29ybGQ="
```

### fromBase64(v)

Decodes the Base64 encoded string `v` back to its original form.

```expr
fromBase64("SGVsbG8gV29ybGQ=") == "Hello World"
```

### toPairs(map)

Converts a map to an array of key-value pairs.

```expr
toPairs({"name": "John", "age": 30}) == [["name", "John"], ["age", 30]]
```

### fromPairs(array)

Converts an array of key-value pairs to a map.

```expr
fromPairs([["name", "John"], ["age", 30]]) == {"name": "John", "age": 30}
```

## Miscellaneous Functions

### len(v)

Returns the length of an array, a map or a string.

### get(v, index)

Retrieves the element at the specified index from an array or map `v`. If the index is out of range, returns `nil`.
Or the key does not exist, returns `nil`.

```expr
get([1, 2, 3], 1) == 2
get({"name": "John", "age": 30}, "name") == "John"
```

## Predicate

The predicate is an expression. It takes one or more arguments and returns a boolean value.
To access the arguments, the `#` symbol is used.

```expr
map(0..9, {# / 2})
```

If items of the array is a struct or a map, it is possible to access fields with
omitted `#` symbol (`#.Value` becomes `.Value`).

```expr
filter(tweets, {len(.Value) > 280})
```

Braces `{` `}` can be omitted:

```expr
filter(tweets, len(.Value) > 280)
```

## `$env` variable

The `$env` variable is a map of all variables passed to the expression.

```expr
foo.Name == $env["foo"].Name
$env["var with spaces"]
```
