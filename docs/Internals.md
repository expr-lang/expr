# Internals

Expr is a stack based virtual machine. It compiles expressions to bytecode and 
runs it. 

Compilation is done in a few steps:
- Parse expression to AST
- Type check AST
  - Apply operator overloading
  - Apply custom AST patching
- Optimize AST
- Compile AST to a bytecode

Compiler has a bunch of optimization which will produce a more optimal program.

## In array

```js
value in ['foo', 'bar', 'baz']
```

If Expr finds an `in` or `not in` expression with an array, it will be 
transformed into:

```js
value in {"foo": true, "bar": true, "baz": true}
```

## Constant folding

Arithmetic expressions with constants is computed on compile step and replaced 
with the result.

```js
-(2-5)**3-2/(+4-3)+-2
```

Will be compiled to just single number:

```js
23
```

## In range

```js
user.Age in 18..32
```

Will be replaced with a binary operator:

```js
18 <= user.Age && user.Age <= 32
```

## Const range

```js
1..10_000
```

Ranges computed on compile stage, replaced with precreated slices.

## Const expr

If some function marked as constant expression with `expr.ConstExpr`. It will be
replaced with result of the call, if all arguments are constants.

```go
expr.ConstExpt("fib")
```

```js
fib(42)
``` 

Will be replaced with result of `fib(42)` on the compile step.

[ConstExpr Example](https://pkg.go.dev/github.com/antonmedv/expr?tab=doc#ConstExpr)
