package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type ProductPlazaResponse struct {
	ID          uuid.UUID  `json:"id"`
	VenueID     uuid.UUID  `json:"venue_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	LocationID  *uuid.UUID `json:"location_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ProductPlazaRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	LocationID  *uuid.UUID `json:"location_id"`
}

func ProductPlazaToResponse(p *domain.ProductPlaza) ProductPlazaResponse {
	return ProductPlazaResponse{
		ID:          p.ID,
		VenueID:     p.VenueID,
		Name:        p.Name,
		Description: p.Description,
		LocationID:  p.LocationID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
