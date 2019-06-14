package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"gopkg.in/antonmedv/expr.v2/compiler"
	"gopkg.in/antonmedv/expr.v2/parser"
	. "gopkg.in/antonmedv/expr.v2/vm"
	"sort"
	"strconv"
	"strings"
	"time"
)

func debugger() {
	tree, err := parser.Parse(input())
	check(err)

	program, err := compiler.Compile(tree)
	check(err)

	vm := NewVM(true)
	go vm.Run(program, nil)

	app := tview.NewApplication()
	table := tview.NewTable()
	stack := tview.NewTable()
	stack.
		SetBorder(true).
		SetTitle("Stack")
	scope := tview.NewTable()
	scope.
		SetBorder(true).
		SetTitle("Scope")
	sub := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(stack, 0, 3, false).
		AddItem(scope, 0, 1, false)
	flex := tview.NewFlex().
		AddItem(table, 0, 1, true).
		AddItem(sub, 0, 1, false)
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

	draw := func(ip int) {
		app.QueueUpdateDraw(func() {
			for row := 0; row < table.GetRowCount(); row++ {
				for col := 0; col < table.GetColumnCount(); col++ {
					table.GetCell(row, col).SetBackgroundColor(tcell.ColorDefault)
				}
			}

			if row, ok := index[ip]; ok {
				table.Select(row, 0)
				for col := 0; col < 5; col++ {
					table.GetCell(row, col).SetBackgroundColor(tcell.ColorMediumBlue)
				}
				table.SetOffset(row-10, 0)

				opcode := table.GetCell(row, 1).Text
				if strings.HasPrefix(opcode, "OpJump") {
					jump := table.GetCell(row, 3).Text
					jump = strings.Trim(jump, "()")
					ip, err := strconv.Atoi(jump)
					if err == nil {
						if row, ok := index[ip]; ok {
							for col := 0; col < 5; col++ {
								table.GetCell(row, col).SetBackgroundColor(tcell.ColorDimGrey)
							}
						}
					}
				}
			}

			stack.Clear()
			for i, value := range vm.Stack() {
				stack.SetCellSimple(i, 0, fmt.Sprintf("% *d: ", 2, i))
				stack.SetCellSimple(i, 1, fmt.Sprintf("%+v", value))
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

	getSelectedPosition := func() int {
		row, _ := table.GetSelection()
		ip, err := strconv.Atoi(strings.TrimSpace(table.GetCell(row, 0).Text))
		check(err)
		return ip
	}

	autostep := false
	var breakpoint int

	draw(0)
	go func() {
		for ip := range vm.Position() {
			draw(ip)

			if autostep {
				if breakpoint != ip {
					time.Sleep(20 * time.Millisecond)
					vm.Step()
				} else {
					autostep = false
				}
			}
		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyUp {
			table.SetSelectable(true, false)
		}
		if event.Key() == tcell.KeyEnter {
			selectable, _ := table.GetSelectable()
			if selectable {
				table.SetSelectable(false, false)
				breakpoint = getSelectedPosition()
				autostep = true
			}
			vm.Step()
		}
		return event
	})

	err = app.Run()
	check(err)
}
