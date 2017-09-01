package autoload

import (
	"github.com/gnue/goutils/exenv"
)

func init() {
	exenv.Load(".env")
}
