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

func (c Contact) ToFormDataValues() map[string]string {
	out := make(map[string]string)
	out["contact-name"] = c.Name
	out["contact-company-name"] = c.CompanyName
	out["contact-title"] = c.Title
	out["contact-phone"] = c.Phone
	out["contact-email"] = c.Email
	return out
}

type ContactModel struct {
	DB *sql.DB
}

type Repository interface {
	CreateContact(contact *Contact, user user.User) error
	GetContactById(id uint, user user.User) (Contact, error)
	UpdateContact(contact Contact) error
	DeleteContact(opptyID uint, contact Contact) error
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

func (cm ContactModel) UpdateContact(contact Contact) error {
	_, err := cm.DB.Exec(`
		UPDATE contacts
		SET
			name = ?,
			company_name = ?,
			title = ?,
			phone = ?,
			email = ?
		WHERE id = ?;
	`,
		contact.Name,
		contact.CompanyName,
		contact.Title,
		contact.Phone,
		contact.Email,
		contact.ID)
	return err
}

func (cm ContactModel) DeleteContact(opptyID uint, contact Contact) error {
	tx, err := cm.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		DELETE FROM contacts
		WHERE id = ?;
	`, contact.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`
		DELETE FROM opportunity_contacts
		WHERE opportunity_id = ?
		AND contact_id = ?;
	`, opptyID, contact.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
