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

// EnterParenthesizedExpression is called when production ParenthesizedExpression is entered.
func (s *BaseExprListener) EnterParenthesizedExpression(ctx *ParenthesizedExpressionContext) {}

// ExitParenthesizedExpression is called when production ParenthesizedExpression is exited.
func (s *BaseExprListener) ExitParenthesizedExpression(ctx *ParenthesizedExpressionContext) {}

// EnterAdditiveExpression is called when production AdditiveExpression is entered.
func (s *BaseExprListener) EnterAdditiveExpression(ctx *AdditiveExpressionContext) {}

// ExitAdditiveExpression is called when production AdditiveExpression is exited.
func (s *BaseExprListener) ExitAdditiveExpression(ctx *AdditiveExpressionContext) {}

// EnterRelationalExpression is called when production RelationalExpression is entered.
func (s *BaseExprListener) EnterRelationalExpression(ctx *RelationalExpressionContext) {}

// ExitRelationalExpression is called when production RelationalExpression is exited.
func (s *BaseExprListener) ExitRelationalExpression(ctx *RelationalExpressionContext) {}

// EnterTernaryExpression is called when production TernaryExpression is entered.
func (s *BaseExprListener) EnterTernaryExpression(ctx *TernaryExpressionContext) {}

// ExitTernaryExpression is called when production TernaryExpression is exited.
func (s *BaseExprListener) ExitTernaryExpression(ctx *TernaryExpressionContext) {}

// EnterContainsExpression is called when production ContainsExpression is entered.
func (s *BaseExprListener) EnterContainsExpression(ctx *ContainsExpressionContext) {}

// ExitContainsExpression is called when production ContainsExpression is exited.
func (s *BaseExprListener) ExitContainsExpression(ctx *ContainsExpressionContext) {}

// EnterMatchesExpression is called when production MatchesExpression is entered.
func (s *BaseExprListener) EnterMatchesExpression(ctx *MatchesExpressionContext) {}

// ExitMatchesExpression is called when production MatchesExpression is exited.
func (s *BaseExprListener) ExitMatchesExpression(ctx *MatchesExpressionContext) {}

// EnterMapLiteralExpression is called when production MapLiteralExpression is entered.
func (s *BaseExprListener) EnterMapLiteralExpression(ctx *MapLiteralExpressionContext) {}

// ExitMapLiteralExpression is called when production MapLiteralExpression is exited.
func (s *BaseExprListener) ExitMapLiteralExpression(ctx *MapLiteralExpressionContext) {}

// EnterLiteralExpression is called when production LiteralExpression is entered.
func (s *BaseExprListener) EnterLiteralExpression(ctx *LiteralExpressionContext) {}

// ExitLiteralExpression is called when production LiteralExpression is exited.
func (s *BaseExprListener) ExitLiteralExpression(ctx *LiteralExpressionContext) {}

// EnterInExpression is called when production InExpression is entered.
func (s *BaseExprListener) EnterInExpression(ctx *InExpressionContext) {}

// ExitInExpression is called when production InExpression is exited.
func (s *BaseExprListener) ExitInExpression(ctx *InExpressionContext) {}

// EnterArrayLiteralExpression is called when production ArrayLiteralExpression is entered.
func (s *BaseExprListener) EnterArrayLiteralExpression(ctx *ArrayLiteralExpressionContext) {}

// ExitArrayLiteralExpression is called when production ArrayLiteralExpression is exited.
func (s *BaseExprListener) ExitArrayLiteralExpression(ctx *ArrayLiteralExpressionContext) {}

// EnterMemberDotExpression is called when production MemberDotExpression is entered.
func (s *BaseExprListener) EnterMemberDotExpression(ctx *MemberDotExpressionContext) {}

// ExitMemberDotExpression is called when production MemberDotExpression is exited.
func (s *BaseExprListener) ExitMemberDotExpression(ctx *MemberDotExpressionContext) {}

