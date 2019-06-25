// Code generated from Expr.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Expr
import "github.com/antlr/antlr4/runtime/Go/antlr"

// ExprListener is a complete listener for a parse tree produced by ExprParser.
type ExprListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterCall is called when entering the Call production.
	EnterCall(c *CallContext)

	// EnterMatches is called when entering the Matches production.
	EnterMatches(c *MatchesContext)

	// EnterTernary is called when entering the Ternary production.
	EnterTernary(c *TernaryContext)

	// EnterPointer is called when entering the Pointer production.
	EnterPointer(c *PointerContext)

	// EnterString is called when entering the String production.
	EnterString(c *StringContext)

	// EnterClosureMemberDot is called when entering the ClosureMemberDot production.
	EnterClosureMemberDot(c *ClosureMemberDotContext)

	// EnterUnary is called when entering the Unary production.
	EnterUnary(c *UnaryContext)

	// EnterNil is called when entering the Nil production.
	EnterNil(c *NilContext)

	// EnterInteger is called when entering the Integer production.
	EnterInteger(c *IntegerContext)

	// EnterArray is called when entering the Array production.
	EnterArray(c *ArrayContext)

	// EnterFloat is called when entering the Float production.
	EnterFloat(c *FloatContext)

	// EnterIdentifier is called when entering the Identifier production.
	EnterIdentifier(c *IdentifierContext)

	// EnterParenthesized is called when entering the Parenthesized production.
	EnterParenthesized(c *ParenthesizedContext)

	// EnterSlice is called when entering the Slice production.
	EnterSlice(c *SliceContext)

	// EnterMemberIndex is called when entering the MemberIndex production.
	EnterMemberIndex(c *MemberIndexContext)

	// EnterBuiltinsList is called when entering the BuiltinsList production.
	EnterBuiltinsList(c *BuiltinsListContext)

	// EnterBinary is called when entering the Binary production.
	EnterBinary(c *BinaryContext)

	// EnterBoolean is called when entering the Boolean production.
	EnterBoolean(c *BooleanContext)

	// EnterMap is called when entering the Map production.
	EnterMap(c *MapContext)

	// EnterMemberDot is called when entering the MemberDot production.
	EnterMemberDot(c *MemberDotContext)

	// EnterBuiltinLen is called when entering the BuiltinLen production.
	EnterBuiltinLen(c *BuiltinLenContext)

	// EnterBuiltin is called when entering the Builtin production.
	EnterBuiltin(c *BuiltinContext)

	// EnterClosure is called when entering the closure production.
	EnterClosure(c *ClosureContext)

	// EnterArguments is called when entering the arguments production.
	EnterArguments(c *ArgumentsContext)

	// EnterArrayLiteral is called when entering the arrayLiteral production.
	EnterArrayLiteral(c *ArrayLiteralContext)

	// EnterMapLiteral is called when entering the mapLiteral production.
	EnterMapLiteral(c *MapLiteralContext)

	// EnterPropertyNameAndValueList is called when entering the propertyNameAndValueList production.
	EnterPropertyNameAndValueList(c *PropertyNameAndValueListContext)

	// EnterPropertyAssignment is called when entering the propertyAssignment production.
	EnterPropertyAssignment(c *PropertyAssignmentContext)

	// EnterPropertyName is called when entering the propertyName production.
	EnterPropertyName(c *PropertyNameContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitCall is called when exiting the Call production.
	ExitCall(c *CallContext)

	// ExitMatches is called when exiting the Matches production.
	ExitMatches(c *MatchesContext)

	// ExitTernary is called when exiting the Ternary production.
	ExitTernary(c *TernaryContext)

	// ExitPointer is called when exiting the Pointer production.
	ExitPointer(c *PointerContext)

	// ExitString is called when exiting the String production.
	ExitString(c *StringContext)

	// ExitClosureMemberDot is called when exiting the ClosureMemberDot production.
	ExitClosureMemberDot(c *ClosureMemberDotContext)

	// ExitUnary is called when exiting the Unary production.
	ExitUnary(c *UnaryContext)

	// ExitNil is called when exiting the Nil production.
	ExitNil(c *NilContext)

	// ExitInteger is called when exiting the Integer production.
	ExitInteger(c *IntegerContext)

	// ExitArray is called when exiting the Array production.
	ExitArray(c *ArrayContext)

	// ExitFloat is called when exiting the Float production.
	ExitFloat(c *FloatContext)

	// ExitIdentifier is called when exiting the Identifier production.
	ExitIdentifier(c *IdentifierContext)

	// ExitParenthesized is called when exiting the Parenthesized production.
	ExitParenthesized(c *ParenthesizedContext)

	// ExitSlice is called when exiting the Slice production.
	ExitSlice(c *SliceContext)

	// ExitMemberIndex is called when exiting the MemberIndex production.
	ExitMemberIndex(c *MemberIndexContext)

	// ExitBuiltinsList is called when exiting the BuiltinsList production.
	ExitBuiltinsList(c *BuiltinsListContext)

	// ExitBinary is called when exiting the Binary production.
	ExitBinary(c *BinaryContext)

	// ExitBoolean is called when exiting the Boolean production.
	ExitBoolean(c *BooleanContext)

	// ExitMap is called when exiting the Map production.
	ExitMap(c *MapContext)

	// ExitMemberDot is called when exiting the MemberDot production.
	ExitMemberDot(c *MemberDotContext)

	// ExitBuiltinLen is called when exiting the BuiltinLen production.
	ExitBuiltinLen(c *BuiltinLenContext)

	// ExitBuiltin is called when exiting the Builtin production.
	ExitBuiltin(c *BuiltinContext)

	// ExitClosure is called when exiting the closure production.
	ExitClosure(c *ClosureContext)

	// ExitArguments is called when exiting the arguments production.
	ExitArguments(c *ArgumentsContext)

	// ExitArrayLiteral is called when exiting the arrayLiteral production.
	ExitArrayLiteral(c *ArrayLiteralContext)

	// ExitMapLiteral is called when exiting the mapLiteral production.
	ExitMapLiteral(c *MapLiteralContext)

	// ExitPropertyNameAndValueList is called when exiting the propertyNameAndValueList production.
	ExitPropertyNameAndValueList(c *PropertyNameAndValueListContext)

	// ExitPropertyAssignment is called when exiting the propertyAssignment production.
	ExitPropertyAssignment(c *PropertyAssignmentContext)

	// ExitPropertyName is called when exiting the propertyName production.
	ExitPropertyName(c *PropertyNameContext)
}
