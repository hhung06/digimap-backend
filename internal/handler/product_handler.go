package handler

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/service"
)

type productHandler struct {
	svc        service.ProductService
	storageSvc service.StorageService
	mediaSvc   service.MediaService
	enrichers  *enricher.Registry
}

func newProductHandler(svc service.ProductService, storageSvc service.StorageService, enrichers *enricher.Registry, mediaSvc ...service.MediaService) *productHandler {
	var media service.MediaService
	if len(mediaSvc) > 0 {
		media = mediaSvc[0]
	}
	return &productHandler{svc: svc, storageSvc: storageSvc, enrichers: enrichers, mediaSvc: media}
}

// ── Product categories ────────────────────────────────────────────────────────

// @Summary     List product categories
// @Description List all product categories for a venue
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=[]dto.ProductCategoryResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /venues/{id}/product-categories [get]
func (h *productHandler) ListCategories(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	cats, err := h.svc.ListCategories(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ProductCategoryResponse, len(cats))
	for i, cat := range cats {
		items[i] = dto.ProductCategoryToResponse(cat)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Create product category
// @Description Create a new product category (requires editor role)
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                      true "Venue ID"
// @Param       body body     dto.ProductCategoryRequest  true "Category details"
// @Success     201  {object} dto.Response{data=dto.ProductCategoryResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/product-categories [post]
func (h *productHandler) CreateCategory(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.ProductCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cat := &domain.ProductCategory{
		VenueID: venueID, ExternalID: req.ExternalID, Name: req.Name,
		Source: req.Source, Localization: req.Localization,
	}
	if err := h.svc.CreateCategory(c.Request.Context(), cat); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ProductCategoryToResponse(cat)))
}

// @Summary     Update product category
// @Description Update a product category (requires editor role)
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path     string                     true "Venue ID"
// @Param       catID path     string                     true "Category ID"
// @Param       body  body     dto.ProductCategoryRequest true "Category details"
// @Success     200   {object} dto.Response{data=dto.ProductCategoryResponse}
// @Failure     400   {object} dto.Response
// @Failure     401   {object} dto.Response
// @Failure     403   {object} dto.Response
// @Router      /venues/{id}/product-categories/{catID} [put]
func (h *productHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("catID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid category id"))
		return
	}
	var req dto.ProductCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cat := &domain.ProductCategory{
		ID: id, ExternalID: req.ExternalID, Name: req.Name,
		Source: req.Source, Localization: req.Localization,
	}
	if err := h.svc.UpdateCategory(c.Request.Context(), cat); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ProductCategoryToResponse(cat)))
}

// @Summary     Delete product category
// @Description Delete a product category (requires editor role)
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id    path string true "Venue ID"
// @Param       catID path string true "Category ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/product-categories/{catID} [delete]
func (h *productHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("catID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid category id"))
		return
	}
	if err := h.svc.DeleteCategory(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Products ──────────────────────────────────────────────────────────────────

// @Summary     List products
// @Description List products for a venue
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/products [get]
func (h *productHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	products, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	ctx := c.Request.Context()
	var extras map[string]any
	if h.enrichers != nil {
		extras, _ = h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceProduct)
	}
	items := make([]any, len(products))
	for i, prod := range products {
		items[i] = enricher.MergeInto(h.productListResponse(ctx, prod), extras)
	}
	c.PureJSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

// @Summary     Get product
// @Description Get a product by ID
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true "Venue ID"
// @Param       productID path     string true "Product ID"
// @Success     200       {object} dto.Response{data=dto.ProductResponse}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Failure     404       {object} dto.Response
// @Router      /venues/{id}/products/{productID} [get]
func (h *productHandler) Get(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	id, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid product id"))
		return
	}
	prod, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	ctx := c.Request.Context()
	var extras map[string]any
	if h.enrichers != nil {
		extras, _ = h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceProduct)
	}
	c.PureJSON(http.StatusOK, dto.OK(enricher.MergeInto(h.productResponse(ctx, prod), extras)))
}

