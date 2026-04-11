package handlers

import (
	"go-api-template/internal/services"
	"go-api-template/internal/utils"

	"github.com/gin-gonic/gin"
)

// Register godoc
//
//	@Summary		Registrar usuario
//	@Description	Crea una nueva cuenta de usuario con nombre, email y contraseña
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		services.RegisterInput				true	"Datos de registro"
//	@Success		201		{object}	utils.Response{data=RegisterData}	"Usuario registrado"
//	@Failure		400		{object}	utils.Response						"Datos inválidos"
//	@Failure		409		{object}	utils.Response						"Email ya registrado"
//	@Failure		429		{object}	utils.Response						"Rate limit excedido"
//	@Router			/api/v1/auth/register [post]
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

// Login godoc
//
//	@Summary		Iniciar sesión
//	@Description	Autentica al usuario y devuelve un par de tokens JWT (access + refresh)
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		services.LoginInput						true	"Credenciales"
//	@Success		200		{object}	utils.Response{data=services.TokenPair}	"Tokens generados"
//	@Failure		400		{object}	utils.Response							"Datos inválidos"
//	@Failure		401		{object}	utils.Response							"Credenciales incorrectas"
//	@Failure		429		{object}	utils.Response							"Rate limit excedido"
//	@Router			/api/v1/auth/login [post]
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

// Refresh godoc
//
//	@Summary		Renovar tokens
//	@Description	Rota el refresh token y devuelve un nuevo par de tokens JWT
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		TokenRefreshRequest						true	"Refresh token"
//	@Success		200		{object}	utils.Response{data=services.TokenPair}	"Tokens renovados"
//	@Failure		400		{object}	utils.Response							"Datos inválidos"
//	@Failure		401		{object}	utils.Response							"Token inválido o expirado"
//	@Failure		429		{object}	utils.Response							"Rate limit excedido"
//	@Router			/api/v1/auth/refresh [post]
func Refresh(c *gin.Context) {
	var body TokenRefreshRequest
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

// Logout godoc
//
//	@Summary		Cerrar sesión
//	@Description	Revoca el refresh token para invalidar la sesión
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LogoutRequest	true	"Refresh token a revocar"
//	@Success		200		{object}	utils.Response	"Sesión cerrada"
//	@Failure		400		{object}	utils.Response	"Datos inválidos"
//	@Failure		500		{object}	utils.Response	"Error interno"
//	@Failure		429		{object}	utils.Response	"Rate limit excedido"
//	@Router			/api/v1/auth/logout [post]
func Logout(c *gin.Context) {
	var body LogoutRequest
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
