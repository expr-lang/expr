// Code generated from Expr.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Expr
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 57, 219,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 47, 10,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5,
	3, 96, 10, 3, 3, 3, 7, 3, 99, 10, 3, 12, 3, 14, 3, 102, 11, 3, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 5, 4, 151, 10,
	4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 7, 6, 160, 10, 6, 12, 6, 14,
	6, 163, 11, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 7, 7, 171, 10, 7, 12,
	7, 14, 7, 174, 11, 7, 3, 7, 5, 7, 177, 10, 7, 3, 7, 3, 7, 5, 7, 181, 10,
	7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 5, 8, 188, 10, 8, 3, 8, 3, 8, 5, 8, 192,
	10, 8, 3, 9, 3, 9, 3, 9, 7, 9, 197, 10, 9, 12, 9, 14, 9, 200, 11, 9, 3,
	10, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12,
	5, 12, 213, 10, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 2, 3, 4, 15, 2,
	4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 2, 10, 3, 2, 23, 25, 3, 2,
	26, 29, 3, 2, 23, 24, 3, 2, 32, 35, 3, 2, 44, 45, 3, 2, 36, 37, 3, 2, 51,
	52, 4, 2, 48, 48, 50, 50, 2, 245, 2, 28, 3, 2, 2, 2, 4, 46, 3, 2, 2, 2,
	6, 150, 3, 2, 2, 2, 8, 152, 3, 2, 2, 2, 10, 156, 3, 2, 2, 2, 12, 180, 3,
	2, 2, 2, 14, 191, 3, 2, 2, 2, 16, 193, 3, 2, 2, 2, 18, 201, 3, 2, 2, 2,
	20, 205, 3, 2, 2, 2, 22, 212, 3, 2, 2, 2, 24, 214, 3, 2, 2, 2, 26, 216,
	3, 2, 2, 2, 28, 29, 5, 4, 3, 2, 29, 30, 7, 2, 2, 3, 30, 3, 3, 2, 2, 2,
	31, 32, 8, 3, 1, 2, 32, 33, 7, 22, 2, 2, 33, 47, 7, 51, 2, 2, 34, 47, 5,
	6, 4, 2, 35, 36, 9, 2, 2, 2, 36, 47, 5, 4, 3, 20, 37, 47, 7, 51, 2, 2,
	38, 47, 7, 38, 2, 2, 39, 47, 5, 22, 12, 2, 40, 47, 5, 12, 7, 2, 41, 47,
	5, 14, 8, 2, 42, 43, 7, 13, 2, 2, 43, 44, 5, 4, 3, 2, 44, 45, 7, 14, 2,
	2, 45, 47, 3, 2, 2, 2, 46, 31, 3, 2, 2, 2, 46, 34, 3, 2, 2, 2, 46, 35,
	3, 2, 2, 2, 46, 37, 3, 2, 2, 2, 46, 38, 3, 2, 2, 2, 46, 39, 3, 2, 2, 2,
	46, 40, 3, 2, 2, 2, 46, 41, 3, 2, 2, 2, 46, 42, 3, 2, 2, 2, 47, 100, 3,
	2, 2, 2, 48, 49, 12, 19, 2, 2, 49, 50, 7, 3, 2, 2, 50, 99, 5, 4, 3, 20,
	51, 52, 12, 18, 2, 2, 52, 53, 9, 3, 2, 2, 53, 99, 5, 4, 3, 19, 54, 55,
	12, 17, 2, 2, 55, 56, 9, 4, 2, 2, 56, 99, 5, 4, 3, 18, 57, 58, 12, 16,
	2, 2, 58, 59, 9, 5, 2, 2, 59, 99, 5, 4, 3, 17, 60, 61, 12, 15, 2, 2, 61,
	62, 7, 42, 2, 2, 62, 99, 5, 4, 3, 16, 63, 64, 12, 14, 2, 2, 64, 65, 7,
	43, 2, 2, 65, 99, 5, 4, 3, 15, 66, 67, 12, 13, 2, 2, 67, 68, 9, 6, 2, 2,
	68, 99, 5, 4, 3, 14, 69, 70, 12, 12, 2, 2, 70, 71, 9, 7, 2, 2, 71, 99,
	5, 4, 3, 13, 72, 73, 12, 11, 2, 2, 73, 74, 7, 39, 2, 2, 74, 99, 5, 4, 3,
	12, 75, 76, 12, 10, 2, 2, 76, 77, 7, 40, 2, 2, 77, 99, 5, 4, 3, 11, 78,
	79, 12, 9, 2, 2, 79, 80, 7, 20, 2, 2, 80, 81, 5, 4, 3, 2, 81, 82, 7, 21,
	2, 2, 82, 83, 5, 4, 3, 10, 83, 99, 3, 2, 2, 2, 84, 85, 12, 24, 2, 2, 85,
	86, 7, 11, 2, 2, 86, 87, 5, 4, 3, 2, 87, 88, 7, 12, 2, 2, 88, 99, 3, 2,
	2, 2, 89, 90, 12, 23, 2, 2, 90, 91, 7, 22, 2, 2, 91, 99, 7, 51, 2, 2, 92,
	93, 12, 21, 2, 2, 93, 95, 7, 13, 2, 2, 94, 96, 5, 10, 6, 2, 95, 94, 3,
	2, 2, 2, 95, 96, 3, 2, 2, 2, 96, 97, 3, 2, 2, 2, 97, 99, 7, 14, 2, 2, 98,
	48, 3, 2, 2, 2, 98, 51, 3, 2, 2, 2, 98, 54, 3, 2, 2, 2, 98, 57, 3, 2, 2,
	2, 98, 60, 3, 2, 2, 2, 98, 63, 3, 2, 2, 2, 98, 66, 3, 2, 2, 2, 98, 69,
	3, 2, 2, 2, 98, 72, 3, 2, 2, 2, 98, 75, 3, 2, 2, 2, 98, 78, 3, 2, 2, 2,
	98, 84, 3, 2, 2, 2, 98, 89, 3, 2, 2, 2, 98, 92, 3, 2, 2, 2, 99, 102, 3,
	2, 2, 2, 100, 98, 3, 2, 2, 2, 100, 101, 3, 2, 2, 2, 101, 5, 3, 2, 2, 2,
	102, 100, 3, 2, 2, 2, 103, 104, 7, 4, 2, 2, 104, 105, 7, 13, 2, 2, 105,
	106, 5, 4, 3, 2, 106, 107, 7, 14, 2, 2, 107, 151, 3, 2, 2, 2, 108, 109,
	7, 5, 2, 2, 109, 110, 7, 13, 2, 2, 110, 111, 5, 4, 3, 2, 111, 112, 7, 18,
	2, 2, 112, 113, 5, 8, 5, 2, 113, 114, 7, 14, 2, 2, 114, 151, 3, 2, 2, 2,
	115, 116, 7, 6, 2, 2, 116, 117, 7, 13, 2, 2, 117, 118, 5, 4, 3, 2, 118,
	119, 7, 18, 2, 2, 119, 120, 5, 8, 5, 2, 120, 121, 7, 14, 2, 2, 121, 151,
	3, 2, 2, 2, 122, 123, 7, 7, 2, 2, 123, 124, 7, 13, 2, 2, 124, 125, 5, 4,
	3, 2, 125, 126, 7, 18, 2, 2, 126, 127, 5, 8, 5, 2, 127, 128, 7, 14, 2,
	2, 128, 151, 3, 2, 2, 2, 129, 130, 7, 8, 2, 2, 130, 131, 7, 13, 2, 2, 131,
	132, 5, 4, 3, 2, 132, 133, 7, 18, 2, 2, 133, 134, 5, 8, 5, 2, 134, 135,
	7, 14, 2, 2, 135, 151, 3, 2, 2, 2, 136, 137, 7, 9, 2, 2, 137, 138, 7, 13,
	2, 2, 138, 139, 5, 4, 3, 2, 139, 140, 7, 18, 2, 2, 140, 141, 5, 8, 5, 2,
	141, 142, 7, 14, 2, 2, 142, 151, 3, 2, 2, 2, 143, 144, 7, 10, 2, 2, 144,
	145, 7, 13, 2, 2, 145, 146, 5, 4, 3, 2, 146, 147, 7, 18, 2, 2, 147, 148,
	5, 8, 5, 2, 148, 149, 7, 14, 2, 2, 149, 151, 3, 2, 2, 2, 150, 103, 3, 2,
	2, 2, 150, 108, 3, 2, 2, 2, 150, 115, 3, 2, 2, 2, 150, 122, 3, 2, 2, 2,
	150, 129, 3, 2, 2, 2, 150, 136, 3, 2, 2, 2, 150, 143, 3, 2, 2, 2, 151,
	7, 3, 2, 2, 2, 152, 153, 7, 15, 2, 2, 153, 154, 5, 4, 3, 2, 154, 155, 7,
	16, 2, 2, 155, 9, 3, 2, 2, 2, 156, 161, 5, 4, 3, 2, 157, 158, 7, 18, 2,
	2, 158, 160, 5, 4, 3, 2, 159, 157, 3, 2, 2, 2, 160, 163, 3, 2, 2, 2, 161,
	159, 3, 2, 2, 2, 161, 162, 3, 2, 2, 2, 162, 11, 3, 2, 2, 2, 163, 161, 3,
	2, 2, 2, 164, 165, 7, 11, 2, 2, 165, 181, 7, 12, 2, 2, 166, 167, 7, 11,
	2, 2, 167, 172, 5, 4, 3, 2, 168, 169, 7, 18, 2, 2, 169, 171, 5, 4, 3, 2,
	170, 168, 3, 2, 2, 2, 171, 174, 3, 2, 2, 2, 172, 170, 3, 2, 2, 2, 172,
	173, 3, 2, 2, 2, 173, 176, 3, 2, 2, 2, 174, 172, 3, 2, 2, 2, 175, 177,
	7, 18, 2, 2, 176, 175, 3, 2, 2, 2, 176, 177, 3, 2, 2, 2, 177, 178, 3, 2,
	2, 2, 178, 179, 7, 12, 2, 2, 179, 181, 3, 2, 2, 2, 180, 164, 3, 2, 2, 2,
	180, 166, 3, 2, 2, 2, 181, 13, 3, 2, 2, 2, 182, 183, 7, 15, 2, 2, 183,
	192, 7, 16, 2, 2, 184, 185, 7, 15, 2, 2, 185, 187, 5, 16, 9, 2, 186, 188,
	7, 18, 2, 2, 187, 186, 3, 2, 2, 2, 187, 188, 3, 2, 2, 2, 188, 189, 3, 2,
	2, 2, 189, 190, 7, 16, 2, 2, 190, 192, 3, 2, 2, 2, 191, 182, 3, 2, 2, 2,
	191, 184, 3, 2, 2, 2, 192, 15, 3, 2, 2, 2, 193, 198, 5, 18, 10, 2, 194,
	195, 7, 18, 2, 2, 195, 197, 5, 18, 10, 2, 196, 194, 3, 2, 2, 2, 197, 200,
	3, 2, 2, 2, 198, 196, 3, 2, 2, 2, 198, 199, 3, 2, 2, 2, 199, 17, 3, 2,
	2, 2, 200, 198, 3, 2, 2, 2, 201, 202, 5, 20, 11, 2, 202, 203, 7, 21, 2,
	2, 203, 204, 5, 4, 3, 2, 204, 19, 3, 2, 2, 2, 205, 206, 9, 8, 2, 2, 206,
	21, 3, 2, 2, 2, 207, 213, 7, 46, 2, 2, 208, 213, 7, 47, 2, 2, 209, 213,
	5, 24, 13, 2, 210, 213, 5, 26, 14, 2, 211, 213, 7, 49, 2, 2, 212, 207,
	3, 2, 2, 2, 212, 208, 3, 2, 2, 2, 212, 209, 3, 2, 2, 2, 212, 210, 3, 2,
	2, 2, 212, 211, 3, 2, 2, 2, 213, 23, 3, 2, 2, 2, 214, 215, 7, 52, 2, 2,
	215, 25, 3, 2, 2, 2, 216, 217, 9, 9, 2, 2, 217, 27, 3, 2, 2, 2, 15, 46,
	95, 98, 100, 150, 161, 172, 176, 180, 187, 191, 198, 212,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'..'", "'len'", "'all'", "'none'", "'any'", "'one'", "'filter'", "'map'",
	"'['", "']'", "'('", "')'", "'{'", "'}'", "';'", "','", "'='", "'?'", "':'",
	"'.'", "'+'", "'-'", "", "'*'", "'**'", "'/'", "'%'", "'>>'", "'<<'", "'<'",
	"'>'", "'<='", "'>='", "'=='", "'!='", "'#'", "", "", "", "'contains'",
	"'matches'", "'in'", "'not in'", "'nil'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "SemiColon", "Comma", "Assign",
	"QuestionMark", "Colon", "Dot", "Plus", "Minus", "Not", "Multiply", "Exponent",
	"Divide", "Modulus", "RightShiftArithmetic", "LeftShiftArithmetic", "LessThan",
	"MoreThan", "LessThanEquals", "GreaterThanEquals", "Equals", "NotEquals",
	"Pointer", "And", "Or", "Builtins", "Contains", "Matches", "In", "NotIn",
	"NilLiteral", "BooleanLiteral", "IntegerLiteral", "FloatLiteral", "HexIntegerLiteral",
	"Identifier", "StringLiteral", "WhiteSpaces", "MultiLineComment", "SingleLineComment",
	"LineTerminator", "UnexpectedCharacter",
}