// EnterUnaryExpression is called when production UnaryExpression is entered.
func (s *BaseExprListener) EnterUnaryExpression(ctx *UnaryExpressionContext) {}

// ExitUnaryExpression is called when production UnaryExpression is exited.
func (s *BaseExprListener) ExitUnaryExpression(ctx *UnaryExpressionContext) {}

// EnterRangeExpression is called when production RangeExpression is entered.
func (s *BaseExprListener) EnterRangeExpression(ctx *RangeExpressionContext) {}

// ExitRangeExpression is called when production RangeExpression is exited.
func (s *BaseExprListener) ExitRangeExpression(ctx *RangeExpressionContext) {}

// EnterMemberIndexExpression is called when production MemberIndexExpression is entered.
func (s *BaseExprListener) EnterMemberIndexExpression(ctx *MemberIndexExpressionContext) {}

// ExitMemberIndexExpression is called when production MemberIndexExpression is exited.
func (s *BaseExprListener) ExitMemberIndexExpression(ctx *MemberIndexExpressionContext) {}

// EnterIdentifierExpression is called when production IdentifierExpression is entered.
func (s *BaseExprListener) EnterIdentifierExpression(ctx *IdentifierExpressionContext) {}

// ExitIdentifierExpression is called when production IdentifierExpression is exited.
func (s *BaseExprListener) ExitIdentifierExpression(ctx *IdentifierExpressionContext) {}

// EnterPointerExpression is called when production PointerExpression is entered.
func (s *BaseExprListener) EnterPointerExpression(ctx *PointerExpressionContext) {}

// ExitPointerExpression is called when production PointerExpression is exited.
func (s *BaseExprListener) ExitPointerExpression(ctx *PointerExpressionContext) {}

// EnterLogicalExpression is called when production LogicalExpression is entered.
func (s *BaseExprListener) EnterLogicalExpression(ctx *LogicalExpressionContext) {}

// ExitLogicalExpression is called when production LogicalExpression is exited.
func (s *BaseExprListener) ExitLogicalExpression(ctx *LogicalExpressionContext) {}

// EnterClosureMemberDotExpression is called when production ClosureMemberDotExpression is entered.
func (s *BaseExprListener) EnterClosureMemberDotExpression(ctx *ClosureMemberDotExpressionContext) {}

// ExitClosureMemberDotExpression is called when production ClosureMemberDotExpression is exited.
func (s *BaseExprListener) ExitClosureMemberDotExpression(ctx *ClosureMemberDotExpressionContext) {}

// EnterEndsWithExpression is called when production EndsWithExpression is entered.
func (s *BaseExprListener) EnterEndsWithExpression(ctx *EndsWithExpressionContext) {}

// ExitEndsWithExpression is called when production EndsWithExpression is exited.
func (s *BaseExprListener) ExitEndsWithExpression(ctx *EndsWithExpressionContext) {}

// EnterStartsWithExpression is called when production StartsWithExpression is entered.
func (s *BaseExprListener) EnterStartsWithExpression(ctx *StartsWithExpressionContext) {}

// ExitStartsWithExpression is called when production StartsWithExpression is exited.
func (s *BaseExprListener) ExitStartsWithExpression(ctx *StartsWithExpressionContext) {}

// EnterEqualityExpression is called when production EqualityExpression is entered.
func (s *BaseExprListener) EnterEqualityExpression(ctx *EqualityExpressionContext) {}

// ExitEqualityExpression is called when production EqualityExpression is exited.
func (s *BaseExprListener) ExitEqualityExpression(ctx *EqualityExpressionContext) {}

// EnterBuiltinLiteralExpression is called when production BuiltinLiteralExpression is entered.
func (s *BaseExprListener) EnterBuiltinLiteralExpression(ctx *BuiltinLiteralExpressionContext) {}

