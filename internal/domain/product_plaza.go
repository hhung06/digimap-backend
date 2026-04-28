package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductPlaza struct {
	ID          uuid.UUID
	VenueID     uuid.UUID
	Name        string
	Description string
	LocationID  *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
