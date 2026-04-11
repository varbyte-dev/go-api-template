package handlers

import (
	"go-api-template/internal/services"
	"go-api-template/internal/utils"

	"github.com/gin-gonic/gin"
)

// Me godoc
//
//	@Summary		Perfil del usuario autenticado
//	@Description	Devuelve los datos del usuario asociado al token JWT
//	@Tags			usuarios
//	@Produce		json
//	@Success		200	{object}	utils.Response{data=UserData}	"Datos del usuario"
//	@Failure		401	{object}	utils.Response					"No autenticado"
//	@Failure		404	{object}	utils.Response					"Usuario no encontrado"
//	@Failure		429	{object}	utils.Response					"Rate limit excedido"
//	@Security		BearerAuth
//	@Router			/api/v1/users/me [get]
func Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"id": user.ID, "name": user.Name, "email": user.Email})
}
