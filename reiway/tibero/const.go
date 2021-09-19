package tibero

const (
	SelectRowOne  = "SELECT ID FROM @p WHERE rownum = 1;"
	SequenceGenID = "CREATE SEQUENCE seq_@p START WITH 1;"
	TriggerGenID  = "CREATE OR REPLACE TRIGGER tg_@p " +
		"BEFORE INSERT ON @p " +
		"FOR EACH ROW " +
		"WHEN (new.id IS NULL) " +
		"BEGIN " +
		"SELECT seq_@p.NEXTVAL " +
		"INTO   :new.id " +
		"FROM   dual; " +
		"END;"
	CreateTable      = "CREATE TABLE %v(%v);"
	AlterTable       = "ALTER TABLE %v(%v);"
	InsertRow        = "INSERT INTO %v(%v) VALUES (%v) RETURNING ID INTO :C;"
	UpdateRow        = "UPDATE %v SET %v WHERE %v=%v;"
	DeleteRow        = "DELETE FROM %v WHERE %v=%v"
	SelectRow        = "SELECT * FROM %v WHERE %v"
	SelectCount      = "SELECT COUNT(*) FROM %v "
	SelectCountWhere = "SELECT COUNT(*) FROM %v WHERE %v"
	SelectCols       = "SELECT %v FROM %v WHERE %v"
	UnsafeSelect     = "SELECT %v FROM %v WHERE %v"
	BeginInsertAll   = "INSERT ALL "
	BodyInsertAll    = "INTO %v(%v) VALUES (%v) "
	EndInsertAll     = "SELECT * FROM DUAL;"
	SelectDistinct   = "SELECT DISTINCT %v  FROM %v WHERE %v;"
)

const (
	Inline = "inline"
)

const (
	TypeNumber = "NUMBER"
	TypeClob   = "CLOB"
	TypeNClob  = "NCLOB"
	TypeDate   = "DATE"
	TypeBool   = "NUMBER(1) DEFAULT 0 NOT NULL"
	Typevachar = "VARCHAR"
	TypeFloat  = "FLOAT"
)
