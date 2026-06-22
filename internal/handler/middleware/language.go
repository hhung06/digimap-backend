package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	HeaderLangCode  = "X-Lang-Code"
	ContextKeyLang  = "lang"
	DefaultLangCode = "en"
)

// LangCode reads the X-Lang-Code request header, normalises it to lowercase,
// and stores it in the Gin context under ContextKeyLang. Falls back to "en"
// when the header is absent or blank.
func LangCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.ToLower(strings.TrimSpace(c.GetHeader(HeaderLangCode)))
		if code == "" {
			code = DefaultLangCode
		}
		c.Set(ContextKeyLang, code)
		c.Next()
	}
}

// GetLang returns the language code stored by the LangCode middleware.
// Returns "en" if the value is absent or not a string.
func GetLang(c *gin.Context) string {
	if v, exists := c.Get(ContextKeyLang); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return DefaultLangCode
}
