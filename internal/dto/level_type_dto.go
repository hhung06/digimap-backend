package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type LevelTypeResponse struct {
	ID        uuid.UUID `json:"id"`
	VenueID   uuid.UUID `json:"venue_id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LevelTypeRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon"`
}

func LevelTypeToResponse(lt *domain.LevelType) LevelTypeResponse {
	return LevelTypeResponse{
		ID:        lt.ID,
		VenueID:   lt.VenueID,
		Name:      lt.Name,
		Icon:      lt.Icon,
		CreatedAt: lt.CreatedAt,
		UpdatedAt: lt.UpdatedAt,
	}
}
