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

// @Summary     List tags
// @Description List all tags
// @Tags        tags
// @Produce     json
// @Security    BearerAuth
// @Param       page      query    int false "Page number"
// @Param       page_size query    int false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     401       {object} dto.Response
// @Router      /tags [get]
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

// @Summary     Get tag
// @Description Get a tag by ID
// @Tags        tags
// @Produce     json
// @Security    BearerAuth
// @Param       tagID path     string true "Tag ID"
// @Success     200   {object} dto.Response{data=dto.TagResponse}
// @Failure     400   {object} dto.Response
// @Failure     401   {object} dto.Response
// @Failure     404   {object} dto.Response
// @Router      /tags/{tagID} [get]
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

// @Summary     Create tag
// @Description Create a new tag
// @Tags        tags
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.TagRequest true "Tag details"
// @Success     201  {object} dto.Response{data=dto.TagResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /tags [post]
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

// @Summary     Update tag
// @Description Update a tag
// @Tags        tags
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       tagID path     string          true "Tag ID"
// @Param       body  body     dto.TagRequest  true "Tag details"
// @Success     200   {object} dto.Response{data=dto.TagResponse}
// @Failure     400   {object} dto.Response
// @Failure     401   {object} dto.Response
// @Failure     404   {object} dto.Response
// @Router      /tags/{tagID} [put]
func (h *tagHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid tag id"))
		return
	}
	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	t, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(t)
	if err := h.svc.Update(c.Request.Context(), t); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.TagToResponse(t)))
}

// @Summary     Delete tag
// @Description Delete a tag
// @Tags        tags
// @Produce     json
// @Security    BearerAuth
// @Param       tagID path string true "Tag ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /tags/{tagID} [delete]
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

// @Summary     Attach tag to entity
// @Description Attach a tag to an entity
// @Tags        tags
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.AttachTagRequest true "Attach tag details"
// @Success     201  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /tags/attach [post]
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

// @Summary     Detach tag from entity
// @Description Detach a tag from an entity
// @Tags        tags
// @Produce     json
// @Security    BearerAuth
// @Param       tagID       path     string true  "Tag ID"
// @Param       entity_type query    string true  "Entity type"
// @Param       entity_id   query    string true  "Entity ID"
// @Success     204
// @Failure     400         {object} dto.Response
// @Failure     401         {object} dto.Response
// @Router      /tags/{tagID}/detach [delete]
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

// @Summary     List entity tags
// @Description List all tags attached to a specific entity
// @Tags        tags
// @Produce     json
// @Security    BearerAuth
// @Param       entity_type query    string true  "Entity type"
// @Param       entity_id   query    string true  "Entity ID"
// @Success     200         {object} dto.Response{data=[]dto.TagResponse}
// @Failure     400         {object} dto.Response
// @Failure     401         {object} dto.Response
// @Router      /tags/entity [get]
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
