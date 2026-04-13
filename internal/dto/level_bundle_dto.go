package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type LevelBundleResponse struct {
	ID         uuid.UUID `json:"id"`
	SnapshotID uuid.UUID `json:"snapshot_id"`
	VenueID    uuid.UUID `json:"venue_id"`
	LevelID    uuid.UUID `json:"level_id"`
	State      int       `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateLevelBundleRequest struct {
	LevelID uuid.UUID `json:"level_id" binding:"required"`
	Bundle  json.RawMessage `json:"bundle" binding:"required"`
}

func LevelBundleToResponse(b *domain.LevelBundle) LevelBundleResponse {
	return LevelBundleResponse{
		ID:         b.ID,
		SnapshotID: b.SnapshotID,
		VenueID:    b.VenueID,
		LevelID:    b.LevelID,
		State:      b.State,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}
}
