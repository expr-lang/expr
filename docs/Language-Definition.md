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
</table>

Examples:

```expr
user.Age in 18..45 and user.Name not in ["admin", "root"]
```

```expr
foo matches "^[A-Z].*"
```

### Membership Operator

Fields of structs and items of maps can be accessed with `.` operator
or `[]` operator. Elements of arrays and slices can be accessed with
`[]` operator. Negative indices are supported with `-1` being
the last element.

The `in` operator can be used to check if an item is in an array or a map.

```python
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

## Built-in Functions

### all(array, predicate)

Returns **true** if all elements satisfies the [predicate](#predicate).
If the array is empty, returns **true**.

```expr
all(Tweets, {.Size < 280})
```

### any(array, predicate)

Returns **true** if any elements satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

### one(array, predicate)

Returns **true** if _exactly one_ element satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

```expr
one(Participants, {.Winner})
```

### none(array, predicate)

Returns **true** if _all elements does not_ satisfy the [predicate](#predicate).
If the array is empty, returns **true**.

### map(array, predicate)

Returns new array by applying the [predicate](#predicate) to each element of
the array.

### filter(array, predicate)

Returns new array by filtering elements of the array by [predicate](#predicate).

### count(array, predicate)

Returns the number of elements what satisfies the [predicate](#predicate).
Equivalent to:

```expr
len(filter(array, predicate))
```

### len(v)

Returns the length of an array, a map or a string.

### abs(v)

Returns the absolute value of a number.

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

### trim(v[, chars])

Removes white spaces from both ends of a string `v`.
If the optional `chars` argument is given, it is a string specifying the set of characters to be removed.

```expr
trim("  Hello  ") == "Hello"
trim("__Hello__", "_") == "Hello"
```

### trimPrefix(v, prefix)

Removes the specified prefix from the string `v` if it starts with that prefix.

```expr
trimPrefix("HelloWorld", "Hello") == "World"
```

### trimSuffix(v, suffix)

Removes the specified suffix from the string `v` if it ends with that suffix.

```expr
trimSuffix("HelloWorld", "World") == "Hello"
```

### upper(v)

Converts all the characters in string `v` to uppercase.

```expr
upper("hello") == "HELLO"
```

### lower(v)

Converts all the characters in string `v` to lowercase.

```expr
lower("HELLO") == "hello"
```

### split(v, delimiter)

Splits the string `v` at each instance of the delimiter and returns an array of substrings.

```expr
split("apple,orange,grape", ",") == ["apple", "orange", "grape"]
```

### splitN(v, delimiter, n)

Splits the string `v` at each instance of the delimiter but limits the result to `n` substrings.

```expr
splitN("apple,orange,grape", ",", 2) == ["apple", "orange,grape"]
```

### splitAfter(v, delimiter)

Splits the string `v` after each instance of the delimiter.

```expr
splitAfter("apple,orange,grape", ",") == ["apple,", "orange,", "grape"]
```

### splitAfterN(v, delimiter, n)

Splits the string `v` after each instance of the delimiter but limits the result to `n` substrings.

```expr
splitAfterN("apple,orange,grape", ",", 2) == ["apple,", "orange,grape"]
```

### replace(v, old, new)

Replaces all occurrences of `old` in string `v` with `new`.

```expr
replace("Hello World", "World", "Universe") == "Hello Universe"
```

### repeat(v, n)

Repeats the string `v` `n` times.

```expr
repeat("Hi", 3) == "HiHiHi"
```

### join(v, delimiter)

Joins an array of strings `v` into a single string with the given delimiter.

```expr
join(["apple", "orange", "grape"], ",") == "apple,orange,grape"
```

### indexOf(v, substring)

Returns the index of the first occurrence of the substring in string `v` or -1 if not found.

```expr
indexOf("apple pie", "pie") == 6
```

### lastIndexOf(v, substring)

Returns the index of the last occurrence of the substring in string `v` or -1 if not found.

```expr
lastIndexOf("apple pie apple", "apple") == 10
```

### hasPrefix(v, prefix)

Returns `true` if string `v` starts with the given prefix.

```expr
hasPrefix("HelloWorld", "Hello") == true
```

### hasSuffix(v, suffix)

Returns `true` if string `v` ends with the given suffix.

```expr
hasSuffix("HelloWorld", "World") == true
```

### max(v1, v2)

Returns the maximum of the two values `v1` and `v2`.

```expr
max(5, 7) == 7
```

### min(v1, v2)

Returns the minimum of the two values `v1` and `v2`.

```expr
min(5, 7) == 5
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

### now()

Returns the current date and time.

```expr
createdAt > now() - duration(1h)
```

### duration(v)

Returns [time.Duration](https://pkg.go.dev/time#Duration) value of the given string `v`.

Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

```expr
duration("1h").Seconds() == 3600
```

### date(v[, format, timezone])

Converts the given value `v` into a date representation.

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

### first(v)

Returns the first element from an array `v`. If the array is empty, returns `nil`.

```expr
first([1, 2, 3]) == 1
```

### last(v)

Returns the last element from an array `v`. If the array is empty, returns `nil`.

```expr
last([1, 2, 3]) == 3
```

### get(v, index)

Retrieves the element at the specified index from an array or map `v`. If the index is out of range, returns `nil`.
Or the key does not exist, returns `nil`.

```expr
get([1, 2, 3], 1) == 2
get({"name": "John", "age": 30}, "name") == "John"
```

## Predicate

The predicate is an expression that accepts a single argument. To access
the argument use the `#` symbol.

```expr
map(0..9, {# / 2})
```

If items of the array is a struct or a map, it is possible to access fields with
omitted `#` symbol (`#.Value` becomes `.Value`).

```expr
filter(Tweets, {len(.Value) > 280})
```

Braces `{` `}` can be omitted:

```expr
filter(Tweets, len(.Value) > 280)
```

## `env` variable

The `env` variable is a map of all variables passed to the expression.

```expr
Foo.Name == env['Foo'].Name
```
