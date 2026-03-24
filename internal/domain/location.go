package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ── Location category ─────────────────────────────────────────────────────────

type LocationCategory struct {
	ID          uuid.UUID
	VenueID     uuid.UUID
	ExternalID  string
	Name        string
	ShortName   string
	Color       string
	Icon        string
	IconDefault string
	SortIndex   int
	Visible     bool
	Description string
	Type        string
	Image       string
	Localization json.RawMessage
	Source      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// ── Amenity (location template / master) ──────────────────────────────────────

type Amenity struct {
	ID                           uuid.UUID
	CommonName                   string
	CommonShortName              string
	CommonDescription            string
	CommonColor                  string
	CommonLocationType           int
	CommonLatitude               float64
	CommonLongitude              float64
	CommonAddress                string
	CommonLocationState          *int
	CommonLocationStateStartDate *time.Time
	CommonLocationStateEndDate   *time.Time
	CommonLogo                   string
	CommonSocialWebsite          string
	CommonSocialTwitter          string
	CommonSocialTiktok           string
	CommonSocialFacebook         string
	CommonSocialInstagram        string
	CommonContactEmail           string
	CommonContactPhone           string
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
	Localization                 json.RawMessage
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	DeletedAt                    *time.Time
}

// VenueAmenity links an amenity to a venue.
type VenueAmenity struct {
	ID        uuid.UUID
	VenueID   uuid.UUID
	AmenityID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// ── Location ──────────────────────────────────────────────────────────────────

type Location struct {
	ID                            uuid.UUID
	VenueID                       uuid.UUID
	LevelID                       *uuid.UUID
	MainCategoryID                *uuid.UUID
	ExternalID                    string
	CommonHidden                  bool
	CommonName                    string
	CommonShortName               string
	CommonDescription             string
	CommonColor                   string
	CommonLocationType            int
	CommonSubType                 int
	CommonLatitude                float64
	CommonLongitude               float64
	CommonAddress                 string
	CommonLocationState           *int
	CommonLocationStateStartDate  *time.Time
	CommonLocationStateEndDate    *time.Time
	CommonLogo                    string
	CommonLargeLogo               string
	CommonMediumLogo              string
	CommonSmallLogo               string
	CommonSocialWebsite           string
	CommonSocialTwitter           string
	CommonSocialTiktok            string
	CommonSocialFacebook          string
	CommonSocialInstagram         string
	CommonContactEmail            string
	CommonContactPhone            string
	CommonShowShortName           bool
	TopLogo                       string
	TopLogoType                   string
	PlaceWorkHours                json.RawMessage
	BoothNumber                   string
	BoothEventDate                *time.Time
	BoothSize                     string
	BoothServicesOffered          string
	BoothProductsShowcased        string
	PersonFullName                string
	PersonJobTitle                string
	RoomNumber                    string
	RoomDepartment                string
	RoomBedCount                  *int
	RoomEquipmentDetails          string
	IsTopLocation                 bool
	TopLocationSortIndex          *int
	IconDefault                   string
	Custom                        json.RawMessage
	Localization                  json.RawMessage
	Source                        string
	StartTime                     *time.Time
	EndTime                       *time.Time
	IsSearchable                  bool
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

// ── Promotion ─────────────────────────────────────────────────────────────────

type Promotion struct {
	ID                  uuid.UUID
	VenueID             uuid.UUID
	LocationID          *uuid.UUID
	ExternalID          string
	PromoImage          string
	Introduction        string
	GiftContent         string
	DetailURL           string
	BoothNumber         string
	ExpectedGiftCount   *int
	DistributionStart   *time.Time
	DistributionEnd     *time.Time
	DisplayType         string
	Localization        json.RawMessage
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}
