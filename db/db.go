package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const DRIVER_NAME = "mysql"
const DB_URL = "root:secret@/devbook?charset=utf8&parseTime=True&loc=Local"

func Connect() (*sql.DB, error) {
	db, err := sql.Open(DRIVER_NAME, DB_URL)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
