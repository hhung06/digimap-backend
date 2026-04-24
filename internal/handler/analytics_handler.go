package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type analyticsHandler struct {
	svc service.AnalyticsService
}

func newAnalyticsHandler(svc service.AnalyticsService) *analyticsHandler {
	return &analyticsHandler{svc: svc}
}

// @Summary     Track event
// @Description Track an analytics event (no auth required)
// @Tags        analytics
// @Accept      json
// @Produce     json
// @Param       id   path     string                 true "Venue ID"
// @Param       body body     dto.TrackEventRequest  true "Event data"
// @Success     201  {object} dto.Response{data=dto.EventLogResponse}
// @Failure     400  {object} dto.Response
// @Router      /public/venues/{id}/events [post]
func (h *analyticsHandler) TrackEvent(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.TrackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	e := &domain.EventLog{
		VenueID:   &venueID,
		Name:      req.Name,
		Params:    req.Params,
		DeviceID:  req.DeviceID,
		UserID:    req.UserID,
		UserAgent: c.Request.UserAgent(),
		IPAddress: clientIP(c),
	}
	if err := h.svc.LogEvent(c.Request.Context(), e); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.EventLogToResponse(e)))
}

// @Summary     Track search
// @Description Track a search query (no auth required)
// @Tags        analytics
// @Accept      json
// @Produce     json
// @Param       id   path     string                  true "Venue ID"
// @Param       body body     dto.TrackSearchRequest  true "Search data"
// @Success     201  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Router      /public/venues/{id}/searches [post]
func (h *analyticsHandler) TrackSearch(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.TrackSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	if err := h.svc.TrackSearch(c.Request.Context(), venueID, req.Term); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(nil))
}

// @Summary     List event logs
// @Description List analytics event logs for a venue
// @Tags        analytics
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/analytics/events [get]
func (h *analyticsHandler) ListEventLogs(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	logs, total, err := h.svc.ListEventLogs(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.EventLogResponse, len(logs))
	for i, e := range logs {
		items[i] = dto.EventLogToResponse(e)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

// @Summary     List search queries
// @Description List search query logs for a venue
// @Tags        analytics
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/analytics/searches [get]
func (h *analyticsHandler) ListSearchQueries(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	queries, total, err := h.svc.ListSearchQueries(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.SearchQueryResponse, len(queries))
	for i, q := range queries {
		items[i] = dto.SearchQueryToResponse(q)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

// clientIP extracts the real client IP from the request, respecting X-Forwarded-For.
func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// Take the first (original client) IP
		for i, ch := range xff {
			if ch == ',' {
				return xff[:i]
			}
		}
		return xff
	}
	return c.ClientIP()
}
