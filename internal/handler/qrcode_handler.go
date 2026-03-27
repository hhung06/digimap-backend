package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type qrcodeHandler struct {
	svc service.QRCodeService
}

func newQRCodeHandler(svc service.QRCodeService) *qrcodeHandler {
	return &qrcodeHandler{svc: svc}
}

func (h *qrcodeHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	codes, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.QRCodeResponse, len(codes))
	for i, q := range codes {
		items[i] = dto.QRCodeToResponse(q)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *qrcodeHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("qrID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid qr code id"))
		return
	}
	q, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.QRCodeToResponse(q)))
}

func (h *qrcodeHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.QRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	q := &domain.QRCode{
		VenueID: &venueID, LevelID: req.LevelID, LocationID: req.LocationID,
		Lat: req.Lat, Lng: req.Lng, Angle: req.Angle, Link: req.Link,
	}
	if err := h.svc.Create(c.Request.Context(), q); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.QRCodeToResponse(q)))
}

func (h *qrcodeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("qrID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid qr code id"))
		return
	}
	var req dto.QRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	q := &domain.QRCode{
		ID: id, LevelID: req.LevelID, LocationID: req.LocationID,
		Lat: req.Lat, Lng: req.Lng, Angle: req.Angle, Link: req.Link,
	}
	if err := h.svc.Update(c.Request.Context(), q); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.QRCodeToResponse(q)))
}

func (h *qrcodeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("qrID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid qr code id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
