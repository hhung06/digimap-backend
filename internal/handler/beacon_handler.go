package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type beaconHandler struct {
	svc service.BeaconService
}

func newBeaconHandler(svc service.BeaconService) *beaconHandler {
	return &beaconHandler{svc: svc}
}

// @Summary     List beacons
// @Description List beacons for a venue
// @Tags        beacons
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/beacons [get]
func (h *beaconHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	beacons, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.BeaconResponse, len(beacons))
	for i, b := range beacons {
		items[i] = dto.BeaconToResponse(b)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

// @Summary     Get beacon
// @Description Get a beacon by ID
// @Tags        beacons
// @Produce     json
// @Security    BearerAuth
// @Param       id       path     string true "Venue ID"
// @Param       beaconID path     string true "Beacon ID"
// @Success     200      {object} dto.Response{data=dto.BeaconResponse}
// @Failure     400      {object} dto.Response
// @Failure     401      {object} dto.Response
// @Failure     404      {object} dto.Response
// @Router      /venues/{id}/beacons/{beaconID} [get]
func (h *beaconHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("beaconID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid beacon id"))
		return
	}
	b, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.BeaconToResponse(b)))
}

// @Summary     Create beacon
// @Description Create a new beacon (requires editor role)
// @Tags        beacons
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string             true "Venue ID"
// @Param       body body     dto.BeaconRequest  true "Beacon details"
// @Success     201  {object} dto.Response{data=dto.BeaconResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/beacons [post]
func (h *beaconHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.BeaconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	b := &domain.Beacon{
		VenueID: &venueID, LevelID: req.LevelID, ElementID: req.ElementID,
		Name: req.Name, HwID: req.HwID, VendorKey: req.VendorKey, LotKey: req.LotKey,
		UUIDVal: req.UUIDVal, MAC: req.MAC,
		Radius: req.Radius, Battery: req.Battery,
		PositionX: req.PositionX, PositionY: req.PositionY, IsEnable: req.IsEnable,
		Major: req.Major, Minor: req.Minor, Voltage: req.Voltage, TxPower: req.TxPower,
	}
	if err := h.svc.Create(c.Request.Context(), b); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.BeaconToResponse(b)))
}

// @Summary     Update beacon
// @Description Update a beacon (requires editor role)
// @Tags        beacons
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id       path     string             true "Venue ID"
// @Param       beaconID path     string             true "Beacon ID"
// @Param       body     body     dto.BeaconRequest  true "Beacon details"
// @Success     200      {object} dto.Response{data=dto.BeaconResponse}
// @Failure     400      {object} dto.Response
// @Failure     401      {object} dto.Response
// @Failure     403      {object} dto.Response
// @Router      /venues/{id}/beacons/{beaconID} [put]
func (h *beaconHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("beaconID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid beacon id"))
		return
	}
	var req dto.UpdateBeaconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	b, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(b)
	if err := h.svc.Update(c.Request.Context(), b); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.BeaconToResponse(b)))
}

// @Summary     Delete beacon
// @Description Delete a beacon (requires editor role)
// @Tags        beacons
// @Produce     json
// @Security    BearerAuth
// @Param       id       path string true "Venue ID"
// @Param       beaconID path string true "Beacon ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/beacons/{beaconID} [delete]
func (h *beaconHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("beaconID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid beacon id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
