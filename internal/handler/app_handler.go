package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type appHandler struct {
	locations service.LocationService
	events    service.EventService
}

func newAppHandler(locations service.LocationService, events service.EventService) *appHandler {
	return &appHandler{locations: locations, events: events}
}

func (h *appHandler) ListLocations(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	locations, total, err := h.locations.List(c.Request.Context(), venueID, nil, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LocationResponse, len(locations))
	for i, l := range locations {
		items[i] = dto.LocationToResponse(l)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *appHandler) ListEvents(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	events, total, err := h.events.List(c.Request.Context(), venueID, p)
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
