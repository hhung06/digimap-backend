package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type LanguageResponse struct {
	ID        uuid.UUID  `json:"id"`
	VenueID   uuid.UUID  `json:"venue_id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	IsDefault bool       `json:"is_default"`
	Enabled   bool       `json:"enabled"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type LanguageRequest struct {
	Code      string `json:"code" binding:"required,max=10"`
	Name      string `json:"name" binding:"required,max=100"`
	IsDefault bool   `json:"is_default"`
	Enabled   *bool  `json:"enabled"`
}

func LanguageToResponse(l *domain.Language) LanguageResponse {
	return LanguageResponse{
		ID: l.ID, VenueID: l.VenueID, Code: l.Code,
		Name: l.Name, IsDefault: l.IsDefault, Enabled: l.Enabled,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}
