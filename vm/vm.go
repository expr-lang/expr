package vm

//go:generate sh -c "go run ./func_types > ./generated.go"

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/antonmedv/expr/file"
	"github.com/antonmedv/expr/vm/runtime"
)

var (
	MemoryBudget int = 1e6
)

func Run(program *Program, env interface{}) (interface{}, error) {
	if program == nil {
		return nil, fmt.Errorf("program is nil")
	}

	vm := VM{}
	return vm.Run(program, env)
}

type VM struct {
	stack        []interface{}
	ip           int
	scopes       []*Scope
	debug        bool
	step         chan struct{}
	curr         chan int
	memory       int
	memoryBudget int
}

type Scope struct {
	Array reflect.Value
	It    int
	Len   int
	Count int
}

func Debug() *VM {
	vm := &VM{
		debug: true,
		step:  make(chan struct{}, 0),
		curr:  make(chan int, 0),
	}
	return vm
}

func (vm *VM) Run(program *Program, env interface{}) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			f := &file.Error{
				Location: program.Locations[vm.ip-1],
				Message:  fmt.Sprintf("%v", r),
			}
			err = f.Bind(program.Source)
		}
	}()

	if vm.stack == nil {
		vm.stack = make([]interface{}, 0, 2)
	} else {
		vm.stack = vm.stack[0:0]
	}

	if vm.scopes != nil {
		vm.scopes = vm.scopes[0:0]
	}

	vm.memoryBudget = MemoryBudget
	vm.memory = 0
	vm.ip = 0

	for vm.ip < len(program.Bytecode) {
		if vm.debug {
			<-vm.step
		}

		op := program.Bytecode[vm.ip]
		arg := program.Arguments[vm.ip]
		vm.ip += 1

		switch op {

		case OpPush:
			vm.push(program.Constants[arg])

		case OpPop:
			vm.pop()

		case OpRot:
			b := vm.pop()
			a := vm.pop()
			vm.push(b)
			vm.push(a)

		case OpLoadConst:
			vm.push(runtime.Fetch(env, program.Constants[arg]))

		case OpLoadField:
			vm.push(runtime.FetchField(env, program.Constants[arg].(*runtime.Field)))

		case OpLoadFast:
			vm.push(env.(map[string]interface{})[program.Constants[arg].(string)])

		case OpLoadMethod:
			vm.push(runtime.FetchMethod(env, program.Constants[arg].(*runtime.Method)))

		case OpFetch:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Fetch(a, b))

		case OpFetchField:
			a := vm.pop()
			vm.push(runtime.FetchField(a, program.Constants[arg].(*runtime.Field)))

		case OpMethod:
			a := vm.pop()
			vm.push(runtime.FetchMethod(a, program.Constants[arg].(*runtime.Method)))

		case OpTrue:
			vm.push(true)

		case OpFalse:
			vm.push(false)

		case OpNil:
			vm.push(nil)

		case OpNegate:
			v := runtime.Negate(vm.pop())
			vm.push(v)

		case OpNot:
			v := vm.pop().(bool)
			vm.push(!v)

		case OpEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Equal(a, b))

		case OpEqualInt:
			b := vm.pop()
			a := vm.pop()
			vm.push(a.(int) == b.(int))

		case OpEqualString:
			b := vm.pop()
			a := vm.pop()
			vm.push(a.(string) == b.(string))

		case OpJump:
			vm.ip += arg

		case OpJumpIfTrue:
			if vm.current().(bool) {
				vm.ip += arg
			}

		case Op1:
			vm.push(1)
		case Op2:
			vm.push(2)
		case Op3:
			vm.push(3)
		case Op4:
			vm.push(4)
		case Op5:
			vm.push(5)
		case Op6:
			vm.push(6)
		case Op7:
			vm.push(7)
		case Op8:
			vm.push(8)
		case Op9:
			vm.push(9)
		case Op10:
			vm.push(10)
		case Op11:
			vm.push(11)
		case Op12:
			vm.push(12)
		case Op13:
			vm.push(13)
		case Op14:
			vm.push(14)
		case Op15:
			vm.push(15)
		case Op16:
			vm.push(16)
		case Op17:
			vm.push(17)
		case Op18:
			vm.push(18)
		case Op19:
			vm.push(19)
		case Op20:
			vm.push(20)
		case Op21:
			vm.push(21)
		case Op22:
			vm.push(22)
		case Op23:
			vm.push(23)
		case Op24:
			vm.push(24)
		case Op25:
			vm.push(25)
		case Op26:
			vm.push(26)
		case Op27:
			vm.push(27)
		case Op28:
			vm.push(28)
		case Op29:
			vm.push(29)
		case Op30:
			vm.push(30)
		case Op31:
			vm.push(31)
		case Op32:
			vm.push(32)
		case Op33:
			vm.push(33)
		case Op34:
			vm.push(34)
		case Op35:
			vm.push(35)
		case Op36:
			vm.push(36)
		case Op37:
			vm.push(37)
		case Op38:
			vm.push(38)
		case Op39:
			vm.push(39)
		case Op40:
			vm.push(40)
		case Op41:
			vm.push(41)
		case Op42:
			vm.push(42)
		case Op43:
			vm.push(43)
		case Op44:
			vm.push(44)
		case Op45:
			vm.push(45)
		case Op46:
			vm.push(46)
		case Op47:
			vm.push(47)
		case Op48:
			vm.push(48)
		case Op49:
			vm.push(49)
		case Op50:
			vm.push(50)
		case Op51:
			vm.push(51)
		case Op52:
			vm.push(52)
		case Op53:
			vm.push(53)
		case Op54:
			vm.push(54)
		case Op55:
			vm.push(55)
		case Op56:
			vm.push(56)
		case Op57:
			vm.push(57)
		case Op58:
			vm.push(58)
		case Op59:
			vm.push(59)
		case Op60:
			vm.push(60)
		case Op61:
			vm.push(61)
		case Op62:
			vm.push(62)
		case Op63:
			vm.push(63)
		case Op64:
			vm.push(64)
		case Op65:
			vm.push(65)
		case Op66:
			vm.push(66)
		case Op67:
			vm.push(67)
		case Op68:
			vm.push(68)
		case Op69:
			vm.push(69)
		case Op70:
			vm.push(70)
		case Op71:
			vm.push(71)
		case Op72:
			vm.push(72)
		case Op73:
			vm.push(73)
		case Op74:
			vm.push(74)
		case Op75:
			vm.push(75)
		case Op76:
			vm.push(76)
		case Op77:
			vm.push(77)
		case Op78:
			vm.push(78)
		case Op79:
			vm.push(79)
		case Op80:
			vm.push(80)
		case Op81:
			vm.push(81)
		case Op82:
			vm.push(82)
		case Op83:
			vm.push(83)
		case Op84:
			vm.push(84)
		case Op85:
			vm.push(85)
		case Op86:
			vm.push(86)
		case Op87:
			vm.push(87)
		case Op88:
			vm.push(88)
		case Op89:
			vm.push(89)
		case Op90:
			vm.push(90)
		case Op91:
			vm.push(91)
		case Op92:
			vm.push(92)
		case Op93:
			vm.push(93)
		case Op94:
			vm.push(94)
		case Op95:
			vm.push(95)
		case Op96:
			vm.push(96)
		case Op97:
			vm.push(97)
		case Op98:
			vm.push(98)
		case Op99:
			vm.push(99)
		case Op100:
			vm.push(100)

		case OpJumpIfFalse:
			if !vm.current().(bool) {
				vm.ip += arg
			}

		case OpJumpIfNil:
			if runtime.IsNil(vm.current()) {
				vm.ip += arg
			}

		case OpJumpIfEnd:
			scope := vm.Scope()
			if scope.It >= scope.Len {
				vm.ip += arg
			}

		case OpJumpBackward:
			vm.ip -= arg

		case OpIn:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.In(a, b))

		case OpLess:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Less(a, b))

		case OpMore:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.More(a, b))

		case OpLessOrEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.LessOrEqual(a, b))

		case OpMoreOrEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.MoreOrEqual(a, b))

		case OpAdd:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Add(a, b))

		case OpSubtract:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Subtract(a, b))

		case OpMultiply:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Multiply(a, b))

		case OpDivide:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Divide(a, b))

		case OpModulo:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Modulo(a, b))

		case OpExponent:
			b := vm.pop()
			a := vm.pop()
			vm.push(runtime.Exponent(a, b))

		case OpRange:
			b := vm.pop()
			a := vm.pop()
			min := runtime.ToInt(a)
			max := runtime.ToInt(b)
			size := max - min + 1
			if vm.memory+size >= vm.memoryBudget {
				panic("memory budget exceeded")
			}
			vm.push(runtime.MakeRange(min, max))
			vm.memory += size

		case OpMatches:
			b := vm.pop()
			a := vm.pop()
			match, err := regexp.MatchString(b.(string), a.(string))
			if err != nil {
				panic(err)
			}

			vm.push(match)

		case OpMatchesConst:
			a := vm.pop()
			r := program.Constants[arg].(*regexp.Regexp)
			vm.push(r.MatchString(a.(string)))

		case OpContains:
			b := vm.pop()
			a := vm.pop()
			vm.push(strings.Contains(a.(string), b.(string)))

		case OpStartsWith:
			b := vm.pop()
			a := vm.pop()
			vm.push(strings.HasPrefix(a.(string), b.(string)))

		case OpEndsWith:
			b := vm.pop()
			a := vm.pop()
			vm.push(strings.HasSuffix(a.(string), b.(string)))

		case OpSlice:
			from := vm.pop()
			to := vm.pop()
			node := vm.pop()
			vm.push(runtime.Slice(node, from, to))

		case OpCall:
			fn := reflect.ValueOf(vm.pop())
			size := arg
			in := make([]reflect.Value, size)
			for i := int(size) - 1; i >= 0; i-- {
				param := vm.pop()
				if param == nil && reflect.TypeOf(param) == nil {
					// In case of nil value and nil type use this hack,
					// otherwise reflect.Call will panic on zero value.
					in[i] = reflect.ValueOf(&param).Elem()
				} else {
					in[i] = reflect.ValueOf(param)
				}
			}
			out := fn.Call(in)
			if len(out) == 2 && !runtime.IsNil(out[1]) {
				panic(out[1].Interface().(error))
			}
			vm.push(out[0].Interface())

		case OpCallFast:
			fn := vm.pop().(func(...interface{}) interface{})
			size := arg
			in := make([]interface{}, size)
			for i := int(size) - 1; i >= 0; i-- {
				in[i] = vm.pop()
			}
			vm.push(fn(in...))

		case OpCallTyped:
			fn := vm.pop()
			out := vm.call(fn, arg)
			vm.push(out)

		case OpArray:
			size := vm.pop().(int)
			array := make([]interface{}, size)
			for i := size - 1; i >= 0; i-- {
				array[i] = vm.pop()
			}
			vm.push(array)
			vm.memory += size
			if vm.memory >= vm.memoryBudget {
				panic("memory budget exceeded")
			}

		case OpMap:
			size := vm.pop().(int)
			m := make(map[string]interface{})
			for i := size - 1; i >= 0; i-- {
				value := vm.pop()
				key := vm.pop()
				m[key.(string)] = value
			}
			vm.push(m)
			vm.memory += size
			if vm.memory >= vm.memoryBudget {
				panic("memory budget exceeded")
			}

		case OpLen:
			vm.push(runtime.Length(vm.current()))

		case OpCast:
			t := arg
			switch t {
			case 0:
				vm.push(runtime.ToInt(vm.pop()))
			case 1:
				vm.push(runtime.ToInt64(vm.pop()))
			case 2:
				vm.push(runtime.ToFloat64(vm.pop()))
			}

		case OpDeref:
			a := vm.pop()
			vm.push(runtime.Deref(a))

		case OpIncrementIt:
			scope := vm.Scope()
			scope.It++

		case OpIncrementCount:
			scope := vm.Scope()
			scope.Count++

		case OpGetCount:
			scope := vm.Scope()
			vm.push(scope.Count)

		case OpGetLen:
			scope := vm.Scope()
			vm.push(scope.Len)

		case OpPointer:
			scope := vm.Scope()
			vm.push(scope.Array.Index(scope.It).Interface())

		case OpBegin:
			a := vm.pop()
			array := reflect.ValueOf(a)
			vm.scopes = append(vm.scopes, &Scope{
				Array: array,
				Len:   array.Len(),
			})

		case OpEnd:
			vm.scopes = vm.scopes[:len(vm.scopes)-1]

		default:
			panic(fmt.Sprintf("unknown bytecode %#x", op))
		}

		if vm.debug {
			vm.curr <- vm.ip
		}
	}

	if vm.debug {
		close(vm.curr)
		close(vm.step)
	}

	if len(vm.stack) > 0 {
		return vm.pop(), nil
	}

	return nil, nil
}

func (vm *VM) push(value interface{}) {
	vm.stack = append(vm.stack, value)
}

func (vm *VM) current() interface{} {
	return vm.stack[len(vm.stack)-1]
}

func (vm *VM) pop() interface{} {
	value := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return value
}

func (vm *VM) Stack() []interface{} {
	return vm.stack
}

func (vm *VM) Scope() *Scope {
	if len(vm.scopes) > 0 {
		return vm.scopes[len(vm.scopes)-1]
	}
	return nil
}

func (vm *VM) Step() {
	vm.step <- struct{}{}
}

func (vm *VM) Position() chan int {
	return vm.curr
}
