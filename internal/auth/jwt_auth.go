package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJwt(userID uuid.UUID, tokenSecret string, expiresIN time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "kurllz",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIN)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil || !token.Valid {
		return uuid.UUID{}, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.UUID{}, errors.New("Invalid claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, errors.New("Invalid user ID")
	}
	return userID, nil
}
