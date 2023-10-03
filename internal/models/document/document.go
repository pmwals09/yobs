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
	"github.com/pmwals09/yobs/internal/models/user"
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
	User        *user.User   `json:"user"`
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

func (d *Document) WithUser(u *user.User) *Document {
	d.User = u
	return d
}

func (d *Document) GetKey() string {
	return d.User.UUID.String() + "/" + d.FileName
}

type DocumentModel struct {
	DB *sql.DB
}

type Repository interface {
	CreateDocument(doc *Document) error
}

func (d *Document) Upload(file multipart.File) error {
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
		Key:         aws.String(d.GetKey()),
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
			content_type,
      user_id
		) VALUES ( $1, $2, $3, $4, $5 ) RETURNING id;
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
		doc.User.ID,
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
	key := d.GetKey()
	req, err := pc.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSeconds * int64(time.Second))
	})

	if err != nil {
		return "", err
	}

	d.URL = req.URL
	return req.URL, nil
}
