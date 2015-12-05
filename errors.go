package kaiju

import "bytes"

type AggregateError struct {
	errors []error
}

func NewAggregateError() AggregateError {
	return AggregateError{
		errors: []error{},
	}
}

func (e *AggregateError) AddError(err error) {
	e.errors = append(e.errors, err)
}

func (e AggregateError) Empty() bool {
	return len(e.errors) == 0
}

func (e AggregateError) Error() string {
	var buffer bytes.Buffer
	for _, err := range e.errors {
		buffer.WriteString(err.Error())
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func (e AggregateError) String() string {
	return e.Error()
}
