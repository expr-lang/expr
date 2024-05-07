module github.com/expr-lang/expr/repl

go 1.20

require (
	github.com/bettercap/readline v0.0.0-20210228151553-655e48bcb7bf
	github.com/expr-lang/expr v1.13.0
	github.com/expr-lang/expr/debug v0.0.0
)

require (
	github.com/chzyer/test v1.0.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.6.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/tview v0.0.0-20230814110005-ccc2c8119703 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/term v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
)

replace github.com/expr-lang/expr => ../

replace github.com/expr-lang/expr/debug => ../debug
