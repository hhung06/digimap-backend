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

// List godoc
// GET /api/v1/venues
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

// Get godoc
// GET /api/v1/venues/:id
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

// Create godoc
// POST /api/v1/venues
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

// Update godoc
// PUT /api/v1/venues/:id
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
	v := venueFromUpdateRequest(id, req)
	if err := h.svc.Update(c.Request.Context(), v); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VenueToResponse(v)))
}

// Delete godoc
// DELETE /api/v1/venues/:id
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

// Publish godoc
// POST /api/v1/venues/:id/publish
func (h *venueHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var body struct {
		Published bool `json:"published"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid request body"))
		return
	}
	if err := h.svc.Publish(c.Request.Context(), id, body.Published); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"published": body.Published}))
}

// Clone godoc
// POST /api/v1/venues/:id/clone
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

// GetKey godoc
// GET /api/v1/venues/:id/key
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

// RegenerateKey godoc
// PUT /api/v1/venues/:id/key
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
		CustomerID: req.CustomerID, Name: req.Name, Slug: req.Slug,
		ExternalID: req.ExternalID, Type: req.Type,
		Address: req.Address, City: req.City, State: req.State,
		Country: req.Country, Postal: req.Postal,
		Lat: req.Lat, Lng: req.Lng, Timezone: req.Timezone,
		Telephone: req.Telephone, WorkHours: req.WorkHours,
		Description: req.Description,
		Theme: req.Theme, Plugins: req.Plugins, Translations: req.Translations,
		Localization: req.Localization, CustomData: req.CustomData,
		AppConfigs: req.AppConfigs, AppDomains: req.AppDomains, SubDomains: req.SubDomains,
		SEOTitle: req.SEOTitle, SEODescription: req.SEODescription, SEOKeywords: req.SEOKeywords,
		HeadTag: req.HeadTag, BodyTag: req.BodyTag,
		StartAt: req.StartAt, EndAt: req.EndAt,
	}
}

func venueFromUpdateRequest(id uuid.UUID, req dto.UpdateVenueRequest) *domain.Venue {
	return &domain.Venue{
		ID: id, Name: req.Name, Slug: req.Slug,
		ExternalID: req.ExternalID, Type: req.Type,
		Address: req.Address, City: req.City, State: req.State,
		Country: req.Country, Postal: req.Postal,
		Lat: req.Lat, Lng: req.Lng, Timezone: req.Timezone,
		Telephone: req.Telephone, WorkHours: req.WorkHours,
		Description: req.Description, IsPublished: req.IsPublished,
		Theme: req.Theme, Plugins: req.Plugins, Translations: req.Translations,
		Localization: req.Localization, CustomData: req.CustomData,
		AppConfigs: req.AppConfigs, AppDomains: req.AppDomains, SubDomains: req.SubDomains,
		SEOTitle: req.SEOTitle, SEODescription: req.SEODescription, SEOKeywords: req.SEOKeywords,
		HeadTag: req.HeadTag, BodyTag: req.BodyTag,
		OriginalLogo: req.OriginalLogo, SmallLogo: req.SmallLogo,
		MediumLogo: req.MediumLogo, LargeLogo: req.LargeLogo,
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
