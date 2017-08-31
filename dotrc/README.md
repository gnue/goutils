# dotrc

execute command with .rc

## Installation

```sh
$ go get github.com/gnue/goutils/dotrc
```

## Usage

```go
import "github.com/gnue/goutils/dotrc"
```

## Examples

### New

```go
package main

import (
	"github.com/gnue/goutils/dotrc"
)

func main() {
	sh := dotrc.New(".rc")
	cmd := sh.Command("command")
	cmd.Run()
}

```

execute command with .rc

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

