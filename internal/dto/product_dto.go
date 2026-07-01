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
	ExternalID   *string         `json:"external_id,omitempty"`
	Name         string          `json:"name"`
	Source       string          `json:"source"`
	Localization json.RawMessage `json:"localization,omitempty" swaggertype:"object"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ProductCategoryRequest struct {
	ExternalID   *string         `json:"external_id"`
	Name         string          `json:"name" binding:"required"`
	Source       string          `json:"source"`
	Localization json.RawMessage `json:"localization" swaggertype:"object"`
}

type ProductAttachmentResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Title     *string   `json:"title,omitempty"`
	FileType  string    `json:"file_type"`
	File      *string   `json:"file,omitempty"`
	FileURL   *string   `json:"file_url,omitempty"`
	SourceURL *string   `json:"source_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductAttachmentRequest struct {
	Title     *string `json:"title"`
	FileType  string  `json:"file_type"`
	File      *string `json:"file"`
	SourceURL *string `json:"source_url"`
}

type ProductResponse struct {
	ID                  uuid.UUID                   `json:"id"`
	VenueID             uuid.UUID                   `json:"venue_id"`
	LocationID          *uuid.UUID                  `json:"location_id,omitempty"`
	LocationName        *string                     `json:"location_name,omitempty"`
	MainCategoryID      *uuid.UUID                  `json:"main_category_id,omitempty"`
	MainCategory        *ProductCategoryResponse    `json:"main_category,omitempty"`
	MainCategoryName    *string                     `json:"main_category_name,omitempty"`
	ExternalID          *string                     `json:"external_id,omitempty"`
	Image               *string                     `json:"image,omitempty"`
	Name                *string                     `json:"name,omitempty"`
	Size                *string                     `json:"size,omitempty"`
	Price               *string                     `json:"price,omitempty"`
	OriginCountry       *string                     `json:"origin_country,omitempty"`
	Expiration          *string                     `json:"expiration,omitempty"`
	Description         *string                     `json:"description,omitempty"`
	Custom              json.RawMessage             `json:"custom,omitempty" swaggertype:"object"`
	Localization        json.RawMessage             `json:"localization,omitempty" swaggertype:"object"`
	Source              string                      `json:"source"`
	Categories          []ProductCategoryResponse   `json:"categories,omitempty"`
	Attachments         []ProductAttachmentResponse `json:"attachments,omitempty"`
	ImageAttachments    []ProductAttachmentResponse `json:"image_attachments,omitempty"`
	DocumentAttachments []ProductAttachmentResponse `json:"document_attachments,omitempty"`
	AttachmentImage     *ProductAttachmentResponse  `json:"attachment_image,omitempty"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
}

type CreateProductRequest struct {
	LocationID        *uuid.UUID                 `json:"location_id"`
	MainCategoryID    *uuid.UUID                 `json:"main_category_id"`
	CategoryIDs       []uuid.UUID                `json:"category_ids"`
	ExternalID        *string                    `json:"external_id"`
	Image             *string                    `json:"image"`
	Name              *string                    `json:"name"`
	Size              *string                    `json:"size"`
	Price             *string                    `json:"price"`
	OriginCountry     *string                    `json:"origin_country"`
	Expiration        *string                    `json:"expiration"`
	Description       *string                    `json:"description"`
	Custom            json.RawMessage            `json:"custom" swaggertype:"object"`
	Localization      json.RawMessage            `json:"localization" swaggertype:"object"`
	Source            string                     `json:"source"`
	Attachments       []ProductAttachmentRequest `json:"attachments"`
	KeepAttachmentIDs []uuid.UUID                `json:"keep_attachment_ids"`
}

type UpdateProductRequest struct {
	LocationID        *uuid.UUID                 `json:"location_id"`
	MainCategoryID    *uuid.UUID                 `json:"main_category_id"`
	CategoryIDs       []uuid.UUID                `json:"category_ids"`
	Image             *string                    `json:"image"`
	Name              *string                    `json:"name"`
	ExternalID        *string                    `json:"external_id"`
	Size              *string                    `json:"size"`
	Price             *string                    `json:"price"`
	OriginCountry     *string                    `json:"origin_country"`
	Expiration        *string                    `json:"expiration"`
	Description       *string                    `json:"description"`
	Custom            json.RawMessage            `json:"custom" swaggertype:"object"`
	Localization      json.RawMessage            `json:"localization" swaggertype:"object"`
	Source            *string                    `json:"source"`
	Attachments       []ProductAttachmentRequest `json:"attachments"`
	KeepAttachmentIDs []uuid.UUID                `json:"keep_attachment_ids"`
}

func (r UpdateProductRequest) ApplyTo(p *domain.Product) {
	if r.LocationID != nil {
		p.LocationID = r.LocationID
	}
	if r.MainCategoryID != nil {
		p.MainCategoryID = r.MainCategoryID
	}
	if r.Image != nil {
		p.Image = r.Image
	}
	if r.Name != nil {
		p.Name = r.Name
	}
	if r.ExternalID != nil {
		p.ExternalID = r.ExternalID
	}
	if r.Size != nil {
		p.Size = r.Size
	}
	if r.Price != nil {
		p.Price = r.Price
	}
	if r.OriginCountry != nil {
		p.Country = r.OriginCountry
	}
	if r.Expiration != nil {
		p.Expiration = r.Expiration
	}
	if r.Description != nil {
		p.Description = r.Description
	}
	if r.Custom != nil {
		p.Custom = r.Custom
	}
	if r.Localization != nil {
		p.Localization = r.Localization
	}
	if r.Source != nil {
		p.Source = *r.Source
	}
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
		LocationName:   p.LocationName,
		MainCategoryID: p.MainCategoryID, Image: p.Image, Name: p.Name,
		ExternalID: p.ExternalID, Size: p.Size, Price: p.Price,
		OriginCountry: p.Country, Expiration: p.Expiration,
		Description: p.Description, Custom: p.Custom,
		Localization: p.Localization, Source: p.Source,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
	if p.MainCategory != nil {
		category := ProductCategoryToResponse(p.MainCategory)
		r.MainCategory = &category
		r.MainCategoryName = &p.MainCategory.Name
	}
	for _, c := range p.Categories {
		r.Categories = append(r.Categories, ProductCategoryToResponse(c))
	}
	for _, a := range p.Attachments {
		attachment := ProductAttachmentResponse{
			ID: a.ID, ProductID: a.ProductID, Title: a.Title,
			FileType: a.FileType, File: a.File, SourceURL: a.SourceURL,
			CreatedAt: a.CreatedAt,
		}
		r.Attachments = append(r.Attachments, attachment)
		if a.FileType == "image" {
			r.ImageAttachments = append(r.ImageAttachments, attachment)
		} else {
			r.DocumentAttachments = append(r.DocumentAttachments, attachment)
		}
		if r.AttachmentImage == nil && a.FileType == "image" {
			image := attachment
			r.AttachmentImage = &image
		}
	}
	return r
}
