package domain

import (
	"time"

	"github.com/google/uuid"
)

type Advertisement struct {
	ID              uuid.UUID
	VenueID         *uuid.UUID
	LocationID      *uuid.UUID
	Type            string
	Status          string
	Navigate        *string
	ContentImageURL *string
	ContentCTAURL   *string
	Placement       string
	SizeWidth       *int
	SizeHeight      *int
	RewardType      *string
	RewardAmount    *int
	DisplayDuration *int
	PublishedAt     *time.Time
	StartAt         *time.Time
	EndAt           *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
