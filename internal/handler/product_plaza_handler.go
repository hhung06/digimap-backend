package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type productPlazaHandler struct{ svc service.ProductPlazaService }

func newProductPlazaHandler(svc service.ProductPlazaService) *productPlazaHandler {
	return &productPlazaHandler{svc: svc}
}

func (h *productPlazaHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	plazas, err := h.svc.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ProductPlazaResponse, len(plazas))
	for i, p := range plazas {
		items[i] = dto.ProductPlazaToResponse(p)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *productPlazaHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.ProductPlazaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	p, err := h.svc.Create(c.Request.Context(), service.CreateProductPlazaInput{
		VenueID:      venueID,
		Name:         req.Name,
		Description:  req.Description,
		Localization: req.Localization,
		LocationID:   req.LocationID,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ProductPlazaToResponse(p)))
}

func (h *productPlazaHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("plazaID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid plaza id"))
		return
	}
	var req dto.ProductPlazaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, service.CreateProductPlazaInput{
		Name:         req.Name,
		Description:  req.Description,
		Localization: req.Localization,
		LocationID:   req.LocationID,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ProductPlazaToResponse(p)))
}

func (h *productPlazaHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("plazaID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid plaza id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
