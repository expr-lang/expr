# ecc

## Install

```bash
go get github.com/antonmedv/expr/cmd/ecc
```

## Usage

Print ast of program.

```bash
echo '1 + 2' | ecc -ast
```

Disassemble bytecode to human readable format.

```bash
echo 'map(0..9, {# * 2})' | ecc -bytecode
```

Run expression.

```bash
echo '2**8' | ecc -run
```

Start interactive debugger.

```bash
echo 'all(1..3, {# > 0})' | ecc -debug
```
