package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type customerHandler struct {
	svc service.CustomerService
}

func newCustomerHandler(svc service.CustomerService) *customerHandler {
	return &customerHandler{svc: svc}
}

// List godoc
// GET /api/v1/customers
func (h *customerHandler) List(c *gin.Context) {
	p := paginationFromQuery(c)
	customers, total, err := h.svc.List(c.Request.Context(), p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.CustomerResponse, len(customers))
	for i, cust := range customers {
		items[i] = dto.CustomerToResponse(cust)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

// Get godoc
// GET /api/v1/customers/:id
func (h *customerHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid customer id"))
		return
	}
	cust, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.CustomerToResponse(cust)))
}

// Create godoc
// POST /api/v1/customers
func (h *customerHandler) Create(c *gin.Context) {
	var req dto.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cust := &domain.Customer{
		Name:        req.Name,
		Image:       req.Image,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		URL:         req.URL,
		Description: req.Description,
	}
	if err := h.svc.Create(c.Request.Context(), cust); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.CustomerToResponse(cust)))
}

// Update godoc
// PUT /api/v1/customers/:id
func (h *customerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid customer id"))
		return
	}
	var req dto.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cust := &domain.Customer{
		ID:          id,
		Name:        req.Name,
		Image:       req.Image,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		URL:         req.URL,
		Description: req.Description,
	}
	if err := h.svc.Update(c.Request.Context(), cust); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.CustomerToResponse(cust)))
}

// Delete godoc
// DELETE /api/v1/customers/:id
func (h *customerHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid customer id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
