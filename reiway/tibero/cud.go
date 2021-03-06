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
	t.GetDB()
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
	t.GetDB()
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
	t.GetDB()
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
	t.GetDB()
	var q = "deleted_at=" + ConvertTimeToDate(time.Now())
	var query = fmt.Sprintf(UpdateRow, t.TableName, q, "id", ID)
	_, err := t.DB.Exec(query)
	return err
}

func (t *Table) DeleteByID(model IModel) error {
	t.GetDB()
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

func (t *Table) SelectOne(model IModel, where map[string]string, orderBy string, skip int, limit int, v interface{}) error {
	t.GetDB()
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}
	var cells, err = t.GetCells(model)
	if err != nil {
		cells = "*"
	}
	var q = fmt.Sprintf(SelectRow, cells, t.TableName, qDate)
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
	if row != nil && row.Err() != nil {
		return row.Err()
	}
	return t.scanRowOne(v, row)
}

func (t *Table) SelectRows(model IModel, where map[string]string, cols []string, orderBy string, skip int, limit int) (*sql.Rows, error) {
	t.DB = GetDB()
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}

	var q string
	if len(cols) == 0 {
		var cells, err = t.GetCells(model)
		if err != nil {
			cells = "*"
		}
		q = fmt.Sprintf(SelectRow, cells, t.TableName, qDate)
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
	fmt.Println(q)
	return t.DB.Query(q)
}

func (t *Table) scanRowOne(v interface{}, row *sql.Rows) error {
	for row.Next() {
		s := reflect.ValueOf(v).Elem()
		rt := reflect.TypeOf(v).Elem()
		columns := t.getColumns(s, rt)
		err := row.Scan(columns...)
		if err != nil {
			return err
		}
		row.Close()
	}
	return nil
}

func (t *Table) UnsafeSelectOne(query string, v interface{}) error {
	t.GetDB()
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
	for r.Next() {
		sliceItem := reflect.New(itemType).Elem()
		cols := t.getColumns(sliceItem, reflect.TypeOf(sliceItem))
		err := r.Scan(cols...)
		if err != nil {
			fmt.Println("===== V??O L???I", err.Error())
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, sliceItem))
		i++
	}
	return nil
}

func (t *Table) SelectSkipLimit(where map[string]string, model IModel, orderBy string, skip int, limit int, v interface{}) error {
	t.GetDB()
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}
	var cells, err = t.GetCells(model)
	if err != nil {
		cells = "*"
	}

	var q = fmt.Sprintf(SelectRow, cells, t.TableName, qDate)
	if orderBy != "" {
		q += " ORDER BY " + orderBy
	}

	if limit > 0 {
		q = "SELECT " + cells + " FROM (SELECT a.*, rownum rowcell FROM ( " + q + " ) a WHERE rownum < " + strconv.Itoa(limit) + ") WHERE rowcell >= " + strconv.Itoa(skip) + ";"
	}

	rows, err := t.DB.Query(q)
	if err != nil {
		return err
	}
	return t.scanRows(v, rows)
}

func (t *Table) SelectMany(model IModel, where map[string]string, orderBy string, skip int, limit int, v interface{}) error {
	t.GetDB()
	var qDate = " deleted_at is null "
	for key, val := range where {
		qDate += " AND " + key + "=" + val
	}
	var cells, err = t.GetCells(model)
	if err != nil {
		cells = "*"
	}

	var q = fmt.Sprintf(SelectRow, cells, t.TableName, qDate)
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

func (t *Table) SelectCustomMany(model IModel, where string, orderBy string, skip int, limit int, v interface{}) error {
	t.GetDB()
	var cells, err = t.GetCells(model)
	if err != nil {
		cells = "*"
	}

	var q = fmt.Sprintf(SelectRow, cells, t.TableName, where)
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

func (t *Table) SelectCustomSkipLimit(where string, model IModel, orderBy string, skip int, limit int, v interface{}) error {
	t.GetDB()
	var cells, err = t.GetCells(model)
	if err != nil {
		cells = "*"
	}
	var q = fmt.Sprintf(SelectRow, cells, t.TableName, where)
	if orderBy != "" {
		q += " ORDER BY " + strings.ToUpper(orderBy)
	}

	if limit > 0 {
		q = "SELECT " + cells + " FROM (SELECT a.*, rownum rowcell FROM ( " + q + " ) a WHERE rownum < " + strconv.Itoa(limit) + ") WHERE rowcell >= " + strconv.Itoa(skip) + ";"
	}
	fmt.Println(q)
	rows, err := t.DB.Query(q)
	if err != nil {
		return err
	}
	return t.scanRows(v, rows)
}

func (t *Table) Count(where map[string]string) (int64, error) {
	t.GetDB()
	var q = fmt.Sprintf(SelectCount, t.TableName)
	if len(where) > 0 {
		var w string
		for i, val := range where {
			if w == "" {
				w += i + "=" + val
			} else {
				w += " AND " + val + "=" + val
			}
		}
		q = fmt.Sprintf(SelectCountWhere, t.TableName, w)
	}

	row := t.DB.QueryRow(q)
	var total int64
	err := row.Scan(&total)
	return total, err
}

func (t *Table) UnsafeSelectMany(cells, where, orderBy string, skip int, limit int, v interface{}) error {
	t.GetDB()
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
	t.GetDB()
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
	fmt.Println("Distinic: ", q)
	if err != nil {
		return err
	}
	return t.scanRows(v, rows)
}
