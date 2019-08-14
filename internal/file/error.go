package file

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/width"
)

// Error type which references a location within source and a message.
type Error struct {
	location Location
	message  string
	source   *Source
}

func NewError(message string, location Location, source *Source) Error {
	return Error{
		message:  message,
		location: location,
		source:   source,
	}
}

const (
	dot = "."
	ind = "^"
)

var (
	wideDot = width.Widen.String(dot)
	wideInd = width.Widen.String(ind)
)

func (e Error) Message() string {
	return e.message
}

func (e Error) Location() (line, column int) {
	return e.location.Line, e.location.Column
}

func (e Error) Source() string {
	if e.source == nil {
		return ""
	}
	return e.source.Content()
}

func (e Error) Error() string {
	if e.location.Empty() {
		return e.message
	}
	var result = fmt.Sprintf(
		"%s (%d:%d)",
		e.message,
		e.location.Line,
		e.location.Column+1, // add one to the 0-based column for display
	)
	if e.source != nil {
		if snippet, found := e.source.Snippet(e.location.Line); found {
			snippet := strings.Replace(snippet, "\t", " ", -1)
			srcLine := "\n | " + snippet
			var bytes = []byte(snippet)
			var indLine = "\n | "
			for i := 0; i < e.location.Column && len(bytes) > 0; i++ {
				_, sz := utf8.DecodeRune(bytes)
				bytes = bytes[sz:]
				if sz > 1 {
					indLine += wideDot
				} else {
					indLine += dot
				}
			}
			if _, sz := utf8.DecodeRune(bytes); sz > 1 {
				indLine += wideInd
			} else {
				indLine += ind
			}
			result += srcLine + indLine
		}
	}
	return result
}
