package tibero

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"database/sql"
)

func (t *Table) InsertOne(model IModel) error {
	model.BeforeCreate()
	var cells, values, err = cellValueInserts(model)
	if err != nil {
		return err
	}
	var q = fmt.Sprintf(InsertRow, t.TableName, cells, values)
	fmt.Println(err, q)
	var idLast int64
	err = t.DB.QueryRow(q).Scan(&idLast)
	model.setID(idLast)
	return err
}

func (t *Table) UpdateByID(model IModel) error {
	model.BeforeUpdate()
	var cells, where, err = cellValueUpdate(model, "id")
	if err != nil {
		return err
	}
	var q = fmt.Sprintf(UpdateRow, t.TableName, cells, "id", where)
	fmt.Println(q)
	_, err = t.DB.Exec(q)
	return err
}

func (t *Table) UnsafeUpdateByID(ID string, cells map[string]string) error {
	var q string
	for key, val := range cells {
		if q == "" {
			q += key + "=" + val
		} else {
			q += "," + key + "=" + val
		}
	}
	var query = fmt.Sprintf(UpdateRow, t.TableName, q, "id", ID)
	_, err := t.DB.Exec(query)
	return err
}

func (t *Table) UnsafeDeleteByID(ID string) error {
	var q = "deleted_at=" + ConvertTimeToDate(time.Now())
	var query = fmt.Sprintf(UpdateRow, t.TableName, q, "id", ID)
	_, err := t.DB.Exec(query)
	return err
}

func (t *Table) DeleteByID(model IModel) error {
	model.BeforeDelete()
	var cells, where, err = cellValueDelete(model, "id")
	if err != nil {
		return err
	}
	var q = fmt.Sprintf(UpdateRow, t.TableName, cells, "id", where)
	fmt.Println(q)
	_, err = t.DB.Exec(q)
	return err
}

func (t *Table) SelectOne(where map[string]string, orderBy string, skip int, limit int, v interface{}) error {
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}
	var q = fmt.Sprintf(SelectRow, t.TableName, qDate)
	if orderBy != "" {
		q += " ORDER BY " + orderBy
	}
	if limit > 0 {
		q = "SELECT * FROM (" + q + ")dt WHERE ROWNUM >= " + strconv.Itoa(skip) + " AND ROWNUM<=" + strconv.Itoa(limit) + ";"
	}
	fmt.Println(q)
	row, err := t.DB.Query(q)
	if err != nil {
		return err
	}
	return t.scanRowOne(v, row)
}

func (t *Table) SelectRows(where map[string]string, cols []string, orderBy string, skip int, limit int) (*sql.Rows, error) {
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}
	var q string
	if len(cols) == 0 {
		q = fmt.Sprintf(SelectRow, t.TableName, qDate)
	} else {
		var colStr string
		for i, val := range cols {
			if i == 0 {
				colStr += val
			} else {
				colStr += "," + val
			}
		}
		q = fmt.Sprintf(SelectCols, colStr, t.TableName, qDate)
	}

	if orderBy != "" {
		q += " ORDER BY " + orderBy
	}
	if limit > 0 {
		q = "SELECT * FROM (" + q + ")dt WHERE ROWNUM >= " + strconv.Itoa(skip) + " AND ROWNUM<=" + strconv.Itoa(limit) + ";"
	}
	return t.DB.Query(q)
}

func (t *Table) scanRowOne(v interface{}, row *sql.Rows) error {
	for row.Next() {
		fmt.Println("===== Vao 0:")
		s := reflect.ValueOf(v).Elem()
		fmt.Println("===== Vao 1:")
		rt := reflect.TypeOf(v).Elem()
		fmt.Println("===== Vao 2:")
		columns := t.getColumns(s, rt)
		err := row.Scan(columns...)
		if err != nil {
			fmt.Println("===== Vao Loi:", err.Error())
			return err
		}
		row.Close()
	}
	return nil
}

func (t *Table) UnsafeSelectOne(query string, v interface{}) error {
	row, err := t.DB.Query(query)
	if err != nil {
		return err
	}
	return t.scanRowOne(v, row)
}

