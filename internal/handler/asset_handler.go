package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type assetHandler struct{ svc service.AssetService }

func newAssetHandler(svc service.AssetService) *assetHandler {
	return &assetHandler{svc: svc}
}

func (h *assetHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	assets, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.AssetResponse, len(assets))
	for i, a := range assets {
		items[i] = dto.AssetToResponse(a)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *assetHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid asset id"))
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.AssetToResponse(a)))
}

func (h *assetHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	userID := middleware.GetUserID(c)
	a, err := h.svc.Create(c.Request.Context(), service.CreateAssetInput{
		VenueID:     venueID,
		Name:        req.Name,
		Key:         req.Key,
		ContentType: req.ContentType,
		SizeBytes:   req.SizeBytes,
		URL:         req.URL,
		CreatedBy:   &userID,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.AssetToResponse(a)))
}

func (h *assetHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid asset id"))
		return
	}
	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a, err := h.svc.Update(c.Request.Context(), id, req.Name)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.AssetToResponse(a)))
}

func (h *assetHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid asset id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
