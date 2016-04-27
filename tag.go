package tag

import (
	"reflect"
	"strings"
)

const noPrefix = "no"

// タグを解析する
func Parse(tag string, opts interface{}) (name string) {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return
	}

	name = trimTag(tags[0])
	updateOptions(tags[1:], opts)

	return
}

// オプションを更新する
func updateOptions(tags []string, opts interface{}) {
	v := reflect.ValueOf(opts).Elem()

	for _, name := range tags {
		flag := true

		name = trimTag(name)
		if strings.HasPrefix(name, noPrefix) {
			name = strings.TrimPrefix(name, noPrefix)
			flag = false
		}

		d := v.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, name)
		})

		if d.CanSet() && d.Kind() == reflect.Bool {
			d.SetBool(flag)
		}
	}
}

// タグの前後空白を削除
func trimTag(tag string) string {
	return strings.Trim(tag, " \t")
}
