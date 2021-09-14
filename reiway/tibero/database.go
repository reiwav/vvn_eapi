package tibero

import (
	"database/sql"
	//
	_ "github.com/alexbrainman/odbc"
)

const (
	KEY_TAG = "rei"
)

type ConfigDB struct {
	DriverName string `json:"driver"` //odbc
	DSN        string `json:"dsn"`    //"DSN=tibero"
	Username   string `json:"username"`
	Password   string `json:"password"`
}

func (d *ConfigDB) Connect() (*sql.DB, error) {
	return sql.Open(d.DriverName, d.DSN)
}
