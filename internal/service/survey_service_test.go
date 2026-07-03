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

// ── Duplicate ─────────────────────────────────────────────────────────────────

func TestSurveyService_Duplicate_DeepCopiesQuestionsAndOptions(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	questionID := uuid.New()
	original := &domain.Survey{
		ID: surveyID, Title: strPtr("Feedback"), Status: domain.SurveyStatusActive,
		Source: domain.SurveySourceCMS, PublishType: domain.SurveyPublishInApp,
		Questions: []*domain.Question{{
			ID: questionID, SurveyID: surveyID, QuestionNumber: 1,
			QuestionType: "single_choice", QuestionText: "Rate us", IsOther: true,
			Options: []*domain.Option{
				{ID: uuid.New(), QuestionID: questionID, OptionNumber: 1, OptionText: "Good"},
				{ID: uuid.New(), QuestionID: questionID, OptionNumber: 2, OptionText: "Bad"},
			},
		}},
	}

	repo.On("FindByID", ctx, surveyID).Return(original, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Survey")).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*domain.Survey).ID = uuid.New()
	})
	repo.On("CreateQuestion", ctx, mock.AnythingOfType("*domain.Question")).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*domain.Question).ID = uuid.New()
	})
	repo.On("CreateOption", ctx, mock.AnythingOfType("*domain.Option")).Return(nil).Times(2)

	copied, err := svc.Duplicate(ctx, surveyID)
	require.NoError(t, err)
	assert.Equal(t, "Feedback (Copy)", *copied.Title)
	assert.NotEqual(t, surveyID, copied.ID)
	assert.Equal(t, domain.SurveyStatusActive, copied.Status, "status copied verbatim, not reset")
	require.Len(t, copied.Questions, 1)
	assert.NotEqual(t, questionID, copied.Questions[0].ID)
	assert.Equal(t, "Rate us", copied.Questions[0].QuestionText)
	assert.Len(t, copied.Questions[0].Options, 2)
	repo.AssertExpectations(t)
	// No notification repo wired: Duplicate must not attempt to notify even
	// though the copied survey is Active.
}

// ── Stats ─────────────────────────────────────────────────────────────────────

func TestSurveyService_Stats_AggregatesChoicesAndParagraphs(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	choiceQ := uuid.New()
	paraQ := uuid.New()
	optA := uuid.New()
	optB := uuid.New()
	survey := &domain.Survey{
		ID: surveyID, Title: strPtr("Feedback"),
		Questions: []*domain.Question{
			{ID: choiceQ, QuestionNumber: 1, QuestionType: "multiple_choice", IsOther: true,
				Options: []*domain.Option{
					{ID: optA, OptionNumber: 1, OptionText: "A"},
					{ID: optB, OptionNumber: 2, OptionText: "B"},
				}},
			{ID: paraQ, QuestionNumber: 2, QuestionType: "paragraph"},
		},
	}
	resp1 := uuid.New()
	resp2 := uuid.New()
	responses := []*domain.SurveyResponse{
		{ID: resp1, SurveyID: surveyID, Answers: []*domain.SurveyAnswer{
			{ResponseID: resp1, QuestionID: choiceQ, OptionID: &optA},
			{ResponseID: resp1, QuestionID: choiceQ, OptionID: &optB},
			{ResponseID: resp1, QuestionID: paraQ, AnswerText: strPtr("great")},
		}},
		{ID: resp2, SurveyID: surveyID, Answers: []*domain.SurveyAnswer{
			{ResponseID: resp2, QuestionID: choiceQ, OptionID: &optA},
			{ResponseID: resp2, QuestionID: choiceQ, AnswerText: strPtr("something else")},
		}},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)
	repo.On("ListResponsesWithAnswers", ctx, surveyID, (*string)(nil)).Return(responses, nil)

	stats, err := svc.Stats(ctx, surveyID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), stats.TotalResponses)
	require.Len(t, stats.Questions, 2)

	choice := stats.Questions[0]
	assert.Equal(t, 2, choice.TotalPeople)
	assert.Equal(t, 3, choice.TotalChoices)
	require.Len(t, choice.Options, 3)                             // A, B, Other
	assert.Equal(t, 2, choice.Options[0].Count)                   // A
	assert.Equal(t, 100.0, choice.Options[0].PercentageByPeople)  // 2/2 people
	assert.Equal(t, 66.67, choice.Options[0].PercentageByChoices) // 2/3 choices
	assert.Equal(t, 1, choice.Options[1].Count)                   // B
	assert.True(t, choice.Options[2].IsOther)
	assert.Equal(t, []string{"something else"}, choice.Options[2].OtherTexts)

	para := stats.Questions[1]
	assert.Equal(t, 1, para.TotalPeople)
	assert.Equal(t, 0, para.TotalChoices)
	assert.Equal(t, []string{"great"}, para.Texts)
	repo.AssertExpectations(t)
}

