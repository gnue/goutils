# tag

parse struct tag

## Installation

```sh
$ go get github.com/gnue/goutils/tag
```

## Usage

```go
import "github.com/gnue/goutils/tag"
```

## Examples

### Parse

```go
package main

import (
	"fmt"
	"reflect"

	"github.com/gnue/goutils/tag"
)

type Opts struct {
	Ignore bool
	Star   bool
}

type Data struct {
	Name   string `👤:",star"`
	Age    int    `👤:""`
	Secret string `👤:",ignore"`
}

// parse struct tag with reflect
//
// * "-" name is ignore field
// * "no~" flag is false
func main() {
	v := reflect.ValueOf(Data{"Alice", 10, "🎩🐰⏰👗"})

	if v.Kind() == reflect.Struct {
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			f := v.Type().Field(i)
			if f.PkgPath != "" && !f.Anonymous { // unexported
				continue
			}

			var opts Opts
			name := tag.Parse(f.Tag.Get("👤"), &opts)
			if name == "-" || opts.Ignore {
				continue
			}
			if name == "" {
				name = f.Name
			}

			fv := v.Field(i).Interface()
			s := fmt.Sprint(fv)
			if opts.Star {
				s += " ⭐️"
			}

			fmt.Printf("%-6s %s\n", name+":", s)
		}
	}
}

```

Output:

```
Name:  Alice ⭐️
Age:   10

```

parse struct tag with reflect

* "-" name is ignore field
* "no~" flag is false

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

