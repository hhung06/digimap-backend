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

// @Summary     List advertisements
// @Description List advertisements for a venue
// @Tags        ads
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/ads [get]
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

// @Summary     Get advertisement
// @Description Get an advertisement by ID
// @Tags        ads
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string true "Venue ID"
// @Param       adID path     string true "Advertisement ID"
// @Success     200  {object} dto.Response{data=dto.AdvertisementResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     404  {object} dto.Response
// @Router      /venues/{id}/ads/{adID} [get]
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

// @Summary     Create advertisement
// @Description Create a new advertisement (requires editor role)
// @Tags        ads
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                      true "Venue ID"
// @Param       body body     dto.AdvertisementRequest    true "Advertisement details"
// @Success     201  {object} dto.Response{data=dto.AdvertisementResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/ads [post]
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

// @Summary     Update advertisement
// @Description Update an advertisement (requires editor role)
// @Tags        ads
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                     true "Venue ID"
// @Param       adID path     string                     true "Advertisement ID"
// @Param       body body     dto.AdvertisementRequest   true "Advertisement details"
// @Success     200  {object} dto.Response{data=dto.AdvertisementResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/ads/{adID} [put]
func (h *adHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("adID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid ad id"))
		return
	}
	var req dto.UpdateAdvertisementRequest
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
	c.JSON(http.StatusOK, dto.OK(dto.AdvertisementToResponse(a)))
}

// @Summary     Delete advertisement
// @Description Delete an advertisement (requires editor role)
// @Tags        ads
// @Produce     json
// @Security    BearerAuth
// @Param       id   path string true "Venue ID"
// @Param       adID path string true "Advertisement ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/ads/{adID} [delete]
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

// @Summary     Publish advertisement
// @Description Publish a draft advertisement (requires editor role)
// @Tags        ads
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string true "Venue ID"
// @Param       adID path     string true "Advertisement ID"
// @Success     200  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/ads/{adID}/publish [post]
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
