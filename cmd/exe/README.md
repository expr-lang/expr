# Expr Debugger

## Install

```bash
go get github.com/antonmedv/expr/cmd/exe
```

## Usage

Print ast of program.

```bash
echo '1 + 2' | exe -ast
```

Disassemble bytecode to human readable format.

```bash
echo 'map(0..9, {# * 2})' | exe -bytecode
```

Run expression.

```bash
echo '2**8' | exe -run
```

Start interactive debugger.

```bash
echo 'all(1..3, {# > 0})' | exe -debug
```
