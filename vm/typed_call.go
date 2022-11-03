package vm

import "fmt"

var FuncTypes = []interface{}{
	new(bool), // The zero element represents no typed func.
	new(func() bool),
	new(func() string),
}

func (vm *VM) call(fn interface{}, kind int) interface{} {
	switch kind {
	case 1:
		return fn.(func() bool)()
	case 2:
		return fn.(func() string)()
	}
	panic(fmt.Sprintf("unknown function kind (%v)", kind))
}
