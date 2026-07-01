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
	"github.com/hhung06/digimap-backend/internal/service"
)

type productServiceStub struct {
	product            *domain.Product
	products           []*domain.Product
	attachments        []*domain.ProductAttachment
	deletedAttachments []uuid.UUID
	createdAttachments []*domain.ProductAttachment
}

func (s *productServiceStub) GetCategory(context.Context, uuid.UUID) (*domain.ProductCategory, error) {
	return nil, nil
}
func (s *productServiceStub) ListCategories(context.Context, uuid.UUID) ([]*domain.ProductCategory, error) {
	return nil, nil
}
func (s *productServiceStub) CreateCategory(context.Context, *domain.ProductCategory) error {
	return nil
}
func (s *productServiceStub) UpdateCategory(context.Context, *domain.ProductCategory) error {
	return nil
}
func (s *productServiceStub) DeleteCategory(context.Context, uuid.UUID) error { return nil }
func (s *productServiceStub) Get(_ context.Context, id uuid.UUID) (*domain.Product, error) {
	if s.product != nil {
		return s.product, nil
	}
	return &domain.Product{ID: id}, nil
}
func (s *productServiceStub) List(context.Context, uuid.UUID, domain.Pagination) ([]*domain.Product, int64, error) {
	return s.products, int64(len(s.products)), nil
}
func (s *productServiceStub) SearchByName(context.Context, uuid.UUID, string, int) ([]*domain.Product, error) {
	return nil, nil
}
func (s *productServiceStub) Create(context.Context, *domain.Product, []uuid.UUID) error {
	return nil
}
func (s *productServiceStub) Update(context.Context, *domain.Product, []uuid.UUID) error {
	return nil
}
func (s *productServiceStub) Delete(context.Context, uuid.UUID) error { return nil }
func (s *productServiceStub) ListAttachments(context.Context, uuid.UUID) ([]*domain.ProductAttachment, error) {
	return s.attachments, nil
}
func (s *productServiceStub) CreateAttachment(_ context.Context, a *domain.ProductAttachment) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.FileType == "" && a.File != nil {
		a.FileType = "image"
	}
	s.createdAttachments = append(s.createdAttachments, a)
	s.attachments = append(s.attachments, a)
	return nil
}
func (s *productServiceStub) DeleteAttachment(_ context.Context, _ uuid.UUID, id uuid.UUID) error {
	s.deletedAttachments = append(s.deletedAttachments, id)
	return nil
}

type productMediaURLStub struct {
	urls       map[string]*string
	uploadKeys []string
	uploads    []service.MediaUpload
}

func (s productMediaURLStub) Upload(context.Context, service.MediaTarget, service.MediaUpload) (string, error) {
	return "", nil
}
func (s productMediaURLStub) URL(_ context.Context, _ service.MediaTarget, key string) *string {
	return s.urls[key]
}
func (s productMediaURLStub) DeleteOwned(context.Context, service.MediaTarget, string) error {
	return nil
}

type productMediaStub struct {
	urls       map[string]*string
	uploadKeys []string
	uploads    []service.MediaUpload
	deleted    []string
}

func (s *productMediaStub) Upload(_ context.Context, _ service.MediaTarget, upload service.MediaUpload) (string, error) {
	s.uploads = append(s.uploads, upload)
	key := "local/media/products/" + uuid.New().String() + "/attachments/" + uuid.New().String() + ".png"
	if len(s.uploadKeys) > 0 {
		key = s.uploadKeys[0]
		s.uploadKeys = s.uploadKeys[1:]
	}
	return key, nil
}
func (s *productMediaStub) URL(_ context.Context, _ service.MediaTarget, key string) *string {
	if s.urls == nil {
		return nil
	}
	return s.urls[key]
}
func (s *productMediaStub) DeleteOwned(_ context.Context, _ service.MediaTarget, key string) error {
	s.deleted = append(s.deleted, key)
	return nil
}

