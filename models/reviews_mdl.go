package models

import (
	"time"
)

type Review struct {
	PartId, Rating                   int
	Subject, ReviewText, Name, Email string
	CreatedDate                      time.Time
}
