package shared

import (
	"reflect"
	"strings"
)

const (
	KeyJson   = "json"
	KeyTibero = "rei"
)

func GetFieldNameTibero(fieldTagJson string, s interface{}) (fieldname string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(KeyJson), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == fieldTagJson {
			return strings.Split(f.Tag.Get(KeyTibero), ",")[0]
		}
	}
	return ""
}
