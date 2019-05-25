package compiler

import (
	"encoding/binary"
	"fmt"
	"github.com/antonmedv/expr/vm"
	"regexp"
)

func Disassemble(program vm.Program) string {
	out := ""
	ip := 0
	for ip < len(program.Bytecode) {
		b := program.Bytecode[ip]
		cp := ip
		ip++

		readArg := func() uint16 {
			if ip+1 >= len(program.Bytecode) {
				return 0
			}

			i := binary.LittleEndian.Uint16([]byte{program.Bytecode[ip], program.Bytecode[ip+1]})
			ip += 2
			return i
		}

		op := func(b string) {
			out += fmt.Sprintf("%v\t%v\n", cp, b)
		}
		arg := func(b string) {
			out += fmt.Sprintf("%v\t%v\t%v\n", cp, b, readArg())
		}
		jump := func(b string) {
			a := readArg()
			out += fmt.Sprintf("%v\t%v\t%v\n", cp, b, cp+1+int(a))
		}
		constant := func(b string) {
			a := readArg()
			var c interface{}
			if int(a) < len(program.Constant) {
				c = program.Constant[a]
			}
			if r, ok := c.(*regexp.Regexp); ok {
				c = r.String()
			}
			out += fmt.Sprintf("%v\t%v\t%v\t%#v\n", cp, b, a, c)
		}

		switch b {
		case vm.OpPush:
			arg("OpPush")

		case vm.OpPop:
			op("OpPop")

		case vm.OpLoad:
			constant("OpLoad")

		case vm.OpFetch:
			constant("OpFetch")

		case vm.OpTrue:
			op("OpTrue")

		case vm.OpFalse:
			op("OpFalse")

		case vm.OpNegate:
			op("OpNegate")

		case vm.OpNot:
			op("OpNot")

		case vm.OpEqual:
			op("OpEqual")

		case vm.OpJumpIfTrue:
			jump("OpJumpIfTrue")

		case vm.OpJumpIfFalse:
			jump("OpJumpIfFalse")

		case vm.OpIn:
			op("OpIn")

		case vm.OpLess:
			op("OpLess")

		case vm.OpMore:
			op("OpMore")

		case vm.OpLessOrEqual:
			op("OpLessOrEqual")

		case vm.OpMoreOrEqual:
			op("OpMoreOrEqual")

		case vm.OpAdd:
			op("OpAdd")

		case vm.OpSubtract:
			op("OpSubtract")

		case vm.OpMultiply:
			op("OpMultiply")

		case vm.OpDivide:
			op("OpDivide")

		case vm.OpModulo:
			op("OpModulo")

		case vm.OpExponent:
			op("OpExponent")

		case vm.OpContains:
			op("OpContains")

		case vm.OpRange:
			op("OpRange")

		case vm.OpMatches:
			op("OpMatches")

		case vm.OpMatchesConst:
			constant("OpMatchesConst")

		case vm.OpField:
			op("OpField")

		case vm.OpFieldConst:
			constant("OpFieldConst")

		default:
			out += fmt.Sprintf("%v\t%#x\n", cp, b)
		}
	}
	return out
}
