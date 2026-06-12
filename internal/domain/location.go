package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ── Location types ────────────────────────────────────────────────────────────

const (
	LocationTypeDefault = 0
	LocationTypeBooth   = 2
	LocationTypeMemo    = 6
)

// ── Location category ─────────────────────────────────────────────────────────

type LocationCategory struct {
	ID            uuid.UUID
	VenueID       uuid.UUID
	ParentID      *uuid.UUID
	ExternalID    string
	Name          string
	ShortName     string
	Color         string
	Icon          string
	IconDefault   string
	SortIndex     int
	Visible       bool
	Description   string
	Type          string
	Image         string
	Localization  json.RawMessage
	Source        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	Parent        *LocationCategory
	Subcategories []*LocationCategory
}

// ── Location ──────────────────────────────────────────────────────────────────

type Location struct {
	ID                           uuid.UUID
	VenueID                      uuid.UUID
	LevelID                      *uuid.UUID
	MainCategoryID               *uuid.UUID
	ExternalID                   string
	CommonHidden                 bool
	CommonName                   string
	CommonShortName              string
	CommonDescription            string
	CommonColor                  string
	CommonLocationType           int
	CommonLocationSubType        int
	CommonLatitude               float64
	CommonLongitude              float64
	CommonAddress                string
	CommonLocationState          *int
	CommonLocationStateStartDate *time.Time
	CommonLocationStateEndDate   *time.Time
	CommonLogo                   string
	CommonLargeLogo              string
	CommonMediumLogo             string
	CommonSmallLogo              string
	CommonSocialWebsite          string
	CommonSocialTwitter          string
	CommonSocialTiktok           string
	CommonSocialFacebook         string
	CommonSocialInstagram        string
	CommonContactEmail           string
	CommonContactPhone           string
	CommonShowShortName          bool
	TopLogo                      string
	TopLogoType                  string
	PlaceWorkHours               json.RawMessage
	BoothNumber                  string
	BoothEventDate               *time.Time
	BoothSize                    string
	BoothServicesOffered         string
	BoothProductsShowcased       string
	PersonFullName               string
	PersonJobTitle               string
	RoomNumber                   string
	RoomDepartment               string
	RoomBedCount                 *int
	RoomEquipmentDetails         string
	IsTopLocation                bool
	TopLocationSortIndex         *int
	IconDefault                  string
	Custom                       json.RawMessage
	Localization                 json.RawMessage
	Source                       string
	StartTime                    *time.Time
	EndTime                      *time.Time
	IsSearchable                 bool
	// Eagerly loaded relations
	Categories []*LocationCategory
	Images     []*LocationImage
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type LocationImage struct {
	ID         uuid.UUID
	LocationID uuid.UUID
	Original   string
	Small      string
	Medium     string
	Large      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
