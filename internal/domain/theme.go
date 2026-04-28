package domain

import (
	"time"

	"github.com/google/uuid"
)

type Theme struct {
	ID             uuid.UUID
	VenueID        uuid.UUID
	Name           string
	PrimaryColor   string
	SecondaryColor string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
