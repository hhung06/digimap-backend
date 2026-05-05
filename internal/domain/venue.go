package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ── Customer ──────────────────────────────────────────────────────────────────

type Customer struct {
	ID          uuid.UUID
	Name        string
	Image       string
	Phone       string
	Email       string
	Address     string
	URL         string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// ── Venue ─────────────────────────────────────────────────────────────────────

type Venue struct {
	ID             uuid.UUID
	CustomerID     uuid.UUID
	ExternalID     string
	Name           string
	Type           int
	PublicKey      string
	PrivateKey     string
	Address        string
	City           string
	State          string
	Country        string
	Postal         string
	Lat            float64
	Lng            float64
	Timezone       string
	Telephone      string
	Description    string
	Localization   json.RawMessage
	AppConfigs     json.RawMessage
	AppDomains     json.RawMessage
	SubDomains     string
	SEOTitle       string
	SEODescription string
	SEOKeywords    string
	HeadTag        string
	BodyTag        string
	OriginalLogo   string
	SmallLogo      string
	MediumLogo     string
	LargeLogo      string
	StartAt        *time.Time
	EndAt          *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
	// Eagerly joined
	CustomerName string
	Languages    []string
	FloorCount   int
}

// ── Level-related ─────────────────────────────────────────────────────────────

type MapGroup struct {
	ID        uuid.UUID
	VenueID   uuid.UUID
	Type      string
	Name      string
	ShortName string
	SortIndex int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Perspective struct {
	ID                    uuid.UUID
	Name                  string
	CameraZoom            float64
	CameraType            int
	CameraMaxZoom         float64
	CameraMinZoom         float64
	CameraTargetCenterLng float64
	CameraTargetCenterLat float64
	CameraTargetZoom      float64
	CameraTargetBearing   float64
	CameraTargetPitch     float64
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

type Level struct {
	ID            uuid.UUID
	VenueID       uuid.UUID
	MapGroupID    *uuid.UUID
	PerspectiveID *uuid.UUID
	Name          string
	ShortName     string
	ExternalID    string
	Type          int
	Latitude      float64
	Longitude     float64
	Bearing       float64
	Width         int
	Height        int
	Scale         float64
	LevelWidth    float64
	LevelHeight   float64
	FileIDs       string
	Elevation     *int
	IsPublished   bool
	Perspective   *Perspective // eagerly loaded when needed
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type GeoReference struct {
	ID        uuid.UUID
	LevelID   uuid.UUID
	ControlX  int
	ControlY  int
	TargetX   float64
	TargetY   float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
