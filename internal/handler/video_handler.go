package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type videoHandler struct {
	svc service.VideoService
}

func newVideoHandler(svc service.VideoService) *videoHandler {
	return &videoHandler{svc: svc}
}

func (h *videoHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	videos, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.VideoResponse, len(videos))
	for i, v := range videos {
		items[i] = dto.VideoToResponse(v)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *videoHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("videoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid video id"))
		return
	}
	v, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VideoToResponse(v)))
}

func (h *videoHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	v := &domain.Video{
		VenueID: &venueID, Title: req.Title, Description: req.Description,
		URL: req.URL, Thumbnail: req.Thumbnail, Duration: req.Duration,
		Status: req.Status, PublishedAt: req.PublishedAt,
	}
	if err := h.svc.Create(c.Request.Context(), v); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.VideoToResponse(v)))
}

func (h *videoHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("videoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid video id"))
		return
	}
	var req dto.VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	v := &domain.Video{
		ID: id, Title: req.Title, Description: req.Description,
		URL: req.URL, Thumbnail: req.Thumbnail, Duration: req.Duration,
		Status: req.Status, PublishedAt: req.PublishedAt,
	}
	if err := h.svc.Update(c.Request.Context(), v); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VideoToResponse(v)))
}

func (h *videoHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("videoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid video id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
