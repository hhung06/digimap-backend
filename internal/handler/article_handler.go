package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type articleHandler struct {
	svc service.ArticleService
}

func newArticleHandler(svc service.ArticleService) *articleHandler {
	return &articleHandler{svc: svc}
}

func (h *articleHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	articles, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ArticleResponse, len(articles))
	for i, a := range articles {
		items[i] = dto.ArticleToResponse(a)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *articleHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("articleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid article id"))
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ArticleToResponse(a)))
}

func (h *articleHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a := &domain.Article{
		VenueID: &venueID, ExternalID: req.ExternalID, LocationID: req.LocationID,
		Placement: req.Placement, Navigate: req.Navigate,
		Title: req.Title, Label: req.Label, Content: req.Content,
		Status: req.Status,
		PublishedAt: req.PublishedAt,
		PublishedPeriodStart: req.PublishedPeriodStart,
		PublishedPeriodEnd:   req.PublishedPeriodEnd,
		Localization: req.Localization,
	}
	if err := h.svc.Create(c.Request.Context(), a); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ArticleToResponse(a)))
}

func (h *articleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("articleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid article id"))
		return
	}
	var req dto.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a := &domain.Article{
		ID: id, ExternalID: req.ExternalID, LocationID: req.LocationID,
		Placement: req.Placement, Navigate: req.Navigate,
		Title: req.Title, Label: req.Label, Content: req.Content,
		Status: req.Status,
		PublishedAt: req.PublishedAt,
		PublishedPeriodStart: req.PublishedPeriodStart,
		PublishedPeriodEnd:   req.PublishedPeriodEnd,
		Localization: req.Localization,
	}
	if err := h.svc.Update(c.Request.Context(), a); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ArticleToResponse(a)))
}

func (h *articleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("articleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid article id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *articleHandler) CreateImage(c *gin.Context) {
	articleID, err := uuid.Parse(c.Param("articleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid article id"))
		return
	}
	var req dto.ArticleImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	img := &domain.ArticleImage{
		ArticleID: articleID, Image: req.Image, SortOrder: req.SortOrder,
	}
	if err := h.svc.CreateImage(c.Request.Context(), img); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ArticleImageToResponse(img)))
}

func (h *articleHandler) DeleteImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("imageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid image id"))
		return
	}
	if err := h.svc.DeleteImage(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
