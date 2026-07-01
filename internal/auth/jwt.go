package auth

import (
	"errors"
	"time"
	"final-project/internal/config" 
	"github.com/golang-jwt/jwt/v5"
)


type Claims struct {
	CustomerID int `json:"customer_id"`
	jwt.RegisteredClaims
}


func getSecret() []byte {
	return []byte(config.AppConfig.JWTSecret)
}


func GenerateToken(customerID int) (string, error) {
	claims := Claims{
		CustomerID: customerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}


func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
