package lexer

import (
	"fmt"
	"github.com/antonmedv/expr/file"
)

type Kind string

const (
	Identifier  Kind = "Identifier"
	Number           = "Number"
	String           = "String"
	Operator         = "Operator"
	Bracket          = "Bracket"
	Punctuation      = "Punctuation"
	EOF              = "EOF"
)

type Token struct {
	file.Location
	Kind  Kind
	Value string
}

func (t Token) String() string {
	if t.Value == "" {
		return string(t.Kind)
	}
	return fmt.Sprintf("%s(%#v)", t.Kind, t.Value)
}
