package vm

import (
	"sync"
)

type VMPool struct {
	pool sync.Pool
}

func NewVMPool() *VMPool {
	return &VMPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &VM{}
			},
		},
	}
}

func (vp *VMPool) Run(program *Program, env interface{}) (out interface{}, err error) {
	vm := vp.pool.Get().(*VM)
	defer vp.pool.Put(vm)

	return vm.Run(program, env)
}
