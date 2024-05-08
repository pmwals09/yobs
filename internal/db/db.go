package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
