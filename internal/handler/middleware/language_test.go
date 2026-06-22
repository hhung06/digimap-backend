package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLangCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name       string
		headerVal  string
		wantLang   string
	}{
		{"header present lowercase", "ja", "ja"},
		{"header present uppercase", "JA", "ja"},
		{"header present mixed case", "En", "en"},
		{"header absent", "", "en"},
		{"header blank whitespace", "  ", "en"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got string

			r := gin.New()
			r.Use(LangCode())
			r.GET("/", func(c *gin.Context) {
				got = GetLang(c)
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.headerVal != "" {
				req.Header.Set(HeaderLangCode, tc.headerVal)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if got != tc.wantLang {
				t.Errorf("GetLang() = %q, want %q", got, tc.wantLang)
			}
		})
	}
}
