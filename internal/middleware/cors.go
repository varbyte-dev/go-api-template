package middleware

import (
	"time"

	"go-api-template/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	origins := config.App.CORSOrigins

	allowAll := len(origins) == 0
	if !allowAll {
		for _, o := range origins {
			if o == "*" {
				allowAll = true
				break
			}
		}
	}

	cfg := cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders: []string{"X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining", "Retry-After"},
		MaxAge:        12 * time.Hour,
	}

	if allowAll {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = origins
		cfg.AllowCredentials = true
	}

	return cors.New(cfg)
}
