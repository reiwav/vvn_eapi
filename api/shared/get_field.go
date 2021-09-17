package shared

import (
	"eapi/reiway/tibero"
	"fmt"
	"reflect"
	"strings"
)

const (
	KeyJson   = "json"
	KeyTibero = "rei"
)

func GetFieldNameTibero(fieldTagJson string, s interface{}) (fieldname string) {
	rt := reflect.TypeOf(s)
	fmt.Println(rt.Kind())
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		v := strings.Split(f.Tag.Get(KeyJson), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == "" {
			var tagTb = strings.Split(f.Tag.Get(KeyTibero), ",")
			if len(tagTb) > 0 && tagTb[0] == tibero.Inline {
				value := tibero.BaseModel{}
				data := GetFieldNameTibero(fieldTagJson, value)
				if data != "" {
					return data
				}
			}

		} else if v == fieldTagJson {
			return strings.Split(f.Tag.Get(KeyTibero), ",")[0]
		}
	}
	return ""
}
