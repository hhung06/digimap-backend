package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

// ── Customer ──────────────────────────────────────────────────────────────────

type CustomerResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
	Address     string    `json:"address,omitempty"`
	URL         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CustomerRequest struct {
	Name        string `json:"name"        binding:"required"`
	Image       string `json:"image"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type UpdateCustomerRequest struct {
	Name        *string `json:"name"`
	Image       *string `json:"image"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Address     *string `json:"address"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
}

func (r UpdateCustomerRequest) ApplyTo(c *domain.Customer) {
	if r.Name != nil {
		c.Name = *r.Name
	}
	if r.Image != nil {
		c.Image = *r.Image
	}
	if r.Phone != nil {
		c.Phone = *r.Phone
	}
	if r.Email != nil {
		c.Email = *r.Email
	}
	if r.Address != nil {
		c.Address = *r.Address
	}
	if r.URL != nil {
		c.URL = *r.URL
	}
	if r.Description != nil {
		c.Description = *r.Description
	}
}

func CustomerToResponse(c *domain.Customer) CustomerResponse {
	return CustomerResponse{
		ID:          c.ID,
		Name:        c.Name,
		Image:       c.Image,
		Phone:       c.Phone,
		Email:       c.Email,
		Address:     c.Address,
		URL:         c.URL,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// ── Venue ─────────────────────────────────────────────────────────────────────

var venueTypeNames = map[int]string{
	1:  "Shopping Mall",
	2:  "University",
	3:  "Office",
	4:  "Hospital",
	5:  "Expo / Exhibition",
	6:  "Stadium / Sport",
	7:  "Airport",
	8:  "Amusement Park",
	9:  "Real Estate",
	10: "Inventory / Warehouse",
	11: "Others",
}

type VenueBriefCustomer struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type VenueResponse struct {
	ID              uuid.UUID          `json:"id"`
	Customer        VenueBriefCustomer `json:"customer"`
	Name            string             `json:"name"`
	ExternalID      string             `json:"external_id,omitempty"`
	Type            int                `json:"type"`
	TypeDisplayName string             `json:"type_display_name"`
	PublicKey       string             `json:"public_key"`
	Address         string             `json:"address,omitempty"`
	City            string             `json:"city,omitempty"`
	State           string             `json:"state,omitempty"`
	Country         string             `json:"country,omitempty"`
	Postal          string             `json:"postal,omitempty"`
	Lat             float64            `json:"lat"`
	Lng             float64            `json:"lng"`
	Timezone        string             `json:"timezone"`
	Telephone       string             `json:"telephone,omitempty"`
	Description     string             `json:"description,omitempty"`
	FloorCount      int                `json:"floor_count"`
	Languages       []string           `json:"languages"`
	Localization    json.RawMessage    `json:"localization,omitempty" swaggertype:"object"`
	AppConfigs      json.RawMessage    `json:"app_configs,omitempty" swaggertype:"object"`
	AppDomains      json.RawMessage    `json:"app_domains,omitempty" swaggertype:"object"`
	SubDomains      string             `json:"sub_domains,omitempty"`
	SEOTitle        string             `json:"seo_title,omitempty"`
	SEODescription  string             `json:"seo_description,omitempty"`
	SEOKeywords     string             `json:"seo_keywords,omitempty"`
	HeadTag         string             `json:"head_tag,omitempty"`
	BodyTag         string             `json:"body_tag,omitempty"`
	OriginalLogo    string             `json:"original_logo,omitempty"`
	SmallLogo       string             `json:"small_logo,omitempty"`
	MediumLogo      string             `json:"medium_logo,omitempty"`
	LargeLogo       string             `json:"large_logo,omitempty"`
	StartAt         *time.Time         `json:"start_at,omitempty"`
	EndAt           *time.Time         `json:"end_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type VenueKeyResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type CreateVenueRequest struct {
	CustomerID     uuid.UUID       `json:"customer_id"     binding:"required"`
	Name           string          `json:"name"            binding:"required"`
	Slug           string          `json:"slug"`
	ExternalID     string          `json:"external_id"`
	Type           int             `json:"type"`
	Address        string          `json:"address"`
	City           string          `json:"city"`
	State          string          `json:"state"`
	Country        string          `json:"country"`
	Postal         string          `json:"postal"`
	Lat            float64         `json:"lat"`
	Lng            float64         `json:"lng"`
	Timezone       string          `json:"timezone"`
	Telephone      string          `json:"telephone"`
	WorkHours      string          `json:"work_hours"`
	Description    string          `json:"description"`
	Translations   json.RawMessage `json:"translations" swaggertype:"object"`
	Localization   json.RawMessage `json:"localization" swaggertype:"object"`
	CustomData     json.RawMessage `json:"custom_data" swaggertype:"object"`
	AppConfigs     json.RawMessage `json:"app_configs" swaggertype:"object"`
	AppDomains     json.RawMessage `json:"app_domains" swaggertype:"object"`
	SubDomains     string          `json:"sub_domains"`
	SEOTitle       string          `json:"seo_title"`
	SEODescription string          `json:"seo_description"`
	SEOKeywords    string          `json:"seo_keywords"`
	HeadTag        string          `json:"head_tag"`
	BodyTag        string          `json:"body_tag"`
	StartAt        *time.Time      `json:"start_at"`
	EndAt          *time.Time      `json:"end_at"`
}