func TestSurveyService_Stats_EmptySurvey(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	optA := uuid.New()
	survey := &domain.Survey{
		ID: surveyID,
		Questions: []*domain.Question{
			{ID: uuid.New(), QuestionNumber: 1, QuestionType: "single_choice",
				Options: []*domain.Option{{ID: optA, OptionText: "A"}}},
		},
	}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)
	repo.On("ListResponsesWithAnswers", ctx, surveyID, (*string)(nil)).Return([]*domain.SurveyResponse{}, nil)

	stats, err := svc.Stats(ctx, surveyID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.TotalResponses)
	require.Len(t, stats.Questions, 1)
	assert.Equal(t, 0, stats.Questions[0].Options[0].Count)
	assert.Equal(t, 0.0, stats.Questions[0].Options[0].PercentageByPeople, "empty denominators must not divide by zero")
}

// ── ExportResponses ───────────────────────────────────────────────────────────

func TestSurveyService_ExportResponses_BuildsCSVRows(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	choiceQ := uuid.New()
	paraQ := uuid.New()
	optA := uuid.New()
	survey := &domain.Survey{
		ID: surveyID,
		Questions: []*domain.Question{
			{ID: choiceQ, QuestionNumber: 1, QuestionType: "multiple_choice", QuestionText: "Pick", IsOther: true,
				Options: []*domain.Option{{ID: optA, OptionNumber: 1, OptionText: "A"}}},
			{ID: paraQ, QuestionNumber: 2, QuestionType: "paragraph", QuestionText: "Comment"},
		},
	}
	submitted := time.Date(2026, 7, 1, 9, 30, 0, 0, time.UTC)
	responses := []*domain.SurveyResponse{
		{ID: uuid.New(), SurveyID: surveyID, ExternalID: strPtr("EXT-1"), SubmittedAt: submitted,
			Answers: []*domain.SurveyAnswer{
				{QuestionID: choiceQ, OptionID: &optA},
				{QuestionID: choiceQ, AnswerText: strPtr("custom")},
				{QuestionID: paraQ, AnswerText: strPtr("nice")},
			}},
		{ID: uuid.New(), SurveyID: surveyID, SubmittedAt: submitted}, // anonymous, no answers
	}
	users := []*domain.AppUser{{
		ID: uuid.New(), ExternalID: "EXT-1",
		FirstName: "太郎", LastName: "山田", FirstNameEn: "Taro", LastNameEn: "Yamada",
	}}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)
	repo.On("ListResponsesWithAnswers", ctx, surveyID, (*string)(nil)).Return(responses, nil)
	repo.On("FindAppUsersByExternalIDs", ctx, []string{"EXT-1"}).Return(users, nil)

	rows, err := svc.ExportResponses(ctx, surveyID, "")
	require.NoError(t, err)
	require.Len(t, rows, 3)
	assert.Equal(t, []string{"No", "User External ID", "User Name (JP)", "User Name (EN)", "Submitted At", "Q1: Pick", "Q2: Comment"}, rows[0])
	assert.Equal(t, []string{"1", "EXT-1", "山田 太郎", "Taro Yamada", "2026-07-01 09:30:00", "A; [Other] custom", "nice"}, rows[1])
	assert.Equal(t, []string{"2", "", "", "", "2026-07-01 09:30:00", "", ""}, rows[2])
	repo.AssertExpectations(t)
}

func TestSurveyService_ExportResponses_FiltersByExternalID(t *testing.T) {
	repo := &mocks.SurveyRepository{}
	svc := newTestSurveyService(repo)

	ctx := context.Background()
	surveyID := uuid.New()
	survey := &domain.Survey{ID: surveyID}

	repo.On("FindByID", ctx, surveyID).Return(survey, nil)
	repo.On("ListResponsesWithAnswers", ctx, surveyID, strPtr("EXT-9")).Return([]*domain.SurveyResponse{}, nil)
	repo.On("FindAppUsersByExternalIDs", ctx, []string{}).Return(nil, nil)

	rows, err := svc.ExportResponses(ctx, surveyID, " EXT-9 ")
	require.NoError(t, err)
	require.Len(t, rows, 1, "header only when no responses")
	repo.AssertExpectations(t)
}
