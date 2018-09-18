package expr

import (
	"strings"
	"testing"
)

type lexTest struct {
	input  string
	tokens []token
}

type lexErrorTest struct {
	input string
	err   string
}

var lexTests = []lexTest{
	{
		"1 02 1e3 1.2e-4",
		[]token{
			{kind: number, value: "1"},
			{kind: number, value: "02"},
			{kind: number, value: "1e3"},
			{kind: number, value: "1.2e-4"},
			{kind: eof},
		},
	},
	{
		`"double" 'single' "abc \n\t\" "`,
		[]token{
			{kind: text, value: "double"},
			{kind: text, value: "single"},
			{kind: text, value: "abc \n\t\" "},
			{kind: eof},
		},
	},
	{
		"+0 != -0",
		[]token{
			{kind: operator, value: "+"},
			{kind: number, value: "0"},
			{kind: operator, value: "!="},
			{kind: operator, value: "-"},
			{kind: number, value: "0"},
			{kind: eof},
		},
	},
	{
		"a and b or not in c not orx",
		[]token{
			{kind: name, value: "a"},
			{kind: operator, value: "and"},
			{kind: name, value: "b"},
			{kind: operator, value: "or"},
			{kind: operator, value: "not in"},
			{kind: name, value: "c"},
			{kind: operator, value: "not"},
			{kind: name, value: "orx"},
			{kind: eof},
		},
	},
	{
		`(3 + 5) ~ foo("bar").baz[4]`,
		[]token{
			{kind: punctuation, value: "("},
			{kind: number, value: "3"},
			{kind: operator, value: "+"},
			{kind: number, value: "5"},
			{kind: punctuation, value: ")"},
			{kind: operator, value: "~"},
			{kind: name, value: "foo"},
			{kind: punctuation, value: "("},
			{kind: text, value: "bar"},
			{kind: punctuation, value: ")"},
			{kind: punctuation, value: "."},
			{kind: name, value: "baz"},
			{kind: punctuation, value: "["},
			{kind: number, value: "4"},
			{kind: punctuation, value: "]"},
			{kind: eof},
		},
	},
	{
		`1..2`,
		[]token{
			{kind: number, value: "1"},
			{kind: operator, value: ".."},
			{kind: number, value: "2"},
			{kind: eof},
		},
	},
	{
		`matches`,
		[]token{
			{kind: operator, value: "matches"},
			{kind: eof},
		},
	},
	{
		`'\.'`,
		[]token{
			{kind: text, value: "\\."},
			{kind: eof},
		},
	},
}

var lexErrorTests = []lexErrorTest{
	{
		`{[}]`,
		`unclosed "["`,
	},
}

func compareTokens(i1, i2 []token) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].kind != i2[k].kind {
			return false
		}
		if i1[k].value != i2[k].value {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		tokens, err := lex(test.input)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		if !compareTokens(tokens, test.tokens) {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, tokens, test.tokens)
		}
	}
}

func TestLex_error(t *testing.T) {
	for _, test := range lexErrorTests {
		_, err := lex(test.input)
		if !strings.HasPrefix(err.Error(), test.err) || test.err == "" {
			t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%v", test.input, err.Error(), test.err)
		}
	}
}
