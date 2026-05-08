package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/service"
	applog "github.com/hhung06/digimap-backend/log"
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

func newRouterTestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Environment:        "local",
			LogLevel:           "error",
			LogFormat:          "text",
			CORSAllowedOrigins: []string{"*"},
		},
	}
}

func newRouterTestLogger() applog.Logger {
	return applog.NewLogger(newRouterTestConfig())
}

func TestNewRouter_DisablesLocationDuplicateRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(newRouterTestConfig(), newRouterTestLogger(), Dependencies{
		AuthService: &authServiceStub{},
	})

	for _, route := range router.Routes() {
		assert.False(t,
			route.Method == http.MethodPost && route.Path == "/api/v1/venues/:id/locations/:locationID/duplicate",
			"duplicate route should not be mounted",
		)
	}
}

func TestNewRouter_ProfileRejectsNonSystemAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(newRouterTestConfig(), newRouterTestLogger(), Dependencies{
		AuthService: &authServiceStub{
			claims: &service.Claims{
				UserID:        uuid.New().String(),
				Email:         "user@example.com",
				IsSystemAdmin: false,
			},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "system admin access required")
}

func TestNewRouter_MountsAppSurveySubmitRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(newRouterTestConfig(), newRouterTestLogger(), Dependencies{
		AuthService: &authServiceStub{},
	})

	found := false
	for _, route := range router.Routes() {
		if route.Method == http.MethodPost && route.Path == "/app/v1/surveys/:surveyID/submit-response" {
			found = true
			break
		}
	}

	assert.True(t, found, "app survey submit route should be mounted")
}

func TestNewRouter_AppSurveySubmitRequiresAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(newRouterTestConfig(), newRouterTestLogger(), Dependencies{
		AuthService: &authServiceStub{},
	})

	req := httptest.NewRequest(http.MethodPost, "/app/v1/surveys/"+uuid.New().String()+"/submit-response", strings.NewReader(`{"external_id":"visitor-1","answers":[]}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "API key required")
}

func TestNewRouter_PublicSurveySubmitDoesNotRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(newRouterTestConfig(), newRouterTestLogger(), Dependencies{
		AuthService: &authServiceStub{},
	})

	req := httptest.NewRequest(http.MethodPost, "/public/v1/surveys/not-a-uuid/submit-response", strings.NewReader(`{"answers":[]}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.NotContains(t, rec.Body.String(), "API key required")
	assert.NotContains(t, rec.Body.String(), "Authorization")
}

func TestNewRouter_MountsVisitorSurveyPostRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(newRouterTestConfig(), newRouterTestLogger(), Dependencies{
		AuthService: &authServiceStub{},
	})

	found := false
	for _, route := range router.Routes() {
		if route.Method == http.MethodPost && route.Path == "/public/v1/visitor-surveys" {
			found = true
			break
		}
	}

	assert.True(t, found, "visitor survey POST route should be mounted")
}

func TestNewRouter_DoesNotMountVisitorSurveyGetRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(newRouterTestConfig(), newRouterTestLogger(), Dependencies{
		AuthService: &authServiceStub{},
	})

	for _, route := range router.Routes() {
		assert.False(t,
			route.Method == http.MethodGet && route.Path == "/public/v1/visitor-surveys",
			"visitor survey GET route should not be mounted",
		)
	}
}
