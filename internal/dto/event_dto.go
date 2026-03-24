package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

// ── Event tags ────────────────────────────────────────────────────────────────

type EventTagResponse struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Localization json.RawMessage `json:"localization,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type EventTagRequest struct {
	Name         string          `json:"name" binding:"required"`
	Localization json.RawMessage `json:"localization"`
}

func EventTagToResponse(t *domain.EventTag) EventTagResponse {
	return EventTagResponse{
		ID: t.ID, Name: t.Name, Localization: t.Localization,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

// ── Event types ───────────────────────────────────────────────────────────────

type EventTypeResponse struct {
	ID           uuid.UUID       `json:"id"`
	VenueID      *uuid.UUID      `json:"venue_id,omitempty"`
	Name         string          `json:"name"`
	Localization json.RawMessage `json:"localization,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type EventTypeRequest struct {
	Name         string          `json:"name" binding:"required"`
	Localization json.RawMessage `json:"localization"`
}

func EventTypeToResponse(t *domain.EventType) EventTypeResponse {
	return EventTypeResponse{
		ID: t.ID, VenueID: t.VenueID, Name: t.Name,
		Localization: t.Localization, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

// ── Event images ──────────────────────────────────────────────────────────────

type EventImageResponse struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	Image     string    `json:"image,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type EventImageRequest struct {
	Image string `json:"image" binding:"required"`
}

// ── Events ────────────────────────────────────────────────────────────────────

type EventResponse struct {
	ID            uuid.UUID           `json:"id"`
	VenueID       uuid.UUID           `json:"venue_id"`
	TypeID        *uuid.UUID          `json:"type_id,omitempty"`
	Title         string              `json:"title,omitempty"`
	Description   string              `json:"description,omitempty"`
	BannerImage   string              `json:"banner_image,omitempty"`
	IconImage     string              `json:"icon_image,omitempty"`
	StartTime     *time.Time          `json:"start_time,omitempty"`
	EndTime       *time.Time          `json:"end_time,omitempty"`
	ShowStartTime *time.Time          `json:"show_start_time,omitempty"`
	ShowEndTime   *time.Time          `json:"show_end_time,omitempty"`
	ContentDetail string              `json:"content_detail,omitempty"`
	ContentURL    string              `json:"content_url,omitempty"`
	Localization  json.RawMessage     `json:"localization,omitempty"`
	Tags          []EventTagResponse  `json:"tags,omitempty"`
	Locations     []uuid.UUID         `json:"locations,omitempty"`
	Images        []EventImageResponse `json:"images,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type CreateEventRequest struct {
	TypeID        *uuid.UUID      `json:"type_id"`
	TagIDs        []uuid.UUID     `json:"tag_ids"`
	LocationIDs   []uuid.UUID     `json:"location_ids"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	BannerImage   string          `json:"banner_image"`
	IconImage     string          `json:"icon_image"`
	StartTime     *time.Time      `json:"start_time"`
	EndTime       *time.Time      `json:"end_time"`
	ShowStartTime *time.Time      `json:"show_start_time"`
	ShowEndTime   *time.Time      `json:"show_end_time"`
	ContentDetail string          `json:"content_detail"`
	ContentURL    string          `json:"content_url"`
	Localization  json.RawMessage `json:"localization"`
}

type UpdateEventRequest struct {
	TypeID        *uuid.UUID      `json:"type_id"`
	TagIDs        []uuid.UUID     `json:"tag_ids"`
	LocationIDs   []uuid.UUID     `json:"location_ids"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	BannerImage   string          `json:"banner_image"`
	IconImage     string          `json:"icon_image"`
	StartTime     *time.Time      `json:"start_time"`
	EndTime       *time.Time      `json:"end_time"`
	ShowStartTime *time.Time      `json:"show_start_time"`
	ShowEndTime   *time.Time      `json:"show_end_time"`
	ContentDetail string          `json:"content_detail"`
	ContentURL    string          `json:"content_url"`
	Localization  json.RawMessage `json:"localization"`
}

func EventToResponse(e *domain.Event) EventResponse {
	r := EventResponse{
		ID: e.ID, VenueID: e.VenueID, TypeID: e.TypeID,
		Title: e.Title, Description: e.Description,
		BannerImage: e.BannerImage, IconImage: e.IconImage,
		StartTime: e.StartTime, EndTime: e.EndTime,
		ShowStartTime: e.ShowStartTime, ShowEndTime: e.ShowEndTime,
		ContentDetail: e.ContentDetail, ContentURL: e.ContentURL,
		Localization: e.Localization, Locations: e.Locations,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
	for _, t := range e.Tags {
		r.Tags = append(r.Tags, EventTagToResponse(t))
	}
	for _, img := range e.Images {
		r.Images = append(r.Images, EventImageResponse{
			ID: img.ID, EventID: img.EventID, Image: img.Image, CreatedAt: img.CreatedAt,
		})
	}
	return r
}
