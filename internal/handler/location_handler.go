package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/service"
)

type locationHandler struct {
	categorySvc service.LocationCategoryService
	locationSvc service.LocationService
	enrichers   *enricher.Registry
}

func newLocationHandler(
	categorySvc service.LocationCategoryService,
	locationSvc service.LocationService,
	enrichers *enricher.Registry,
) *locationHandler {
	return &locationHandler{categorySvc: categorySvc, locationSvc: locationSvc, enrichers: enrichers}
}

// ── Location categories ───────────────────────────────────────────────────────

// @Summary     List location categories
// @Description List all location categories for a venue
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=[]dto.LocationCategoryResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /venues/{id}/categories [get]
func (h *locationHandler) ListCategories(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	cats, err := h.categorySvc.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LocationCategoryResponse, len(cats))
	for i, cat := range cats {
		items[i] = dto.LocationCategoryToResponse(cat)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Get location category
// @Description Get a location category by ID
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id    path     string true "Venue ID"
// @Param       catID path     string true "Category ID"
// @Success     200   {object} dto.Response{data=dto.LocationCategoryResponse}
// @Failure     400   {object} dto.Response
// @Failure     401   {object} dto.Response
// @Failure     404   {object} dto.Response
// @Router      /venues/{id}/categories/{catID} [get]
func (h *locationHandler) GetCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("catID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid category id"))
		return
	}
	cat, err := h.categorySvc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationCategoryToResponse(cat)))
}

// @Summary     Create location category
// @Description Create a new location category (requires editor role)
// @Tags        locations
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                      true "Venue ID"
// @Param       body body     dto.LocationCategoryRequest true "Category details"
// @Success     201  {object} dto.Response{data=dto.LocationCategoryResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/categories [post]
func (h *locationHandler) CreateCategory(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	var req dto.LocationCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cat := &domain.LocationCategory{
		VenueID: venueID, ExternalID: req.ExternalID, Name: req.Name,
		ShortName: req.ShortName, Color: req.Color, Icon: req.Icon,
		IconDefault: req.IconDefault, SortIndex: req.SortIndex,
		Visible: req.Visible, Description: req.Description,
		Type: req.Type, Image: req.Image, Localization: req.Localization,
		Source: req.Source,
	}
	if err := h.categorySvc.Create(c.Request.Context(), cat); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LocationCategoryToResponse(cat)))
}

// @Summary     Update location category
// @Description Update a location category (requires editor role)
// @Tags        locations
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path     string                      true "Venue ID"
// @Param       catID path     string                      true "Category ID"
// @Param       body  body     dto.LocationCategoryRequest true "Category details"
// @Success     200   {object} dto.Response{data=dto.LocationCategoryResponse}
// @Failure     400   {object} dto.Response
// @Failure     401   {object} dto.Response
// @Failure     403   {object} dto.Response
// @Failure     404   {object} dto.Response
// @Router      /venues/{id}/categories/{catID} [put]
func (h *locationHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("catID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid category id"))
		return
	}
	var req dto.LocationCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cat := &domain.LocationCategory{
		ID: id, ExternalID: req.ExternalID, Name: req.Name,
		ShortName: req.ShortName, Color: req.Color, Icon: req.Icon,
		IconDefault: req.IconDefault, SortIndex: req.SortIndex,
		Visible: req.Visible, Description: req.Description,
		Type: req.Type, Image: req.Image, Localization: req.Localization,
		Source: req.Source,
	}
	if err := h.categorySvc.Update(c.Request.Context(), cat); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationCategoryToResponse(cat)))
}

// @Summary     Delete location category
// @Description Delete a location category (requires editor role)
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id    path string true "Venue ID"
// @Param       catID path string true "Category ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/categories/{catID} [delete]
func (h *locationHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("catID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid category id"))
		return
	}
	if err := h.categorySvc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Locations ─────────────────────────────────────────────────────────────────