// UpdateVenueRequest uses pointer fields so only explicitly provided fields are applied.
type UpdateVenueRequest struct {
	Name           *string         `json:"name"`
	ExternalID     *string         `json:"external_id"`
	Type           *int            `json:"type"`
	Address        *string         `json:"address"`
	City           *string         `json:"city"`
	State          *string         `json:"state"`
	Country        *string         `json:"country"`
	Postal         *string         `json:"postal"`
	Lat            *float64        `json:"lat"`
	Lng            *float64        `json:"lng"`
	Timezone       *string         `json:"timezone"`
	Telephone      *string         `json:"telephone"`
	Description    *string         `json:"description"`
	Localization   json.RawMessage `json:"localization" swaggertype:"object"`
	AppConfigs     json.RawMessage `json:"app_configs" swaggertype:"object"`
	AppDomains     json.RawMessage `json:"app_domains" swaggertype:"object"`
	SubDomains     *string         `json:"sub_domains"`
	SEOTitle       *string         `json:"seo_title"`
	SEODescription *string         `json:"seo_description"`
	SEOKeywords    *string         `json:"seo_keywords"`
	HeadTag        *string         `json:"head_tag"`
	BodyTag        *string         `json:"body_tag"`
	OriginalLogo   *string         `json:"original_logo"`
	SmallLogo      *string         `json:"small_logo"`
	MediumLogo     *string         `json:"medium_logo"`
	LargeLogo      *string         `json:"large_logo"`
	StartAt        *time.Time      `json:"start_at"`
	EndAt          *time.Time      `json:"end_at"`
}

