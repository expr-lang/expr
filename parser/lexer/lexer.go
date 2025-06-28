package lexer

import (
	"fmt"
	"strings"

	"github.com/expr-lang/expr/file"
)

func Lex(source file.Source) ([]Token, error) {
	raw := source.String()
	l := &lexer{
		raw:    raw,
		runes:  []rune(raw),
		tokens: make([]Token, 0),
	}

	for state := root; state != nil; {
		state = state(l)
	}

	if l.err != nil {
		return nil, l.err.Bind(source)
	}

	return l.tokens, nil
}

type lexer struct {
	raw        string
	runes      []rune
	tokens     []Token
	err        *file.Error
	start, end pos
	eof        bool
}

type pos struct {
	byte, rune int
}

const eof rune = -1

func (l *lexer) commit() {
	l.start = l.end
}

func (l *lexer) next() rune {
	if l.end.rune >= len(l.runes) {
		l.eof = true
		return eof
	}
	r := l.runes[l.end.rune]
	l.end.rune++
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	if l.eof {
		l.eof = false
	} else {
		l.end.rune--
	}
}

func (l *lexer) emit(t Kind) {
	l.emitValue(t, l.word())
}

func (l *lexer) emitValue(t Kind, value string) {
	l.tokens = append(l.tokens, Token{
		Location: file.Location{From: l.start.rune, To: l.end.rune},
		Kind:     t,
		Value:    value,
	})
	l.commit()
}

func (l *lexer) emitEOF() {
	from := l.end.rune - 1
	if from < 0 {
		from = 0
	}
	to := l.end.rune - 0
	if to < 0 {
		to = 0
	}
	l.tokens = append(l.tokens, Token{
		Location: file.Location{From: from, To: to},
		Kind:     EOF,
	})
	l.commit()
}

func (l *lexer) skip() {
	l.commit()
}

func (l *lexer) word() string {
	// TODO: boundary check is NOT needed here, but for some reason CI fuzz tests are failing.
	if l.start.rune > len(l.runes) || l.end.rune > len(l.runes) {
		return "__invalid__"
	}
	return string(l.runes[l.start.rune:l.end.rune])
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

func (l *lexer) skipSpaces() {
	r := l.peek()
	for ; r == ' '; r = l.peek() {
		l.next()
	}
	l.skip()
}

func (l *lexer) acceptWord(word string) bool {
	pos := l.end

	l.skipSpaces()

	for _, ch := range word {
		if l.next() != ch {
			l.end = pos
			return false
		}
	}
	if r := l.peek(); r != ' ' && r != eof {
		l.end = pos
		return false
	}

	return true
}

func (l *lexer) error(format string, args ...any) stateFn {
	if l.err == nil { // show first error
		end := l.end.rune
		if l.eof {
			end++
		}
		l.err = &file.Error{
			Location: file.Location{
				From: end - 1,
				To:   end,
			},
			Message: fmt.Sprintf(format, args...),
		}
	}
	return nil
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16 // larger than any legal digit val
}

func lower(ch rune) rune { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter

func (l *lexer) scanDigits(ch rune, base, n int) rune {
	for n > 0 && digitVal(ch) < base {
		ch = l.next()
		n--
	}
	if n > 0 {
		l.error("invalid char escape")
	}
	return ch
}

func (l *lexer) scanEscape(quote rune) rune {
	ch := l.next() // read character after '/'
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		// nothing to do
		ch = l.next()
	case '0', '1', '2', '3', '4', '5', '6', '7':
		ch = l.scanDigits(ch, 8, 3)
	case 'x':
		ch = l.scanDigits(l.next(), 16, 2)
	case 'u':
		ch = l.scanDigits(l.next(), 16, 4)
	case 'U':
		ch = l.scanDigits(l.next(), 16, 8)
	default:
		l.error("invalid char escape")
	}
	return ch
}

func (l *lexer) scanString(quote rune) (n int) {
	ch := l.next() // read character after quote
	for ch != quote {
		if ch == '\n' || ch == eof {
			l.error("literal not terminated")
			return
		}
		if ch == '\\' {
			ch = l.scanEscape(quote)
		} else {
			ch = l.next()
		}
		n++
	}
	return
}

func (l *lexer) scanRawString(quote rune) (n int) {
	ch := l.next() // read character after back tick
	for ch != quote {
		if ch == eof {
			l.error("literal not terminated")
			return
		}
		ch = l.next()
		n++
	}
	l.emitValue(String, string(l.runes[l.start.rune+1:l.end.rune-1]))
	return
}
