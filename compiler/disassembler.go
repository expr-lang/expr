package compiler

import (
	"encoding/binary"
	"fmt"
	"github.com/antonmedv/expr/vm"
)

func Disassemble(program vm.Program) string {
	out := ""
	ip := 0
	for ip < len(program.Bytecode) {
		b := program.Bytecode[ip]
		cp := ip
		ip++

		arg := func() uint16 {
			if ip+1 >= len(program.Bytecode) {
				return 0
			}

			i := binary.LittleEndian.Uint16([]byte{program.Bytecode[ip], program.Bytecode[ip+1]})
			ip += 2
			return i
		}

		printOp := func(b string) {
			out += fmt.Sprintf("%v\t%v\n", cp, b)
		}

		printArg := func(b string) {
			out += fmt.Sprintf("%v\t%v\t%v\n", cp, b, arg())
		}

		printConst := func(b string) {
			a := arg()
			var c interface{}
			if int(a) < len(program.Constant) {
				c = program.Constant[a]
			}
			out += fmt.Sprintf("%v\t%v\t%v\t%#v\n", cp, b, a, c)
		}

		switch b {
		case vm.OpPush:
			printArg("OpPush")
		case vm.OpPop:
			printOp("OpPop")
		case vm.OpLoad:
			printConst("OpLoad")
		case vm.OpFetch:
			printConst("OpFetch")
		case vm.OpTrue:
			printOp("OpTrue")
		case vm.OpFalse:
			printOp("OpFalse")
		case vm.OpNegate:
			printOp("OpNegate")
		case vm.OpNot:
			printOp("OpNot")
		case vm.OpJumpIfTrue:
			printArg("OpJumpIfTrue")
		case vm.OpJumpIfFalse:
			printArg("OpJumpIfFalse")

		default:
			out += fmt.Sprintf("%v\t%#v\n", cp, b)
		}
	}
	return out
}
