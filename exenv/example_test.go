package exenv_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gnue/goutils/exenv"
)

func TestMain(m *testing.M) {
	os.Setenv("USER", "gopher")
	os.Exit(m.Run())
}

// replaces ${var} or $var in struct/map
func ExampleExpandEnv() {
	var s = struct {
		User string
		Home string `expandenv:"-"`
	}{"USER=$USER", "HOME=$HOME"}

	exenv.ExpandEnv(&s)

	fmt.Println(s.User)
	fmt.Println(s.Home)

	// Output:
	//
	// USER=gopher
	// HOME=$HOME
}

// load file and os.Setenv
func ExampleLoad() {
	err := exenv.Load(".env")
	if err != nil {
		// ...
	}
}

// create Env struct and load file
func ExampleEnv() {
	e := exenv.New()
	err := e.Load(".env")
	if err != nil {
		// ...
	}
}
