package vm

const (
	OpPush byte = iota
	OpPop
	OpConst
	OpFetch
	OpFetchMap
	OpTrue
	OpFalse
	OpNil
	OpNegate
	OpNot
	OpEqual
	OpJump
	OpJumpIfTrue
	OpJumpIfFalse
	OpJumpBackward
	OpIn
	OpLess
	OpMore
	OpLessOrEqual
	OpMoreOrEqual
	OpAdd
	OpInc
	OpSubtract
	OpMultiply
	OpDivide
	OpModulo
	OpExponent
	OpContains
	OpRange
	OpMatches
	OpMatchesConst
	OpIndex
	OpProperty
	OpCall
	OpMethod
	OpArray
	OpMap
	OpLen
	OpBegin
	OpEnd
	OpStore
	OpLoad
)
