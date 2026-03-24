package domain

import (
	"time"

	"github.com/google/uuid"
)

type Connection struct {
	ID         uuid.UUID
	VenueID    *uuid.UUID
	ExternalID string
	Name       string
	Type       int
	X          float64
	Y          float64
	State      int
	Status     int
	Accessible bool
	Active     bool
	Levels     []*ConnectionLevel
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ConnectionLevel struct {
	ID           uuid.UUID
	ConnectionID uuid.UUID
	LevelID      *uuid.UUID
	ElementID    *uuid.UUID
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type QRCode struct {
	ID          uuid.UUID
	VenueID     *uuid.UUID
	LevelID     *uuid.UUID
	LocationID  *uuid.UUID
	Lat         *float64
	Lng         *float64
	Angle       *float64
	Link        string
	Base64Image string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
