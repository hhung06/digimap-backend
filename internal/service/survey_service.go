package service

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
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
	Duplicate(ctx context.Context, id uuid.UUID) (*domain.Survey, error)
	Stats(ctx context.Context, id uuid.UUID) (*domain.SurveyStats, error)

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
	ListParticipants(ctx context.Context, surveyID uuid.UUID, p domain.Pagination) ([]*domain.AppUser, int64, error)
	// ExportResponses builds CSV rows (header first) for all responses,
	// optionally filtered by participant external_id.
	ExportResponses(ctx context.Context, surveyID uuid.UUID, externalID string) ([][]string, error)
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

// Duplicate deep-copies a survey with its questions and options (not
// responses). All fields are copied verbatim except the title, which gets a
// " (Copy)" suffix. No activation notification is created — mirrors Django's
// surveys duplicate action.
func (s *surveyService) Duplicate(ctx context.Context, id uuid.UUID) (*domain.Survey, error) {
	original, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	copied := *original
	copied.ID = uuid.Nil
	copied.Questions = nil
	title := ""
	if original.Title != nil {
		title = *original.Title
	}
	title = strings.TrimSpace(title + " (Copy)")
	copied.Title = &title
	if err := s.repo.Create(ctx, &copied); err != nil {
		return nil, err
	}

	for _, q := range original.Questions {
		newQ := &domain.Question{
			SurveyID:       copied.ID,
			QuestionNumber: q.QuestionNumber,
			QuestionType:   q.QuestionType,
			QuestionText:   q.QuestionText,
			IsRequired:     q.IsRequired,
			IsOther:        q.IsOther,
		}
		if err := s.repo.CreateQuestion(ctx, newQ); err != nil {
			return nil, err
		}
		for _, o := range q.Options {
			newO := &domain.Option{
				QuestionID:   newQ.ID,
				OptionNumber: o.OptionNumber,
				OptionText:   o.OptionText,
			}
			if err := s.repo.CreateOption(ctx, newO); err != nil {
				return nil, err
			}
			newQ.Options = append(newQ.Options, newO)
		}
		copied.Questions = append(copied.Questions, newQ)
	}
	return &copied, nil
}

// Stats aggregates all answers per question: option selection counts with
// percentages by people/choices, "Other" free texts, and paragraph texts.
func (s *surveyService) Stats(ctx context.Context, id uuid.UUID) (*domain.SurveyStats, error) {
	survey, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	responses, err := s.repo.ListResponsesWithAnswers(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	answersByQuestion := make(map[uuid.UUID][]*domain.SurveyAnswer)
	for _, r := range responses {
		for _, a := range r.Answers {
			answersByQuestion[a.QuestionID] = append(answersByQuestion[a.QuestionID], a)
		}
	}

	stats := &domain.SurveyStats{Survey: survey, TotalResponses: int64(len(responses))}
	for _, q := range survey.Questions {
		answers := answersByQuestion[q.ID]
		people := make(map[uuid.UUID]struct{})
		totalChoices := 0
		for _, a := range answers {
			people[a.ResponseID] = struct{}{}
			if a.OptionID != nil {
				totalChoices++
			}
		}
		qs := &domain.QuestionStats{Question: q, TotalPeople: len(people), TotalChoices: totalChoices}

		switch q.QuestionType {
		case "single_choice", "multiple_choice", "rating":
			optionCounts := make(map[uuid.UUID]int)
			var otherTexts []string
			for _, a := range answers {
				switch {
				case a.OptionID != nil:
					optionCounts[*a.OptionID]++
				case q.IsOther && a.AnswerText != nil && *a.AnswerText != "":
					otherTexts = append(otherTexts, *a.AnswerText)
				}
			}
			for _, opt := range q.Options {
				optID := opt.ID
				qs.Options = append(qs.Options, &domain.OptionStats{
					OptionID:            &optID,
					Text:                opt.OptionText,
					Count:               optionCounts[opt.ID],
					PercentageByPeople:  percentage(optionCounts[opt.ID], qs.TotalPeople),
					PercentageByChoices: percentage(optionCounts[opt.ID], qs.TotalChoices),
				})
			}
			if q.IsOther {
				qs.Options = append(qs.Options, &domain.OptionStats{
					Text:                "Other",
					Count:               len(otherTexts),
					PercentageByPeople:  percentage(len(otherTexts), qs.TotalPeople),
					PercentageByChoices: percentage(len(otherTexts), qs.TotalChoices),
					OtherTexts:          otherTexts,
					IsOther:             true,
				})
			}
		default: // paragraph and free-text types
			qs.Texts = []string{}
			for _, a := range answers {
				if a.AnswerText != nil && *a.AnswerText != "" {
					qs.Texts = append(qs.Texts, *a.AnswerText)
				}
			}
		}
		stats.Questions = append(stats.Questions, qs)
	}
	return stats, nil
}

func percentage(count, total int) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(float64(count)*100*100/float64(total)) / 100
}

