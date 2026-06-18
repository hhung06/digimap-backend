package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Notification send-type constants.
const (
	NotifTypeDraft     = 1
	NotifTypeScheduled = 2
	NotifTypeImmediate = 3

	NotifStatusSent   = 1
	NotifStatusUnsent = 2

	NotifSendPending = 0
	NotifSendSuccess = 1
	NotifSendFailed  = 2

	NotifKindNormal = 1
	NotifKindSurvey = 2
)

// Notification represents a push notification record.
type Notification struct {
	ID             uuid.UUID
	VenueID        *uuid.UUID
	SurveyID       *uuid.UUID
	Title          *string
	Content        *string
	Topic          *string
	Kind           int // 1=normal 2=survey (maps to `type` column)
	Status         int // 1=sent 2=unsent
	SendStatus     int // 0=pending 1=success 2=failed
	SendType       int // 1=draft 2=scheduled 3=immediate
	Data           json.RawMessage
	LinkURL        *string
	ScheduledAt    *time.Time
	TargetApp      string
	SegmentFilters json.RawMessage
	DeviceTokens   json.RawMessage
	ErrorInfos     json.RawMessage
	RetryCount     int
	RetryAt        *time.Time
	PublishedAt    *time.Time
	CreatedBy      *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
