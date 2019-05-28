# exp

## Install

```bash
go get github.com/antonmedv/expr/cmd/exp
```

## Usage

Print ast of program.

```bash
echo '1 + 2' | exp -ast
```

Disassemble bytecode to human readable format.

```bash
echo 'map(0..9, {# * 2})' | exp -bytecode
```

Run expression.

```bash
echo '2**8' | exp -run
```

Start interactive debugger.

```bash
echo 'all(1..3, {# > 0})' | exp -debug
```
