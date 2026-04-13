package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type snapshotHandler struct {
	svc service.SnapshotService
}

func newSnapshotHandler(svc service.SnapshotService) *snapshotHandler {
	return &snapshotHandler{svc: svc}
}

func (h *snapshotHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	snapshots, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.SnapshotResponse, len(snapshots))
	for i, s := range snapshots {
		items[i] = dto.SnapshotToResponse(s)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *snapshotHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("snapshotID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid snapshot id"))
		return
	}
	s, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SnapshotToResponse(s)))
}

func (h *snapshotHandler) LatestPublished(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	s, err := h.svc.LatestPublished(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SnapshotToResponse(s)))
}

func (h *snapshotHandler) CreateDraft(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	userID := middleware.GetUserID(c)
	s, err := h.svc.CreateDraft(c.Request.Context(), venueID, userID, req.Bundle)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.SnapshotToResponse(s)))
}

func (h *snapshotHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("snapshotID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid snapshot id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
