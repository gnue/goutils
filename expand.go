package exenv

import (
	"os"
	"reflect"

	"github.com/gnue/tag"
)

const tagKey = "expandenv"

type options struct {
	Ignore bool // 無視する
}

func Expand(value interface{}, mapping func(string) string) {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			f := v.Type().Field(i)

			var opts options

			name := tag.Parse(f.Tag.Get(tagKey), &opts)
			if name == "-" || opts.Ignore {
				continue
			}

			fv := v.Field(i)
			switch fv.Kind() {
			case reflect.String:
				fv.SetString(os.Expand(fv.String(), mapping))
			case reflect.Struct, reflect.Map:
				Expand(fv.Interface(), mapping)
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if key.Kind() == reflect.Interface {
				key = key.Elem()
			}

			fv := v.MapIndex(key)
			if fv.Kind() == reflect.Interface {
				fv = fv.Elem()
			}

			switch fv.Kind() {
			case reflect.String:
				s := os.Expand(fv.String(), mapping)
				v.SetMapIndex(key, reflect.ValueOf(s))
			case reflect.Struct, reflect.Map:
				Expand(fv.Interface(), mapping)
			}
		}
	}
}

func ExpandEnv(value interface{}) {
	Expand(value, os.Getenv)
}
