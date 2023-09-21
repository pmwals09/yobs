package document

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
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
	URL         string       `json:"url"`
}

func CreateTable(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// NOTE: We don't store the URL. We use the URL when creating something
	// clickable in order to download a document, but it's time-limited by AWS
	// so we don't try to save it with the document
	_, docErr := tx.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			file_name TEXT,
			title TEXT UNIQUE,
			type TEXT,
			content_type TEXT
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
		fmt.Println("DOCERR", docErr.Error())
		fmt.Println("OPPDOCERR", oppDocErr.Error())
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
	// TODO: eventually preface with UUID of user

	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_S3_BUCKET")

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return err
	}
	c := s3.NewFromConfig(cfg)

	_, err = c.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(d.FileName),
		Body:        file,
		ContentType: aws.String(d.ContentType),
	})

	if err != nil {
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
			type,
			content_type
		) VALUES ( $1, $2, $3, $4 ) RETURNING id;
	`)
	if err != nil {
		return err
	}

	var id uint
	err = stmt.QueryRow(
		doc.FileName,
		doc.Title,
		doc.Type,
		doc.ContentType,
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

func (d *Document) GetPresignedDownloadUrl() (string, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET")
	region := os.Getenv("AWS_REGION")
	lifetimeSeconds := int64(60 * 30)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return "", err
	}

	c := s3.NewFromConfig(cfg)
	pc := s3.NewPresignClient(c)
	req, err := pc.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		// TODO: Preface with user's UUID
		Key: &d.FileName,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSeconds * int64(time.Second))
	})

	if err != nil {
		return "", err
	}

	d.URL = req.URL
	return req.URL, nil
}
