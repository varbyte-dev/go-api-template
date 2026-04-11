package handlers

import (
	"go-api-template/internal/services"
	"go-api-template/internal/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	user, err := services.Register(input)
	if err != nil {
		utils.Conflict(c, err.Error())
		return
	}
	utils.Created(c, gin.H{"user_id": user.ID})
}

func Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	tokens, err := services.Login(input)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}
	utils.OK(c, tokens)
}

func Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	tokens, err := services.RefreshTokens(body.RefreshToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}
	utils.OK(c, tokens)
}

func Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := services.Logout(body.RefreshToken); err != nil {
		utils.InternalError(c, "could not logout")
		return
	}
	utils.OK(c, gin.H{"message": "logged out"})
}
