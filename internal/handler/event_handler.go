package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type eventHandler struct {
	svc service.EventService
}

func newEventHandler(svc service.EventService) *eventHandler {
	return &eventHandler{svc: svc}
}

// ── Event tags ────────────────────────────────────────────────────────────────

func (h *eventHandler) ListTags(c *gin.Context) {
	tags, err := h.svc.ListTags(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.EventTagResponse, len(tags))
	for i, t := range tags {
		items[i] = dto.EventTagToResponse(t)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *eventHandler) CreateTag(c *gin.Context) {
	var req dto.EventTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t := &domain.EventTag{Name: req.Name, Localization: req.Localization}
	if err := h.svc.CreateTag(c.Request.Context(), t); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.EventTagToResponse(t)))
}

func (h *eventHandler) UpdateTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid tag id"))
		return
	}
	var req dto.EventTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t := &domain.EventTag{ID: id, Name: req.Name, Localization: req.Localization}
	if err := h.svc.UpdateTag(c.Request.Context(), t); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.EventTagToResponse(t)))
}

func (h *eventHandler) DeleteTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid tag id"))
		return
	}
	if err := h.svc.DeleteTag(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Event types ───────────────────────────────────────────────────────────────

func (h *eventHandler) ListEventTypes(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	types, err := h.svc.ListEventTypes(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.EventTypeResponse, len(types))
	for i, t := range types {
		items[i] = dto.EventTypeToResponse(t)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *eventHandler) CreateEventType(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.EventTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t := &domain.EventType{VenueID: &venueID, Name: req.Name, Localization: req.Localization}
	if err := h.svc.CreateEventType(c.Request.Context(), t); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.EventTypeToResponse(t)))
}

func (h *eventHandler) UpdateEventType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("typeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid event type id"))
		return
	}
	var req dto.EventTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t := &domain.EventType{ID: id, Name: req.Name, Localization: req.Localization}
	if err := h.svc.UpdateEventType(c.Request.Context(), t); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.EventTypeToResponse(t)))
}

func (h *eventHandler) DeleteEventType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("typeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid event type id"))
		return
	}
	if err := h.svc.DeleteEventType(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Events ────────────────────────────────────────────────────────────────────

func (h *eventHandler) ListEvents(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	events, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.EventResponse, len(events))
	for i, e := range events {
		items[i] = dto.EventToResponse(e)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *eventHandler) GetEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid event id"))
		return
	}
	e, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.EventToResponse(e)))
}

func (h *eventHandler) CreateEvent(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	e := &domain.Event{
		VenueID: venueID, TypeID: req.TypeID, Title: req.Title,
		Description: req.Description, BannerImage: req.BannerImage,
		IconImage: req.IconImage, StartTime: req.StartTime, EndTime: req.EndTime,
		ShowStartTime: req.ShowStartTime, ShowEndTime: req.ShowEndTime,
		ContentDetail: req.ContentDetail, ContentURL: req.ContentURL,
		Localization: req.Localization,
	}
	if err := h.svc.Create(c.Request.Context(), e, req.TagIDs, req.LocationIDs); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.EventToResponse(e)))
}

func (h *eventHandler) UpdateEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid event id"))
		return
	}
	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	e := &domain.Event{
		ID: id, TypeID: req.TypeID, Title: req.Title,
		Description: req.Description, BannerImage: req.BannerImage,
		IconImage: req.IconImage, StartTime: req.StartTime, EndTime: req.EndTime,
		ShowStartTime: req.ShowStartTime, ShowEndTime: req.ShowEndTime,
		ContentDetail: req.ContentDetail, ContentURL: req.ContentURL,
		Localization: req.Localization,
	}
	if err := h.svc.Update(c.Request.Context(), e, req.TagIDs, req.LocationIDs); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.EventToResponse(e)))
}

func (h *eventHandler) DeleteEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid event id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Event images ──────────────────────────────────────────────────────────────

func (h *eventHandler) CreateEventImage(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid event id"))
		return
	}
	var req dto.EventImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	img := &domain.EventImage{EventID: eventID, Image: req.Image}
	if err := h.svc.CreateImage(c.Request.Context(), img); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.EventImageResponse{
		ID: img.ID, EventID: img.EventID, Image: img.Image, CreatedAt: img.CreatedAt,
	}))
}

func (h *eventHandler) DeleteEventImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("imageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid image id"))
		return
	}
	if err := h.svc.DeleteImage(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
