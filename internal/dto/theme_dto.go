package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type ThemeResponse struct {
	ID          uuid.UUID       `json:"id"`
	VenueID     *uuid.UUID      `json:"venue_id"`
	Scope       string          `json:"scope"`
	Name        string          `json:"name"`
	Data        json.RawMessage `json:"data"`
	StoragePath string          `json:"storage_path"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ThemeRequest struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type SetThemeRequest struct {
	ThemeID *uuid.UUID `json:"theme_id"`
}

func ThemeToResponse(t *domain.Theme) ThemeResponse {
	data := t.Data
	if len(data) == 0 {
		data = json.RawMessage("{}")
	}
	return ThemeResponse{
		ID:          t.ID,
		VenueID:     t.VenueID,
		Scope:       t.Scope,
		Name:        t.Name,
		Data:        data,
		StoragePath: t.StoragePath,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
