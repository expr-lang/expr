package vm_test

import (
	"github.com/antonmedv/expr/vm"
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
			t.Errorf("cannot disassemble all opcodes")
		}
	}
}
