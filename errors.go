// multiple error
package merr

import (
	"fmt"
)

type Errors interface {
	Errors() []error
}

type multiError struct {
	errs []error
}

func New(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return &multiError{errs}
}

func (e *multiError) Error() string {
	errs := e.errs
	tail := "..."

	switch len(errs) {
	case 0:
		return ""
	case 1:
		tail = ""
	}

	return fmt.Sprint("errors: ", len(errs), errs[0], tail)
}

func (e *multiError) Errors() []error {
	return e.errs
}