var ruleNames = []string{
	"start", "expr", "builtins", "closure", "arguments", "arrayLiteral", "mapLiteral",
	"propertyNameAndValueList", "propertyAssignment", "propertyName", "literal",
	"stringLiteral", "integerLiteral",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type ExprParser struct {
	*antlr.BaseParser
}

func NewExprParser(input antlr.TokenStream) *ExprParser {
	this := new(ExprParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Expr.g4"

	return this
}

// ExprParser tokens.
const (
	ExprParserEOF                  = antlr.TokenEOF
	ExprParserT__0                 = 1
	ExprParserT__1                 = 2
	ExprParserT__2                 = 3
	ExprParserT__3                 = 4
	ExprParserT__4                 = 5
	ExprParserT__5                 = 6
	ExprParserT__6                 = 7
	ExprParserT__7                 = 8
	ExprParserOpenBracket          = 9
	ExprParserCloseBracket         = 10
	ExprParserOpenParen            = 11
	ExprParserCloseParen           = 12
	ExprParserOpenBrace            = 13
	ExprParserCloseBrace           = 14
	ExprParserSemiColon            = 15
	ExprParserComma                = 16
	ExprParserAssign               = 17
	ExprParserQuestionMark         = 18
	ExprParserColon                = 19
	ExprParserDot                  = 20
	ExprParserPlus                 = 21
	ExprParserMinus                = 22
	ExprParserNot                  = 23
	ExprParserMultiply             = 24
	ExprParserExponent             = 25
	ExprParserDivide               = 26
	ExprParserModulus              = 27
	ExprParserRightShiftArithmetic = 28
	ExprParserLeftShiftArithmetic  = 29
	ExprParserLessThan             = 30
	ExprParserMoreThan             = 31
	ExprParserLessThanEquals       = 32
	ExprParserGreaterThanEquals    = 33
	ExprParserEquals               = 34
	ExprParserNotEquals            = 35
	ExprParserPointer              = 36
	ExprParserAnd                  = 37
	ExprParserOr                   = 38
	ExprParserBuiltins             = 39
	ExprParserContains             = 40
	ExprParserMatches              = 41
	ExprParserIn                   = 42
	ExprParserNotIn                = 43
	ExprParserNilLiteral           = 44
	ExprParserBooleanLiteral       = 45
	ExprParserIntegerLiteral       = 46
	ExprParserFloatLiteral         = 47
	ExprParserHexIntegerLiteral    = 48
	ExprParserIdentifier           = 49
	ExprParserStringLiteral        = 50
	ExprParserWhiteSpaces          = 51
	ExprParserMultiLineComment     = 52
	ExprParserSingleLineComment    = 53
	ExprParserLineTerminator       = 54
	ExprParserUnexpectedCharacter  = 55
)

// ExprParser rules.
const (
	ExprParserRULE_start                    = 0
	ExprParserRULE_expr                     = 1
	ExprParserRULE_builtins                 = 2
	ExprParserRULE_closure                  = 3
	ExprParserRULE_arguments                = 4
	ExprParserRULE_arrayLiteral             = 5
	ExprParserRULE_mapLiteral               = 6
	ExprParserRULE_propertyNameAndValueList = 7
	ExprParserRULE_propertyAssignment       = 8
	ExprParserRULE_propertyName             = 9
	ExprParserRULE_literal                  = 10
	ExprParserRULE_stringLiteral            = 11
	ExprParserRULE_integerLiteral           = 12
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IExprContext
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) GetE() IExprContext { return s.e }

func (s *StartContext) SetE(v IExprContext) { s.e = v }

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(ExprParserEOF, 0)
}

func (s *StartContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitStart(s)
	}
}