// ExitBuiltinLiteralExpression is called when production BuiltinLiteralExpression is exited.
func (s *BaseExprListener) ExitBuiltinLiteralExpression(ctx *BuiltinLiteralExpressionContext) {}

// EnterMultiplicativeExpression is called when production MultiplicativeExpression is entered.
func (s *BaseExprListener) EnterMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {}

// ExitMultiplicativeExpression is called when production MultiplicativeExpression is exited.
func (s *BaseExprListener) ExitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {}

// EnterCallExpression is called when production CallExpression is entered.
func (s *BaseExprListener) EnterCallExpression(ctx *CallExpressionContext) {}

// ExitCallExpression is called when production CallExpression is exited.
func (s *BaseExprListener) ExitCallExpression(ctx *CallExpressionContext) {}

// EnterLenBuiltinExpression is called when production LenBuiltinExpression is entered.
func (s *BaseExprListener) EnterLenBuiltinExpression(ctx *LenBuiltinExpressionContext) {}

// ExitLenBuiltinExpression is called when production LenBuiltinExpression is exited.
func (s *BaseExprListener) ExitLenBuiltinExpression(ctx *LenBuiltinExpressionContext) {}

// EnterBuiltinExpression is called when production BuiltinExpression is entered.
func (s *BaseExprListener) EnterBuiltinExpression(ctx *BuiltinExpressionContext) {}

// ExitBuiltinExpression is called when production BuiltinExpression is exited.
func (s *BaseExprListener) ExitBuiltinExpression(ctx *BuiltinExpressionContext) {}

// EnterClosureExpression is called when production ClosureExpression is entered.
func (s *BaseExprListener) EnterClosureExpression(ctx *ClosureExpressionContext) {}

// ExitClosureExpression is called when production ClosureExpression is exited.
func (s *BaseExprListener) ExitClosureExpression(ctx *ClosureExpressionContext) {}

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

// EnterNilExpression is called when production NilExpression is entered.
func (s *BaseExprListener) EnterNilExpression(ctx *NilExpressionContext) {}

// ExitNilExpression is called when production NilExpression is exited.
func (s *BaseExprListener) ExitNilExpression(ctx *NilExpressionContext) {}

// EnterBooleanExpression is called when production BooleanExpression is entered.
func (s *BaseExprListener) EnterBooleanExpression(ctx *BooleanExpressionContext) {}

// ExitBooleanExpression is called when production BooleanExpression is exited.
func (s *BaseExprListener) ExitBooleanExpression(ctx *BooleanExpressionContext) {}

// EnterStringLiteralExpression is called when production StringLiteralExpression is entered.
func (s *BaseExprListener) EnterStringLiteralExpression(ctx *StringLiteralExpressionContext) {}

// ExitStringLiteralExpression is called when production StringLiteralExpression is exited.
func (s *BaseExprListener) ExitStringLiteralExpression(ctx *StringLiteralExpressionContext) {}

// EnterIntegerExpression is called when production IntegerExpression is entered.
func (s *BaseExprListener) EnterIntegerExpression(ctx *IntegerExpressionContext) {}

// ExitIntegerExpression is called when production IntegerExpression is exited.
func (s *BaseExprListener) ExitIntegerExpression(ctx *IntegerExpressionContext) {}

// EnterFloatExpression is called when production FloatExpression is entered.
func (s *BaseExprListener) EnterFloatExpression(ctx *FloatExpressionContext) {}

// ExitFloatExpression is called when production FloatExpression is exited.
func (s *BaseExprListener) ExitFloatExpression(ctx *FloatExpressionContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseExprListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseExprListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterIntegerLiteral is called when production integerLiteral is entered.
func (s *BaseExprListener) EnterIntegerLiteral(ctx *IntegerLiteralContext) {}

// ExitIntegerLiteral is called when production integerLiteral is exited.
func (s *BaseExprListener) ExitIntegerLiteral(ctx *IntegerLiteralContext) {}
