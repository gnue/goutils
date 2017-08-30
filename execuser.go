package execuser

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

type User struct {
	Uid      uint32 // user id
	Gid      uint32 // primary group id
	Username string
	Name     string
	HomeDir  string
}

func Current() (*User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	return New(u)
}

func Lookup(username string) (*User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}

	return New(u)
}

func LookupId(uid int) (*User, error) {
	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return nil, err
	}

	return New(u)
}

func LookupPath(path string) (*User, error) {
	var s syscall.Stat_t

	err := syscall.Stat(path, &s)
	if err != nil {
		return nil, err
	}

	return LookupId(int(s.Uid))
}

func New(u *user.User) (*User, error) {
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return nil, err
	}

	return &User{Uid: uint32(uid), Gid: uint32(gid), Username: u.Username, Name: u.Name, HomeDir: u.HomeDir}, nil
}

func (u *User) RunUser(cmd *exec.Cmd) {
	cmd.SysProcAttr = u.SysProcAttr()
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	cmd.Env = append(cmd.Env, "USER="+u.Username, "HOME="+u.HomeDir)
}

func (u *User) SysProcAttr() *syscall.SysProcAttr {
	attr := &syscall.SysProcAttr{}
	attr.Credential = &syscall.Credential{Uid: u.Uid, Gid: u.Gid}
	return attr
}
