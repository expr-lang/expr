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
