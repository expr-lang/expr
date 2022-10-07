# Expr Debugger

<p align="center">
    <img src="/docs/images/debug.gif" alt="debugger" width="600">
</p>

## Install

```bash
go get github.com/antonmedv/expr/cmd/exprdbg
```

## Usage

Print ast of program.

```bash
echo '1 + 2' | exprdbg -ast
```

Disassemble bytecode to human-readable format.

```bash
echo 'map(0..9, {# * 2})' | exprdbg -bytecode
```

Run expression.

```bash
echo '2**8' | exprdbg -run
```

Start interactive debugger.

```bash
echo 'all(1..3, {# > 0})' | exprdbg -debug
```
