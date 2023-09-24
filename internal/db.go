package helpers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/user"
)

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	// TODO: Ideally we would have some kind of migrations system so we could
	// just run them all with one command, instead of relying on the practice of 
	// having a CreateTable method for each model in use
	opportunity.CreateTable(db)
	document.CreateTable(db)
	user.CreateTable(db)
	return db, nil
}
