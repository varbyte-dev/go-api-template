package router

import (
	"time"

	"go-api-template/internal/config"
	"go-api-template/internal/handlers"
	"go-api-template/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	_ "go-api-template/docs" // swagger docs — generado con: make swagger
)

func Setup(r *gin.Engine) {
	// ── Global middlewares ───────────────────────────────────────────────────
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// Hard ceiling para todo el servidor: 500 req/s, burst de 1000.
	if config.App.RateLimitEnabled {
		r.Use(middleware.GlobalRateLimit(rate.Limit(500), 1_000))
	}

	// ── Documentación Swagger UI ─────────────────────────────────────────────
	// Accesible en: GET /docs/index.html
	// Desactivar en producción con SWAGGER_ENABLED=false
	if config.App.SwaggerEnabled {
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// ── Público ──────────────────────────────────────────────────────────────
	r.GET("/health", handlers.HealthCheck)

	v1 := r.Group("/api/v1")

	// Auth — límite estricto por IP: 10 req/min, burst de 5
	auth := v1.Group("/auth")
	if config.App.RateLimitEnabled {
		auth.Use(middleware.RateLimit(rate.Every(6*time.Second), 5))
	}
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/refresh", handlers.Refresh)
	auth.POST("/logout", handlers.Logout)

	// ── Protegido (requiere JWT válido) ──────────────────────────────────────
	// Límite generoso por IP: 60 req/s, burst de 120
	protected := v1.Group("", middleware.AuthRequired())
	if config.App.RateLimitEnabled {
		protected.Use(middleware.RateLimit(rate.Limit(60), 120))
	}
	protected.GET("/users/me", handlers.Me)
}
