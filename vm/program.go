package vm

import (
	"encoding/binary"
	"fmt"
	"regexp"
)

type Program struct {
	Constant []interface{}
	Bytecode []byte
}

func (program Program) Disassemble() string {
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
			out += fmt.Sprintf("%v\t%v\t%v (%v)\n", cp, b, a, ip+int(a))
		}
		back := func(b string) {
			a := readArg()
			out += fmt.Sprintf("%v\t%v\t%v (%v)\n", cp, b, a, ip-int(a))
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
		case OpPush:
			arg("OpPush")

		case OpPop:
			op("OpPop")

		case OpConst:
			constant("OpConst")

		case OpFetch:
			constant("OpFetch")

		case OpTrue:
			op("OpTrue")

		case OpFalse:
			op("OpFalse")

		case OpNegate:
			op("OpNegate")

		case OpNot:
			op("OpNot")

		case OpEqual:
			op("OpEqual")

		case OpJump:
			jump("OpJump")

		case OpJumpIfTrue:
			jump("OpJumpIfTrue")

		case OpJumpIfFalse:
			jump("OpJumpIfFalse")

		case OpJumpBackward:
			back("OpJumpBackward")

		case OpIn:
			op("OpIn")

		case OpLess:
			op("OpLess")

		case OpMore:
			op("OpMore")

		case OpLessOrEqual:
			op("OpLessOrEqual")

		case OpMoreOrEqual:
			op("OpMoreOrEqual")

		case OpAdd:
			op("OpAdd")

		case OpInc:
			op("OpInc")

		case OpSubtract:
			op("OpSubtract")

		case OpMultiply:
			op("OpMultiply")

		case OpDivide:
			op("OpDivide")

		case OpModulo:
			op("OpModulo")

		case OpExponent:
			op("OpExponent")

		case OpContains:
			op("OpContains")

		case OpRange:
			op("OpRange")

		case OpMatches:
			op("OpMatches")

		case OpMatchesConst:
			constant("OpMatchesConst")

		case OpIndex:
			op("OpIndex")

		case OpProperty:
			constant("OpProperty")

		case OpCall:
			constant("OpCall")

		case OpMethod:
			constant("OpMethod")

		case OpArray:
			op("OpArray")

		case OpMap:
			op("OpMap")

		case OpLen:
			op("OpLen")

		case OpBegin:
			op("OpBegin")

		case OpEnd:
			op("OpEnd")

		case OpStore:
			constant("OpStore")

		case OpLoad:
			constant("OpLoad")

		default:
			out += fmt.Sprintf("%v\t%#x\n", cp, b)
		}
	}
	return out
}
