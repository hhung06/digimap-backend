package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

// ── Options ───────────────────────────────────────────────────────────────────

type OptionResponse struct {
	ID           uuid.UUID `json:"id"`
	QuestionID   uuid.UUID `json:"question_id"`
	OptionNumber int       `json:"option_number"`
	OptionText   string    `json:"option_text"`
	CreatedAt    time.Time `json:"created_at"`
}

type OptionRequest struct {
	OptionNumber int    `json:"option_number"`
	OptionText   string `json:"option_text" binding:"required"`
}

// ── Questions ─────────────────────────────────────────────────────────────────

type QuestionResponse struct {
	ID             uuid.UUID        `json:"id"`
	SurveyID       uuid.UUID        `json:"survey_id"`
	QuestionNumber int              `json:"question_number"`
	QuestionType   string           `json:"question_type"`
	QuestionText   string           `json:"question_text"`
	IsRequired     bool             `json:"is_required"`
	IsOther        bool             `json:"is_other"`
	Options        []OptionResponse `json:"options,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type QuestionRequest struct {
	QuestionNumber int            `json:"question_number"`
	QuestionType   string         `json:"question_type" binding:"required"`
	QuestionText   string         `json:"question_text" binding:"required"`
	IsRequired     bool           `json:"is_required"`
	IsOther        bool           `json:"is_other"`
	Options        []OptionRequest `json:"options"`
}

// ── Surveys ───────────────────────────────────────────────────────────────────

type SurveyResponse struct {
	ID             uuid.UUID          `json:"id"`
	VenueID        *uuid.UUID         `json:"venue_id,omitempty"`
	ExternalID     string             `json:"external_id,omitempty"`
	Title          string             `json:"title,omitempty"`
	Content        string             `json:"content,omitempty"`
	StartDate      *time.Time         `json:"start_date,omitempty"`
	EndDate        *time.Time         `json:"end_date,omitempty"`
	Status         int                `json:"status"`
	PublishType    int                `json:"publish_type"`
	IsForced       bool               `json:"is_forced"`
	Source         int                `json:"source"`
	App            string             `json:"app"`
	SegmentFilters json.RawMessage    `json:"segment_filters,omitempty" swaggertype:"object"`
	Questions      []QuestionResponse `json:"questions,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type CreateSurveyRequest struct {
	ExternalID     string          `json:"external_id"`
	Title          string          `json:"title"`
	Content        string          `json:"content"`
	StartDate      *time.Time      `json:"start_date"`
	EndDate        *time.Time      `json:"end_date"`
	Status         int             `json:"status"`
	PublishType    int             `json:"publish_type"`
	IsForced       bool            `json:"is_forced"`
	App            string          `json:"app"`
	SegmentFilters json.RawMessage `json:"segment_filters" swaggertype:"object"`
}

type UpdateSurveyRequest struct {
	ExternalID     *string         `json:"external_id"`
	Title          *string         `json:"title"`
	Content        *string         `json:"content"`
	StartDate      *time.Time      `json:"start_date"`
	EndDate        *time.Time      `json:"end_date"`
	Status         *int            `json:"status"`
	PublishType    *int            `json:"publish_type"`
	IsForced       *bool           `json:"is_forced"`
	App            *string         `json:"app"`
	SegmentFilters json.RawMessage `json:"segment_filters" swaggertype:"object"`
}

func (r UpdateSurveyRequest) ApplyTo(s *domain.Survey) {
	if r.ExternalID != nil {
		s.ExternalID = *r.ExternalID
	}
	if r.Title != nil {
		s.Title = *r.Title
	}
	if r.Content != nil {
		s.Content = *r.Content
	}
	if r.StartDate != nil {
		s.StartDate = r.StartDate
	}
	if r.EndDate != nil {
		s.EndDate = r.EndDate
	}
	if r.Status != nil {
		s.Status = *r.Status
	}
	if r.PublishType != nil {
		s.PublishType = *r.PublishType
	}
	if r.IsForced != nil {
		s.IsForced = *r.IsForced
	}
	if r.App != nil {
		s.App = *r.App
	}
	if r.SegmentFilters != nil {
		s.SegmentFilters = r.SegmentFilters
	}
}

// ── Survey responses ──────────────────────────────────────────────────────────

type SurveyResponseResponse struct {
	ID          uuid.UUID              `json:"id"`
	SurveyID    uuid.UUID              `json:"survey_id"`
	SubmittedAt time.Time              `json:"submitted_at"`
	Answers     []SurveyAnswerResponse `json:"answers,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

type SurveyAnswerResponse struct {
	ID         uuid.UUID  `json:"id"`
	QuestionID uuid.UUID  `json:"question_id"`
	OptionID   *uuid.UUID `json:"option_id,omitempty"`
	AnswerText string     `json:"answer_text,omitempty"`
}

type SubmitSurveyRequest struct {
	Answers []SurveyAnswerRequest `json:"answers" binding:"required"`
}

type SurveyAnswerRequest struct {
	QuestionID uuid.UUID  `json:"question_id" binding:"required"`
	OptionID   *uuid.UUID `json:"option_id"`
	AnswerText string     `json:"answer_text"`
}

// ── mappers ───────────────────────────────────────────────────────────────────

func SurveyToResponse(s *domain.Survey) SurveyResponse {
	r := SurveyResponse{
		ID: s.ID, VenueID: s.VenueID, ExternalID: s.ExternalID,
		Title: s.Title, Content: s.Content,
		StartDate: s.StartDate, EndDate: s.EndDate,
		Status: s.Status, PublishType: s.PublishType,
		IsForced: s.IsForced, Source: s.Source, App: s.App,
		SegmentFilters: s.SegmentFilters,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
	for _, q := range s.Questions {
		r.Questions = append(r.Questions, questionToResponse(q))
	}
	return r
}

func questionToResponse(q *domain.Question) QuestionResponse {
	r := QuestionResponse{
		ID: q.ID, SurveyID: q.SurveyID,
		QuestionNumber: q.QuestionNumber, QuestionType: q.QuestionType,
		QuestionText: q.QuestionText, IsRequired: q.IsRequired, IsOther: q.IsOther,
		CreatedAt: q.CreatedAt, UpdatedAt: q.UpdatedAt,
	}
	for _, o := range q.Options {
		r.Options = append(r.Options, OptionResponse{
			ID: o.ID, QuestionID: o.QuestionID,
			OptionNumber: o.OptionNumber, OptionText: o.OptionText, CreatedAt: o.CreatedAt,
		})
	}
	return r
}

func SurveyResponseToResponse(sr *domain.SurveyResponse) SurveyResponseResponse {
	r := SurveyResponseResponse{
		ID: sr.ID, SurveyID: sr.SurveyID,
		SubmittedAt: sr.SubmittedAt, CreatedAt: sr.CreatedAt,
	}
	for _, a := range sr.Answers {
		r.Answers = append(r.Answers, SurveyAnswerResponse{
			ID: a.ID, QuestionID: a.QuestionID,
			OptionID: a.OptionID, AnswerText: a.AnswerText,
		})
	}
	return r
}
