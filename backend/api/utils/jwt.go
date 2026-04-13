package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GenerateJWT(userID string) (string, error) {
	// Implement JWT generation logic here, e.g., using the "github.com/dgrijalva/jwt-go" package
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
	})
	return token.SignedString([]byte(viper.GetString("JWT_SECRET")))

}

func ValidateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	// Implement JWT validation logic here
	token, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
