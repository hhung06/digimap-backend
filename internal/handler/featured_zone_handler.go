package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type featuredZoneHandler struct {
	svc service.FeaturedZoneService
}

func newFeaturedZoneHandler(svc service.FeaturedZoneService) *featuredZoneHandler {
	return &featuredZoneHandler{svc: svc}
}

func (h *featuredZoneHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	zones, err := h.svc.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.FeaturedZoneResponse, len(zones))
	for i, z := range zones {
		items[i] = dto.FeaturedZoneToResponse(z)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *featuredZoneHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	var req dto.FeaturedZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	z := &domain.FeaturedZone{
		VenueID:     venueID,
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		SortIndex:   req.SortIndex,
		IsActive:    isActive,
	}
	if err := h.svc.Create(c.Request.Context(), z); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.FeaturedZoneToResponse(z)))
}

func (h *featuredZoneHandler) Update(c *gin.Context) {
	zoneID, err := uuid.Parse(c.Param("zoneID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid zone id"))
		return
	}
	var req dto.FeaturedZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	z, err := h.svc.Get(c.Request.Context(), zoneID)
	if err != nil {
		respondError(c, err)
		return
	}
	z.Name = req.Name
	z.Description = req.Description
	z.ImageURL = req.ImageURL
	z.SortIndex = req.SortIndex
	if req.IsActive != nil {
		z.IsActive = *req.IsActive
	}
	if err := h.svc.Update(c.Request.Context(), z); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.FeaturedZoneToResponse(z)))
}

func (h *featuredZoneHandler) Delete(c *gin.Context) {
	zoneID, err := uuid.Parse(c.Param("zoneID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid zone id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), zoneID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
