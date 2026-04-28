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

// @Summary     List snapshots
// @Description List snapshots for a venue (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Failure     403       {object} dto.Response
// @Router      /venues/{id}/snapshots [get]
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

// @Summary     Get snapshot
// @Description Get a snapshot by ID (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string true "Venue ID"
// @Param       snapshotID path     string true "Snapshot ID"
// @Success     200        {object} dto.Response{data=dto.SnapshotResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Failure     404        {object} dto.Response
// @Router      /venues/{id}/snapshots/{snapshotID} [get]
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

// @Summary     Get latest published snapshot
// @Description Get the most recent published snapshot for a venue (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=dto.SnapshotResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /venues/{id}/snapshots/recent [get]
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

// @Summary     Create snapshot draft
// @Description Create a new snapshot draft for a venue (system admin only)
// @Tags        snapshots
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                     true "Venue ID"
// @Param       body body     dto.CreateSnapshotRequest  true "Snapshot details"
// @Success     201  {object} dto.Response{data=dto.SnapshotResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/snapshots [post]
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

// @Summary     Publish snapshot
// @Description Publish a snapshot (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string true "Venue ID"
// @Param       snapshotID path     string true "Snapshot ID"
// @Success     200        {object} dto.Response{data=dto.SnapshotResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Failure     404        {object} dto.Response
// @Router      /venues/{id}/snapshots/{snapshotID}/publish [post]
func (h *snapshotHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("snapshotID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid snapshot id"))
		return
	}
	s, err := h.svc.Publish(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SnapshotToResponse(s)))
}

// @Summary     Revert snapshot
// @Description Revert a snapshot to published state (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string true "Venue ID"
// @Param       snapshotID path     string true "Snapshot ID"
// @Success     200        {object} dto.Response{data=dto.SnapshotResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Failure     404        {object} dto.Response
// @Router      /venues/{id}/snapshots/{snapshotID}/revert [post]
func (h *snapshotHandler) Revert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("snapshotID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid snapshot id"))
		return
	}
	s, err := h.svc.Revert(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SnapshotToResponse(s)))
}

// @Summary     Auto-publish snapshot
// @Description Create a new auto snapshot and immediately publish it (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=dto.SnapshotResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/snapshots/auto-publish [post]
func (h *snapshotHandler) AutoPublish(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	userID := middleware.GetUserID(c)
	s, err := h.svc.AutoPublish(c.Request.Context(), venueID, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SnapshotToResponse(s)))
}

// @Summary     Delete snapshot
// @Description Delete a snapshot (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id         path string true "Venue ID"
// @Param       snapshotID path string true "Snapshot ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/snapshots/{snapshotID} [delete]
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
