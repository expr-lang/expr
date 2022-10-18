package vm

const (
	OpPush byte = iota
	OpPop
	OpRot
	OpFetch
	OpFetchField
	OpFetchEnv
	OpFetchEnvField
	OpFetchEnvFast
	OpTrue
	OpFalse
	OpNil
	OpNegate
	OpNot
	OpEqual
	OpEqualInt
	OpEqualString
	OpJump
	OpJumpIfTrue
	OpJumpIfFalse
	OpJumpIfNil
	OpJumpBackward
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
	OpRange
	OpMatches
	OpMatchesConst
	OpContains
	OpStartsWith
	OpEndsWith
	OpSlice
	OpCall
	OpCallFast
	OpMethod
	OpArray
	OpMap
	OpLen
	OpCast
	OpStore
	OpLoad
	OpInc
	OpDeref
	OpBegin
	OpEnd // This opcode must be at the end of this list.
)