func (s *StartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitStart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) Start() (localctx IStartContext) {
	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ExprParserRULE_start)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(26)

		var _x = p.expr(0)

		localctx.(*StartContext).e = _x
	}
	{
		p.SetState(27)
		p.Match(ExprParserEOF)
	}

	return localctx
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) CopyFrom(ctx *ExprContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ParenthesizedExpressionContext struct {
	*ExprContext
}

func NewParenthesizedExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenthesizedExpressionContext {
	var p = new(ParenthesizedExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *ParenthesizedExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenthesizedExpressionContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *ParenthesizedExpressionContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ParenthesizedExpressionContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *ParenthesizedExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterParenthesizedExpression(s)
	}
}

func (s *ParenthesizedExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitParenthesizedExpression(s)
	}
}

func (s *ParenthesizedExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitParenthesizedExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type AdditiveExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewAdditiveExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AdditiveExpressionContext {
	var p = new(AdditiveExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *AdditiveExpressionContext) GetOp() antlr.Token { return s.op }

func (s *AdditiveExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *AdditiveExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditiveExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *AdditiveExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *AdditiveExpressionContext) Plus() antlr.TerminalNode {
	return s.GetToken(ExprParserPlus, 0)
}

func (s *AdditiveExpressionContext) Minus() antlr.TerminalNode {
	return s.GetToken(ExprParserMinus, 0)
}

func (s *AdditiveExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterAdditiveExpression(s)
	}
}

func (s *AdditiveExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitAdditiveExpression(s)
	}
}

func (s *AdditiveExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitAdditiveExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type RelationalExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewRelationalExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RelationalExpressionContext {
	var p = new(RelationalExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *RelationalExpressionContext) GetOp() antlr.Token { return s.op }

func (s *RelationalExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *RelationalExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationalExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *RelationalExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *RelationalExpressionContext) LessThan() antlr.TerminalNode {
	return s.GetToken(ExprParserLessThan, 0)
}

func (s *RelationalExpressionContext) MoreThan() antlr.TerminalNode {
	return s.GetToken(ExprParserMoreThan, 0)
}

func (s *RelationalExpressionContext) LessThanEquals() antlr.TerminalNode {
	return s.GetToken(ExprParserLessThanEquals, 0)
}

func (s *RelationalExpressionContext) GreaterThanEquals() antlr.TerminalNode {
	return s.GetToken(ExprParserGreaterThanEquals, 0)
}

func (s *RelationalExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterRelationalExpression(s)
	}
}

func (s *RelationalExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitRelationalExpression(s)
	}
}

func (s *RelationalExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitRelationalExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type TernaryExpressionContext struct {
	*ExprContext
	e1 IExprContext
	e2 IExprContext
}

func NewTernaryExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TernaryExpressionContext {
	var p = new(TernaryExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *TernaryExpressionContext) GetE1() IExprContext { return s.e1 }

func (s *TernaryExpressionContext) GetE2() IExprContext { return s.e2 }

func (s *TernaryExpressionContext) SetE1(v IExprContext) { s.e1 = v }

func (s *TernaryExpressionContext) SetE2(v IExprContext) { s.e2 = v }

func (s *TernaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TernaryExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *TernaryExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *TernaryExpressionContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(ExprParserQuestionMark, 0)
}

func (s *TernaryExpressionContext) Colon() antlr.TerminalNode {
	return s.GetToken(ExprParserColon, 0)
}

func (s *TernaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterTernaryExpression(s)
	}
}

func (s *TernaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitTernaryExpression(s)
	}
}

func (s *TernaryExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitTernaryExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type ContainsExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewContainsExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ContainsExpressionContext {
	var p = new(ContainsExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *ContainsExpressionContext) GetOp() antlr.Token { return s.op }

func (s *ContainsExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *ContainsExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContainsExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *ContainsExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ContainsExpressionContext) Contains() antlr.TerminalNode {
	return s.GetToken(ExprParserContains, 0)
}

func (s *ContainsExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterContainsExpression(s)
	}
}

func (s *ContainsExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitContainsExpression(s)
	}
}

func (s *ContainsExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitContainsExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MatchesExpressionContext struct {
	*ExprContext
	op      antlr.Token
	pattern IExprContext
}

func NewMatchesExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MatchesExpressionContext {
	var p = new(MatchesExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MatchesExpressionContext) GetOp() antlr.Token { return s.op }

func (s *MatchesExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *MatchesExpressionContext) GetPattern() IExprContext { return s.pattern }

func (s *MatchesExpressionContext) SetPattern(v IExprContext) { s.pattern = v }

func (s *MatchesExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MatchesExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *MatchesExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MatchesExpressionContext) Matches() antlr.TerminalNode {
	return s.GetToken(ExprParserMatches, 0)
}

func (s *MatchesExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMatchesExpression(s)
	}
}

func (s *MatchesExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMatchesExpression(s)
	}
}

func (s *MatchesExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMatchesExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MapLiteralExpressionContext struct {
	*ExprContext
}

func NewMapLiteralExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MapLiteralExpressionContext {
	var p = new(MapLiteralExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MapLiteralExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapLiteralExpressionContext) MapLiteral() IMapLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMapLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMapLiteralContext)
}

func (s *MapLiteralExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMapLiteralExpression(s)
	}
}

func (s *MapLiteralExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMapLiteralExpression(s)
	}
}

func (s *MapLiteralExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMapLiteralExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type LiteralExpressionContext struct {
	*ExprContext
}

func NewLiteralExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LiteralExpressionContext {
	var p = new(LiteralExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *LiteralExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralExpressionContext) Literal() ILiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *LiteralExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterLiteralExpression(s)
	}
}

func (s *LiteralExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitLiteralExpression(s)
	}
}

func (s *LiteralExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitLiteralExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type InExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewInExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InExpressionContext {
	var p = new(InExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *InExpressionContext) GetOp() antlr.Token { return s.op }

func (s *InExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *InExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *InExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *InExpressionContext) In() antlr.TerminalNode {
	return s.GetToken(ExprParserIn, 0)
}

func (s *InExpressionContext) NotIn() antlr.TerminalNode {
	return s.GetToken(ExprParserNotIn, 0)
}

func (s *InExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterInExpression(s)
	}
}

func (s *InExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitInExpression(s)
	}
}

func (s *InExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitInExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type ArrayLiteralExpressionContext struct {
	*ExprContext
}

func NewArrayLiteralExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ArrayLiteralExpressionContext {
	var p = new(ArrayLiteralExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *ArrayLiteralExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayLiteralExpressionContext) ArrayLiteral() IArrayLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *ArrayLiteralExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterArrayLiteralExpression(s)
	}
}

func (s *ArrayLiteralExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitArrayLiteralExpression(s)
	}
}

func (s *ArrayLiteralExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitArrayLiteralExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MemberDotExpressionContext struct {
	*ExprContext
	name antlr.Token
}

func NewMemberDotExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberDotExpressionContext {
	var p = new(MemberDotExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MemberDotExpressionContext) GetName() antlr.Token { return s.name }

func (s *MemberDotExpressionContext) SetName(v antlr.Token) { s.name = v }

func (s *MemberDotExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberDotExpressionContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MemberDotExpressionContext) Dot() antlr.TerminalNode {
	return s.GetToken(ExprParserDot, 0)
}

func (s *MemberDotExpressionContext) Identifier() antlr.TerminalNode {
	return s.GetToken(ExprParserIdentifier, 0)
}

func (s *MemberDotExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMemberDotExpression(s)
	}
}

func (s *MemberDotExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMemberDotExpression(s)
	}
}

func (s *MemberDotExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMemberDotExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type UnaryExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewUnaryExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UnaryExpressionContext {
	var p = new(UnaryExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *UnaryExpressionContext) GetOp() antlr.Token { return s.op }

func (s *UnaryExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *UnaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryExpressionContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *UnaryExpressionContext) Plus() antlr.TerminalNode {
	return s.GetToken(ExprParserPlus, 0)
}

func (s *UnaryExpressionContext) Minus() antlr.TerminalNode {
	return s.GetToken(ExprParserMinus, 0)
}

func (s *UnaryExpressionContext) Not() antlr.TerminalNode {
	return s.GetToken(ExprParserNot, 0)
}

func (s *UnaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterUnaryExpression(s)
	}
}

func (s *UnaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitUnaryExpression(s)
	}
}

