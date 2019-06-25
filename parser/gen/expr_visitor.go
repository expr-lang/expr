// Code generated from Expr.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Expr
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ExprParser.
type ExprVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ExprParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by ExprParser#Call.
	VisitCall(ctx *CallContext) interface{}

	// Visit a parse tree produced by ExprParser#Matches.
	VisitMatches(ctx *MatchesContext) interface{}

	// Visit a parse tree produced by ExprParser#Ternary.
	VisitTernary(ctx *TernaryContext) interface{}

	// Visit a parse tree produced by ExprParser#Pointer.
	VisitPointer(ctx *PointerContext) interface{}

	// Visit a parse tree produced by ExprParser#String.
	VisitString(ctx *StringContext) interface{}

	// Visit a parse tree produced by ExprParser#ClosureMemberDot.
	VisitClosureMemberDot(ctx *ClosureMemberDotContext) interface{}

	// Visit a parse tree produced by ExprParser#Unary.
	VisitUnary(ctx *UnaryContext) interface{}

	// Visit a parse tree produced by ExprParser#Nil.
	VisitNil(ctx *NilContext) interface{}

	// Visit a parse tree produced by ExprParser#Integer.
	VisitInteger(ctx *IntegerContext) interface{}

	// Visit a parse tree produced by ExprParser#Array.
	VisitArray(ctx *ArrayContext) interface{}

	// Visit a parse tree produced by ExprParser#Float.
	VisitFloat(ctx *FloatContext) interface{}

	// Visit a parse tree produced by ExprParser#Identifier.
	VisitIdentifier(ctx *IdentifierContext) interface{}

	// Visit a parse tree produced by ExprParser#Parenthesized.
	VisitParenthesized(ctx *ParenthesizedContext) interface{}

	// Visit a parse tree produced by ExprParser#Slice.
	VisitSlice(ctx *SliceContext) interface{}

	// Visit a parse tree produced by ExprParser#MemberIndex.
	VisitMemberIndex(ctx *MemberIndexContext) interface{}

	// Visit a parse tree produced by ExprParser#BuiltinsList.
	VisitBuiltinsList(ctx *BuiltinsListContext) interface{}

	// Visit a parse tree produced by ExprParser#Binary.
	VisitBinary(ctx *BinaryContext) interface{}

	// Visit a parse tree produced by ExprParser#Boolean.
	VisitBoolean(ctx *BooleanContext) interface{}

	// Visit a parse tree produced by ExprParser#Map.
	VisitMap(ctx *MapContext) interface{}

	// Visit a parse tree produced by ExprParser#MemberDot.
	VisitMemberDot(ctx *MemberDotContext) interface{}

	// Visit a parse tree produced by ExprParser#BuiltinLen.
	VisitBuiltinLen(ctx *BuiltinLenContext) interface{}

	// Visit a parse tree produced by ExprParser#Builtin.
	VisitBuiltin(ctx *BuiltinContext) interface{}

	// Visit a parse tree produced by ExprParser#closure.
	VisitClosure(ctx *ClosureContext) interface{}

	// Visit a parse tree produced by ExprParser#arguments.
	VisitArguments(ctx *ArgumentsContext) interface{}

	// Visit a parse tree produced by ExprParser#arrayLiteral.
	VisitArrayLiteral(ctx *ArrayLiteralContext) interface{}

	// Visit a parse tree produced by ExprParser#mapLiteral.
	VisitMapLiteral(ctx *MapLiteralContext) interface{}

	// Visit a parse tree produced by ExprParser#propertyNameAndValueList.
	VisitPropertyNameAndValueList(ctx *PropertyNameAndValueListContext) interface{}

	// Visit a parse tree produced by ExprParser#propertyAssignment.
	VisitPropertyAssignment(ctx *PropertyAssignmentContext) interface{}

	// Visit a parse tree produced by ExprParser#propertyName.
	VisitPropertyName(ctx *PropertyNameContext) interface{}
}
