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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 61, 223,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 3, 2, 3, 2, 3, 2, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 46, 10, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 5, 3, 98, 10, 3, 3, 3, 3, 3, 5, 3, 102, 10, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 5, 3, 111, 10, 3, 3, 3, 7, 3, 114, 10, 3, 12, 3, 14,
	3, 117, 11, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 5, 4, 166, 10, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 7,
	6, 175, 10, 6, 12, 6, 14, 6, 178, 11, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7,
	3, 7, 7, 7, 186, 10, 7, 12, 7, 14, 7, 189, 11, 7, 3, 7, 5, 7, 192, 10,
	7, 3, 7, 3, 7, 5, 7, 196, 10, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 5, 8, 203,
	10, 8, 3, 8, 3, 8, 5, 8, 207, 10, 8, 3, 9, 3, 9, 3, 9, 7, 9, 212, 10, 9,
	12, 9, 14, 9, 215, 11, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3,
	11, 2, 3, 4, 12, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 2, 9, 3, 2, 29, 32,
	3, 2, 34, 37, 3, 2, 29, 30, 3, 2, 40, 43, 3, 2, 8, 9, 3, 2, 44, 45, 3,
	2, 55, 56, 2, 258, 2, 22, 3, 2, 2, 2, 4, 45, 3, 2, 2, 2, 6, 165, 3, 2,
	2, 2, 8, 167, 3, 2, 2, 2, 10, 171, 3, 2, 2, 2, 12, 195, 3, 2, 2, 2, 14,
	206, 3, 2, 2, 2, 16, 208, 3, 2, 2, 2, 18, 216, 3, 2, 2, 2, 20, 220, 3,
	2, 2, 2, 22, 23, 5, 4, 3, 2, 23, 24, 7, 2, 2, 3, 24, 3, 3, 2, 2, 2, 25,
	26, 8, 3, 1, 2, 26, 27, 7, 28, 2, 2, 27, 46, 7, 55, 2, 2, 28, 46, 5, 6,
	4, 2, 29, 30, 9, 2, 2, 2, 30, 46, 5, 4, 3, 27, 31, 46, 7, 33, 2, 2, 32,
	46, 7, 51, 2, 2, 33, 46, 7, 56, 2, 2, 34, 46, 7, 52, 2, 2, 35, 46, 7, 54,
	2, 2, 36, 46, 7, 53, 2, 2, 37, 46, 7, 55, 2, 2, 38, 46, 7, 46, 2, 2, 39,
	46, 5, 12, 7, 2, 40, 46, 5, 14, 8, 2, 41, 42, 7, 19, 2, 2, 42, 43, 5, 4,
	3, 2, 43, 44, 7, 20, 2, 2, 44, 46, 3, 2, 2, 2, 45, 25, 3, 2, 2, 2, 45,
	28, 3, 2, 2, 2, 45, 29, 3, 2, 2, 2, 45, 31, 3, 2, 2, 2, 45, 32, 3, 2, 2,
	2, 45, 33, 3, 2, 2, 2, 45, 34, 3, 2, 2, 2, 45, 35, 3, 2, 2, 2, 45, 36,
	3, 2, 2, 2, 45, 37, 3, 2, 2, 2, 45, 38, 3, 2, 2, 2, 45, 39, 3, 2, 2, 2,
	45, 40, 3, 2, 2, 2, 45, 41, 3, 2, 2, 2, 46, 115, 3, 2, 2, 2, 47, 48, 12,
	26, 2, 2, 48, 49, 7, 3, 2, 2, 49, 114, 5, 4, 3, 27, 50, 51, 12, 25, 2,
	2, 51, 52, 9, 3, 2, 2, 52, 114, 5, 4, 3, 26, 53, 54, 12, 24, 2, 2, 54,
	55, 9, 4, 2, 2, 55, 114, 5, 4, 3, 25, 56, 57, 12, 23, 2, 2, 57, 58, 9,
	5, 2, 2, 58, 114, 5, 4, 3, 24, 59, 60, 12, 22, 2, 2, 60, 61, 7, 4, 2, 2,
	61, 114, 5, 4, 3, 23, 62, 63, 12, 21, 2, 2, 63, 64, 7, 5, 2, 2, 64, 114,
	5, 4, 3, 22, 65, 66, 12, 20, 2, 2, 66, 67, 7, 6, 2, 2, 67, 114, 5, 4, 3,
	21, 68, 69, 12, 19, 2, 2, 69, 70, 7, 7, 2, 2, 70, 114, 5, 4, 3, 20, 71,
	72, 12, 18, 2, 2, 72, 73, 9, 6, 2, 2, 73, 114, 5, 4, 3, 19, 74, 75, 12,
	17, 2, 2, 75, 76, 9, 7, 2, 2, 76, 114, 5, 4, 3, 18, 77, 78, 12, 16, 2,
	2, 78, 79, 7, 47, 2, 2, 79, 114, 5, 4, 3, 17, 80, 81, 12, 15, 2, 2, 81,
	82, 7, 48, 2, 2, 82, 114, 5, 4, 3, 16, 83, 84, 12, 14, 2, 2, 84, 85, 7,
	26, 2, 2, 85, 86, 5, 4, 3, 2, 86, 87, 7, 27, 2, 2, 87, 88, 5, 4, 3, 15,
	88, 114, 3, 2, 2, 2, 89, 90, 12, 32, 2, 2, 90, 91, 7, 17, 2, 2, 91, 92,
	5, 4, 3, 2, 92, 93, 7, 18, 2, 2, 93, 114, 3, 2, 2, 2, 94, 95, 12, 31, 2,
	2, 95, 97, 7, 17, 2, 2, 96, 98, 5, 4, 3, 2, 97, 96, 3, 2, 2, 2, 97, 98,
	3, 2, 2, 2, 98, 99, 3, 2, 2, 2, 99, 101, 7, 27, 2, 2, 100, 102, 5, 4, 3,
	2, 101, 100, 3, 2, 2, 2, 101, 102, 3, 2, 2, 2, 102, 103, 3, 2, 2, 2, 103,
	114, 7, 18, 2, 2, 104, 105, 12, 30, 2, 2, 105, 106, 7, 28, 2, 2, 106, 114,
	7, 55, 2, 2, 107, 108, 12, 28, 2, 2, 108, 110, 7, 19, 2, 2, 109, 111, 5,
	10, 6, 2, 110, 109, 3, 2, 2, 2, 110, 111, 3, 2, 2, 2, 111, 112, 3, 2, 2,
	2, 112, 114, 7, 20, 2, 2, 113, 47, 3, 2, 2, 2, 113, 50, 3, 2, 2, 2, 113,
	53, 3, 2, 2, 2, 113, 56, 3, 2, 2, 2, 113, 59, 3, 2, 2, 2, 113, 62, 3, 2,
	2, 2, 113, 65, 3, 2, 2, 2, 113, 68, 3, 2, 2, 2, 113, 71, 3, 2, 2, 2, 113,
	74, 3, 2, 2, 2, 113, 77, 3, 2, 2, 2, 113, 80, 3, 2, 2, 2, 113, 83, 3, 2,
	2, 2, 113, 89, 3, 2, 2, 2, 113, 94, 3, 2, 2, 2, 113, 104, 3, 2, 2, 2, 113,
	107, 3, 2, 2, 2, 114, 117, 3, 2, 2, 2, 115, 113, 3, 2, 2, 2, 115, 116,
	3, 2, 2, 2, 116, 5, 3, 2, 2, 2, 117, 115, 3, 2, 2, 2, 118, 119, 7, 10,
	2, 2, 119, 120, 7, 19, 2, 2, 120, 121, 5, 4, 3, 2, 121, 122, 7, 20, 2,
	2, 122, 166, 3, 2, 2, 2, 123, 124, 7, 11, 2, 2, 124, 125, 7, 19, 2, 2,
	125, 126, 5, 4, 3, 2, 126, 127, 7, 24, 2, 2, 127, 128, 5, 8, 5, 2, 128,
	129, 7, 20, 2, 2, 129, 166, 3, 2, 2, 2, 130, 131, 7, 12, 2, 2, 131, 132,
	7, 19, 2, 2, 132, 133, 5, 4, 3, 2, 133, 134, 7, 24, 2, 2, 134, 135, 5,
	8, 5, 2, 135, 136, 7, 20, 2, 2, 136, 166, 3, 2, 2, 2, 137, 138, 7, 13,
	2, 2, 138, 139, 7, 19, 2, 2, 139, 140, 5, 4, 3, 2, 140, 141, 7, 24, 2,
	2, 141, 142, 5, 8, 5, 2, 142, 143, 7, 20, 2, 2, 143, 166, 3, 2, 2, 2, 144,
	145, 7, 14, 2, 2, 145, 146, 7, 19, 2, 2, 146, 147, 5, 4, 3, 2, 147, 148,
	7, 24, 2, 2, 148, 149, 5, 8, 5, 2, 149, 150, 7, 20, 2, 2, 150, 166, 3,
	2, 2, 2, 151, 152, 7, 15, 2, 2, 152, 153, 7, 19, 2, 2, 153, 154, 5, 4,
	3, 2, 154, 155, 7, 24, 2, 2, 155, 156, 5, 8, 5, 2, 156, 157, 7, 20, 2,
	2, 157, 166, 3, 2, 2, 2, 158, 159, 7, 16, 2, 2, 159, 160, 7, 19, 2, 2,
	160, 161, 5, 4, 3, 2, 161, 162, 7, 24, 2, 2, 162, 163, 5, 8, 5, 2, 163,
	164, 7, 20, 2, 2, 164, 166, 3, 2, 2, 2, 165, 118, 3, 2, 2, 2, 165, 123,
	3, 2, 2, 2, 165, 130, 3, 2, 2, 2, 165, 137, 3, 2, 2, 2, 165, 144, 3, 2,
	2, 2, 165, 151, 3, 2, 2, 2, 165, 158, 3, 2, 2, 2, 166, 7, 3, 2, 2, 2, 167,
	168, 7, 21, 2, 2, 168, 169, 5, 4, 3, 2, 169, 170, 7, 22, 2, 2, 170, 9,
	3, 2, 2, 2, 171, 176, 5, 4, 3, 2, 172, 173, 7, 24, 2, 2, 173, 175, 5, 4,
	3, 2, 174, 172, 3, 2, 2, 2, 175, 178, 3, 2, 2, 2, 176, 174, 3, 2, 2, 2,
	176, 177, 3, 2, 2, 2, 177, 11, 3, 2, 2, 2, 178, 176, 3, 2, 2, 2, 179, 180,
	7, 17, 2, 2, 180, 196, 7, 18, 2, 2, 181, 182, 7, 17, 2, 2, 182, 187, 5,
	4, 3, 2, 183, 184, 7, 24, 2, 2, 184, 186, 5, 4, 3, 2, 185, 183, 3, 2, 2,
	2, 186, 189, 3, 2, 2, 2, 187, 185, 3, 2, 2, 2, 187, 188, 3, 2, 2, 2, 188,
	191, 3, 2, 2, 2, 189, 187, 3, 2, 2, 2, 190, 192, 7, 24, 2, 2, 191, 190,
	3, 2, 2, 2, 191, 192, 3, 2, 2, 2, 192, 193, 3, 2, 2, 2, 193, 194, 7, 18,
	2, 2, 194, 196, 3, 2, 2, 2, 195, 179, 3, 2, 2, 2, 195, 181, 3, 2, 2, 2,
	196, 13, 3, 2, 2, 2, 197, 198, 7, 21, 2, 2, 198, 207, 7, 22, 2, 2, 199,
	200, 7, 21, 2, 2, 200, 202, 5, 16, 9, 2, 201, 203, 7, 24, 2, 2, 202, 201,
	3, 2, 2, 2, 202, 203, 3, 2, 2, 2, 203, 204, 3, 2, 2, 2, 204, 205, 7, 22,
	2, 2, 205, 207, 3, 2, 2, 2, 206, 197, 3, 2, 2, 2, 206, 199, 3, 2, 2, 2,
	207, 15, 3, 2, 2, 2, 208, 213, 5, 18, 10, 2, 209, 210, 7, 24, 2, 2, 210,
	212, 5, 18, 10, 2, 211, 209, 3, 2, 2, 2, 212, 215, 3, 2, 2, 2, 213, 211,
	3, 2, 2, 2, 213, 214, 3, 2, 2, 2, 214, 17, 3, 2, 2, 2, 215, 213, 3, 2,
	2, 2, 216, 217, 5, 20, 11, 2, 217, 218, 7, 27, 2, 2, 218, 219, 5, 4, 3,
	2, 219, 19, 3, 2, 2, 2, 220, 221, 9, 8, 2, 2, 221, 21, 3, 2, 2, 2, 16,
	45, 97, 101, 110, 113, 115, 165, 176, 187, 191, 195, 202, 206, 213,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'..'", "'startsWith'", "'endsWith'", "'contains'", "'matches'", "'in'",
	"'not in'", "'len'", "'all'", "'none'", "'any'", "'one'", "'filter'", "'map'",
	"'['", "']'", "'('", "')'", "'{'", "'}'", "';'", "','", "'='", "'?'", "':'",
	"'.'", "'+'", "'-'", "'!'", "'not'", "'nil'", "'*'", "'**'", "'/'", "'%'",
	"'>>'", "'<<'", "'<'", "'>'", "'<='", "'>='", "'=='", "'!='", "'#'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "OpenBracket",
	"CloseBracket", "OpenParen", "CloseParen", "OpenBrace", "CloseBrace", "SemiColon",
	"Comma", "Assign", "QuestionMark", "Colon", "Dot", "Plus", "Minus", "Negate",
	"Not", "Nil", "Multiply", "Exponent", "Divide", "Modulus", "RightShiftArithmetic",
	"LeftShiftArithmetic", "LessThan", "MoreThan", "LessThanEquals", "GreaterThanEquals",
	"Equals", "NotEquals", "Pointer", "And", "Or", "Builtins", "Ops", "BooleanLiteral",
	"IntegerLiteral", "FloatLiteral", "HexIntegerLiteral", "Identifier", "StringLiteral",
	"WhiteSpaces", "MultiLineComment", "SingleLineComment", "LineTerminator",
	"UnexpectedCharacter",
}

