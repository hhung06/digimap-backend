package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AppUser struct {
	ID              uuid.UUID
	VenueID         *uuid.UUID
	ExternalID      string
	Source          string
	Type            int
	FirstName       string
	LastName        string
	FirstNameEn     string
	LastNameEn      string
	Email           string
	Phone           string
	CompanyName     string
	CompanyNameEn   string
	Department      string
	Position        string
	PositionEn      string
	Section         *int
	App             string
	Token           string
	TokenExpiration *time.Time
	SurveyFlag      int
	StaffLeadFlag   int
	VisitorType     *int
	Interests       json.RawMessage
	OtherInterests  string
	IsConsented     bool
	IPAddress       string
	UserAgent       string
	BusinessName    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
