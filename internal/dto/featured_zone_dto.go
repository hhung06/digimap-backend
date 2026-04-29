package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type FeaturedZoneResponse struct {
	ID          uuid.UUID  `json:"id"`
	VenueID     uuid.UUID  `json:"venue_id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ImageURL    string     `json:"image_url,omitempty"`
	SortIndex   int        `json:"sort_index"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type FeaturedZoneRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	SortIndex   int    `json:"sort_index"`
	IsActive    *bool  `json:"is_active"`
}

func FeaturedZoneToResponse(z *domain.FeaturedZone) FeaturedZoneResponse {
	return FeaturedZoneResponse{
		ID: z.ID, VenueID: z.VenueID, Name: z.Name,
		Description: z.Description, ImageURL: z.ImageURL,
		SortIndex: z.SortIndex, IsActive: z.IsActive,
		CreatedAt: z.CreatedAt, UpdatedAt: z.UpdatedAt,
	}
}