var ruleNames = []string{
	"start", "expr", "builtins", "closure", "arguments", "arrayLiteral", "mapLiteral",
	"propertyNameAndValueList", "propertyAssignment", "propertyName",
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
	ExprParserT__8                 = 9
	ExprParserT__9                 = 10
	ExprParserT__10                = 11
	ExprParserT__11                = 12
	ExprParserT__12                = 13
	ExprParserT__13                = 14
	ExprParserOpenBracket          = 15
	ExprParserCloseBracket         = 16
	ExprParserOpenParen            = 17
	ExprParserCloseParen           = 18
	ExprParserOpenBrace            = 19
	ExprParserCloseBrace           = 20
	ExprParserSemiColon            = 21
	ExprParserComma                = 22
	ExprParserAssign               = 23
	ExprParserQuestionMark         = 24
	ExprParserColon                = 25
	ExprParserDot                  = 26
	ExprParserPlus                 = 27
	ExprParserMinus                = 28
	ExprParserNegate               = 29
	ExprParserNot                  = 30
	ExprParserNil                  = 31
	ExprParserMultiply             = 32
	ExprParserExponent             = 33
	ExprParserDivide               = 34
	ExprParserModulus              = 35
	ExprParserRightShiftArithmetic = 36
	ExprParserLeftShiftArithmetic  = 37
	ExprParserLessThan             = 38
	ExprParserMoreThan             = 39
	ExprParserLessThanEquals       = 40
	ExprParserGreaterThanEquals    = 41
	ExprParserEquals               = 42
	ExprParserNotEquals            = 43
	ExprParserPointer              = 44
	ExprParserAnd                  = 45
	ExprParserOr                   = 46
	ExprParserBuiltins             = 47
	ExprParserOps                  = 48
	ExprParserBooleanLiteral       = 49
	ExprParserIntegerLiteral       = 50
	ExprParserFloatLiteral         = 51
	ExprParserHexIntegerLiteral    = 52
	ExprParserIdentifier           = 53
	ExprParserStringLiteral        = 54
	ExprParserWhiteSpaces          = 55
	ExprParserMultiLineComment     = 56
	ExprParserSingleLineComment    = 57
	ExprParserLineTerminator       = 58
	ExprParserUnexpectedCharacter  = 59
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
		p.SetState(20)

		var _x = p.expr(0)

		localctx.(*StartContext).e = _x
	}
	{
		p.SetState(21)
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

type CallContext struct {
	*ExprContext
	args IArgumentsContext
}

func NewCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CallContext {
	var p = new(CallContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *CallContext) GetArgs() IArgumentsContext { return s.args }

func (s *CallContext) SetArgs(v IArgumentsContext) { s.args = v }

func (s *CallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CallContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *CallContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *CallContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *CallContext) Arguments() IArgumentsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgumentsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgumentsContext)
}