func (s *UnaryExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitUnaryExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type RangeExpressionContext struct {
	*ExprContext
}

func NewRangeExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RangeExpressionContext {
	var p = new(RangeExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *RangeExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RangeExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *RangeExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *RangeExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterRangeExpression(s)
	}
}

func (s *RangeExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitRangeExpression(s)
	}
}

func (s *RangeExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitRangeExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MemberIndexExpressionContext struct {
	*ExprContext
	index IExprContext
}

func NewMemberIndexExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberIndexExpressionContext {
	var p = new(MemberIndexExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MemberIndexExpressionContext) GetIndex() IExprContext { return s.index }

func (s *MemberIndexExpressionContext) SetIndex(v IExprContext) { s.index = v }

func (s *MemberIndexExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberIndexExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *MemberIndexExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MemberIndexExpressionContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenBracket, 0)
}

func (s *MemberIndexExpressionContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseBracket, 0)
}

func (s *MemberIndexExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMemberIndexExpression(s)
	}
}

func (s *MemberIndexExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMemberIndexExpression(s)
	}
}

func (s *MemberIndexExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMemberIndexExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdentifierExpressionContext struct {
	*ExprContext
}

func NewIdentifierExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentifierExpressionContext {
	var p = new(IdentifierExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *IdentifierExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierExpressionContext) Identifier() antlr.TerminalNode {
	return s.GetToken(ExprParserIdentifier, 0)
}

func (s *IdentifierExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterIdentifierExpression(s)
	}
}

func (s *IdentifierExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitIdentifierExpression(s)
	}
}

func (s *IdentifierExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitIdentifierExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type PointerExpressionContext struct {
	*ExprContext
}

func NewPointerExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PointerExpressionContext {
	var p = new(PointerExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *PointerExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PointerExpressionContext) Pointer() antlr.TerminalNode {
	return s.GetToken(ExprParserPointer, 0)
}

func (s *PointerExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterPointerExpression(s)
	}
}

func (s *PointerExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitPointerExpression(s)
	}
}

func (s *PointerExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitPointerExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type LogicalExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewLogicalExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalExpressionContext {
	var p = new(LogicalExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *LogicalExpressionContext) GetOp() antlr.Token { return s.op }

func (s *LogicalExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *LogicalExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *LogicalExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *LogicalExpressionContext) And() antlr.TerminalNode {
	return s.GetToken(ExprParserAnd, 0)
}

func (s *LogicalExpressionContext) Or() antlr.TerminalNode {
	return s.GetToken(ExprParserOr, 0)
}

func (s *LogicalExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterLogicalExpression(s)
	}
}

func (s *LogicalExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitLogicalExpression(s)
	}
}

func (s *LogicalExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitLogicalExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type ClosureMemberDotExpressionContext struct {
	*ExprContext
	name antlr.Token
}

func NewClosureMemberDotExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ClosureMemberDotExpressionContext {
	var p = new(ClosureMemberDotExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *ClosureMemberDotExpressionContext) GetName() antlr.Token { return s.name }

func (s *ClosureMemberDotExpressionContext) SetName(v antlr.Token) { s.name = v }

func (s *ClosureMemberDotExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClosureMemberDotExpressionContext) Dot() antlr.TerminalNode {
	return s.GetToken(ExprParserDot, 0)
}

func (s *ClosureMemberDotExpressionContext) Identifier() antlr.TerminalNode {
	return s.GetToken(ExprParserIdentifier, 0)
}

func (s *ClosureMemberDotExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterClosureMemberDotExpression(s)
	}
}

func (s *ClosureMemberDotExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitClosureMemberDotExpression(s)
	}
}

func (s *ClosureMemberDotExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitClosureMemberDotExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type EqualityExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewEqualityExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EqualityExpressionContext {
	var p = new(EqualityExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *EqualityExpressionContext) GetOp() antlr.Token { return s.op }

func (s *EqualityExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *EqualityExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *EqualityExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *EqualityExpressionContext) Equals() antlr.TerminalNode {
	return s.GetToken(ExprParserEquals, 0)
}

func (s *EqualityExpressionContext) NotEquals() antlr.TerminalNode {
	return s.GetToken(ExprParserNotEquals, 0)
}

func (s *EqualityExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterEqualityExpression(s)
	}
}

func (s *EqualityExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitEqualityExpression(s)
	}
}

func (s *EqualityExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitEqualityExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type BuiltinLiteralExpressionContext struct {
	*ExprContext
}

func NewBuiltinLiteralExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BuiltinLiteralExpressionContext {
	var p = new(BuiltinLiteralExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *BuiltinLiteralExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BuiltinLiteralExpressionContext) Builtins() IBuiltinsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBuiltinsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBuiltinsContext)
}

func (s *BuiltinLiteralExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBuiltinLiteralExpression(s)
	}
}

func (s *BuiltinLiteralExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBuiltinLiteralExpression(s)
	}
}

func (s *BuiltinLiteralExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBuiltinLiteralExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MultiplicativeExpressionContext struct {
	*ExprContext
	op antlr.Token
}

func NewMultiplicativeExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MultiplicativeExpressionContext {
	var p = new(MultiplicativeExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MultiplicativeExpressionContext) GetOp() antlr.Token { return s.op }

func (s *MultiplicativeExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *MultiplicativeExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicativeExpressionContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *MultiplicativeExpressionContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MultiplicativeExpressionContext) Multiply() antlr.TerminalNode {
	return s.GetToken(ExprParserMultiply, 0)
}

func (s *MultiplicativeExpressionContext) Exponent() antlr.TerminalNode {
	return s.GetToken(ExprParserExponent, 0)
}

func (s *MultiplicativeExpressionContext) Divide() antlr.TerminalNode {
	return s.GetToken(ExprParserDivide, 0)
}

func (s *MultiplicativeExpressionContext) Modulus() antlr.TerminalNode {
	return s.GetToken(ExprParserModulus, 0)
}

func (s *MultiplicativeExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMultiplicativeExpression(s)
	}
}

func (s *MultiplicativeExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMultiplicativeExpression(s)
	}
}

func (s *MultiplicativeExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMultiplicativeExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type CallExpressionContext struct {
	*ExprContext
	args IArgumentsContext
}

func NewCallExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CallExpressionContext {
	var p = new(CallExpressionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *CallExpressionContext) GetArgs() IArgumentsContext { return s.args }

func (s *CallExpressionContext) SetArgs(v IArgumentsContext) { s.args = v }

func (s *CallExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CallExpressionContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *CallExpressionContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *CallExpressionContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *CallExpressionContext) Arguments() IArgumentsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgumentsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgumentsContext)
}

func (s *CallExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterCallExpression(s)
	}
}

func (s *CallExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitCallExpression(s)
	}
}

