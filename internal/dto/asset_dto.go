package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type AssetResponse struct {
	ID          uuid.UUID  `json:"id"`
	VenueID     uuid.UUID  `json:"venue_id"`
	Name        string     `json:"name"`
	ContentType string     `json:"content_type"`
	SizeBytes   int64      `json:"size_bytes"`
	URL         string     `json:"url"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateAssetRequest struct {
	Name        string `json:"name"         binding:"required"`
	Key         string `json:"key"          binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	SizeBytes   int64  `json:"size_bytes"   binding:"required"`
	URL         string `json:"url"          binding:"required"`
}

type UpdateAssetRequest struct {
	Name string `json:"name" binding:"required"`
}

func AssetToResponse(a *domain.Asset) AssetResponse {
	return AssetResponse{
		ID:          a.ID,
		VenueID:     a.VenueID,
		Name:        a.Name,
		ContentType: a.ContentType,
		SizeBytes:   a.SizeBytes,
		URL:         a.URL,
		CreatedBy:   a.CreatedBy,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
