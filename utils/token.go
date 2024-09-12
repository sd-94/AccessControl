package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(secret string, email string) (string, error) {
	claims := jwt.MapClaims{
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"expired": "false",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(secret string, signedToken string) (string, error) {
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		username := claims["email"].(string)
		return username, nil
	}
	return "", errors.New("invalid token")
}
