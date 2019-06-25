// Code generated from Expr.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Expr
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExprListener is a complete listener for a parse tree produced by ExprParser.
type BaseExprListener struct{}

var _ ExprListener = &BaseExprListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExprListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExprListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExprListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExprListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseExprListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseExprListener) ExitStart(ctx *StartContext) {}

// EnterCall is called when production Call is entered.
func (s *BaseExprListener) EnterCall(ctx *CallContext) {}

// ExitCall is called when production Call is exited.
func (s *BaseExprListener) ExitCall(ctx *CallContext) {}

// EnterMatches is called when production Matches is entered.
func (s *BaseExprListener) EnterMatches(ctx *MatchesContext) {}

// ExitMatches is called when production Matches is exited.
func (s *BaseExprListener) ExitMatches(ctx *MatchesContext) {}

// EnterTernary is called when production Ternary is entered.
func (s *BaseExprListener) EnterTernary(ctx *TernaryContext) {}

// ExitTernary is called when production Ternary is exited.
func (s *BaseExprListener) ExitTernary(ctx *TernaryContext) {}

// EnterPointer is called when production Pointer is entered.
func (s *BaseExprListener) EnterPointer(ctx *PointerContext) {}

// ExitPointer is called when production Pointer is exited.
func (s *BaseExprListener) ExitPointer(ctx *PointerContext) {}

// EnterString is called when production String is entered.
func (s *BaseExprListener) EnterString(ctx *StringContext) {}

// ExitString is called when production String is exited.
func (s *BaseExprListener) ExitString(ctx *StringContext) {}

// EnterClosureMemberDot is called when production ClosureMemberDot is entered.
func (s *BaseExprListener) EnterClosureMemberDot(ctx *ClosureMemberDotContext) {}

// ExitClosureMemberDot is called when production ClosureMemberDot is exited.
func (s *BaseExprListener) ExitClosureMemberDot(ctx *ClosureMemberDotContext) {}

// EnterUnary is called when production Unary is entered.
func (s *BaseExprListener) EnterUnary(ctx *UnaryContext) {}

// ExitUnary is called when production Unary is exited.
func (s *BaseExprListener) ExitUnary(ctx *UnaryContext) {}

// EnterNil is called when production Nil is entered.
func (s *BaseExprListener) EnterNil(ctx *NilContext) {}

// ExitNil is called when production Nil is exited.
func (s *BaseExprListener) ExitNil(ctx *NilContext) {}

// EnterInteger is called when production Integer is entered.
func (s *BaseExprListener) EnterInteger(ctx *IntegerContext) {}

// ExitInteger is called when production Integer is exited.
func (s *BaseExprListener) ExitInteger(ctx *IntegerContext) {}

// EnterArray is called when production Array is entered.
func (s *BaseExprListener) EnterArray(ctx *ArrayContext) {}

// ExitArray is called when production Array is exited.
func (s *BaseExprListener) ExitArray(ctx *ArrayContext) {}

// EnterFloat is called when production Float is entered.
func (s *BaseExprListener) EnterFloat(ctx *FloatContext) {}

// ExitFloat is called when production Float is exited.
func (s *BaseExprListener) ExitFloat(ctx *FloatContext) {}

// EnterIdentifier is called when production Identifier is entered.
func (s *BaseExprListener) EnterIdentifier(ctx *IdentifierContext) {}

// ExitIdentifier is called when production Identifier is exited.
func (s *BaseExprListener) ExitIdentifier(ctx *IdentifierContext) {}

// EnterParenthesized is called when production Parenthesized is entered.
func (s *BaseExprListener) EnterParenthesized(ctx *ParenthesizedContext) {}

// ExitParenthesized is called when production Parenthesized is exited.
func (s *BaseExprListener) ExitParenthesized(ctx *ParenthesizedContext) {}

// EnterSlice is called when production Slice is entered.
func (s *BaseExprListener) EnterSlice(ctx *SliceContext) {}

// ExitSlice is called when production Slice is exited.
func (s *BaseExprListener) ExitSlice(ctx *SliceContext) {}

// EnterMemberIndex is called when production MemberIndex is entered.
func (s *BaseExprListener) EnterMemberIndex(ctx *MemberIndexContext) {}

// ExitMemberIndex is called when production MemberIndex is exited.
func (s *BaseExprListener) ExitMemberIndex(ctx *MemberIndexContext) {}

