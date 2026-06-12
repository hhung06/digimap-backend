package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

// ── Location category ─────────────────────────────────────────────────────────

type LocationCategoryResponse struct {
	ID            uuid.UUID                         `json:"id"`
	VenueID       uuid.UUID                         `json:"venue_id"`
	ParentID      *uuid.UUID                        `json:"parent_id,omitempty"`
	Parent        *LocationCategorySummaryResponse  `json:"parent,omitempty"`
	Subcategories []LocationCategorySummaryResponse `json:"subcategories,omitempty"`
	ExternalID    string                            `json:"external_id,omitempty"`
	Name          string                            `json:"name,omitempty"`
	ShortName     string                            `json:"short_name,omitempty"`
	Color         string                            `json:"color,omitempty"`
	Icon          string                            `json:"icon,omitempty"`
	IconDefault   string                            `json:"icon_default,omitempty"`
	SortIndex     int                               `json:"sort_index"`
	Visible       bool                              `json:"visible"`
	Description   string                            `json:"description,omitempty"`
	Type          string                            `json:"type,omitempty"`
	Image         string                            `json:"image,omitempty"`
	Localization  json.RawMessage                   `json:"localization,omitempty" swaggertype:"object"`
	Source        string                            `json:"source"`
	CreatedAt     time.Time                         `json:"created_at"`
	UpdatedAt     time.Time                         `json:"updated_at"`
}

type LocationCategorySummaryResponse struct {
	ID         uuid.UUID  `json:"id"`
	VenueID    uuid.UUID  `json:"venue_id"`
	ParentID   *uuid.UUID `json:"parent_id,omitempty"`
	ExternalID string     `json:"external_id,omitempty"`
	Name       string     `json:"name,omitempty"`
	ShortName  string     `json:"short_name,omitempty"`
	Color      string     `json:"color,omitempty"`
	Icon       string     `json:"icon,omitempty"`
	SortIndex  int        `json:"sort_index"`
	Visible    bool       `json:"visible"`
	Type       string     `json:"type,omitempty"`
}

type LocationCategoryRequest struct {
	ParentID     *uuid.UUID      `json:"parent_id"`
	ExternalID   string          `json:"external_id"`
	Name         string          `json:"name"`
	ShortName    string          `json:"short_name"`
	Color        string          `json:"color"`
	Icon         string          `json:"icon"`
	IconDefault  string          `json:"icon_default"`
	SortIndex    int             `json:"sort_index"`
	Visible      bool            `json:"visible"`
	Description  string          `json:"description"`
	Type         string          `json:"type"`
	Image        string          `json:"image"`
	Localization json.RawMessage `json:"localization" swaggertype:"object"`
	Source       string          `json:"source"`
}

type UpdateLocationCategoryRequest struct {
	ParentID     *uuid.UUID      `json:"parent_id"`
	ExternalID   *string         `json:"external_id"`
	Name         *string         `json:"name"`
	ShortName    *string         `json:"short_name"`
	Color        *string         `json:"color"`
	Icon         *string         `json:"icon"`
	IconDefault  *string         `json:"icon_default"`
	SortIndex    *int            `json:"sort_index"`
	Visible      *bool           `json:"visible"`
	Description  *string         `json:"description"`
	Type         *string         `json:"type"`
	Image        *string         `json:"image"`
	Localization json.RawMessage `json:"localization" swaggertype:"object"`
	Source       *string         `json:"source"`
}

func (r UpdateLocationCategoryRequest) ApplyTo(c *domain.LocationCategory) {
	if r.ParentID != nil {
		c.ParentID = r.ParentID
	}
	if r.ExternalID != nil {
		c.ExternalID = *r.ExternalID
	}
	if r.Name != nil {
		c.Name = *r.Name
	}
	if r.ShortName != nil {
		c.ShortName = *r.ShortName
	}
	if r.Color != nil {
		c.Color = *r.Color
	}
	if r.Icon != nil {
		c.Icon = *r.Icon
	}
	if r.IconDefault != nil {
		c.IconDefault = *r.IconDefault
	}
	if r.SortIndex != nil {
		c.SortIndex = *r.SortIndex
	}
	if r.Visible != nil {
		c.Visible = *r.Visible
	}
	if r.Description != nil {
		c.Description = *r.Description
	}
	if r.Type != nil {
		c.Type = *r.Type
	}
	if r.Image != nil {
		c.Image = *r.Image
	}
	if r.Localization != nil {
		c.Localization = r.Localization
	}
	if r.Source != nil {
		c.Source = *r.Source
	}
}

