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

type adHandler struct {
	svc      service.AdvertisementService
	enrichers *enricher.Registry
}

func newAdHandler(svc service.AdvertisementService, enrichers *enricher.Registry) *adHandler {
	return &adHandler{svc: svc, enrichers: enrichers}
}

func (h *adHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	ads, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceAd)
	items := make([]any, len(ads))
	for i, a := range ads {
		items[i] = enricher.MergeInto(dto.AdvertisementToResponse(a), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *adHandler) Get(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	id, err := uuid.Parse(c.Param("adID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid ad id"))
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceAd)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.AdvertisementToResponse(a), extras)))
}

func (h *adHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.AdvertisementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a := &domain.Advertisement{
		VenueID: &venueID, LocationID: req.LocationID,
		Type: req.Type, Status: "draft", Navigate: req.Navigate,
		ContentImageURL: req.ContentImageURL, ContentCTAURL: req.ContentCTAURL,
		Placement: req.Placement,
		SizeWidth: req.SizeWidth, SizeHeight: req.SizeHeight,
		RewardType: req.RewardType, RewardAmount: req.RewardAmount,
		DisplayDuration: req.DisplayDuration,
		StartAt: req.StartAt, EndAt: req.EndAt,
	}
	if err := h.svc.Create(c.Request.Context(), a); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.AdvertisementToResponse(a)))
}

func (h *adHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("adID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid ad id"))
		return
	}
	var req dto.AdvertisementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a := &domain.Advertisement{
		ID: id, LocationID: req.LocationID,
		Type: req.Type, Navigate: req.Navigate,
		ContentImageURL: req.ContentImageURL, ContentCTAURL: req.ContentCTAURL,
		Placement: req.Placement,
		SizeWidth: req.SizeWidth, SizeHeight: req.SizeHeight,
		RewardType: req.RewardType, RewardAmount: req.RewardAmount,
		DisplayDuration: req.DisplayDuration,
		StartAt: req.StartAt, EndAt: req.EndAt,
	}
	if err := h.svc.Update(c.Request.Context(), a); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.AdvertisementToResponse(a)))
}

func (h *adHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("adID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid ad id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *adHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("adID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid ad id"))
		return
	}
	if err := h.svc.Publish(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
}
