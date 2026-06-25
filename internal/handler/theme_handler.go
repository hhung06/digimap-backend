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

// ListGlobal returns all global themes.
func (h *themeHandler) ListGlobal(c *gin.Context) {
	themes, err := h.svc.ListGlobal(c.Request.Context())
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

// CreateGlobal creates a new global theme.
func (h *themeHandler) CreateGlobal(c *gin.Context) {
	var req dto.ThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t, err := h.svc.CreateGlobal(c.Request.Context(), req.Name, req.Data)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ThemeToResponse(t)))
}

// List returns custom themes for a venue plus all global themes.
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

// Create creates a custom theme for a venue.
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
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "name is required"))
		return
	}
	t, err := h.svc.Create(c.Request.Context(), venueID, req.Name, req.Data)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.ThemeToResponse(t)))
}

// Update updates name and data of any theme (global or custom).
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
	t, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Data)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ThemeToResponse(t)))
}

// Delete soft-deletes a theme; fails if the theme is in use by any venue.
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

// SetTheme assigns (or clears) the active theme for a venue.
func (h *themeHandler) SetTheme(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.SetThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	if err := h.svc.SetVenueTheme(c.Request.Context(), venueID, req.ThemeID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
}
