package handler

import (
	"net/http"

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
	enrichers  *enricher.Registry
}

func newProductHandler(svc service.ProductService, storageSvc service.StorageService, enrichers *enricher.Registry) *productHandler {
	return &productHandler{svc: svc, storageSvc: storageSvc, enrichers: enrichers}
}

// ── Product categories ────────────────────────────────────────────────────────

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
	extras, _ := h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceProduct)
	items := make([]any, len(products))
	for i, prod := range products {
		items[i] = enricher.MergeInto(dto.ProductToResponse(prod), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

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
	extras, _ := h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceProduct)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.ProductToResponse(prod), extras)))
}

func (h *productHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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
	c.JSON(http.StatusCreated, dto.OK(dto.ProductToResponse(prod)))
}

func (h *productHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid product id"))
		return
	}
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	prod := &domain.Product{
		ID: id, LocationID: req.LocationID, MainCategoryID: req.MainCategoryID,
		Image: req.Image, Name: req.Name, ExternalID: req.ExternalID, Size: req.Size,
		Price: req.Price, Country: req.OriginCountry, Expiration: req.Expiration,
		Description: req.Description, Custom: req.Custom, Localization: req.Localization,
		Source: req.Source,
	}
	if err := h.svc.Update(c.Request.Context(), prod, req.CategoryIDs); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ProductToResponse(prod)))
}

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
	att := &domain.ProductAttachment{
		ProductID: productID, Title: req.Title, File: req.File, SourceURL: req.SourceURL,
	}
	if err := h.svc.CreateAttachment(c.Request.Context(), att); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ProductAttachmentResponse{
		ID: att.ID, ProductID: att.ProductID, Title: att.Title,
		FileType: att.FileType, File: att.File, SourceURL: att.SourceURL,
		CreatedAt: att.CreatedAt,
	}))
}

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
	if err := h.svc.DeleteAttachment(c.Request.Context(), productID, attID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Storage ───────────────────────────────────────────────────────────────────

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
