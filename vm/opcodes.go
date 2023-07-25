package vm

type Opcode byte

const (
	OpInvalid Opcode = iota
	OpPush
	OpPushInt
	OpPop
	OpLoadConst
	OpLoadField
	OpLoadFast
	OpLoadMethod
	OpLoadFunc
	OpLoadEnv
	OpFetch
	OpFetchField
	OpMethod
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
	OpJumpIfNotNil
	OpJumpIfSaveCommon // if common value save skip common calculate
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
	OpCall0
	OpCall1
	OpCall2
	OpCall3
	OpCallN
	OpCallFast
	OpCallTyped
	OpBuiltin
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
	OpSaveCommon // save common sub-expr value
	OpLoadCommon // load common sub-expr value
	OpBegin
	OpEnd // This opcode must be at the end of this list.
)