// EnterBuiltinsList is called when production BuiltinsList is entered.
func (s *BaseExprListener) EnterBuiltinsList(ctx *BuiltinsListContext) {}

// ExitBuiltinsList is called when production BuiltinsList is exited.
func (s *BaseExprListener) ExitBuiltinsList(ctx *BuiltinsListContext) {}

// EnterBinary is called when production Binary is entered.
func (s *BaseExprListener) EnterBinary(ctx *BinaryContext) {}

// ExitBinary is called when production Binary is exited.
func (s *BaseExprListener) ExitBinary(ctx *BinaryContext) {}

// EnterBoolean is called when production Boolean is entered.
func (s *BaseExprListener) EnterBoolean(ctx *BooleanContext) {}

// ExitBoolean is called when production Boolean is exited.
func (s *BaseExprListener) ExitBoolean(ctx *BooleanContext) {}

// EnterMap is called when production Map is entered.
func (s *BaseExprListener) EnterMap(ctx *MapContext) {}

// ExitMap is called when production Map is exited.
func (s *BaseExprListener) ExitMap(ctx *MapContext) {}

// EnterMemberDot is called when production MemberDot is entered.
func (s *BaseExprListener) EnterMemberDot(ctx *MemberDotContext) {}

// ExitMemberDot is called when production MemberDot is exited.
func (s *BaseExprListener) ExitMemberDot(ctx *MemberDotContext) {}

// EnterBuiltinLen is called when production BuiltinLen is entered.
func (s *BaseExprListener) EnterBuiltinLen(ctx *BuiltinLenContext) {}

// ExitBuiltinLen is called when production BuiltinLen is exited.
func (s *BaseExprListener) ExitBuiltinLen(ctx *BuiltinLenContext) {}

// EnterBuiltin is called when production Builtin is entered.
func (s *BaseExprListener) EnterBuiltin(ctx *BuiltinContext) {}

// ExitBuiltin is called when production Builtin is exited.
func (s *BaseExprListener) ExitBuiltin(ctx *BuiltinContext) {}

// EnterClosure is called when production closure is entered.
func (s *BaseExprListener) EnterClosure(ctx *ClosureContext) {}

// ExitClosure is called when production closure is exited.
func (s *BaseExprListener) ExitClosure(ctx *ClosureContext) {}

// EnterArguments is called when production arguments is entered.
func (s *BaseExprListener) EnterArguments(ctx *ArgumentsContext) {}

// ExitArguments is called when production arguments is exited.
func (s *BaseExprListener) ExitArguments(ctx *ArgumentsContext) {}

// EnterArrayLiteral is called when production arrayLiteral is entered.
func (s *BaseExprListener) EnterArrayLiteral(ctx *ArrayLiteralContext) {}

// ExitArrayLiteral is called when production arrayLiteral is exited.
func (s *BaseExprListener) ExitArrayLiteral(ctx *ArrayLiteralContext) {}

// EnterMapLiteral is called when production mapLiteral is entered.
func (s *BaseExprListener) EnterMapLiteral(ctx *MapLiteralContext) {}

// ExitMapLiteral is called when production mapLiteral is exited.
func (s *BaseExprListener) ExitMapLiteral(ctx *MapLiteralContext) {}

// EnterPropertyNameAndValueList is called when production propertyNameAndValueList is entered.
func (s *BaseExprListener) EnterPropertyNameAndValueList(ctx *PropertyNameAndValueListContext) {}

// ExitPropertyNameAndValueList is called when production propertyNameAndValueList is exited.
func (s *BaseExprListener) ExitPropertyNameAndValueList(ctx *PropertyNameAndValueListContext) {}

// EnterPropertyAssignment is called when production propertyAssignment is entered.
func (s *BaseExprListener) EnterPropertyAssignment(ctx *PropertyAssignmentContext) {}

// ExitPropertyAssignment is called when production propertyAssignment is exited.
func (s *BaseExprListener) ExitPropertyAssignment(ctx *PropertyAssignmentContext) {}

// EnterPropertyName is called when production propertyName is entered.
func (s *BaseExprListener) EnterPropertyName(ctx *PropertyNameContext) {}

// ExitPropertyName is called when production propertyName is exited.
func (s *BaseExprListener) ExitPropertyName(ctx *PropertyNameContext) {}
