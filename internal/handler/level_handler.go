package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type levelHandler struct {
	svc service.LevelService
}

func newLevelHandler(svc service.LevelService) *levelHandler {
	return &levelHandler{svc: svc}
}

// ── Map groups ────────────────────────────────────────────────────────────────

// @Summary     List map groups
// @Description List all map groups for a venue
// @Tags        levels
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=[]dto.MapGroupResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /venues/{id}/map-groups [get]
func (h *levelHandler) ListMapGroups(c *gin.Context) {
	venueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	groups, err := h.svc.ListMapGroups(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.MapGroupResponse, len(groups))
	for i, mg := range groups {
		items[i] = dto.MapGroupToResponse(mg)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Create map group
// @Description Create a new map group for a venue (requires editor role)
// @Tags        levels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string              true "Venue ID"
// @Param       body body     dto.MapGroupRequest true "Map group details"
// @Success     201  {object} dto.Response{data=dto.MapGroupResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/map-groups [post]
func (h *levelHandler) CreateMapGroup(c *gin.Context) {
	venueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.MapGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	mg := &domain.MapGroup{
		VenueID: venueID, Type: req.Type, Name: req.Name,
		ShortName: req.ShortName, SortIndex: req.SortIndex,
	}
	if err := h.svc.CreateMapGroup(c.Request.Context(), mg); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.MapGroupToResponse(mg)))
}

// @Summary     Update map group
// @Description Update a map group (requires editor role)
// @Tags        levels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string              true "Venue ID"
// @Param       mgID path     string              true "Map Group ID"
// @Param       body body     dto.MapGroupRequest true "Map group details"
// @Success     200  {object} dto.Response{data=dto.MapGroupResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Failure     404  {object} dto.Response
// @Router      /venues/{id}/map-groups/{mgID} [put]
func (h *levelHandler) UpdateMapGroup(c *gin.Context) {
	mgID, err := uuid.Parse(c.Param("mgID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid map group id"))
		return
	}
	var req dto.UpdateMapGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	mg, err := h.svc.GetMapGroup(c.Request.Context(), mgID)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(mg)
	if err := h.svc.UpdateMapGroup(c.Request.Context(), mg); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.MapGroupToResponse(mg)))
}

// @Summary     Delete map group
// @Description Delete a map group (requires editor role)
// @Tags        levels
// @Produce     json
// @Security    BearerAuth
// @Param       id   path string true "Venue ID"
// @Param       mgID path string true "Map Group ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/map-groups/{mgID} [delete]
func (h *levelHandler) DeleteMapGroup(c *gin.Context) {
	mgID, err := uuid.Parse(c.Param("mgID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid map group id"))
		return
	}
	if err := h.svc.DeleteMapGroup(c.Request.Context(), mgID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Levels ────────────────────────────────────────────────────────────────────

// @Summary     List levels
// @Description List all levels for a venue
// @Tags        levels
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=[]dto.LevelResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /venues/{id}/levels [get]
func (h *levelHandler) ListLevels(c *gin.Context) {
	venueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	levels, err := h.svc.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LevelResponse, len(levels))
	for i, l := range levels {
		items[i] = dto.LevelToResponse(l)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Get level
// @Description Get a level by ID
// @Tags        levels
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string true "Venue ID"
// @Param       levelID path     string true "Level ID"
// @Success     200     {object} dto.Response{data=dto.LevelResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     404     {object} dto.Response
// @Router      /venues/{id}/levels/{levelID} [get]
func (h *levelHandler) GetLevel(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("levelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level id"))
		return
	}
	l, err := h.svc.Get(c.Request.Context(), levelID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LevelToResponse(l)))
}

// @Summary     Create level
// @Description Create a new level for a venue (requires editor role)
// @Tags        levels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                true "Venue ID"
// @Param       body body     dto.CreateLevelRequest true "Level details"
// @Success     201  {object} dto.Response{data=dto.LevelResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/levels [post]
func (h *levelHandler) CreateLevel(c *gin.Context) {
	venueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	l := &domain.Level{
		VenueID: venueID, MapGroupID: req.MapGroupID,
		Name: req.Name, ShortName: req.ShortName, ExternalID: req.ExternalID,
		Type: req.Type, Latitude: req.Latitude, Longitude: req.Longitude, Bearing: req.Bearing,
		Width: req.Width, Height: req.Height, Scale: req.Scale,
		LevelWidth: req.LevelWidth, LevelHeight: req.LevelHeight,
		FileIDs: req.FileIDs, Elevation: req.Elevation, IsPublished: req.IsPublished,
	}
	if err := h.svc.Create(c.Request.Context(), l); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LevelToResponse(l)))
}

