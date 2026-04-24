package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

// ── EventLog ─────────────────────────────────────────────────────────────────

type EventLogResponse struct {
	ID        uuid.UUID       `json:"id"`
	VenueID   *uuid.UUID      `json:"venue_id,omitempty"`
	Name      string          `json:"name"`
	Params    json.RawMessage `json:"params,omitempty" swaggertype:"object"`
	DeviceID  string          `json:"device_id,omitempty"`
	UserID    string          `json:"user_id,omitempty"`
	UserAgent string          `json:"user_agent,omitempty"`
	IPAddress string          `json:"ip_address,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// TrackEventRequest is the public-API payload for logging a single event.
type TrackEventRequest struct {
	Name     string          `json:"name" binding:"required"`
	Params   json.RawMessage `json:"params" swaggertype:"object"`
	DeviceID string          `json:"device_id"`
	UserID   string          `json:"user_id"`
}

func EventLogToResponse(e *domain.EventLog) EventLogResponse {
	return EventLogResponse{
		ID: e.ID, VenueID: e.VenueID, Name: e.Name, Params: e.Params,
		DeviceID: e.DeviceID, UserID: e.UserID,
		UserAgent: e.UserAgent, IPAddress: e.IPAddress,
		CreatedAt: e.CreatedAt,
	}
}

// ── SearchQuery ───────────────────────────────────────────────────────────────

type SearchQueryResponse struct {
	ID           uuid.UUID       `json:"id"`
	VenueID      *uuid.UUID      `json:"venue_id,omitempty"`
	AppID        string          `json:"app_id,omitempty"`
	SearchTerm   string          `json:"search_term"`
	SearchCount  int             `json:"search_count"`
	LastSearched *time.Time      `json:"last_searched,omitempty"`
	IsPromoted   bool            `json:"is_promoted"`
	Reference    json.RawMessage `json:"reference,omitempty" swaggertype:"object"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// TrackSearchRequest is the public-API payload for recording a search term.
type TrackSearchRequest struct {
	Term     string `json:"term" binding:"required"`
	DeviceID string `json:"device_id"`
	UserID   string `json:"user_id"`
}

func SearchQueryToResponse(q *domain.SearchQuery) SearchQueryResponse {
	return SearchQueryResponse{
		ID: q.ID, VenueID: q.VenueID, AppID: q.AppID,
		SearchTerm: q.SearchTerm, SearchCount: q.SearchCount,
		LastSearched: q.LastSearched, IsPromoted: q.IsPromoted,
		Reference: q.Reference, Status: q.Status,
		CreatedAt: q.CreatedAt, UpdatedAt: q.UpdatedAt,
	}
}
