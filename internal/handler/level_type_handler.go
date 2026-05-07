package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type levelTypeHandler struct{ svc service.LevelTypeService }

func newLevelTypeHandler(svc service.LevelTypeService) *levelTypeHandler {
	return &levelTypeHandler{svc: svc}
}

func (h *levelTypeHandler) List(c *gin.Context) {
	types, err := h.svc.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LevelTypeResponse, len(types))
	for i, lt := range types {
		items[i] = dto.LevelTypeToResponse(lt)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *levelTypeHandler) Create(c *gin.Context) {
	var req dto.LevelTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	lt, err := h.svc.Create(c.Request.Context(), req.Name, req.Icon)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LevelTypeToResponse(lt)))
}

func (h *levelTypeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("typeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level type id"))
		return
	}
	var req dto.LevelTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	lt, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Icon)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LevelTypeToResponse(lt)))
}

func (h *levelTypeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("typeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level type id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