func LocationCategoryToResponse(c *domain.LocationCategory) LocationCategoryResponse {
	r := LocationCategoryResponse{
		ID: c.ID, VenueID: c.VenueID, ExternalID: c.ExternalID,
		ParentID: c.ParentID,
		Name:     c.Name, ShortName: c.ShortName, Color: c.Color,
		Icon: c.Icon, IconDefault: c.IconDefault, SortIndex: c.SortIndex,
		Visible: c.Visible, Description: c.Description, Type: c.Type,
		Image: c.Image, Localization: c.Localization, Source: c.Source,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
	if c.Parent != nil {
		parent := LocationCategoryToSummaryResponse(c.Parent)
		r.Parent = &parent
	}
	for _, sub := range c.Subcategories {
		r.Subcategories = append(r.Subcategories, LocationCategoryToSummaryResponse(sub))
	}
	return r
}

func LocationCategoryToSummaryResponse(c *domain.LocationCategory) LocationCategorySummaryResponse {
	return LocationCategorySummaryResponse{
		ID: c.ID, VenueID: c.VenueID, ParentID: c.ParentID, ExternalID: c.ExternalID,
		Name: c.Name, ShortName: c.ShortName, Color: c.Color, Icon: c.Icon,
		SortIndex: c.SortIndex, Visible: c.Visible, Type: c.Type,
	}
}

// ── Amenity ───────────────────────────────────────────────────────────────────

type AmenityResponse struct {
	ID                    uuid.UUID       `json:"id"`
	CommonName            string          `json:"common_name"`
	CommonShortName       string          `json:"common_short_name,omitempty"`
	CommonDescription     string          `json:"common_description,omitempty"`
	CommonColor           string          `json:"common_color,omitempty"`
	CommonLocationType    int             `json:"common_location_type"`
	CommonLatitude        float64         `json:"common_latitude"`
	CommonLongitude       float64         `json:"common_longitude"`
	CommonAddress         string          `json:"common_address,omitempty"`
	CommonLogo            string          `json:"common_logo,omitempty"`
	CommonSocialWebsite   string          `json:"common_social_website,omitempty"`
	CommonSocialTwitter   string          `json:"common_social_twitter,omitempty"`
	CommonSocialTiktok    string          `json:"common_social_tiktok,omitempty"`
	CommonSocialFacebook  string          `json:"common_social_facebook,omitempty"`
	CommonSocialInstagram string          `json:"common_social_instagram,omitempty"`
	CommonContactEmail    string          `json:"common_contact_email,omitempty"`
	CommonContactPhone    string          `json:"common_contact_phone,omitempty"`
	PlaceWorkHours        json.RawMessage `json:"place_work_hours,omitempty" swaggertype:"object"`
	Localization          json.RawMessage `json:"localization,omitempty" swaggertype:"object"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

type AmenityRequest struct {
	CommonName            string          `json:"common_name" binding:"required"`
	CommonShortName       string          `json:"common_short_name"`
	CommonDescription     string          `json:"common_description"`
	CommonColor           string          `json:"common_color"`
	CommonLocationType    int             `json:"common_location_type"`
	CommonLatitude        float64         `json:"common_latitude"`
	CommonLongitude       float64         `json:"common_longitude"`
	CommonAddress         string          `json:"common_address"`
	CommonLogo            string          `json:"common_logo"`
	CommonSocialWebsite   string          `json:"common_social_website"`
	CommonSocialTwitter   string          `json:"common_social_twitter"`
	CommonSocialTiktok    string          `json:"common_social_tiktok"`
	CommonSocialFacebook  string          `json:"common_social_facebook"`
	CommonSocialInstagram string          `json:"common_social_instagram"`
	CommonContactEmail    string          `json:"common_contact_email"`
	CommonContactPhone    string          `json:"common_contact_phone"`
	PlaceWorkHours        json.RawMessage `json:"place_work_hours" swaggertype:"object"`
	Localization          json.RawMessage `json:"localization" swaggertype:"object"`
}

// ── Location ──────────────────────────────────────────────────────────────────

type LocationResponse struct {
	ID                    uuid.UUID                  `json:"id"`
	VenueID               uuid.UUID                  `json:"venue_id"`
	LevelID               *uuid.UUID                 `json:"level_id,omitempty"`
	MainCategoryID        *uuid.UUID                 `json:"main_category_id,omitempty"`
	ExternalID            string                     `json:"external_id,omitempty"`
	CommonHidden          bool                       `json:"common_hidden"`
	CommonName            string                     `json:"common_name"`
	CommonShortName       string                     `json:"common_short_name,omitempty"`
	CommonDescription     string                     `json:"common_description,omitempty"`
	CommonColor           string                     `json:"common_color,omitempty"`
	CommonLocationType    int                        `json:"common_location_type"`
	CommonLocationSubType int                        `json:"common_location_sub_type"`
	CommonLatitude        float64                    `json:"common_latitude"`
	CommonLongitude       float64                    `json:"common_longitude"`
	CommonAddress         string                     `json:"common_address,omitempty"`
	CommonLogo            string                     `json:"common_logo,omitempty"`
	CommonLargeLogo       string                     `json:"common_large_logo,omitempty"`
	CommonMediumLogo      string                     `json:"common_medium_logo,omitempty"`
	CommonSmallLogo       string                     `json:"common_small_logo,omitempty"`
	CommonContactEmail    string                     `json:"common_contact_email,omitempty"`
	CommonContactPhone    string                     `json:"common_contact_phone,omitempty"`
	IsTopLocation         bool                       `json:"is_top_location"`
	IsSearchable          bool                       `json:"is_searchable"`
	Source                string                     `json:"source"`
	PlaceWorkHours        json.RawMessage            `json:"place_work_hours,omitempty" swaggertype:"object"`
	Custom                json.RawMessage            `json:"custom,omitempty" swaggertype:"object"`
	Localization          json.RawMessage            `json:"localization,omitempty" swaggertype:"object"`
	StartTime             *time.Time                 `json:"start_time,omitempty"`
	EndTime               *time.Time                 `json:"end_time,omitempty"`
	Categories            []LocationCategoryResponse `json:"categories,omitempty"`
	Images                []LocationImageResponse    `json:"images,omitempty"`
	CreatedAt             time.Time                  `json:"created_at"`
	UpdatedAt             time.Time                  `json:"updated_at"`
}

type LocationImageResponse struct {
	ID         uuid.UUID `json:"id"`
	LocationID uuid.UUID `json:"location_id"`
	Original   string    `json:"original,omitempty"`
	Small      string    `json:"small,omitempty"`
	Medium     string    `json:"medium,omitempty"`
	Large      string    `json:"large,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateLocationRequest struct {
	LevelID               *uuid.UUID      `json:"level_id"`
	MainCategoryID        *uuid.UUID      `json:"main_category_id"`
	CategoryIDs           []uuid.UUID     `json:"category_ids"`
	ExternalID            string          `json:"external_id"`
	CommonHidden          bool            `json:"common_hidden"`
	CommonName            string          `json:"common_name" binding:"required"`
	CommonShortName       string          `json:"common_short_name"`
	CommonDescription     string          `json:"common_description"`
	CommonColor           string          `json:"common_color"`
	CommonLocationType    int             `json:"common_location_type"`
	CommonLocationSubType int             `json:"common_location_sub_type"`
	CommonLatitude        float64         `json:"common_latitude"`
	CommonLongitude       float64         `json:"common_longitude"`
	CommonAddress         string          `json:"common_address"`
	CommonContactEmail    string          `json:"common_contact_email"`
	CommonContactPhone    string          `json:"common_contact_phone"`
	PlaceWorkHours        json.RawMessage `json:"place_work_hours" swaggertype:"object"`
	Custom                json.RawMessage `json:"custom" swaggertype:"object"`
	Localization          json.RawMessage `json:"localization" swaggertype:"object"`
	Source                string          `json:"source"`
	StartTime             *time.Time      `json:"start_time"`
	EndTime               *time.Time      `json:"end_time"`
	IsSearchable          bool            `json:"is_searchable"`
}

type UpdateLocationRequest struct {
	LevelID               *uuid.UUID      `json:"level_id"`
	MainCategoryID        *uuid.UUID      `json:"main_category_id"`
	CategoryIDs           []uuid.UUID     `json:"category_ids"`
	ExternalID            *string         `json:"external_id"`
	CommonHidden          *bool           `json:"common_hidden"`
	CommonName            *string         `json:"common_name"`
	CommonShortName       *string         `json:"common_short_name"`
	CommonDescription     *string         `json:"common_description"`
	CommonColor           *string         `json:"common_color"`
	CommonLocationType    *int            `json:"common_location_type"`
	CommonLocationSubType *int            `json:"common_location_sub_type"`
	CommonLatitude        *float64        `json:"common_latitude"`
	CommonLongitude       *float64        `json:"common_longitude"`
	CommonAddress         *string         `json:"common_address"`
	CommonLogo            *string         `json:"common_logo"`
	CommonContactEmail    *string         `json:"common_contact_email"`
	CommonContactPhone    *string         `json:"common_contact_phone"`
	PlaceWorkHours        json.RawMessage `json:"place_work_hours" swaggertype:"object"`
	Custom                json.RawMessage `json:"custom" swaggertype:"object"`
	Localization          json.RawMessage `json:"localization" swaggertype:"object"`
	Source                *string         `json:"source"`
	StartTime             *time.Time      `json:"start_time"`
	EndTime               *time.Time      `json:"end_time"`
	IsSearchable          *bool           `json:"is_searchable"`
}

func (r UpdateLocationRequest) ApplyTo(l *domain.Location) {
	if r.LevelID != nil {
		l.LevelID = r.LevelID
	}
	if r.MainCategoryID != nil {
		l.MainCategoryID = r.MainCategoryID
	}
	if r.ExternalID != nil {
		l.ExternalID = *r.ExternalID
	}
	if r.CommonHidden != nil {
		l.CommonHidden = *r.CommonHidden
	}
	if r.CommonName != nil {
		l.CommonName = *r.CommonName
	}
	if r.CommonShortName != nil {
		l.CommonShortName = *r.CommonShortName
	}
	if r.CommonDescription != nil {
		l.CommonDescription = *r.CommonDescription
	}
	if r.CommonColor != nil {
		l.CommonColor = *r.CommonColor
	}
	if r.CommonLocationType != nil {
		l.CommonLocationType = *r.CommonLocationType
	}
	if r.CommonLocationSubType != nil {
		l.CommonLocationSubType = *r.CommonLocationSubType
	}
	if r.CommonLatitude != nil {
		l.CommonLatitude = *r.CommonLatitude
	}
	if r.CommonLongitude != nil {
		l.CommonLongitude = *r.CommonLongitude
	}
	if r.CommonAddress != nil {
		l.CommonAddress = *r.CommonAddress
	}
	if r.CommonLogo != nil {
		l.CommonLogo = *r.CommonLogo
	}
	if r.CommonContactEmail != nil {
		l.CommonContactEmail = *r.CommonContactEmail
	}
	if r.CommonContactPhone != nil {
		l.CommonContactPhone = *r.CommonContactPhone
	}
	if r.PlaceWorkHours != nil {
		l.PlaceWorkHours = r.PlaceWorkHours
	}
	if r.Custom != nil {
		l.Custom = r.Custom
	}
	if r.Localization != nil {
		l.Localization = r.Localization
	}
	if r.Source != nil {
		l.Source = *r.Source
	}
	if r.StartTime != nil {
		l.StartTime = r.StartTime
	}
	if r.EndTime != nil {
		l.EndTime = r.EndTime
	}
	if r.IsSearchable != nil {
		l.IsSearchable = *r.IsSearchable
	}
}

type SetTopLocationRequest struct {
	IsTop     bool `json:"is_top"`
	SortIndex *int `json:"sort_index"`
}

type LocationImageRequest struct {
	Original string `json:"original"`
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
}

// ── Promotion ─────────────────────────────────────────────────────────────────

type PromotionResponse struct {
	ID                uuid.UUID       `json:"id"`
	VenueID           uuid.UUID       `json:"venue_id"`
	LocationID        *uuid.UUID      `json:"location_id,omitempty"`
	ExternalID        string          `json:"external_id,omitempty"`
	PromoImage        string          `json:"promo_image,omitempty"`
	Introduction      string          `json:"introduction,omitempty"`
	GiftContent       string          `json:"gift_content,omitempty"`
	DetailURL         string          `json:"detail_url,omitempty"`
	BoothNumber       string          `json:"booth_number,omitempty"`
	ExpectedGiftCount *int            `json:"expected_gift_count,omitempty"`
	DistributionStart *time.Time      `json:"distribution_start,omitempty"`
	DistributionEnd   *time.Time      `json:"distribution_end,omitempty"`
	DisplayType       string          `json:"display_type"`
	Localization      json.RawMessage `json:"localization,omitempty" swaggertype:"object"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type PromotionRequest struct {
	LocationID        *uuid.UUID      `json:"location_id"`
	ExternalID        string          `json:"external_id"`
	PromoImage        string          `json:"promo_image"`
	Introduction      string          `json:"introduction"`
	GiftContent       string          `json:"gift_content"`
	DetailURL         string          `json:"detail_url"`
	BoothNumber       string          `json:"booth_number"`
	ExpectedGiftCount *int            `json:"expected_gift_count"`
	DistributionStart *time.Time      `json:"distribution_start"`
	DistributionEnd   *time.Time      `json:"distribution_end"`
	DisplayType       string          `json:"display_type"`
	Localization      json.RawMessage `json:"localization" swaggertype:"object"`
}

func LocationToResponse(l *domain.Location) LocationResponse {
	r := LocationResponse{
		ID: l.ID, VenueID: l.VenueID, LevelID: l.LevelID,
		MainCategoryID: l.MainCategoryID, ExternalID: l.ExternalID,
		CommonHidden: l.CommonHidden, CommonName: l.CommonName,
		CommonShortName: l.CommonShortName, CommonDescription: l.CommonDescription,
		CommonColor: l.CommonColor, CommonLocationType: l.CommonLocationType,
		CommonLocationSubType: l.CommonLocationSubType, CommonLatitude: l.CommonLatitude,
		CommonLongitude: l.CommonLongitude, CommonAddress: l.CommonAddress,
		CommonLogo: l.CommonLogo, CommonLargeLogo: l.CommonLargeLogo,
		CommonMediumLogo: l.CommonMediumLogo, CommonSmallLogo: l.CommonSmallLogo,
		CommonContactEmail: l.CommonContactEmail, CommonContactPhone: l.CommonContactPhone,
		IsTopLocation: l.IsTopLocation, IsSearchable: l.IsSearchable,
		Source: l.Source, PlaceWorkHours: l.PlaceWorkHours,
		Custom: l.Custom, Localization: l.Localization,
		StartTime: l.StartTime, EndTime: l.EndTime,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
	for _, c := range l.Categories {
		r.Categories = append(r.Categories, LocationCategoryToResponse(c))
	}
	for _, img := range l.Images {
		r.Images = append(r.Images, LocationImageResponse{
			ID: img.ID, LocationID: img.LocationID,
			Original: img.Original, Small: img.Small,
			Medium: img.Medium, Large: img.Large, CreatedAt: img.CreatedAt,
		})
	}
	return r
}
