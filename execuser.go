package execuser

import (
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

	return convert(u)
}

func Lookup(username string) (*User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}

	return convert(u)
}

func LookupId(uid int) (*User, error) {
	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return nil, err
	}

	return convert(u)
}

func LookupPath(path string) (*User, error) {
	var s syscall.Stat_t

	err := syscall.Stat(path, &s)
	if err != nil {
		return nil, err
	}

	return LookupId(int(s.Uid))
}

func convert(u *user.User) (*User, error) {
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

func (u *User) SysProcAttr() *syscall.SysProcAttr {
	attr := &syscall.SysProcAttr{}
	attr.Credential = &syscall.Credential{Uid: u.Uid, Gid: u.Gid}
	return attr
}
