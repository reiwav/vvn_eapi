package tibero

import (
	"errors"
	"strconv"
	"strings"
)

func dataTypeCells(model IModel) (string, error) {
	var fields, err = GetValueFields(KEY_TAG, model)
	if err != nil {
		return "", err
	}
	var cells string
	for _, val := range fields {
		var contructorCell, err = getDataTypeCell(val) // check datatype field
		if err != nil {
			return "", err
		}
		if cells == "" {
			cells += contructorCell
		} else {
			cells += "," + contructorCell
		}
	}

	return cells, err
}

func getDataTypeCell(field Field) (string, error) {
	var leng string
	if field.Name == "deleted_at" {
		field.DataType = TypeDate
	}
	if field.Length != "" {
		if field.Length == "primary key" {
			leng = " " + strings.ToUpper("primary key") + " NOT NULL"
		} else {
			var _, err = strconv.Atoi(field.Length)
			if err != nil {
				return "", errors.New("field " + field.Name + " not valid ," + field.Length)
			}
			leng = "(" + field.Length + ")"
		}
	}
	return field.Name + " " + strings.ToUpper(field.DataType) + leng, nil
}
