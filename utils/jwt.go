package utils

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(publicID string, alias string) (string, error) {
	secret := config.AppConfig.JWTSecret
	duration, _ := time.ParseDuration(config.AppConfig.JWTExpire)

	claims := jwt.MapClaims{
		"public_id": publicID,
		"alias":     alias,
		"exp":       time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(publicID string) (string, error) {
	secret := config.AppConfig.JWTSecret
	duration, _ := time.ParseDuration(config.AppConfig.JWTResfreshToken)

	claims := jwt.MapClaims{
		"public_id": publicID,
		"exp":       time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
