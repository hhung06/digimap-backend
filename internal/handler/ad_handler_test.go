package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/service"
)

type adServiceStub struct {
	ad            *domain.Advertisement
	mediaFilename string
	replacement   service.AdvertisementMediaReplacement
}

func (s *adServiceStub) List(context.Context, uuid.UUID, domain.Pagination) ([]*domain.Advertisement, int, error) {
	if s.ad == nil {
		return nil, 0, nil
	}
	return []*domain.Advertisement{s.ad}, 1, nil
}
func (s *adServiceStub) Get(_ context.Context, id uuid.UUID) (*domain.Advertisement, error) {
	if s.ad != nil {
		return s.ad, nil
	}
	return &domain.Advertisement{ID: id, Type: "banner", Status: "draft", Placement: "home_screen"}, nil
}
func (s *adServiceStub) Create(_ context.Context, a *domain.Advertisement) error {
	s.ad = a
	return nil
}
func (s *adServiceStub) CreateWithMedia(_ context.Context, a *domain.Advertisement, upload *service.MediaUpload) error {
	if upload != nil {
		s.mediaFilename = upload.Filename
		if upload.Reader != nil {
			_, _ = io.ReadAll(upload.Reader)
		}
		key := "local/media/ads/" + a.ID.String() + "/content_image/upload.png"
		a.ContentImage = &key
	}
	s.ad = a
	return nil
}
func (s *adServiceStub) Update(_ context.Context, a *domain.Advertisement) error {
	s.ad = a
	return nil
}
func (s *adServiceStub) UpdateWithMedia(_ context.Context, a *domain.Advertisement, replacement service.AdvertisementMediaReplacement) error {
	s.replacement = replacement
	if replacement.Upload != nil {
		s.mediaFilename = replacement.Upload.Filename
		if replacement.Upload.Reader != nil {
			_, _ = io.ReadAll(replacement.Upload.Reader)
		}
		key := "local/media/ads/" + a.ID.String() + "/content_image/new.png"
		a.ContentImage = &key
	}
	s.ad = a
	return nil
}
func (s *adServiceStub) Delete(context.Context, uuid.UUID) error  { return nil }
func (s *adServiceStub) Publish(context.Context, uuid.UUID) error { return nil }

type adMediaURLStub struct {
	urls map[string]*string
}

func (s adMediaURLStub) Upload(context.Context, service.MediaTarget, service.MediaUpload) (string, error) {
	return "", nil
}
func (s adMediaURLStub) URL(_ context.Context, _ service.MediaTarget, key string) *string {
	return s.urls[key]
}
func (s adMediaURLStub) DeleteOwned(context.Context, service.MediaTarget, string) error {
	return nil
}

func TestAdHandlerGetReturnsContentImageAndPresignedURL(t *testing.T) {
	adID := uuid.New()
	venueID := uuid.New()
	key := "local/media/ads/" + adID.String() + "/content_image/upload.png"
	presignedURL := "https://bucket.s3.ap-northeast-1.amazonaws.com/" + key + "?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=access-key"
	svc := &adServiceStub{ad: &domain.Advertisement{
		ID:           adID,
		VenueID:      &venueID,
		Type:         "banner",
		Status:       "published",
		Placement:    "home_screen",
		ContentImage: &key,
	}}
	h := newAdHandler(svc, enricher.NewRegistry(nil), adMediaURLStub{urls: map[string]*string{key: &presignedURL}})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{
		{Key: "id", Value: venueID.String()},
		{Key: "adID", Value: adID.String()},
	}
	c.Request = httptest.NewRequest(http.MethodGet, "/ads/"+adID.String(), nil)

	h.Get(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "&X-Amz-Credential=")
	assert.NotContains(t, rec.Body.String(), `\u0026X-Amz-Credential=`)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	assert.Equal(t, key, data["content_image"])
	assert.Equal(t, presignedURL, data["content_image_url"])
}

func TestAdHandlerMultipartUpdateUsesContentImageFile(t *testing.T) {
	adID := uuid.New()
	oldKey := "local/media/ads/" + adID.String() + "/content_image/old.png"
	svc := &adServiceStub{ad: &domain.Advertisement{
		ID:           adID,
		Type:         "banner",
		Status:       "draft",
		Placement:    "home_screen",
		ContentImage: &oldKey,
	}}
	h := newAdHandler(svc, enricher.NewRegistry(nil))
	rec, c := multipartAdUpdateContext(t, adID, `{"placement":"map_banner"}`, "banner.png")

	h.Update(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "banner.png", svc.mediaFilename)
	assert.Equal(t, oldKey, svc.replacement.OldKey)
	require.NotNil(t, svc.ad.ContentImage)
	assert.Contains(t, *svc.ad.ContentImage, "/content_image/new.png")
}

func multipartAdUpdateContext(t *testing.T, adID uuid.UUID, data string, filename string) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("data", data))
	part, err := writer.CreateFormFile("content_image", filename)
	require.NoError(t, err)
	_, err = part.Write([]byte("image"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "adID", Value: adID.String()}}
	c.Request = httptest.NewRequest(http.MethodPut, "/ads/"+adID.String(), &body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return rec, c
}
