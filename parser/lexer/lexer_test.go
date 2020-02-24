package lexer_test

import (
	"github.com/antonmedv/expr/file"
	"strings"
	"testing"

	. "github.com/antonmedv/expr/parser/lexer"
)

type lexTest struct {
	input  string
	tokens []Token
}

type lexErrorTest struct {
	input string
	err   string
}

var lexTests = []lexTest{
	{
		".5 1 02 1e3 0xFF 1.2e-4 1_000_000 _42 -.5",
		[]Token{
			{Kind: Number, Value: ".5"},
			{Kind: Number, Value: "1"},
			{Kind: Number, Value: "02"},
			{Kind: Number, Value: "1e3"},
			{Kind: Number, Value: "0xFF"},
			{Kind: Number, Value: "1.2e-4"},
			{Kind: Number, Value: "1_000_000"},
			{Kind: Identifier, Value: "_42"},
			{Kind: Operator, Value: "-"},
			{Kind: Number, Value: ".5"},
			{Kind: EOF},
		},
	},
	{
		`"double" 'single' "abc \n\t\"\\" '"\'' "'\"" "\xC3\xBF\u263A\U000003A8" '❤️'`,
		[]Token{
			{Kind: String, Value: "double"},
			{Kind: String, Value: "single"},
			{Kind: String, Value: "abc \n\t\"\\"},
			{Kind: String, Value: "\"'"},
			{Kind: String, Value: "'\""},
			{Kind: String, Value: "Ã¿☺Ψ"},
			{Kind: String, Value: "❤️"},
			{Kind: EOF},
		},
	},
	{
		"a and orb().val #.",
		[]Token{
			{Kind: Identifier, Value: "a"},
			{Kind: Operator, Value: "and"},
			{Kind: Identifier, Value: "orb"},
			{Kind: Bracket, Value: "("},
			{Kind: Bracket, Value: ")"},
			{Kind: Operator, Value: "."},
			{Kind: Identifier, Value: "val"},
			{Kind: Operator, Value: "#"},
			{Kind: Operator, Value: "."},
			{Kind: EOF},
		},
	},
	{
		`not in not abc not i not(false) not  in`,
		[]Token{
			{Kind: Operator, Value: "not in"},
			{Kind: Operator, Value: "not"},
			{Kind: Identifier, Value: "abc"},
			{Kind: Operator, Value: "not"},
			{Kind: Identifier, Value: "i"},
			{Kind: Operator, Value: "not"},
			{Kind: Bracket, Value: "("},
			{Kind: Identifier, Value: "false"},
			{Kind: Bracket, Value: ")"},
			{Kind: Operator, Value: "not"},
			{Kind: Operator, Value: "in"},
			{Kind: EOF},
		},
	},
	{
		`1..5`,
		[]Token{
			{Kind: Number, Value: "1"},
			{Kind: Operator, Value: ".."},
			{Kind: Number, Value: "5"},
			{Kind: EOF},
		},
	},
}

var lexErrorTests = []lexErrorTest{
	{
		`
		"\xQA"
		`,
		`invalid char escape (2:6)`,
	},
	{
		`id "hello`,
		`literal not terminated (1:9)`,
	},
}

func compareTokens(i1, i2 []Token) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].Kind != i2[k].Kind {
			return false
		}
		if i1[k].Value != i2[k].Value {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		tokens, err := Lex(file.NewSource(test.input))
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			return
		}
		if !compareTokens(tokens, test.tokens) {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, tokens, test.tokens)
		}
	}
}

func TestLex_error(t *testing.T) {
	for _, test := range lexErrorTests {
		out, err := Lex(file.NewSource(test.input))
		if err == nil {
			t.Errorf("%s:\nexpected error\n%v", test.input, out)
			continue
		}
		if !strings.HasPrefix(err.Error(), test.err) || test.err == "" {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, err.Error(), test.err)
		}
	}
}
