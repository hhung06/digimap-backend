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
