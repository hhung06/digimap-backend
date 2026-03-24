package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

type ProductCategoryResponse struct {
	ID           uuid.UUID       `json:"id"`
	VenueID      uuid.UUID       `json:"venue_id"`
	ExternalID   string          `json:"external_id,omitempty"`
	Name         string          `json:"name"`
	Source       string          `json:"source"`
	Localization json.RawMessage `json:"localization,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ProductCategoryRequest struct {
	ExternalID   string          `json:"external_id"`
	Name         string          `json:"name" binding:"required"`
	Source       string          `json:"source"`
	Localization json.RawMessage `json:"localization"`
}

type ProductAttachmentResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Title     string    `json:"title,omitempty"`
	FileType  string    `json:"file_type"`
	File      string    `json:"file,omitempty"`
	SourceURL string    `json:"source_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductAttachmentRequest struct {
	Title     string `json:"title"`
	File      string `json:"file"`
	SourceURL string `json:"source_url"`
}

type ProductResponse struct {
	ID              uuid.UUID                   `json:"id"`
	VenueID         uuid.UUID                   `json:"venue_id"`
	LocationID      *uuid.UUID                  `json:"location_id,omitempty"`
	MainCategoryID  *uuid.UUID                  `json:"main_category_id,omitempty"`
	Image           string                      `json:"image,omitempty"`
	Name            string                      `json:"name,omitempty"`
	Code            string                      `json:"code,omitempty"`
	Size            string                      `json:"size,omitempty"`
	Price           string                      `json:"price,omitempty"`
	OriginCountry   string                      `json:"origin_country,omitempty"`
	Expiration      string                      `json:"expiration,omitempty"`
	Description     string                      `json:"description,omitempty"`
	Custom          json.RawMessage             `json:"custom,omitempty"`
	Localization    json.RawMessage             `json:"localization,omitempty"`
	Source          string                      `json:"source"`
	Categories      []ProductCategoryResponse   `json:"categories,omitempty"`
	Attachments     []ProductAttachmentResponse `json:"attachments,omitempty"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
}

type CreateProductRequest struct {
	LocationID     *uuid.UUID      `json:"location_id"`
	MainCategoryID *uuid.UUID      `json:"main_category_id"`
	CategoryIDs    []uuid.UUID     `json:"category_ids"`
	Image          string          `json:"image"`
	Name           string          `json:"name"`
	Code           string          `json:"code"`
	Size           string          `json:"size"`
	Price          string          `json:"price"`
	OriginCountry  string          `json:"origin_country"`
	Expiration     string          `json:"expiration"`
	Description    string          `json:"description"`
	Custom         json.RawMessage `json:"custom"`
	Localization   json.RawMessage `json:"localization"`
	Source         string          `json:"source"`
}

type UpdateProductRequest struct {
	LocationID     *uuid.UUID      `json:"location_id"`
	MainCategoryID *uuid.UUID      `json:"main_category_id"`
	CategoryIDs    []uuid.UUID     `json:"category_ids"`
	Image          string          `json:"image"`
	Name           string          `json:"name"`
	Code           string          `json:"code"`
	Size           string          `json:"size"`
	Price          string          `json:"price"`
	OriginCountry  string          `json:"origin_country"`
	Expiration     string          `json:"expiration"`
	Description    string          `json:"description"`
	Custom         json.RawMessage `json:"custom"`
	Localization   json.RawMessage `json:"localization"`
	Source         string          `json:"source"`
}

// Storage
type PresignUploadRequest struct {
	Key         string `json:"key"          binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type PresignUploadResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func ProductCategoryToResponse(c *domain.ProductCategory) ProductCategoryResponse {
	return ProductCategoryResponse{
		ID: c.ID, VenueID: c.VenueID, ExternalID: c.ExternalID,
		Name: c.Name, Source: c.Source, Localization: c.Localization,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

func ProductToResponse(p *domain.Product) ProductResponse {
	r := ProductResponse{
		ID: p.ID, VenueID: p.VenueID, LocationID: p.LocationID,
		MainCategoryID: p.MainCategoryID, Image: p.Image, Name: p.Name,
		Code: p.Code, Size: p.Size, Price: p.Price,
		OriginCountry: p.OriginCountry, Expiration: p.Expiration,
		Description: p.Description, Custom: p.Custom,
		Localization: p.Localization, Source: p.Source,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
	for _, c := range p.Categories {
		r.Categories = append(r.Categories, ProductCategoryToResponse(c))
	}
	for _, a := range p.Attachments {
		r.Attachments = append(r.Attachments, ProductAttachmentResponse{
			ID: a.ID, ProductID: a.ProductID, Title: a.Title,
			FileType: a.FileType, File: a.File, SourceURL: a.SourceURL,
			CreatedAt: a.CreatedAt,
		})
	}
	return r
}