// @Summary     Create product
// @Description Create a new product (requires editor role)
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                    true "Venue ID"
// @Param       body body     dto.CreateProductRequest  true "Product details"
// @Success     201  {object} dto.Response{data=dto.ProductResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/products [post]
func (h *productHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateProductRequest
	if isMultipartRequest(c) {
		if err := bindMultipartData(c, &req, 1<<20); err != nil {
			respondError(c, err)
			return
		}
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	prod := &domain.Product{
		VenueID: venueID, LocationID: req.LocationID, MainCategoryID: req.MainCategoryID,
		Image: req.Image, Name: req.Name, ExternalID: req.ExternalID, Size: req.Size,
		Price: req.Price, Country: req.OriginCountry, Expiration: req.Expiration,
		Description: req.Description, Custom: req.Custom, Localization: req.Localization,
		Source: req.Source,
	}
	if err := h.svc.Create(c.Request.Context(), prod, req.CategoryIDs); err != nil {
		respondError(c, err)
		return
	}
	if isMultipartRequest(c) {
		if err := h.replaceProductMultipartAttachments(c.Request.Context(), c, prod.ID, nil); err != nil {
			respondError(c, err)
			return
		}
	} else if err := h.replaceProductAttachments(c.Request.Context(), prod.ID, req.Attachments); err != nil {
		respondError(c, err)
		return
	}
	prod.Attachments, _ = h.svc.ListAttachments(c.Request.Context(), prod.ID)
	c.PureJSON(http.StatusCreated, dto.OK(h.productResponse(c.Request.Context(), prod)))
}

// @Summary     Update product
// @Description Update a product (requires editor role)
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string                   true "Venue ID"
// @Param       productID path     string                   true "Product ID"
// @Param       body      body     dto.UpdateProductRequest true "Product details"
// @Success     200       {object} dto.Response{data=dto.ProductResponse}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Failure     403       {object} dto.Response
// @Failure     404       {object} dto.Response
// @Router      /venues/{id}/products/{productID} [put]
func (h *productHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid product id"))
		return
	}
	var req dto.UpdateProductRequest
	if isMultipartRequest(c) {
		if err := bindMultipartData(c, &req, 1<<20); err != nil {
			respondError(c, err)
			return
		}
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	prod, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(prod)
	categoryIDs := req.CategoryIDs
	if err := h.svc.Update(c.Request.Context(), prod, categoryIDs); err != nil {
		respondError(c, err)
		return
	}
	if isMultipartRequest(c) {
		if err := h.replaceProductMultipartAttachments(c.Request.Context(), c, prod.ID, req.KeepAttachmentIDs); err != nil {
			respondError(c, err)
			return
		}
	} else if req.Attachments != nil {
		if err := h.replaceProductAttachments(c.Request.Context(), prod.ID, req.Attachments); err != nil {
			respondError(c, err)
			return
		}
	}
	prod.Attachments, _ = h.svc.ListAttachments(c.Request.Context(), prod.ID)
	c.PureJSON(http.StatusOK, dto.OK(h.productResponse(c.Request.Context(), prod)))
}

// @Summary     Delete product
// @Description Delete a product (requires editor role)
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id        path string true "Venue ID"
// @Param       productID path string true "Product ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/products/{productID} [delete]
func (h *productHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid product id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Attachments ───────────────────────────────────────────────────────────────

// @Summary     Add product attachment
// @Description Add an attachment to a product (requires editor role)
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string                       true "Venue ID"
// @Param       productID path     string                       true "Product ID"
// @Param       body      body     dto.ProductAttachmentRequest true "Attachment details"
// @Success     201       {object} dto.Response{data=dto.ProductAttachmentResponse}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Failure     403       {object} dto.Response
// @Router      /venues/{id}/products/{productID}/attachments [post]
func (h *productHandler) CreateAttachment(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid product id"))
		return
	}
	var req dto.ProductAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	if req.File != nil && strings.HasPrefix(strings.TrimSpace(*req.File), "data:") {
		respondError(c, domain.NewValidation(map[string]string{"file": "base64 attachments must be uploaded as files"}))
		return
	}
	att := &domain.ProductAttachment{
		ProductID: productID, Title: req.Title, File: req.File, SourceURL: req.SourceURL,
	}
	if err := h.svc.CreateAttachment(c.Request.Context(), att); err != nil {
		respondError(c, err)
		return
	}
	resp := dto.ProductAttachmentResponse{
		ID: att.ID, ProductID: att.ProductID, Title: att.Title,
		FileType: att.FileType, File: att.File, SourceURL: att.SourceURL,
		CreatedAt: att.CreatedAt,
	}
	if h.mediaSvc != nil && att.File != nil {
		target := service.MediaTarget{Entity: "products", RecordID: att.ProductID, Field: "attachments"}
		resp.FileURL = h.mediaSvc.URL(c.Request.Context(), target, *att.File)
	}
	c.PureJSON(http.StatusCreated, dto.OK(resp))
}

// @Summary     Delete product attachment
// @Description Delete a product attachment (requires editor role)
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id        path string true "Venue ID"
// @Param       productID path string true "Product ID"
// @Param       attID     path string true "Attachment ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/products/{productID}/attachments/{attID} [delete]
func (h *productHandler) DeleteAttachment(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid product id"))
		return
	}
	attID, err := uuid.Parse(c.Param("attID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid attachment id"))
		return
	}
	if h.mediaSvc != nil {
		attachments, err := h.svc.ListAttachments(c.Request.Context(), productID)
		if err != nil {
			respondError(c, err)
			return
		}
		target := service.MediaTarget{Entity: "products", RecordID: productID, Field: "attachments"}
		for _, attachment := range attachments {
			if attachment.ID == attID && attachment.File != nil && *attachment.File != "" {
				if err := h.deleteOwnedProductAttachmentMedia(c.Request.Context(), target, *attachment.File); err != nil {
					respondError(c, err)
					return
				}
				break
			}
		}
	}
	if err := h.svc.DeleteAttachment(c.Request.Context(), productID, attID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Storage ───────────────────────────────────────────────────────────────────

// @Summary     Presign upload URL
// @Description Generate a pre-signed S3 upload URL
// @Tags        storage
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.PresignUploadRequest true "Upload request"
// @Success     200  {object} dto.Response{data=dto.PresignUploadResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /storage/presign-upload [post]
func (h *productHandler) PresignUpload(c *gin.Context) {
	var req dto.PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	url, err := h.storageSvc.PresignUpload(c.Request.Context(), req.Key, req.ContentType)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PresignUploadResponse{URL: url, Key: req.Key}))
}

