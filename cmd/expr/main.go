package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/sanity-io/litter"
	"io/ioutil"
	"os"
)

func main() {
	ast := flag.Bool("ast", false, "prints ast")
	hex2 := flag.Bool("hex", false, "prints bytecode")
	bytecode := flag.Bool("bytecode", false, "disassemble bytecode")
	run := flag.Bool("run", false, "runs program")
	flag.Parse()

	if *ast {
		b, _ := ioutil.ReadAll(os.Stdin)
		printAst(string(b))
	} else if *hex2 {
		b, _ := ioutil.ReadAll(os.Stdin)
		printHex(string(b))
	} else if *bytecode {
		b, _ := ioutil.ReadAll(os.Stdin)
		printDisassemble(string(b))
	} else if *run {
		b, _ := ioutil.ReadAll(os.Stdin)
		runProgram(string(b))
	} else {
		flag.Usage()
		os.Exit(2)
	}
}

func printAst(source string) {
	node, err := parser.Parse(source)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	litter.Dump(node)
}

func printHex(source string) {
	node, err := parser.Parse(source)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	program, err := compiler.Compile(node)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, hex.Dump(program.Bytecode))
}

func printDisassemble(source string) {
	node, err := parser.Parse(source)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	program, err := compiler.Compile(node)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, program.Disassemble())
}

func runProgram(source string) {
	tree, err := parser.Parse(source)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	program, err := compiler.Compile(tree)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	out, err := vm.RunSafe(program, nil)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	litter.Dump(out)
}
