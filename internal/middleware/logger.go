package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		requestID, _ := c.Get("request_id")
		reqID, _ := requestID.(string)

		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.Duration("latency", latency),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("request_id", reqID),
		}

		msg := "request"
		switch {
		case status >= 500:
			slog.LogAttrs(c.Request.Context(), slog.LevelError, msg, attrs...)
		case status >= 400:
			slog.LogAttrs(c.Request.Context(), slog.LevelWarn, msg, attrs...)
		default:
			slog.LogAttrs(c.Request.Context(), slog.LevelInfo, msg, attrs...)
		}
	}
}
