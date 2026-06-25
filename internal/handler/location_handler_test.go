package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/service"
)

type locationCategoryServiceStub struct{}

func (s locationCategoryServiceStub) Get(context.Context, uuid.UUID) (*domain.LocationCategory, error) {
	return nil, nil
}
func (s locationCategoryServiceStub) List(context.Context, uuid.UUID) ([]*domain.LocationCategory, error) {
	return nil, nil
}
func (s locationCategoryServiceStub) Create(context.Context, *domain.LocationCategory) error {
	return nil
}
func (s locationCategoryServiceStub) Update(context.Context, *domain.LocationCategory) error {
	return nil
}
func (s locationCategoryServiceStub) Delete(context.Context, uuid.UUID) error { return nil }

type locationServiceStub struct {
	location *domain.Location
}

func (s *locationServiceStub) Get(_ context.Context, id uuid.UUID) (*domain.Location, error) {
	if s.location != nil {
		return s.location, nil
	}
	return &domain.Location{ID: id}, nil
}
func (s *locationServiceStub) List(context.Context, uuid.UUID, *int, *uuid.UUID, domain.Pagination) ([]*domain.Location, int64, error) {
	return nil, 0, nil
}
func (s *locationServiceStub) SearchByName(context.Context, uuid.UUID, string, int) ([]*domain.Location, error) {
	return nil, nil
}
func (s *locationServiceStub) Create(context.Context, *domain.Location) error { return nil }
func (s *locationServiceStub) Update(context.Context, *domain.Location, []uuid.UUID) error {
	return nil
}
func (s *locationServiceStub) Delete(context.Context, uuid.UUID) error { return nil }
func (s *locationServiceStub) Duplicate(context.Context, uuid.UUID) (*domain.Location, error) {
	return nil, nil
}
func (s *locationServiceStub) SetTop(context.Context, uuid.UUID, bool, *int) error { return nil }
func (s *locationServiceStub) SetTopWithMedia(context.Context, uuid.UUID, bool, *int, *service.MediaUpload) error {
	return nil
}
func (s *locationServiceStub) GeoSearch(context.Context, float64, float64, float64, *uuid.UUID) ([]*domain.Location, error) {
	return nil, nil
}
func (s *locationServiceStub) ListImages(context.Context, uuid.UUID) ([]*domain.LocationImage, error) {
	return nil, nil
}
func (s *locationServiceStub) CreateImage(context.Context, *domain.LocationImage) error {
	return nil
}
func (s *locationServiceStub) DeleteImage(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}

type locationMediaURLStub struct {
	urls map[string]*string
}

func (s locationMediaURLStub) Upload(context.Context, service.MediaTarget, service.MediaUpload) (string, error) {
	return "", nil
}
func (s locationMediaURLStub) URL(_ context.Context, _ service.MediaTarget, key string) *string {
	return s.urls[key]
}
func (s locationMediaURLStub) DeleteOwned(context.Context, service.MediaTarget, string) error {
	return nil
}

func TestLocationHandlerGetReturnsTopLogoURL(t *testing.T) {
	locationID := uuid.New()
	venueID := uuid.New()
	key := "local/media/locations/" + locationID.String() + "/top_logo/marker.png"
	presignedURL := "https://bucket.s3.ap-northeast-1.amazonaws.com/" + key + "?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=access-key"
	svc := &locationServiceStub{location: &domain.Location{
		ID:      locationID,
		VenueID: venueID,
		TopLogo: key,
	}}
	h := newLocationHandler(locationCategoryServiceStub{}, svc, nil, locationMediaURLStub{
		urls: map[string]*string{key: &presignedURL},
	})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{
		{Key: "id", Value: venueID.String()},
		{Key: "locationID", Value: locationID.String()},
	}
	c.Request = httptest.NewRequest(http.MethodGet, "/locations/"+locationID.String(), nil)

	h.GetLocation(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "&X-Amz-Credential=")
	assert.NotContains(t, rec.Body.String(), `\u0026X-Amz-Credential=`)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	assert.Equal(t, key, data["top_logo"])
	assert.Equal(t, presignedURL, data["top_logo_url"])
}
