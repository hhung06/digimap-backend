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

// @Summary     List customers
// @Description List all customers (system admin only)
// @Tags        customers
// @Produce     json
// @Security    BearerAuth
// @Param       page      query    int false "Page number"
// @Param       page_size query    int false "Page size"
// @Success     200       {object} dto.Response{data=dto.PaginatedData}
// @Failure     401       {object} dto.Response
// @Failure     403       {object} dto.Response
// @Router      /customers [get]
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

// @Summary     Get customer
// @Description Get a customer by ID (system admin only)
// @Tags        customers
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Customer ID"
// @Success     200 {object} dto.Response{data=dto.CustomerResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /customers/{id} [get]
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

// @Summary     Create customer
// @Description Create a new customer (system admin only)
// @Tags        customers
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.CustomerRequest true "Customer details"
// @Success     201  {object} dto.Response{data=dto.CustomerResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /customers [post]
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

// @Summary     Update customer
// @Description Update an existing customer (system admin only)
// @Tags        customers
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string              true "Customer ID"
// @Param       body body     dto.CustomerRequest true "Customer details"
// @Success     200  {object} dto.Response{data=dto.CustomerResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Failure     404  {object} dto.Response
// @Router      /customers/{id} [put]
func (h *customerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid customer id"))
		return
	}
	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	cust, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(cust)
	if err := h.svc.Update(c.Request.Context(), cust); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.CustomerToResponse(cust)))
}

// @Summary     Delete customer
// @Description Delete a customer by ID (system admin only)
// @Tags        customers
// @Produce     json
// @Security    BearerAuth
// @Param       id  path string true "Customer ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Router      /customers/{id} [delete]
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
