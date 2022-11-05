package vm

type Opcode byte

const (
	OpPush Opcode = iota
	OpPushInt
	OpPop
	OpRot
	OpFetch
	OpFetchField
	OpFetchEnv
	OpFetchEnvField
	OpFetchEnvFast
	OpFetcher
	OpFetcherEnv
	OpMethod
	OpMethodEnv
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
	OpJumpIfEnd
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
	OpArray
	OpMap
	OpLen
	OpCast
	OpDeref
	OpIncrementIt
	OpIncrementCount
	OpGetCount
	OpGetLen
	OpPointer
	OpBegin
	OpEnd // This opcode must be at the end of this list.
)
