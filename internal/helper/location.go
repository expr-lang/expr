package helper

type Location struct {
	line   int
	column int
}

func NewLocation(line, column int) Location {
	return Location{
		line:   line,
		column: column,
	}
}

// Line returns the 1-based line of the location.
func (l Location) Line() int {
	return l.line
}

// Column returns the 0-based column number of the location.
func (l Location) Column() int {
	return l.column
}

func (l Location) Empty() bool {
	return l.column == 0 && l.line == 0
}
