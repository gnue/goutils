package parser

import (
	"errors"
	"fmt"
)

var (
	errEOS           = errors.New("end of strng")
	errInvalidEscape = errors.New("invalid escape")
	errInvalidQuote  = errors.New("invalid quote")
)

type expectError struct {
	token rune
}

func (e *expectError) Error() string {
	return fmt.Sprintf("expect: unexpected token %c", e.token)
}

type parseError struct {
	val string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("parse: error %q", e.val)
}
