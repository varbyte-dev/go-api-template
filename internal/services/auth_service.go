package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"go-api-template/internal/config"
	"go-api-template/internal/database"
	"go-api-template/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- DTOs ---

type RegisterInput struct {
	Name     string `json:"name"     binding:"required,min=2"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// --- Claims ---

type AccessClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// --- Service ---

func Register(input RegisterInput) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, errors.New("email already registered")
	}

	return user, nil
}

func Login(input LoginInput) (*TokenPair, error) {
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return generateTokenPair(user)
}

func RefreshTokens(rawToken string) (*TokenPair, error) {
	var rt models.RefreshToken
	err := database.DB.Where("token = ? AND revoked = false", rawToken).First(&rt).Error
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// Rotate: revoke old token
	database.DB.Model(&rt).Update("revoked", true)

	var user models.User
	if err := database.DB.First(&user, rt.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return generateTokenPair(user)
}

func Logout(rawToken string) error {
	return database.DB.
		Model(&models.RefreshToken{}).
		Where("token = ?", rawToken).
		Update("revoked", true).Error
}

func ValidateAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.App.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

// --- Helpers ---

func generateTokenPair(user models.User) (*TokenPair, error) {
	accessToken, err := createAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := createRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func createAccessToken(user models.User) (string, error) {
	claims := AccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.App.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   string(rune(user.ID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.App.JWTSecret))
}

func createRefreshToken(user models.User) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(raw)

	rt := models.RefreshToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(config.App.RefreshTokenExpiry),
	}

	if err := database.DB.Create(&rt).Error; err != nil {
		return "", err
	}

	return tokenStr, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
