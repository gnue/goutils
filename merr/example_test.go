package merr_test

import (
	"fmt"
	"os"

	"github.com/gnue/goutils/merr"
)

var err error
var files = []string{"Alice", "Bob", "Charlie"}

// create multiple error
func ExampleNew() {
	var errs []error

	for _, s := range files {
		if _, err := os.Stat(s); err != nil {
			errs = append(errs, err)
		}
	}

	err = merr.New(errs...)
}

// get error list
func ExampleErrors() {
	if err != nil {
		if errs, ok := err.(merr.Errors); ok {
			for _, err := range errs.Errors() {
				fmt.Println(err)
			}
		}
	}
}
