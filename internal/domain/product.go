package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProductCategory struct {
	ID           uuid.UUID
	VenueID      uuid.UUID
	ExternalID   *string
	Name         string
	Source       string
	Localization json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type Product struct {
	ID             uuid.UUID
	VenueID        uuid.UUID
	LocationID     *uuid.UUID
	MainCategoryID *uuid.UUID
	ExternalID     *string
	Image          *string
	Name           *string
	Size           *string
	Price          *string
	Country        *string
	Expiration     *string
	Description    *string
	Custom         json.RawMessage
	Localization   json.RawMessage
	Source         string
	// Eagerly loaded
	Categories  []*ProductCategory
	Attachments []*ProductAttachment
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type ProductAttachment struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Title     *string
	FileType  string
	File      *string
	SourceURL *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
