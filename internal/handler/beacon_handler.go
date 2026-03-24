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

func (h *beaconHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("beaconID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid beacon id"))
		return
	}
	var req dto.BeaconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	b := &domain.Beacon{
		ID: id, LevelID: req.LevelID, ElementID: req.ElementID,
		Name: req.Name, HwID: req.HwID, VendorKey: req.VendorKey, LotKey: req.LotKey,
		UUIDVal: req.UUIDVal, MAC: req.MAC,
		Radius: req.Radius, Battery: req.Battery,
		PositionX: req.PositionX, PositionY: req.PositionY, IsEnable: req.IsEnable,
		Major: req.Major, Minor: req.Minor, Voltage: req.Voltage, TxPower: req.TxPower,
	}
	if err := h.svc.Update(c.Request.Context(), b); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.BeaconToResponse(b)))
}

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