func (s *CallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterCall(s)
	}
}

func (s *CallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitCall(s)
	}
}

func (s *CallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type MatchesContext struct {
	*ExprContext
	op      antlr.Token
	pattern IExprContext
}

func NewMatchesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MatchesContext {
	var p = new(MatchesContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MatchesContext) GetOp() antlr.Token { return s.op }

func (s *MatchesContext) SetOp(v antlr.Token) { s.op = v }

func (s *MatchesContext) GetPattern() IExprContext { return s.pattern }

func (s *MatchesContext) SetPattern(v IExprContext) { s.pattern = v }

func (s *MatchesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MatchesContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *MatchesContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MatchesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMatches(s)
	}
}

func (s *MatchesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMatches(s)
	}
}

func (s *MatchesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMatches(s)

	default:
		return t.VisitChildren(s)
	}
}

type TernaryContext struct {
	*ExprContext
	e1 IExprContext
	e2 IExprContext
}

func NewTernaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TernaryContext {
	var p = new(TernaryContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *TernaryContext) GetE1() IExprContext { return s.e1 }

func (s *TernaryContext) GetE2() IExprContext { return s.e2 }

func (s *TernaryContext) SetE1(v IExprContext) { s.e1 = v }

func (s *TernaryContext) SetE2(v IExprContext) { s.e2 = v }

func (s *TernaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TernaryContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *TernaryContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *TernaryContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(ExprParserQuestionMark, 0)
}

func (s *TernaryContext) Colon() antlr.TerminalNode {
	return s.GetToken(ExprParserColon, 0)
}

func (s *TernaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterTernary(s)
	}
}

func (s *TernaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitTernary(s)
	}
}

func (s *TernaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitTernary(s)

	default:
		return t.VisitChildren(s)
	}
}

type PointerContext struct {
	*ExprContext
}

func NewPointerContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PointerContext {
	var p = new(PointerContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *PointerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PointerContext) Pointer() antlr.TerminalNode {
	return s.GetToken(ExprParserPointer, 0)
}

func (s *PointerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterPointer(s)
	}
}

func (s *PointerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitPointer(s)
	}
}

func (s *PointerContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitPointer(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringContext struct {
	*ExprContext
}

func NewStringContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringContext {
	var p = new(StringContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *StringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringContext) StringLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserStringLiteral, 0)
}

func (s *StringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterString(s)
	}
}

func (s *StringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitString(s)
	}
}

func (s *StringContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitString(s)

	default:
		return t.VisitChildren(s)
	}
}

type ClosureMemberDotContext struct {
	*ExprContext
	name antlr.Token
}

func NewClosureMemberDotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ClosureMemberDotContext {
	var p = new(ClosureMemberDotContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *ClosureMemberDotContext) GetName() antlr.Token { return s.name }

func (s *ClosureMemberDotContext) SetName(v antlr.Token) { s.name = v }

func (s *ClosureMemberDotContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClosureMemberDotContext) Dot() antlr.TerminalNode {
	return s.GetToken(ExprParserDot, 0)
}

func (s *ClosureMemberDotContext) Identifier() antlr.TerminalNode {
	return s.GetToken(ExprParserIdentifier, 0)
}

func (s *ClosureMemberDotContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterClosureMemberDot(s)
	}
}

func (s *ClosureMemberDotContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitClosureMemberDot(s)
	}
}

func (s *ClosureMemberDotContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitClosureMemberDot(s)

	default:
		return t.VisitChildren(s)
	}
}

type UnaryContext struct {
	*ExprContext
	op antlr.Token
}

func NewUnaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UnaryContext {
	var p = new(UnaryContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *UnaryContext) GetOp() antlr.Token { return s.op }

func (s *UnaryContext) SetOp(v antlr.Token) { s.op = v }

func (s *UnaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *UnaryContext) Plus() antlr.TerminalNode {
	return s.GetToken(ExprParserPlus, 0)
}

func (s *UnaryContext) Minus() antlr.TerminalNode {
	return s.GetToken(ExprParserMinus, 0)
}

func (s *UnaryContext) Negate() antlr.TerminalNode {
	return s.GetToken(ExprParserNegate, 0)
}

func (s *UnaryContext) Not() antlr.TerminalNode {
	return s.GetToken(ExprParserNot, 0)
}

func (s *UnaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterUnary(s)
	}
}

func (s *UnaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitUnary(s)
	}
}

func (s *UnaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitUnary(s)

	default:
		return t.VisitChildren(s)
	}
}

type NilContext struct {
	*ExprContext
}

func NewNilContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NilContext {
	var p = new(NilContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *NilContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NilContext) Nil() antlr.TerminalNode {
	return s.GetToken(ExprParserNil, 0)
}

func (s *NilContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterNil(s)
	}
}

func (s *NilContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitNil(s)
	}
}

func (s *NilContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitNil(s)

	default:
		return t.VisitChildren(s)
	}
}

type IntegerContext struct {
	*ExprContext
}

func NewIntegerContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntegerContext {
	var p = new(IntegerContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *IntegerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntegerContext) IntegerLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserIntegerLiteral, 0)
}

func (s *IntegerContext) HexIntegerLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserHexIntegerLiteral, 0)
}

func (s *IntegerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterInteger(s)
	}
}

func (s *IntegerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitInteger(s)
	}
}

func (s *IntegerContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitInteger(s)

	default:
		return t.VisitChildren(s)
	}
}

type ArrayContext struct {
	*ExprContext
}

func NewArrayContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ArrayContext {
	var p = new(ArrayContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *ArrayContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayContext) ArrayLiteral() IArrayLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *ArrayContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterArray(s)
	}
}

func (s *ArrayContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitArray(s)
	}
}

func (s *ArrayContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitArray(s)

	default:
		return t.VisitChildren(s)
	}
}

type FloatContext struct {
	*ExprContext
}

func NewFloatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FloatContext {
	var p = new(FloatContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *FloatContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FloatContext) FloatLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserFloatLiteral, 0)
}

func (s *FloatContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterFloat(s)
	}
}

func (s *FloatContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitFloat(s)
	}
}

func (s *FloatContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitFloat(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdentifierContext struct {
	*ExprContext
}

func NewIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentifierContext {
	var p = new(IdentifierContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *IdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierContext) Identifier() antlr.TerminalNode {
	return s.GetToken(ExprParserIdentifier, 0)
}

func (s *IdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterIdentifier(s)
	}
}

func (s *IdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitIdentifier(s)
	}
}

func (s *IdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

type ParenthesizedContext struct {
	*ExprContext
}

func NewParenthesizedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenthesizedContext {
	var p = new(ParenthesizedContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *ParenthesizedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenthesizedContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *ParenthesizedContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ParenthesizedContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *ParenthesizedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterParenthesized(s)
	}
}

func (s *ParenthesizedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitParenthesized(s)
	}
}

func (s *ParenthesizedContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitParenthesized(s)

	default:
		return t.VisitChildren(s)
	}
}

type SliceContext struct {
	*ExprContext
	a IExprContext
	b IExprContext
}

func NewSliceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SliceContext {
	var p = new(SliceContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *SliceContext) GetA() IExprContext { return s.a }

func (s *SliceContext) GetB() IExprContext { return s.b }

func (s *SliceContext) SetA(v IExprContext) { s.a = v }

func (s *SliceContext) SetB(v IExprContext) { s.b = v }

func (s *SliceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SliceContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *SliceContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SliceContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenBracket, 0)
}

func (s *SliceContext) Colon() antlr.TerminalNode {
	return s.GetToken(ExprParserColon, 0)
}

func (s *SliceContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseBracket, 0)
}

func (s *SliceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterSlice(s)
	}
}

func (s *SliceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitSlice(s)
	}
}

func (s *SliceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitSlice(s)

	default:
		return t.VisitChildren(s)
	}
}

type MemberIndexContext struct {
	*ExprContext
	index IExprContext
}

func NewMemberIndexContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberIndexContext {
	var p = new(MemberIndexContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MemberIndexContext) GetIndex() IExprContext { return s.index }

func (s *MemberIndexContext) SetIndex(v IExprContext) { s.index = v }

func (s *MemberIndexContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberIndexContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *MemberIndexContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MemberIndexContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenBracket, 0)
}

func (s *MemberIndexContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseBracket, 0)
}

func (s *MemberIndexContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMemberIndex(s)
	}
}

func (s *MemberIndexContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMemberIndex(s)
	}
}

func (s *MemberIndexContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMemberIndex(s)

	default:
		return t.VisitChildren(s)
	}
}

type BuiltinsListContext struct {
	*ExprContext
}

func NewBuiltinsListContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BuiltinsListContext {
	var p = new(BuiltinsListContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *BuiltinsListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BuiltinsListContext) Builtins() IBuiltinsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBuiltinsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBuiltinsContext)
}

func (s *BuiltinsListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBuiltinsList(s)
	}
}

func (s *BuiltinsListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBuiltinsList(s)
	}
}

func (s *BuiltinsListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBuiltinsList(s)

	default:
		return t.VisitChildren(s)
	}
}

type BinaryContext struct {
	*ExprContext
	op antlr.Token
}

func NewBinaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BinaryContext {
	var p = new(BinaryContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *BinaryContext) GetOp() antlr.Token { return s.op }

func (s *BinaryContext) SetOp(v antlr.Token) { s.op = v }

func (s *BinaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *BinaryContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BinaryContext) Multiply() antlr.TerminalNode {
	return s.GetToken(ExprParserMultiply, 0)
}

func (s *BinaryContext) Exponent() antlr.TerminalNode {
	return s.GetToken(ExprParserExponent, 0)
}

func (s *BinaryContext) Divide() antlr.TerminalNode {
	return s.GetToken(ExprParserDivide, 0)
}

func (s *BinaryContext) Modulus() antlr.TerminalNode {
	return s.GetToken(ExprParserModulus, 0)
}

func (s *BinaryContext) Plus() antlr.TerminalNode {
	return s.GetToken(ExprParserPlus, 0)
}

func (s *BinaryContext) Minus() antlr.TerminalNode {
	return s.GetToken(ExprParserMinus, 0)
}

func (s *BinaryContext) LessThan() antlr.TerminalNode {
	return s.GetToken(ExprParserLessThan, 0)
}

func (s *BinaryContext) MoreThan() antlr.TerminalNode {
	return s.GetToken(ExprParserMoreThan, 0)
}

func (s *BinaryContext) LessThanEquals() antlr.TerminalNode {
	return s.GetToken(ExprParserLessThanEquals, 0)
}

func (s *BinaryContext) GreaterThanEquals() antlr.TerminalNode {
	return s.GetToken(ExprParserGreaterThanEquals, 0)
}

func (s *BinaryContext) Equals() antlr.TerminalNode {
	return s.GetToken(ExprParserEquals, 0)
}

func (s *BinaryContext) NotEquals() antlr.TerminalNode {
	return s.GetToken(ExprParserNotEquals, 0)
}

func (s *BinaryContext) And() antlr.TerminalNode {
	return s.GetToken(ExprParserAnd, 0)
}

func (s *BinaryContext) Or() antlr.TerminalNode {
	return s.GetToken(ExprParserOr, 0)
}

func (s *BinaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBinary(s)
	}
}

func (s *BinaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBinary(s)
	}
}

func (s *BinaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBinary(s)

	default:
		return t.VisitChildren(s)
	}
}

type BooleanContext struct {
	*ExprContext
}

func NewBooleanContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BooleanContext {
	var p = new(BooleanContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *BooleanContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(ExprParserBooleanLiteral, 0)
}

func (s *BooleanContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBoolean(s)
	}
}

func (s *BooleanContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBoolean(s)
	}
}

func (s *BooleanContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBoolean(s)

	default:
		return t.VisitChildren(s)
	}
}

type MapContext struct {
	*ExprContext
}

func NewMapContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MapContext {
	var p = new(MapContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MapContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapContext) MapLiteral() IMapLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMapLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMapLiteralContext)
}

func (s *MapContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMap(s)
	}
}

func (s *MapContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMap(s)
	}
}

func (s *MapContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMap(s)

	default:
		return t.VisitChildren(s)
	}
}

type MemberDotContext struct {
	*ExprContext
	name antlr.Token
}

func NewMemberDotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberDotContext {
	var p = new(MemberDotContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *MemberDotContext) GetName() antlr.Token { return s.name }

func (s *MemberDotContext) SetName(v antlr.Token) { s.name = v }

func (s *MemberDotContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberDotContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MemberDotContext) Dot() antlr.TerminalNode {
	return s.GetToken(ExprParserDot, 0)
}

func (s *MemberDotContext) Identifier() antlr.TerminalNode {
	return s.GetToken(ExprParserIdentifier, 0)
}

func (s *MemberDotContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterMemberDot(s)
	}
}

func (s *MemberDotContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitMemberDot(s)
	}
}

