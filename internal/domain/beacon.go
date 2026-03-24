package domain

import (
	"time"

	"github.com/google/uuid"
)

// Beacon represents a BLE beacon device associated with a venue/level.
type Beacon struct {
	ID        uuid.UUID
	VenueID   *uuid.UUID
	LevelID   *uuid.UUID
	ElementID *uuid.UUID
	Name      string
	HwID      string
	VendorKey string
	LotKey    string
	UUIDVal   string
	MAC       string
	Radius    int
	Battery   int
	PositionX float64
	PositionY float64
	IsEnable  bool
	Major     *int
	Minor     *int
	Voltage   *int
	TxPower   *int
	CreatedAt time.Time
	UpdatedAt time.Time
}
