package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	"github.com/sanity-io/litter"
	"io/ioutil"
	"os"
)

func main() {
	ast := flag.Bool("ast", false, "prints ast")
	hex := flag.Bool("hex", false, "prints bytecode")
	bytecode := flag.Bool("bytecode", false, "disassemble bytecode")
	flag.Parse()

	if *ast {
		b, _ := ioutil.ReadAll(os.Stdin)
		printAst(string(b))
	} else if *hex {
		b, _ := ioutil.ReadAll(os.Stdin)
		printHex(string(b))
	} else if *bytecode {
		b, _ := ioutil.ReadAll(os.Stdin)
		printDisassemble(string(b))
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

	_, _ = fmt.Fprintf(os.Stdout, compiler.Disassemble(program))
}