func (s *MemberDotContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitMemberDot(s)

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
	p.SetState(43)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ExprParserDot:
		localctx = NewClosureMemberDotContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(24)
			p.Match(ExprParserDot)
		}
		{
			p.SetState(25)

			var _m = p.Match(ExprParserIdentifier)

			localctx.(*ClosureMemberDotContext).name = _m
		}

	case ExprParserT__7, ExprParserT__8, ExprParserT__9, ExprParserT__10, ExprParserT__11, ExprParserT__12, ExprParserT__13:
		localctx = NewBuiltinsListContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(26)
			p.Builtins()
		}

	case ExprParserPlus, ExprParserMinus, ExprParserNegate, ExprParserNot:
		localctx = NewUnaryContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(27)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*UnaryContext).op = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExprParserPlus)|(1<<ExprParserMinus)|(1<<ExprParserNegate)|(1<<ExprParserNot))) != 0) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*UnaryContext).op = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(28)
			p.expr(25)
		}

	case ExprParserNil:
		localctx = NewNilContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(29)
			p.Match(ExprParserNil)
		}

	case ExprParserBooleanLiteral:
		localctx = NewBooleanContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(30)
			p.Match(ExprParserBooleanLiteral)
		}

	case ExprParserStringLiteral:
		localctx = NewStringContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(31)
			p.Match(ExprParserStringLiteral)
		}

	case ExprParserIntegerLiteral:
		localctx = NewIntegerContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(32)
			p.Match(ExprParserIntegerLiteral)
		}

	case ExprParserHexIntegerLiteral:
		localctx = NewIntegerContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(33)
			p.Match(ExprParserHexIntegerLiteral)
		}

	case ExprParserFloatLiteral:
		localctx = NewFloatContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(34)
			p.Match(ExprParserFloatLiteral)
		}

	case ExprParserIdentifier:
		localctx = NewIdentifierContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(35)
			p.Match(ExprParserIdentifier)
		}

	case ExprParserPointer:
		localctx = NewPointerContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(36)
			p.Match(ExprParserPointer)
		}

	case ExprParserOpenBracket:
		localctx = NewArrayContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(37)
			p.ArrayLiteral()
		}

	case ExprParserOpenBrace:
		localctx = NewMapContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(38)
			p.MapLiteral()
		}

	case ExprParserOpenParen:
		localctx = NewParenthesizedContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(39)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(40)
			p.expr(0)
		}
		{
			p.SetState(41)
			p.Match(ExprParserCloseParen)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(111)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
			case 1:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(45)

				if !(p.Precpred(p.GetParserRuleContext(), 24)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 24)", ""))
				}
				{
					p.SetState(46)

					var _m = p.Match(ExprParserT__0)

					localctx.(*BinaryContext).op = _m
				}
				{
					p.SetState(47)
					p.expr(25)
				}

			case 2:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(48)

				if !(p.Precpred(p.GetParserRuleContext(), 23)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 23)", ""))
				}
				{
					p.SetState(49)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*BinaryContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(ExprParserMultiply-32))|(1<<(ExprParserExponent-32))|(1<<(ExprParserDivide-32))|(1<<(ExprParserModulus-32)))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*BinaryContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(50)
					p.expr(24)
				}

			case 3:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(51)

				if !(p.Precpred(p.GetParserRuleContext(), 22)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 22)", ""))
				}
				{
					p.SetState(52)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*BinaryContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserPlus || _la == ExprParserMinus) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*BinaryContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(53)
					p.expr(23)
				}

			case 4:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(54)

				if !(p.Precpred(p.GetParserRuleContext(), 21)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 21)", ""))
				}
				{
					p.SetState(55)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*BinaryContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la-38)&-(0x1f+1)) == 0 && ((1<<uint((_la-38)))&((1<<(ExprParserLessThan-38))|(1<<(ExprParserMoreThan-38))|(1<<(ExprParserLessThanEquals-38))|(1<<(ExprParserGreaterThanEquals-38)))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*BinaryContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(56)
					p.expr(22)
				}

			case 5:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(57)

				if !(p.Precpred(p.GetParserRuleContext(), 20)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 20)", ""))
				}
				{
					p.SetState(58)

					var _m = p.Match(ExprParserT__1)

					localctx.(*BinaryContext).op = _m
				}
				{
					p.SetState(59)
					p.expr(21)
				}

			case 6:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(60)

				if !(p.Precpred(p.GetParserRuleContext(), 19)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 19)", ""))
				}
				{
					p.SetState(61)

					var _m = p.Match(ExprParserT__2)

					localctx.(*BinaryContext).op = _m
				}
				{
					p.SetState(62)
					p.expr(20)
				}

			case 7:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(63)

				if !(p.Precpred(p.GetParserRuleContext(), 18)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 18)", ""))
				}
				{
					p.SetState(64)

					var _m = p.Match(ExprParserT__3)

					localctx.(*BinaryContext).op = _m
				}
				{
					p.SetState(65)
					p.expr(19)
				}

			case 8:
				localctx = NewMatchesContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(66)

				if !(p.Precpred(p.GetParserRuleContext(), 17)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 17)", ""))
				}
				{
					p.SetState(67)

					var _m = p.Match(ExprParserT__4)

					localctx.(*MatchesContext).op = _m
				}
				{
					p.SetState(68)

					var _x = p.expr(18)

					localctx.(*MatchesContext).pattern = _x
				}

			case 9:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(69)

				if !(p.Precpred(p.GetParserRuleContext(), 16)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 16)", ""))
				}
				{
					p.SetState(70)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*BinaryContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserT__5 || _la == ExprParserT__6) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*BinaryContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(71)
					p.expr(17)
				}

			case 10:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(72)

				if !(p.Precpred(p.GetParserRuleContext(), 15)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 15)", ""))
				}
				{
					p.SetState(73)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*BinaryContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserEquals || _la == ExprParserNotEquals) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*BinaryContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(74)
					p.expr(16)
				}

			case 11:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(75)

				if !(p.Precpred(p.GetParserRuleContext(), 14)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 14)", ""))
				}
				{
					p.SetState(76)

					var _m = p.Match(ExprParserAnd)

					localctx.(*BinaryContext).op = _m
				}
				{
					p.SetState(77)
					p.expr(15)
				}

			case 12:
				localctx = NewBinaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(78)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
				}
				{
					p.SetState(79)

					var _m = p.Match(ExprParserOr)

					localctx.(*BinaryContext).op = _m
				}
				{
					p.SetState(80)
					p.expr(14)
				}

			case 13:
				localctx = NewTernaryContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(81)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(82)
					p.Match(ExprParserQuestionMark)
				}
				{
					p.SetState(83)

					var _x = p.expr(0)

					localctx.(*TernaryContext).e1 = _x
				}
				{
					p.SetState(84)
					p.Match(ExprParserColon)
				}
				{
					p.SetState(85)

					var _x = p.expr(13)

					localctx.(*TernaryContext).e2 = _x
				}

			case 14:
				localctx = NewMemberIndexContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(87)

				if !(p.Precpred(p.GetParserRuleContext(), 30)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 30)", ""))
				}
				{
					p.SetState(88)
					p.Match(ExprParserOpenBracket)
				}
				{
					p.SetState(89)

					var _x = p.expr(0)

					localctx.(*MemberIndexContext).index = _x
				}
				{
					p.SetState(90)
					p.Match(ExprParserCloseBracket)
				}

			case 15:
				localctx = NewSliceContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(92)

				if !(p.Precpred(p.GetParserRuleContext(), 29)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 29)", ""))
				}
				{
					p.SetState(93)
					p.Match(ExprParserOpenBracket)
				}
				p.SetState(95)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExprParserT__7)|(1<<ExprParserT__8)|(1<<ExprParserT__9)|(1<<ExprParserT__10)|(1<<ExprParserT__11)|(1<<ExprParserT__12)|(1<<ExprParserT__13)|(1<<ExprParserOpenBracket)|(1<<ExprParserOpenParen)|(1<<ExprParserOpenBrace)|(1<<ExprParserDot)|(1<<ExprParserPlus)|(1<<ExprParserMinus)|(1<<ExprParserNegate)|(1<<ExprParserNot)|(1<<ExprParserNil))) != 0) || (((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(ExprParserPointer-44))|(1<<(ExprParserBooleanLiteral-44))|(1<<(ExprParserIntegerLiteral-44))|(1<<(ExprParserFloatLiteral-44))|(1<<(ExprParserHexIntegerLiteral-44))|(1<<(ExprParserIdentifier-44))|(1<<(ExprParserStringLiteral-44)))) != 0) {
					{
						p.SetState(94)

						var _x = p.expr(0)

						localctx.(*SliceContext).a = _x
					}

				}
				{
					p.SetState(97)
					p.Match(ExprParserColon)
				}
				p.SetState(99)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExprParserT__7)|(1<<ExprParserT__8)|(1<<ExprParserT__9)|(1<<ExprParserT__10)|(1<<ExprParserT__11)|(1<<ExprParserT__12)|(1<<ExprParserT__13)|(1<<ExprParserOpenBracket)|(1<<ExprParserOpenParen)|(1<<ExprParserOpenBrace)|(1<<ExprParserDot)|(1<<ExprParserPlus)|(1<<ExprParserMinus)|(1<<ExprParserNegate)|(1<<ExprParserNot)|(1<<ExprParserNil))) != 0) || (((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(ExprParserPointer-44))|(1<<(ExprParserBooleanLiteral-44))|(1<<(ExprParserIntegerLiteral-44))|(1<<(ExprParserFloatLiteral-44))|(1<<(ExprParserHexIntegerLiteral-44))|(1<<(ExprParserIdentifier-44))|(1<<(ExprParserStringLiteral-44)))) != 0) {
					{
						p.SetState(98)

						var _x = p.expr(0)

						localctx.(*SliceContext).b = _x
					}

				}
				{
					p.SetState(101)
					p.Match(ExprParserCloseBracket)
				}

			case 16:
				localctx = NewMemberDotContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(102)

				if !(p.Precpred(p.GetParserRuleContext(), 28)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 28)", ""))
				}
				{
					p.SetState(103)
					p.Match(ExprParserDot)
				}
				{
					p.SetState(104)

					var _m = p.Match(ExprParserIdentifier)

					localctx.(*MemberDotContext).name = _m
				}

			case 17:
				localctx = NewCallContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_expr)
				p.SetState(105)

				if !(p.Precpred(p.GetParserRuleContext(), 26)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 26)", ""))
				}
				{
					p.SetState(106)
					p.Match(ExprParserOpenParen)
				}
				p.SetState(108)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExprParserT__7)|(1<<ExprParserT__8)|(1<<ExprParserT__9)|(1<<ExprParserT__10)|(1<<ExprParserT__11)|(1<<ExprParserT__12)|(1<<ExprParserT__13)|(1<<ExprParserOpenBracket)|(1<<ExprParserOpenParen)|(1<<ExprParserOpenBrace)|(1<<ExprParserDot)|(1<<ExprParserPlus)|(1<<ExprParserMinus)|(1<<ExprParserNegate)|(1<<ExprParserNot)|(1<<ExprParserNil))) != 0) || (((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(ExprParserPointer-44))|(1<<(ExprParserBooleanLiteral-44))|(1<<(ExprParserIntegerLiteral-44))|(1<<(ExprParserFloatLiteral-44))|(1<<(ExprParserHexIntegerLiteral-44))|(1<<(ExprParserIdentifier-44))|(1<<(ExprParserStringLiteral-44)))) != 0) {
					{
						p.SetState(107)

						var _x = p.Arguments()

						localctx.(*CallContext).args = _x
					}

				}
				{
					p.SetState(110)
					p.Match(ExprParserCloseParen)
				}

			}

		}
		p.SetState(115)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext())
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

