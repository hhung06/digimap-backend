package domain

import (
	"time"

	"github.com/google/uuid"
)

type FeaturedZone struct {
	ID          uuid.UUID
	VenueID     uuid.UUID
	Name        string
	Description string
	ImageURL    string
	SortIndex   int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
