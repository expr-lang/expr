package debug

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/antonmedv/expr/vm"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func StartDebugger(program *Program, env any) {
	vm := Debug()
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

	done := false
	go func() {
		out, err := vm.Run(program, env)
		done = true
		app.QueueUpdateDraw(func() {
			sub.RemoveItem(stack)
			sub.RemoveItem(scope)
			result := tview.NewTextView()
			result.
				SetBorder(true).
				SetTitle("Output")
			result.SetText(fmt.Sprintf("%#v", out))
			sub.AddItem(result, 0, 1, false)
			if err != nil {
				errorView := tview.NewTextView()
				errorView.
					SetBorder(true).
					SetTitle("Error")
				errorView.SetText(err.Error())
				sub.AddItem(errorView, 0, 1, false)
			}
		})
	}()

	index := make(map[int]int)
	var buf strings.Builder
	program.Opcodes(&buf)

	for row, line := range strings.Split(buf.String(), "\n") {
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
					table.GetCell(row, col).SetBackgroundColor(tcell.ColorBlack)
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
				stack.SetCellSimple(i, 1, fmt.Sprintf("%#v", value))
			}
			stack.ScrollToEnd()

			scope.Clear()
			s := vm.Scope()
			if s != nil {
				type pair struct {
					key   string
					value any
				}
				var keys []pair
				keys = append(keys, pair{"Array", s.Array})
				keys = append(keys, pair{"Index", s.Index})
				keys = append(keys, pair{"Len", s.Len})
				keys = append(keys, pair{"Count", s.Count})
				if s.GroupBy != nil {
					keys = append(keys, pair{"GroupBy", s.GroupBy})
				}
				if s.Acc != nil {
					keys = append(keys, pair{"Acc", s.Acc})
				}
				row := 0
				for _, pair := range keys {
					scope.SetCellSimple(row, 0, fmt.Sprintf("%v: ", pair.key))
					scope.SetCellSimple(row, 1, fmt.Sprintf("%v", pair.value))
					row++
				}
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

	go func() {
		draw(0)
		for ip := range vm.Position() {
			draw(ip)

			if autostep {
				if breakpoint != ip {
					time.Sleep(20 * time.Millisecond)
					if !done {
						vm.Step()
					}
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
			if !done {
				vm.Step()
			}
		}
		return event
	})

	err := app.Run()
	check(err)
}

func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
