package vm

const (
	OpPush byte = iota
	OpPop
	OpLoad
	OpFetch
	OpTrue
	OpFalse
	OpNil
	OpNegate
	OpNot
	OpEqual
	OpNotEqual
	OpJumpIfTrue
	OpJumpIfFalse
)
