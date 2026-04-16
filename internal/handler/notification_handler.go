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
		Title: req.Title, Content: req.Content, Topic: req.Topic,
		Kind: kind, SendType: sendType,
		Data: req.Data, LinkURL: req.LinkURL, ScheduledAt: req.ScheduledAt,
		TargetApp: targetApp, SegmentFilters: req.SegmentFilters,
		DeviceTokens: req.DeviceTokens,
		CreatedBy: &userID,
	}
	if err := h.svc.Create(c.Request.Context(), n); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.NotificationToResponse(n)))
}

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
	n := &domain.Notification{
		ID: id, SurveyID: req.SurveyID,
		Title: req.Title, Content: req.Content, Topic: req.Topic,
		Kind: req.Kind, SendType: req.SendType,
		Data: req.Data, LinkURL: req.LinkURL, ScheduledAt: req.ScheduledAt,
		TargetApp: req.TargetApp, SegmentFilters: req.SegmentFilters,
		DeviceTokens: req.DeviceTokens,
	}
	if err := h.svc.Update(c.Request.Context(), n); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.NotificationToResponse(n)))
}

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
