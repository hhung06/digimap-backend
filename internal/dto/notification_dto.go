package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

type NotificationResponse struct {
	ID             uuid.UUID       `json:"id"`
	VenueID        *uuid.UUID      `json:"venue_id,omitempty"`
	SurveyID       *uuid.UUID      `json:"survey_id,omitempty"`
	Title          string          `json:"title,omitempty"`
	Content        string          `json:"content,omitempty"`
	Topic          string          `json:"topic,omitempty"`
	Kind           int             `json:"type"`
	Status         int             `json:"status"`
	SendStatus     int             `json:"send_status"`
	SendType       int             `json:"send_type"`
	Data           json.RawMessage `json:"data,omitempty" swaggertype:"object"`
	LinkURL        string          `json:"link_url,omitempty"`
	ScheduledAt    *time.Time      `json:"scheduled_at,omitempty"`
	TargetApp      string          `json:"target_app"`
	SegmentFilters json.RawMessage `json:"segment_filters,omitempty" swaggertype:"object"`
	RetryCount     int             `json:"retry_count"`
	PublishedAt    *time.Time      `json:"published_at,omitempty"`
	CreatedBy      *uuid.UUID      `json:"created_by,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type CreateNotificationRequest struct {
	SurveyID       *uuid.UUID      `json:"survey_id"`
	Title          string          `json:"title"`
	Content        string          `json:"content"`
	Topic          string          `json:"topic"`
	Kind           int             `json:"type"`
	SendType       int             `json:"send_type"`
	Data           json.RawMessage `json:"data" swaggertype:"object"`
	LinkURL        string          `json:"link_url"`
	ScheduledAt    *time.Time      `json:"scheduled_at"`
	TargetApp      string          `json:"target_app"`
	SegmentFilters json.RawMessage `json:"segment_filters" swaggertype:"object"`
	DeviceTokens   json.RawMessage `json:"device_tokens" swaggertype:"object"`
}

type UpdateNotificationRequest struct {
	SurveyID       *uuid.UUID      `json:"survey_id"`
	Title          string          `json:"title"`
	Content        string          `json:"content"`
	Topic          string          `json:"topic"`
	Kind           int             `json:"type"`
	SendType       int             `json:"send_type"`
	Data           json.RawMessage `json:"data" swaggertype:"object"`
	LinkURL        string          `json:"link_url"`
	ScheduledAt    *time.Time      `json:"scheduled_at"`
	TargetApp      string          `json:"target_app"`
	SegmentFilters json.RawMessage `json:"segment_filters" swaggertype:"object"`
	DeviceTokens   json.RawMessage `json:"device_tokens" swaggertype:"object"`
}

func NotificationToResponse(n *domain.Notification) NotificationResponse {
	return NotificationResponse{
		ID: n.ID, VenueID: n.VenueID, SurveyID: n.SurveyID,
		Title: n.Title, Content: n.Content, Topic: n.Topic,
		Kind: n.Kind, Status: n.Status, SendStatus: n.SendStatus, SendType: n.SendType,
		Data: n.Data, LinkURL: n.LinkURL, ScheduledAt: n.ScheduledAt,
		TargetApp: n.TargetApp, SegmentFilters: n.SegmentFilters,
		RetryCount: n.RetryCount, PublishedAt: n.PublishedAt, CreatedBy: n.CreatedBy,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}
}
