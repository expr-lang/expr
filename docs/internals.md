# Internals

Expr is a stack-based virtual machine. It compiles expressions to bytecode and
runs it. The compilation is done in a few steps:

- Parse expression to AST
- Type check AST
    - Apply operator overloading
    - Apply custom AST patching
- Optimize AST
- Compile AST to a bytecode

The compiler has a bunch of optimization which will produce a more optimal program.

## In array

```
value in ['foo', 'bar', 'baz']
```

If Expr finds an `in` or `not in` expression with an array, it will be
transformed into:

```
value in {"foo": true, "bar": true, "baz": true}
```

## Constant folding

Arithmetic expressions with constants are computed on compile step and replaced
with the result.

```
-(2 - 5) ** 3 - 2 / (+4 - 3) + -2
```

Will be compiled into just a single number:

```
23
```

## In range

```
user.Age in 18..32
```

Will be replaced with a binary operator:

```
18 <= user.Age && user.Age <= 32
```

## Const range

```
1..10_000
```

Ranges computed on the compile stage, are replaced with pre-created slices.
