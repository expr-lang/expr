package file

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Error struct {
	Location
	Message string
}

const (
	dot = "."
	ind = "^"
)

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Format(source *Source) string {
	if e.Location.Empty() {
		return e.Message
	}
	var result = fmt.Sprintf(
		"%s (%d:%d)",
		e.Message,
		e.Location.Line,
		e.Location.Column+1, // add one to the 0-based column for display
	)
	if snippet, found := source.Snippet(e.Location.Line); found {
		snippet := strings.Replace(snippet, "\t", " ", -1)
		srcLine := "\n | " + snippet
		var bytes = []byte(snippet)
		var indLine = "\n | "
		for i := 0; i < e.Location.Column && len(bytes) > 0; i++ {
			_, sz := utf8.DecodeRune(bytes)
			bytes = bytes[sz:]
			if sz > 1 {
				goto noind
			} else {
				indLine += dot
			}
		}
		if _, sz := utf8.DecodeRune(bytes); sz > 1 {
			goto noind
		} else {
			indLine += ind
		}
		srcLine += indLine

	noind:
		result += srcLine
	}
	return result
}
