package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type VideoResponse struct {
	ID          uuid.UUID  `json:"id"`
	VenueID     *uuid.UUID `json:"venue_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	URL         string     `json:"url,omitempty"`
	Thumbnail   string     `json:"thumbnail,omitempty"`
	Duration    *int       `json:"duration,omitempty"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type VideoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	Thumbnail   string     `json:"thumbnail"`
	Duration    *int       `json:"duration"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
}

func VideoToResponse(v *domain.Video) VideoResponse {
	return VideoResponse{
		ID: v.ID, VenueID: v.VenueID,
		Title: v.Title, Description: v.Description,
		URL: v.URL, Thumbnail: v.Thumbnail,
		Duration: v.Duration, Status: v.Status, PublishedAt: v.PublishedAt,
		CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}
