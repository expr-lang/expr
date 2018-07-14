package expr

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenKind int

const (
	name tokenKind = iota
	number
	text
	operator
	punctuation
	eof = -1
)

type token struct {
	kind  tokenKind
	value string
	pos   int
}

// is tests if token kind and value matches
func (t token) is(kind tokenKind, values ...string) bool {
	var value *string
	if len(values) == 1 {
		value = &values[0]
	}
	return t.kind == kind && (value == nil || *value == t.value)
}

func (t token) String() string {
	switch t.kind {
	case name:
		return fmt.Sprintf("name(%s)", t.value)
	case number:
		return fmt.Sprintf("number(%s)", t.value)
	case text:
		return fmt.Sprintf("text(%q)", t.value)
	case operator:
		return fmt.Sprintf("operator(%s)", t.value)
	case punctuation:
		return fmt.Sprintf("punctuation(%q)", t.value)
	case eof:
		return "EOF"
	default:
		return t.value
	}
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input    string  // the string being scanned
	pos      int     // current position in the input
	start    int     // start position of this token
	width    int     // width of last rune read from input
	brackets []rune  // stack of brackets
	tokens   []token // slice of scanned tokens
	err      error   // last error
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) word() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) emit(t tokenKind) {
	l.emitValue(t, l.word())
}

func (l *lexer) emitValue(t tokenKind, value string) {
	c := len(l.tokens)
	// Special case for joining "not" and ".." operators
	if c > 0 && l.tokens[c-1].is(operator, "not") && t == operator && value == "in" {
		l.tokens[c-1].value = "not in"
	} else if c > 0 && l.tokens[c-1].is(punctuation, ".") && t == punctuation && value == "." {
		l.tokens[c-1].kind = operator
		l.tokens[c-1].value = ".."
	} else {
		l.tokens = append(l.tokens, token{
			kind:  t,
			value: value,
			pos:   l.start,
		})
	}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.err = &syntaxError{
		message: fmt.Sprintf(format, args...),
		input:   l.input,
		pos:     l.start,
	}
	return nil
}

func lex(input string) ([]token, error) {
	l := &lexer{
		input:  input,
		tokens: make([]token, 0),
	}
	for state := lexRoot; state != nil; {
		state = state(l)
	}
	return l.tokens, l.err
}

// state functions

func lexRoot(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		if len(l.brackets) > 0 {
			return l.errorf("unclosed %q", string(l.brackets[0]))
		}
		l.emit(eof)
		return nil
	case isSpace(r):
		l.ignore()
		return lexRoot
	case r == '\'' || r == '"':
		l.backup()
		return lexQuote
	case '0' <= r && r <= '9':
		l.backup()
		return lexNumber
	case strings.ContainsRune("([{", r):
		l.emit(punctuation)
		l.brackets = append(l.brackets, r)
	case strings.ContainsRune(")]}", r):
		if len(l.brackets) > 0 {
			bracket := l.brackets[len(l.brackets)-1]
			l.brackets = l.brackets[:len(l.brackets)-1]
			if isBracketMatch(bracket, r) {
				l.emit(punctuation)
			} else {
				return l.errorf("unclosed %q", string(bracket))
			}
		} else {
			return l.errorf("unexpected %q", string(r))
		}
	case strings.ContainsRune(".,?:", r):
		l.emit(punctuation)
	case strings.ContainsRune("!%&*+-/<=>^|~", r):
		l.backup()
		return lexOperator
	case isAlphaNumeric(r):
		l.backup()
		return lexName
	default:
		return l.errorf("unrecognized character: %#U", r)
	}
	return lexRoot
}

func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.word())
	}
	l.emit(number)
	return lexRoot
}

func (l *lexer) scanNumber() bool {
	// Is it hex?
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		// Lookup for .. operator.
		if l.peek() == '.' {
			l.backup()
			return true
		}
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func lexQuote(l *lexer) stateFn {
	quote := l.next()
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof:
			return l.errorf("unterminated string")
		case quote:
			break Loop
		}
	}
	word := strings.Trim(l.word(), `"'`)
	value, err := strconv.Unquote(`"` + word + `"`)
	if err != nil {
		return l.errorf("unquote error: %v", err)
	}
	l.emitValue(text, value)
	return lexRoot
}

func lexOperator(l *lexer) stateFn {
	l.next()
	l.accept("|&=*")
	l.emit(operator)
	return lexRoot
}

func lexName(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			switch l.word() {
			case "not":
				l.emit(operator)
			case "in":
				l.emit(operator)
			case "or":
				l.emit(operator)
			case "and":
				l.emit(operator)
			case "matches":
				l.emit(operator)
			default:
				l.emit(name)
			}
			break Loop
		}
	}
	return lexRoot
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return isAlphabetic(r) || unicode.IsDigit(r)
}

func isAlphabetic(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isBracketMatch(open, close rune) bool {
	switch string([]rune{open, close}) {
	case "()":
		return true
	case "[]":
		return true
	case "{}":
		return true
	}
	return false
}
