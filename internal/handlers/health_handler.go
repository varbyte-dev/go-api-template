package handlers

import (
	"go-api-template/internal/config"
	"go-api-template/internal/utils"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	utils.OK(c, gin.H{
		"status": "ok",
		"env":    config.App.AppEnv,
	})
}
