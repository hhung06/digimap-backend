package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type TagResponse struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Localization json.RawMessage `json:"localization,omitempty" swaggertype:"object"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type TagRequest struct {
	Name         string          `json:"name" binding:"required"`
	Localization json.RawMessage `json:"localization" swaggertype:"object"`
}

type UpdateTagRequest struct {
	Name         *string         `json:"name"`
	Localization json.RawMessage `json:"localization" swaggertype:"object"`
}

func (r UpdateTagRequest) ApplyTo(t *domain.Tag) {
	if r.Name != nil {
		t.Name = *r.Name
	}
	if r.Localization != nil {
		t.Localization = r.Localization
	}
}

type AttachTagRequest struct {
	TagID      uuid.UUID `json:"tag_id" binding:"required"`
	EntityType string    `json:"entity_type" binding:"required"`
	EntityID   uuid.UUID `json:"entity_id" binding:"required"`
}

func TagToResponse(t *domain.Tag) TagResponse {
	return TagResponse{
		ID:           t.ID,
		Name:         t.Name,
		Localization: t.Localization,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
