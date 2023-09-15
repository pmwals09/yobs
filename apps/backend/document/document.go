package document

import (
	"database/sql"
	"io"
	"mime/multipart"
	"os"
)

type DocumentType string

const (
	Resume        DocumentType = "Resume"
	CoverLetter                = "CoverLetter"
	Communication              = "Communication"
)

type Document struct {
	ID       uint         `json:"id"`
	FileName string       `json:"fileName"`
	Title    string       `json:"title"`
	Type     DocumentType `json:"type"`
}

func CreateTable(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, docErr := tx.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			file_name TEXT,
			title TEXT UNIQUE,
			type TEXT
		);
	`)
	_, oppDocErr := tx.Exec(`
		CREATE TABLE IF NOT EXISTS opportunity_documents (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			opportunity_id INTEGER NOT NULL,
			document_id INTEGER NOT NULL,
			FOREIGN KEY (opportunity_id) REFERENCES opportunities(id),
			FOREIGN KEY (document_id) REFERENCES documents(id)
		);
		`)
	if docErr != nil || oppDocErr != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func New(fileName string, documentType DocumentType) *Document {
	return &Document{FileName: fileName, Type: documentType}
}

func (d *Document) WithTitle(title string) *Document {
	d.Title = title
	return d
}

type DocumentModel struct {
	DB *sql.DB
}

type Repository interface {
	CreateDocument()
}

func (d *Document) Upload(file multipart.File) error {
	destination, err := os.Create(d.FileName)
	defer destination.Close()
	if err != nil {
		return err
	}

	if _, err := io.Copy(destination, file); err != nil {
		return err
	}

	return nil
}

func (d *DocumentModel) CreateDocument(doc *Document) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO documents (
			file_name,
			title,
			type
		) VALUES ( $1, $2, $3 ) RETURNING id;
	`)
	if err != nil {
		return err
	}

	var id uint
	err = stmt.QueryRow(
		doc.FileName,
		doc.Title,
		doc.Type,
	).Scan(&id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	doc.ID = id
	return err
}
