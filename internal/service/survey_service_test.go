package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

type mockNotificationSender struct{ mock.Mock }

func (m *mockNotificationSender) Send(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func newTestSurveyService(repo *mocks.SurveyRepository, extras ...any) service.SurveyService {
	var notifRepo *mocks.NotificationRepository
	var sender *mockNotificationSender
	if len(extras) > 0 && extras[0] != nil {
		notifRepo = extras[0].(*mocks.NotificationRepository)
	}
	if len(extras) > 1 && extras[1] != nil {
		sender = extras[1].(*mockNotificationSender)
	}
	return service.NewSurveyService(repo, notifRepo, sender)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestSurveyService_Create_DefaultsStatus(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	s := &domain.Survey{Title: strPtr("Customer Feedback")}

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
	s := &domain.Survey{Title: strPtr("Customer Feedback")}

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
	expected := &domain.Survey{ID: id, Title: strPtr("Customer Feedback")}

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
	filter := domain.SurveyListFilter{Keyword: "TEST", Status: intPtr(domain.SurveyStatusActive), PublishType: intPtr(domain.SurveyPublishPromo)}
	expected := []*domain.Survey{{Title: strPtr("Survey A")}}

	repo.On("List", ctx, venueID, filter, p).Return(expected, int64(1), nil)

	got, total, err := svc.List(ctx, venueID, filter, p)
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

func TestSurveyService_SubmitVenueResponse_Success(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	surveyID := uuid.New()
	questionID := uuid.New()
	r := &domain.SurveyResponse{
		SurveyID: surveyID,
		Answers:  []*domain.SurveyAnswer{{QuestionID: questionID, AnswerText: strPtr("hello")}},
	}
	repo.On("FindByID", ctx, surveyID).Return(&domain.Survey{
		ID:      surveyID,
		VenueID: &venueID,
		Status:  domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{ID: questionID, SurveyID: surveyID, QuestionType: "paragraph", IsRequired: true},
		},
	}, nil)

	repo.On("CreateResponse", ctx, r).Return(nil)

	err := svc.SubmitVenueResponse(ctx, venueID, r)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitPublicResponse_ValidatesSurveyAndAnswers(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	questionID := uuid.New()
	optionID := uuid.New()
	resp := &domain.SurveyResponse{
		SurveyID: surveyID,
		Answers: []*domain.SurveyAnswer{
			{QuestionID: questionID, OptionID: &optionID},
		},
	}
	survey := &domain.Survey{
		ID:     surveyID,
		Status: domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{
				ID:           questionID,
				SurveyID:     surveyID,
				QuestionType: "single_choice",
				IsRequired:   true,
				Options: []*domain.Option{
					{ID: optionID, QuestionID: questionID},
				},
			},
		},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)
	repo.On("CreateResponse", ctx, resp).Return(nil)

	err := svc.SubmitPublicResponse(ctx, resp)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitPublicResponse_RejectsMissingRequiredAnswer(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	survey := &domain.Survey{
		ID:     surveyID,
		Status: domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{
				ID:           uuid.New(),
				SurveyID:     surveyID,
				QuestionType: "paragraph",
				IsRequired:   true,
			},
		},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)

	err := svc.SubmitPublicResponse(ctx, &domain.SurveyResponse{SurveyID: surveyID})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitPublicResponse_RejectsQuestionOutsideSurvey(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	otherQuestionID := uuid.New()
	survey := &domain.Survey{
		ID:     surveyID,
		Status: domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{
				ID:           uuid.New(),
				SurveyID:     surveyID,
				QuestionType: "paragraph",
			},
		},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)

	err := svc.SubmitPublicResponse(ctx, &domain.SurveyResponse{
		SurveyID: surveyID,
		Answers:  []*domain.SurveyAnswer{{QuestionID: otherQuestionID, AnswerText: strPtr("hello")}},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitPublicResponse_RejectsOptionOutsideQuestion(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	questionID := uuid.New()
	validOptionID := uuid.New()
	invalidOptionID := uuid.New()
	survey := &domain.Survey{
		ID:     surveyID,
		Status: domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{
				ID:           questionID,
				SurveyID:     surveyID,
				QuestionType: "single_choice",
				Options: []*domain.Option{
					{ID: validOptionID, QuestionID: questionID},
				},
			},
		},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)

	err := svc.SubmitPublicResponse(ctx, &domain.SurveyResponse{
		SurveyID: surveyID,
		Answers:  []*domain.SurveyAnswer{{QuestionID: questionID, OptionID: &invalidOptionID}},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitPublicResponse_RejectsInvalidTextOptionCombination(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	questionID := uuid.New()
	survey := &domain.Survey{
		ID:     surveyID,
		Status: domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{
				ID:           questionID,
				SurveyID:     surveyID,
				QuestionType: "paragraph",
			},
		},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)

	err := svc.SubmitPublicResponse(ctx, &domain.SurveyResponse{
		SurveyID: surveyID,
		Answers:  []*domain.SurveyAnswer{{QuestionID: questionID}},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitPublicResponse_AllowsRepeatedPayloads(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	questionID := uuid.New()
	survey := &domain.Survey{
		ID:     surveyID,
		Status: domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{
				ID:           questionID,
				SurveyID:     surveyID,
				QuestionType: "paragraph",
				IsRequired:   true,
			},
		},
	}
	resp1 := &domain.SurveyResponse{
		SurveyID: surveyID,
		Answers:  []*domain.SurveyAnswer{{QuestionID: questionID, AnswerText: strPtr("first")}},
	}
	resp2 := &domain.SurveyResponse{
		SurveyID: surveyID,
		Answers:  []*domain.SurveyAnswer{{QuestionID: questionID, AnswerText: strPtr("first")}},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil).Twice()
	repo.On("CreateResponse", ctx, resp1).Return(nil).Once()
	repo.On("CreateResponse", ctx, resp2).Return(nil).Once()

	require.NoError(t, svc.SubmitPublicResponse(ctx, resp1))
	require.NoError(t, svc.SubmitPublicResponse(ctx, resp2))
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitAppResponse_RequiresExternalID(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	err := svc.SubmitAppResponse(context.Background(), uuid.New(), &domain.SurveyResponse{
		SurveyID: uuid.New(),
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestSurveyService_SubmitAppResponse_RejectsSurveyFromDifferentVenue(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	otherVenueID := uuid.New()
	surveyID := uuid.New()
	questionID := uuid.New()
	survey := &domain.Survey{
		ID:      surveyID,
		VenueID: &otherVenueID,
		Status:  domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{ID: questionID, SurveyID: surveyID, QuestionType: "paragraph", IsRequired: true},
		},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)

	err := svc.SubmitAppResponse(ctx, venueID, &domain.SurveyResponse{
		SurveyID:   surveyID,
		ExternalID: strPtr("visitor-1"),
		Answers:    []*domain.SurveyAnswer{{QuestionID: questionID, AnswerText: strPtr("hello")}},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertExpectations(t)
}

func TestSurveyService_SubmitAppResponse_PassesExternalIDToRepository(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	surveyID := uuid.New()
	questionID := uuid.New()
	survey := &domain.Survey{
		ID:      surveyID,
		VenueID: &venueID,
		Status:  domain.SurveyStatusActive,
		Questions: []*domain.Question{
			{ID: questionID, SurveyID: surveyID, QuestionType: "paragraph", IsRequired: true},
		},
	}
	resp := &domain.SurveyResponse{
		SurveyID:   surveyID,
		ExternalID: strPtr("visitor-1"),
		Answers:    []*domain.SurveyAnswer{{QuestionID: questionID, AnswerText: strPtr("hello")}},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)
	repo.On("CreateResponse", ctx, resp).Return(nil)

	err := svc.SubmitAppResponse(ctx, venueID, resp)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSurveyService_Create_ActiveCMSCreatesImmediateNotification(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	notifRepo := &mocks.NotificationRepository{}
	sender := &mockNotificationSender{}
	svc := newTestSurveyService(repo, notifRepo, sender)

	ctx := context.Background()
	venueID := uuid.New()
	userID := uuid.New()
	s := &domain.Survey{
		Title:          strPtr("Customer Feedback"),
		Content:        strPtr("Please answer"),
		VenueID:        &venueID,
		CreatedBy:      &userID,
		Status:         domain.SurveyStatusActive,
		Source:         domain.SurveySourceCMS,
		PublishType:    domain.SurveyPublishBoth,
		StartDate:      ptrTime(time.Now().Add(-time.Minute)),
		SegmentFilters: []byte(`[{"key":"visitors","type":"text","value":"vip"}]`),
	}

	repo.On("Create", ctx, s).Return(nil)
	notifRepo.On("Create", ctx, mock.MatchedBy(func(n *domain.Notification) bool {
		return n.SurveyID != nil &&
			*n.SurveyID == s.ID &&
			n.SendType == domain.NotifTypeImmediate &&
			n.Kind == domain.NotifKindSurvey &&
			n.Title == s.Title
	})).Return(nil)
	sender.On("Send", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil)

	err := svc.Create(ctx, s)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	notifRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestSurveyService_Create_ActiveCMSRejectsUnsupportedSegmentFilters(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	notifRepo := &mocks.NotificationRepository{}
	sender := &mockNotificationSender{}
	svc := newTestSurveyService(repo, notifRepo, sender)

	ctx := context.Background()
	venueID := uuid.New()
	userID := uuid.New()
	s := &domain.Survey{
		Title:          strPtr("Customer Feedback"),
		Content:        strPtr("Please answer"),
		VenueID:        &venueID,
		CreatedBy:      &userID,
		Status:         domain.SurveyStatusActive,
		Source:         domain.SurveySourceCMS,
		PublishType:    domain.SurveyPublishBoth,
		StartDate:      ptrTime(time.Now().Add(-time.Minute)),
		SegmentFilters: []byte(`[{"foo":"bar"}]`),
	}

	repo.On("Create", ctx, s).Return(nil)

	err := svc.Create(ctx, s)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertExpectations(t)
	notifRepo.AssertNotCalled(t, "Create")
	sender.AssertNotCalled(t, "Send")
}

func TestSurveyService_Create_ActiveCMSCreatesScheduledNotification(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	notifRepo := &mocks.NotificationRepository{}
	sender := &mockNotificationSender{}
	svc := newTestSurveyService(repo, notifRepo, sender)

	ctx := context.Background()
	venueID := uuid.New()
	userID := uuid.New()
	startAt := time.Now().Add(time.Hour)
	s := &domain.Survey{
		Title:          strPtr("Customer Feedback"),
		Content:        strPtr("Please answer"),
		VenueID:        &venueID,
		CreatedBy:      &userID,
		Status:         domain.SurveyStatusActive,
		Source:         domain.SurveySourceCMS,
		PublishType:    domain.SurveyPublishBoth,
		StartDate:      &startAt,
		SegmentFilters: []byte(`[{"key":"visitors","type":"text","value":"vip"}]`),
	}

	repo.On("Create", ctx, s).Return(nil)
	notifRepo.On("Create", ctx, mock.MatchedBy(func(n *domain.Notification) bool {
		return n.SurveyID != nil &&
			*n.SurveyID == s.ID &&
			n.SendType == domain.NotifTypeScheduled &&
			n.ScheduledAt != nil &&
			n.ScheduledAt.Equal(startAt)
	})).Return(nil)

	err := svc.Create(ctx, s)
	require.NoError(t, err)
	sender.AssertNotCalled(t, "Send")
	repo.AssertExpectations(t)
	notifRepo.AssertExpectations(t)
}

func TestSurveyService_Update_TransitionToActiveCreatesNotification(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	notifRepo := &mocks.NotificationRepository{}
	sender := &mockNotificationSender{}
	svc := newTestSurveyService(repo, notifRepo, sender)

	ctx := context.Background()
	venueID := uuid.New()
	userID := uuid.New()
	surveyID := uuid.New()
	before := &domain.Survey{
		ID:        surveyID,
		VenueID:   &venueID,
		CreatedBy: &userID,
		Status:    domain.SurveyStatusInactive,
		Source:    domain.SurveySourceCMS,
	}
	after := &domain.Survey{
		ID:             surveyID,
		VenueID:        &venueID,
		CreatedBy:      &userID,
		Title:          strPtr("Activated survey"),
		Content:        strPtr("Now live"),
		Status:         domain.SurveyStatusActive,
		Source:         domain.SurveySourceCMS,
		StartDate:      ptrTime(time.Now().Add(-time.Minute)),
		SegmentFilters: []byte(`[{"key":"visitors","type":"text","value":"vip"}]`),
	}

	repo.On("FindByID", ctx, surveyID).Return(before, nil)
	repo.On("Update", ctx, after).Return(nil)
	notifRepo.On("Create", ctx, mock.MatchedBy(func(n *domain.Notification) bool {
		return n.SurveyID != nil && *n.SurveyID == surveyID
	})).Return(nil)
	sender.On("Send", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil)

	err := svc.Update(ctx, after)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	notifRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestSurveyService_ProcessScheduledTransitions_ActivatesAndCloses(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	notifRepo := &mocks.NotificationRepository{}
	sender := &mockNotificationSender{}
	svc := newTestSurveyService(repo, notifRepo, sender)

	ctx := context.Background()
	now := time.Now()
	venueID := uuid.New()
	userID := uuid.New()
	activate := &domain.Survey{
		ID:             uuid.New(),
		VenueID:        &venueID,
		CreatedBy:      &userID,
		Title:          strPtr("Activate me"),
		Content:        strPtr("Now"),
		Status:         domain.SurveyStatusInactive,
		Source:         domain.SurveySourceCMS,
		StartDate:      &now,
		SegmentFilters: []byte(`[{"key":"visitors","type":"text","value":"vip"}]`),
	}
	closeSurvey := &domain.Survey{
		ID:      uuid.New(),
		VenueID: &venueID,
		Status:  domain.SurveyStatusActive,
		EndDate: &now,
	}

	repo.On("ListDueActivation", ctx, now).Return([]*domain.Survey{activate}, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(s *domain.Survey) bool {
		return s.ID == activate.ID && s.Status == domain.SurveyStatusActive
	})).Return(nil).Once()
	notifRepo.On("Create", ctx, mock.MatchedBy(func(n *domain.Notification) bool {
		return n.SurveyID != nil && *n.SurveyID == activate.ID
	})).Return(nil)
	sender.On("Send", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil)
	repo.On("ListDueClosure", ctx, now).Return([]*domain.Survey{closeSurvey}, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(s *domain.Survey) bool {
		return s.ID == closeSurvey.ID && s.Status == domain.SurveyStatusClosed
	})).Return(nil).Once()

	activated, closed, err := svc.ProcessScheduledTransitions(ctx, now)
	require.NoError(t, err)
	assert.Equal(t, 1, activated)
	assert.Equal(t, 1, closed)
	repo.AssertExpectations(t)
	notifRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

// ── PBT: create defaults invariant ───────────────────────────────────────────

func TestSurveyService_Create_DefaultsInvariant(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		repo := &mocks.SurveyRepository{}
		svc := newTestSurveyService(repo)

		ctx := context.Background()
		s := &domain.Survey{
			Title:  func() *string { s := rapid.StringN(1, 50, 50).Draw(rt, "title"); return &s }(),
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

func ptrTime(v time.Time) *time.Time {
	return &v
}

func intPtr(v int) *int {
	return &v
}
