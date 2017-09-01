package autoload

import (
	"github.com/gnue/exenv"
)

func init() {
	exenv.Load(".env")
}
