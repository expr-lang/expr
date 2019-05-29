// Code generated from Expr.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Expr
import "github.com/antlr/antlr4/runtime/Go/antlr"

// ExprListener is a complete listener for a parse tree produced by ExprParser.
type ExprListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterParenthesizedExpression is called when entering the ParenthesizedExpression production.
	EnterParenthesizedExpression(c *ParenthesizedExpressionContext)

	// EnterAdditiveExpression is called when entering the AdditiveExpression production.
	EnterAdditiveExpression(c *AdditiveExpressionContext)

	// EnterRelationalExpression is called when entering the RelationalExpression production.
	EnterRelationalExpression(c *RelationalExpressionContext)

	// EnterTernaryExpression is called when entering the TernaryExpression production.
	EnterTernaryExpression(c *TernaryExpressionContext)

	// EnterContainsExpression is called when entering the ContainsExpression production.
	EnterContainsExpression(c *ContainsExpressionContext)

	// EnterMatchesExpression is called when entering the MatchesExpression production.
	EnterMatchesExpression(c *MatchesExpressionContext)

	// EnterMapLiteralExpression is called when entering the MapLiteralExpression production.
	EnterMapLiteralExpression(c *MapLiteralExpressionContext)

	// EnterLiteralExpression is called when entering the LiteralExpression production.
	EnterLiteralExpression(c *LiteralExpressionContext)

	// EnterInExpression is called when entering the InExpression production.
	EnterInExpression(c *InExpressionContext)

	// EnterArrayLiteralExpression is called when entering the ArrayLiteralExpression production.
	EnterArrayLiteralExpression(c *ArrayLiteralExpressionContext)

	// EnterMemberDotExpression is called when entering the MemberDotExpression production.
	EnterMemberDotExpression(c *MemberDotExpressionContext)

	// EnterUnaryExpression is called when entering the UnaryExpression production.
	EnterUnaryExpression(c *UnaryExpressionContext)

	// EnterRangeExpression is called when entering the RangeExpression production.
	EnterRangeExpression(c *RangeExpressionContext)

	// EnterMemberIndexExpression is called when entering the MemberIndexExpression production.
	EnterMemberIndexExpression(c *MemberIndexExpressionContext)

	// EnterIdentifierExpression is called when entering the IdentifierExpression production.
	EnterIdentifierExpression(c *IdentifierExpressionContext)

	// EnterPointerExpression is called when entering the PointerExpression production.
	EnterPointerExpression(c *PointerExpressionContext)

	// EnterLogicalExpression is called when entering the LogicalExpression production.
	EnterLogicalExpression(c *LogicalExpressionContext)

	// EnterClosureMemberDotExpression is called when entering the ClosureMemberDotExpression production.
	EnterClosureMemberDotExpression(c *ClosureMemberDotExpressionContext)

	// EnterEndsWithExpression is called when entering the EndsWithExpression production.
	EnterEndsWithExpression(c *EndsWithExpressionContext)

	// EnterStartsWithExpression is called when entering the StartsWithExpression production.
	EnterStartsWithExpression(c *StartsWithExpressionContext)

	// EnterEqualityExpression is called when entering the EqualityExpression production.
	EnterEqualityExpression(c *EqualityExpressionContext)

	// EnterBuiltinLiteralExpression is called when entering the BuiltinLiteralExpression production.
	EnterBuiltinLiteralExpression(c *BuiltinLiteralExpressionContext)

	// EnterMultiplicativeExpression is called when entering the MultiplicativeExpression production.
	EnterMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// EnterCallExpression is called when entering the CallExpression production.
	EnterCallExpression(c *CallExpressionContext)

	// EnterLenBuiltinExpression is called when entering the LenBuiltinExpression production.
	EnterLenBuiltinExpression(c *LenBuiltinExpressionContext)

	// EnterBuiltinExpression is called when entering the BuiltinExpression production.
	EnterBuiltinExpression(c *BuiltinExpressionContext)

	// EnterClosureExpression is called when entering the ClosureExpression production.
	EnterClosureExpression(c *ClosureExpressionContext)

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

	// EnterNilExpression is called when entering the NilExpression production.
	EnterNilExpression(c *NilExpressionContext)

	// EnterBooleanExpression is called when entering the BooleanExpression production.
	EnterBooleanExpression(c *BooleanExpressionContext)

	// EnterStringLiteralExpression is called when entering the StringLiteralExpression production.
	EnterStringLiteralExpression(c *StringLiteralExpressionContext)

	// EnterIntegerExpression is called when entering the IntegerExpression production.
	EnterIntegerExpression(c *IntegerExpressionContext)

	// EnterFloatExpression is called when entering the FloatExpression production.
	EnterFloatExpression(c *FloatExpressionContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterIntegerLiteral is called when entering the integerLiteral production.
	EnterIntegerLiteral(c *IntegerLiteralContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitParenthesizedExpression is called when exiting the ParenthesizedExpression production.
	ExitParenthesizedExpression(c *ParenthesizedExpressionContext)

	// ExitAdditiveExpression is called when exiting the AdditiveExpression production.
	ExitAdditiveExpression(c *AdditiveExpressionContext)

	// ExitRelationalExpression is called when exiting the RelationalExpression production.
	ExitRelationalExpression(c *RelationalExpressionContext)

	// ExitTernaryExpression is called when exiting the TernaryExpression production.
	ExitTernaryExpression(c *TernaryExpressionContext)

	// ExitContainsExpression is called when exiting the ContainsExpression production.
	ExitContainsExpression(c *ContainsExpressionContext)

	// ExitMatchesExpression is called when exiting the MatchesExpression production.
	ExitMatchesExpression(c *MatchesExpressionContext)

	// ExitMapLiteralExpression is called when exiting the MapLiteralExpression production.
	ExitMapLiteralExpression(c *MapLiteralExpressionContext)

	// ExitLiteralExpression is called when exiting the LiteralExpression production.
	ExitLiteralExpression(c *LiteralExpressionContext)

	// ExitInExpression is called when exiting the InExpression production.
	ExitInExpression(c *InExpressionContext)

	// ExitArrayLiteralExpression is called when exiting the ArrayLiteralExpression production.
	ExitArrayLiteralExpression(c *ArrayLiteralExpressionContext)

	// ExitMemberDotExpression is called when exiting the MemberDotExpression production.
	ExitMemberDotExpression(c *MemberDotExpressionContext)

	// ExitUnaryExpression is called when exiting the UnaryExpression production.
	ExitUnaryExpression(c *UnaryExpressionContext)

	// ExitRangeExpression is called when exiting the RangeExpression production.
	ExitRangeExpression(c *RangeExpressionContext)

	// ExitMemberIndexExpression is called when exiting the MemberIndexExpression production.
	ExitMemberIndexExpression(c *MemberIndexExpressionContext)

	// ExitIdentifierExpression is called when exiting the IdentifierExpression production.
	ExitIdentifierExpression(c *IdentifierExpressionContext)

	// ExitPointerExpression is called when exiting the PointerExpression production.
	ExitPointerExpression(c *PointerExpressionContext)

	// ExitLogicalExpression is called when exiting the LogicalExpression production.
	ExitLogicalExpression(c *LogicalExpressionContext)

	// ExitClosureMemberDotExpression is called when exiting the ClosureMemberDotExpression production.
	ExitClosureMemberDotExpression(c *ClosureMemberDotExpressionContext)

	// ExitEndsWithExpression is called when exiting the EndsWithExpression production.
	ExitEndsWithExpression(c *EndsWithExpressionContext)

	// ExitStartsWithExpression is called when exiting the StartsWithExpression production.
	ExitStartsWithExpression(c *StartsWithExpressionContext)

	// ExitEqualityExpression is called when exiting the EqualityExpression production.
	ExitEqualityExpression(c *EqualityExpressionContext)

	// ExitBuiltinLiteralExpression is called when exiting the BuiltinLiteralExpression production.
	ExitBuiltinLiteralExpression(c *BuiltinLiteralExpressionContext)

	// ExitMultiplicativeExpression is called when exiting the MultiplicativeExpression production.
	ExitMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// ExitCallExpression is called when exiting the CallExpression production.
	ExitCallExpression(c *CallExpressionContext)

	// ExitLenBuiltinExpression is called when exiting the LenBuiltinExpression production.
	ExitLenBuiltinExpression(c *LenBuiltinExpressionContext)

	// ExitBuiltinExpression is called when exiting the BuiltinExpression production.
	ExitBuiltinExpression(c *BuiltinExpressionContext)

	// ExitClosureExpression is called when exiting the ClosureExpression production.
	ExitClosureExpression(c *ClosureExpressionContext)

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

	// ExitNilExpression is called when exiting the NilExpression production.
	ExitNilExpression(c *NilExpressionContext)

	// ExitBooleanExpression is called when exiting the BooleanExpression production.
	ExitBooleanExpression(c *BooleanExpressionContext)

	// ExitStringLiteralExpression is called when exiting the StringLiteralExpression production.
	ExitStringLiteralExpression(c *StringLiteralExpressionContext)

	// ExitIntegerExpression is called when exiting the IntegerExpression production.
	ExitIntegerExpression(c *IntegerExpressionContext)

	// ExitFloatExpression is called when exiting the FloatExpression production.
	ExitFloatExpression(c *FloatExpressionContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitIntegerLiteral is called when exiting the integerLiteral production.
	ExitIntegerLiteral(c *IntegerLiteralContext)
}
