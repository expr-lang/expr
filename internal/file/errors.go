package file

import (
	"fmt"
)

type Errors struct {
	errors []Error
	source *Source
}

func NewErrors(source *Source) *Errors {
	return &Errors{
		errors: []Error{},
		source: source,
	}
}

func (e *Errors) ReportError(l Location, format string, args ...interface{}) {
	err := Error{
		Location: l,
		Message:  fmt.Sprintf(format, args...),
	}
	e.errors = append(e.errors, err)
}

func (e *Errors) GetErrors() []Error {
	return e.errors[:]
}

func (e *Errors) HasError() bool {
	return len(e.errors) > 0
}

func (e *Errors) First() string {
	return e.errors[0].Format(e.source)
}

func (e *Errors) Error() string {
	var result = ""
	for i, err := range e.errors {
		if i >= 1 {
			result += "\n"
		}
		result += err.Format(e.source)
	}
	return result
}
