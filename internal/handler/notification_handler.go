package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type notificationHandler struct {
	svc      service.NotificationService
	enrichers *enricher.Registry
}

func newNotificationHandler(svc service.NotificationService, enrichers *enricher.Registry) *notificationHandler {
	return &notificationHandler{svc: svc, enrichers: enrichers}
}

// @Summary     List notifications
// @Description List notifications for a venue
// @Tags        notifications
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/notifications [get]
func (h *notificationHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	ns, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceNotification)
	items := make([]any, len(ns))
	for i, n := range ns {
		items[i] = enricher.MergeInto(dto.NotificationToResponse(n), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

// @Summary     Get notification
// @Description Get a notification by ID
// @Tags        notifications
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string true "Venue ID"
// @Param       notifID path     string true "Notification ID"
// @Success     200     {object} dto.Response{data=dto.NotificationResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     404     {object} dto.Response
// @Router      /venues/{id}/notifications/{notifID} [get]
func (h *notificationHandler) Get(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	id, err := uuid.Parse(c.Param("notifID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid notification id"))
		return
	}
	n, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceNotification)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.NotificationToResponse(n), extras)))
}

// @Summary     Create notification
// @Description Create a new notification (requires editor role)
// @Tags        notifications
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                          true "Venue ID"
// @Param       body body     dto.CreateNotificationRequest  true "Notification details"
// @Success     201  {object} dto.Response{data=dto.NotificationResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/notifications [post]
func (h *notificationHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	userID := middleware.GetUserID(c)
	kind := req.Kind
	if kind == 0 {
		kind = domain.NotifKindNormal
	}
	sendType := req.SendType
	if sendType == 0 {
		sendType = domain.NotifTypeDraft
	}
	targetApp := req.TargetApp
	if targetApp == "" {
		targetApp = "all"
	}
	n := &domain.Notification{
		VenueID: &venueID, SurveyID: req.SurveyID,
		Title: &req.Title, Content: &req.Content, Topic: req.Topic,
		Kind: kind, SendType: sendType,
		Data: req.Data, LinkURL: req.LinkURL, ScheduledAt: req.ScheduledAt,
		TargetApp: targetApp, SegmentFilters: req.SegmentFilters,
		DeviceTokens: req.DeviceTokens,
		CreatedBy:    &userID,
	}
	if err := h.svc.Create(c.Request.Context(), n); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.NotificationToResponse(n)))
}

// @Summary     Update notification
// @Description Update a notification (requires editor role)
// @Tags        notifications
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string                          true "Venue ID"
// @Param       notifID path     string                          true "Notification ID"
// @Param       body    body     dto.UpdateNotificationRequest  true "Notification details"
// @Success     200     {object} dto.Response{data=dto.NotificationResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     403     {object} dto.Response
// @Router      /venues/{id}/notifications/{notifID} [put]
func (h *notificationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("notifID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid notification id"))
		return
	}
	var req dto.UpdateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	n, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(n)
	if err := h.svc.Update(c.Request.Context(), n); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.NotificationToResponse(n)))
}

// @Summary     Delete notification
// @Description Delete a notification (requires editor role)
// @Tags        notifications
// @Produce     json
// @Security    BearerAuth
// @Param       id      path string true "Venue ID"
// @Param       notifID path string true "Notification ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/notifications/{notifID} [delete]
func (h *notificationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("notifID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid notification id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary     Send notification
// @Description Send a notification to its recipients (requires editor role)
// @Tags        notifications
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string true "Venue ID"
// @Param       notifID path     string true "Notification ID"
// @Success     200     {object} dto.Response
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     403     {object} dto.Response
// @Router      /venues/{id}/notifications/{notifID}/send [post]
func (h *notificationHandler) Send(c *gin.Context) {
	id, err := uuid.Parse(c.Param("notifID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid notification id"))
		return
	}
	if err := h.svc.Send(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
}
