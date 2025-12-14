package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	. "github.com/expr-lang/expr/parser/lexer"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		tokens []Token
	}{
		{
			"1",
			[]Token{
				{Kind: Number, Value: "1"},
				{Kind: EOF},
			},
		},
		{
			".5 0.025 1 02 1e3 0xFF 0b0101 0o600 1.2e-4 1_000_000 _42 -.5",
			[]Token{
				{Kind: Number, Value: ".5"},
				{Kind: Number, Value: "0.025"},
				{Kind: Number, Value: "1"},
				{Kind: Number, Value: "02"},
				{Kind: Number, Value: "1e3"},
				{Kind: Number, Value: "0xFF"},
				{Kind: Number, Value: "0b0101"},
				{Kind: Number, Value: "0o600"},
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
			"`backtick` `hello\u263Aworld` `hello\n\tworld` `hello\"world'`  `\xC3\xBF\u263A\U000003A8` `❤️`",
			[]Token{
				{Kind: String, Value: `backtick`},
				{Kind: String, Value: `hello☺world`},
				{Kind: String, Value: `hello
	world`},
				{Kind: String, Value: `hello"world'`},
				{Kind: String, Value: `ÿ☺Ψ`},
				{Kind: String, Value: "❤️"},
				{Kind: EOF},
			},
		},
		{
			"`escaped backticks` `` `a``b` ```` `a``` ```b` ```a````b``` ```````` ```a````` `````b```",
			[]Token{
				{Kind: String, Value: "escaped backticks"},
				{Kind: String, Value: ""},
				{Kind: String, Value: "a`b"},
				{Kind: String, Value: "`"},
				{Kind: String, Value: "a`"},
				{Kind: String, Value: "`b"},
				{Kind: String, Value: "`a``b`"},
				{Kind: String, Value: "```"},
				{Kind: String, Value: "`a``"},
				{Kind: String, Value: "``b`"},
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
			"foo?.bar",
			[]Token{

				{Kind: Identifier, Value: "foo"},
				{Kind: Operator, Value: "?."},
				{Kind: Identifier, Value: "bar"},
				{Kind: EOF},
			},
		},
		{
			"foo ? .bar : .baz",
			[]Token{

				{Kind: Identifier, Value: "foo"},
				{Kind: Operator, Value: "?"},
				{Kind: Operator, Value: "."},
				{Kind: Identifier, Value: "bar"},
				{Kind: Operator, Value: ":"},
				{Kind: Operator, Value: "."},
				{Kind: Identifier, Value: "baz"},
				{Kind: EOF},
			},
		},
		{
			"func?()",
			[]Token{

				{Kind: Identifier, Value: "func"},
				{Kind: Operator, Value: "?"},
				{Kind: Bracket, Value: "("},
				{Kind: Bracket, Value: ")"},
				{Kind: EOF},
			},
		},
		{
			"array?[]",
			[]Token{

				{Kind: Identifier, Value: "array"},
				{Kind: Operator, Value: "?"},
				{Kind: Bracket, Value: "["},
				{Kind: Bracket, Value: "]"},
				{Kind: EOF},
			},
		},
		{
			`not in not abc not i not(false) not  in not   in`,
			[]Token{
				{Kind: Operator, Value: "not"},
				{Kind: Operator, Value: "in"},
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
				{Kind: Operator, Value: "not"},
				{Kind: Operator, Value: "in"},
				{Kind: EOF},
			},
		},
		{
			"not in_var",
			[]Token{
				{Kind: Operator, Value: "not"},
				{Kind: Identifier, Value: "in_var"},
				{Kind: EOF},
			},
		},
		{
			"not in",
			[]Token{
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
		{
			`$i _0 früh`,
			[]Token{
				{Kind: Identifier, Value: "$i"},
				{Kind: Identifier, Value: "_0"},
				{Kind: Identifier, Value: "früh"},
				{Kind: EOF},
			},
		},
		{
			`foo // comment
		bar // comment`,
			[]Token{
				{Kind: Identifier, Value: "foo"},
				{Kind: Identifier, Value: "bar"},
				{Kind: EOF},
			},
		},
		{
			`foo /* comment */ bar`,
			[]Token{
				{Kind: Identifier, Value: "foo"},
				{Kind: Identifier, Value: "bar"},
				{Kind: EOF},
			},
		},
		{
			`foo ?? bar`,
			[]Token{
				{Kind: Identifier, Value: "foo"},
				{Kind: Operator, Value: "??"},
				{Kind: Identifier, Value: "bar"},
				{Kind: EOF},
			},
		},
		{
			`let foo = bar;`,
			[]Token{
				{Kind: Operator, Value: "let"},
				{Kind: Identifier, Value: "foo"},
				{Kind: Operator, Value: "="},
				{Kind: Identifier, Value: "bar"},
				{Kind: Operator, Value: ";"},
				{Kind: EOF},
			},
		},
		{
			`#index #1 #`,
			[]Token{
				{Kind: Operator, Value: "#"},
				{Kind: Identifier, Value: "index"},
				{Kind: Operator, Value: "#"},
				{Kind: Identifier, Value: "1"},
				{Kind: Operator, Value: "#"},
				{Kind: EOF},
			},
		},
		{
			`: ::`,
			[]Token{
				{Kind: Operator, Value: ":"},
				{Kind: Operator, Value: "::"},
				{Kind: EOF},
			},
		},
		{
			`if a>b {x1+x2} else {x2}`,
			[]Token{
				{Kind: Operator, Value: "if"},
				{Kind: Identifier, Value: "a"},
				{Kind: Operator, Value: ">"},
				{Kind: Identifier, Value: "b"},
				{Kind: Bracket, Value: "{"},
				{Kind: Identifier, Value: "x1"},
				{Kind: Operator, Value: "+"},
				{Kind: Identifier, Value: "x2"},
				{Kind: Bracket, Value: "}"},
				{Kind: Operator, Value: "else"},
				{Kind: Bracket, Value: "{"},
				{Kind: Identifier, Value: "x2"},
				{Kind: Bracket, Value: "}"},
				{Kind: EOF},
			},
		},
		{
			`a>b if {x1} else {x2}`,
			[]Token{
				{Kind: Identifier, Value: "a"},
				{Kind: Operator, Value: ">"},
				{Kind: Identifier, Value: "b"},
				{Kind: Operator, Value: "if"},
				{Kind: Bracket, Value: "{"},
				{Kind: Identifier, Value: "x1"},
				{Kind: Bracket, Value: "}"},
				{Kind: Operator, Value: "else"},
				{Kind: Bracket, Value: "{"},
				{Kind: Identifier, Value: "x2"},
				{Kind: Bracket, Value: "}"},
				{Kind: EOF},
			},
		},
		{
			"\"\\u{61}\\u{1F600}\" '\\u{61}\\u{1F600}'",
			[]Token{
				{Kind: String, Value: "a😀"},
				{Kind: String, Value: "a😀"},
				{Kind: EOF},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := Lex(file.NewSource(test.input))
			if err != nil {
				t.Errorf("%s:\n%v", test.input, err)
				return
			}
			if !compareTokens(tokens, test.tokens) {
				t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, tokens, test.tokens)
			}
		})
	}
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

func TestLex_location(t *testing.T) {
	source := file.NewSource("1..2\n3..4")
	tokens, err := Lex(source)
	require.NoError(t, err)
	require.Equal(t, []Token{
		{Location: file.Location{From: 0, To: 1}, Kind: Number, Value: "1"},
		{Location: file.Location{From: 1, To: 3}, Kind: Operator, Value: ".."},
		{Location: file.Location{From: 3, To: 4}, Kind: Number, Value: "2"},
		{Location: file.Location{From: 5, To: 6}, Kind: Number, Value: "3"},
		{Location: file.Location{From: 6, To: 8}, Kind: Operator, Value: ".."},
		{Location: file.Location{From: 8, To: 9}, Kind: Number, Value: "4"},
		{Location: file.Location{From: 8, To: 9}, Kind: EOF, Value: ""},
	}, tokens)
}

const errorTests = `
"\xQA"
invalid char escape (1:4)
 | "\xQA"
 | ...^

id "hello
literal not terminated (1:10)
 | id "hello
 | .........^

id ` + "`" + `hello
literal not terminated (1:10)
 | id ` + "`" + `hello
 | .........^

id ` + "`" + `hello` + "``" + `
literal not terminated (1:12)
 | id ` + "`" + `hello` + "``" + `
 | ...........^

id ` + "```" + `hello
literal not terminated (1:12)
 | id ` + "```" + `hello
 | ...........^

id ` + "`" + `hello` + "``" + ` world
literal not terminated (1:18)
 | id ` + "`" + `hello` + "``" + ` world
 | .................^

früh ♥︎
unrecognized character: U+2665 '♥' (1:6)
 | früh ♥︎
 | .....^
`

func TestLex_error(t *testing.T) {
	tests := strings.Split(strings.Trim(errorTests, "\n"), "\n\n")

	for _, test := range tests {
		input := strings.SplitN(test, "\n", 2)
		if len(input) != 2 {
			t.Errorf("syntax error in test: %q", test)
			break
		}

		_, err := Lex(file.NewSource(input[0]))
		if err == nil {
			err = fmt.Errorf("<nil>")
		}

		assert.Equal(t, input[1], err.Error(), input[0])
	}
}