// @Summary     List locations
// @Description List locations for a venue
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/locations [get]
func (h *locationHandler) ListLocations(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	p := paginationFromQuery(c)

	// Parse optional ?type= filter
	var typeFilter *int
	if t := c.Query("type"); t != "" {
		if v, err := strconv.Atoi(t); err == nil {
			typeFilter = &v
		}
	}

	ctx := c.Request.Context()
	locations, total, err := h.locationSvc.List(ctx, venueID, typeFilter, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceLocation)
	items := make([]any, len(locations))
	for i, loc := range locations {
		items[i] = enricher.MergeInto(dto.LocationToResponse(loc), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

// @Summary     Get location
// @Description Get a location by ID
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string true "Venue ID"
// @Param       locationID path     string true "Location ID"
// @Success     200        {object} dto.Response{data=dto.LocationResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     404        {object} dto.Response
// @Router      /venues/{id}/locations/{locationID} [get]
func (h *locationHandler) GetLocation(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	ctx := c.Request.Context()
	l, err := h.locationSvc.Get(ctx, id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceLocation)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.LocationToResponse(l), extras)))
}

// @Summary     Create location
// @Description Create a new location (requires editor role)
// @Tags        locations
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                    true "Venue ID"
// @Param       body body     dto.CreateLocationRequest true "Location details"
// @Success     201  {object} dto.Response{data=dto.LocationResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/locations [post]
func (h *locationHandler) CreateLocation(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	l := locationFromCreateRequest(venueID, req)
	if err := h.locationSvc.Create(c.Request.Context(), l); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LocationToResponse(l)))
}

// @Summary     Update location
// @Description Update a location (requires editor role)
// @Tags        locations
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string                    true "Venue ID"
// @Param       locationID path     string                    true "Location ID"
// @Param       body       body     dto.UpdateLocationRequest true "Location details"
// @Success     200        {object} dto.Response{data=dto.LocationResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Failure     404        {object} dto.Response
// @Router      /venues/{id}/locations/{locationID} [put]
func (h *locationHandler) UpdateLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	l := locationFromUpdateRequest(id, req)
	if err := h.locationSvc.Update(c.Request.Context(), l, req.CategoryIDs); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationToResponse(l)))
}

// @Summary     Delete location
// @Description Delete a location (requires editor role)
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id         path string true "Venue ID"
// @Param       locationID path string true "Location ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/locations/{locationID} [delete]
func (h *locationHandler) DeleteLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	if err := h.locationSvc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary     Duplicate location
// @Description Duplicate a location (requires editor role)
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string true "Venue ID"
// @Param       locationID path     string true "Location ID"
// @Success     201        {object} dto.Response{data=dto.LocationResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/locations/{locationID}/duplicate [post]
func (h *locationHandler) DuplicateLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	clone, err := h.locationSvc.Duplicate(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LocationToResponse(clone)))
}

// @Summary     Set top location
// @Description Pin or unpin a location as a top result (requires editor role)
// @Tags        locations
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string                    true "Venue ID"
// @Param       locationID path     string                    true "Location ID"
// @Param       body       body     dto.SetTopLocationRequest true "Set top details"
// @Success     200        {object} dto.Response
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/locations/{locationID}/set-top [put]
func (h *locationHandler) SetTop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	var req dto.SetTopLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid request body"))
		return
	}
	if err := h.locationSvc.SetTop(c.Request.Context(), id, req.IsTop, req.SortIndex); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"is_top": req.IsTop}))
}

// ── Location images ───────────────────────────────────────────────────────────

// @Summary     Add location image
// @Description Add an image to a location (requires editor role)
// @Tags        locations
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string                    true "Venue ID"
// @Param       locationID path     string                    true "Location ID"
// @Param       body       body     dto.LocationImageRequest  true "Image details"
// @Success     201        {object} dto.Response{data=dto.LocationImageResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/locations/{locationID}/images [post]
func (h *locationHandler) CreateImage(c *gin.Context) {
	locationID, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	var req dto.LocationImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	img := &domain.LocationImage{
		LocationID: locationID, Original: req.Original,
		Small: req.Small, Medium: req.Medium, Large: req.Large,
	}
	if err := h.locationSvc.CreateImage(c.Request.Context(), img); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.LocationImageResponse{
		ID: img.ID, LocationID: img.LocationID,
		Original: img.Original, Small: img.Small,
		Medium: img.Medium, Large: img.Large, CreatedAt: img.CreatedAt,
	}))
}

