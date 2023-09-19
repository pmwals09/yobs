package document

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type DocumentType string

const (
	Resume        DocumentType = "Resume"
	CoverLetter                = "CoverLetter"
	Communication              = "Communication"
)

type Document struct {
	ID          uint         `json:"id"`
	FileName    string       `json:"fileName"`
	Title       string       `json:"title"`
	Type        DocumentType `json:"type"`
	ContentType string       `json:"contentType"`
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

func New(handler *multipart.FileHeader, documentType DocumentType) *Document {
	fmt.Println("CONTENTTYPE", handler.Header.Get("Content-Type"))
	return &Document{
		FileName:    handler.Filename,
		Type:        documentType,
		ContentType: handler.Header.Get("Content-Type"),
	}
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
	// how to use aws?
	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
	// Content-MD5 header for md5 of the file
	// IAM must have s3:PutObject
	// eventually preface with UUID of user

	region := os.Getenv("AWS_REGION")
	key := os.Getenv("AWS_ACCESS_KEY")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET")

	s3Config := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	}
	s3Session, err := session.NewSession(s3Config)

	uploader := s3manager.NewUploader(s3Session)

	f, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	input := &s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(d.FileName), // TODO: add user's UUID in front
		Body:        bytes.NewReader(f),
		ContentType: aws.String(d.ContentType),
	}

	// TODO: write file download link to database
	output, err := uploader.UploadWithContext(context.Background(), input)
	if err != nil {
		return err
	}
	fmt.Println(output)

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
