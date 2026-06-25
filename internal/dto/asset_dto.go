package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type AssetResponse struct {
	ID           uuid.UUID  `json:"id"`
	VenueID      *uuid.UUID `json:"venue_id,omitempty"`
	Name         string     `json:"name"`
	ContentType  string     `json:"content_type"`
	SizeBytes    int64      `json:"size_bytes"`
	URL          string     `json:"url"`
	AssetType    string     `json:"asset_type"`
	Description  string     `json:"description"`
	FileType     string     `json:"file_type"`
	ThumbnailURL string     `json:"thumbnail"`
	MaterialURL  string     `json:"material"`
	Width        *float64   `json:"width,omitempty"`
	Height       *float64   `json:"height,omitempty"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateAssetRequest struct {
	Name        string `json:"name"         binding:"required"`
	Key         string `json:"key"          binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	SizeBytes   int64  `json:"size_bytes"   binding:"required"`
	URL         string `json:"url"          binding:"required"`
}

type UpdateAssetRequest struct {
	Name string `json:"name" binding:"required"`
}

// Upload3DAssetRequest is parsed from multipart/form-data for 3D asset uploads.
type Upload3DAssetRequest struct {
	ID          string   `form:"id"`
	Name        string   `form:"name"`
	Description string   `form:"description"`
	FileType    string   `form:"file_type"`
	Width       *float64 `form:"width"`
	Height      *float64 `form:"height"`
}

// UploadLibraryAssetRequest is parsed from multipart/form-data for library asset uploads.
type UploadLibraryAssetRequest struct {
	Status    string `form:"status"`
	AssetType string `form:"type"`
}

func AssetToResponse(a *domain.Asset, thumbnailURL, materialURL string) AssetResponse {
	return AssetResponse{
		ID:           a.ID,
		VenueID:      a.VenueID,
		Name:         a.Name,
		ContentType:  a.ContentType,
		SizeBytes:    a.SizeBytes,
		URL:          a.URL,
		AssetType:    a.AssetType,
		Description:  a.Description,
		FileType:     a.FileType,
		ThumbnailURL: thumbnailURL,
		MaterialURL:  materialURL,
		Width:        a.Width,
		Height:       a.Height,
		CreatedBy:    a.CreatedBy,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
