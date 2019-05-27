// Code generated from Expr.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Expr
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseExprVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseExprVisitor) VisitStart(ctx *StartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitParenthesizedExpression(ctx *ParenthesizedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitAdditiveExpression(ctx *AdditiveExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitRelationalExpression(ctx *RelationalExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitTernaryExpression(ctx *TernaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitContainsExpression(ctx *ContainsExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitMatchesExpression(ctx *MatchesExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitMapLiteralExpression(ctx *MapLiteralExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitLiteralExpression(ctx *LiteralExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitInExpression(ctx *InExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitArrayLiteralExpression(ctx *ArrayLiteralExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitMemberDotExpression(ctx *MemberDotExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitUnaryExpression(ctx *UnaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitRangeExpression(ctx *RangeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitMemberIndexExpression(ctx *MemberIndexExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitIdentifierExpression(ctx *IdentifierExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitPointerExpression(ctx *PointerExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitLogicalExpression(ctx *LogicalExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitClosureMemberDotExpression(ctx *ClosureMemberDotExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitEqualityExpression(ctx *EqualityExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitBuiltinLiteralExpression(ctx *BuiltinLiteralExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitCallExpression(ctx *CallExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitLenBuiltinExpression(ctx *LenBuiltinExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitBuiltinExpression(ctx *BuiltinExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitClosureExpression(ctx *ClosureExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitArrayLiteral(ctx *ArrayLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitMapLiteral(ctx *MapLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitPropertyNameAndValueList(ctx *PropertyNameAndValueListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitPropertyAssignment(ctx *PropertyAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitPropertyName(ctx *PropertyNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitNilExpression(ctx *NilExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitBooleanExpression(ctx *BooleanExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitStringLiteralExpression(ctx *StringLiteralExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitIntegerExpression(ctx *IntegerExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitFloatExpression(ctx *FloatExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExprVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}
