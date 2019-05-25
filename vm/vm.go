package vm

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

func Run(program Program, env interface{}) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	vm := vm{
		env:      env,
		bytecode: program.Bytecode,
		constant: program.Constant,
	}
	out = vm.run()

	return
}

type vm struct {
	env      interface{}
	stack    []interface{}
	bytecode []byte
	ip       int
	constant []interface{}
}

func (vm *vm) run() interface{} {
	for vm.ip < len(vm.bytecode) {

		b := vm.bytecode[vm.ip]
		vm.ip++

		switch b {

		case OpPush:
			vm.push(int64(vm.arg()))

		case OpPop:
			vm.pop()

		case OpLoad:
			vm.push(vm.constant[vm.arg()])

		case OpFetch:
			vm.push(fetch(vm.env, vm.constant[vm.arg()]))

		case OpTrue:
			vm.push(true)

		case OpFalse:
			vm.push(false)

		case OpNil:
			vm.push(new(struct{}))

		case OpNegate:
			v := negate(vm.pop())
			vm.push(v)

		case OpNot:
			v := vm.pop().(bool)
			vm.push(!v)

		case OpEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(equal(a, b))

		case OpJumpIfTrue:
			offset := vm.arg()
			if vm.current().(bool) {
				vm.ip += int(offset)
			}

		case OpJumpIfFalse:
			offset := vm.arg()
			if !vm.current().(bool) {
				vm.ip += int(offset)
			}

		case OpIn:
			b := vm.pop()
			a := vm.pop()
			vm.push(in(a, b))

		case OpLess:
			b := vm.pop()
			a := vm.pop()
			vm.push(less(a, b))

		case OpMore:
			b := vm.pop()
			a := vm.pop()
			vm.push(more(a, b))

		case OpLessOrEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(lessOrEqual(a, b))

		case OpMoreOrEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(moreOrEqual(a, b))

		case OpAdd:
			b := vm.pop()
			a := vm.pop()
			vm.push(add(a, b))

		case OpSubtract:
			b := vm.pop()
			a := vm.pop()
			vm.push(subtract(a, b))

		case OpMultiply:
			b := vm.pop()
			a := vm.pop()
			vm.push(multiply(a, b))

		case OpDivide:
			b := vm.pop()
			a := vm.pop()
			vm.push(divide(a, b))

		case OpModulo:
			b := vm.pop()
			a := vm.pop()
			vm.push(modulo(a, b))

		case OpExponent:
			b := vm.pop()
			a := vm.pop()
			vm.push(exponent(a, b))

		case OpContains:
			b := vm.pop()
			a := vm.pop()
			vm.push(strings.Contains(a.(string), b.(string)))

		case OpRange:
			b := vm.pop()
			a := vm.pop()
			vm.push(makeRange(a, b))

		case OpMatches:
			b := vm.pop()
			a := vm.pop()

			match, err := regexp.MatchString(a.(string), b.(string))
			if err != nil {
				panic(err)
			}

			vm.push(match)

		case OpMatchesConst:
			a := vm.pop()
			r := vm.constant[vm.arg()].(*regexp.Regexp)
			vm.push(r.MatchString(a.(string)))

		default:
			panic(fmt.Sprintf("unknown bytecode %v", b))
		}
	}

	if len(vm.stack) > 0 {
		return vm.pop()
	}

	return nil
}

func (vm *vm) push(value interface{}) {
	vm.stack = append(vm.stack, value)
}

func (vm *vm) current() interface{} {
	return vm.stack[len(vm.stack)-1]
}

func (vm *vm) pop() interface{} {
	value := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return value
}

func (vm *vm) arg() uint16 {
	arg := binary.LittleEndian.Uint16([]byte{vm.bytecode[vm.ip], vm.bytecode[vm.ip+1]})
	vm.ip += 2
	return arg
}
