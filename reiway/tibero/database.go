package tibero

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

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

var dbCache *sql.DB
var m sync.Mutex

func GetDB() *sql.DB {
	m.Lock()
	defer m.Unlock()
	return dbCache
}

func (d *ConfigDB) Connect(isLoop bool) (*sql.DB, error) {
	db, err := sql.Open(d.DriverName, d.DSN)
	dbCache = db
	if err != nil {
		dbCache = nil
	}
	if isLoop {
		go d.launch()
	}

	return dbCache, err
}

func (d *ConfigDB) launch() {
	var t = time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-t.C:
			d.ping()
		}
	}

}

func (d *ConfigDB) ping() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if GetDB() == nil {
		d.Connect(false)
	} else {
		_, err := GetDB().Exec("select * from v$instance;")
		fmt.Println("========= PING ====", err)
		//db.SetConnMaxLifetime(3600 * time.Hour)
		if err != nil {
			GetDB().Close()
			d.Connect(false)
		}
	}
}
