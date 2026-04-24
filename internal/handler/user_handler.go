package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type userHandler struct {
	svc service.UserService
}

func newUserHandler(svc service.UserService) *userHandler {
	return &userHandler{svc: svc}
}

// ── Profile ───────────────────────────────────────────────────────────────────

// @Summary     Get profile
// @Description Get the authenticated user's profile
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} dto.Response{data=dto.UserResponse}
// @Failure     401 {object} dto.Response
// @Router      /profile [get]
func (h *userHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	u, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.UserToResponse(u)))
}

// @Summary     Update profile
// @Description Update the authenticated user's profile
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.UpdateProfileRequest true "Profile details"
// @Success     200  {object} dto.Response{data=dto.UserResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /profile [put]
func (h *userHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	u := &domain.User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		AvatarURL: req.AvatarURL,
		IsActive:  true,
	}
	if err := h.svc.UpdateProfile(c.Request.Context(), u); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.UserToResponse(u)))
}

// ── Venue users ───────────────────────────────────────────────────────────────

// @Summary     List venue users
// @Description List all users in a venue
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=[]dto.UserResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Router      /venues/{id}/users [get]
func (h *userHandler) ListVenueUsers(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	users, err := h.svc.ListVenueUsers(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.UserResponse, len(users))
	for i, u := range users {
		items[i] = dto.UserToResponse(u)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Change user role
// @Description Change a user's role in a venue (requires owner role)
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id     path     string                true "Venue ID"
// @Param       userID path     string                true "User ID"
// @Param       body   body     dto.ChangeRoleRequest true "Role details"
// @Success     200    {object} dto.Response
// @Failure     400    {object} dto.Response
// @Failure     401    {object} dto.Response
// @Failure     403    {object} dto.Response
// @Router      /venues/{id}/users/{userID}/role [put]
func (h *userHandler) ChangeRole(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid user id"))
		return
	}
	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	role := domain.RoleFromString(req.Role)
	if err := h.svc.ChangeRole(c.Request.Context(), venueID, userID, role); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
}

// @Summary     Remove user from venue
// @Description Remove a user from a venue (requires owner role)
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       id     path string true "Venue ID"
// @Param       userID path string true "User ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/users/{userID} [delete]
func (h *userHandler) RemoveFromVenue(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid user id"))
		return
	}
	if err := h.svc.RemoveFromVenue(c.Request.Context(), venueID, userID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Invitations ───────────────────────────────────────────────────────────────

// @Summary     List invitations
// @Description List pending invitations for a venue (requires owner role)
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     string true "Venue ID"
// @Success     200 {object} dto.Response{data=[]dto.InvitationResponse}
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/invitations [get]
func (h *userHandler) ListInvitations(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	invs, err := h.svc.ListInvitations(c.Request.Context(), venueID)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.InvitationResponse, len(invs))
	for i, inv := range invs {
		items[i] = dto.InvitationToResponse(inv)
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// @Summary     Invite user
// @Description Invite a user to a venue (requires owner role)
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                true "Venue ID"
// @Param       body body     dto.InviteUserRequest true "Invitation details"
// @Success     201  {object} dto.Response{data=dto.InvitationResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/users/invite [post]
func (h *userHandler) InviteUser(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	invitedBy := middleware.GetUserID(c)
	role := domain.RoleFromString(req.Role)
	inv, err := h.svc.InviteUser(c.Request.Context(), venueID, invitedBy, req.Email, role)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.InvitationToResponse(inv)))
}

// @Summary     Accept invitation
// @Description Accept a venue invitation using a token
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.AcceptInvitationRequest true "Invitation token"
// @Success     200  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /invitations/accept [post]
func (h *userHandler) AcceptInvitation(c *gin.Context) {
	var req dto.AcceptInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	userID := middleware.GetUserID(c)
	if err := h.svc.AcceptInvitation(c.Request.Context(), req.Token, userID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
}

// @Summary     Cancel invitation
// @Description Cancel a pending invitation (requires owner role)
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       id           path string true "Venue ID"
// @Param       invitationID path string true "Invitation ID"
// @Success     200          {object} dto.Response
// @Failure     400          {object} dto.Response
// @Failure     401          {object} dto.Response
// @Failure     403          {object} dto.Response
// @Router      /venues/{id}/invitations/{invitationID}/cancel [post]
func (h *userHandler) CancelInvitation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("invitationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid invitation id"))
		return
	}
	if err := h.svc.CancelInvitation(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
}
