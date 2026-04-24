package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type connectionHandler struct {
	svc service.ConnectionService
}

func newConnectionHandler(svc service.ConnectionService) *connectionHandler {
	return &connectionHandler{svc: svc}
}

// @Summary     List connections
// @Description List connections for a venue
// @Tags        connections
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/connections [get]
func (h *connectionHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	conns, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ConnectionResponse, len(conns))
	for i, conn := range conns {
		items[i] = dto.ConnectionToResponse(conn)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

// @Summary     Get connection
// @Description Get a connection by ID
// @Tags        connections
// @Produce     json
// @Security    BearerAuth
// @Param       id           path     string true "Venue ID"
// @Param       connectionID path     string true "Connection ID"
// @Success     200          {object} dto.Response{data=dto.ConnectionResponse}
// @Failure     400          {object} dto.Response
// @Failure     401          {object} dto.Response
// @Failure     404          {object} dto.Response
// @Router      /venues/{id}/connections/{connectionID} [get]
func (h *connectionHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("connectionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid connection id"))
		return
	}
	conn, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ConnectionToResponse(conn)))
}

// @Summary     Create connection
// @Description Create a new connection (requires editor role)
// @Tags        connections
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                  true "Venue ID"
// @Param       body body     dto.ConnectionRequest   true "Connection details"
// @Success     201  {object} dto.Response{data=dto.ConnectionResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/connections [post]
func (h *connectionHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	conn := &domain.Connection{
		VenueID: &venueID, ExternalID: req.ExternalID, Name: req.Name,
		Type: req.Type, X: req.X, Y: req.Y,
		State: req.State, Status: req.Status,
		Accessible: req.Accessible, Active: req.Active,
	}
	if err := h.svc.Create(c.Request.Context(), conn); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ConnectionToResponse(conn)))
}

// @Summary     Update connection
// @Description Update a connection (requires editor role)
// @Tags        connections
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id           path     string                 true "Venue ID"
// @Param       connectionID path     string                 true "Connection ID"
// @Param       body         body     dto.ConnectionRequest  true "Connection details"
// @Success     200          {object} dto.Response{data=dto.ConnectionResponse}
// @Failure     400          {object} dto.Response
// @Failure     401          {object} dto.Response
// @Failure     403          {object} dto.Response
// @Router      /venues/{id}/connections/{connectionID} [put]
func (h *connectionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("connectionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid connection id"))
		return
	}
	var req dto.ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	conn := &domain.Connection{
		ID: id, ExternalID: req.ExternalID, Name: req.Name,
		Type: req.Type, X: req.X, Y: req.Y,
		State: req.State, Status: req.Status,
		Accessible: req.Accessible, Active: req.Active,
	}
	if err := h.svc.Update(c.Request.Context(), conn); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ConnectionToResponse(conn)))
}

// @Summary     Delete connection
// @Description Delete a connection (requires editor role)
// @Tags        connections
// @Produce     json
// @Security    BearerAuth
// @Param       id           path string true "Venue ID"
// @Param       connectionID path string true "Connection ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/connections/{connectionID} [delete]
func (h *connectionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("connectionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid connection id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary     List connection levels
// @Description List levels attached to a connection
// @Tags        connections
// @Produce     json
// @Security    BearerAuth
// @Param       id           path     string true "Venue ID"
// @Param       connectionID path     string true "Connection ID"
// @Success     200          {object} dto.Response{data=[]dto.ConnectionLevelResponse}
// @Failure     400          {object} dto.Response
// @Failure     401          {object} dto.Response
// @Router      /venues/{id}/connections/{connectionID}/levels [get]
func (h *connectionHandler) ListLevels(c *gin.Context) {
	id, err := uuid.Parse(c.Param("connectionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid connection id"))
		return
	}
	levels, err := h.svc.ListLevels(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ConnectionLevelResponse, len(levels))
	for i, cl := range levels {
		items[i] = dto.ConnectionLevelToResponse(cl)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Add level to connection
// @Description Attach a level to a connection (requires editor role)
// @Tags        connections
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id           path     string                       true "Venue ID"
// @Param       connectionID path     string                       true "Connection ID"
// @Param       body         body     dto.ConnectionLevelRequest   true "Level details"
// @Success     201          {object} dto.Response{data=dto.ConnectionLevelResponse}
// @Failure     400          {object} dto.Response
// @Failure     401          {object} dto.Response
// @Failure     403          {object} dto.Response
// @Router      /venues/{id}/connections/{connectionID}/levels [post]
func (h *connectionHandler) AddLevel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("connectionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid connection id"))
		return
	}
	var req dto.ConnectionLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cl := &domain.ConnectionLevel{
		ConnectionID: id, LevelID: req.LevelID, ElementID: req.ElementID, Active: req.Active,
	}
	if err := h.svc.AddLevel(c.Request.Context(), cl); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ConnectionLevelToResponse(cl)))
}

// @Summary     Remove level from connection
// @Description Detach a level from a connection (requires editor role)
// @Tags        connections
// @Produce     json
// @Security    BearerAuth
// @Param       id           path string true "Venue ID"
// @Param       connectionID path string true "Connection ID"
// @Param       clID         path string true "Connection Level ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/connections/{connectionID}/levels/{clID} [delete]
func (h *connectionHandler) RemoveLevel(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("clID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid connection level id"))
		return
	}
	if err := h.svc.RemoveLevel(c.Request.Context(), levelID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
