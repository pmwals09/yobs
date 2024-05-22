package status

import (
	"database/sql"
	"time"
)

const (
	None          = "None"
	Applied       = "Applied"
	FollowedUp    = "Followed Up"
	Interview     = "Interview"
	Pending       = "Pending"
	OfferExtended = "Offer Extended"
	OfferAccepted = "Offer Accepted"
	Rejected      = "Rejected"
	Closed        = "Closed"
	Archived      = "Archived"
)

type Status struct {
	ID uint
	Name string
	Note string
	Date time.Time
}

func (s Status) IsEmpty() bool {
	return s.Name == "" || s.Date.IsZero()
}

type StatusModel struct {
	DB *sql.DB
}

type Repository interface {
	DeleteStatusByID(id uint) error
	GetStatusByID(id uint) (Status, error)
	UpdateStatus(s Status) error
}

func (sm StatusModel) DeleteStatusByID(id uint) error {
	_, err := sm.DB.Exec(`DELETE FROM statuses WHERE id = ?`, id)
	return err
} 

func (sm StatusModel) GetStatusByID(id uint) (Status, error) {
	var out Status
	res := sm.DB.QueryRow(`
		SELECT
			id,
			name,
			note,
			date
		FROM statuses WHERE id = ?;
	`, id)
	err := res.Scan(
	&out.ID,
	&out.Name,
	&out.Note,
	&out.Date)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (sm StatusModel) UpdateStatus(s Status) error {
	_, err := sm.DB.Exec(`
		UPDATE statuses
		SET
			name = ?,
			date = ?,
			note = ?
		WHERE id = ?;
	`, s.Name, s.Date, s.Note, s.ID)
	return err
}
