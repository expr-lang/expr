package main

import (
	"flag"
	"fmt"
	"github.com/antonmedv/expr/parser"
	"github.com/sanity-io/litter"
	"io/ioutil"
	"os"
)

func main() {
	ast := flag.Bool("ast", false, "prints ast of given expression")
	flag.Parse()

	if *ast {
		b, _ := ioutil.ReadAll(os.Stdin)
		printAst(string(b))
	} else {
		flag.Usage()
		os.Exit(2)
	}
}

func printAst(source string) {
	ast, err := parser.Parse(source)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	litter.Dump(ast)
}
