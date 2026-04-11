// Package main is the entry point for go-api-template.
//
//	@title			go-api-template
//	@version		1.0
//	@description	Boilerplate de API REST en Go con Gin, GORM, JWT y rate limiting.
//
//	@contact.name	Soporte API
//	@contact.email	support@example.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Escribe "Bearer " seguido de tu access token JWT.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-template/internal/config"
	"go-api-template/internal/database"
	"go-api-template/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	if config.App.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	database.Connect()

	r := gin.New()
	router.Setup(r)

	addr := ":" + config.App.AppPort
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		slog.Info("Server starting", "addr", addr, "env", config.App.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "err", err)
		os.Exit(1)
	}

	slog.Info("Server exited")
}
