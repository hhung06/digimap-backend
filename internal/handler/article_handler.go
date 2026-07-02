package handler

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type articleHandler struct {
	svc      service.ArticleService
	mediaSvc service.MediaService
}

func newArticleHandler(svc service.ArticleService, mediaSvc ...service.MediaService) *articleHandler {
	var media service.MediaService
	if len(mediaSvc) > 0 {
		media = mediaSvc[0]
	}
	return &articleHandler{svc: svc, mediaSvc: media}
}

// @Summary     List articles
// @Description List articles for a venue
// @Tags        articles
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
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
		items[i] = h.articleResponse(c.Request.Context(), a)
	}
	c.PureJSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
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
	c.PureJSON(http.StatusOK, dto.OK(h.articleResponse(c.Request.Context(), a)))
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
	if isMultipartRequest(c) {
		h.createMultipart(c, venueID)
		return
	}
	var req dto.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a := articleFromRequest(venueID, req)
	if err := h.svc.Create(c.Request.Context(), a); err != nil {
		respondError(c, err)
		return
	}
	c.PureJSON(http.StatusCreated, dto.OK(h.articleResponse(c.Request.Context(), a)))
}

func (h *articleHandler) createMultipart(c *gin.Context, venueID uuid.UUID) {
	var req dto.ArticleRequest
	if err := bindMultipartData(c, &req, 1<<20); err != nil {
		respondError(c, err)
		return
	}
	files, hasFiles, err := multipartFiles(c, "images")
	if err != nil {
		respondError(c, err)
		return
	}
	uploads, closers, err := articleMediaUploads(files)
	if err != nil {
		respondError(c, err)
		return
	}
	defer closeAll(closers)
	a := articleFromRequest(venueID, req)
	if hasFiles {
		err = h.svc.CreateWithMedia(c.Request.Context(), a, uploads)
	} else {
		err = h.svc.Create(c.Request.Context(), a)
	}
	if err != nil {
		respondError(c, err)
		return
	}
	c.PureJSON(http.StatusCreated, dto.OK(h.articleResponse(c.Request.Context(), a)))
}

// @Summary     Update article
// @Description Update an article (requires editor role)
// @Tags        articles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string               true "Venue ID"
// @Param       articleID path     string               true "Article ID"
// @Param       body      body     dto.UpdateArticleRequest true "Article details"
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
	if isMultipartRequest(c) {
		h.updateMultipart(c, id)
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
	c.PureJSON(http.StatusOK, dto.OK(h.articleResponse(c.Request.Context(), a)))
}

func (h *articleHandler) updateMultipart(c *gin.Context, id uuid.UUID) {
	var req dto.UpdateArticleRequest
	if err := bindMultipartData(c, &req, 1<<20); err != nil {
		respondError(c, err)
		return
	}
	files, hasFiles, err := multipartFiles(c, "images")
	if err != nil {
		respondError(c, err)
		return
	}
	removeImages := req.RemoveImages != nil && *req.RemoveImages
	if removeImages && hasFiles {
		respondError(c, domain.NewValidation(map[string]string{"images": "cannot upload images when remove_images is true"}))
		return
	}

	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(a)

	if removeImages {
		if err := h.svc.UpdateWithMedia(c.Request.Context(), a, service.ArticleMediaReplacement{Replace: true, Remove: true}); err != nil {
			respondError(c, err)
			return
		}
	} else if hasFiles || req.KeepImageIDs != nil {
		uploads, closers, err := articleMediaUploads(files)
		if err != nil {
			respondError(c, err)
			return
		}
		defer closeAll(closers)
		keepImageIDs := []uuid.UUID(nil)
		if req.KeepImageIDs != nil {
			keepImageIDs = *req.KeepImageIDs
		}
		if err := h.svc.UpdateWithMedia(c.Request.Context(), a, service.ArticleMediaReplacement{Replace: true, KeepImageIDs: keepImageIDs, Uploads: uploads}); err != nil {
			respondError(c, err)
			return
		}
	} else if err := h.svc.Update(c.Request.Context(), a); err != nil {
		respondError(c, err)
		return
	}

	c.PureJSON(http.StatusOK, dto.OK(h.articleResponse(c.Request.Context(), a)))
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
	resp := dto.ArticleImageToResponse(img)
	if h.mediaSvc != nil {
		target := service.MediaTarget{Entity: "articles", RecordID: img.ArticleID, Field: "images"}
		resp.ImageURL = h.mediaSvc.URL(c.Request.Context(), target, img.Image)
	}
	c.PureJSON(http.StatusCreated, dto.OK(resp))
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

func isMultipartRequest(c *gin.Context) bool {
	return strings.HasPrefix(strings.ToLower(c.GetHeader("Content-Type")), "multipart/form-data")
}

func articleFromRequest(venueID uuid.UUID, req dto.ArticleRequest) *domain.Article {
	return &domain.Article{
		VenueID: &venueID, ExternalID: req.ExternalID, LocationID: req.LocationID,
		Placement: req.Placement, Navigate: req.Navigate,
		Title: req.Title, Label: req.Label, Content: req.Content,
		Status:               req.Status,
		PublishedAt:          req.PublishedAt,
		PublishedPeriodStart: dto.DateToTimePtr(req.PublishedPeriodStart),
		PublishedPeriodEnd:   dto.DateToTimePtr(req.PublishedPeriodEnd),
		Localization:         req.Localization,
		RelatedProducts:      req.RelatedProducts,
	}
}

func articleMediaUploads(files []*multipart.FileHeader) ([]service.MediaUpload, []io.Closer, error) {
	uploads := make([]service.MediaUpload, 0, len(files))
	closers := make([]io.Closer, 0, len(files))
	for _, file := range files {
		upload, closer, err := mediaUpload(file)
		if err != nil {
			closeAll(closers)
			return nil, nil, err
		}
		uploads = append(uploads, upload)
		closers = append(closers, closer)
	}
	return uploads, closers, nil
}

func closeAll(closers []io.Closer) {
	for _, closer := range closers {
		_ = closer.Close()
	}
}

func (h *articleHandler) articleResponse(ctx context.Context, a *domain.Article) dto.ArticleResponse {
	resp := dto.ArticleToResponse(a)
	if h.mediaSvc == nil || a == nil {
		return resp
	}
	target := service.MediaTarget{Entity: "articles", RecordID: a.ID, Field: "images"}
	for i := range resp.Images {
		resp.Images[i].ImageURL = h.mediaSvc.URL(ctx, target, resp.Images[i].Image)
	}
	return resp
}
