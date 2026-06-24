package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type venueHandler struct {
	svc service.VenueService
}

func newVenueHandler(svc service.VenueService) *venueHandler {
	return &venueHandler{svc: svc}
}

// @Summary     List venues
// @Description List venues; system admins see all, others must supply customer_id
// @Tags        venues
// @Produce     json
// @Security    BearerAuth
// @Param       customer_id query    string false "Customer ID (required for non-admins)"
// @Param       page        query    int    false "Page number"
// @Param       page_size   query    int    false "Page size"
// @Success     200         {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400         {object} dto.Response
// @Failure     401         {object} dto.Response
// @Router      /venues [get]
func (h *venueHandler) List(c *gin.Context) {
	p := paginationFromQuery(c)

	// System admins see all venues; others are scoped to their customer.
	if middleware.GetIsSystemAdmin(c) {
		venues, total, err := h.svc.ListAll(c.Request.Context(), p)
		if err != nil {
			respondError(c, err)
			return
		}
		items := venuesToResponse(venues)
		c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
		return
	}

	customerIDStr := c.Query("customer_id")
	if customerIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "customer_id query param required"))
		return
	}
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid customer_id"))
		return
	}

	venues, total, err := h.svc.List(c.Request.Context(), customerID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Paginated(venuesToResponse(venues), total, p.Page, p.PageSize))
}

// @Summary     Get venue
// @Description Get a venue by ID
// @Tags        venues
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=dto.VenueResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /venues/{id} [get]
func (h *venueHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	v, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VenueToResponse(v)))
}

// @Summary     Create venue
// @Description Create a new venue
// @Tags        venues
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.CreateVenueRequest true "Venue details"
// @Success     201  {object} dto.Response{data=dto.VenueResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /venues [post]
func (h *venueHandler) Create(c *gin.Context) {
	var req dto.CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	v := venueFromCreateRequest(req)
	if err := h.svc.Create(c.Request.Context(), v); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.VenueToResponse(v)))
}

// @Summary     Update venue
// @Description Update an existing venue (requires editor role)
// @Tags        venues
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                true "Venue ID"
// @Param       body body     dto.UpdateVenueRequest true "Venue details"
// @Success     200  {object} dto.Response{data=dto.VenueResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Failure     404  {object} dto.Response
// @Router      /venues/{id} [put]
func (h *venueHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	v, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(v)
	if err := h.svc.Update(c.Request.Context(), v); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VenueToResponse(v)))
}

// @Summary     Delete venue
// @Description Delete a venue (requires owner role)
// @Tags        venues
// @Produce     json
// @Security    BearerAuth
// @Param       id  path string true "Venue ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /venues/{id} [delete]
func (h *venueHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary     Clone venue
// @Description Clone a venue (requires owner role)
// @Tags        venues
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     201 {object} dto.Response{data=dto.VenueResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /venues/{id}/clone [post]
func (h *venueHandler) Clone(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	cloned, err := h.svc.Clone(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.VenueToResponse(cloned)))
}

// @Summary     Get venue API keys
// @Description Get public and private API keys for a venue (requires owner role)
// @Tags        venues
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=dto.VenueKeyResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /venues/{id}/key [get]
func (h *venueHandler) GetKey(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	v, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VenueKeyResponse{
		PublicKey:  v.PublicKey,
		PrivateKey: v.PrivateKey,
	}))
}

// @Summary     Regenerate venue API keys
// @Description Regenerate API keys for a venue (requires owner role)
// @Tags        venues
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=dto.VenueKeyResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /venues/{id}/key [put]
func (h *venueHandler) RegenerateKey(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	v, err := h.svc.RegenerateKeys(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VenueKeyResponse{
		PublicKey:  v.PublicKey,
		PrivateKey: v.PrivateKey,
	}))
}

// ── helpers ───────────────────────────────────────────────────────────────────

func venueFromCreateRequest(req dto.CreateVenueRequest) *domain.Venue {
	return &domain.Venue{
		CustomerID: req.CustomerID, Name: req.Name,
		ExternalID: req.ExternalID, Type: req.Type,
		Address: req.Address, City: req.City, State: req.State,
		Country: req.Country, Postal: req.Postal,
		Lat: req.Lat, Lng: req.Lng, Timezone: req.Timezone,
		Telephone:    req.Telephone,
		Description:  req.Description,
		Localization: req.Localization,
		AppConfigs:   req.AppConfigs, AppDomains: req.AppDomains, SubDomains: req.SubDomains,
		SEOTitle: req.SEOTitle, SEODescription: req.SEODescription, SEOKeywords: req.SEOKeywords,
		HeadTag: req.HeadTag, BodyTag: req.BodyTag,
		StartAt: req.StartAt, EndAt: req.EndAt,
	}
}

func venuesToResponse(venues []*domain.Venue) []dto.VenueResponse {
	items := make([]dto.VenueResponse, len(venues))
	for i, v := range venues {
		items[i] = dto.VenueToResponse(v)
	}
	return items
}
