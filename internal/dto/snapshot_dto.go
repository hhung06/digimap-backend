package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type SnapshotResponse struct {
	ID        uuid.UUID  `json:"id"`
	VenueID   uuid.UUID  `json:"venue_id"`
	State     int        `json:"state"`
	Method    int        `json:"method"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	PublishAt *time.Time `json:"publish_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateSnapshotRequest struct {
	Bundle json.RawMessage `json:"bundle" binding:"required"`
}

func SnapshotToResponse(s *domain.Snapshot) SnapshotResponse {
	return SnapshotResponse{
		ID:        s.ID,
		VenueID:   s.VenueID,
		State:     s.State,
		Method:    s.Method,
		CreatedBy: s.CreatedBy,
		PublishAt: s.PublishAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
