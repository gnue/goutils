package tag_test

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
func ExampleParse() {
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
	// Output:
	//
	// Name:  Alice ⭐️
	// Age:   10
}
