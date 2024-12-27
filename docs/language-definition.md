# Language Definition

**Expr** is a simple expression language that can be used to evaluate expressions.

## Literals

<table>
    <tr>
        <td><strong>Comment</strong></td>
        <td>
             <code>/* */</code> or <code>//</code>
        </td>
    </tr>
    <tr>
        <td><strong>Boolean</strong></td>
        <td>
            <code>true</code>, <code>false</code>
        </td>
    </tr>
    <tr>
        <td><strong>Integer</strong></td>
        <td>
            <code>42</code>, <code>0x2A</code>, <code>0o52</code>, <code>0b101010</code>
        </td>
    </tr>
    <tr>
        <td><strong>Float</strong></td>
        <td>
            <code>0.5</code>, <code>.5</code>
        </td>
    </tr>
    <tr>
        <td><strong>String</strong></td>
        <td>
            <code>"foo"</code>, <code>'bar'</code>
        </td>
    </tr>
    <tr>
        <td><strong>Array</strong></td>
        <td>
            <code>[1, 2, 3]</code>
        </td>
    </tr>
    <tr>
        <td><strong>Map</strong></td>
        <td>
            <code>&#123;a: 1, b: 2, c: 3&#125;</code>
        </td>
    </tr>
    <tr>
        <td><strong>Nil</strong></td>
        <td>
            <code>nil</code>
        </td>
    </tr>
</table>

### Strings

Strings can be enclosed in single quotes or double quotes. Strings can contain escape sequences, like `\n` for newline,
`\t` for tab, `\uXXXX` for Unicode code points.

```expr
"Hello\nWorld"
```

For multiline strings, use backticks:

```expr
`Hello
World`
```

Backticks strings are raw strings, they do not support escape sequences.

## Operators

<table>
    <tr>
        <td><strong>Arithmetic</strong></td>
        <td>
            <code>+</code>, <code>-</code>, <code>*</code>, <code>/</code>, <code>%</code> (modulus), <code>^</code> or <code>**</code> (exponent)
        </td>
    </tr>
    <tr>
        <td><strong>Comparison</strong></td>
        <td>
            <code>==</code>, <code>!=</code>, <code>&lt;</code>, <code>&gt;</code>, <code>&lt;=</code>, <code>&gt;=</code>
        </td>
    </tr>
    <tr>
        <td><strong>Logical</strong></td>
        <td>
            <code>not</code> or <code>!</code>, <code>and</code> or <code>&amp;&amp;</code>, <code>or</code> or <code>||</code>
        </td>
    </tr>
    <tr>
        <td><strong>Conditional</strong></td>
        <td>
            <code>?:</code> (ternary), <code>??</code> (nil coalescing)
        </td>
    </tr>
    <tr>
        <td><strong>Membership</strong></td>
        <td>
            <code>[]</code>, <code>.</code>, <code>?.</code>, <code>in</code>
        </td>
    </tr>
    <tr>
        <td><strong>String</strong></td>
        <td>
            <code>+</code> (concatenation), <code>contains</code>, <code>startsWith</code>, <code>endsWith</code>
        </td>
    </tr>
    <tr>
        <td><strong>Regex</strong></td>
        <td>
            <code>matches</code>
        </td>
    </tr>
    <tr>
        <td><strong>Range</strong></td>
        <td>
            <code>..</code>
        </td>
    </tr>
    <tr>
        <td><strong>Slice</strong></td>
        <td>
            <code>[:]</code>
        </td>
    </tr>
    <tr>
        <td><strong>Pipe</strong></td>
        <td>
            <code>|</code>
        </td>
    </tr>
</table>

### Membership Operator

Fields of structs and items of maps can be accessed with `.` operator
or `[]` operator. Next two expressions are equivalent:

```expr
user.Name
user["Name"]
``` 

Elements of arrays and slices can be accessed with
`[]` operator. Negative indices are supported with `-1` being
the last element.

```expr
array[0] // first element
array[-1] // last element
```

The `in` operator can be used to check if an item is in an array or a map.

```expr
"John" in ["John", "Jane"]
"name" in {"name": "John", "age": 30}
```

