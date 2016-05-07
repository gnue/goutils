package merr

import (
	"fmt"
)

type MultiError struct {
	Errors []error
}

func New(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return &MultiError{errs}
}

func (e *MultiError) Error() string {
	errs := e.Errors
	tail := "..."

	switch len(errs) {
	case 0:
		return ""
	case 1:
		tail = ""
	}

	return fmt.Sprint("errors: ", len(errs), errs[0], tail)
}