func TestProductHandlerListReturnsTableFieldsAndAttachmentURL(t *testing.T) {
	productID := uuid.New()
	venueID := uuid.New()
	categoryID := uuid.New()
	locationName := "FutureTech Solutions"
	key := "local/media/products/" + productID.String() + "/attachments/" + uuid.New().String() + ".png"
	presignedURL := "https://bucket.s3.ap-northeast-1.amazonaws.com/" + key + "?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=access-key"
	name := "Smart Display"
	now := time.Now()
	svc := &productServiceStub{products: []*domain.Product{{
		ID:             productID,
		VenueID:        venueID,
		Name:           &name,
		LocationName:   &locationName,
		MainCategoryID: &categoryID,
		MainCategory:   &domain.ProductCategory{ID: categoryID, VenueID: venueID, Name: "Food & Beverage", Source: "internal", CreatedAt: now, UpdatedAt: now},
		Attachments:    []*domain.ProductAttachment{{ID: uuid.New(), ProductID: productID, FileType: "image", File: &key, CreatedAt: now}},
		CreatedAt:      now,
		UpdatedAt:      now,
	}}}
	h := newProductHandler(svc, nil, nil, productMediaURLStub{urls: map[string]*string{key: &presignedURL}})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: venueID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/products", nil)

	h.List(c)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "&X-Amz-Credential=")
	assert.NotContains(t, rec.Body.String(), `\u0026X-Amz-Credential=`)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	data := body["data"].([]any)
	require.Len(t, data, 1)
	first := data[0].(map[string]any)
	assert.Equal(t, locationName, first["location_name"])
	assert.Equal(t, "Food & Beverage", first["main_category_name"])
	assert.Equal(t, "Food & Beverage", first["main_category"].(map[string]any)["name"])
	attachments := first["attachments"].([]any)
	require.Len(t, attachments, 1)
	assert.Equal(t, key, attachments[0].(map[string]any)["file"])
	assert.Equal(t, presignedURL, attachments[0].(map[string]any)["file_url"])
	attachmentImage := first["attachment_image"].(map[string]any)
	assert.Equal(t, key, attachmentImage["file"])
	assert.Equal(t, presignedURL, attachmentImage["file_url"])
	imageAttachments := first["image_attachments"].([]any)
	require.Len(t, imageAttachments, 1)
	assert.Equal(t, key, imageAttachments[0].(map[string]any)["file"])
	assert.NotContains(t, first, "document_attachments")
}

func TestProductHandlerGetReturnsAttachmentsSplitByType(t *testing.T) {
	productID := uuid.New()
	docKey := "local/media/products/" + productID.String() + "/attachments/" + uuid.New().String() + ".pdf"
	imageKey := "local/media/products/" + productID.String() + "/attachments/" + uuid.New().String() + ".png"
	svc := &productServiceStub{product: &domain.Product{
		ID: productID,
		Attachments: []*domain.ProductAttachment{
			{ID: uuid.New(), ProductID: productID, FileType: "document", File: &docKey},
			{ID: uuid.New(), ProductID: productID, FileType: "image", File: &imageKey},
		},
	}}
	h := newProductHandler(svc, nil, nil, productMediaURLStub{})
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}, {Key: "productID", Value: productID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)

	h.Get(c)

	require.Equal(t, http.StatusOK, rec.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	attachments := data["attachments"].([]any)
	require.Len(t, attachments, 2)
	assert.Equal(t, docKey, attachments[0].(map[string]any)["file"])
	assert.Equal(t, imageKey, attachments[1].(map[string]any)["file"])
	imageAttachments := data["image_attachments"].([]any)
	require.Len(t, imageAttachments, 1)
	assert.Equal(t, imageKey, imageAttachments[0].(map[string]any)["file"])
	documentAttachments := data["document_attachments"].([]any)
	require.Len(t, documentAttachments, 1)
	assert.Equal(t, docKey, documentAttachments[0].(map[string]any)["file"])
}

