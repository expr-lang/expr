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
	err := NewError(fmt.Sprintf(format, args...), l, e.source)
	e.errors = append(e.errors, err)
}

func (e *Errors) GetErrors() []Error {
	return e.errors[:]
}

func (e *Errors) HasError() bool {
	return len(e.errors) > 0
}

func (e *Errors) First() error {
	return e.errors[0]
}

func (e *Errors) Error() string {
	var result = ""
	for i, err := range e.errors {
		if i >= 1 {
			result += "\n"
		}
		result += err.Error()
	}
	return result
}
