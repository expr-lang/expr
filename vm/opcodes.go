package vm

type Opcode byte

const (
	OpPush Opcode = iota
	OpPushInt
	OpPop
	OpRot
	OpLoadConst
	OpLoadField
	OpLoadFast
	OpLoadMethod
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
	OpCallTyped
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
	Op1
	Op2
	Op3
	Op4
	Op5
	Op6
	Op7
	Op8
	Op9
	Op10
	Op11
	Op12
	Op13
	Op14
	Op15
	Op16
	Op17
	Op18
	Op19
	Op20
	Op21
	Op22
	Op23
	Op24
	Op25
	Op26
	Op27
	Op28
	Op29
	Op30
	Op31
	Op32
	Op33
	Op34
	Op35
	Op36
	Op37
	Op38
	Op39
	Op40
	Op41
	Op42
	Op43
	Op44
	Op45
	Op46
	Op47
	Op48
	Op49
	Op50
	Op51
	Op52
	Op53
	Op54
	Op55
	Op56
	Op57
	Op58
	Op59
	Op60
	Op61
	Op62
	Op63
	Op64
	Op65
	Op66
	Op67
	Op68
	Op69
	Op70
	Op71
	Op72
	Op73
	Op74
	Op75
	Op76
	Op77
	Op78
	Op79
	Op80
	Op81
	Op82
	Op83
	Op84
	Op85
	Op86
	Op87
	Op88
	Op89
	Op90
	Op91
	Op92
	Op93
	Op94
	Op95
	Op96
	Op97
	Op98
	Op99
	Op100
)