func TestProductHandlerUpdateRejectsBase64AttachmentPayload(t *testing.T) {
	productID := uuid.New()
	oldAttachmentID := uuid.New()
	newKey := "data:image/png;base64,photo"
	svc := &productServiceStub{
		product:     &domain.Product{ID: productID},
		attachments: []*domain.ProductAttachment{{ID: oldAttachmentID, ProductID: productID, FileType: "image"}},
	}
	h := newProductHandler(svc, nil, nil)
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "productID", Value: productID.String()}}
	c.Request = httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewBufferString(`{"attachments":[{"title":"Photo","file_type":"image","file":"`+newKey+`"}]}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Update(c)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Empty(t, svc.deletedAttachments)
	require.Empty(t, svc.createdAttachments)
}

func TestProductHandlerUpdateMultipartUploadsAttachmentsAndKeepsSelected(t *testing.T) {
	productID := uuid.New()
	keepID := uuid.New()
	dropID := uuid.New()
	keepKey := "local/media/products/" + productID.String() + "/attachments/" + uuid.New().String() + ".png"
	dropKey := "local/media/products/" + productID.String() + "/attachments/" + uuid.New().String() + ".pdf"
	imageKey := "local/media/products/" + productID.String() + "/attachments/" + uuid.New().String() + ".png"
	docKey := "local/media/products/" + productID.String() + "/attachments/" + uuid.New().String() + ".pdf"
	svc := &productServiceStub{
		product: &domain.Product{ID: productID},
		attachments: []*domain.ProductAttachment{
			{ID: keepID, ProductID: productID, FileType: "image", File: &keepKey},
			{ID: dropID, ProductID: productID, FileType: "document", File: &dropKey},
		},
	}
	media := &productMediaStub{uploadKeys: []string{imageKey, docKey}}
	h := newProductHandler(svc, nil, nil, media)
	rec, c := multipartProductUpdateContext(t, productID, `{"name":"Updated","keep_attachment_ids":["`+keepID.String()+`"]}`, []productMultipartFile{
		{field: "image_attachments", name: "photo.png", contentType: "image/png", body: []byte("png")},
		{field: "document_attachments", name: "spec.pdf", contentType: "application/pdf", body: []byte("%PDF-1.4")},
	})

	h.Update(c)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, []uuid.UUID{dropID}, svc.deletedAttachments)
	require.Equal(t, []string{dropKey}, media.deleted)
	require.Len(t, svc.createdAttachments, 2)
	assert.Equal(t, "photo.png", *svc.createdAttachments[0].Title)
	assert.Equal(t, "image", svc.createdAttachments[0].FileType)
	assert.Equal(t, imageKey, *svc.createdAttachments[0].File)
	assert.Equal(t, "spec.pdf", *svc.createdAttachments[1].Title)
	assert.Equal(t, "document", svc.createdAttachments[1].FileType)
	assert.Equal(t, docKey, *svc.createdAttachments[1].File)
	require.Len(t, media.uploads, 2)
	assert.Equal(t, "image/png", media.uploads[0].ContentType)
	assert.Equal(t, "application/pdf", media.uploads[1].ContentType)
}

func TestProductHandlerUpdateMultipartIgnoresLegacyNonOwnedAttachmentOnDrop(t *testing.T) {
	productID := uuid.New()
	dropID := uuid.New()
	legacyValue := "data:application/pdf;base64,JVBERi0x"
	svc := &productServiceStub{
		product:     &domain.Product{ID: productID},
		attachments: []*domain.ProductAttachment{{ID: dropID, ProductID: productID, FileType: "document", File: &legacyValue}},
	}
	media := &productMediaStub{}
	h := newProductHandler(svc, nil, nil, media)
	rec, c := multipartProductUpdateContext(t, productID, `{"name":"Updated","keep_attachment_ids":[]}`, nil)

	h.Update(c)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, []uuid.UUID{dropID}, svc.deletedAttachments)
	require.Empty(t, media.deleted)
}

type productMultipartFile struct {
	field       string
	name        string
	contentType string
	body        []byte
}

func multipartProductUpdateContext(t *testing.T, productID uuid.UUID, data string, files []productMultipartFile) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.WriteField("data", data))
	for _, file := range files {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="`+file.field+`"; filename="`+file.name+`"`)
		header.Set("Content-Type", file.contentType)
		part, err := writer.CreatePart(header)
		require.NoError(t, err)
		_, err = io.Copy(part, bytes.NewReader(file.body))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())

	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "productID", Value: productID.String()}}
	c.Request = httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return rec, c
}
