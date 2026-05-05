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

// @Summary     List event tags
// @Description List all global event tags
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} dto.Response{data=[]dto.EventTagResponse}
// @Failure     401 {object} dto.Response
// @Router      /event-tags [get]
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

// @Summary     Create event tag
// @Description Create a new global event tag
// @Tags        events
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.EventTagRequest true "Tag details"
// @Success     201  {object} dto.Response{data=dto.EventTagResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /event-tags [post]
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

// @Summary     Update event tag
// @Description Update a global event tag
// @Tags        events
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       tagID path     string              true "Tag ID"
// @Param       body  body     dto.EventTagRequest true "Tag details"
// @Success     200   {object} dto.Response{data=dto.EventTagResponse}
// @Failure     400   {object} dto.Response
// @Failure     401   {object} dto.Response
// @Failure     404   {object} dto.Response
// @Router      /event-tags/{tagID} [put]
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

// @Summary     Delete event tag
// @Description Delete a global event tag
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Param       tagID path string true "Tag ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /event-tags/{tagID} [delete]
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

// @Summary     List event types
// @Description List all event types for a venue
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=[]dto.EventTypeResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /venues/{id}/event-types [get]
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

// @Summary     Create event type
// @Description Create a new event type for a venue (requires editor role)
// @Tags        events
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string               true "Venue ID"
// @Param       body body     dto.EventTypeRequest true "Event type details"
// @Success     201  {object} dto.Response{data=dto.EventTypeResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/event-types [post]
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

// @Summary     Update event type
// @Description Update an event type (requires editor role)
// @Tags        events
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id     path     string               true "Venue ID"
// @Param       typeID path     string               true "Event Type ID"
// @Param       body   body     dto.EventTypeRequest true "Event type details"
// @Success     200    {object} dto.Response{data=dto.EventTypeResponse}
// @Failure     400    {object} dto.Response
// @Failure     401    {object} dto.Response
// @Failure     403    {object} dto.Response
// @Router      /venues/{id}/event-types/{typeID} [put]
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

// @Summary     Delete event type
// @Description Delete an event type (requires editor role)
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Param       id     path string true "Venue ID"
// @Param       typeID path string true "Event Type ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/event-types/{typeID} [delete]
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

// @Summary     List events
// @Description List events for a venue
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/events [get]
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

// @Summary     Get event
// @Description Get an event by ID
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string true "Venue ID"
// @Param       eventID path     string true "Event ID"
// @Success     200     {object} dto.Response{data=dto.EventResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     404     {object} dto.Response
// @Router      /venues/{id}/events/{eventID} [get]
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

// @Summary     Create event
// @Description Create a new event (requires editor role)
// @Tags        events
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                  true "Venue ID"
// @Param       body body     dto.CreateEventRequest  true "Event details"
// @Success     201  {object} dto.Response{data=dto.EventResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/events [post]
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

// @Summary     Update event
// @Description Update an event (requires editor role)
// @Tags        events
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string                 true "Venue ID"
// @Param       eventID path     string                 true "Event ID"
// @Param       body    body     dto.UpdateEventRequest true "Event details"
// @Success     200     {object} dto.Response{data=dto.EventResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     403     {object} dto.Response
// @Failure     404     {object} dto.Response
// @Router      /venues/{id}/events/{eventID} [put]
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

// @Summary     Delete event
// @Description Delete an event (requires editor role)
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Param       id      path string true "Venue ID"
// @Param       eventID path string true "Event ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/events/{eventID} [delete]
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

// @Summary     Add event image
// @Description Add an image to an event (requires editor role)
// @Tags        events
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string                 true "Venue ID"
// @Param       eventID path     string                 true "Event ID"
// @Param       body    body     dto.EventImageRequest  true "Image details"
// @Success     201     {object} dto.Response{data=dto.EventImageResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     403     {object} dto.Response
// @Router      /venues/{id}/events/{eventID}/images [post]
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

// @Summary     Delete event image
// @Description Delete an image from an event (requires editor role)
// @Tags        events
// @Produce     json
// @Security    BearerAuth
// @Param       id      path string true "Venue ID"
// @Param       eventID path string true "Event ID"
// @Param       imageID path string true "Image ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/events/{eventID}/images/{imageID} [delete]
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
