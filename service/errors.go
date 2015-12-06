package service

import "bytes"

// AggregateError is an error implementation that aggregates a list of errors into a single error.
type AggregateError struct {
	errors []error
}

// Creates a new, empty AggregateError instance.
func NewAggregateError() AggregateError {
	return AggregateError{
		errors: []error{},
	}
}

// Add an error to the list of errors to aggregate.
func (e *AggregateError) AddError(err error) {
	e.errors = append(e.errors, err)
}

// Returns whether the AggregateError is currently empty.
func (e AggregateError) Empty() bool {
	return len(e.errors) == 0
}

// Returns an aggregate error message of all of contained errors. The aggregate message is just the message of each
// contained error, separated by a newline character.
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
