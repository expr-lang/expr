package lexer

import (
	"strings"
)

type stateFn func(*lexer) stateFn

func root(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(EOF)
		return nil
	case isSpace(r):
		l.ignore()
		return root
	case r == '\'' || r == '"':
		l.scanString(r)
		str, err := unescape(l.word())
		if err != nil {
			l.error("%v", err)
		}
		l.emitValue(String, str)
	case '0' <= r && r <= '9':
		l.backup()
		return lexNumber
	case strings.ContainsRune("([{", r):
		l.emit(Bracket)
	case strings.ContainsRune(")]}", r):
		l.emit(Bracket)
	case strings.ContainsRune(".,?:", r):
		l.emit(Punctuation)
	case strings.ContainsRune("!%&*+-/<=>^|~", r):
		l.backup()
		return lexOperator
	case isAlphaNumeric(r):
		l.backup()
		return lexName
	default:
		return l.error("unrecognized character: %#U", r)
	}
	return root
}

func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.error("bad number syntax: %q", l.word())
	}
	l.emit(Number)
	return root
}

func (l *lexer) scanNumber() bool {
	// Is it hex?
	digits := "0123456789_"
	l.acceptRun(digits)
	if l.accept(".") {
		// Lookup for .. operator: if after dot there is another dot (1..2), it maybe a range operator.
		if l.peek() == '.' {
			l.backup()
			return true
		}
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func lexOperator(l *lexer) stateFn {
	l.next()
	l.accept("|&=*")
	l.emit(Operator)
	return root
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
				l.emit(Operator)
			case "in":
				l.emit(Operator)
			case "or":
				l.emit(Operator)
			case "and":
				l.emit(Operator)
			case "matches":
				l.emit(Operator)
			default:
				l.emit(Identifier)
			}
			break Loop
		}
	}
	return root
}