func (s *CallExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitCallExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *ExprParser) expr(_p int) (localctx IExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, ExprParserRULE_expr, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(44)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ExprParserDot:
		localctx = NewClosureMemberDotExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(30)
			p.Match(ExprParserDot)
		}
		{
			p.SetState(31)

			var _m = p.Match(ExprParserIdentifier)

			localctx.(*ClosureMemberDotExpressionContext).name = _m
		}

	case ExprParserT__1, ExprParserT__2, ExprParserT__3, ExprParserT__4, ExprParserT__5, ExprParserT__6, ExprParserT__7:
		localctx = NewBuiltinLiteralExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(32)
			p.Builtins()
		}

	case ExprParserPlus, ExprParserMinus, ExprParserNot:
		localctx = NewUnaryExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(33)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*UnaryExpressionContext).op = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExprParserPlus)|(1<<ExprParserMinus)|(1<<ExprParserNot))) != 0) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*UnaryExpressionContext).op = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(34)
			p.expr(18)
		}

	case ExprParserIdentifier:
		localctx = NewIdentifierExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(35)
			p.Match(ExprParserIdentifier)
		}

	case ExprParserPointer:
		localctx = NewPointerExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(36)
			p.Match(ExprParserPointer)
		}

	case ExprParserNilLiteral, ExprParserBooleanLiteral, ExprParserIntegerLiteral, ExprParserFloatLiteral, ExprParserHexIntegerLiteral, ExprParserStringLiteral:
		localctx = NewLiteralExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(37)
			p.Literal()
		}

	case ExprParserOpenBracket:
		localctx = NewArrayLiteralExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(38)
			p.ArrayLiteral()
		}

	case ExprParserOpenBrace:
		localctx = NewMapLiteralExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(39)
			p.MapLiteral()
		}

	case ExprParserOpenParen:
		localctx = NewParenthesizedExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(40)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(41)
			p.expr(0)
		}
		{
			p.SetState(42)
			p.Match(ExprParserCloseParen)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(98)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(96)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
			case 1:
				localctx = NewRangeExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(46)

				if !(p.Precpred(p.GetParserRuleContext(), 17)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 17)", ""))
				}
				{
					p.SetState(47)
					p.Match(ExprParserT__0)
				}
				{
					p.SetState(48)
					p.expr(18)
				}

			case 2:
				localctx = NewMultiplicativeExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(49)

				if !(p.Precpred(p.GetParserRuleContext(), 16)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 16)", ""))
				}
				{
					p.SetState(50)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*MultiplicativeExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExprParserMultiply)|(1<<ExprParserExponent)|(1<<ExprParserDivide)|(1<<ExprParserModulus))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*MultiplicativeExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(51)
					p.expr(17)
				}

			case 3:
				localctx = NewAdditiveExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(52)

				if !(p.Precpred(p.GetParserRuleContext(), 15)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 15)", ""))
				}
				{
					p.SetState(53)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*AdditiveExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserPlus || _la == ExprParserMinus) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*AdditiveExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(54)
					p.expr(16)
				}

			case 4:
				localctx = NewRelationalExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(55)

				if !(p.Precpred(p.GetParserRuleContext(), 14)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 14)", ""))
				}
				{
					p.SetState(56)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*RelationalExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la-30)&-(0x1f+1)) == 0 && ((1<<uint((_la-30)))&((1<<(ExprParserLessThan-30))|(1<<(ExprParserMoreThan-30))|(1<<(ExprParserLessThanEquals-30))|(1<<(ExprParserGreaterThanEquals-30)))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*RelationalExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(57)
					p.expr(15)
				}

			case 5:
				localctx = NewContainsExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(58)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
				}
				{
					p.SetState(59)

					var _m = p.Match(ExprParserContains)

					localctx.(*ContainsExpressionContext).op = _m
				}
				{
					p.SetState(60)
					p.expr(14)
				}

			case 6:
				localctx = NewMatchesExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(61)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(62)

					var _m = p.Match(ExprParserMatches)

					localctx.(*MatchesExpressionContext).op = _m
				}
				{
					p.SetState(63)

					var _x = p.expr(13)

					localctx.(*MatchesExpressionContext).pattern = _x
				}

			case 7:
				localctx = NewInExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(64)

				if !(p.Precpred(p.GetParserRuleContext(), 11)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 11)", ""))
				}
				{
					p.SetState(65)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*InExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserIn || _la == ExprParserNotIn) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*InExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(66)
					p.expr(12)
				}

			case 8:
				localctx = NewEqualityExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(67)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
				}
				{
					p.SetState(68)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*EqualityExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserEquals || _la == ExprParserNotEquals) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*EqualityExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(69)
					p.expr(11)
				}

			case 9:
				localctx = NewLogicalExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(70)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
				}
				{
					p.SetState(71)

					var _m = p.Match(ExprParserAnd)

					localctx.(*LogicalExpressionContext).op = _m
				}
				{
					p.SetState(72)
					p.expr(10)
				}

			case 10:
				localctx = NewLogicalExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(73)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(74)

					var _m = p.Match(ExprParserOr)

					localctx.(*LogicalExpressionContext).op = _m
				}
				{
					p.SetState(75)
					p.expr(9)
				}

			case 11:
				localctx = NewTernaryExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(76)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(77)
					p.Match(ExprParserQuestionMark)
				}
				{
					p.SetState(78)

					var _x = p.expr(0)

					localctx.(*TernaryExpressionContext).e1 = _x
				}
				{
					p.SetState(79)
					p.Match(ExprParserColon)
				}
				{
					p.SetState(80)

					var _x = p.expr(8)

					localctx.(*TernaryExpressionContext).e2 = _x
				}

			case 12:
				localctx = NewMemberIndexExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(82)

				if !(p.Precpred(p.GetParserRuleContext(), 22)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 22)", ""))
				}
				{
					p.SetState(83)
					p.Match(ExprParserOpenBracket)
				}
				{
					p.SetState(84)

					var _x = p.expr(0)

					localctx.(*MemberIndexExpressionContext).index = _x
				}
				{
					p.SetState(85)
					p.Match(ExprParserCloseBracket)
				}

			case 13:
				localctx = NewMemberDotExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(87)

				if !(p.Precpred(p.GetParserRuleContext(), 21)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 21)", ""))
				}
				{
					p.SetState(88)
					p.Match(ExprParserDot)
				}
				{
					p.SetState(89)

					var _m = p.Match(ExprParserIdentifier)

					localctx.(*MemberDotExpressionContext).name = _m
				}

			case 14:
				localctx = NewCallExpressionContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(90)

				if !(p.Precpred(p.GetParserRuleContext(), 19)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 19)", ""))
				}
				{
					p.SetState(91)
					p.Match(ExprParserOpenParen)
				}
				p.SetState(93)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExprParserT__1)|(1<<ExprParserT__2)|(1<<ExprParserT__3)|(1<<ExprParserT__4)|(1<<ExprParserT__5)|(1<<ExprParserT__6)|(1<<ExprParserT__7)|(1<<ExprParserOpenBracket)|(1<<ExprParserOpenParen)|(1<<ExprParserOpenBrace)|(1<<ExprParserDot)|(1<<ExprParserPlus)|(1<<ExprParserMinus)|(1<<ExprParserNot))) != 0) || (((_la-36)&-(0x1f+1)) == 0 && ((1<<uint((_la-36)))&((1<<(ExprParserPointer-36))|(1<<(ExprParserNilLiteral-36))|(1<<(ExprParserBooleanLiteral-36))|(1<<(ExprParserIntegerLiteral-36))|(1<<(ExprParserFloatLiteral-36))|(1<<(ExprParserHexIntegerLiteral-36))|(1<<(ExprParserIdentifier-36))|(1<<(ExprParserStringLiteral-36)))) != 0) {
					{
						p.SetState(92)

						var _x = p.Arguments()

						localctx.(*CallExpressionContext).args = _x
					}

				}
				{
					p.SetState(95)
					p.Match(ExprParserCloseParen)
				}

			}

		}
		p.SetState(100)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}

	return localctx
}

// IBuiltinsContext is an interface to support dynamic dispatch.
type IBuiltinsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBuiltinsContext differentiates from other interfaces.
	IsBuiltinsContext()
}

type BuiltinsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBuiltinsContext() *BuiltinsContext {
	var p = new(BuiltinsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_builtins
	return p
}

func (*BuiltinsContext) IsBuiltinsContext() {}

func NewBuiltinsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BuiltinsContext {
	var p = new(BuiltinsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_builtins

	return p
}

func (s *BuiltinsContext) GetParser() antlr.Parser { return s.parser }

func (s *BuiltinsContext) CopyFrom(ctx *BuiltinsContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *BuiltinsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BuiltinsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LenBuiltinExpressionContext struct {
	*BuiltinsContext
	e IExprContext
}

func NewLenBuiltinExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LenBuiltinExpressionContext {
	var p = new(LenBuiltinExpressionContext)

	p.BuiltinsContext = NewEmptyBuiltinsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*BuiltinsContext))

	return p
}

func (s *LenBuiltinExpressionContext) GetE() IExprContext { return s.e }

func (s *LenBuiltinExpressionContext) SetE(v IExprContext) { s.e = v }

func (s *LenBuiltinExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LenBuiltinExpressionContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *LenBuiltinExpressionContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *LenBuiltinExpressionContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *LenBuiltinExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterLenBuiltinExpression(s)
	}
}

func (s *LenBuiltinExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitLenBuiltinExpression(s)
	}
}

func (s *LenBuiltinExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitLenBuiltinExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type BuiltinExpressionContext struct {
	*BuiltinsContext
	name antlr.Token
	e    IExprContext
	c    IClosureContext
}

func NewBuiltinExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BuiltinExpressionContext {
	var p = new(BuiltinExpressionContext)

	p.BuiltinsContext = NewEmptyBuiltinsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*BuiltinsContext))

	return p
}

func (s *BuiltinExpressionContext) GetName() antlr.Token { return s.name }

func (s *BuiltinExpressionContext) SetName(v antlr.Token) { s.name = v }

