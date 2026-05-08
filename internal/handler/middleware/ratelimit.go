package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	limiterredis "github.com/ulule/limiter/v3/drivers/store/redis"

	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"

	"github.com/hhung06/digimap-backend/internal/dto"
)

// RateLimitByIP returns a Gin middleware that rate-limits requests per client IP.
// The rate string follows the ulule/limiter format: "<count>-<period>" e.g. "60-M" (60/min).
// Counters are stored in Redis so limits apply across all replicas.
func RateLimitByIP(client *redis.Client, rate string) gin.HandlerFunc {
	if client == nil {
		return func(c *gin.Context) { c.Next() }
	}

	r, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		panic("ratelimit: invalid rate format: " + rate)
	}

	store, err := limiterredis.NewStore(client)
	if err != nil {
		panic("ratelimit: could not create Redis store: " + err.Error())
	}

	lim := limiter.New(store, r)

	middleware := ginlimiter.NewMiddleware(lim,
		ginlimiter.WithKeyGetter(func(c *gin.Context) string {
			return clientIPKey(c)
		}),
		ginlimiter.WithLimitReachedHandler(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				dto.Fail(dto.CodeRateLimited, "too many requests, please slow down"))
		}),
	)

	return middleware
}

// RateLimitByUser returns a Gin middleware that rate-limits per authenticated user UUID.
// Falls back to IP-based keying for unauthenticated requests.
// Must be placed after AuthRequired in the middleware chain.
func RateLimitByUser(client *redis.Client, rate string) gin.HandlerFunc {
	if client == nil {
		return func(c *gin.Context) { c.Next() }
	}

	r, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		panic("ratelimit: invalid rate format: " + rate)
	}

	store, err := limiterredis.NewStore(client)
	if err != nil {
		panic("ratelimit: could not create Redis store: " + err.Error())
	}

	lim := limiter.New(store, r)

	middleware := ginlimiter.NewMiddleware(lim,
		ginlimiter.WithKeyGetter(func(c *gin.Context) string {
			userID := GetUserID(c)
			if userID.String() != "00000000-0000-0000-0000-000000000000" {
				return "user:" + userID.String()
			}
			return clientIPKey(c)
		}),
		ginlimiter.WithLimitReachedHandler(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				dto.Fail(dto.CodeRateLimited, "too many requests, please slow down"))
		}),
	)

	return middleware
}

// clientIPKey returns the real client IP for use as a rate-limit key.
func clientIPKey(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		for i, ch := range xff {
			if ch == ',' {
				return "ip:" + xff[:i]
			}
		}
		return "ip:" + xff
	}
	return "ip:" + c.ClientIP()
}
