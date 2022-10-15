package main

import (
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm/debug"
)

func debugger() {
	tree, err := parser.Parse(input())
	check(err)

	_, err = checker.Check(tree, nil)
	check(err)

	if optFlag {
		err = optimizer.Optimize(&tree.Node, nil)
		check(err)
	}

	program, err := compiler.Compile(tree, nil)
	check(err)

	debug.StartDebugger(program, nil)
}
