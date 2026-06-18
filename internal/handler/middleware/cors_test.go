package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hhung06/digimap-backend/config"
)

func TestCORS_AllowsFrontendRequestHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORS(config.AppConfig{CORSAllowedOrigins: []string{"*"}}))
	router.GET("/api/v1/profile", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/profile", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "authorization,cache-control,content-type,x-requested-with")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	allowedHeaders := strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers"))
	for _, header := range []string{"cache-control", "x-requested-with"} {
		if !strings.Contains(allowedHeaders, header) {
			t.Errorf("Access-Control-Allow-Headers = %q, want %s", allowedHeaders, header)
		}
	}
}
