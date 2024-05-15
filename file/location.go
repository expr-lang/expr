package file

type Location struct {
	Line   int `json:"line"`   // The 1-based line of the location.
	Column int `json:"column"` // The 0-based column number of the location.
	From   int `json:"from"`   // The 0-based byte offset from the beginning.
	To     int `json:"to"`     // The 0-based byte offset to the end.
}

func (l Location) Empty() bool {
	return l.Column == 0 && l.Line == 0
}
