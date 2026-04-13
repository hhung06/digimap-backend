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
