package helper

import (
	"strings"
	"unicode/utf8"
)

type Source struct {
	contents    []rune
	lineOffsets []int32
}

func NewSource(contents string) *Source {
	// Compute line offsets up front as they are referred to frequently.
	lines := strings.Split(contents, "\n")
	offsets := make([]int32, len(lines))
	var offset int32
	for i, line := range lines {
		offset = offset + int32(utf8.RuneCountInString(line)) + 1
		offsets[int32(i)] = offset
	}
	return &Source{
		contents:    []rune(contents),
		lineOffsets: offsets,
	}
}

func (s *Source) Content() string {
	return string(s.contents)
}

func (s *Source) Snippet(line int) (string, bool) {
	charStart, found := s.findLineOffset(line)
	if !found || len(s.contents) == 0 {
		return "", false
	}
	charEnd, found := s.findLineOffset(line + 1)
	if found {
		return string(s.contents[charStart : charEnd-1]), true
	}
	return string(s.contents[charStart:]), true
}

// findLineOffset returns the offset where the (1-indexed) line begins,
// or false if line doesn't exist.
func (s *Source) findLineOffset(line int) (int32, bool) {
	if line == 1 {
		return 0, true
	} else if line > 1 && line <= int(len(s.lineOffsets)) {
		offset := s.lineOffsets[line-2]
		return offset, true
	}
	return -1, false
}

// findLine finds the line that contains the given character offset and
// returns the line number and offset of the beginning of that line.
// Note that the last line is treated as if it contains all offsets
// beyond the end of the actual source.
func (s *Source) findLine(characterOffset int32) (int32, int32) {
	var line int32 = 1
	for _, lineOffset := range s.lineOffsets {
		if lineOffset > characterOffset {
			break
		} else {
			line++
		}
	}
	if line == 1 {
		return line, 0
	}
	return line, s.lineOffsets[line-2]
}
