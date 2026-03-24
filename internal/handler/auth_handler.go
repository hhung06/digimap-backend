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

// POST /api/v1/auth/login
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

// POST /api/v1/auth/register
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

// POST /api/v1/auth/refresh
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
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		TokenType:    "Bearer",
	}))
}

// POST /api/v1/auth/logout  [protected]
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

	c.JSON(http.StatusOK, dto.OK(gin.H{"message": "logged out"}))
}

// POST /api/v1/auth/password-reset
func (h *authHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	// Always returns 200 — prevents email enumeration
	_ = h.auth.RequestPasswordReset(c.Request.Context(), req.Email)
	c.JSON(http.StatusOK, dto.OK(gin.H{"message": "if the email exists a reset link has been sent"}))
}

// POST /api/v1/auth/password-reset/confirm
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

	c.JSON(http.StatusOK, dto.OK(gin.H{"message": "password updated"}))
}

// PUT /api/v1/auth/password-change  [protected]
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

	c.JSON(http.StatusOK, dto.OK(gin.H{"message": "password changed"}))
}

// ── helpers ───────────────────────────────────────────────────────────────────

func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserToResponse(u)
}
