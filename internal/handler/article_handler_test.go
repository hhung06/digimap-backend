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
	articles       []*domain.Article
	createCalled   bool
	updateCalled   bool
	mediaChange    service.ArticleMediaReplacement
	mediaFilenames []string
}

func (s *articleServiceStub) List(context.Context, uuid.UUID, domain.Pagination) ([]*domain.Article, int, error) {
	return s.articles, len(s.articles), nil
}

func (s *articleServiceStub) Get(_ context.Context, id uuid.UUID) (*domain.Article, error) {
	if s.article != nil {
		return s.article, nil
	}
	return &domain.Article{ID: id, Title: "Existing"}, nil
}

func (s *articleServiceStub) Create(_ context.Context, a *domain.Article) error {
	s.createCalled = true
	s.article = a
	return nil
}

func (s *articleServiceStub) CreateWithMedia(_ context.Context, a *domain.Article, uploads []service.MediaUpload) error {
	s.createCalled = true
	s.mediaFilenames = nil
	for _, upload := range uploads {
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

func TestArticleHandlerListReturnsImageURL(t *testing.T) {
	venueID := uuid.New()
	articleID := uuid.New()
	locationID := uuid.New()
	firstProductID := uuid.New()
	secondProductID := uuid.New()
	key := "develop/media/articles/" + articleID.String() + "/images/cover.png"
	imageURL := "https://cdn.example.com/cover.png"
	svc := &articleServiceStub{articles: []*domain.Article{{
		ID:         articleID,
		VenueID:    &venueID,
		LocationID: &locationID,
		Location:   &domain.ArticleLocation{ID: locationID, Name: "Premium Lounge"},
		RelatedProducts: []uuid.UUID{
			firstProductID,
			secondProductID,
		},
		Title:  "Listed",
		Images: []*domain.ArticleImage{{ArticleID: articleID, Image: key}},
	}}}
	h := newArticleHandler(svc, articleMediaURLStub{urls: map[string]*string{key: &imageURL}})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: venueID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/venues/"+venueID.String()+"/articles", nil)

	h.List(c)

	require.Equal(t, http.StatusOK, rec.Code)
	var body dto.Response
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	items := body.Data.([]any)
	require.Len(t, items, 1)
	firstArticle := items[0].(map[string]any)
	location := firstArticle["location"].(map[string]any)
	assert.Equal(t, locationID.String(), location["id"])
	assert.Equal(t, "Premium Lounge", location["name"])
	relatedProducts := firstArticle["related_products"].([]any)
	assert.Equal(t, []any{firstProductID.String(), secondProductID.String()}, relatedProducts)
	images := firstArticle["images"].([]any)
	require.Len(t, images, 1)
	firstImage := images[0].(map[string]any)
	assert.Equal(t, key, firstImage["image"])
	assert.Equal(t, imageURL, firstImage["image_url"])
}

func TestArticleHandlerGetDoesNotHTMLEscapePresignedImageURL(t *testing.T) {
	articleID := uuid.New()
	key := "local/media/articles/" + articleID.String() + "/images/old.png"
	presignedURL := "https://bucket.s3.ap-northeast-1.amazonaws.com/" + key + "?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=access-key"
	svc := &articleServiceStub{article: &domain.Article{
		ID:     articleID,
		Title:  "Existing",
		Images: []*domain.ArticleImage{{ArticleID: articleID, Image: key}},
	}}
	h := newArticleHandler(svc, articleMediaURLStub{urls: map[string]*string{key: &presignedURL}})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "articleID", Value: articleID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/articles/"+articleID.String(), nil)

	h.Get(c)

	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "&X-Amz-Credential=")
	assert.NotContains(t, body, `\u0026X-Amz-Credential=`)
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

func TestArticleHandlerMultipartUpdateKeepsSelectedImagesWithoutUploads(t *testing.T) {
	articleID := uuid.New()
	keepID := uuid.New()
	dropID := uuid.New()
	svc := &articleServiceStub{article: &domain.Article{
		ID: articleID,
		Images: []*domain.ArticleImage{
			{ID: keepID, ArticleID: articleID, Image: "keep.png"},
			{ID: dropID, ArticleID: articleID, Image: "drop.png"},
		},
	}}
	h := newArticleHandler(svc, nil)
	rec, c := multipartArticleUpdateContext(t, articleID, `{"keep_image_ids":["`+keepID.String()+`"]}`, nil)

	h.Update(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.False(t, svc.updateCalled)
	assert.True(t, svc.mediaChange.Replace)
	assert.False(t, svc.mediaChange.Remove)
	assert.Equal(t, []uuid.UUID{keepID}, svc.mediaChange.KeepImageIDs)
	assert.Empty(t, svc.mediaChange.Uploads)
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
