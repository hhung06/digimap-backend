package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

type BeaconResponse struct {
	ID        uuid.UUID  `json:"id"`
	VenueID   *uuid.UUID `json:"venue_id,omitempty"`
	LevelID   *uuid.UUID `json:"level_id,omitempty"`
	ElementID *uuid.UUID `json:"element_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	HwID      string     `json:"hw_id,omitempty"`
	VendorKey string     `json:"vendor_key,omitempty"`
	LotKey    string     `json:"lot_key,omitempty"`
	UUIDVal   string     `json:"uuid,omitempty"`
	MAC       string     `json:"mac,omitempty"`
	Radius    int        `json:"radius"`
	Battery   int        `json:"battery"`
	PositionX float64    `json:"position_x"`
	PositionY float64    `json:"position_y"`
	IsEnable  bool       `json:"is_enable"`
	Major     *int       `json:"major,omitempty"`
	Minor     *int       `json:"minor,omitempty"`
	Voltage   *int       `json:"voltage,omitempty"`
	TxPower   *int       `json:"tx_power,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type BeaconRequest struct {
	LevelID   *uuid.UUID `json:"level_id"`
	ElementID *uuid.UUID `json:"element_id"`
	Name      string     `json:"name"`
	HwID      string     `json:"hw_id"`
	VendorKey string     `json:"vendor_key"`
	LotKey    string     `json:"lot_key"`
	UUIDVal   string     `json:"uuid"`
	MAC       string     `json:"mac"`
	Radius    int        `json:"radius"`
	Battery   int        `json:"battery"`
	PositionX float64    `json:"position_x"`
	PositionY float64    `json:"position_y"`
	IsEnable  bool       `json:"is_enable"`
	Major     *int       `json:"major"`
	Minor     *int       `json:"minor"`
	Voltage   *int       `json:"voltage"`
	TxPower   *int       `json:"tx_power"`
}

type UpdateBeaconRequest struct {
	LevelID   *uuid.UUID `json:"level_id"`
	ElementID *uuid.UUID `json:"element_id"`
	Name      *string    `json:"name"`
	HwID      *string    `json:"hw_id"`
	VendorKey *string    `json:"vendor_key"`
	LotKey    *string    `json:"lot_key"`
	UUIDVal   *string    `json:"uuid"`
	MAC       *string    `json:"mac"`
	Radius    *int       `json:"radius"`
	Battery   *int       `json:"battery"`
	PositionX *float64   `json:"position_x"`
	PositionY *float64   `json:"position_y"`
	IsEnable  *bool      `json:"is_enable"`
	Major     *int       `json:"major"`
	Minor     *int       `json:"minor"`
	Voltage   *int       `json:"voltage"`
	TxPower   *int       `json:"tx_power"`
}

func (r UpdateBeaconRequest) ApplyTo(b *domain.Beacon) {
	if r.LevelID != nil {
		b.LevelID = r.LevelID
	}
	if r.ElementID != nil {
		b.ElementID = r.ElementID
	}
	if r.Name != nil {
		b.Name = *r.Name
	}
	if r.HwID != nil {
		b.HwID = *r.HwID
	}
	if r.VendorKey != nil {
		b.VendorKey = *r.VendorKey
	}
	if r.LotKey != nil {
		b.LotKey = *r.LotKey
	}
	if r.UUIDVal != nil {
		b.UUIDVal = *r.UUIDVal
	}
	if r.MAC != nil {
		b.MAC = *r.MAC
	}
	if r.Radius != nil {
		b.Radius = *r.Radius
	}
	if r.Battery != nil {
		b.Battery = *r.Battery
	}
	if r.PositionX != nil {
		b.PositionX = *r.PositionX
	}
	if r.PositionY != nil {
		b.PositionY = *r.PositionY
	}
	if r.IsEnable != nil {
		b.IsEnable = *r.IsEnable
	}
	if r.Major != nil {
		b.Major = r.Major
	}
	if r.Minor != nil {
		b.Minor = r.Minor
	}
	if r.Voltage != nil {
		b.Voltage = r.Voltage
	}
	if r.TxPower != nil {
		b.TxPower = r.TxPower
	}
}

func BeaconToResponse(b *domain.Beacon) BeaconResponse {
	return BeaconResponse{
		ID: b.ID, VenueID: b.VenueID, LevelID: b.LevelID, ElementID: b.ElementID,
		Name: b.Name, HwID: b.HwID, VendorKey: b.VendorKey, LotKey: b.LotKey,
		UUIDVal: b.UUIDVal, MAC: b.MAC,
		Radius: b.Radius, Battery: b.Battery, PositionX: b.PositionX, PositionY: b.PositionY,
		IsEnable: b.IsEnable, Major: b.Major, Minor: b.Minor, Voltage: b.Voltage, TxPower: b.TxPower,
		CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt,
	}
}
