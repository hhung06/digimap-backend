package domain

import (
	"time"

	"github.com/google/uuid"
)

// Asset represents a media file tracked in the database (metadata only — actual file is in S3).
type Asset struct {
	ID          uuid.UUID
	VenueID     uuid.UUID
	Name        string
	Key         string
	ContentType string
	SizeBytes   int64
	URL         string
	CreatedBy   *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
