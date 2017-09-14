# exenv

extend Env

 * Env struct
 * ExpandEnv for struct/map
 * Load env file
 * Autoload

## Installation

```sh
$ go get github.com/gnue/goutils/exenv
```

## Usage

```go
import "github.com/gnue/goutils/exenv"
```

## Examples

### Env

```go
package main

import (
	"github.com/gnue/goutils/exenv"
)

func main() {
	e := exenv.New()
	err := e.Load(".env")
	if err != nil {
		// ...
	}
}

```

create Env struct and load file

### ExpandEnv

```go
package main

import (
	"fmt"
	"github.com/gnue/goutils/exenv"
)

func main() {
	var s = struct {
		User string
		Home string `expandenv:"-"`
	}{"USER=$USER", "HOME=$HOME"}

	exenv.ExpandEnv(&s)

	fmt.Println(s.User)
	fmt.Println(s.Home)

}

```

Output:

```
USER=gopher
HOME=$HOME

```

replaces ${var} or $var in struct/map

### Load

```go
package main

import (
	"github.com/gnue/goutils/exenv"
)

func main() {
	err := exenv.Load(".env")
	if err != nil {
		// ...
	}
}

```

load file and os.Setenv

### Autload

```go
import "github.com/gnue/goutils/exenv/autoload"
```

load `.env` and os.Setenv

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