#### Optional chaining

The `?.` operator can be used to access a field of a struct or an item of a map
without checking if the struct or the map is `nil`. If the struct or the map is
`nil`, the result of the expression is `nil`.

```expr
author.User?.Name
```

Is equivalent to:

```expr
author.User != nil ? author.User.Name : nil
```

#### Nil coalescing

The `??` operator can be used to return the left-hand side if it is not `nil`,
otherwise the right-hand side is returned.

```expr
author.User?.Name ?? "Anonymous"
```

Is equivalent to:

```expr
author.User != nil ? author.User.Name : "Anonymous"
```

### Slice Operator

The slice operator `[:]` can be used to access a slice of an array.

For example, variable **array** is `[1, 2, 3, 4, 5]`:

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

```expr
user.Name | lower() | split(" ")
```

Is equivalent to:

```expr
split(lower(user.Name), " ")
```

### Range Operator

The range operator `..` can be used to create a range of integers.

```expr
1..3 == [1, 2, 3]
```

## Variables

Variables can be declared with the `let` keyword. The variable name must start with a letter or an underscore.
The variable name can contain letters, digits and underscores. After the variable is declared, it can be used in the
expression.

```expr
let x = 42; x * 2
```

A few variables can be declared by a few `let` statements separated by a semicolon.

```expr
let x = 42; 
let y = 2; 
x * y
```

Here is an example of variable with pipe operator:

```expr
let name = user.Name | lower() | split(" "); 
"Hello, " + name[0] + "!"
```

### $env

The `$env` variable is a map of all variables passed to the expression.

```expr
foo.Name == $env["foo"].Name
$env["var with spaces"]
```

Think of `$env` as a global variable that contains all variables.

The `$env` can be used to check if a variable is defined:

```expr
'foo' in $env
```

## Predicate

The predicate is an expression. Predicates can be used in functions like `filter`, `all`, `any`, `one`, `none`, etc.
For example, next expression creates a new array from 0 to 9 and then filters it by even numbers:

```expr
filter(0..9, {# % 2 == 0})
```

If items of the array is a struct or a map, it is possible to access fields with
omitted `#` symbol (`#.Value` becomes `.Value`).

```expr
filter(tweets, {len(.Content) > 240})
```

Braces `{` `}` can be omitted:

```expr
filter(tweets, len(.Content) > 240)
```

