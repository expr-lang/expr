module github.com/antonmedv/expr/repl

go 1.20

require (
	github.com/antonmedv/expr v0.0.0
	github.com/antonmedv/expr/debug v0.0.0
	github.com/chzyer/readline v1.5.1
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell v1.3.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/rivo/tview v0.0.0-20200219210816-cd38d7432498 // indirect
	github.com/rivo/uniseg v0.1.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.3.8 // indirect
)

replace github.com/antonmedv/expr => ../

replace github.com/antonmedv/expr/debug => ../debug
