// Code generated from Expr.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Expr
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ExprParser.
type ExprVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ExprParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by ExprParser#ParenthesizedExpression.
	VisitParenthesizedExpression(ctx *ParenthesizedExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#AdditiveExpression.
	VisitAdditiveExpression(ctx *AdditiveExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#RelationalExpression.
	VisitRelationalExpression(ctx *RelationalExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#TernaryExpression.
	VisitTernaryExpression(ctx *TernaryExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#ContainsExpression.
	VisitContainsExpression(ctx *ContainsExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#MatchesExpression.
	VisitMatchesExpression(ctx *MatchesExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#MapLiteralExpression.
	VisitMapLiteralExpression(ctx *MapLiteralExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#LiteralExpression.
	VisitLiteralExpression(ctx *LiteralExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#InExpression.
	VisitInExpression(ctx *InExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#ArrayLiteralExpression.
	VisitArrayLiteralExpression(ctx *ArrayLiteralExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#MemberDotExpression.
	VisitMemberDotExpression(ctx *MemberDotExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#UnaryExpression.
	VisitUnaryExpression(ctx *UnaryExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#RangeExpression.
	VisitRangeExpression(ctx *RangeExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#MemberIndexExpression.
	VisitMemberIndexExpression(ctx *MemberIndexExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#IdentifierExpression.
	VisitIdentifierExpression(ctx *IdentifierExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#PointerExpression.
	VisitPointerExpression(ctx *PointerExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#LogicalExpression.
	VisitLogicalExpression(ctx *LogicalExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#ClosureMemberDotExpression.
	VisitClosureMemberDotExpression(ctx *ClosureMemberDotExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#EndsWithExpression.
	VisitEndsWithExpression(ctx *EndsWithExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#StartsWithExpression.
	VisitStartsWithExpression(ctx *StartsWithExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#EqualityExpression.
	VisitEqualityExpression(ctx *EqualityExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#BuiltinLiteralExpression.
	VisitBuiltinLiteralExpression(ctx *BuiltinLiteralExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#MultiplicativeExpression.
	VisitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#CallExpression.
	VisitCallExpression(ctx *CallExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#LenBuiltinExpression.
	VisitLenBuiltinExpression(ctx *LenBuiltinExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#BuiltinExpression.
	VisitBuiltinExpression(ctx *BuiltinExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#ClosureExpression.
	VisitClosureExpression(ctx *ClosureExpressionContext) interface{}

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

	// Visit a parse tree produced by ExprParser#NilExpression.
	VisitNilExpression(ctx *NilExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#BooleanExpression.
	VisitBooleanExpression(ctx *BooleanExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#StringLiteralExpression.
	VisitStringLiteralExpression(ctx *StringLiteralExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#IntegerExpression.
	VisitIntegerExpression(ctx *IntegerExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#FloatExpression.
	VisitFloatExpression(ctx *FloatExpressionContext) interface{}

	// Visit a parse tree produced by ExprParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by ExprParser#integerLiteral.
	VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{}
}
