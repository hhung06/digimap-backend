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

type UpdateConnectionRequest struct {
	ExternalID *string  `json:"external_id"`
	Name       *string  `json:"name"`
	Type       *int     `json:"type"`
	X          *float64 `json:"x"`
	Y          *float64 `json:"y"`
	State      *int     `json:"state"`
	Status     *int     `json:"status"`
	Accessible *bool    `json:"accessible"`
	Active     *bool    `json:"active"`
}

func (r UpdateConnectionRequest) ApplyTo(c *domain.Connection) {
	if r.ExternalID != nil {
		c.ExternalID = *r.ExternalID
	}
	if r.Name != nil {
		c.Name = *r.Name
	}
	if r.Type != nil {
		c.Type = *r.Type
	}
	if r.X != nil {
		c.X = *r.X
	}
	if r.Y != nil {
		c.Y = *r.Y
	}
	if r.State != nil {
		c.State = *r.State
	}
	if r.Status != nil {
		c.Status = *r.Status
	}
	if r.Accessible != nil {
		c.Accessible = *r.Accessible
	}
	if r.Active != nil {
		c.Active = *r.Active
	}
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
