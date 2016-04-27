package exenv

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

var reExport = regexp.MustCompile(`^export\s+`)

type Env struct {
	Data map[string]string
	Keys []string
}

func New(environ ...string) *Env {
	env := &Env{Data: make(map[string]string), Keys: make([]string, 0)}
	env.SetEnviron(environ...)
	return env
}

func LoadEnv(files ...string) (*Env, error) {
	env := New()

	err := env.Load(files...)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (env *Env) Getenv(key string) string {
	return env.Data[key]
}

func (env *Env) Setenv(key string, val string) error {
	if _, ok := env.Data[key]; !ok {
		env.Keys = append(env.Keys, key)
	}
	env.Data[key] = val
	return nil
}

func (env *Env) Environ() []string {
	d := env.Data
	a := make([]string, 0, len(d))

	for _, key := range env.Keys {
		a = append(a, key+"="+d[key])
	}

	return a
}

func (env *Env) SetEnviron(environ ...string) {
	for _, s := range environ {
		a := strings.SplitN(s, "=", 2)
		if len(a) == 2 {
			env.Setenv(a[0], a[1])
		}
	}
}

func (env *Env) ExpandEnv(s string) string {
	return os.Expand(s, env.Getenv)
}

func (env *Env) Expand(value interface{}) {
	Expand(value, env.Getenv)
}

func (env *Env) Load(files ...string) error {
	for _, fname := range files {
		err := env.read(fname)
		if err != nil {
			return err
		}
	}

	return nil
}

func (env *Env) read(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, val, err := parse(scanner.Text())
		if err != nil {
			continue
		}

		env.Setenv(key, val)
	}

	return nil
}

func parse(s string) (string, string, error) {
	s = strings.SplitN(s, "#", 2)[0]
	s = strings.Trim(s, " \t")
	if s == "" {
		return "", "", errors.New("skip line")
	}

	a := strings.SplitN(s, "=", 2)
	key := strings.Trim(reExport.ReplaceAllString(a[0], ""), " \t")
	val := unqoute(strings.Trim(a[1], " \t"))

	return key, val, nil
}

func unqoute(s string) string {
	if 0 < len(s) {
		switch s[0] {
		case '"':
			return strings.Trim(s, `"`)
		case '\'':
			return strings.Trim(s, `'`)
		}
	}

	return s
}
