package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID                  uuid.UUID
	VenueID             *uuid.UUID
	ExternalID          string
	LocationID          *uuid.UUID
	Placement           string
	Navigate            string
	Title               string
	Label               string
	Content             string
	Status              string
	PublishedAt         *time.Time
	PublishedPeriodStart *time.Time
	PublishedPeriodEnd   *time.Time
	Localization        json.RawMessage
	Images              []*ArticleImage
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ArticleImage struct {
	ID        uuid.UUID
	ArticleID uuid.UUID
	Image     string
	SortOrder int
	CreatedAt time.Time
	UpdatedAt time.Time
}
