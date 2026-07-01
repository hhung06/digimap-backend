package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// SurveyService manages surveys, questions, options, and responses.
type SurveyService interface {
	List(ctx context.Context, venueID uuid.UUID, filter domain.SurveyListFilter, p domain.Pagination) ([]*domain.Survey, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Survey, error)
	Create(ctx context.Context, s *domain.Survey) error
	Update(ctx context.Context, s *domain.Survey) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Questions
	CreateQuestion(ctx context.Context, q *domain.Question) error
	UpdateQuestion(ctx context.Context, q *domain.Question) error
	DeleteQuestion(ctx context.Context, id uuid.UUID) error

	// Options
	CreateOption(ctx context.Context, o *domain.Option) error
	UpdateOption(ctx context.Context, o *domain.Option) error
	DeleteOption(ctx context.Context, id uuid.UUID) error

	// Active surveys for visitor/promo consumption
	ListPromo(ctx context.Context, venueID uuid.UUID) ([]*domain.Survey, error)
	ListActive(ctx context.Context, venueID uuid.UUID, publishTypes []int) ([]*domain.Survey, error)

	// Responses
	ListResponses(ctx context.Context, surveyID uuid.UUID, p domain.Pagination) ([]*domain.SurveyResponse, int64, error)
	SubmitPublicResponse(ctx context.Context, r *domain.SurveyResponse) error
	SubmitAppResponse(ctx context.Context, venueID uuid.UUID, r *domain.SurveyResponse) error
	SubmitVenueResponse(ctx context.Context, venueID uuid.UUID, r *domain.SurveyResponse) error
	ProcessScheduledTransitions(ctx context.Context, now time.Time) (int, int, error)
}

type notificationSender interface {
	Send(ctx context.Context, id uuid.UUID) error
}

type surveyService struct {
	repo          repository.SurveyRepository
	notifRepo     repository.NotificationRepository
	notifications notificationSender
}

// NewSurveyService creates a SurveyService.
func NewSurveyService(repo repository.SurveyRepository, notifRepo repository.NotificationRepository, notifications notificationSender) SurveyService {
	return &surveyService{repo: repo, notifRepo: notifRepo, notifications: notifications}
}

func (s *surveyService) List(ctx context.Context, venueID uuid.UUID, filter domain.SurveyListFilter, p domain.Pagination) ([]*domain.Survey, int64, error) {
	return s.repo.List(ctx, venueID, filter, p)
}

func (s *surveyService) Get(ctx context.Context, id uuid.UUID) (*domain.Survey, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *surveyService) Create(ctx context.Context, survey *domain.Survey) error {
	if survey.Status == 0 {
		survey.Status = domain.SurveyStatusDraft
	}
	if survey.Source == 0 {
		survey.Source = domain.SurveySourceCMS
	}
	if err := s.repo.Create(ctx, survey); err != nil {
		return err
	}
	return s.handleActivationNotification(ctx, survey, 0)
}

func (s *surveyService) Update(ctx context.Context, survey *domain.Survey) error {
	before, err := s.repo.FindByID(ctx, survey.ID)
	if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, survey); err != nil {
		return err
	}
	return s.handleActivationNotification(ctx, survey, before.Status)
}

func (s *surveyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *surveyService) CreateQuestion(ctx context.Context, q *domain.Question) error {
	return s.repo.CreateQuestion(ctx, q)
}

func (s *surveyService) UpdateQuestion(ctx context.Context, q *domain.Question) error {
	return s.repo.UpdateQuestion(ctx, q)
}

func (s *surveyService) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteQuestion(ctx, id)
}

func (s *surveyService) CreateOption(ctx context.Context, o *domain.Option) error {
	return s.repo.CreateOption(ctx, o)
}

func (s *surveyService) UpdateOption(ctx context.Context, o *domain.Option) error {
	return s.repo.UpdateOption(ctx, o)
}

func (s *surveyService) DeleteOption(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteOption(ctx, id)
}

func (s *surveyService) ListPromo(ctx context.Context, venueID uuid.UUID) ([]*domain.Survey, error) {
	return s.repo.ListActive(ctx, venueID, []int{domain.SurveyPublishPromo})
}

func (s *surveyService) ListActive(ctx context.Context, venueID uuid.UUID, publishTypes []int) ([]*domain.Survey, error) {
	return s.repo.ListActive(ctx, venueID, publishTypes)
}

func (s *surveyService) ListResponses(ctx context.Context, surveyID uuid.UUID, p domain.Pagination) ([]*domain.SurveyResponse, int64, error) {
	return s.repo.ListResponses(ctx, surveyID, p)
}

func (s *surveyService) SubmitPublicResponse(ctx context.Context, r *domain.SurveyResponse) error {
	if err := s.validateResponse(ctx, r, nil, false); err != nil {
		return err
	}
	return s.repo.CreateResponse(ctx, r)
}

func (s *surveyService) SubmitAppResponse(ctx context.Context, venueID uuid.UUID, r *domain.SurveyResponse) error {
	if err := s.validateResponse(ctx, r, &venueID, true); err != nil {
		return err
	}
	return s.repo.CreateResponse(ctx, r)
}

func (s *surveyService) SubmitVenueResponse(ctx context.Context, venueID uuid.UUID, r *domain.SurveyResponse) error {
	if err := s.validateResponse(ctx, r, &venueID, false); err != nil {
		return err
	}
	return s.repo.CreateResponse(ctx, r)
}