func (s *BuiltinExpressionContext) GetE() IExprContext { return s.e }

func (s *BuiltinExpressionContext) GetC() IClosureContext { return s.c }

func (s *BuiltinExpressionContext) SetE(v IExprContext) { s.e = v }

func (s *BuiltinExpressionContext) SetC(v IClosureContext) { s.c = v }

func (s *BuiltinExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BuiltinExpressionContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *BuiltinExpressionContext) Comma() antlr.TerminalNode {
	return s.GetToken(ExprParserComma, 0)
}

func (s *BuiltinExpressionContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *BuiltinExpressionContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BuiltinExpressionContext) Closure() IClosureContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IClosureContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IClosureContext)
}

func (s *BuiltinExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBuiltinExpression(s)
	}
}

func (s *BuiltinExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBuiltinExpression(s)
	}
}

func (s *BuiltinExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBuiltinExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) Builtins() (localctx IBuiltinsContext) {
	localctx = NewBuiltinsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, ExprParserRULE_builtins)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(148)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ExprParserT__1:
		localctx = NewLenBuiltinExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(101)
			p.Match(ExprParserT__1)
		}
		{
			p.SetState(102)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(103)

			var _x = p.expr(0)

			localctx.(*LenBuiltinExpressionContext).e = _x
		}
		{
			p.SetState(104)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__2:
		localctx = NewBuiltinExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(106)

			var _m = p.Match(ExprParserT__2)

			localctx.(*BuiltinExpressionContext).name = _m
		}
		{
			p.SetState(107)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(108)

			var _x = p.expr(0)

			localctx.(*BuiltinExpressionContext).e = _x
		}
		{
			p.SetState(109)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(110)

			var _x = p.Closure()

			localctx.(*BuiltinExpressionContext).c = _x
		}
		{
			p.SetState(111)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__3:
		localctx = NewBuiltinExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(113)

			var _m = p.Match(ExprParserT__3)

			localctx.(*BuiltinExpressionContext).name = _m
		}
		{
			p.SetState(114)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(115)

			var _x = p.expr(0)

			localctx.(*BuiltinExpressionContext).e = _x
		}
		{
			p.SetState(116)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(117)

			var _x = p.Closure()

			localctx.(*BuiltinExpressionContext).c = _x
		}
		{
			p.SetState(118)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__4:
		localctx = NewBuiltinExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(120)

			var _m = p.Match(ExprParserT__4)

			localctx.(*BuiltinExpressionContext).name = _m
		}
		{
			p.SetState(121)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(122)

			var _x = p.expr(0)

			localctx.(*BuiltinExpressionContext).e = _x
		}
		{
			p.SetState(123)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(124)

			var _x = p.Closure()

			localctx.(*BuiltinExpressionContext).c = _x
		}
		{
			p.SetState(125)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__5:
		localctx = NewBuiltinExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(127)

			var _m = p.Match(ExprParserT__5)

			localctx.(*BuiltinExpressionContext).name = _m
		}
		{
			p.SetState(128)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(129)

			var _x = p.expr(0)

			localctx.(*BuiltinExpressionContext).e = _x
		}
		{
			p.SetState(130)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(131)

			var _x = p.Closure()

			localctx.(*BuiltinExpressionContext).c = _x
		}
		{
			p.SetState(132)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__6:
		localctx = NewBuiltinExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(134)

			var _m = p.Match(ExprParserT__6)

			localctx.(*BuiltinExpressionContext).name = _m
		}
		{
			p.SetState(135)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(136)

			var _x = p.expr(0)

			localctx.(*BuiltinExpressionContext).e = _x
		}
		{
			p.SetState(137)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(138)

			var _x = p.Closure()

			localctx.(*BuiltinExpressionContext).c = _x
		}
		{
			p.SetState(139)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__7:
		localctx = NewBuiltinExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(141)

			var _m = p.Match(ExprParserT__7)

			localctx.(*BuiltinExpressionContext).name = _m
		}
		{
			p.SetState(142)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(143)

			var _x = p.expr(0)

			localctx.(*BuiltinExpressionContext).e = _x
		}
		{
			p.SetState(144)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(145)

			var _x = p.Closure()

			localctx.(*BuiltinExpressionContext).c = _x
		}
		{
			p.SetState(146)
			p.Match(ExprParserCloseParen)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IClosureContext is an interface to support dynamic dispatch.
type IClosureContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsClosureContext differentiates from other interfaces.
	IsClosureContext()
}

type ClosureContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClosureContext() *ClosureContext {
	var p = new(ClosureContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_closure
	return p
}

func (*ClosureContext) IsClosureContext() {}

func NewClosureContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClosureContext {
	var p = new(ClosureContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_closure

	return p
}

func (s *ClosureContext) GetParser() antlr.Parser { return s.parser }

func (s *ClosureContext) CopyFrom(ctx *ClosureContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ClosureContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClosureContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ClosureExpressionContext struct {
	*ClosureContext
	body IExprContext
}

func NewClosureExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ClosureExpressionContext {
	var p = new(ClosureExpressionContext)

	p.ClosureContext = NewEmptyClosureContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ClosureContext))

	return p
}

func (s *ClosureExpressionContext) GetBody() IExprContext { return s.body }

func (s *ClosureExpressionContext) SetBody(v IExprContext) { s.body = v }

func (s *ClosureExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClosureExpressionContext) OpenBrace() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenBrace, 0)
}

func (s *ClosureExpressionContext) CloseBrace() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseBrace, 0)
}

func (s *ClosureExpressionContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ClosureExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterClosureExpression(s)
	}
}

func (s *ClosureExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitClosureExpression(s)
	}
}

func (s *ClosureExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitClosureExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) Closure() (localctx IClosureContext) {
	localctx = NewClosureContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, ExprParserRULE_closure)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	localctx = NewClosureExpressionContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(150)
		p.Match(ExprParserOpenBrace)
	}
	{
		p.SetState(151)

		var _x = p.expr(0)

		localctx.(*ClosureExpressionContext).body = _x
	}
	{
		p.SetState(152)
		p.Match(ExprParserCloseBrace)
	}

	return localctx
}

// IArgumentsContext is an interface to support dynamic dispatch.
type IArgumentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_expr returns the _expr rule contexts.
	Get_expr() IExprContext

	// Set_expr sets the _expr rule contexts.
	Set_expr(IExprContext)

	// GetList returns the list rule context list.
	GetList() []IExprContext

	// SetList sets the list rule context list.
	SetList([]IExprContext)

	// IsArgumentsContext differentiates from other interfaces.
	IsArgumentsContext()
}

type ArgumentsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	_expr  IExprContext
	list   []IExprContext
}

func NewEmptyArgumentsContext() *ArgumentsContext {
	var p = new(ArgumentsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_arguments
	return p
}

func (*ArgumentsContext) IsArgumentsContext() {}

func NewArgumentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentsContext {
	var p = new(ArgumentsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_arguments

	return p
}

func (s *ArgumentsContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentsContext) Get_expr() IExprContext { return s._expr }

func (s *ArgumentsContext) Set_expr(v IExprContext) { s._expr = v }

func (s *ArgumentsContext) GetList() []IExprContext { return s.list }

func (s *ArgumentsContext) SetList(v []IExprContext) { s.list = v }

func (s *ArgumentsContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *ArgumentsContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ArgumentsContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(ExprParserComma)
}

func (s *ArgumentsContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(ExprParserComma, i)
}

func (s *ArgumentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterArguments(s)
	}
}

func (s *ArgumentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitArguments(s)
	}
}

func (s *ArgumentsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitArguments(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) Arguments() (localctx IArgumentsContext) {
	localctx = NewArgumentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, ExprParserRULE_arguments)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(154)

		var _x = p.expr(0)

		localctx.(*ArgumentsContext)._expr = _x
	}
	localctx.(*ArgumentsContext).list = append(localctx.(*ArgumentsContext).list, localctx.(*ArgumentsContext)._expr)
	p.SetState(159)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ExprParserComma {
		{
			p.SetState(155)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(156)

			var _x = p.expr(0)

			localctx.(*ArgumentsContext)._expr = _x
		}
		localctx.(*ArgumentsContext).list = append(localctx.(*ArgumentsContext).list, localctx.(*ArgumentsContext)._expr)

		p.SetState(161)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IArrayLiteralContext is an interface to support dynamic dispatch.
type IArrayLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_expr returns the _expr rule contexts.
	Get_expr() IExprContext

	// Set_expr sets the _expr rule contexts.
	Set_expr(IExprContext)

	// GetList returns the list rule context list.
	GetList() []IExprContext

	// SetList sets the list rule context list.
	SetList([]IExprContext)

	// IsArrayLiteralContext differentiates from other interfaces.
	IsArrayLiteralContext()
}

type ArrayLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	_expr  IExprContext
	list   []IExprContext
}

func NewEmptyArrayLiteralContext() *ArrayLiteralContext {
	var p = new(ArrayLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_arrayLiteral
	return p
}

func (*ArrayLiteralContext) IsArrayLiteralContext() {}

func NewArrayLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayLiteralContext {
	var p = new(ArrayLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_arrayLiteral

	return p
}

func (s *ArrayLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayLiteralContext) Get_expr() IExprContext { return s._expr }

func (s *ArrayLiteralContext) Set_expr(v IExprContext) { s._expr = v }

func (s *ArrayLiteralContext) GetList() []IExprContext { return s.list }

func (s *ArrayLiteralContext) SetList(v []IExprContext) { s.list = v }

func (s *ArrayLiteralContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenBracket, 0)
}

func (s *ArrayLiteralContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseBracket, 0)
}

func (s *ArrayLiteralContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *ArrayLiteralContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ArrayLiteralContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(ExprParserComma)
}

func (s *ArrayLiteralContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(ExprParserComma, i)
}

func (s *ArrayLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterArrayLiteral(s)
	}
}

func (s *ArrayLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitArrayLiteral(s)
	}
}

func (s *ArrayLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitArrayLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) ArrayLiteral() (localctx IArrayLiteralContext) {
	localctx = NewArrayLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, ExprParserRULE_arrayLiteral)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.SetState(178)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(162)
			p.Match(ExprParserOpenBracket)
		}
		{
			p.SetState(163)
			p.Match(ExprParserCloseBracket)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(164)
			p.Match(ExprParserOpenBracket)
		}
		{
			p.SetState(165)

			var _x = p.expr(0)

			localctx.(*ArrayLiteralContext)._expr = _x
		}
		localctx.(*ArrayLiteralContext).list = append(localctx.(*ArrayLiteralContext).list, localctx.(*ArrayLiteralContext)._expr)
		p.SetState(170)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(166)
					p.Match(ExprParserComma)
				}
				{
					p.SetState(167)

					var _x = p.expr(0)

					localctx.(*ArrayLiteralContext)._expr = _x
				}
				localctx.(*ArrayLiteralContext).list = append(localctx.(*ArrayLiteralContext).list, localctx.(*ArrayLiteralContext)._expr)

			}
			p.SetState(172)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())
		}
		p.SetState(174)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == ExprParserComma {
			{
				p.SetState(173)
				p.Match(ExprParserComma)
			}

		}
		{
			p.SetState(176)
			p.Match(ExprParserCloseBracket)
		}

	}

	return localctx
}

// IMapLiteralContext is an interface to support dynamic dispatch.
type IMapLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetE returns the e rule contexts.
	GetE() IPropertyNameAndValueListContext

	// SetE sets the e rule contexts.
	SetE(IPropertyNameAndValueListContext)

	// IsMapLiteralContext differentiates from other interfaces.
	IsMapLiteralContext()
}

type MapLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IPropertyNameAndValueListContext
}

func NewEmptyMapLiteralContext() *MapLiteralContext {
	var p = new(MapLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_mapLiteral
	return p
}

func (*MapLiteralContext) IsMapLiteralContext() {}

func NewMapLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapLiteralContext {
	var p = new(MapLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_mapLiteral

	return p
}

func (s *MapLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *MapLiteralContext) GetE() IPropertyNameAndValueListContext { return s.e }

func (s *MapLiteralContext) SetE(v IPropertyNameAndValueListContext) { s.e = v }

func (s *MapLiteralContext) OpenBrace() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenBrace, 0)
}

func (s *MapLiteralContext) CloseBrace() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseBrace, 0)
}

func (s *MapLiteralContext) PropertyNameAndValueList() IPropertyNameAndValueListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyNameAndValueListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPropertyNameAndValueListContext)
}

func (s *MapLiteralContext) Comma() antlr.TerminalNode {
	return s.GetToken(ExprParserComma, 0)
}

func (s *MapLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MapLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMapLiteral(s)
	}
}

func (s *MapLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMapLiteral(s)
	}
}

func (s *MapLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMapLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) MapLiteral() (localctx IMapLiteralContext) {
	localctx = NewMapLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, ExprParserRULE_mapLiteral)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(180)
			p.Match(ExprParserOpenBrace)
		}
		{
			p.SetState(181)
			p.Match(ExprParserCloseBrace)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(182)
			p.Match(ExprParserOpenBrace)
		}
		{
			p.SetState(183)

			var _x = p.PropertyNameAndValueList()

			localctx.(*MapLiteralContext).e = _x
		}
		p.SetState(185)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == ExprParserComma {
			{
				p.SetState(184)
				p.Match(ExprParserComma)
			}

		}
		{
			p.SetState(187)
			p.Match(ExprParserCloseBrace)
		}

	}

	return localctx
}

// IPropertyNameAndValueListContext is an interface to support dynamic dispatch.
type IPropertyNameAndValueListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_propertyAssignment returns the _propertyAssignment rule contexts.
	Get_propertyAssignment() IPropertyAssignmentContext

	// Set_propertyAssignment sets the _propertyAssignment rule contexts.
	Set_propertyAssignment(IPropertyAssignmentContext)

	// GetList returns the list rule context list.
	GetList() []IPropertyAssignmentContext

	// SetList sets the list rule context list.
	SetList([]IPropertyAssignmentContext)

	// IsPropertyNameAndValueListContext differentiates from other interfaces.
	IsPropertyNameAndValueListContext()
}

type PropertyNameAndValueListContext struct {
	*antlr.BaseParserRuleContext
	parser              antlr.Parser
	_propertyAssignment IPropertyAssignmentContext
	list                []IPropertyAssignmentContext
}

func NewEmptyPropertyNameAndValueListContext() *PropertyNameAndValueListContext {
	var p = new(PropertyNameAndValueListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_propertyNameAndValueList
	return p
}

func (*PropertyNameAndValueListContext) IsPropertyNameAndValueListContext() {}

func NewPropertyNameAndValueListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyNameAndValueListContext {
	var p = new(PropertyNameAndValueListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_propertyNameAndValueList

	return p
}

func (s *PropertyNameAndValueListContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyNameAndValueListContext) Get_propertyAssignment() IPropertyAssignmentContext {
	return s._propertyAssignment
}

func (s *PropertyNameAndValueListContext) Set_propertyAssignment(v IPropertyAssignmentContext) {
	s._propertyAssignment = v
}

func (s *PropertyNameAndValueListContext) GetList() []IPropertyAssignmentContext { return s.list }

func (s *PropertyNameAndValueListContext) SetList(v []IPropertyAssignmentContext) { s.list = v }

func (s *PropertyNameAndValueListContext) AllPropertyAssignment() []IPropertyAssignmentContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPropertyAssignmentContext)(nil)).Elem())
	var tst = make([]IPropertyAssignmentContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPropertyAssignmentContext)
		}
	}

	return tst
}

func (s *PropertyNameAndValueListContext) PropertyAssignment(i int) IPropertyAssignmentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyAssignmentContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPropertyAssignmentContext)
}

func (s *PropertyNameAndValueListContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(ExprParserComma)
}

func (s *PropertyNameAndValueListContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(ExprParserComma, i)
}

func (s *PropertyNameAndValueListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyNameAndValueListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyNameAndValueListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterPropertyNameAndValueList(s)
	}
}

func (s *PropertyNameAndValueListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitPropertyNameAndValueList(s)
	}
}

func (s *PropertyNameAndValueListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitPropertyNameAndValueList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) PropertyNameAndValueList() (localctx IPropertyNameAndValueListContext) {
	localctx = NewPropertyNameAndValueListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, ExprParserRULE_propertyNameAndValueList)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(191)

		var _x = p.PropertyAssignment()

		localctx.(*PropertyNameAndValueListContext)._propertyAssignment = _x
	}
	localctx.(*PropertyNameAndValueListContext).list = append(localctx.(*PropertyNameAndValueListContext).list, localctx.(*PropertyNameAndValueListContext)._propertyAssignment)
	p.SetState(196)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(192)
				p.Match(ExprParserComma)
			}
			{
				p.SetState(193)

				var _x = p.PropertyAssignment()

				localctx.(*PropertyNameAndValueListContext)._propertyAssignment = _x
			}
			localctx.(*PropertyNameAndValueListContext).list = append(localctx.(*PropertyNameAndValueListContext).list, localctx.(*PropertyNameAndValueListContext)._propertyAssignment)

		}
		p.SetState(198)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())
	}

	return localctx
}

