package status

import (
    "time"
)

const (
	None          = "None"
	Applied       = "Applied"
	Rejected      = "Rejected"
	FollowedUp    = "Followed Up"
	Pending       = "Pending"
	Offer         = "Offer"
	AcceptedOffer = "Accepted Offer"
)

type Status struct {
	Name string
	Note string
    Date time.Time
}