// @Summary     Update level
// @Description Update a level (requires editor role)
// @Tags        levels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string                true "Venue ID"
// @Param       levelID path     string                true "Level ID"
// @Param       body    body     dto.UpdateLevelRequest true "Level details"
// @Success     200     {object} dto.Response{data=dto.LevelResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     403     {object} dto.Response
// @Failure     404     {object} dto.Response
// @Router      /venues/{id}/levels/{levelID} [put]
func (h *levelHandler) UpdateLevel(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("levelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level id"))
		return
	}
	var req dto.UpdateLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	l, err := h.svc.Get(c.Request.Context(), levelID)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(l)
	if err := h.svc.Update(c.Request.Context(), l); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LevelToResponse(l)))
}

// @Summary     Delete level
// @Description Delete a level (requires editor role)
// @Tags        levels
// @Produce     json
// @Security    BearerAuth
// @Param       id      path string true "Venue ID"
// @Param       levelID path string true "Level ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/levels/{levelID} [delete]
func (h *levelHandler) DeleteLevel(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("levelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), levelID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Perspectives ──────────────────────────────────────────────────────────────

// @Summary     Upsert perspective
// @Description Create or update the camera perspective for a level (requires editor role)
// @Tags        levels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string                  true "Venue ID"
// @Param       levelID path     string                  true "Level ID"
// @Param       body    body     dto.PerspectiveRequest  true "Perspective details"
// @Success     200     {object} dto.Response{data=dto.PerspectiveResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     403     {object} dto.Response
// @Router      /venues/{id}/levels/{levelID}/perspective [put]
func (h *levelHandler) UpsertPerspective(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("levelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level id"))
		return
	}
	var req dto.PerspectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	p := &domain.Perspective{
		Name: req.Name, CameraZoom: req.CameraZoom, CameraType: req.CameraType,
		CameraMaxZoom: req.CameraMaxZoom, CameraMinZoom: req.CameraMinZoom,
		CameraTargetCenterLng: req.CameraTargetCenterLng,
		CameraTargetCenterLat: req.CameraTargetCenterLat,
		CameraTargetZoom:      req.CameraTargetZoom,
		CameraTargetBearing:   req.CameraTargetBearing,
		CameraTargetPitch:     req.CameraTargetPitch,
	}
	if err := h.svc.UpsertPerspective(c.Request.Context(), levelID, p); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PerspectiveResponse{
		ID: p.ID, Name: p.Name, CameraZoom: p.CameraZoom, CameraType: p.CameraType,
		CameraMaxZoom: p.CameraMaxZoom, CameraMinZoom: p.CameraMinZoom,
		CameraTargetCenterLng: p.CameraTargetCenterLng,
		CameraTargetCenterLat: p.CameraTargetCenterLat,
		CameraTargetZoom:      p.CameraTargetZoom,
		CameraTargetBearing:   p.CameraTargetBearing,
		CameraTargetPitch:     p.CameraTargetPitch,
	}))
}

// ── Geo references ────────────────────────────────────────────────────────────

// @Summary     List geo references
// @Description List geo references for a level
// @Tags        levels
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string true "Venue ID"
// @Param       levelID path     string true "Level ID"
// @Success     200     {object} dto.Response{data=[]dto.GeoReferenceResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Router      /venues/{id}/levels/{levelID}/geo-references [get]
func (h *levelHandler) ListGeoReferences(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("levelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level id"))
		return
	}
	refs, err := h.svc.ListGeoReferences(c.Request.Context(), levelID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.GeoReferenceResponse, len(refs))
	for i, g := range refs {
		items[i] = dto.GeoRefToResponse(g)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Create geo reference
// @Description Create a geo reference for a level (requires editor role)
// @Tags        levels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string                  true "Venue ID"
// @Param       levelID path     string                  true "Level ID"
// @Param       body    body     dto.GeoReferenceRequest true "Geo reference details"
// @Success     201     {object} dto.Response{data=dto.GeoReferenceResponse}
// @Failure     400     {object} dto.Response
// @Failure     401     {object} dto.Response
// @Failure     403     {object} dto.Response
// @Router      /venues/{id}/levels/{levelID}/geo-references [post]
func (h *levelHandler) CreateGeoReference(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("levelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level id"))
		return
	}
	var req dto.GeoReferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	g := &domain.GeoReference{
		LevelID: levelID, ControlX: req.ControlX, ControlY: req.ControlY,
		TargetX: req.TargetX, TargetY: req.TargetY,
	}
	if err := h.svc.CreateGeoReference(c.Request.Context(), g); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.GeoRefToResponse(g)))
}

// @Summary     Delete geo reference
// @Description Delete a geo reference (requires editor role)
// @Tags        levels
// @Produce     json
// @Security    BearerAuth
// @Param       id      path string true "Venue ID"
// @Param       levelID path string true "Level ID"
// @Param       refID   path string true "Geo Reference ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/levels/{levelID}/geo-references/{refID} [delete]
func (h *levelHandler) DeleteGeoReference(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("levelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid level id"))
		return
	}
	refID, err := uuid.Parse(c.Param("refID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid geo reference id"))
		return
	}
	if err := h.svc.DeleteGeoReference(c.Request.Context(), levelID, refID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
