package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventTag struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Localization json.RawMessage `json:"localization,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type EventType struct {
	ID           uuid.UUID       `json:"id"`
	VenueID      *uuid.UUID      `json:"venue_id,omitempty"`
	Name         string          `json:"name"`
	Localization json.RawMessage `json:"localization,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type Event struct {
	ID            uuid.UUID       `json:"id"`
	VenueID       uuid.UUID       `json:"venue_id"`
	TypeID        *uuid.UUID      `json:"type_id,omitempty"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	BannerImage   string          `json:"banner_image"`
	IconImage     string          `json:"icon_image"`
	StartTime     *time.Time      `json:"start_time,omitempty"`
	EndTime       *time.Time      `json:"end_time,omitempty"`
	ShowStartTime *time.Time      `json:"show_start_time,omitempty"`
	ShowEndTime   *time.Time      `json:"show_end_time,omitempty"`
	ContentDetail string          `json:"content_detail"`
	ContentURL    string          `json:"content_url"`
	Localization  json.RawMessage `json:"localization,omitempty"`
	Tags          []*EventTag     `json:"tags,omitempty"`
	Locations     []uuid.UUID     `json:"locations,omitempty"` // location IDs
	Images        []*EventImage   `json:"images,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type EventImage struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