// @Summary     Delete location image
// @Description Delete an image from a location (requires editor role)
// @Tags        locations
// @Produce     json
// @Security    BearerAuth
// @Param       id         path string true "Venue ID"
// @Param       locationID path string true "Location ID"
// @Param       imageID    path string true "Image ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/locations/{locationID}/images/{imageID} [delete]
func (h *locationHandler) DeleteImage(c *gin.Context) {
	locationID, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	imageID, err := uuid.Parse(c.Param("imageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid image id"))
		return
	}
	if err := h.locationSvc.DeleteImage(c.Request.Context(), locationID, imageID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GeoSearch returns locations within radiusKm of the given lat/lng.
func (h *locationHandler) GeoSearch(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")

	if latStr == "" || lngStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "lat, lng, and radius are required"))
		return
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)
	radius, err3 := strconv.ParseFloat(radiusStr, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "lat, lng, and radius must be valid numbers"))
		return
	}

	var venueID *uuid.UUID
	if v := c.Query("venue_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue_id"))
			return
		}
		venueID = &id
	}

	locs, err := h.locationSvc.GeoSearch(c.Request.Context(), lat, lng, radius, venueID)
	if err != nil {
		respondError(c, err)
		return
	}

	items := make([]dto.LocationResponse, len(locs))
	for i, loc := range locs {
		items[i] = dto.LocationToResponse(loc)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// ── helpers ───────────────────────────────────────────────────────────────────

func parseVenueID(c *gin.Context) (uuid.UUID, error) {
	return uuid.Parse(c.Param("id"))
}

func locationFromCreateRequest(venueID uuid.UUID, req dto.CreateLocationRequest) *domain.Location {
	return &domain.Location{
		VenueID: venueID, LevelID: req.LevelID, MainCategoryID: req.MainCategoryID,
		ExternalID: req.ExternalID, CommonHidden: req.CommonHidden,
		CommonName: req.CommonName, CommonShortName: req.CommonShortName,
		CommonDescription: req.CommonDescription, CommonColor: req.CommonColor,
		CommonLocationType: req.CommonLocationType, CommonLocationSubType: req.CommonLocationSubType,
		CommonLatitude: req.CommonLatitude, CommonLongitude: req.CommonLongitude,
		CommonAddress:      req.CommonAddress,
		CommonContactEmail: req.CommonContactEmail, CommonContactPhone: req.CommonContactPhone,
		PlaceWorkHours: req.PlaceWorkHours, Custom: req.Custom, Localization: req.Localization,
		Source: req.Source, StartTime: req.StartTime, EndTime: req.EndTime,
		IsSearchable: req.IsSearchable,
	}
}

func locationFromUpdateRequest(id uuid.UUID, req dto.UpdateLocationRequest) *domain.Location {
	return &domain.Location{
		ID: id, LevelID: req.LevelID, MainCategoryID: req.MainCategoryID,
		ExternalID: req.ExternalID, CommonHidden: req.CommonHidden,
		CommonName: req.CommonName, CommonShortName: req.CommonShortName,
		CommonDescription: req.CommonDescription, CommonColor: req.CommonColor,
		CommonLocationType: req.CommonLocationType, CommonLocationSubType: req.CommonLocationSubType,
		CommonLatitude: req.CommonLatitude, CommonLongitude: req.CommonLongitude,
		CommonAddress: req.CommonAddress, CommonLogo: req.CommonLogo,
		CommonContactEmail: req.CommonContactEmail, CommonContactPhone: req.CommonContactPhone,
		PlaceWorkHours: req.PlaceWorkHours, Custom: req.Custom, Localization: req.Localization,
		Source: req.Source, StartTime: req.StartTime, EndTime: req.EndTime,
		IsSearchable: req.IsSearchable,
	}
}
