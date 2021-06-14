package vm

const (
	OpPush byte = iota
	OpPop
	OpRot
	OpFetch
	OpFetchNilSafe
	OpFetchMap
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
	OpIndex
	OpSlice
	OpProperty
	OpPropertyNilSafe
	OpCall
	OpCallFast
	OpMethod
	OpMethodNilSafe
	OpArray
	OpMap
	OpLen
	OpCast
	OpStore
	OpLoad
	OpInc
	OpNp      // Numpy builtin function
	OpNpUnary // Numpy builtin function for single vector input
	OpBegin
	OpEnd // This opcode must be at the end of this list.
)
