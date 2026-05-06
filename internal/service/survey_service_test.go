package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestSurveyService(repo *mocks.SurveyRepository) service.SurveyService {
	return service.NewSurveyService(repo)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestSurveyService_Create_DefaultsStatus(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	s := &domain.Survey{Title: "Customer Feedback"}

	repo.On("Create", ctx, s).Return(nil)

	err := svc.Create(ctx, s)
	require.NoError(t, err)
	assert.Equal(t, domain.SurveyStatusDraft, s.Status, "status must default to Draft")
	repo.AssertExpectations(t)
}

func TestSurveyService_Create_DefaultsSource(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	s := &domain.Survey{Title: "Customer Feedback"}

	repo.On("Create", ctx, s).Return(nil)

	err := svc.Create(ctx, s)
	require.NoError(t, err)
	assert.Equal(t, domain.SurveySourceCMS, s.Source, "source must default to CMS")
	repo.AssertExpectations(t)
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestSurveyService_Get_Success(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Survey{ID: id, Title: "Customer Feedback"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestSurveyService_Get_NotFound(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Survey)(nil), domain.NewNotFound("survey not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestSurveyService_List_Success(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Survey{{Title: "Survey A"}}

	repo.On("List", ctx, venueID, p).Return(expected, int64(1), nil)

	got, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestSurveyService_Delete_Success(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Questions ─────────────────────────────────────────────────────────────────

func TestSurveyService_CreateQuestion_Success(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	q := &domain.Question{SurveyID: uuid.New(), QuestionText: "How was your visit?"}

	repo.On("CreateQuestion", ctx, q).Return(nil)

	err := svc.CreateQuestion(ctx, q)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSurveyService_DeleteQuestion_Success(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("DeleteQuestion", ctx, id).Return(nil)

	err := svc.DeleteQuestion(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Responses ─────────────────────────────────────────────────────────────────

func TestSurveyService_SubmitResponse_Success(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	r := &domain.SurveyResponse{SurveyID: uuid.New()}

	repo.On("CreateResponse", ctx, r).Return(nil)

	err := svc.SubmitResponse(ctx, r)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── PBT: create defaults invariant ───────────────────────────────────────────

func TestSurveyService_Create_DefaultsInvariant(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		repo := &mocks.SurveyRepository{}
		svc := newTestSurveyService(repo)

		ctx := context.Background()
		s := &domain.Survey{
			Title:  rapid.StringN(1, 50, 50).Draw(rt, "title"),
			Status: 0, // always zero to test defaulting
			Source: 0,
		}

		repo.On("Create", ctx, s).Return(nil)

		err := svc.Create(ctx, s)
		require.NoError(rt, err)
		assert.NotZero(rt, s.Status, "status must always be non-zero after Create")
		assert.NotZero(rt, s.Source, "source must always be non-zero after Create")
	})
}