func (h *productHandler) productResponse(ctx context.Context, p *domain.Product) dto.ProductResponse {
	if p == nil {
		return dto.ProductResponse{}
	}
	resp := dto.ProductToResponse(p)
	if h.mediaSvc == nil {
		return resp
	}
	target := service.MediaTarget{Entity: "products", RecordID: p.ID, Field: "attachments"}
	for i := range resp.Attachments {
		if resp.Attachments[i].File != nil {
			resp.Attachments[i].FileURL = h.mediaSvc.URL(ctx, target, *resp.Attachments[i].File)
			if resp.AttachmentImage != nil && resp.AttachmentImage.ID == resp.Attachments[i].ID {
				resp.AttachmentImage.FileURL = resp.Attachments[i].FileURL
			}
		}
	}
	for i := range resp.ImageAttachments {
		if resp.ImageAttachments[i].File != nil {
			resp.ImageAttachments[i].FileURL = h.mediaSvc.URL(ctx, target, *resp.ImageAttachments[i].File)
		}
	}
	for i := range resp.DocumentAttachments {
		if resp.DocumentAttachments[i].File != nil {
			resp.DocumentAttachments[i].FileURL = h.mediaSvc.URL(ctx, target, *resp.DocumentAttachments[i].File)
		}
	}
	return resp
}

func (h *productHandler) productListResponse(ctx context.Context, p *domain.Product) dto.ProductResponse {
	resp := h.productResponse(ctx, p)
	if resp.AttachmentImage != nil {
		resp.Attachments = []dto.ProductAttachmentResponse{*resp.AttachmentImage}
		resp.ImageAttachments = []dto.ProductAttachmentResponse{*resp.AttachmentImage}
	} else {
		resp.Attachments = nil
		resp.ImageAttachments = nil
	}
	resp.DocumentAttachments = nil
	return resp
}

