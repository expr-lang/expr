package main

import (
	"flag"
	"fmt"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/sanity-io/litter"
	"io/ioutil"
	"os"
)

var (
	bytecode bool
	debug    bool
	run      bool
	ast      bool
)

func init() {
	flag.BoolVar(&bytecode, "bytecode", false, "disassemble bytecode")
	flag.BoolVar(&debug, "debug", false, "debug program")
	flag.BoolVar(&run, "run", false, "run program")
	flag.BoolVar(&ast, "ast", false, "print ast")
}

func main() {
	flag.Parse()

	if ast {
		printAst()
		os.Exit(0)
	}
	if bytecode {
		printDisassemble()
		os.Exit(0)
	}
	if run {
		runProgram()
		os.Exit(0)
	}
	if debug {
		debugger()
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(2)
}

func input() string {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func printAst() {
	tree, err := parser.Parse(input())
	check(err)
	litter.Dump(tree.Node)
}

func printDisassemble() {
	node, err := parser.Parse(input())
	check(err)

	program, err := compiler.Compile(node)
	check(err)

	_, _ = fmt.Fprintf(os.Stdout, program.Disassemble())
}

func runProgram() {
	tree, err := parser.Parse(input())
	check(err)

	program, err := compiler.Compile(tree)
	check(err)

	out, err := vm.Run(program, nil)
	check(err)

	litter.Dump(out)
}