func (t *Table) scanRows(v interface{}, r *sql.Rows) error {
	vType := reflect.TypeOf(v)
	if k := vType.Kind(); k != reflect.Ptr {
		return fmt.Errorf("%q must be a pointer: %w", k.String(), "not is pointer")
	}
	sliceType := vType.Elem()
	if reflect.Slice != sliceType.Kind() {
		return fmt.Errorf("%q must be a slice: %w", sliceType.String(), "not is slice")
	}

	sliceVal := reflect.Indirect(reflect.ValueOf(v))
	itemType := sliceType.Elem()
	var i = 1
	fmt.Println("== Vòng ", "ĐẾN ĐÂY")
	for r.Next() {
		fmt.Println("== Vòng ", fmt.Sprintf("%v", i))
		sliceItem := reflect.New(itemType).Elem()
		cols := t.getColumns(sliceItem, reflect.TypeOf(sliceItem))
		err := r.Scan(cols...)
		if err != nil {
			fmt.Println("===== VÀO LỖI", err.Error())
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, sliceItem))
		i++
	}
	return nil
}

func (t *Table) SelectMany(where map[string]string, orderBy string, skip int, limit int, v interface{}) error {
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}
	var q = fmt.Sprintf(SelectRow, t.TableName, qDate)
	if orderBy != "" {
		q += " ORDER BY " + orderBy
	}
	if limit > 0 {
		q = "SELECT * FROM (" + q + ")dt WHERE ROWNUM >= " + strconv.Itoa(skip) + " AND ROWNUM<=" + strconv.Itoa(limit) + ";"
	}
	fmt.Println("SelectMany:", q)
	rows, err := t.DB.Query(q)
	if err != nil {
		return err
	}
	return t.scanRows(v, rows)
}

func (t *Table) SelectCustomMany(where string, orderBy string, skip int, limit int, v interface{}) error {
	var q = fmt.Sprintf(SelectRow, t.TableName, where)
	if orderBy != "" {
		q += " ORDER BY " + strings.ToUpper(orderBy)
	}
	if limit > 0 {
		q = "SELECT * FROM (" + q + ")dt WHERE ROWNUM >= " + strconv.Itoa(skip) + " AND ROWNUM<=" + strconv.Itoa(limit) + ";"
	}
	fmt.Println("SelectCustomMany:", q)
	rows, err := t.DB.Query(q)
	if err != nil {
		fmt.Println("ERROR GET SCAN ROWS")
		return err
	}
	return t.scanRows(v, rows)
}

func (t *Table) UnsafeSelectMany(cells, where, orderBy string, skip int, limit int, v interface{}) error {
	var q = fmt.Sprintf(UnsafeSelect, t.TableName, cells, where)
	if orderBy != "" {
		q += " ORDER BY " + orderBy
	}
	if limit > 0 {
		q = "SELECT * FROM (" + q + ")dt WHERE ROWNUM >= " + strconv.Itoa(skip) + " AND ROWNUM<=" + strconv.Itoa(limit) + ";"
	}
	row, err := t.DB.Query(q)
	if err != nil {
		return err
	}
	return t.scanRows(v, row)
}

func (t *Table) getColumns(s reflect.Value, rt reflect.Type) []interface{} {
	var message_value = s
	if message_value.Kind() == reflect.Ptr {
		message_value = message_value.Elem()
	}
	numCols := message_value.NumField()
	columns := make([]interface{}, 0)
	for i := 0; i < numCols; i++ {
		field := message_value.Field(i)
		typeName := field.Type().Name()
		if field.Kind() == reflect.Slice {
			continue
		}
		// if typeName != "BaseModel" && tag == "" || tag == "-" {
		// 	continue
		// }
		switch typeName {
		case "BaseModel":
			var cl = t.getColumns(field, reflect.TypeOf(field))
			columns = append(columns, cl...)
			continue
		}

		columns = append(columns, field.Addr().Interface())
	}
	return columns
}

func (t *Table) SelectDistinct(where map[string]string, columns []string, v interface{}) error {
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}
	var cols string
	for i, val := range columns {
		if i == 0 {
			cols += val
			continue
		}
		cols += "," + val
	}
	var q = fmt.Sprintf(SelectDistinct, cols, t.TableName, qDate)
	rows, err := t.DB.Query(q)
	if err != nil {
		return err
	}
	return t.scanRows(v, rows)
}
