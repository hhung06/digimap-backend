package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/service"
)

type authServiceStub struct {
	claims *service.Claims
	err    error
}

func (s *authServiceStub) Login(context.Context, string, string) (*domain.User, service.TokenPair, error) {
	return nil, service.TokenPair{}, nil
}

func (s *authServiceStub) Register(context.Context, service.RegisterRequest) (*domain.User, error) {
	return nil, nil
}

func (s *authServiceStub) RefreshToken(context.Context, string) (service.TokenPair, error) {
	return service.TokenPair{}, nil
}

func (s *authServiceStub) Logout(context.Context, string) error {
	return nil
}

func (s *authServiceStub) RequestPasswordReset(context.Context, string) error {
	return nil
}

func (s *authServiceStub) ConfirmPasswordReset(context.Context, string, string) error {
	return nil
}

func (s *authServiceStub) ChangePassword(context.Context, uuid.UUID, string, string) error {
	return nil
}

func (s *authServiceStub) ValidateClaims(string) (*service.Claims, error) {
	return s.claims, s.err
}

func TestAdminJWTChain_AllowsSystemAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/profile",
		AuthRequired(&authServiceStub{
			claims: &service.Claims{
				UserID:        uuid.New().String(),
				Email:         "admin@example.com",
				IsSystemAdmin: true,
			},
		}),
		SystemAdminRequired(),
		func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}
