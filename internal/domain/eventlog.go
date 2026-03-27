package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventLog represents a single analytics event emitted by a client app.
type EventLog struct {
	ID        uuid.UUID
	VenueID   *uuid.UUID
	Name      string // click, filter, view, search, ads_impression, etc.
	Params    json.RawMessage
	DeviceID  string
	UserID    string
	UserAgent string
	IPAddress string
	CreatedAt time.Time
}

// SearchQuery tracks search terms used within a venue.
type SearchQuery struct {
	ID           uuid.UUID
	VenueID      *uuid.UUID
	AppID        string
	Origin       string
	SearchTerm   string
	SearchCount  int
	LastSearched *time.Time
	IsPromoted   bool
	Reference    json.RawMessage
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
