package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type ThemeResponse struct {
	ID             uuid.UUID `json:"id"`
	VenueID        uuid.UUID `json:"venue_id"`
	Name           string    `json:"name"`
	PrimaryColor   string    `json:"primary_color"`
	SecondaryColor string    `json:"secondary_color"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ThemeRequest struct {
	Name           string `json:"name" binding:"required"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
}

func ThemeToResponse(t *domain.Theme) ThemeResponse {
	return ThemeResponse{
		ID:             t.ID,
		VenueID:        t.VenueID,
		Name:           t.Name,
		PrimaryColor:   t.PrimaryColor,
		SecondaryColor: t.SecondaryColor,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
