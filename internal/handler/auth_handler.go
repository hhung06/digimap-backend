package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type authHandler struct {
	auth service.AuthService
	cfg  authConfig
}

type authConfig struct {
	accessExpirySeconds  int
	refreshExpirySeconds int
}

func newAuthHandler(authSvc service.AuthService, cfg authConfig) *authHandler {
	return &authHandler{auth: authSvc, cfg: cfg}
}

// @Summary     Login
// @Description Authenticate with email and password, returns JWT tokens
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     dto.LoginRequest true "Login credentials"
// @Success     200  {object} dto.Response{data=dto.LoginResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	user, pair, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.OK(dto.LoginResponse{
		AccessToken:      pair.AccessToken,
		AccessExpiresIn:  h.cfg.accessExpirySeconds,
		RefreshToken:     pair.RefreshToken,
		RefreshExpiresIn: h.cfg.refreshExpirySeconds,
		TokenType:        "Bearer",
		User:             toUserResponse(user),
	}))
}

// @Summary     Register
// @Description Register a new user account
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     dto.RegisterRequest true "Registration details"
// @Success     201  {object} dto.Response{data=dto.UserResponse}
// @Failure     400  {object} dto.Response
// @Failure     409  {object} dto.Response
// @Router      /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	user, err := h.auth.Register(c.Request.Context(), service.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.OK(toUserResponse(user)))
}

// @Summary     Refresh token
// @Description Exchange a refresh token for a new token pair
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     dto.RefreshTokenRequest true "Refresh token"
// @Success     200  {object} dto.Response{data=dto.TokenResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /auth/refresh [post]
func (h *authHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	pair, err := h.auth.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.OK(dto.TokenResponse{
		AccessToken:      pair.AccessToken,
		AccessExpiresIn:  pair.AccessExpiresIn,
		RefreshToken:     pair.RefreshToken,
		RefreshExpiresIn: pair.RefreshExpiresIn,
		TokenType:        "Bearer",
	}))
}

// @Summary     Logout
// @Description Invalidate the current refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.LogoutRequest true "Refresh token to invalidate"
// @Success     200  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /auth/logout [post]
func (h *authHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	if err := h.auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.OKMessage("logged out"))
}

// @Summary     Request password reset
// @Description Send a password reset email
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     dto.PasswordResetRequest true "Email address"
// @Success     200  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Router      /auth/password-reset [post]
func (h *authHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	// Always returns 200 — prevents email enumeration
	_ = h.auth.RequestPasswordReset(c.Request.Context(), req.Email)
	c.JSON(http.StatusOK, dto.OKMessage("if the email exists a reset link has been sent"))
}

// @Summary     Confirm password reset
// @Description Reset password using token from email
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     dto.PasswordResetConfirmRequest true "Token and new password"
// @Success     200  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Router      /auth/password-reset/confirm [post]
func (h *authHandler) ConfirmPasswordReset(c *gin.Context) {
	var req dto.PasswordResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	if err := h.auth.ConfirmPasswordReset(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.OKMessage("password updated"))
}

// @Summary     Change password
// @Description Change password for the authenticated user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body     dto.ChangePasswordRequest true "Old and new password"
// @Success     200  {object} dto.Response
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Router      /auth/password-change [put]
func (h *authHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.auth.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.OKMessage("password changed"))
}

// ── helpers ───────────────────────────────────────────────────────────────────

func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserToResponse(u)
}
