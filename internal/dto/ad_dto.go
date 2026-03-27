package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type AdvertisementResponse struct {
	ID              uuid.UUID  `json:"id"`
	VenueID         *uuid.UUID `json:"venue_id,omitempty"`
	LocationID      *uuid.UUID `json:"location_id,omitempty"`
	Type            string     `json:"type"`
	Status          string     `json:"status"`
	Navigate        string     `json:"navigate,omitempty"`
	ContentImageURL string     `json:"content_image_url,omitempty"`
	ContentCTAURL   string     `json:"content_cta_url,omitempty"`
	Placement       string     `json:"placement"`
	SizeWidth       *int       `json:"size_width,omitempty"`
	SizeHeight      *int       `json:"size_height,omitempty"`
	RewardType      string     `json:"reward_type,omitempty"`
	RewardAmount    *int       `json:"reward_amount,omitempty"`
	DisplayDuration *int       `json:"display_duration,omitempty"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
	StartAt         *time.Time `json:"start_at,omitempty"`
	EndAt           *time.Time `json:"end_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type AdvertisementRequest struct {
	LocationID      *uuid.UUID `json:"location_id"`
	Type            string     `json:"type"`
	Navigate        string     `json:"navigate"`
	ContentImageURL string     `json:"content_image_url"`
	ContentCTAURL   string     `json:"content_cta_url"`
	Placement       string     `json:"placement"`
	SizeWidth       *int       `json:"size_width"`
	SizeHeight      *int       `json:"size_height"`
	RewardType      string     `json:"reward_type"`
	RewardAmount    *int       `json:"reward_amount"`
	DisplayDuration *int       `json:"display_duration"`
	StartAt         *time.Time `json:"start_at"`
	EndAt           *time.Time `json:"end_at"`
}

func AdvertisementToResponse(a *domain.Advertisement) AdvertisementResponse {
	return AdvertisementResponse{
		ID: a.ID, VenueID: a.VenueID, LocationID: a.LocationID,
		Type: a.Type, Status: a.Status, Navigate: a.Navigate,
		ContentImageURL: a.ContentImageURL, ContentCTAURL: a.ContentCTAURL,
		Placement: a.Placement,
		SizeWidth: a.SizeWidth, SizeHeight: a.SizeHeight,
		RewardType: a.RewardType, RewardAmount: a.RewardAmount,
		DisplayDuration: a.DisplayDuration,
		PublishedAt: a.PublishedAt, StartAt: a.StartAt, EndAt: a.EndAt,
		CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
	}
}
