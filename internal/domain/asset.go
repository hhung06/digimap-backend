package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	AssetType2D = "2d"
	AssetType3D = "3d"

	AssetStatusPublished   = "published"
	AssetStatusUnpublished = "unpublished"
)

// Asset represents a media file tracked in the database (metadata only — actual file is in S3).
type Asset struct {
	ID          uuid.UUID
	VenueID     *uuid.UUID // nullable — nil for global (library) assets
	Name        string
	Key         string
	ContentType string
	SizeBytes   int64
	URL         string
	AssetType   string
	Description string
	FileType    string
	Thumbnail   string   // S3 key for thumbnail image (3D assets)
	Material    string   // S3 key for material file (3D assets)
	Width       *float64
	Height      *float64
	Status      string   // "published" | "unpublished"
	CreatedBy   *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