func (s *surveyService) ListParticipants(ctx context.Context, surveyID uuid.UUID, p domain.Pagination) ([]*domain.AppUser, int64, error) {
	if _, err := s.repo.FindByID(ctx, surveyID); err != nil {
		return nil, 0, err
	}
	return s.repo.ListParticipants(ctx, surveyID, p)
}

// ExportResponses builds the CSV rows for a survey export. Layout mirrors the
// Django export: header "No, User External ID, User Name (JP), User Name (EN),
// Submitted At, Q{n}: {text}…", one row per response ordered by submitted_at,
// choice answers joined by "; " with other-texts prefixed "[Other] ".
func (s *surveyService) ExportResponses(ctx context.Context, surveyID uuid.UUID, externalID string) ([][]string, error) {
	survey, err := s.repo.FindByID(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	var extFilter *string
	if trimmed := strings.TrimSpace(externalID); trimmed != "" {
		extFilter = &trimmed
	}
	responses, err := s.repo.ListResponsesWithAnswers(ctx, surveyID, extFilter)
	if err != nil {
		return nil, err
	}

	externalIDs := make([]string, 0, len(responses))
	seen := make(map[string]struct{})
	for _, r := range responses {
		if r.ExternalID != nil && *r.ExternalID != "" {
			if _, ok := seen[*r.ExternalID]; !ok {
				seen[*r.ExternalID] = struct{}{}
				externalIDs = append(externalIDs, *r.ExternalID)
			}
		}
	}
	users, err := s.repo.FindAppUsersByExternalIDs(ctx, externalIDs)
	if err != nil {
		return nil, err
	}
	userByExternalID := make(map[string]*domain.AppUser, len(users))
	for _, u := range users {
		userByExternalID[u.ExternalID] = u
	}

	optionTextByID := make(map[uuid.UUID]string)
	questionByID := make(map[uuid.UUID]*domain.Question, len(survey.Questions))
	for _, q := range survey.Questions {
		questionByID[q.ID] = q
		for _, o := range q.Options {
			optionTextByID[o.ID] = o.OptionText
		}
	}

	header := []string{"No", "User External ID", "User Name (JP)", "User Name (EN)", "Submitted At"}
	for _, q := range survey.Questions {
		header = append(header, "Q"+strconv.Itoa(q.QuestionNumber)+": "+q.QuestionText)
	}
	rows := [][]string{header}

	for i, r := range responses {
		extID, nameJP, nameEN := "", "", ""
		if r.ExternalID != nil {
			extID = *r.ExternalID
			if u, ok := userByExternalID[extID]; ok {
				nameJP = strings.TrimSpace(u.LastName + " " + u.FirstName)
				nameEN = strings.TrimSpace(u.FirstNameEn + " " + u.LastNameEn)
			}
		}
		row := []string{
			strconv.Itoa(i + 1), extID, nameJP, nameEN,
			r.SubmittedAt.Format("2006-01-02 15:04:05"),
		}
		answersByQuestion := make(map[uuid.UUID][]*domain.SurveyAnswer)
		for _, a := range r.Answers {
			answersByQuestion[a.QuestionID] = append(answersByQuestion[a.QuestionID], a)
		}
		for _, q := range survey.Questions {
			row = append(row, exportAnswerCell(q, answersByQuestion[q.ID], optionTextByID))
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func exportAnswerCell(q *domain.Question, answers []*domain.SurveyAnswer, optionTextByID map[uuid.UUID]string) string {
	if len(answers) == 0 {
		return ""
	}
	switch q.QuestionType {
	case "single_choice", "multiple_choice", "rating":
		var parts []string
		for _, a := range answers {
			if a.OptionID != nil {
				parts = append(parts, optionTextByID[*a.OptionID])
			}
		}
		for _, a := range answers {
			if a.OptionID == nil && q.IsOther && a.AnswerText != nil && *a.AnswerText != "" {
				parts = append(parts, "[Other] "+*a.AnswerText)
			}
		}
		return strings.Join(parts, "; ")
	case "paragraph":
		var parts []string
		for _, a := range answers {
			if a.AnswerText != nil && *a.AnswerText != "" {
				parts = append(parts, *a.AnswerText)
			}
		}
		return strings.Join(parts, "; ")
	default:
		return ""
	}
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
