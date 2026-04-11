package handlers

import (
	"go-api-template/internal/services"
	"go-api-template/internal/utils"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"id": user.ID, "name": user.Name, "email": user.Email})
}
