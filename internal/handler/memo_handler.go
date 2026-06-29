package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type memoHandler struct {
	svc service.MemoService
}

func newMemoHandler(svc service.MemoService) *memoHandler {
	return &memoHandler{svc: svc}
}

func (h *memoHandler) List(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	memos, err := h.svc.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LocationResponse, len(memos))
	for i, m := range memos {
		items[i] = dto.LocationToResponse(m)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *memoHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("memoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid memo id"))
		return
	}
	memo, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationToResponse(memo)))
}

func (h *memoHandler) Create(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, err.Error()))
		return
	}
	loc := &domain.Location{
		VenueID:               venueID,
		LevelID:               req.LevelID,
		MainCategoryID:        req.MainCategory,
		ExternalID:            req.ExternalID,
		CommonHidden:          req.CommonHidden,
		CommonName:            req.CommonName,
		CommonShortName:       req.CommonShortName,
		CommonDescription:     req.CommonDescription,
		CommonColor:           req.CommonColor,
		CommonLocationType:    domain.LocationTypeMemo,
		CommonLocationSubType: req.CommonLocationSubType,
		CommonLatitude:        req.CommonLatitude,
		CommonLongitude:       req.CommonLongitude,
		CommonAddress:         req.CommonAddress,
		CommonContactEmail:    req.CommonContactEmail,
		CommonContactPhone:    req.CommonContactPhone,
		PlaceWorkHours:        req.PlaceWorkHours,
		Custom:                req.Custom,
		Localization:          req.Localization,
		Source:                req.Source,
		StartTime:             req.StartTime,
		EndTime:               req.EndTime,
		IsSearchable:          req.IsSearchable,
	}
	if err := h.svc.Create(c.Request.Context(), loc); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LocationToResponse(loc)))
}

func (h *memoHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("memoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid memo id"))
		return
	}
	memo, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, err.Error()))
		return
	}
	req.ApplyTo(memo)
	memo.CommonLocationType = domain.LocationTypeMemo // never allow changing the type
	if err := h.svc.Update(c.Request.Context(), memo); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationToResponse(memo)))
}

func (h *memoHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("memoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid memo id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OKMessage("memo deleted"))
}
