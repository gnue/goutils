package execuser_test

import (
	"os/exec"
	"os/user"

	"github.com/gnue/goutils/execuser"
)

func isRootUser() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	return u.Uid == "0"
}

// run user by username
func ExampleLookup() {
	cmd := exec.Command("ps", "u")

	if isRootUser() {
		if u, err := execuser.Lookup("username"); err == nil {
			u.RunUser(cmd)
		}
	}

	cmd.Run()
}

// run user by path
func ExampleLookupPath() {
	cmd := exec.Command("ps", "u")

	if isRootUser() {
		if u, err := execuser.LookupPath("."); err == nil {
			u.RunUser(cmd)
		}
	}

	cmd.Run()
}
