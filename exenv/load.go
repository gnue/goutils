package exenv

import (
	"os"

	"github.com/gnue/goutils/exenv/parser"
)

func Load(files ...string) error {
	for _, fname := range files {
		err := read(fname)
		if err != nil {
			return err
		}
	}

	return nil
}

func read(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	p := parser.NewParser(f)
	p.Parse(func(key, val string) {
		os.Setenv(key, os.ExpandEnv(val))
	})

	return nil
}
