package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type levelBundleHandler struct {
	svc service.LevelBundleService
}

func newLevelBundleHandler(svc service.LevelBundleService) *levelBundleHandler {
	return &levelBundleHandler{svc: svc}
}

// @Summary     List level bundles
// @Description List level bundles for a snapshot (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string true "Venue ID"
// @Param       snapshotID path     string true "Snapshot ID"
// @Success     200        {object} dto.Response{data=[]dto.LevelBundleResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/snapshots/{snapshotID}/bundles [get]
func (h *levelBundleHandler) List(c *gin.Context) {
	snapshotID, err := uuid.Parse(c.Param("snapshotID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid snapshot id"))
		return
	}
	bundles, err := h.svc.ListBySnapshot(c.Request.Context(), snapshotID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LevelBundleResponse, len(bundles))
	for i, b := range bundles {
		items[i] = dto.LevelBundleToResponse(b)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Create level bundle
// @Description Create a level bundle within a snapshot (system admin only)
// @Tags        snapshots
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string                        true "Venue ID"
// @Param       snapshotID path     string                        true "Snapshot ID"
// @Param       body       body     dto.CreateLevelBundleRequest  true "Bundle details"
// @Success     201        {object} dto.Response{data=dto.LevelBundleResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/snapshots/{snapshotID}/bundles [post]
func (h *levelBundleHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	snapshotID, err := uuid.Parse(c.Param("snapshotID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid snapshot id"))
		return
	}
	var req dto.CreateLevelBundleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	b, err := h.svc.Create(c.Request.Context(), snapshotID, venueID, req.LevelID, req.Bundle)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LevelBundleToResponse(b)))
}

// @Summary     Delete level bundle
// @Description Delete a level bundle (system admin only)
// @Tags        snapshots
// @Produce     json
// @Security    BearerAuth
// @Param       id         path string true "Venue ID"
// @Param       snapshotID path string true "Snapshot ID"
// @Param       bundleID   path string true "Bundle ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/snapshots/{snapshotID}/bundles/{bundleID} [delete]
func (h *levelBundleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("bundleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid bundle id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
