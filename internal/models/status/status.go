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
	ID   uint
	Name string
	Note string
	Date time.Time
}

func (s Status) IsEmpty() bool {
	return s.Name == "" || s.Date.IsZero()
}

func (s Status) ToFormDataValues() map[string]string {
	out := make(map[string]string)
	out["status-name"] = s.Name
	out["status-date"] = s.Date.Format(time.DateOnly)
	out["status-note"] = s.Note
	return out
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
	var dateStr string
	err := res.Scan(
		&out.ID,
		&out.Name,
		&out.Note,
		&dateStr)
	if err != nil {
		return out, err
	}

	if t, err := time.Parse(time.DateOnly, dateStr); err != nil {
		return out, err
	} else {
		out.Date = t
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
	`, s.Name, s.Date.Format(time.DateOnly), s.Note, s.ID)
	return err
}
