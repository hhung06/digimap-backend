package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type articleServiceStub struct {
	article        *domain.Article
	updateCalled   bool
	mediaChange    service.ArticleMediaReplacement
	mediaFilenames []string
}

func (s *articleServiceStub) List(context.Context, uuid.UUID, domain.Pagination) ([]*domain.Article, int, error) {
	return nil, 0, nil
}

func (s *articleServiceStub) Get(_ context.Context, id uuid.UUID) (*domain.Article, error) {
	if s.article != nil {
		return s.article, nil
	}
	return &domain.Article{ID: id, Title: "Existing"}, nil
}

func (s *articleServiceStub) Create(context.Context, *domain.Article) error { return nil }

func (s *articleServiceStub) Update(_ context.Context, a *domain.Article) error {
	s.updateCalled = true
	s.article = a
	return nil
}

func (s *articleServiceStub) UpdateWithMedia(_ context.Context, a *domain.Article, replacement service.ArticleMediaReplacement) error {
	s.mediaChange = replacement
	s.mediaFilenames = nil
	for _, upload := range replacement.Uploads {
		s.mediaFilenames = append(s.mediaFilenames, upload.Filename)
		if upload.Reader != nil {
			_, _ = io.ReadAll(upload.Reader)
		}
	}
	a.Images = []*domain.ArticleImage{}
	for i, name := range s.mediaFilenames {
		a.Images = append(a.Images, &domain.ArticleImage{ArticleID: a.ID, Image: "key-" + name, SortOrder: i})
	}
	s.article = a
	return nil
}

func (s *articleServiceStub) Delete(context.Context, uuid.UUID) error { return nil }
func (s *articleServiceStub) CreateImage(context.Context, *domain.ArticleImage) error {
	return nil
}
func (s *articleServiceStub) DeleteImage(context.Context, uuid.UUID) error { return nil }

type articleMediaURLStub struct {
	urls map[string]*string
}

func (s articleMediaURLStub) Upload(context.Context, service.MediaTarget, service.MediaUpload) (string, error) {
	return "", nil
}
func (s articleMediaURLStub) URL(_ context.Context, _ service.MediaTarget, key string) *string {
	return s.urls[key]
}
func (s articleMediaURLStub) DeleteOwned(context.Context, service.MediaTarget, string) error {
	return nil
}

func TestArticleHandlerJSONUpdatePreservesImagesAndReturnsImageURL(t *testing.T) {
	articleID := uuid.New()
	oldKey := "develop/media/articles/" + articleID.String() + "/images/old.png"
	oldURL := "https://cdn.example.com/old.png"
	svc := &articleServiceStub{article: &domain.Article{
		ID:     articleID,
		Title:  "Existing",
		Images: []*domain.ArticleImage{{ArticleID: articleID, Image: oldKey}},
	}}
	h := newArticleHandler(svc, articleMediaURLStub{urls: map[string]*string{oldKey: &oldURL}})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "articleID", Value: articleID.String()}}
	c.Request = httptest.NewRequest(http.MethodPut, "/articles/"+articleID.String(), bytes.NewBufferString(`{"title":"Updated"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Update(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, svc.updateCalled)
	require.Len(t, svc.article.Images, 1)
	assert.Equal(t, oldKey, svc.article.Images[0].Image)
	body := articleResponseBody(t, rec.Body.Bytes())
	assert.Equal(t, "Updated", body["title"])
	images := body["images"].([]any)
	require.Len(t, images, 1)
	first := images[0].(map[string]any)
	assert.Equal(t, oldKey, first["image"])
	assert.Equal(t, oldURL, first["image_url"])
}

func TestArticleHandlerMultipartUpdateWithDateOnlyAndRepeatedImagesReplacesImages(t *testing.T) {
	articleID := uuid.New()
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	svc := &articleServiceStub{article: &domain.Article{ID: articleID, Title: "Existing"}}
	h := newArticleHandler(svc, nil)
	rec, c := multipartArticleUpdateContext(t, articleID, `{"title":"Updated","published_period_start":"2026-06-01"}`, []multipartTestFile{
		{field: "images", filename: "one.png", contentType: "image/png", body: "one"},
		{field: "images", filename: "two.png", contentType: "image/png", body: "two"},
	})

	h.Update(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.False(t, svc.updateCalled)
	assert.True(t, svc.mediaChange.Replace)
	assert.False(t, svc.mediaChange.Remove)
	assert.Equal(t, []string{"one.png", "two.png"}, svc.mediaFilenames)
	require.NotNil(t, svc.article.PublishedPeriodStart)
	assert.Equal(t, start, *svc.article.PublishedPeriodStart)
}

func TestArticleHandlerMultipartRemoveImagesClearsImages(t *testing.T) {
	articleID := uuid.New()
	svc := &articleServiceStub{article: &domain.Article{ID: articleID, Title: "Existing"}}
	h := newArticleHandler(svc, nil)
	rec, c := multipartArticleUpdateContext(t, articleID, `{"remove_images":true}`, nil)

	h.Update(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, svc.mediaChange.Replace)
	assert.True(t, svc.mediaChange.Remove)
	assert.Empty(t, svc.mediaChange.Uploads)
}

func TestArticleHandlerMultipartRejectsRemoveImagesWithUploads(t *testing.T) {
	articleID := uuid.New()
	svc := &articleServiceStub{article: &domain.Article{ID: articleID, Title: "Existing"}}
	h := newArticleHandler(svc, nil)
	rec, c := multipartArticleUpdateContext(t, articleID, `{"remove_images":true}`, []multipartTestFile{
		{field: "images", filename: "one.png", contentType: "image/png", body: "one"},
	})

	h.Update(c)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, svc.updateCalled)
	assert.False(t, svc.mediaChange.Replace)
}

func TestArticleHandlerGetPresignFailureLeavesImageURLNull(t *testing.T) {
	articleID := uuid.New()
	key := "develop/media/articles/" + articleID.String() + "/images/old.png"
	svc := &articleServiceStub{article: &domain.Article{
		ID:     articleID,
		Title:  "Existing",
		Images: []*domain.ArticleImage{{ArticleID: articleID, Image: key}},
	}}
	h := newArticleHandler(svc, articleMediaURLStub{urls: map[string]*string{}})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "articleID", Value: articleID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/articles/"+articleID.String(), nil)

	h.Get(c)

	require.Equal(t, http.StatusOK, rec.Code)
	body := articleResponseBody(t, rec.Body.Bytes())
	images := body["images"].([]any)
	require.Len(t, images, 1)
	first := images[0].(map[string]any)
	assert.Equal(t, key, first["image"])
	assert.Nil(t, first["image_url"])
}

func multipartArticleUpdateContext(t *testing.T, articleID uuid.UUID, data string, files []multipartTestFile) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("data", data))
	for _, file := range files {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="`+file.field+`"; filename="`+file.filename+`"`)
		header.Set("Content-Type", file.contentType)
		part, err := writer.CreatePart(header)
		require.NoError(t, err)
		_, err = part.Write([]byte(file.body))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())

	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "articleID", Value: articleID.String()}}
	c.Request = httptest.NewRequest(http.MethodPut, "/articles/"+articleID.String(), &body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return rec, c
}

func articleResponseBody(t *testing.T, body []byte) map[string]any {
	t.Helper()

	var resp dto.Response
	require.NoError(t, json.Unmarshal(body, &resp))
	data, ok := resp.Data.(map[string]any)
	require.True(t, ok, "data should be an object: %#v", resp.Data)
	return data
}
