package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type languageHandler struct {
	svc service.LanguageService
}

func newLanguageHandler(svc service.LanguageService) *languageHandler {
	return &languageHandler{svc: svc}
}

func (h *languageHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	langs, err := h.svc.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LanguageResponse, len(langs))
	for i, l := range langs {
		items[i] = dto.LanguageToResponse(l)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *languageHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	var req dto.LanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	l := &domain.Language{
		VenueID:   venueID,
		Code:      req.Code,
		Name:      req.Name,
		IsDefault: req.IsDefault,
	}
	if err := h.svc.Create(c.Request.Context(), l); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LanguageToResponse(l)))
}

func (h *languageHandler) Update(c *gin.Context) {
	langID, err := uuid.Parse(c.Param("langID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid language id"))
		return
	}
	var req dto.LanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	l, err := h.svc.Get(c.Request.Context(), langID)
	if err != nil {
		respondError(c, err)
		return
	}
	l.Code = req.Code
	l.Name = req.Name
	l.IsDefault = req.IsDefault
	if err := h.svc.Update(c.Request.Context(), l); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LanguageToResponse(l)))
}

func (h *languageHandler) Delete(c *gin.Context) {
	langID, err := uuid.Parse(c.Param("langID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid language id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), langID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