// ApplyTo merges non-nil fields from the request onto an existing venue.
func (r UpdateVenueRequest) ApplyTo(v *domain.Venue) {
	if r.Name != nil {
		v.Name = *r.Name
	}
	if r.ExternalID != nil {
		v.ExternalID = *r.ExternalID
	}
	if r.Type != nil {
		v.Type = *r.Type
	}
	if r.Address != nil {
		v.Address = *r.Address
	}
	if r.City != nil {
		v.City = *r.City
	}
	if r.State != nil {
		v.State = *r.State
	}
	if r.Country != nil {
		v.Country = *r.Country
	}
	if r.Postal != nil {
		v.Postal = *r.Postal
	}
	if r.Lat != nil {
		v.Lat = *r.Lat
	}
	if r.Lng != nil {
		v.Lng = *r.Lng
	}
	if r.Timezone != nil {
		v.Timezone = *r.Timezone
	}
	if r.Telephone != nil {
		v.Telephone = *r.Telephone
	}
	if r.Description != nil {
		v.Description = *r.Description
	}
	if r.Localization != nil {
		v.Localization = r.Localization
	}
	if r.AppConfigs != nil {
		v.AppConfigs = r.AppConfigs
	}
	if r.AppDomains != nil {
		v.AppDomains = r.AppDomains
	}
	if r.SubDomains != nil {
		v.SubDomains = *r.SubDomains
	}
	if r.SEOTitle != nil {
		v.SEOTitle = *r.SEOTitle
	}
	if r.SEODescription != nil {
		v.SEODescription = *r.SEODescription
	}
	if r.SEOKeywords != nil {
		v.SEOKeywords = *r.SEOKeywords
	}
	if r.HeadTag != nil {
		v.HeadTag = *r.HeadTag
	}
	if r.BodyTag != nil {
		v.BodyTag = *r.BodyTag
	}
	if r.OriginalLogo != nil {
		v.OriginalLogo = *r.OriginalLogo
	}
	if r.SmallLogo != nil {
		v.SmallLogo = *r.SmallLogo
	}
	if r.MediumLogo != nil {
		v.MediumLogo = *r.MediumLogo
	}
	if r.LargeLogo != nil {
		v.LargeLogo = *r.LargeLogo
	}
	if r.StartAt != nil {
		v.StartAt = r.StartAt
	}
	if r.EndAt != nil {
		v.EndAt = r.EndAt
	}
}

