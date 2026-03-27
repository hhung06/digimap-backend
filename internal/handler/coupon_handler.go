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

func (h *couponHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("couponID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid coupon id"))
		return
	}
	var req dto.CouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cp := &domain.Coupon{
		ID: id, ExternalID: req.ExternalID,
		CouponName: req.CouponName, CouponCode: req.CouponCode,
		Status: req.Status, IssuedAt: req.IssuedAt, ExpiredAt: req.ExpiredAt,
		Localization: req.Localization,
	}
	if err := h.svc.Update(c.Request.Context(), cp); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.CouponToResponse(cp)))
}

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
