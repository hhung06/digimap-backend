package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type tagHandler struct {
	svc service.TagService
}

func newTagHandler(svc service.TagService) *tagHandler {
	return &tagHandler{svc: svc}
}

func (h *tagHandler) List(c *gin.Context) {
	p := paginationFromQuery(c)
	tags, total, err := h.svc.List(c.Request.Context(), p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.TagResponse, len(tags))
	for i, t := range tags {
		items[i] = dto.TagToResponse(t)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *tagHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid tag id"))
		return
	}
	t, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.TagToResponse(t)))
}

func (h *tagHandler) Create(c *gin.Context) {
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t := &domain.Tag{Name: req.Name, Localization: req.Localization}
	if err := h.svc.Create(c.Request.Context(), t); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.TagToResponse(t)))
}

func (h *tagHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid tag id"))
		return
	}
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t := &domain.Tag{ID: id, Name: req.Name, Localization: req.Localization}
	if err := h.svc.Update(c.Request.Context(), t); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.TagToResponse(t)))
}

func (h *tagHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid tag id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *tagHandler) AttachTag(c *gin.Context) {
	var req dto.AttachTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	et := &domain.EntityTag{TagID: req.TagID, EntityType: req.EntityType, EntityID: req.EntityID}
	if err := h.svc.AttachTag(c.Request.Context(), et); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(nil))
}

func (h *tagHandler) DetachTag(c *gin.Context) {
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid tag id"))
		return
	}
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid entity id"))
		return
	}
	if err := h.svc.DetachTag(c.Request.Context(), tagID, entityType, entityID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *tagHandler) ListEntityTags(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid entity id"))
		return
	}
	tags, err := h.svc.ListEntityTags(c.Request.Context(), entityType, entityID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.TagResponse, len(tags))
	for i, t := range tags {
		items[i] = dto.TagToResponse(t)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}