type BuiltinLenContext struct {
	*BuiltinsContext
	name antlr.Token
	e    IExprContext
}

func NewBuiltinLenContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BuiltinLenContext {
	var p = new(BuiltinLenContext)

	p.BuiltinsContext = NewEmptyBuiltinsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*BuiltinsContext))

	return p
}

func (s *BuiltinLenContext) GetName() antlr.Token { return s.name }

func (s *BuiltinLenContext) SetName(v antlr.Token) { s.name = v }

func (s *BuiltinLenContext) GetE() IExprContext { return s.e }

func (s *BuiltinLenContext) SetE(v IExprContext) { s.e = v }

func (s *BuiltinLenContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BuiltinLenContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *BuiltinLenContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *BuiltinLenContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BuiltinLenContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBuiltinLen(s)
	}
}

func (s *BuiltinLenContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBuiltinLen(s)
	}
}

func (s *BuiltinLenContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBuiltinLen(s)

	default:
		return t.VisitChildren(s)
	}
}

type BuiltinContext struct {
	*BuiltinsContext
	name antlr.Token
	e    IExprContext
	c    IClosureContext
}

func NewBuiltinContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BuiltinContext {
	var p = new(BuiltinContext)

	p.BuiltinsContext = NewEmptyBuiltinsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*BuiltinsContext))

	return p
}

func (s *BuiltinContext) GetName() antlr.Token { return s.name }

func (s *BuiltinContext) SetName(v antlr.Token) { s.name = v }

func (s *BuiltinContext) GetE() IExprContext { return s.e }

func (s *BuiltinContext) GetC() IClosureContext { return s.c }

func (s *BuiltinContext) SetE(v IExprContext) { s.e = v }

func (s *BuiltinContext) SetC(v IClosureContext) { s.c = v }

func (s *BuiltinContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BuiltinContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenParen, 0)
}

func (s *BuiltinContext) Comma() antlr.TerminalNode {
	return s.GetToken(ExprParserComma, 0)
}

func (s *BuiltinContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseParen, 0)
}

func (s *BuiltinContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BuiltinContext) Closure() IClosureContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IClosureContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IClosureContext)
}

func (s *BuiltinContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterBuiltin(s)
	}
}

func (s *BuiltinContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitBuiltin(s)
	}
}

func (s *BuiltinContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitBuiltin(s)

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

	p.SetState(163)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ExprParserT__7:
		localctx = NewBuiltinLenContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(116)

			var _m = p.Match(ExprParserT__7)

			localctx.(*BuiltinLenContext).name = _m
		}
		{
			p.SetState(117)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(118)

			var _x = p.expr(0)

			localctx.(*BuiltinLenContext).e = _x
		}
		{
			p.SetState(119)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__8:
		localctx = NewBuiltinContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(121)

			var _m = p.Match(ExprParserT__8)

			localctx.(*BuiltinContext).name = _m
		}
		{
			p.SetState(122)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(123)

			var _x = p.expr(0)

			localctx.(*BuiltinContext).e = _x
		}
		{
			p.SetState(124)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(125)

			var _x = p.Closure()

			localctx.(*BuiltinContext).c = _x
		}
		{
			p.SetState(126)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__9:
		localctx = NewBuiltinContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(128)

			var _m = p.Match(ExprParserT__9)

			localctx.(*BuiltinContext).name = _m
		}
		{
			p.SetState(129)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(130)

			var _x = p.expr(0)

			localctx.(*BuiltinContext).e = _x
		}
		{
			p.SetState(131)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(132)

			var _x = p.Closure()

			localctx.(*BuiltinContext).c = _x
		}
		{
			p.SetState(133)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__10:
		localctx = NewBuiltinContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(135)

			var _m = p.Match(ExprParserT__10)

			localctx.(*BuiltinContext).name = _m
		}
		{
			p.SetState(136)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(137)

			var _x = p.expr(0)

			localctx.(*BuiltinContext).e = _x
		}
		{
			p.SetState(138)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(139)

			var _x = p.Closure()

			localctx.(*BuiltinContext).c = _x
		}
		{
			p.SetState(140)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__11:
		localctx = NewBuiltinContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(142)

			var _m = p.Match(ExprParserT__11)

			localctx.(*BuiltinContext).name = _m
		}
		{
			p.SetState(143)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(144)

			var _x = p.expr(0)

			localctx.(*BuiltinContext).e = _x
		}
		{
			p.SetState(145)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(146)

			var _x = p.Closure()

			localctx.(*BuiltinContext).c = _x
		}
		{
			p.SetState(147)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__12:
		localctx = NewBuiltinContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(149)

			var _m = p.Match(ExprParserT__12)

			localctx.(*BuiltinContext).name = _m
		}
		{
			p.SetState(150)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(151)

			var _x = p.expr(0)

			localctx.(*BuiltinContext).e = _x
		}
		{
			p.SetState(152)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(153)

			var _x = p.Closure()

			localctx.(*BuiltinContext).c = _x
		}
		{
			p.SetState(154)
			p.Match(ExprParserCloseParen)
		}

	case ExprParserT__13:
		localctx = NewBuiltinContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(156)

			var _m = p.Match(ExprParserT__13)

			localctx.(*BuiltinContext).name = _m
		}
		{
			p.SetState(157)
			p.Match(ExprParserOpenParen)
		}
		{
			p.SetState(158)

			var _x = p.expr(0)

			localctx.(*BuiltinContext).e = _x
		}
		{
			p.SetState(159)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(160)

			var _x = p.Closure()

			localctx.(*BuiltinContext).c = _x
		}
		{
			p.SetState(161)
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

	// GetBody returns the body rule contexts.
	GetBody() IExprContext

	// SetBody sets the body rule contexts.
	SetBody(IExprContext)

	// IsClosureContext differentiates from other interfaces.
	IsClosureContext()
}

type ClosureContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	body   IExprContext
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

func (s *ClosureContext) GetBody() IExprContext { return s.body }

func (s *ClosureContext) SetBody(v IExprContext) { s.body = v }

func (s *ClosureContext) OpenBrace() antlr.TerminalNode {
	return s.GetToken(ExprParserOpenBrace, 0)
}

func (s *ClosureContext) CloseBrace() antlr.TerminalNode {
	return s.GetToken(ExprParserCloseBrace, 0)
}

func (s *ClosureContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ClosureContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClosureContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ClosureContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterClosure(s)
	}
}

func (s *ClosureContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitClosure(s)
	}
}

func (s *ClosureContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExprVisitor:
		return t.VisitClosure(s)

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

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(165)
		p.Match(ExprParserOpenBrace)
	}
	{
		p.SetState(166)

		var _x = p.expr(0)

		localctx.(*ClosureContext).body = _x
	}
	{
		p.SetState(167)
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
		p.SetState(169)

		var _x = p.expr(0)

		localctx.(*ArgumentsContext)._expr = _x
	}
	localctx.(*ArgumentsContext).list = append(localctx.(*ArgumentsContext).list, localctx.(*ArgumentsContext)._expr)
	p.SetState(174)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ExprParserComma {
		{
			p.SetState(170)
			p.Match(ExprParserComma)
		}
		{
			p.SetState(171)

			var _x = p.expr(0)

			localctx.(*ArgumentsContext)._expr = _x
		}
		localctx.(*ArgumentsContext).list = append(localctx.(*ArgumentsContext).list, localctx.(*ArgumentsContext)._expr)

		p.SetState(176)
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

	p.SetState(193)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(177)
			p.Match(ExprParserOpenBracket)
		}
		{
			p.SetState(178)
			p.Match(ExprParserCloseBracket)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(179)
			p.Match(ExprParserOpenBracket)
		}
		{
			p.SetState(180)

			var _x = p.expr(0)

			localctx.(*ArrayLiteralContext)._expr = _x
		}
		localctx.(*ArrayLiteralContext).list = append(localctx.(*ArrayLiteralContext).list, localctx.(*ArrayLiteralContext)._expr)
		p.SetState(185)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(181)
					p.Match(ExprParserComma)
				}
				{
					p.SetState(182)

					var _x = p.expr(0)

					localctx.(*ArrayLiteralContext)._expr = _x
				}
				localctx.(*ArrayLiteralContext).list = append(localctx.(*ArrayLiteralContext).list, localctx.(*ArrayLiteralContext)._expr)

			}
			p.SetState(187)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
		}
		p.SetState(189)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == ExprParserComma {
			{
				p.SetState(188)
				p.Match(ExprParserComma)
			}

		}
		{
			p.SetState(191)
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

	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(195)
			p.Match(ExprParserOpenBrace)
		}
		{
			p.SetState(196)
			p.Match(ExprParserCloseBrace)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(197)
			p.Match(ExprParserOpenBrace)
		}
		{
			p.SetState(198)

			var _x = p.PropertyNameAndValueList()

			localctx.(*MapLiteralContext).e = _x
		}
		p.SetState(200)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == ExprParserComma {
			{
				p.SetState(199)
				p.Match(ExprParserComma)
			}

		}
		{
			p.SetState(202)
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
		p.SetState(206)

		var _x = p.PropertyAssignment()

		localctx.(*PropertyNameAndValueListContext)._propertyAssignment = _x
	}
	localctx.(*PropertyNameAndValueListContext).list = append(localctx.(*PropertyNameAndValueListContext).list, localctx.(*PropertyNameAndValueListContext)._propertyAssignment)
	p.SetState(211)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(207)
				p.Match(ExprParserComma)
			}
			{
				p.SetState(208)

				var _x = p.PropertyAssignment()

				localctx.(*PropertyNameAndValueListContext)._propertyAssignment = _x
			}
			localctx.(*PropertyNameAndValueListContext).list = append(localctx.(*PropertyNameAndValueListContext).list, localctx.(*PropertyNameAndValueListContext)._propertyAssignment)

		}
		p.SetState(213)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())
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
		p.SetState(214)

		var _x = p.PropertyName()

		localctx.(*PropertyAssignmentContext).name = _x
	}
	{
		p.SetState(215)
		p.Match(ExprParserColon)
	}
	{
		p.SetState(216)

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
		p.SetState(218)
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
		return p.Precpred(p.GetParserRuleContext(), 24)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 23)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 22)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 21)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 20)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 19)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 18)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 17)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 16)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 15)

	case 10:
		return p.Precpred(p.GetParserRuleContext(), 14)

	case 11:
		return p.Precpred(p.GetParserRuleContext(), 13)

	case 12:
		return p.Precpred(p.GetParserRuleContext(), 12)

	case 13:
		return p.Precpred(p.GetParserRuleContext(), 30)

	case 14:
		return p.Precpred(p.GetParserRuleContext(), 29)

	case 15:
		return p.Precpred(p.GetParserRuleContext(), 28)

	case 16:
		return p.Precpred(p.GetParserRuleContext(), 26)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
