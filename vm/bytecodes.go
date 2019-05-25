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
	OpJumpIfTrue
	OpJumpIfFalse
	OpContains
	OpLess
	OpMore
	OpLessOrEqual
	OpMoreOrEqual
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpModulo
	OpExponent
)
