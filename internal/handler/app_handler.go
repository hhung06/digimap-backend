package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type appHandler struct {
	locations          service.LocationService
	events             service.EventService
	products           service.ProductService
	productPlazas      service.ProductPlazaService
	articles           service.ArticleService
	notifications      service.NotificationService
	coupons            service.CouponService
	ads                service.AdvertisementService
	surveys            service.SurveyService
	locationCategories service.LocationCategoryService
	analytics          service.AnalyticsService
}

func newAppHandler(
	locations service.LocationService,
	events service.EventService,
	products service.ProductService,
	productPlazas service.ProductPlazaService,
	articles service.ArticleService,
	notifications service.NotificationService,
	coupons service.CouponService,
	ads service.AdvertisementService,
	surveys service.SurveyService,
	locationCategories service.LocationCategoryService,
	analytics service.AnalyticsService,
) *appHandler {
	return &appHandler{
		locations:          locations,
		events:             events,
		products:           products,
		productPlazas:      productPlazas,
		articles:           articles,
		notifications:      notifications,
		coupons:            coupons,
		ads:                ads,
		surveys:            surveys,
		locationCategories: locationCategories,
		analytics:          analytics,
	}
}

func (h *appHandler) ListLocations(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	locations, total, err := h.locations.List(c.Request.Context(), venueID, nil, p)
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

func (h *appHandler) GetLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid location id"))
		return
	}
	l, err := h.locations.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationToResponse(l)))
}

func (h *appHandler) ListEvents(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	events, total, err := h.events.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.EventResponse, len(events))
	for i, e := range events {
		items[i] = dto.EventToResponse(e)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *appHandler) ListProducts(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	products, total, err := h.products.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ProductResponse, len(products))
	for i, prod := range products {
		items[i] = dto.ProductToResponse(prod)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *appHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid product id"))
		return
	}
	prod, err := h.products.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ProductToResponse(prod)))
}

func (h *appHandler) ListProductCategories(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	cats, err := h.products.ListCategories(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ProductCategoryResponse, len(cats))
	for i, cat := range cats {
		items[i] = dto.ProductCategoryToResponse(cat)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *appHandler) ListProductPlazas(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	plazas, err := h.productPlazas.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ProductPlazaResponse, len(plazas))
	for i, p := range plazas {
		items[i] = dto.ProductPlazaToResponse(p)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *appHandler) GetProductPlaza(c *gin.Context) {
	id, err := uuid.Parse(c.Param("plazaID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid plaza id"))
		return
	}
	plaza, err := h.productPlazas.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ProductPlazaToResponse(plaza)))
}

func (h *appHandler) ListArticles(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	articles, total, err := h.articles.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.ArticleResponse, len(articles))
	for i, a := range articles {
		items[i] = dto.ArticleToResponse(a)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *appHandler) GetArticle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("articleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid article id"))
		return
	}
	a, err := h.articles.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.ArticleToResponse(a)))
}

func (h *appHandler) ListFeaturedZones(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	cats, err := h.locationCategories.List(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.LocationCategoryResponse, 0, len(cats))
	for _, cat := range cats {
		if cat.Source == "external" {
			items = append(items, dto.LocationCategoryToResponse(cat))
		}
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *appHandler) GetFeaturedZone(c *gin.Context) {
	id, err := uuid.Parse(c.Param("zoneID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid zone id"))
		return
	}
	cat, err := h.locationCategories.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.LocationCategoryToResponse(cat)))
}

func (h *appHandler) ListNotifications(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	notifs, total, err := h.notifications.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.NotificationResponse, len(notifs))
	for i, n := range notifs {
		items[i] = dto.NotificationToResponse(n)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *appHandler) GetNotification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("notifID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid notification id"))
		return
	}
	n, err := h.notifications.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.NotificationToResponse(n)))
}

func (h *appHandler) ListCoupons(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	coupons, total, err := h.coupons.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		items[i] = dto.CouponToResponse(coupon)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *appHandler) GetCoupon(c *gin.Context) {
	id, err := uuid.Parse(c.Param("couponID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid coupon id"))
		return
	}
	coupon, err := h.coupons.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.CouponToResponse(coupon)))
}

func (h *appHandler) RedeemCoupon(c *gin.Context) {
	id, err := uuid.Parse(c.Param("couponID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid coupon id"))
		return
	}
	var body struct {
		AppUserID string `json:"app_user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	appUserID, err := uuid.Parse(body.AppUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid app_user_id"))
		return
	}
	if err := h.coupons.RedeemCoupon(c.Request.Context(), id, appUserID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
}

func (h *appHandler) ListAds(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	ads, _, err := h.ads.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.AdvertisementResponse, len(ads))
	for i, a := range ads {
		items[i] = dto.AdvertisementToResponse(a)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

func (h *appHandler) GetSurvey(c *gin.Context) {
	id, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid survey id"))
		return
	}
	s, err := h.surveys.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SurveyToResponse(s)))
}

func (h *appHandler) TopSearch(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	p := paginationFromQuery(c)
	queries, total, err := h.analytics.ListSearchQueries(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.SearchQueryResponse, len(queries))
	for i, q := range queries {
		items[i] = dto.SearchQueryToResponse(q)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

func (h *appHandler) SearchOptions(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	q := c.Query("q")
	const limit = 20
	locs, err := h.locations.SearchByName(c.Request.Context(), venueID, q, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	prods, err := h.products.SearchByName(c.Request.Context(), venueID, q, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	locItems := make([]dto.LocationResponse, len(locs))
	for i, l := range locs {
		locItems[i] = dto.LocationToResponse(l)
	}
	prodItems := make([]dto.ProductResponse, len(prods))
	for i, p := range prods {
		prodItems[i] = dto.ProductToResponse(p)
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{
		"locations": locItems,
		"products":  prodItems,
	}))
}

func (h *appHandler) Promotions(c *gin.Context) {
	venueID := middleware.GetVenueID(c)
	surveys, err := h.surveys.ListPromo(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	p := paginationFromQuery(c)
	coupons, total, err := h.coupons.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	surveyItems := make([]dto.SurveyResponse, len(surveys))
	for i, s := range surveys {
		surveyItems[i] = dto.SurveyToResponse(s)
	}
	couponItems := make([]dto.CouponResponse, len(coupons))
	for i, cp := range coupons {
		couponItems[i] = dto.CouponToResponse(cp)
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{
		"surveys": surveyItems,
		"coupons": gin.H{
			"items": couponItems,
			"total": total,
		},
	}))
}