:::tip
In nested predicates, to access the outer variable, use [variables](#variables).

```expr
filter(posts, {
    let post = #; 
    any(.Comments, .Author == post.Author)
}) 
```

:::

## String Functions

### trim(str[, chars]) {#trim}

Removes white spaces from both ends of a string `str`.
If the optional `chars` argument is given, it is a string specifying the set of characters to be removed.

```expr
trim("  Hello  ") == "Hello"
trim("__Hello__", "_") == "Hello"
```

### trimPrefix(str, prefix) {#trimPrefix}

Removes the specified prefix from the string `str` if it starts with that prefix.

```expr
trimPrefix("HelloWorld", "Hello") == "World"
```

### trimSuffix(str, suffix) {#trimSuffix}

Removes the specified suffix from the string `str` if it ends with that suffix.

```expr
trimSuffix("HelloWorld", "World") == "Hello"
```

### upper(str) {#upper}

Converts all the characters in string `str` to uppercase.

```expr
upper("hello") == "HELLO"
```

### lower(str) {#lower}

Converts all the characters in string `str` to lowercase.

```expr
lower("HELLO") == "hello"
```

### split(str, delimiter[, n]) {#split}

Splits the string `str` at each instance of the delimiter and returns an array of substrings.

```expr
split("apple,orange,grape", ",") == ["apple", "orange", "grape"]
split("apple,orange,grape", ",", 2) == ["apple", "orange,grape"]
```

### splitAfter(str, delimiter[, n]) {#splitAfter}

Splits the string `str` after each instance of the delimiter.

```expr
splitAfter("apple,orange,grape", ",") == ["apple,", "orange,", "grape"]
splitAfter("apple,orange,grape", ",", 2) == ["apple,", "orange,grape"]
```

### replace(str, old, new) {#replace}

Replaces all occurrences of `old` in string `str` with `new`.

```expr
replace("Hello World", "World", "Universe") == "Hello Universe"
```

### repeat(str, n) {#repeat}

Repeats the string `str` `n` times.

```expr
repeat("Hi", 3) == "HiHiHi"
```

### indexOf(str, substring) {#indexOf}

Returns the index of the first occurrence of the substring in string `str` or -1 if not found.

```expr
indexOf("apple pie", "pie") == 6
```

### lastIndexOf(str, substring) {#lastIndexOf}

Returns the index of the last occurrence of the substring in string `str` or -1 if not found.

```expr
lastIndexOf("apple pie apple", "apple") == 10
```

### hasPrefix(str, prefix) {#hasPrefix}

Returns `true` if string `str` starts with the given prefix.

```expr
hasPrefix("HelloWorld", "Hello") == true
```

### hasSuffix(str, suffix) {#hasSuffix}

Returns `true` if string `str` ends with the given suffix.

```expr
hasSuffix("HelloWorld", "World") == true
```

## Date Functions

Expr has a built-in support for Go's [time package](https://pkg.go.dev/time).
It is possible to subtract two dates and get the duration between them:

```expr
createdAt - now()
```

It is possible to add a duration to a date:

```expr
createdAt + duration("1h")
```

And it is possible to compare dates:

```expr
createdAt > now() - duration("1h")
```

### now() {#now}

Returns the current date as a [time.Time](https://pkg.go.dev/time#Time) value.

```expr
now().Year() == 2024
```

### duration(str) {#duration}

Returns [time.Duration](https://pkg.go.dev/time#Duration) value of the given string `str`.

Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

```expr
duration("1h").Seconds() == 3600
```

### date(str[, format[, timezone]]) {#date}

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

Available methods on the date:

- `Year()` - returns the year 
- `Month()` - returns the month (starting from 1)
- `Day()` - returns the day of the month
- `Hour()` - returns the hour
- `Minute()` - returns the minute
- `Second()` - returns the second
- `Weekday()` - returns the day of the week
- `YearDay()` - returns the day of the year
- and [more](https://pkg.go.dev/time#Time).

```expr
date("2023-08-14").Year() == 2023
```

### timezone(str) {#timezone}

Returns the timezone of the given string `str`. List of available timezones can be
found [here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).

```expr
timezone("Europe/Zurich")
timezone("UTC")
```

To convert a date to a different timezone, use the [`In()`](https://pkg.go.dev/time#Time.In) method:

```expr
date("2023-08-14 00:00:00").In(timezone("Europe/Zurich"))
```

## Number Functions

### max(n1, n2) {#max}

Returns the maximum of the two numbers `n1` and `n2`.

```expr
max(5, 7) == 7
```

### min(n1, n2) {#min}

Returns the minimum of the two numbers `n1` and `n2`.

```expr
min(5, 7) == 5
```

### abs(n) {#abs}

Returns the absolute value of a number.

```expr
abs(-5) == 5
```

### ceil(n) {#ceil}

Returns the least integer value greater than or equal to x.

```expr
ceil(1.5) == 2.0
```

### floor(n) {#floor}

Returns the greatest integer value less than or equal to x.

```expr
floor(1.5) == 1.0
```

### round(n) {#round}

Returns the nearest integer, rounding half away from zero.

```expr
round(1.5) == 2.0
```

## Array Functions

### all(array, predicate) {#all}

Returns **true** if all elements satisfies the [predicate](#predicate).
If the array is empty, returns **true**.

```expr
all(tweets, {.Size < 280})
```

### any(array, predicate) {#any}

Returns **true** if any elements satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

```expr
any(tweets, {.Size > 280})
```

### one(array, predicate) {#one}

Returns **true** if _exactly one_ element satisfies the [predicate](#predicate).
If the array is empty, returns **false**.

```expr
one(participants, {.Winner})
```

### none(array, predicate) {#none}

Returns **true** if _all elements does not_ satisfy the [predicate](#predicate).
If the array is empty, returns **true**.

```expr
none(tweets, {.Size > 280})
```

### map(array, predicate) {#map}

Returns new array by applying the [predicate](#predicate) to each element of
the array.

```expr
map(tweets, {.Size})
```

### filter(array, predicate) {#filter}

Returns new array by filtering elements of the array by [predicate](#predicate).

```expr
filter(users, .Name startsWith "J")
```

### find(array, predicate) {#find}

Finds the first element in an array that satisfies the [predicate](#predicate).

```expr
find([1, 2, 3, 4], # > 2) == 3
```

### findIndex(array, predicate) {#findIndex}

Finds the index of the first element in an array that satisfies the [predicate](#predicate).

```expr
findIndex([1, 2, 3, 4], # > 2) == 2
```

### findLast(array, predicate) {#findLast}

Finds the last element in an array that satisfies the [predicate](#predicate).

```expr
findLast([1, 2, 3, 4], # > 2) == 4
```

### findLastIndex(array, predicate) {#findLastIndex}

Finds the index of the last element in an array that satisfies the [predicate](#predicate).

```expr
findLastIndex([1, 2, 3, 4], # > 2) == 3
```

### groupBy(array, predicate) {#groupBy}

Groups the elements of an array by the result of the [predicate](#predicate).

```expr
groupBy(users, .Age)
```

### count(array[, predicate]) {#count}

Returns the number of elements what satisfies the [predicate](#predicate).

```expr
count(users, .Age > 18)
```

Equivalent to:

```expr
len(filter(users, .Age > 18))
```

If the predicate is not given, returns the number of `true` elements in the array.

```expr
count([true, false, true]) == 2
```

### concat(array1, array2[, ...]) {#concat}

Concatenates two or more arrays.

```expr
concat([1, 2], [3, 4]) == [1, 2, 3, 4]
```

### flatten(array) {#flatten}

Flattens given array into one-dimensional array.

```expr
flatten([1, 2, [3, 4]]) == [1, 2, 3, 4]
```

### join(array[, delimiter]) {#join}

Joins an array of strings into a single string with the given delimiter.
If no delimiter is given, an empty string is used.

```expr
join(["apple", "orange", "grape"], ",") == "apple,orange,grape"
join(["apple", "orange", "grape"]) == "appleorangegrape"
```

### reduce(array, predicate[, initialValue]) {#reduce}

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

### sum(array[, predicate]) {#sum}

Returns the sum of all numbers in the array.

```expr
sum([1, 2, 3]) == 6
```

If the optional `predicate` argument is given, it is a predicate that is applied on each element
of the array before summing.

```expr
sum(accounts, .Balance)
```

Equivalent to:

```expr
reduce(accounts, #acc + .Balance, 0)
// or
sum(map(accounts, .Balance))
```

### mean(array) {#mean}

Returns the average of all numbers in the array.

```expr
mean([1, 2, 3]) == 2.0
```

### median(array) {#median}

Returns the median of all numbers in the array.

```expr
median([1, 2, 3]) == 2.0
```

### first(array) {#first}

Returns the first element from an array. If the array is empty, returns `nil`.

```expr
first([1, 2, 3]) == 1
```

### last(array) {#last}

Returns the last element from an array. If the array is empty, returns `nil`.

```expr
last([1, 2, 3]) == 3
```

### take(array, n) {#take}

Returns the first `n` elements from an array. If the array has fewer than `n` elements, returns the whole array.

```expr
take([1, 2, 3, 4], 2) == [1, 2]
```

### reverse(array) {#reverse}

Return new reversed copy of the array.

```expr
reverse([3, 1, 4]) == [4, 1, 3]
reverse(reverse([3, 1, 4])) == [3, 1, 4]
```

### sort(array[, order]) {#sort}

Sorts an array in ascending order. Optional `order` argument can be used to specify the order of sorting: `asc`
or `desc`.

```expr
sort([3, 1, 4]) == [1, 3, 4]
sort([3, 1, 4], "desc") == [4, 3, 1]
```

### sortBy(array[, predicate, order]) {#sortBy}

Sorts an array by the result of the [predicate](#predicate). Optional `order` argument can be used to specify the order
of sorting: `asc` or `desc`.

```expr
sortBy(users, .Age)
sortBy(users, .Age, "desc")
```

## Map Functions

### keys(map) {#keys}

Returns an array containing the keys of the map.

```expr
keys({"name": "John", "age": 30}) == ["name", "age"]
```

### values(map) {#values}

Returns an array containing the values of the map.

```expr
values({"name": "John", "age": 30}) == ["John", 30]
```

## Type Conversion Functions

### type(v) {#type}

Returns the type of the given value `v`.

Returns on of the following types:

- `nil`
- `bool`
- `int`
- `uint`
- `float`
- `string`
- `array`
- `map`.

For named types and structs, the type name is returned.

```expr
type(42) == "int"
type("hello") == "string"
type(now()) == "time.Time"
```

### int(v) {#int}

Returns the integer value of a number or a string.

```expr
int("123") == 123
```

### float(v) {#float}

Returns the float value of a number or a string.

```expr
float("123.45") == 123.45
```

### string(v) {#string}

Converts the given value `v` into a string representation.

```expr
string(123) == "123"
```

### toJSON(v) {#toJSON}

Converts the given value `v` to its JSON string representation.

```expr
toJSON({"name": "John", "age": 30})
```

### fromJSON(v) {#fromJSON}

Parses the given JSON string `v` and returns the corresponding value.

```expr
fromJSON('{"name": "John", "age": 30}')
```

### toBase64(v) {#toBase64}

Encodes the string `v` into Base64 format.

```expr
toBase64("Hello World") == "SGVsbG8gV29ybGQ="
```

### fromBase64(v) {#fromBase64}

Decodes the Base64 encoded string `v` back to its original form.

```expr
fromBase64("SGVsbG8gV29ybGQ=") == "Hello World"
```

### toPairs(map) {#toPairs}

Converts a map to an array of key-value pairs.

```expr
toPairs({"name": "John", "age": 30}) == [["name", "John"], ["age", 30]]
```

### fromPairs(array) {#fromPairs}

Converts an array of key-value pairs to a map.

```expr
fromPairs([["name", "John"], ["age", 30]]) == {"name": "John", "age": 30}
```

## Miscellaneous Functions

### len(v) {#len}

Returns the length of an array, a map or a string.

```expr
len([1, 2, 3]) == 3
len({"name": "John", "age": 30}) == 2
len("Hello") == 5
```

### get(v, index) {#get}

Retrieves the element at the specified index from an array or map `v`. If the index is out of range, returns `nil`.
Or the key does not exist, returns `nil`.

```expr
get([1, 2, 3], 1) == 2
get({"name": "John", "age": 30}, "name") == "John"
```

## Bitwise Functions

### bitand(int, int) {#bitand}

Returns the values resulting from the bitwise AND operation.

```expr
bitand(0b1010, 0b1100) == 0b1000
```

### bitor(int, int) {#bitor}

Returns the values resulting from the bitwise OR operation.

```expr
bitor(0b1010, 0b1100) == 0b1110
```

### bitxor(int, int) {#bitxor}

Returns the values resulting from the bitwise XOR operation.

```expr
bitxor(0b1010, 0b1100) == 0b110
```

### bitnand(int, int) {#bitnand}

Returns the values resulting from the bitwise AND NOT operation.

```expr
bitnand(0b1010, 0b1100) == 0b10
```

### bitnot(int) {#bitnot}

Returns the values resulting from the bitwise NOT operation.

```expr
bitnot(0b1010) == -0b1011
```

### bitshl(int, int) {#bitshl}

Returns the values resulting from the Left Shift operation.

```expr
bitshl(0b101101, 2) == 0b10110100
```

### bitshr(int, int) {#bitshr}

Returns the values resulting from the Right Shift operation.

```expr
bitshr(0b101101, 2) == 0b1011
```

### bitushr(int, int) {#bitushr}

Returns the values resulting from the unsigned Right Shift operation.

```expr
bitushr(-0b101, 2) == 4611686018427387902
```
