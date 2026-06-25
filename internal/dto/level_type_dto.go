package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type LevelTypeResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LevelTypeRequest struct {
	Name  string `json:"name" binding:"required"`
	Icon  string `json:"icon"`
	Value int    `json:"value"`
}

func LevelTypeToResponse(lt *domain.LevelType) LevelTypeResponse {
	return LevelTypeResponse{
		ID:        lt.ID,
		Name:      lt.Name,
		Icon:      lt.Icon,
		Value:     lt.Value,
		CreatedAt: lt.CreatedAt,
		UpdatedAt: lt.UpdatedAt,
	}
}
