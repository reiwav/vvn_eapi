package tibero

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Field struct {
	Name     string
	Value    string
	DataType string
	Length   string
	Insert   bool
}

func GetValueFields(key string, obj IModel) ([]Field, error) {
	var mapField = make([]Field, 0)
	var values = reflect.ValueOf(obj).Elem()
	mapField, err := getFieldCell(key, values)
	return mapField, err
}

func getFieldCell(key string, values reflect.Value) ([]Field, error) {
	var fileds = make([]Field, 0)
	var numField = values.NumField()
	for i := 0; i < numField; i++ {
		f := values.Type().Field(i)
		var tags = strings.Split(f.Tag.Get(key), ",")
		v := tags[0] // use split to ignore tag "options"
		if v == "" || v == "-" {
			continue
		}
		value := values.Field(i)
		if v == Inline {
			var fs, err = getFieldCell(key, value)
			if err != nil {
				return nil, err
			}
			fileds = append(fileds, fs...)
			continue
		}
		var curValue = value.Interface()
		curStr, len, dataType, isInsert := checkField(tags, value, curValue, f.Type)
		var field = Field{
			Name:     v,
			Value:    curStr,
			DataType: dataType,
			Length:   len,
			Insert:   isInsert,
		}
		fileds = append(fileds, field)
	}
	return fileds, nil
}

func checkField(tags []string, value reflect.Value, curValue interface{}, vType reflect.Type) (string, string, string, bool) {
	var isInsert = true
	var length string
	if len(tags) > 1 {
		var tag1 = tags[1]
		if strings.Contains(tag1, "omitempty") {
			isInsert = false
		} else if strings.Contains(tag1, "primary key") {
			length = tag1
			isInsert = false
		} else { //case tag1 is number
			length = tag1
		}
	}

	switch vType.String() {
	case "string":
		var typeD = TypeNvachar
		if len(length) == 0 {
			typeD = TypeClob
		}
		v := curValue.(string)
		if v == "" {
			return "''", length, typeD, isInsert
		}
		res := "'" + v + "'"

		return res, length, typeD, true
	case "String":
		var typeD = TypeNvachar
		if len(length) == 0 {
			typeD = TypeClob
		}
		v := curValue.(String)
		if v == "" {
			return "''", length, typeD, isInsert
		}
		res := "'" + v + "'"

		return string(res), length, typeD, true
	case "Bool":
		var res string
		v := curValue.(Bool)
		if v {
			res = "1"
		} else {
			res = "0"
		}
		return res, length, TypeBool, true
	case "bool":
		var res string
		v := curValue.(bool)
		if v {
			res = "1"
		} else {
			res = "0"
		}
		return res, length, TypeBool, true
	case "time.Time":
		var t time.Time
		v := curValue.(time.Time)
		if v == t {
			return "", length, TypeDate, false
		}

		res := ConvertTimeToDate(v)
		return res, length, TypeDate, isInsert
	case "int":
		return fmt.Sprintf("%v", curValue), length, TypeNumber, isInsert
	case "int64":
		return fmt.Sprintf("%v", curValue), length, TypeNumber, isInsert
	case "float32":
		return fmt.Sprintf("%v", curValue), length, TypeFloat, isInsert
	case "float64":
		return fmt.Sprintf("%v", curValue), length, TypeFloat, isInsert
	case "interface {}":
		return checkInsert(curValue, isInsert)
	default:
		return checkInsert(curValue, isInsert)
	}
}

func checkInsert(curValue interface{}, isInsert bool) (string, string, string, bool) {
	var isCheck bool
	var resValue string
	var dataType string
	var length string
	switch v := curValue.(type) {
	case string:
		resValue = "'" + v + "'"
		if v == "" && !isInsert {
			isCheck = false
		} else {
			isCheck = true
		}
		dataType = TypeNvachar
		length = "5000"
	case String:
		resValue = "'" + string(v) + "'"
		if v == "" && !isInsert {
			isCheck = false
		} else {
			isCheck = true
		}
		dataType = TypeNvachar
		length = "5000"
	case time.Time:
		var t = time.Time{}
		if v != t {
			isCheck = true
			resValue = ConvertTimeToDate(v)
		} else {
			isCheck = false
		}
		dataType = TypeDate
	case Bool:
		isCheck = true
		resValue = "0"
		if v {
			resValue = "1"
		}
		dataType = TypeBool
	case bool:
		isCheck = true
		resValue = "0"
		if v {
			resValue = "1"
		}
		dataType = TypeBool
	default:
		isCheck = isInsert
		dataType = TypeNClob
		resValue = fmt.Sprintf("%v", curValue)
	}
	return resValue, length, dataType, isCheck
}

func ConvertTimeToDate(v time.Time) string {
	var format = "2006-01-02 15:04:05"
	valFormat := v.Format(format)
	return "TO_DATE('" + valFormat + "', 'YY/MM/DD HH24:MI:SS')"
}
