# Optimizations

Expr has a bunch of optimization which will produce more optimal program during compile step.

## In array

```js
value in ['foo', 'bar', 'baz']
```

If expr finds an `in` or `not in` expression with an array, it will be transformed into:

```js
value in {"foo": true, "bar": true, "baz": true}
```

## Constant folding

Arithmetic expressions with constants is computed on compile step and replaced with result.

```js
-(2-5)**3-2/(+4-3)+-2
```

Will be compiled to just single number:

```js
23
```

So in expressions it's safe to use some arithmetics for better readability:

```js
percentage > 0.3 * 100
```

As it will be simpified to:

```js
percentage > 30
```

## In range

```js
user.Age in 18..32
```

Will be replaced with binary operator:

```js
18 <= user.Age && user.Age <= 32
```

`not in` operator will also work.

## Const range

```js
1..10_000
```

Ranges computed on compile stage, repleced with preallocated slices.

## Const expr

If some function marked as constant expression with `expr.ConstExpr`. It will be replaced with result
of call, if all arguments are constants.

```go
expr.ConstExpt("fib")
```

```js
fib(42)
``` 

Will be replaced with result of `fib(42)` on compile step. No need to calculate it during runtime.
