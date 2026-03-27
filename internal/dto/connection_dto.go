package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

// ── Connection ────────────────────────────────────────────────────────────────

type ConnectionResponse struct {
	ID         uuid.UUID                 `json:"id"`
	VenueID    *uuid.UUID                `json:"venue_id,omitempty"`
	ExternalID string                    `json:"external_id,omitempty"`
	Name       string                    `json:"name,omitempty"`
	Type       int                       `json:"type"`
	X          float64                   `json:"x"`
	Y          float64                   `json:"y"`
	State      int                       `json:"state"`
	Status     int                       `json:"status"`
	Accessible bool                      `json:"accessible"`
	Active     bool                      `json:"active"`
	Levels     []ConnectionLevelResponse `json:"levels,omitempty"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
}

type ConnectionLevelResponse struct {
	ID           uuid.UUID  `json:"id"`
	ConnectionID uuid.UUID  `json:"connection_id"`
	LevelID      *uuid.UUID `json:"level_id,omitempty"`
	ElementID    *uuid.UUID `json:"element_id,omitempty"`
	Active       bool       `json:"active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ConnectionRequest struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	Type       int     `json:"type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	State      int     `json:"state"`
	Status     int     `json:"status"`
	Accessible bool    `json:"accessible"`
	Active     bool    `json:"active"`
}

type ConnectionLevelRequest struct {
	LevelID   *uuid.UUID `json:"level_id"`
	ElementID *uuid.UUID `json:"element_id"`
	Active    bool       `json:"active"`
}

func ConnectionToResponse(c *domain.Connection) ConnectionResponse {
	resp := ConnectionResponse{
		ID: c.ID, VenueID: c.VenueID, ExternalID: c.ExternalID, Name: c.Name,
		Type: c.Type, X: c.X, Y: c.Y, State: c.State, Status: c.Status,
		Accessible: c.Accessible, Active: c.Active,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
	for _, cl := range c.Levels {
		resp.Levels = append(resp.Levels, ConnectionLevelToResponse(cl))
	}
	return resp
}

func ConnectionLevelToResponse(cl *domain.ConnectionLevel) ConnectionLevelResponse {
	return ConnectionLevelResponse{
		ID: cl.ID, ConnectionID: cl.ConnectionID,
		LevelID: cl.LevelID, ElementID: cl.ElementID,
		Active: cl.Active, CreatedAt: cl.CreatedAt, UpdatedAt: cl.UpdatedAt,
	}
}

// ── QRCode ────────────────────────────────────────────────────────────────────

type QRCodeResponse struct {
	ID          uuid.UUID  `json:"id"`
	VenueID     *uuid.UUID `json:"venue_id,omitempty"`
	LevelID     *uuid.UUID `json:"level_id,omitempty"`
	LocationID  *uuid.UUID `json:"location_id,omitempty"`
	Lat         *float64   `json:"lat,omitempty"`
	Lng         *float64   `json:"lng,omitempty"`
	Angle       *float64   `json:"angle,omitempty"`
	Link        string     `json:"link,omitempty"`
	Base64Image string     `json:"base64_image,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type QRCodeRequest struct {
	LevelID    *uuid.UUID `json:"level_id"`
	LocationID *uuid.UUID `json:"location_id"`
	Lat        *float64   `json:"lat"`
	Lng        *float64   `json:"lng"`
	Angle      *float64   `json:"angle"`
	Link       string     `json:"link"`
}

func QRCodeToResponse(q *domain.QRCode) QRCodeResponse {
	return QRCodeResponse{
		ID: q.ID, VenueID: q.VenueID, LevelID: q.LevelID, LocationID: q.LocationID,
		Lat: q.Lat, Lng: q.Lng, Angle: q.Angle,
		Link: q.Link, Base64Image: q.Base64Image,
		CreatedAt: q.CreatedAt, UpdatedAt: q.UpdatedAt,
	}
}
