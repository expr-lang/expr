package main

import (
	"fmt"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	. "github.com/antonmedv/expr/vm"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

func debugger() {
	tree, err := parser.Parse(input())
	check(err)

	program, err := compiler.Compile(tree)
	check(err)

	vm := NewVM(program, nil, true)

	app := tview.NewApplication()
	table := tview.NewTable()

	index := make(map[int]int)
	for row, line := range strings.Split(program.Disassemble(), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")

		ip, err := strconv.Atoi(parts[0])
		check(err)
		index[ip] = row
		table.SetCellSimple(row, 0, fmt.Sprintf("% *d", 5, ip))

		for col := 1; col < len(parts); col++ {
			table.SetCellSimple(row, col, parts[col])
		}
		for col := len(parts); col < 4; col++ {
			table.SetCellSimple(row, col, "")
		}
		table.SetCell(row, 4, tview.NewTableCell("").SetExpansion(1))
	}

	app.QueueUpdateDraw(func() {
		row, ok := index[57]
		if !ok {
			panic("missing ip pointer")
		}
		for cel := 0; cel < 5; cel++ {
			table.GetCell(row, cel).SetBackgroundColor(tcell.ColorBlueViolet)
		}
		table.SetOffset(row-10, 0)
	})

	flex := tview.NewFlex()
	flex.AddItem(table, 0, 1, true)

	app.SetRoot(flex, true)

	err = app.Run()
	check(err)
}
