package contact

import (
	"database/sql"
	"errors"

	"github.com/pmwals09/yobs/internal/models/user"
)

type Contact struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CompanyName string `json:"Company"`
	Title       string `json:"Title"`
	Phone       string `json:"Phone"`
	Email       string `json:"Email"`
}

func (c Contact) IsEmpty() bool {
	return c.Name == "" && c.CompanyName == "" && c.Title == "" && c.Phone == "" && c.Email == ""
}

type ContactModel struct {
	DB *sql.DB
}

type Repository interface {
	CreateContact(contact *Contact, user user.User) error
	GetContactById(id uint, user user.User) (Contact, error)
}

func (cm ContactModel) CreateContact(contact *Contact, user user.User) error {
	if contact.IsEmpty() {
		return errors.New("Contact is empty")
	}
	res, err := cm.DB.Exec(`
        INSERT INTO contacts (
            name,
            company_name,
            title,
            phone,
            email,
			user_id
        ) VALUES (?, ?, ?, ?, ?, ?);
`,
		contact.Name,
		contact.CompanyName,
		contact.Title,
		contact.Phone,
		contact.Email,
		user.ID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	contact.ID = uint(id)
	return nil
}

func (cm ContactModel) GetContactById(id uint, user user.User) (Contact, error) {
	row := cm.DB.QueryRow(`
        SELECT
            id,
            name,
            company_name,
            title,
            phone,
            email
        FROM contacts WHERE id = ? AND user_id = ?;
`,
		id,
		user.ID)
	var contact Contact
	err := row.Scan(
		&contact.ID,
		&contact.Name,
		&contact.CompanyName,
		&contact.Title,
		&contact.Phone,
		&contact.Email)
	if err != nil {
		return contact, err
	}
	return contact, nil
}
