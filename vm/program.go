package vm

import (
	"fmt"
	"regexp"

	"github.com/antonmedv/expr/file"
	"github.com/antonmedv/expr/vm/runtime"
)

type Program struct {
	Source    *file.Source
	Locations []file.Location
	Constants []interface{}
	Bytecode  []Opcode
}

func (program *Program) Disassemble() string {
	out := ""
	ip := 0
	for ip < len(program.Bytecode) {
		pp := ip
		op, arg := program.Bytecode[ip], program.Bytecode[ip+1]
		ip += 2

		code := func(label string) {
			out += fmt.Sprintf("%v\t%v\n", pp, label)
		}
		jump := func(label string) {
			out += fmt.Sprintf("%v\t%v\t%v\t(%v)\n", pp, label, arg, ip+int(arg))
		}
		back := func(label string) {
			out += fmt.Sprintf("%v\t%v\t%v\t(%v)\n", pp, label, arg, ip-int(arg))
		}
		argument := func(label string) {
			out += fmt.Sprintf("%v\t%v\t%v\n", pp, label, arg)
		}
		constant := func(label string) {
			var c interface{}
			if int(arg) < len(program.Constants) {
				c = program.Constants[arg]
			} else {
				c = "out of range"
			}
			if r, ok := c.(*regexp.Regexp); ok {
				c = r.String()
			}
			if field, ok := c.(*runtime.Field); ok {
				c = fmt.Sprintf("%v %v", field.Path, field.Index)
			}
			out += fmt.Sprintf("%v\t%v\t%v\t%v\n", pp, label, arg, c)
		}

		switch op {
		case OpPush:
			constant("OpPush")

		case OpPop:
			code("OpPop")

		case OpRot:
			code("OpRot")

		case OpFetch:
			code("OpFetch")

		case OpFetchField:
			constant("OpFetchField")

		case OpFetchEnv:
			constant("OpFetchEnv")

		case OpFetchEnvField:
			constant("OpFetchEnvField")

		case OpFetchEnvFast:
			constant("OpFetchEnvFast")

		case OpTrue:
			code("OpTrue")

		case OpFalse:
			code("OpFalse")

		case OpNil:
			code("OpNil")

		case OpNegate:
			code("OpNegate")

		case OpNot:
			code("OpNot")

		case OpEqual:
			code("OpEqual")

		case OpEqualInt:
			code("OpEqualInt")

		case OpEqualString:
			code("OpEqualString")

		case OpJump:
			jump("OpJump")

		case OpJumpIfTrue:
			jump("OpJumpIfTrue")

		case OpJumpIfFalse:
			jump("OpJumpIfFalse")

		case OpJumpIfNil:
			jump("OpJumpIfNil")

		case OpJumpBackward:
			back("OpJumpBackward")

		case OpIn:
			code("OpIn")

		case OpLess:
			code("OpLess")

		case OpMore:
			code("OpMore")

		case OpLessOrEqual:
			code("OpLessOrEqual")

		case OpMoreOrEqual:
			code("OpMoreOrEqual")

		case OpAdd:
			code("OpAdd")

		case OpSubtract:
			code("OpSubtract")

		case OpMultiply:
			code("OpMultiply")

		case OpDivide:
			code("OpDivide")

		case OpModulo:
			code("OpModulo")

		case OpExponent:
			code("OpExponent")

		case OpRange:
			code("OpRange")

		case OpMatches:
			code("OpMatches")

		case OpMatchesConst:
			constant("OpMatchesConst")

		case OpContains:
			code("OpContains")

		case OpStartsWith:
			code("OpStartsWith")

		case OpEndsWith:
			code("OpEndsWith")

		case OpSlice:
			code("OpSlice")

		case OpCall:
			argument("OpCall")

		case OpCallFast:
			constant("OpCallFast")

		case OpMethod:
			constant("OpMethod")

		case OpArray:
			code("OpArray")

		case OpMap:
			code("OpMap")

		case OpLen:
			code("OpLen")

		case OpCast:
			argument("OpCast")

		case OpStore:
			constant("OpStore")

		case OpLoad:
			constant("OpLoad")

		case OpInc:
			constant("OpInc")

		case OpDeref:
			code("OpDeref")

		case OpBegin:
			code("OpBegin")

		case OpEnd:
			code("OpEnd")

		default:
			out += fmt.Sprintf("%v\t%#x\n", ip, op)
		}
	}
	return out
}
