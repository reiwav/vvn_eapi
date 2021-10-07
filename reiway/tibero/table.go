package tibero

import (
	"fmt"
	"strings"

	"database/sql"
)

type Table struct {
	TableName string
	DB        *sql.DB
}

func (tb *Table) GetDB() {
	tb.DB = GetDB()
}

func (tb *Table) NewTable(model IModel) error {
	tb.GetDB()
	var res interface{}
	var selectRow = strings.ReplaceAll(SelectRowOne, "@p", tb.TableName)
	err := tb.DB.QueryRow(selectRow).Scan(&res)
	fmt.Println("===== SLECT ROW ", err)
	if err != nil && err.Error() != "sql: no rows in result set" {
		contructorTable, err := dataTypeCells(model)
		if err != nil {
			return err
		}

		var f = fmt.Sprintf(CreateTable, tb.TableName, contructorTable)
		fmt.Println(f)
		_, err = tb.DB.Exec(f)
		if err != nil {
			return err
		}
		err = tb.CreateSequanceGenID()
		if err != nil {
			return err
		}
		err = tb.CreateTriggerGenID()
		if err != nil {
			return err
		}
	}
	return nil
}

func NewTable(tbName string, db *sql.DB, model IModel) error {
	var t = Table{
		TableName: tbName,
		DB:        GetDB(),
	}
	var res interface{}
	var selectRow = strings.ReplaceAll(SelectRowOne, "@p", tbName)
	err := db.QueryRow(selectRow).Scan(&res)
	fmt.Println("===== SLECT ROW ", err)
	if err != nil && err.Error() != "sql: no rows in result set" {
		contructorTable, err := dataTypeCells(model)
		if err != nil {
			return err
		}
		var f = fmt.Sprintf(CreateTable, tbName, contructorTable)
		fmt.Println(f)
		_, err = db.Exec(f)
		if err != nil {
			return err
		}
		err = t.CreateSequanceGenID()
		if err != nil {
			return err
		}
		err = t.CreateTriggerGenID()
		if err != nil {
			return err
		}
	}
	return nil
}

type LstModel []IModel

func (t *Table) InsertAll(models LstModel) error {
	t.GetDB()
	var q = BeginInsertAll
	for _, model := range models {
		model.BeforeCreate()
		var cells, values, err = cellValueInserts(model)
		if err != nil {
			return err
		}
		q += fmt.Sprintf(BodyInsertAll, t.TableName, cells, values)
	}
	q += EndInsertAll
	_, err := t.DB.Exec(q)
	return err
}

func (t *Table) CreateSequanceGenID() error {
	t.GetDB()
	query := strings.ReplaceAll(SequenceGenID, "@p", t.TableName)
	fmt.Println("===== ", query)
	res, err := t.DB.Exec(query)

	if err != nil {
		fmt.Println("ERR CREATE SEQUANCE TABLE", err)
		return err
	} else {
		fmt.Println(res)
	}
	return nil
}

func (t *Table) CreateTriggerGenID() error {
	t.GetDB()
	var query = strings.ReplaceAll(TriggerGenID, "@p", t.TableName)
	fmt.Println("===== ", query)
	res, err := t.DB.Exec(query)
	if err != nil {
		fmt.Println("ERR CREATE TRIGGER TABLE", err)
		return err
	} else {
		fmt.Println(res)
	}
	return nil
}
