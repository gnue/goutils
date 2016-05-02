package exenv

import (
	"bufio"
	"os"

	"github.com/gnue/exenv/parser"
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

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p := parser.NewParser(scanner.Text())
		p.Parse(func(key, val string) {
			os.Setenv(key, val)
		})
	}

	return nil
}
