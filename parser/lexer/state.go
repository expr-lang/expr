package lexer

import (
	"strings"
)

type stateFn func(*lexer) stateFn

func root(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emitEOF()
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
		return number
	case strings.ContainsRune("([{", r):
		l.emit(Bracket)
	case strings.ContainsRune(")]}", r):
		l.emit(Bracket)
	case strings.ContainsRune("#,?:%+-/", r): // single rune operator
		l.emit(Operator)
	case strings.ContainsRune("&|!=*<>", r): // possible double rune operator
		l.accept("&|=*")
		l.emit(Operator)
	case r == '.':
		l.backup()
		return dot
	case isAlphaNumeric(r):
		l.backup()
		return identifier
	default:
		return l.error("unrecognized character: %#U", r)
	}
	return root
}

func number(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.error("bad number syntax: %q", l.word())
	}
	l.emit(Number)
	return root
}

func (l *lexer) scanNumber() bool {
	digits := "0123456789_"
	// Is it hex?
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}
	l.acceptRun(digits)
	if l.accept(".") {
		// Lookup for .. operator: if after dot there is another dot (1..2), it maybe a range operator.
		if l.peek() == '.' {
			l.backup()
			return true
		}
		l.accept(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun(digits)
	}
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func dot(l *lexer) stateFn {
	l.next()
	if l.accept("0123456789") {
		l.backup()
		return number
	}
	l.accept(".")
	l.emit(Operator)
	return root
}

func identifier(l *lexer) stateFn {
loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb
		default:
			l.backup()
			switch l.word() {
			case "not":
				return not
			case "in", "or", "and", "matches", "contains", "startsWith", "endsWith":
				l.emit(Operator)
			default:
				l.emit(Identifier)
			}
			break loop
		}
	}
	return root
}

func not(l *lexer) stateFn {
	if l.acceptWord(" in") {
		l.emit(Operator)
	} else {
		l.emit(Operator)
	}
	return root
}