func (s *surveyService) ProcessScheduledTransitions(ctx context.Context, now time.Time) (int, int, error) {
	toActivate, err := s.repo.ListDueActivation(ctx, now)
	if err != nil {
		return 0, 0, err
	}
	activated := 0
	for _, survey := range toActivate {
		survey.Status = domain.SurveyStatusActive
		if err := s.repo.Update(ctx, survey); err != nil {
			return activated, 0, err
		}
		if err := s.handleActivationNotification(ctx, survey, domain.SurveyStatusInactive); err != nil {
			return activated, 0, err
		}
		activated++
	}

	toClose, err := s.repo.ListDueClosure(ctx, now)
	if err != nil {
		return activated, 0, err
	}
	closed := 0
	for _, survey := range toClose {
		survey.Status = domain.SurveyStatusClosed
		if err := s.repo.Update(ctx, survey); err != nil {
			return activated, closed, err
		}
		closed++
	}

	return activated, closed, nil
}

func (s *surveyService) handleActivationNotification(ctx context.Context, survey *domain.Survey, previousStatus int) error {
	if s.notifRepo == nil || survey.Source != domain.SurveySourceCMS || survey.Status != domain.SurveyStatusActive || survey.StartDate == nil {
		return nil
	}
	if previousStatus != 0 && previousStatus != domain.SurveyStatusDraft && previousStatus != domain.SurveyStatusInactive {
		return nil
	}

	sendType := domain.NotifTypeImmediate
	scheduledAt := ptrSurveyTime(time.Now())
	if survey.StartDate.After(time.Now()) {
		sendType = domain.NotifTypeScheduled
		scheduledAt = survey.StartDate
	}

	data, _ := json.Marshal(map[string]any{
		"survey_id":    survey.ID.String(),
		"title":        survey.Title,
		"content":      survey.Content,
		"type":         "survey",
		"publish_type": survey.PublishType,
	})

	n := &domain.Notification{
		VenueID:        survey.VenueID,
		SurveyID:       &survey.ID,
		Title:          survey.Title,
		Content:        survey.Content,
		Kind:           domain.NotifKindSurvey,
		Status:         domain.NotifStatusUnsent,
		SendStatus:     domain.NotifSendPending,
		SendType:       sendType,
		Data:           data,
		ScheduledAt:    scheduledAt,
		TargetApp:      "all",
		SegmentFilters: survey.SegmentFilters,
		CreatedBy:      survey.CreatedBy,
	}
	if err := validateNotificationDelivery(n); err != nil {
		return err
	}
	if err := s.notifRepo.Create(ctx, n); err != nil {
		return err
	}
	if sendType == domain.NotifTypeImmediate && s.notifications != nil {
		_ = s.notifications.Send(ctx, n.ID)
	}
	return nil
}

func (s *surveyService) validateResponse(ctx context.Context, r *domain.SurveyResponse, expectedVenueID *uuid.UUID, requireExternalID bool) error {
	if requireExternalID && (r.ExternalID == nil || strings.TrimSpace(*r.ExternalID) == "") {
		return domain.NewValidation(map[string]string{"external_id": "is required"})
	}

	survey, err := s.repo.FindByID(ctx, r.SurveyID)
	if err != nil {
		return err
	}
	if expectedVenueID != nil {
		if survey.VenueID == nil || *survey.VenueID != *expectedVenueID {
			return domain.NewValidation(map[string]string{"survey_id": "does not belong to the expected venue"})
		}
	}
	if survey.Status != domain.SurveyStatusActive {
		return domain.NewValidation(map[string]string{"survey_id": "survey is not accepting responses"})
	}

	questionByID := make(map[uuid.UUID]*domain.Question, len(survey.Questions))
	for _, question := range survey.Questions {
		questionByID[question.ID] = question
	}

	requiredAnswers := make(map[uuid.UUID]int)
	for _, answer := range r.Answers {
		question, ok := questionByID[answer.QuestionID]
		if !ok {
			return domain.NewValidation(map[string]string{"question_id": "question does not belong to survey"})
		}
		if err := validateSurveyAnswer(question, answer); err != nil {
			return err
		}
		if question.IsRequired {
			requiredAnswers[question.ID]++
		}
	}

	for _, question := range survey.Questions {
		if question.IsRequired && requiredAnswers[question.ID] == 0 {
			return domain.NewValidation(map[string]string{"question_id": "missing required answer"})
		}
	}

	return nil
}

func validateSurveyAnswer(question *domain.Question, answer *domain.SurveyAnswer) error {
	answerText := ""
	if answer.AnswerText != nil {
		answerText = *answer.AnswerText
	}
	text := strings.TrimSpace(answerText)

	if answer.OptionID != nil {
		validOption := false
		for _, option := range question.Options {
			if option.ID == *answer.OptionID {
				validOption = true
				break
			}
		}
		if !validOption {
			return domain.NewValidation(map[string]string{"option_id": "option does not belong to question"})
		}
	}

	switch question.QuestionType {
	case "paragraph":
		if answer.OptionID != nil {
			return domain.NewValidation(map[string]string{"option_id": "option answers are not allowed for this question type"})
		}
		if text == "" {
			return domain.NewValidation(map[string]string{"answer_text": "answer text is required"})
		}
	case "single_choice", "multiple_choice", "rating":
		if answer.OptionID == nil {
			return domain.NewValidation(map[string]string{"option_id": "option is required"})
		}
		if text != "" {
			return domain.NewValidation(map[string]string{"answer_text": "answer text is not allowed for this question type"})
		}
	default:
		if answer.OptionID == nil && text == "" {
			return domain.NewValidation(map[string]string{"answer": "either option_id or answer_text is required"})
		}
	}

	return nil
}

func ptrSurveyTime(v time.Time) *time.Time {
	return &v
}
