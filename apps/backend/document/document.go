package document

import (
	"gorm.io/gorm"
)

type DocumentType string

const (
	Resume        DocumentType = "Resume"
	CoverLetter                = "CoverLetter"
	Communication              = "Communication"
)

type Document struct {
	ID    uint         `gorm:"primary_key" json:"id"`
	Title string       `json:"title"`
	Type  DocumentType `json:"type"`
}

func NewDocument(title string, documentType DocumentType) Document {
	return Document{Title: title, Type: documentType}
}

type GormRepository struct {
	DB *gorm.DB
}

type Repository interface {
	Upload() error
}

func (*Document) Upload() error {
	// Upload method as required
	// Add to database as a document entry
	// Return id, err?
	return nil
}
