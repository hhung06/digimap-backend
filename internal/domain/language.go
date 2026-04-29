package domain

import (
	"time"

	"github.com/google/uuid"
)

type Language struct {
	ID        uuid.UUID
	VenueID   uuid.UUID
	Code      string
	Name      string
	IsDefault bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
