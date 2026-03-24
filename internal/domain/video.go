package domain

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID          uuid.UUID
	VenueID     *uuid.UUID
	Title       string
	Description string
	URL         string
	Thumbnail   string
	Duration    *int
	Status      string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
