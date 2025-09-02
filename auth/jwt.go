package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	jwtKey = []byte(secret)
}

type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string, userID int64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
