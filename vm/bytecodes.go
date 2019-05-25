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
	OpJump
	OpJumpIfTrue
	OpJumpIfFalse
	OpIn
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
	OpContains
	OpRange
	OpMatches
	OpMatchesConst
	OpField
	OpFieldConst
	OpCall
	OpMethod
	OpArray
	OpMap
)
