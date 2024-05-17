package status

import (
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
	Name string
	Note string
	Date time.Time
}

func (s Status) IsEmpty() bool {
	return s.Name == "" || s.Date.IsZero()
}