func VenueToResponse(v *domain.Venue) VenueResponse {
	langs := v.Languages
	if langs == nil {
		langs = []string{}
	}
	return VenueResponse{
		ID:              v.ID,
		Customer:        VenueBriefCustomer{ID: v.CustomerID, Name: v.CustomerName},
		Name:            v.Name,
		ExternalID:      v.ExternalID,
		Type:            v.Type,
		TypeDisplayName: venueTypeNames[v.Type],
		PublicKey:       v.PublicKey,
		Address:         v.Address, City: v.City, State: v.State, Country: v.Country,
		Postal:          v.Postal, Lat: v.Lat, Lng: v.Lng, Timezone: v.Timezone,
		Telephone:       v.Telephone, Description: v.Description,
		FloorCount:      v.FloorCount,
		Languages:       langs,
		Localization:    v.Localization,
		AppConfigs:      v.AppConfigs, AppDomains: v.AppDomains, SubDomains: v.SubDomains,
		SEOTitle:        v.SEOTitle, SEODescription: v.SEODescription, SEOKeywords: v.SEOKeywords,
		HeadTag:         v.HeadTag, BodyTag: v.BodyTag,
		OriginalLogo:    v.OriginalLogo, SmallLogo: v.SmallLogo,
		MediumLogo:      v.MediumLogo, LargeLogo: v.LargeLogo,
		StartAt:         v.StartAt, EndAt: v.EndAt,
		CreatedAt:       v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}

// ── Level ─────────────────────────────────────────────────────────────────────

type MapGroupResponse struct {
	ID        uuid.UUID `json:"id"`
	VenueID   uuid.UUID `json:"venue_id"`
	Type      string    `json:"type,omitempty"`
	Name      string    `json:"name,omitempty"`
	ShortName string    `json:"short_name,omitempty"`
	SortIndex int       `json:"sort_index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MapGroupRequest struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	SortIndex int    `json:"sort_index"`
}

type UpdateMapGroupRequest struct {
	Type      *string `json:"type"`
	Name      *string `json:"name"`
	ShortName *string `json:"short_name"`
	SortIndex *int    `json:"sort_index"`
}

func (r UpdateMapGroupRequest) ApplyTo(mg *domain.MapGroup) {
	if r.Type != nil {
		mg.Type = *r.Type
	}
	if r.Name != nil {
		mg.Name = *r.Name
	}
	if r.ShortName != nil {
		mg.ShortName = *r.ShortName
	}
	if r.SortIndex != nil {
		mg.SortIndex = *r.SortIndex
	}
}

type PerspectiveResponse struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name,omitempty"`
	CameraZoom            float64   `json:"camera_zoom"`
	CameraType            int       `json:"camera_type"`
	CameraMaxZoom         float64   `json:"camera_max_zoom"`
	CameraMinZoom         float64   `json:"camera_min_zoom"`
	CameraTargetCenterLng float64   `json:"camera_target_center_lng"`
	CameraTargetCenterLat float64   `json:"camera_target_center_lat"`
	CameraTargetZoom      float64   `json:"camera_target_zoom"`
	CameraTargetBearing   float64   `json:"camera_target_bearing"`
	CameraTargetPitch     float64   `json:"camera_target_pitch"`
}

type PerspectiveRequest struct {
	Name                  string  `json:"name"`
	CameraZoom            float64 `json:"camera_zoom"`
	CameraType            int     `json:"camera_type"`
	CameraMaxZoom         float64 `json:"camera_max_zoom"`
	CameraMinZoom         float64 `json:"camera_min_zoom"`
	CameraTargetCenterLng float64 `json:"camera_target_center_lng"`
	CameraTargetCenterLat float64 `json:"camera_target_center_lat"`
	CameraTargetZoom      float64 `json:"camera_target_zoom"`
	CameraTargetBearing   float64 `json:"camera_target_bearing"`
	CameraTargetPitch     float64 `json:"camera_target_pitch"`
}

type LevelResponse struct {
	ID            uuid.UUID            `json:"id"`
	VenueID       uuid.UUID            `json:"venue_id"`
	MapGroupID    *uuid.UUID           `json:"map_group_id,omitempty"`
	PerspectiveID *uuid.UUID           `json:"perspective_id,omitempty"`
	Name          string               `json:"name,omitempty"`
	ShortName     string               `json:"short_name,omitempty"`
	ExternalID    string               `json:"external_id,omitempty"`
	Type          int                  `json:"type"`
	Latitude      float64              `json:"latitude"`
	Longitude     float64              `json:"longitude"`
	Bearing       float64              `json:"bearing"`
	Width         int                  `json:"width"`
	Height        int                  `json:"height"`
	Scale         float64              `json:"scale"`
	LevelWidth    float64              `json:"level_width"`
	LevelHeight   float64              `json:"level_height"`
	FileIDs       string               `json:"file_ids,omitempty"`
	Elevation     *int                 `json:"elevation,omitempty"`
	IsPublished   bool                 `json:"is_published"`
	Perspective   *PerspectiveResponse `json:"perspective,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

type CreateLevelRequest struct {
	MapGroupID  *uuid.UUID `json:"map_group_id"`
	Name        string     `json:"name"`
	ShortName   string     `json:"short_name"`
	ExternalID  string     `json:"external_id"`
	Type        int        `json:"type"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Bearing     float64    `json:"bearing"`
	Width       int        `json:"width"`
	Height      int        `json:"height"`
	Scale       float64    `json:"scale"`
	LevelWidth  float64    `json:"level_width"`
	LevelHeight float64    `json:"level_height"`
	FileIDs     string     `json:"file_ids"`
	Elevation   *int       `json:"elevation"`
	IsPublished bool       `json:"is_published"`
}

type UpdateLevelRequest struct {
	MapGroupID  *uuid.UUID `json:"map_group_id"`
	Name        *string    `json:"name"`
	ShortName   *string    `json:"short_name"`
	ExternalID  *string    `json:"external_id"`
	Type        *int       `json:"type"`
	Latitude    *float64   `json:"latitude"`
	Longitude   *float64   `json:"longitude"`
	Bearing     *float64   `json:"bearing"`
	Width       *int       `json:"width"`
	Height      *int       `json:"height"`
	Scale       *float64   `json:"scale"`
	LevelWidth  *float64   `json:"level_width"`
	LevelHeight *float64   `json:"level_height"`
	FileIDs     *string    `json:"file_ids"`
	Elevation   *int       `json:"elevation"`
	IsPublished *bool      `json:"is_published"`
}

func (r UpdateLevelRequest) ApplyTo(l *domain.Level) {
	if r.MapGroupID != nil {
		l.MapGroupID = r.MapGroupID
	}
	if r.Name != nil {
		l.Name = *r.Name
	}
	if r.ShortName != nil {
		l.ShortName = *r.ShortName
	}
	if r.ExternalID != nil {
		l.ExternalID = *r.ExternalID
	}
	if r.Type != nil {
		l.Type = *r.Type
	}
	if r.Latitude != nil {
		l.Latitude = *r.Latitude
	}
	if r.Longitude != nil {
		l.Longitude = *r.Longitude
	}
	if r.Bearing != nil {
		l.Bearing = *r.Bearing
	}
	if r.Width != nil {
		l.Width = *r.Width
	}
	if r.Height != nil {
		l.Height = *r.Height
	}
	if r.Scale != nil {
		l.Scale = *r.Scale
	}
	if r.LevelWidth != nil {
		l.LevelWidth = *r.LevelWidth
	}
	if r.LevelHeight != nil {
		l.LevelHeight = *r.LevelHeight
	}
	if r.FileIDs != nil {
		l.FileIDs = *r.FileIDs
	}
	if r.Elevation != nil {
		l.Elevation = r.Elevation
	}
	if r.IsPublished != nil {
		l.IsPublished = *r.IsPublished
	}
}

type GeoReferenceResponse struct {
	ID        uuid.UUID `json:"id"`
	LevelID   uuid.UUID `json:"level_id"`
	ControlX  int       `json:"control_x"`
	ControlY  int       `json:"control_y"`
	TargetX   float64   `json:"target_x"`
	TargetY   float64   `json:"target_y"`
	CreatedAt time.Time `json:"created_at"`
}

type GeoReferenceRequest struct {
	ControlX int     `json:"control_x"`
	ControlY int     `json:"control_y"`
	TargetX  float64 `json:"target_x"`
	TargetY  float64 `json:"target_y"`
}

func MapGroupToResponse(mg *domain.MapGroup) MapGroupResponse {
	return MapGroupResponse{
		ID: mg.ID, VenueID: mg.VenueID, Type: mg.Type,
		Name: mg.Name, ShortName: mg.ShortName, SortIndex: mg.SortIndex,
		CreatedAt: mg.CreatedAt, UpdatedAt: mg.UpdatedAt,
	}
}

func LevelToResponse(l *domain.Level) LevelResponse {
	r := LevelResponse{
		ID: l.ID, VenueID: l.VenueID,
		MapGroupID: l.MapGroupID, PerspectiveID: l.PerspectiveID,
		Name: l.Name, ShortName: l.ShortName, ExternalID: l.ExternalID,
		Type: l.Type, Latitude: l.Latitude, Longitude: l.Longitude, Bearing: l.Bearing,
		Width: l.Width, Height: l.Height, Scale: l.Scale,
		LevelWidth: l.LevelWidth, LevelHeight: l.LevelHeight,
		FileIDs: l.FileIDs, Elevation: l.Elevation, IsPublished: l.IsPublished,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
	if l.Perspective != nil {
		p := PerspectiveResponse{
			ID: l.Perspective.ID, Name: l.Perspective.Name,
			CameraZoom: l.Perspective.CameraZoom, CameraType: l.Perspective.CameraType,
			CameraMaxZoom: l.Perspective.CameraMaxZoom, CameraMinZoom: l.Perspective.CameraMinZoom,
			CameraTargetCenterLng: l.Perspective.CameraTargetCenterLng,
			CameraTargetCenterLat: l.Perspective.CameraTargetCenterLat,
			CameraTargetZoom:      l.Perspective.CameraTargetZoom,
			CameraTargetBearing:   l.Perspective.CameraTargetBearing,
			CameraTargetPitch:     l.Perspective.CameraTargetPitch,
		}
		r.Perspective = &p
	}
	return r
}

func GeoRefToResponse(g *domain.GeoReference) GeoReferenceResponse {
	return GeoReferenceResponse{
		ID: g.ID, LevelID: g.LevelID,
		ControlX: g.ControlX, ControlY: g.ControlY,
		TargetX: g.TargetX, TargetY: g.TargetY,
		CreatedAt: g.CreatedAt,
	}
}
