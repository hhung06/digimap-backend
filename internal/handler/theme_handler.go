package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type themeHandler struct{ svc service.ThemeService }

func newThemeHandler(svc service.ThemeService) *themeHandler {
	return &themeHandler{svc: svc}
}

func (h *themeHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	themes, err := h.svc.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ThemeResponse, len(themes))
	for i, t := range themes {
		items[i] = dto.ThemeToResponse(t)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *themeHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.ThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t, err := h.svc.Create(c.Request.Context(), venueID, req.Name, req.PrimaryColor, req.SecondaryColor)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ThemeToResponse(t)))
}

func (h *themeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("themeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid theme id"))
		return
	}
	var req dto.ThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t, err := h.svc.Update(c.Request.Context(), id, req.Name, req.PrimaryColor, req.SecondaryColor)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ThemeToResponse(t)))
}

func (h *themeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("themeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid theme id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
