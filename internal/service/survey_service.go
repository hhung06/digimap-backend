package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// SurveyService manages surveys, questions, options, and responses.
type SurveyService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Survey, int64, error)
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
	SubmitResponse(ctx context.Context, r *domain.SurveyResponse) error
}

type surveyService struct {
	repo repository.SurveyRepository
}

// NewSurveyService creates a SurveyService.
func NewSurveyService(repo repository.SurveyRepository) SurveyService {
	return &surveyService{repo: repo}
}

func (s *surveyService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Survey, int64, error) {
	return s.repo.List(ctx, venueID, p)
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
	return s.repo.Create(ctx, survey)
}

func (s *surveyService) Update(ctx context.Context, survey *domain.Survey) error {
	return s.repo.Update(ctx, survey)
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

func (s *surveyService) SubmitResponse(ctx context.Context, r *domain.SurveyResponse) error {
	return s.repo.CreateResponse(ctx, r)
}
