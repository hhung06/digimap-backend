package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type ArticleResponse struct {
	ID                   uuid.UUID              `json:"id"`
	VenueID              *uuid.UUID             `json:"venue_id,omitempty"`
	ExternalID           string                 `json:"external_id,omitempty"`
	LocationID           *uuid.UUID             `json:"location_id,omitempty"`
	Placement            string                 `json:"placement"`
	Navigate             string                 `json:"navigate,omitempty"`
	Title                string                 `json:"title"`
	Label                string                 `json:"label,omitempty"`
	Content              string                 `json:"content,omitempty"`
	Status               string                 `json:"status"`
	PublishedAt          *time.Time             `json:"published_at,omitempty"`
	PublishedPeriodStart *time.Time             `json:"published_period_start,omitempty"`
	PublishedPeriodEnd   *time.Time             `json:"published_period_end,omitempty"`
	Localization         json.RawMessage        `json:"localization,omitempty" swaggertype:"object"`
	Images               []ArticleImageResponse `json:"images,omitempty"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

type ArticleImageResponse struct {
	ID        uuid.UUID `json:"id"`
	ArticleID uuid.UUID `json:"article_id"`
	Image     string    `json:"image"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleRequest struct {
	ExternalID           string          `json:"external_id"`
	LocationID           *uuid.UUID      `json:"location_id"`
	Placement            string          `json:"placement"`
	Navigate             string          `json:"navigate"`
	Title                string          `json:"title" binding:"required"`
	Label                string          `json:"label"`
	Content              string          `json:"content"`
	Status               string          `json:"status"`
	PublishedAt          *time.Time      `json:"published_at"`
	PublishedPeriodStart *time.Time      `json:"published_period_start"`
	PublishedPeriodEnd   *time.Time      `json:"published_period_end"`
	Localization         json.RawMessage `json:"localization" swaggertype:"object"`
}

type UpdateArticleRequest struct {
	ExternalID           *string         `json:"external_id"`
	LocationID           *uuid.UUID      `json:"location_id"`
	Placement            *string         `json:"placement"`
	Navigate             *string         `json:"navigate"`
	Title                *string         `json:"title"`
	Label                *string         `json:"label"`
	Content              *string         `json:"content"`
	Status               *string         `json:"status"`
	PublishedAt          *time.Time      `json:"published_at"`
	PublishedPeriodStart *time.Time      `json:"published_period_start"`
	PublishedPeriodEnd   *time.Time      `json:"published_period_end"`
	Localization         json.RawMessage `json:"localization" swaggertype:"object"`
}

func (r UpdateArticleRequest) ApplyTo(a *domain.Article) {
	if r.ExternalID != nil {
		a.ExternalID = *r.ExternalID
	}
	if r.LocationID != nil {
		a.LocationID = r.LocationID
	}
	if r.Placement != nil {
		a.Placement = *r.Placement
	}
	if r.Navigate != nil {
		a.Navigate = *r.Navigate
	}
	if r.Title != nil {
		a.Title = *r.Title
	}
	if r.Label != nil {
		a.Label = *r.Label
	}
	if r.Content != nil {
		a.Content = *r.Content
	}
	if r.Status != nil {
		a.Status = *r.Status
	}
	if r.PublishedAt != nil {
		a.PublishedAt = r.PublishedAt
	}
	if r.PublishedPeriodStart != nil {
		a.PublishedPeriodStart = r.PublishedPeriodStart
	}
	if r.PublishedPeriodEnd != nil {
		a.PublishedPeriodEnd = r.PublishedPeriodEnd
	}
	if r.Localization != nil {
		a.Localization = r.Localization
	}
}

type ArticleImageRequest struct {
	Image     string `json:"image" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

func ArticleToResponse(a *domain.Article) ArticleResponse {
	resp := ArticleResponse{
		ID: a.ID, VenueID: a.VenueID, ExternalID: a.ExternalID, LocationID: a.LocationID,
		Placement: a.Placement, Navigate: a.Navigate,
		Title: a.Title, Label: a.Label, Content: a.Content, Status: a.Status,
		PublishedAt:          a.PublishedAt,
		PublishedPeriodStart: a.PublishedPeriodStart,
		PublishedPeriodEnd:   a.PublishedPeriodEnd,
		Localization:         a.Localization,
		CreatedAt:            a.CreatedAt, UpdatedAt: a.UpdatedAt,
	}
	for _, img := range a.Images {
		resp.Images = append(resp.Images, ArticleImageToResponse(img))
	}
	return resp
}

func ArticleImageToResponse(img *domain.ArticleImage) ArticleImageResponse {
	return ArticleImageResponse{
		ID: img.ID, ArticleID: img.ArticleID,
		Image: img.Image, SortOrder: img.SortOrder,
		CreatedAt: img.CreatedAt, UpdatedAt: img.UpdatedAt,
	}
}
