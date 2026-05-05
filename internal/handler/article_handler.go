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

// @Summary     List articles
// @Description List articles for a venue
// @Tags        articles
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/articles [get]
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

// @Summary     Get article
// @Description Get an article by ID
// @Tags        articles
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true "Venue ID"
// @Param       articleID path     string true "Article ID"
// @Success     200       {object} dto.Response{data=dto.ArticleResponse}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Failure     404       {object} dto.Response
// @Router      /venues/{id}/articles/{articleID} [get]
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

// @Summary     Create article
// @Description Create a new article (requires editor role)
// @Tags        articles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string               true "Venue ID"
// @Param       body body     dto.ArticleRequest   true "Article details"
// @Success     201  {object} dto.Response{data=dto.ArticleResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/articles [post]
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

// @Summary     Update article
// @Description Update an article (requires editor role)
// @Tags        articles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string               true "Venue ID"
// @Param       articleID path     string               true "Article ID"
// @Param       body      body     dto.ArticleRequest   true "Article details"
// @Success     200       {object} dto.Response{data=dto.ArticleResponse}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Failure     403       {object} dto.Response
// @Router      /venues/{id}/articles/{articleID} [put]
func (h *articleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("articleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid article id"))
		return
	}
	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(a)
	if err := h.svc.Update(c.Request.Context(), a); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ArticleToResponse(a)))
}

// @Summary     Delete article
// @Description Delete an article (requires editor role)
// @Tags        articles
// @Produce     json
// @Security    BearerAuth
// @Param       id        path string true "Venue ID"
// @Param       articleID path string true "Article ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/articles/{articleID} [delete]
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

// @Summary     Add article image
// @Description Add an image to an article (requires editor role)
// @Tags        articles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string                    true "Venue ID"
// @Param       articleID path     string                    true "Article ID"
// @Param       body      body     dto.ArticleImageRequest   true "Image details"
// @Success     201       {object} dto.Response{data=dto.ArticleImageResponse}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Failure     403       {object} dto.Response
// @Router      /venues/{id}/articles/{articleID}/images [post]
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

// @Summary     Delete article image
// @Description Delete an image from an article (requires editor role)
// @Tags        articles
// @Produce     json
// @Security    BearerAuth
// @Param       id        path string true "Venue ID"
// @Param       articleID path string true "Article ID"
// @Param       imageID   path string true "Image ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/articles/{articleID}/images/{imageID} [delete]
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
