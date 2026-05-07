package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestPublicHandler(venueRepo *mocks.VenueRepository) *publicHandler {
	venueSvc := service.NewVenueService(venueRepo, nil, nil)
	return newPublicHandler(venueSvc, nil, nil, nil)
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
