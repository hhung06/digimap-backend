package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Survey status constants.
const (
	SurveyStatusDraft    = 1
	SurveyStatusActive   = 2
	SurveyStatusInactive = 3
	SurveyStatusClosed   = 4

	SurveyPublishPush  = 1
	SurveyPublishInApp = 2
	SurveyPublishBoth  = 3
	SurveyPublishPromo = 4

	SurveySourceCMS    = 1
	SurveySourceApp    = 2
	SurveySourceImport = 3
)

// Survey is a questionnaire associated with a venue.
type Survey struct {
	ID             uuid.UUID
	VenueID        *uuid.UUID
	ExternalID     *string
	Title          *string
	Content        *string
	StartDate      *time.Time
	EndDate        *time.Time
	Status         int
	PublishType    int
	IsForced       bool
	Source         int
	App            string
	SegmentFilters json.RawMessage
	CreatedBy      *uuid.UUID
	Questions      []*Question // eagerly loaded on Get
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SurveyListFilter struct {
	Status      *int
	PublishType *int
	Keyword     string
}

// Question is a single question within a Survey.
type Question struct {
	ID             uuid.UUID
	SurveyID       uuid.UUID
	QuestionNumber int
	QuestionType   string // paragraph, single_choice, multiple_choice, rating, …
	QuestionText   string
	IsRequired     bool
	IsOther        bool
	Options        []*Option // eagerly loaded
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Option is a selectable choice for a Question.
type Option struct {
	ID           uuid.UUID
	QuestionID   uuid.UUID
	OptionNumber int
	OptionText   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SurveyResponse records that a participant submitted a survey.
type SurveyResponse struct {
	ID          uuid.UUID
	SurveyID    uuid.UUID
	ExternalID  *string
	SubmittedAt time.Time
	Answers     []*SurveyAnswer
	CreatedAt   time.Time
}

// SurveyAnswer records a single answer within a SurveyResponse.
type SurveyAnswer struct {
	ID         uuid.UUID
	ResponseID uuid.UUID
	QuestionID uuid.UUID
	OptionID   *uuid.UUID
	AnswerText *string
	CreatedAt  time.Time
}

// SurveyStats aggregates response counts for a survey.
type SurveyStats struct {
	Survey         *Survey
	TotalResponses int64
	Questions      []*QuestionStats
}

// QuestionStats aggregates answers for one question. Options is populated for
// choice-type questions; Texts holds free-text answers for paragraph questions.
type QuestionStats struct {
	Question     *Question
	TotalPeople  int
	TotalChoices int
	Options      []*OptionStats
	Texts        []string
}

// OptionStats counts selections of one option. A nil OptionID with OtherTexts
// set represents the synthetic "Other" entry of an is_other question.
type OptionStats struct {
	OptionID            *uuid.UUID
	Text                string
	Count               int
	PercentageByPeople  float64
	PercentageByChoices float64
	OtherTexts          []string
	IsOther             bool
}
