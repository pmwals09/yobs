package helpers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
)

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	opportunity.CreateTable(db)
	document.CreateTable(db)
	return db, nil
}
