package main

import (
	"fmt"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	. "github.com/antonmedv/expr/vm"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"sort"
	"strconv"
	"strings"
)

func debugger() {
	tree, err := parser.Parse(input())
	check(err)

	program, err := compiler.Compile(tree)
	check(err)

	vm := NewVM(program, nil, true)
	go vm.Run()

	app := tview.NewApplication()
	table := tview.NewTable()
	table.
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				vm.Step()
			}
		})
	stack := tview.NewTable()
	stack.
		SetBorder(true).
		SetTitle("Stack")
	scope := tview.NewTable()
	scope.
		SetBorder(true).
		SetTitle("Scope")
	sub := tview.NewFlex()
	sub.SetDirection(tview.FlexRow)
	sub.AddItem(stack, 0, 3, false)
	sub.AddItem(scope, 0, 1, false)
	flex := tview.NewFlex()
	flex.AddItem(table, 0, 1, true)
	flex.AddItem(sub, 0, 1, false)
	app.SetRoot(flex, true)

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

	pp := 0
	draw := func(ip int) {
		app.QueueUpdateDraw(func() {
			if row, ok := index[pp]; ok {
				for cel := 0; cel < 5; cel++ {
					table.GetCell(row, cel).SetBackgroundColor(tcell.ColorDefault)
				}
			}

			if row, ok := index[ip]; ok {
				for cel := 0; cel < 5; cel++ {
					table.GetCell(row, cel).SetBackgroundColor(tcell.ColorBlueViolet)
				}
				table.SetOffset(row-10, 0)
				pp = ip
			}

			stack.Clear()
			for i, value := range vm.Stack() {
				stack.SetCellSimple(i, 0, fmt.Sprintf("% *d: ", 2, i))
				stack.SetCellSimple(i, 1, fmt.Sprintf("%#v", value))
			}
			stack.ScrollToEnd()

			scope.Clear()
			var keys []string
			for k := range vm.Scope() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			row := 0
			for _, name := range keys {
				scope.SetCellSimple(row, 0, fmt.Sprintf("%v: ", name))
				scope.SetCellSimple(row, 1, fmt.Sprintf("%v", vm.Scope()[name]))
				row++
			}
		})
	}

	draw(0)
	go func() {
		for ip := range vm.Position() {
			draw(ip)
		}
	}()

	err = app.Run()
	check(err)
}
