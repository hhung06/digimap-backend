package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type couponHandler struct {
	svc service.CouponService
}

func newCouponHandler(svc service.CouponService) *couponHandler {
	return &couponHandler{svc: svc}
}

// @Summary     List coupons
// @Description List coupons for a venue
// @Tags        coupons
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/coupons [get]
func (h *couponHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	coupons, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.CouponResponse, len(coupons))
	for i, cp := range coupons {
		items[i] = dto.CouponToResponse(cp)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}

// @Summary     Get coupon
// @Description Get a coupon by ID
// @Tags        coupons
// @Produce     json
// @Security    BearerAuth
// @Param       id       path     string true "Venue ID"
// @Param       couponID path     string true "Coupon ID"
// @Success     200      {object} dto.Response{data=dto.CouponResponse}
// @Failure     400      {object} dto.Response
// @Failure     401      {object} dto.Response
// @Failure     404      {object} dto.Response
// @Router      /venues/{id}/coupons/{couponID} [get]
func (h *couponHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("couponID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid coupon id"))
		return
	}
	cp, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.CouponToResponse(cp)))
}

// @Summary     Create coupon
// @Description Create a new coupon (requires editor role)
// @Tags        coupons
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string              true "Venue ID"
// @Param       body body     dto.CouponRequest   true "Coupon details"
// @Success     201  {object} dto.Response{data=dto.CouponResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/coupons [post]
func (h *couponHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cp := &domain.Coupon{
		VenueID: &venueID, ExternalID: req.ExternalID,
		CouponName: req.CouponName, CouponCode: req.CouponCode,
		Status: req.Status, IssuedAt: req.IssuedAt, ExpiredAt: req.ExpiredAt,
		Localization: req.Localization,
	}
	if err := h.svc.Create(c.Request.Context(), cp); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.CouponToResponse(cp)))
}

// @Summary     Update coupon
// @Description Update a coupon (requires editor role)
// @Tags        coupons
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id       path     string             true "Venue ID"
// @Param       couponID path     string             true "Coupon ID"
// @Param       body     body     dto.CouponRequest  true "Coupon details"
// @Success     200      {object} dto.Response{data=dto.CouponResponse}
// @Failure     400      {object} dto.Response
// @Failure     401      {object} dto.Response
// @Failure     403      {object} dto.Response
// @Router      /venues/{id}/coupons/{couponID} [put]
func (h *couponHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("couponID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid coupon id"))
		return
	}
	var req dto.UpdateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cp, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(cp)
	if err := h.svc.Update(c.Request.Context(), cp); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.CouponToResponse(cp)))
}

// @Summary     Delete coupon
// @Description Delete a coupon (requires editor role)
// @Tags        coupons
// @Produce     json
// @Security    BearerAuth
// @Param       id       path string true "Venue ID"
// @Param       couponID path string true "Coupon ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/coupons/{couponID} [delete]
func (h *couponHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("couponID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid coupon id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