// IPropertyAssignmentContext is an interface to support dynamic dispatch.
type IPropertyAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name rule contexts.
	GetName() IPropertyNameContext

	// GetValue returns the value rule contexts.
	GetValue() IExprContext

	// SetName sets the name rule contexts.
	SetName(IPropertyNameContext)

	// SetValue sets the value rule contexts.
	SetValue(IExprContext)

	// IsPropertyAssignmentContext differentiates from other interfaces.
	IsPropertyAssignmentContext()
}

type PropertyAssignmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	name   IPropertyNameContext
	value  IExprContext
}

func NewEmptyPropertyAssignmentContext() *PropertyAssignmentContext {
	var p = new(PropertyAssignmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_propertyAssignment
	return p
}

func (*PropertyAssignmentContext) IsPropertyAssignmentContext() {}

func NewPropertyAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyAssignmentContext {
	var p = new(PropertyAssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_propertyAssignment

	return p
}

func (s *PropertyAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyAssignmentContext) GetName() IPropertyNameContext { return s.name }

func (s *PropertyAssignmentContext) GetValue() IExprContext { return s.value }

func (s *PropertyAssignmentContext) SetName(v IPropertyNameContext) { s.name = v }

func (s *PropertyAssignmentContext) SetValue(v IExprContext) { s.value = v }

func (s *PropertyAssignmentContext) Colon() antlr.TerminalNode {
	return s.GetToken(ExprParserColon, 0)
}

func (s *PropertyAssignmentContext) PropertyName() IPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPropertyNameContext)
}

func (s *PropertyAssignmentContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *PropertyAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterPropertyAssignment(s)
	}
}

func (s *PropertyAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitPropertyAssignment(s)
	}
}

func (s *PropertyAssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitPropertyAssignment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) PropertyAssignment() (localctx IPropertyAssignmentContext) {
	localctx = NewPropertyAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, ExprParserRULE_propertyAssignment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(199)

		var _x = p.PropertyName()

		localctx.(*PropertyAssignmentContext).name = _x
	}
	{
		p.SetState(200)
		p.Match(ExprParserColon)
	}
	{
		p.SetState(201)

		var _x = p.expr(0)

		localctx.(*PropertyAssignmentContext).value = _x
	}

	return localctx
}

// IPropertyNameContext is an interface to support dynamic dispatch.
type IPropertyNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPropertyNameContext differentiates from other interfaces.
	IsPropertyNameContext()
}

type PropertyNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyNameContext() *PropertyNameContext {
	var p = new(PropertyNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_propertyName
	return p
}

func (*PropertyNameContext) IsPropertyNameContext() {}

func NewPropertyNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyNameContext {
	var p = new(PropertyNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_propertyName

	return p
}

func (s *PropertyNameContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyNameContext) Identifier() antlr.TerminalNode {
	return s.GetToken(ExprParserIdentifier, 0)
}

func (s *PropertyNameContext) StringLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserStringLiteral, 0)
}

func (s *PropertyNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterPropertyName(s)
	}
}

func (s *PropertyNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitPropertyName(s)
	}
}

func (s *PropertyNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitPropertyName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) PropertyName() (localctx IPropertyNameContext) {
	localctx = NewPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, ExprParserRULE_propertyName)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(203)
		_la = p.GetTokenStream().LA(1)

		if !(_la == ExprParserIdentifier || _la == ExprParserStringLiteral) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) CopyFrom(ctx *LiteralContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type BooleanExpressionContext struct {
	*LiteralContext
}

func NewBooleanExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BooleanExpressionContext {
	var p = new(BooleanExpressionContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *BooleanExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanExpressionContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserBooleanLiteral, 0)
}

func (s *BooleanExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBooleanExpression(s)
	}
}

func (s *BooleanExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBooleanExpression(s)
	}
}

func (s *BooleanExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBooleanExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringLiteralExpressionContext struct {
	*LiteralContext
}

func NewStringLiteralExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralExpressionContext {
	var p = new(StringLiteralExpressionContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *StringLiteralExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralExpressionContext) StringLiteral() IStringLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStringLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *StringLiteralExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterStringLiteralExpression(s)
	}
}

func (s *StringLiteralExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitStringLiteralExpression(s)
	}
}

func (s *StringLiteralExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitStringLiteralExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type NilExpressionContext struct {
	*LiteralContext
}

func NewNilExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NilExpressionContext {
	var p = new(NilExpressionContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *NilExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NilExpressionContext) NilLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserNilLiteral, 0)
}

func (s *NilExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterNilExpression(s)
	}
}

func (s *NilExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitNilExpression(s)
	}
}

func (s *NilExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitNilExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type FloatExpressionContext struct {
	*LiteralContext
}

func NewFloatExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FloatExpressionContext {
	var p = new(FloatExpressionContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *FloatExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FloatExpressionContext) FloatLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserFloatLiteral, 0)
}

func (s *FloatExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterFloatExpression(s)
	}
}

func (s *FloatExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitFloatExpression(s)
	}
}

func (s *FloatExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitFloatExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type IntegerExpressionContext struct {
	*LiteralContext
}

func NewIntegerExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntegerExpressionContext {
	var p = new(IntegerExpressionContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *IntegerExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntegerExpressionContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *IntegerExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterIntegerExpression(s)
	}
}

func (s *IntegerExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitIntegerExpression(s)
	}
}

func (s *IntegerExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitIntegerExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, ExprParserRULE_literal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(210)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ExprParserNilLiteral:
		localctx = NewNilExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(205)
			p.Match(ExprParserNilLiteral)
		}

	case ExprParserBooleanLiteral:
		localctx = NewBooleanExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(206)
			p.Match(ExprParserBooleanLiteral)
		}

	case ExprParserStringLiteral:
		localctx = NewStringLiteralExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(207)
			p.StringLiteral()
		}

	case ExprParserIntegerLiteral, ExprParserHexIntegerLiteral:
		localctx = NewIntegerExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(208)
			p.IntegerLiteral()
		}

	case ExprParserFloatLiteral:
		localctx = NewFloatExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(209)
			p.Match(ExprParserFloatLiteral)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IStringLiteralContext is an interface to support dynamic dispatch.
type IStringLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStringLiteralContext differentiates from other interfaces.
	IsStringLiteralContext()
}

type StringLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringLiteralContext() *StringLiteralContext {
	var p = new(StringLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_stringLiteral
	return p
}

func (*StringLiteralContext) IsStringLiteralContext() {}

func NewStringLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringLiteralContext {
	var p = new(StringLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_stringLiteral

	return p
}

func (s *StringLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *StringLiteralContext) StringLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserStringLiteral, 0)
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) StringLiteral() (localctx IStringLiteralContext) {
	localctx = NewStringLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, ExprParserRULE_stringLiteral)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(212)
		p.Match(ExprParserStringLiteral)
	}

	return localctx
}

// IIntegerLiteralContext is an interface to support dynamic dispatch.
type IIntegerLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntegerLiteralContext differentiates from other interfaces.
	IsIntegerLiteralContext()
}

type IntegerLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntegerLiteralContext() *IntegerLiteralContext {
	var p = new(IntegerLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExprParserRULE_integerLiteral
	return p
}

func (*IntegerLiteralContext) IsIntegerLiteralContext() {}

func NewIntegerLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntegerLiteralContext {
	var p = new(IntegerLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_integerLiteral

	return p
}

func (s *IntegerLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *IntegerLiteralContext) IntegerLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserIntegerLiteral, 0)
}

func (s *IntegerLiteralContext) HexIntegerLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserHexIntegerLiteral, 0)
}

func (s *IntegerLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntegerLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntegerLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitIntegerLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExprParser) IntegerLiteral() (localctx IIntegerLiteralContext) {
	localctx = NewIntegerLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, ExprParserRULE_integerLiteral)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(214)
		_la = p.GetTokenStream().LA(1)

		if !(_la == ExprParserIntegerLiteral || _la == ExprParserHexIntegerLiteral) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

func (p *ExprParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *ExprParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 17)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 16)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 15)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 14)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 13)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 12)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 11)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 10)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 8)

	case 10:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 11:
		return p.Precpred(p.GetParserRuleContext(), 22)

	case 12:
		return p.Precpred(p.GetParserRuleContext(), 21)

	case 13:
		return p.Precpred(p.GetParserRuleContext(), 19)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
