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
}

func (sm StatusModel) DeleteStatusByID(id uint) error {
	_, err := sm.DB.Exec(`DELETE FROM statuses WHERE id = ?`, id)
	return err
} 
