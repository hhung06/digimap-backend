package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestPublicHandler(venueRepo *mocks.VenueRepository) *publicHandler {
	venueSvc := service.NewVenueService(venueRepo, nil, nil)
	return newPublicHandler(venueSvc, nil, nil, nil, nil)
}

type visitorSurveyServiceStub struct {
	req  service.VisitorSurveySubmission
	user *domain.AppUser
	err  error
}

func (s *visitorSurveyServiceStub) Submit(_ context.Context, req service.VisitorSurveySubmission) (*domain.AppUser, error) {
	s.req = req
	return s.user, s.err
}

func TestPublicHandler_VenueInformation_RequiresPublicKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := newTestPublicHandler(&mocks.VenueRepository{})
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	id := uuid.New()
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/public/v1/venues/"+id.String()+"/information", nil)

	h.VenueInformation(c)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "public_key is required")
}

func TestPublicHandler_VenueInformation_InvalidPublicKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	venueRepo := &mocks.VenueRepository{}
	h := newTestPublicHandler(venueRepo)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	id := uuid.New()
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/public/v1/venues/"+id.String()+"/information?public_key=wrong", nil)

	venueRepo.
		On("FindByID", mock.Anything, id).
		Return(&domain.Venue{ID: id, PublicKey: "expected"}, nil)

	h.VenueInformation(c)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid public_key")
	venueRepo.AssertExpectations(t)
}

func TestPublicHandler_VenueInformation_ValidPublicKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	venueRepo := &mocks.VenueRepository{}
	h := newTestPublicHandler(venueRepo)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	id := uuid.New()
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/public/v1/venues/"+id.String()+"/information?public_key=expected", nil)

	venue := &domain.Venue{ID: id, Name: "Expo", PublicKey: "expected"}
	venueRepo.
		On("FindByID", mock.Anything, id).
		Return(venue, nil)

	h.VenueInformation(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Expo")
	venueRepo.AssertExpectations(t)
}

func TestPublicHandler_SubmitVisitorSurvey_RequiresPublicKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &visitorSurveyServiceStub{err: domain.NewValidation(map[string]string{"public_key": "is required"})}
	h := newPublicHandler(nil, nil, nil, nil, stub)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/public/v1/visitor-surveys", bytes.NewBufferString(`{"visitor_type":1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.SubmitVisitorSurvey(c)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "public_key")
}

func TestPublicHandler_SubmitVisitorSurvey_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &visitorSurveyServiceStub{user: &domain.AppUser{ID: uuid.New()}}
	h := newPublicHandler(nil, nil, nil, nil, stub)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/public/v1/visitor-surveys", bytes.NewBufferString(`{"public_key":"pub","full_name":"John Doe","email":"john@example.com","phone_number":"+123","visitor_type":2,"business_name":"Acme","interests":[1,2,1],"other_interests":"other","is_consented":true}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", "test-agent")
	c.Request.RemoteAddr = "127.0.0.1:1234"

	h.SubmitVisitorSurvey(c)

	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "pub", stub.req.PublicKey)
	assert.Equal(t, "test-agent", stub.req.UserAgent)
	assert.Equal(t, "127.0.0.1", stub.req.IPAddress)
	assert.Contains(t, rec.Body.String(), stub.user.ID.String())
}
