package dotrc_test

import (
	"github.com/gnue/goutils/dotrc"
)

// execute command with .rc
func ExampleNew() {
	sh := dotrc.New(".rc")
	cmd := sh.Command("command")
	cmd.Run()
}