func (h *productHandler) replaceProductAttachments(ctx context.Context, productID uuid.UUID, attachments []dto.ProductAttachmentRequest) error {
	for _, attachment := range attachments {
		if attachment.File != nil && strings.HasPrefix(strings.TrimSpace(*attachment.File), "data:") {
			return domain.NewValidation(map[string]string{"attachments": "base64 attachments must be uploaded as files"})
		}
	}
	existing, err := h.svc.ListAttachments(ctx, productID)
	if err != nil {
		return err
	}
	for _, attachment := range existing {
		if err := h.svc.DeleteAttachment(ctx, productID, attachment.ID); err != nil {
			return err
		}
	}
	for _, attachment := range attachments {
		if attachment.File == nil || *attachment.File == "" {
			continue
		}
		next := &domain.ProductAttachment{
			ProductID: productID,
			Title:     attachment.Title,
			FileType:  attachment.FileType,
			File:      attachment.File,
			SourceURL: attachment.SourceURL,
		}
		if err := h.svc.CreateAttachment(ctx, next); err != nil {
			return err
		}
	}
	return nil
}

func (h *productHandler) replaceProductMultipartAttachments(ctx context.Context, c *gin.Context, productID uuid.UUID, keepAttachmentIDs []uuid.UUID) error {
	if h.mediaSvc == nil {
		return domain.NewValidation(map[string]string{"attachments": "media service is not configured"})
	}
	existing, err := h.svc.ListAttachments(ctx, productID)
	if err != nil {
		return err
	}
	keep := map[uuid.UUID]struct{}{}
	for _, id := range keepAttachmentIDs {
		keep[id] = struct{}{}
	}
	target := service.MediaTarget{Entity: "products", RecordID: productID, Field: "attachments"}
	for _, attachment := range existing {
		if _, ok := keep[attachment.ID]; ok {
			continue
		}
		if attachment.File != nil && *attachment.File != "" {
			if err := h.deleteOwnedProductAttachmentMedia(ctx, target, *attachment.File); err != nil {
				return err
			}
		}
		if err := h.svc.DeleteAttachment(ctx, productID, attachment.ID); err != nil {
			return err
		}
	}
	imageFiles, _, err := multipartFiles(c, "image_attachments")
	if err != nil {
		return err
	}
	documentFiles, _, err := multipartFiles(c, "document_attachments")
	if err != nil {
		return err
	}
	if err := h.createProductAttachmentUploads(ctx, productID, target, "image", imageFiles); err != nil {
		return err
	}
	return h.createProductAttachmentUploads(ctx, productID, target, "document", documentFiles)
}

func (h *productHandler) createProductAttachmentUploads(ctx context.Context, productID uuid.UUID, target service.MediaTarget, fileType string, files []*multipart.FileHeader) error {
	for _, file := range files {
		if fileType == "document" && !isSupportedProductDocument(file) {
			return domain.NewValidation(map[string]string{"document_attachments": "only pdf and docx files are supported"})
		}
		upload, closer, err := mediaUpload(file)
		if err != nil {
			return err
		}
		defer closer.Close()
		key, err := h.mediaSvc.Upload(ctx, target, upload)
		if err != nil {
			return err
		}
		title := file.Filename
		next := &domain.ProductAttachment{
			ProductID: productID,
			Title:     &title,
			FileType:  fileType,
			File:      &key,
		}
		if err := h.svc.CreateAttachment(ctx, next); err != nil {
			_ = h.mediaSvc.DeleteOwned(ctx, target, key)
			return err
		}
	}
	return nil
}

func isSupportedProductDocument(file *multipart.FileHeader) bool {
	contentType := strings.ToLower(strings.TrimSpace(file.Header.Get("Content-Type")))
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return (ext == ".pdf" && contentType == "application/pdf") ||
		(ext == ".docx" && contentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
}

func (h *productHandler) deleteOwnedProductAttachmentMedia(ctx context.Context, target service.MediaTarget, key string) error {
	key = strings.TrimSpace(key)
	if h.mediaSvc == nil || key == "" || strings.HasPrefix(key, "data:") || strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return nil
	}
	err := h.mediaSvc.DeleteOwned(ctx, target, key)
	if errors.Is(err, domain.ErrValidation) {
		return nil
	}
	return err
}
