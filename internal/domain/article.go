package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID                   uuid.UUID
	VenueID              *uuid.UUID
	ExternalID           *string
	LocationID           *uuid.UUID
	Location             *ArticleLocation
	Placement            string
	Navigate             *string
	Title                string
	Label                *string
	Content              *string
	Status               string
	PublishedAt          *time.Time
	PublishedPeriodStart *time.Time
	PublishedPeriodEnd   *time.Time
	Localization         json.RawMessage
	RelatedProducts      []uuid.UUID
	Images               []*ArticleImage
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ArticleLocation struct {
	ID   uuid.UUID
	Name string
}

type ArticleImage struct {
	ID        uuid.UUID
	ArticleID uuid.UUID
	Image     string
	SortOrder int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ArticleMediaChange struct {
	Replace bool
	Keys    []string
}
