package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type locationHandler struct {
	categorySvc service.LocationCategoryService
	amenitySvc  service.AmenityService
	locationSvc service.LocationService
}

func newLocationHandler(
	categorySvc service.LocationCategoryService,
	amenitySvc service.AmenityService,
	locationSvc service.LocationService,
) *locationHandler {
	return &locationHandler{categorySvc: categorySvc, amenitySvc: amenitySvc, locationSvc: locationSvc}
}

// ── Location categories ───────────────────────────────────────────────────────

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

// ── Amenities ─────────────────────────────────────────────────────────────────

func (h *locationHandler) ListAmenitiesByVenue(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	amenities, err := h.amenitySvc.ListByVenue(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.AmenityResponse, len(amenities))
	for i, a := range amenities {
		items[i] = dto.AmenityToResponse(a)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *locationHandler) LinkAmenity(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	amenityID, err := uuid.Parse(c.Param("amenityID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid amenity id"))
		return
	}
	if err := h.amenitySvc.LinkToVenue(c.Request.Context(), venueID, amenityID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"linked": true}))
}

func (h *locationHandler) UnlinkAmenity(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	amenityID, err := uuid.Parse(c.Param("amenityID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid amenity id"))
		return
	}
	if err := h.amenitySvc.UnlinkFromVenue(c.Request.Context(), venueID, amenityID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Locations ─────────────────────────────────────────────────────────────────

func (h *locationHandler) ListLocations(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	p := paginationFromQuery(c)
	locations, total, err := h.locationSvc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LocationResponse, len(locations))
	for i, l := range locations {
		items[i] = dto.LocationToResponse(l)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *locationHandler) GetLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	l, err := h.locationSvc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationToResponse(l)))
}

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

// ── Promotions ────────────────────────────────────────────────────────────────

func (h *locationHandler) ListPromotions(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	p := paginationFromQuery(c)
	promos, total, err := h.locationSvc.ListPromotions(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.PromotionResponse, len(promos))
	for i, promo := range promos {
		items[i] = dto.PromotionToResponse(promo)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *locationHandler) GetPromotion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid promotion id"))
		return
	}
	promo, err := h.locationSvc.GetPromotion(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PromotionToResponse(promo)))
}

func (h *locationHandler) CreatePromotion(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	var req dto.PromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	promo := &domain.Promotion{
		VenueID: venueID, LocationID: req.LocationID, ExternalID: req.ExternalID,
		PromoImage: req.PromoImage, Introduction: req.Introduction,
		GiftContent: req.GiftContent, DetailURL: req.DetailURL,
		BoothNumber: req.BoothNumber, ExpectedGiftCount: req.ExpectedGiftCount,
		DistributionStart: req.DistributionStart, DistributionEnd: req.DistributionEnd,
		DisplayType: req.DisplayType, Localization: req.Localization,
	}
	if err := h.locationSvc.CreatePromotion(c.Request.Context(), promo); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.PromotionToResponse(promo)))
}

func (h *locationHandler) UpdatePromotion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid promotion id"))
		return
	}
	var req dto.PromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	promo := &domain.Promotion{
		ID: id, LocationID: req.LocationID, ExternalID: req.ExternalID,
		PromoImage: req.PromoImage, Introduction: req.Introduction,
		GiftContent: req.GiftContent, DetailURL: req.DetailURL,
		BoothNumber: req.BoothNumber, ExpectedGiftCount: req.ExpectedGiftCount,
		DistributionStart: req.DistributionStart, DistributionEnd: req.DistributionEnd,
		DisplayType: req.DisplayType, Localization: req.Localization,
	}
	if err := h.locationSvc.UpdatePromotion(c.Request.Context(), promo); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PromotionToResponse(promo)))
}

func (h *locationHandler) DeletePromotion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid promotion id"))
		return
	}
	if err := h.locationSvc.DeletePromotion(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
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
		CommonLocationType: req.CommonLocationType, CommonSubType: req.CommonSubType,
		CommonLatitude: req.CommonLatitude, CommonLongitude: req.CommonLongitude,
		CommonAddress: req.CommonAddress,
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
		CommonLocationType: req.CommonLocationType, CommonSubType: req.CommonSubType,
		CommonLatitude: req.CommonLatitude, CommonLongitude: req.CommonLongitude,
		CommonAddress: req.CommonAddress, CommonLogo: req.CommonLogo,
		CommonContactEmail: req.CommonContactEmail, CommonContactPhone: req.CommonContactPhone,
		PlaceWorkHours: req.PlaceWorkHours, Custom: req.Custom, Localization: req.Localization,
		Source: req.Source, StartTime: req.StartTime, EndTime: req.EndTime,
		IsSearchable: req.IsSearchable,
	}
}
