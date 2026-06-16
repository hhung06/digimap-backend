package service_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func TestVisitorSurveySubmissionService_Submit_CreatesFreshAppUser(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	appUserRepo := &mocks.AppUserRepository{}
	svc := service.NewVisitorSurveySubmissionService(venueRepo, appUserRepo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	venueRepo.On("FindByPublicKey", ctx, "pub-key").Return(&domain.Venue{ID: venueID}, nil)
	appUserRepo.On("Create", ctx, mock.MatchedBy(func(u *domain.AppUser) bool {
		if u.VenueID == nil || *u.VenueID != venueID {
			return false
		}
		if u.Source != "visitor_survey" || u.FirstName != "John" || u.LastName != "Doe" {
			return false
		}
		if u.Phone != "+123" || u.Email != "john@example.com" || u.BusinessName != "Acme" {
			return false
		}
		if u.IPAddress != "127.0.0.1" || u.UserAgent != "test-agent" || !u.IsConsented {
			return false
		}
		if u.VisitorType == nil || *u.VisitorType != 2 {
			return false
		}
		if u.OtherInterests != "other" {
			return false
		}
		var interests []int
		if err := json.Unmarshal(u.Interests, &interests); err != nil {
			return false
		}
		return assert.ObjectsAreEqual([]int{1, 2}, interests)
	})).Return(nil)

	u, err := svc.Submit(ctx, service.VisitorSurveySubmission{
		PublicKey:      "pub-key",
		FullName:       "John Doe",
		Email:          "john@example.com",
		PhoneNumber:    "+123",
		VisitorType:    2,
		BusinessName:   "Acme",
		Interests:      []int{1, 2, 1},
		OtherInterests: "other",
		IsConsented:    true,
		IPAddress:      "127.0.0.1",
		UserAgent:      "test-agent",
	})
	require.NoError(t, err)
	require.NotNil(t, u)
	venueRepo.AssertExpectations(t)
	appUserRepo.AssertExpectations(t)
}

func TestVisitorSurveySubmissionService_Submit_RequiresPublicKey(t *testing.T) {
	svc := service.NewVisitorSurveySubmissionService(&mocks.VenueRepository{}, &mocks.AppUserRepository{}, nil)

	_, err := svc.Submit(context.Background(), service.VisitorSurveySubmission{
		VisitorType: 1,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestVisitorSurveySubmissionService_Submit_RequiresVisitorType(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	appUserRepo := &mocks.AppUserRepository{}
	svc := service.NewVisitorSurveySubmissionService(venueRepo, appUserRepo, nil)

	_, err := svc.Submit(context.Background(), service.VisitorSurveySubmission{
		PublicKey: "pub-key",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}
