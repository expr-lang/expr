package vm

import "fmt"

var FuncTypes = []interface{}{
	new(byte), // The zero element represents no typed func.
	new(func() bool),
	new(func() int),
	new(func() uint),
	new(func() int64),
	new(func() uint64),
	new(func() float64),
	new(func() string),
	new(func() interface{}),
}

func (vm *VM) call(fn interface{}, kind int) interface{} {
	switch kind {
	case 1:
		return fn.(func() bool)()
	case 2:
		return fn.(func() int)()
	case 3:
		return fn.(func() uint)()
	case 4:
		return fn.(func() int64)()
	case 5:
		return fn.(func() uint64)()
	case 6:
		return fn.(func() float64)()
	case 7:
		return fn.(func() string)()
	case 8:
		return fn.(func() interface{})()
	}
	panic(fmt.Sprintf("unknown function kind (%v)", kind))
}
