package file

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

type Source struct {
	contents []rune
}

func NewSource(contents string) *Source {
	return &Source{
		contents: []rune(contents),
	}
}

func (s *Source) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.contents)
}

func (s *Source) UnmarshalJSON(b []byte) error {
	contents := make([]rune, 0)
	err := json.Unmarshal(b, &contents)
	if err != nil {
		return err
	}

	s.contents = contents
	return nil
}

func (s *Source) String() string {
	return string(s.contents)
}

func (s *Source) Snippet(line int) (string, bool) {
	if s == nil {
		return "", false
	}
	lines := strings.Split(string(s.contents), "\n")
	lineOffsets := make([]int, len(lines))
	var offset int
	for i, line := range lines {
		offset = offset + utf8.RuneCountInString(line) + 1
		lineOffsets[i] = offset
	}
	charStart, found := getLineOffset(lineOffsets, line)
	if !found || len(s.contents) == 0 {
		return "", false
	}
	charEnd, found := getLineOffset(lineOffsets, line+1)
	if found {
		return string(s.contents[charStart : charEnd-1]), true
	}
	return string(s.contents[charStart:]), true
}

func getLineOffset(lineOffsets []int, line int) (int, bool) {
	if line == 1 {
		return 0, true
	} else if line > 1 && line <= len(lineOffsets) {
		offset := lineOffsets[line-2]
		return offset, true
	}
	return -1, false
}
