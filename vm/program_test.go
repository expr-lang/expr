package vm_test

import (
	"gopkg.in/antonmedv/expr.v2/vm"
	"strings"
	"testing"
)

func TestProgram_Disassemble(t *testing.T) {
	for op := vm.OpPush; op < vm.OpLoad; op++ {
		program := vm.Program{
			Constants: []interface{}{true},
			Bytecode:  []byte{op},
		}
		d := program.Disassemble()
		if strings.Contains(d, "\t0x") {
			t.Errorf("cannot disassemble opcode %#x", op)
		}
	}
}
