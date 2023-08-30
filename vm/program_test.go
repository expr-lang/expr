package vm_test

import (
	"strings"
	"testing"

	"github.com/antonmedv/expr/vm"
)

func TestProgram_Disassemble(t *testing.T) {
	for op := vm.OpPush; op < vm.OpEnd; op++ {
		program := vm.Program{
			Constants: []any{1, 2},
			Bytecode:  []vm.Opcode{op},
			Arguments: []int{1},
		}
		d := program.Disassemble()
		if strings.Contains(d, "(unknown)") {
			t.Errorf("cannot disassemble all opcodes")
		}
	}
}
