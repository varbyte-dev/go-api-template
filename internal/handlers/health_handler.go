package handlers

import (
	"go-api-template/internal/config"
	"go-api-template/internal/utils"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Devuelve el estado del servidor y el entorno de ejecución
//	@Tags			sistema
//	@Produce		json
//	@Success		200	{object}	utils.Response{data=HealthData}	"Servidor operativo"
//	@Router			/health [get]
func HealthCheck(c *gin.Context) {
	utils.OK(c, gin.H{
		"status": "ok",
		"env":    config.App.AppEnv,
	})
}
