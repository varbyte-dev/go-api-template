package router

import (
	"time"

	"go-api-template/internal/config"
	"go-api-template/internal/handlers"
	"go-api-template/internal/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func Setup(r *gin.Engine) {
	// ── Global middlewares (every request goes through these) ────────────────
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// Hard ceiling for the whole server: 500 req/s, burst of 1000.
	// Protects against sudden traffic spikes regardless of origin IP.
	if config.App.RateLimitEnabled {
		r.Use(middleware.GlobalRateLimit(rate.Limit(500), 1_000))
	}

	// ── Public ───────────────────────────────────────────────────────────────
	r.GET("/health", handlers.HealthCheck)

	v1 := r.Group("/api/v1")

	// Auth routes — stricter per-IP limit to slow brute-force / enumeration.
	// Allows 10 requests per minute per IP (one every 6 s), burst of 5.
	auth := v1.Group("/auth")
	if config.App.RateLimitEnabled {
		auth.Use(middleware.RateLimit(rate.Every(6*time.Second), 5))
	}
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/refresh", handlers.Refresh)
	auth.POST("/logout", handlers.Logout)

	// ── Protected (requires valid JWT) ──────────────────────────────────────
	// Generous per-IP limit for authenticated users: 60 req/s, burst of 120.
	protected := v1.Group("", middleware.AuthRequired())
	if config.App.RateLimitEnabled {
		protected.Use(middleware.RateLimit(rate.Limit(60), 120))
	}
	protected.GET("/users/me", handlers.Me)
}
