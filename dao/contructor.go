package dao

import (
	"database/sql"
)

var database *sql.DB

func Database() *sql.DB {
	return database
}

func SetDatabase(db *sql.DB) {
	database = db
}
