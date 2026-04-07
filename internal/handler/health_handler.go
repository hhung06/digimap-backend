package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/hhung06/digimap-backend/internal/dto"
)

const healthCheckTimeout = 3 * time.Second

// serviceStatus is the health state of a single dependency.
type serviceStatus struct {
	Status  string `json:"status"`           // "ok" or "error"
	Latency string `json:"latency"`          // round-trip time as a human string
	Error   string `json:"error,omitempty"`  // only present when status == "error"
}

// healthResponse is the full /health payload.
type healthResponse struct {
	Status   string                   `json:"status"`   // "ok" only when ALL checks pass
	Services map[string]serviceStatus `json:"services"`
}

// newHealthHandler returns a Gin handler that checks all critical dependencies.
// Returns HTTP 200 only when every dependency is healthy; 503 otherwise.
// Traefik's backend health check and Docker HEALTHCHECK both rely on this.
func newHealthHandler(db *pgxpool.Pool, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), healthCheckTimeout)
		defer cancel()

		services := make(map[string]serviceStatus, 2)
		allOK := true

		// ── PostgreSQL ─────────────────────────────────────────────────────────
		services["postgres"] = checkPostgres(ctx, db)
		if services["postgres"].Status != "ok" {
			allOK = false
		}

		// ── Redis ──────────────────────────────────────────────────────────────
		services["redis"] = checkRedis(ctx, cache)
		if services["redis"].Status != "ok" {
			allOK = false
		}

		overall := "ok"
		if !allOK {
			overall = "degraded"
		}

		resp := healthResponse{Status: overall, Services: services}

		// 503 signals Traefik and Docker Swarm to pull this replica from rotation.
		// The detail payload is always included so operators can see which service failed.
		status := http.StatusOK
		if !allOK {
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, dto.OK(resp))
	}
}

func checkPostgres(ctx context.Context, db *pgxpool.Pool) serviceStatus {
	start := time.Now()
	if err := db.Ping(ctx); err != nil {
		return serviceStatus{
			Status:  "error",
			Latency: time.Since(start).String(),
			Error:   err.Error(),
		}
	}
	return serviceStatus{Status: "ok", Latency: time.Since(start).String()}
}

func checkRedis(ctx context.Context, cache *redis.Client) serviceStatus {
	start := time.Now()
	if err := cache.Ping(ctx).Err(); err != nil {
		return serviceStatus{
			Status:  "error",
			Latency: time.Since(start).String(),
			Error:   err.Error(),
		}
	}
	return serviceStatus{Status: "ok", Latency: time.Since(start).String()}
}
